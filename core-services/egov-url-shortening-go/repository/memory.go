package repository

import (
	"context"
	"fmt"
	"sync"
	"sync/atomic"

	"egov-url-shortening-go/models"
	"github.com/sirupsen/logrus"
)

// MemoryRepository implements URLRepository using in-memory storage for testing
type MemoryRepository struct {
	data    map[string]*models.ShortenRequest
	counter int64
	mu      sync.RWMutex
	log     *logrus.Logger
}

// NewMemoryRepository creates a new in-memory repository instance
func NewMemoryRepository(log *logrus.Logger) *MemoryRepository {
	log.Info("Using in-memory repository for testing")
	return &MemoryRepository{
		data:    make(map[string]*models.ShortenRequest),
		counter: 0,
		log:     log,
	}
}

// IncrementID increments and returns the next available ID
func (m *MemoryRepository) IncrementID(ctx context.Context) (int64, error) {
	id := atomic.AddInt64(&m.counter, 1) - 1
	m.log.WithField("id", id).Debug("Incremented ID")
	return id, nil
}

// SaveURL saves a URL with the given key and request
func (m *MemoryRepository) SaveURL(ctx context.Context, key string, request *models.ShortenRequest) error {
	if key == "" {
		return fmt.Errorf("key cannot be empty")
	}
	if request == nil {
		return fmt.Errorf("request cannot be nil")
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	// Make a copy to avoid data races
	requestCopy := *request
	m.data[key] = &requestCopy

	m.log.WithFields(logrus.Fields{
		"url": request.URL,
		"key": key,
	}).Info("URL saved successfully in memory")

	return nil
}

// GetURL retrieves a URL by ID
func (m *MemoryRepository) GetURL(ctx context.Context, id int64) (string, error) {
	request, err := m.GetURLDetails(ctx, id)
	if err != nil {
		return "", err
	}

	// Check if URL is active
	if !request.IsActive() {
		if request.IsExpired() {
			m.log.WithField("id", id).Warn("URL has expired")
			return "", fmt.Errorf("URL at key %d has expired", id)
		}
		m.log.WithField("id", id).Warn("URL is not yet active")
		return "", fmt.Errorf("URL at key %d is not yet active", id)
	}

	return request.URL, nil
}

// GetURLDetails retrieves full URL details by ID
func (m *MemoryRepository) GetURLDetails(ctx context.Context, id int64) (*models.ShortenRequest, error) {
	key := fmt.Sprintf("url:%d", id)
	
	m.mu.RLock()
	defer m.mu.RUnlock()

	request, exists := m.data[key]
	if !exists {
		m.log.WithField("id", id).Debug("URL not found")
		return nil, fmt.Errorf("URL at key %d does not exist", id)
	}

	m.log.WithFields(logrus.Fields{
		"url": request.URL,
		"id":  id,
	}).Debug("Retrieved URL details from memory")

	// Return a copy to avoid data races
	requestCopy := *request
	return &requestCopy, nil
}

// DeleteURL deletes a URL by ID
func (m *MemoryRepository) DeleteURL(ctx context.Context, id int64) error {
	key := fmt.Sprintf("url:%d", id)
	
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, exists := m.data[key]; !exists {
		return fmt.Errorf("URL at key %d does not exist", id)
	}

	delete(m.data, key)
	m.log.WithField("id", id).Info("URL deleted successfully from memory")
	return nil
}

// CheckURLExists checks if a URL exists for the given ID
func (m *MemoryRepository) CheckURLExists(ctx context.Context, id int64) (bool, error) {
	key := fmt.Sprintf("url:%d", id)
	
	m.mu.RLock()
	defer m.mu.RUnlock()

	_, exists := m.data[key]
	return exists, nil
}

// HealthCheck performs a health check on the repository
func (m *MemoryRepository) HealthCheck(ctx context.Context) error {
	// Memory repository is always healthy
	return nil
}

// Close closes the repository (no-op for memory)
func (m *MemoryRepository) Close() error {
	m.mu.Lock()
	defer m.mu.Unlock()
	
	// Clear data
	m.data = make(map[string]*models.ShortenRequest)
	m.log.Info("Memory repository closed and data cleared")
	return nil
}

// GetStats returns memory repository statistics
func (m *MemoryRepository) GetStats() map[string]interface{} {
	m.mu.RLock()
	defer m.mu.RUnlock()

	return map[string]interface{}{
		"type":  "memory",
		"count": len(m.data),
		"urls":  m.data,
	}
}