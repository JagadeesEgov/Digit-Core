package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"egov-url-shortening-go/config"
	"egov-url-shortening-go/models"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sirupsen/logrus"
)

// PostgresRepository implements URLRepository using PostgreSQL
type PostgresRepository struct {
	db  *pgxpool.Pool
	log *logrus.Logger
}

// NewPostgresRepository creates a new PostgreSQL repository instance
func NewPostgresRepository(cfg *config.DatabaseConfig, log *logrus.Logger) (*PostgresRepository, error) {
	connString := fmt.Sprintf("postgres://%s:%s@%s:%d/%s",
		cfg.Username, cfg.Password, cfg.Host, cfg.Port, cfg.Name)

	db, err := pgxpool.New(context.Background(), connString)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Test the connection
	err = db.Ping(context.Background())
	if err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	log.Info("Connected to PostgreSQL database")

	return &PostgresRepository{
		db:  db,
		log: log,
	}, nil
}

// IncrementID increments and returns the next available ID using the sequence
func (p *PostgresRepository) IncrementID() (int64, error) {
	ctx := context.Background()
	
	var id int64
	err := p.db.QueryRow(ctx, "SELECT nextval('eg_url_shorter_id')").Scan(&id)
	if err != nil {
		p.log.WithError(err).Error("Failed to get next sequence value")
		return 0, err
	}

	// Return id-1 to match Java implementation
	result := id - 1
	p.log.WithField("id", result).Info("Incrementing ID")
	return result, nil
}

// SaveURL saves a URL with the given key and request
func (p *PostgresRepository) SaveURL(key string, request *models.ShortenRequest) error {
	ctx := context.Background()
	
	// Extract ID from key (format: "url:123")
	var id string
	if len(key) > 4 && key[:4] == "url:" {
		id = key[4:]
	} else {
		return fmt.Errorf("invalid key format: %s", key)
	}

	// Convert Unix timestamps to time.Time if provided
	var validFrom, validTo *time.Time
	if request.ValidFrom != nil {
		t := time.Unix(*request.ValidFrom/1000, 0) // Assuming milliseconds
		validFrom = &t
	}
	if request.ValidTill != nil {
		t := time.Unix(*request.ValidTill/1000, 0) // Assuming milliseconds
		validTo = &t
	}

	// Insert or update the URL entry
	query := `
		INSERT INTO eg_url_shortener (id, url, validform, validto) 
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (id) 
		DO UPDATE SET url = $2, validform = $3, validto = $4`

	var validFromMs, validToMs *int64
	if validFrom != nil {
		ms := validFrom.UnixMilli()
		validFromMs = &ms
	}
	if validTo != nil {
		ms := validTo.UnixMilli()
		validToMs = &ms
	}

	_, err := p.db.Exec(ctx, query, id, request.URL, validFromMs, validToMs)
	if err != nil {
		p.log.WithError(err).WithField("key", key).Error("Failed to save URL")
		return err
	}

	p.log.WithFields(logrus.Fields{
		"url": request.URL,
		"key": key,
	}).Info("Saving URL")

	return nil
}

// GetURL retrieves a URL by ID
func (p *PostgresRepository) GetURL(id int64) (string, error) {
	ctx := context.Background()
	
	p.log.WithField("id", id).Info("Retrieving URL")

	var url string
	var validFrom, validTo sql.NullInt64
	
	query := "SELECT url, validform, validto FROM eg_url_shortener WHERE id = $1"
	err := p.db.QueryRow(ctx, query, fmt.Sprintf("%d", id)).Scan(&url, &validFrom, &validTo)
	if err != nil {
		if err == pgx.ErrNoRows {
			p.log.WithField("id", id).Warn("URL not found")
			return "", fmt.Errorf("URL at key %d does not exist", id)
		}
		p.log.WithError(err).WithField("id", id).Error("Failed to retrieve URL")
		return "", err
	}

	// Check if URL is still valid (if expiry is set)
	if validTo.Valid {
		expiryTime := time.Unix(validTo.Int64/1000, 0)
		if time.Now().After(expiryTime) {
			p.log.WithField("id", id).Warn("URL has expired")
			return "", fmt.Errorf("URL at key %d has expired", id)
		}
	}

	p.log.WithFields(logrus.Fields{
		"url": url,
		"id":  id,
	}).Info("Retrieved URL")

	return url, nil
}

// Close closes the database connection pool
func (p *PostgresRepository) Close() {
	p.db.Close()
}