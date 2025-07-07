package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"egov-url-shortening-go/config"
	"egov-url-shortening-go/models"
	"egov-url-shortening-go/repository"
	"egov-url-shortening-go/utils"

	"github.com/sirupsen/logrus"
)

// URLService provides URL shortening and conversion functionality
type URLService struct {
	repository    repository.URLRepository
	hashConverter *utils.HashIDConverter
	validator     *utils.URLValidator
	config        *config.Config
	log           *logrus.Logger
	httpClient    *http.Client
}

// NewURLService creates a new URL service instance
func NewURLService(
	repo repository.URLRepository,
	hashConverter *utils.HashIDConverter,
	validator *utils.URLValidator,
	cfg *config.Config,
	log *logrus.Logger,
) *URLService {
	return &URLService{
		repository:    repo,
		hashConverter: hashConverter,
		validator:     validator,
		config:        cfg,
		log:           log,
		httpClient:    &http.Client{},
	}
}

// ShortenURL creates a shortened URL from the given request
func (s *URLService) ShortenURL(request *models.ShortenRequest, tenantID string, multiInstance bool) (string, error) {
	s.log.WithField("url", request.URL).Info("Shortening URL")

	// Get next ID
	id, err := s.repository.IncrementID()
	if err != nil {
		return "", fmt.Errorf("failed to get next ID: %w", err)
	}

	// Create unique ID using hash
	uniqueID, err := s.hashConverter.CreateHashStringForID(id)
	if err != nil {
		return "", fmt.Errorf("failed to create hash for ID %d: %w", id, err)
	}

	// Save URL
	key := fmt.Sprintf("url:%d", id)
	err = s.repository.SaveURL(key, request)
	if err != nil {
		return "", fmt.Errorf("failed to save URL: %w", err)
	}

	// Build shortened URL
	shortenedURL, err := s.buildShortenedURL(uniqueID, tenantID, multiInstance)
	if err != nil {
		return "", fmt.Errorf("failed to build shortened URL: %w", err)
	}

	return shortenedURL, nil
}

// GetLongURLFromID retrieves the original URL from a shortened ID
func (s *URLService) GetLongURLFromID(uniqueID string) (string, error) {
	s.log.WithField("uniqueID", uniqueID).Info("Converting shortened URL back")

	// Convert hash back to ID
	id, err := s.hashConverter.GetIDForString(uniqueID)
	if err != nil || id == 0 {
		// If modern hash conversion fails, this could be a legacy ID
		// For now, just return error - can implement legacy support later if needed
		return "", fmt.Errorf("invalid shortened URL ID: %s", uniqueID)
	}

	// Get URL from repository
	longURL, err := s.repository.GetURL(id)
	if err != nil {
		return "", err
	}

	if longURL == "" {
		return "", fmt.Errorf("invalid request: URL not found")
	}

	s.log.WithField("longURL", longURL).Info("Successfully retrieved original URL")
	return longURL, nil
}

// ValidateURL validates a URL using the configured validator
func (s *URLService) ValidateURL(url string) bool {
	return s.validator.ValidateURL(url)
}

// buildShortenedURL constructs the complete shortened URL
func (s *URLService) buildShortenedURL(uniqueID, tenantID string, multiInstance bool) (string, error) {
	var hostName string

	if multiInstance {
		// Multi-instance: get hostname from tenant mapping
		if s.config.App.UIAppHostMapParsed == nil {
			return "", fmt.Errorf("UI app host map not configured for multi-instance")
		}

		var exists bool
		hostName, exists = s.config.App.UIAppHostMapParsed[tenantID]
		if !exists {
			return "", fmt.Errorf("hostname for provided state level tenant has not been configured for tenantId: %s", tenantID)
		}
	} else {
		// Single instance: use configured hostname
		hostName = s.config.App.HostName
	}

	// Clean up hostname (remove trailing slash)
	if strings.HasSuffix(hostName, "/") {
		hostName = hostName[:len(hostName)-1]
	}

	// Clean up context path (remove leading slash)
	contextPath := s.config.Server.ContextPath
	if strings.HasPrefix(contextPath, "/") {
		contextPath = contextPath[1:]
	}

	// Build the complete URL
	var shortenedURL strings.Builder
	shortenedURL.WriteString(hostName)
	shortenedURL.WriteString("/")
	shortenedURL.WriteString(contextPath)
	if !strings.HasSuffix(contextPath, "/") {
		shortenedURL.WriteString("/")
	}
	shortenedURL.WriteString(uniqueID)

	return shortenedURL.String(), nil
}

// GetUserUUID retrieves user UUID by mobile number (optional functionality)
func (s *URLService) GetUserUUID(mobileNumber string) (string, error) {
	request := models.UserSearchRequest{
		Type:     "CITIZEN",
		TenantID: s.config.App.StateLevelTenantID,
		UserName: mobileNumber,
	}

	jsonData, err := json.Marshal(request)
	if err != nil {
		return "", fmt.Errorf("failed to marshal user search request: %w", err)
	}

	url := s.config.App.UserHost + s.config.App.UserSearchPath
	resp, err := s.httpClient.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		s.log.WithError(err).Error("Exception while fetching user")
		return "", err
	}
	defer resp.Body.Close()

	var response models.UserSearchResponse
	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		s.log.WithError(err).Error("Failed to decode user search response")
		return "", err
	}

	if len(response.User) > 0 {
		return response.User[0].UUID, nil
	}

	return "", nil // User not found
}