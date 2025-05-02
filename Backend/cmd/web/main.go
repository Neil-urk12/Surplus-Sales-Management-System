// Package main implements the web server for the Cortes Surplus Inventory Management System.
package main

import (
	"fmt"
	"log"
	"oop/internal/config"
	"oop/internal/repositories"
	"os"
	"os/signal"

	"github.com/gofiber/fiber/v2"
)

func main() {
	// Load db config
	dbClient, err := initDatabase()
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer dbClient.Close()

	// Handle graceful shutdown on interrupt signal
	go handleShutdown(dbClient)

	// Initialize Fiber app
	app := fiber.New()

	// Start server
	log.Fatal(app.Listen(":8080"))
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
// It closes the database connection and exits the program.
func handleShutdown(dbClient *repositories.DatabaseClient) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c
	log.Println("Shutting down...")
	dbClient.Close()
	os.Exit(0)
}
