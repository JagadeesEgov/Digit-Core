package repository

import (
	"context"
	"encoding/json"
	"fmt"

	"egov-url-shortening-go/config"
	"egov-url-shortening-go/models"

	"github.com/go-redis/redis/v8"
	"github.com/sirupsen/logrus"
)

// RedisRepository implements URLRepository using Redis
type RedisRepository struct {
	client *redis.Client
	idKey  string
	urlKey string
	log    *logrus.Logger
}

// NewRedisRepository creates a new Redis repository instance
func NewRedisRepository(cfg *config.RedisConfig, log *logrus.Logger) *RedisRepository {
	client := redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
	})

	return &RedisRepository{
		client: client,
		idKey:  "id",
		urlKey: "url:",
		log:    log,
	}
}

// IncrementID increments and returns the next available ID
func (r *RedisRepository) IncrementID() (int64, error) {
	ctx := context.Background()
	id, err := r.client.Incr(ctx, r.idKey).Result()
	if err != nil {
		r.log.WithError(err).Error("Failed to increment ID")
		return 0, err
	}

	// Return id-1 to match Java implementation (Java uses pre-increment logic)
	result := id - 1
	r.log.WithField("id", result).Info("Incrementing ID")
	return result, nil
}

// SaveURL saves a URL with the given key and request
func (r *RedisRepository) SaveURL(key string, request *models.ShortenRequest) error {
	ctx := context.Background()
	
	// Serialize the request to JSON (matching Java ObjectMapper behavior)
	data, err := json.Marshal(request)
	if err != nil {
		r.log.WithError(err).WithField("key", key).Error("Failed to marshal request")
		return err
	}

	// Use HSET to store in a hash (matching Java Jedis.hset behavior)
	err = r.client.HSet(ctx, r.urlKey, key, string(data)).Err()
	if err != nil {
		r.log.WithError(err).WithField("key", key).Error("Failed to save URL")
		return err
	}

	r.log.WithFields(logrus.Fields{
		"url": request.URL,
		"key": key,
	}).Info("Saving URL")

	return nil
}

// GetURL retrieves a URL by ID
func (r *RedisRepository) GetURL(id int64) (string, error) {
	ctx := context.Background()
	key := fmt.Sprintf("url:%d", id)
	
	r.log.WithField("id", id).Info("Retrieving URL")
	
	// Get from hash
	data, err := r.client.HGet(ctx, r.urlKey, key).Result()
	if err != nil {
		if err == redis.Nil {
			r.log.WithField("id", id).Warn("URL not found")
			return "", fmt.Errorf("URL at key %d does not exist", id)
		}
		r.log.WithError(err).WithField("id", id).Error("Failed to retrieve URL")
		return "", err
	}

	// Deserialize the JSON data
	var request models.ShortenRequest
	err = json.Unmarshal([]byte(data), &request)
	if err != nil {
		r.log.WithError(err).WithField("id", id).Error("Failed to unmarshal URL data")
		return "", err
	}

	r.log.WithFields(logrus.Fields{
		"url": request.URL,
		"id":  id,
	}).Info("Retrieved URL")

	return request.URL, nil
}

// Close closes the Redis connection
func (r *RedisRepository) Close() error {
	return r.client.Close()
}