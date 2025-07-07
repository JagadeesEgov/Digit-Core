package models

import "time"

// ShortenRequest represents the request payload for URL shortening
type ShortenRequest struct {
	ID        string `json:"id,omitempty"`
	URL       string `json:"url" validate:"required"`
	ValidFrom *int64 `json:"validFrom,omitempty"`
	ValidTill *int64 `json:"validTill,omitempty"`
}

// ShortenResponse represents the response for URL shortening
type ShortenResponse struct {
	ShortenedURL string `json:"shortenedUrl"`
}

// URLEntry represents a URL entry in the database
type URLEntry struct {
	ID        string     `json:"id" db:"id"`
	URL       string     `json:"url" db:"url"`
	ValidFrom *time.Time `json:"validFrom" db:"validform"`
	ValidTo   *time.Time `json:"validTo" db:"validto"`
	CreatedAt time.Time  `json:"createdAt" db:"created_at"`
	UpdatedAt time.Time  `json:"updatedAt" db:"updated_at"`
}

// ErrorResponse represents an error response
type ErrorResponse struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

// UserSearchRequest represents a user search request
type UserSearchRequest struct {
	Type     string `json:"type"`
	TenantID string `json:"tenantId"`
	UserName string `json:"userName"`
}

// UserSearchResponse represents a user search response
type UserSearchResponse struct {
	User []UserInfo `json:"user"`
}

// UserInfo represents user information
type UserInfo struct {
	UUID string `json:"uuid"`
}