// Package main implements the web server for the Cortes Surplus Inventory Management System.
package main

import (
	"context"
	"fmt"
	"log"
	"oop/internal/config"
	"oop/internal/repositories"
	"os"
	"os/signal"
	"time"

	"github.com/gofiber/fiber/v2"
)

func main() {
	// Load db config
	dbClient, err := initDatabase()
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	// Create a shutdown channel
	shutdown := make(chan struct{})

	// Handle graceful shutdown on interrupt signal
	go handleShutdown(dbClient, shutdown)

	// Initialize Fiber app
	app := fiber.New()

	// Start server in a goroutine so we can listen for shutdown signal
	go func() {
		if err := app.Listen(":8080"); err != nil {
			log.Printf("Server error: %v", err)
			close(shutdown) // Signal shutdown if server fails
		}
	}()

	// Wait for shutdown signal
	<-shutdown
	log.Println("Shutting down gracefully...")

	// Shutdown the server with a timeout
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := app.ShutdownWithContext(ctx); err != nil {
		log.Printf("Error during server shutdown: %v", err)
	}

	log.Println("Server shutdown complete")
}

// initDatabase loads the database configuration, connects to the database,
// It returns a DatabaseClient pointer and an error if initialization fails.
func initDatabase() (*repositories.DatabaseClient, error) {
	dbConfig, err := config.LoadDatabaseConfig()
	if err != nil {
		return nil, fmt.Errorf("load database config: %v", err)
	}

	dbClient, err := repositories.NewDatabaseClient(dbConfig)
	if err != nil {
		return nil, fmt.Errorf("connect to database: %v", err)
	}

	log.Println("Database connection test successful.")
	return dbClient, nil
}

// handleShutdown listens for interrupt signals (like Ctrl+C) to gracefully shut down the application.
// It closes the database connection and signals the main goroutine to shut down.
func handleShutdown(dbClient *repositories.DatabaseClient, shutdown chan struct{}) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c
	log.Println("Interrupt received, initiating shutdown...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second) // Adjust timeout as needed
	defer cancel()

	if err := dbClient.Close(ctx); err != nil {
		log.Printf("Error closing database connection: %v", err)
	}

	// Signal the main goroutine to shut down
	close(shutdown)
}
