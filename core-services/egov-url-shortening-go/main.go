package main

import (
	"fmt"
	"log"

	"egov-url-shortening-go/config"
	"egov-url-shortening-go/handlers"
	"egov-url-shortening-go/repository"
	"egov-url-shortening-go/service"
	"egov-url-shortening-go/utils"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func main() {
	// Initialize logger
	logger := logrus.New()
	logger.SetFormatter(&logrus.JSONFormatter{})
	logger.Info("Starting eGov URL Shortening Service")

	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Initialize HashID converter
	hashConverter, err := utils.NewHashIDConverter(&cfg.HashIDs)
	if err != nil {
		log.Fatalf("Failed to initialize HashID converter: %v", err)
	}

	// Initialize URL validator
	urlValidator := utils.NewURLValidator()

	// Initialize repository based on configuration
	var urlRepo repository.URLRepository
	if cfg.Database.Enabled {
		// Use PostgreSQL
		urlRepo, err = repository.NewPostgresRepository(&cfg.Database, logger)
		if err != nil {
			log.Fatalf("Failed to initialize PostgreSQL repository: %v", err)
		}
		logger.Info("Using PostgreSQL repository")
	} else {
		// Use Redis
		urlRepo = repository.NewRedisRepository(&cfg.Redis, logger)
		logger.Info("Using Redis repository")
	}

	// Initialize service
	urlService := service.NewURLService(urlRepo, hashConverter, urlValidator, cfg, logger)

	// Initialize handlers
	handler := handlers.NewHandler(urlService, cfg, logger)

	// Setup Gin router
	router := gin.New()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	// Setup routes
	setupRoutes(router, handler, cfg)

	// Start server
	port := fmt.Sprintf(":%d", cfg.Server.Port)
	logger.WithField("port", cfg.Server.Port).Info("Starting HTTP server")
	
	if err := router.Run(port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

// setupRoutes configures the HTTP routes
func setupRoutes(router *gin.Engine, handler *handlers.Handler, cfg *config.Config) {
	// Health check endpoint
	router.GET("/health", handler.HealthCheck)

	// API routes with context path
	api := router.Group(cfg.Server.ContextPath)
	{
		// URL shortening endpoint
		api.POST("/shortener", handler.ShortenURL)
		
		// URL redirection endpoint
		api.GET("/:id", handler.RedirectURL)
	}

	// Also add routes without context path for direct access
	router.POST("/shortener", handler.ShortenURL)
	router.GET("/:id", handler.RedirectURL)
}