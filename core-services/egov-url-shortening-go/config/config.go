package config

import (
	"encoding/json"
	"log"

	"github.com/kelseyhightower/envconfig"
)

// Config holds all configuration for the URL shortening service
type Config struct {
	Server   ServerConfig   `envconfig:"SERVER"`
	Redis    RedisConfig    `envconfig:"REDIS"`
	Database DatabaseConfig `envconfig:"DATABASE"`
	Kafka    KafkaConfig    `envconfig:"KAFKA"`
	App      AppConfig      `envconfig:"APP"`
	HashIDs  HashIDsConfig  `envconfig:"HASHIDS"`
}

// ServerConfig holds server configuration
type ServerConfig struct {
	Port        int    `envconfig:"PORT" default:"8091"`
	ContextPath string `envconfig:"CONTEXT_PATH" default:"/egov-url-shortening"`
}

// RedisConfig holds Redis configuration
type RedisConfig struct {
	Host string `envconfig:"HOST" default:"localhost"`
	Port int    `envconfig:"PORT" default:"6379"`
}

// DatabaseConfig holds database configuration
type DatabaseConfig struct {
	Host     string `envconfig:"HOST" default:"localhost"`
	Port     int    `envconfig:"PORT" default:"5432"`
	Name     string `envconfig:"NAME" default:"devdb"`
	Username string `envconfig:"USERNAME" default:"postgres"`
	Password string `envconfig:"PASSWORD" default:"postgres"`
	Enabled  bool   `envconfig:"ENABLED" default:"true"`
}

// KafkaConfig holds Kafka configuration
type KafkaConfig struct {
	BootstrapServers string `envconfig:"BOOTSTRAP_SERVERS" default:"localhost:9092"`
	Topic            string `envconfig:"TOPIC" default:"save-url-shortening-details"`
}

// AppConfig holds application-specific configuration
type AppConfig struct {
	HostName                 string            `envconfig:"HOST_NAME" default:"https://qa.digit.org/"`
	StateLevelTenantID       string            `envconfig:"STATE_LEVEL_TENANT_ID" default:"pb"`
	UserHost                 string            `envconfig:"USER_HOST" default:"http://egov-user.egov:8080/"`
	UserSearchPath           string            `envconfig:"USER_SEARCH_PATH" default:"user/_search"`
	UIAppHostMap             string            `envconfig:"UI_APP_HOST_MAP" default:"{\"in\":\"https://central-instance.digit.org\",\"in.statea\":\"https://statea.digit.org\"}"`
	UIAppHostMapParsed       map[string]string
	IsMultiInstance          bool `envconfig:"IS_MULTI_INSTANCE" default:"false"`
	IsCentralInstance        bool `envconfig:"IS_CENTRAL_INSTANCE" default:"true"`
	StateLevelTenantIDLength int  `envconfig:"STATE_LEVEL_TENANT_ID_LENGTH" default:"2"`
}

// HashIDsConfig holds HashIDs configuration
type HashIDsConfig struct {
	Salt      string `envconfig:"SALT" default:"randomsalt"`
	MinLength int    `envconfig:"MIN_LENGTH" default:"3"`
}

// LoadConfig loads configuration from environment variables
func LoadConfig() (*Config, error) {
	var cfg Config
	err := envconfig.Process("", &cfg)
	if err != nil {
		return nil, err
	}

	// Parse the UI app host map JSON
	if cfg.App.UIAppHostMap != "" {
		err = json.Unmarshal([]byte(cfg.App.UIAppHostMap), &cfg.App.UIAppHostMapParsed)
		if err != nil {
			log.Printf("Warning: Failed to parse UI app host map: %v", err)
			cfg.App.UIAppHostMapParsed = make(map[string]string)
		}
	}

	return &cfg, nil
}