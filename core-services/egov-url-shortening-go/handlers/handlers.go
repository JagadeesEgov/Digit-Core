package handlers

import (
	"net/http"
	"strings"

	"egov-url-shortening-go/config"
	"egov-url-shortening-go/models"
	"egov-url-shortening-go/service"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// Handler holds the HTTP handlers for the URL shortening service
type Handler struct {
	urlService *service.URLService
	config     *config.Config
	log        *logrus.Logger
}

// NewHandler creates a new handler instance
func NewHandler(urlService *service.URLService, cfg *config.Config, log *logrus.Logger) *Handler {
	return &Handler{
		urlService: urlService,
		config:     cfg,
		log:        log,
	}
}

// ShortenURL handles URL shortening requests
func (h *Handler) ShortenURL(c *gin.Context) {
	var request models.ShortenRequest
	
	// Bind JSON request body
	if err := c.ShouldBindJSON(&request); err != nil {
		h.log.WithError(err).Error("Failed to bind request")
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Code:    "INVALID_REQUEST",
			Message: "Invalid request format",
		})
		return
	}

	// Log headers for debugging (matching Java implementation)
	headers := make(map[string]string)
	for key, values := range c.Request.Header {
		if len(values) > 0 {
			headers[strings.ToLower(key)] = values[0]
		}
	}
	h.log.WithField("headers", headers).Info("Request headers")

	// Determine tenant ID based on multi-instance configuration
	var tenantID string
	if !h.config.App.IsMultiInstance {
		tenantID = h.config.App.StateLevelTenantID
	} else {
		// Extract state-specific tenant ID from ULB level tenant
		tenantIDHeader := c.GetHeader("tenantid")
		if tenantIDHeader == "" {
			c.JSON(http.StatusBadRequest, models.ErrorResponse{
				Code:    "INVALID_TENANTID",
				Message: "TenantId not present in header",
			})
			return
		}

		// Extract state level tenant from ULB tenant
		// This is a simplified version - you might need to implement the actual logic
		// based on your tenant ID structure
		tenantID = h.extractStateLevelTenant(tenantIDHeader)
	}

	// Validate URL
	if !h.urlService.ValidateURL(request.URL) {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Code:    "URL_SHORTENING_INVALID_URL",
			Message: "Please enter a valid URL",
		})
		return
	}

	// Shorten URL
	shortenedURL, err := h.urlService.ShortenURL(&request, tenantID, h.config.App.IsMultiInstance)
	if err != nil {
		h.log.WithError(err).Error("Failed to shorten URL")
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Code:    "URL_SHORTENING_FAILED",
			Message: "Failed to shorten URL",
		})
		return
	}

	// Return shortened URL as plain text (matching Java implementation)
	c.String(http.StatusOK, shortenedURL)
}

// RedirectURL handles URL redirection requests
func (h *Handler) RedirectURL(c *gin.Context) {
	id := c.Param("id")
	
	if id == "" {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Code:    "INVALID_REQUEST",
			Message: "ID parameter is required",
		})
		return
	}

	// Get original URL
	longURL, err := h.urlService.GetLongURLFromID(id)
	if err != nil {
		h.log.WithError(err).WithField("id", id).Error("Failed to get long URL")
		c.JSON(http.StatusNotFound, models.ErrorResponse{
			Code:    "INVALID_REQUEST",
			Message: "Invalid Key",
		})
		return
	}

	// Redirect to original URL
	c.Redirect(http.StatusFound, longURL)
}

// HealthCheck provides a health check endpoint
func (h *Handler) HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status": "UP",
		"service": "egov-url-shortening-go",
	})
}

// extractStateLevelTenant extracts state level tenant from ULB level tenant
// This is a simplified implementation - adjust based on your actual tenant structure
func (h *Handler) extractStateLevelTenant(ulbTenantID string) string {
	// Example: if ULB tenant is "pb.amritsar", state tenant is "pb"
	parts := strings.Split(ulbTenantID, ".")
	if len(parts) > 0 && len(parts[0]) == h.config.App.StateLevelTenantIDLength {
		return parts[0]
	}
	return h.config.App.StateLevelTenantID // fallback
}