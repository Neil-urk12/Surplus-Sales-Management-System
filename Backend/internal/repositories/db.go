package repositories

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"oop/internal/config"

	_ "github.com/go-sql-driver/mysql"
)

// DatabaseClient represents a client for interacting with the database.
type DatabaseClient struct {
	DB *sql.DB
}

// NewDatabaseClient creates a new DatabaseClient and establishes a database connection
// using the provided database configuration. It returns a pointer to the client and an error.
func NewDatabaseClient(config config.DatabaseConfig) (*DatabaseClient, error) {
	// Enable parsing of MySQL TIMESTAMP fields into time.Time and set charset to utf8mb4
	connStr := fmt.Sprintf(
		"%s:%s@tcp(%s:%d)/%s?parseTime=true&charset=utf8mb4&loc=Local",
		config.Username, config.Password, config.Host, config.Port, config.DatabaseName,
	)

	db, err := sql.Open("mysql", connStr)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %v", err)
	}

	return &DatabaseClient{DB: db}, nil
}

// Close closes the database connection held by the DatabaseClient.
// It accepts a context for timeout control and returns an error if closing fails.
func (c *DatabaseClient) Close(ctx context.Context) error {
	log.Println("Closing database connection...")
	if err := ctx.Err(); err != nil {
		return err
	}

	log.Println("Database connection closed successfully.")
	return c.DB.Close()
}
