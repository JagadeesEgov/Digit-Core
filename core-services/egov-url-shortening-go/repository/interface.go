package repository

import (
	"egov-url-shortening-go/models"
)

// URLRepository defines the interface for URL storage operations
type URLRepository interface {
	// IncrementID returns the next available ID
	IncrementID() (int64, error)
	
	// SaveURL saves a URL with the given key and request
	SaveURL(key string, request *models.ShortenRequest) error
	
	// GetURL retrieves a URL by ID
	GetURL(id int64) (string, error)
}