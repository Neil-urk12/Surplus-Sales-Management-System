// Package main implements the web server for the Cortes Surplus Inventory Management System.
package main

import (
	"context"
	"fmt"
	"log"
	"oop/internal/config"
	"oop/internal/handlers"
	"oop/internal/middleware"
	"oop/internal/repositories"
	"os"
	"os/signal"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

var jwtSecret = []byte(getEnv("JWT_SECRET", "7b324dbe6535a315ef300b1c79bc35324c8c2fb1495195176364d8e12379f22c469bea496c1336c85724175d83d566d2271df2f49a14a09d5fddecc85962435a17196950ad693461e38c6645148a15726626c35e0ead273673a8f4a98547d015c89dfdf0d6bc5332ad9cbc1a180363d881ab320ef0f825c8cc83286aea871562f442c71c05e44298b1e83f43e1e7a57a101718bdd58489c05978317afd1feaa7fa2d2898f6bd41e09823172c87d676b025c488eb849e71adc2408cbe36ca9814bd38091dd14b6d790c68e4210350f7bdd365fbfa903fe59a9744f70021943b2c4f65101695ad1c1d25468bb4589fefcaa863fdb57b8321b2a5cbfa4cb6ac275a"))

// Helper function to get environment variable or default value
func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	log.Printf("Warning: Environment variable %s not set. Using default.", key)
	return fallback
}

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
	app := fiber.New(fiber.Config{
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			// Default error handling
			code := fiber.StatusInternalServerError
			if e, ok := err.(*fiber.Error); ok {
				code = e.Code
			}
			log.Printf("Error: %v", err) // Log the error
			return c.Status(code).JSON(fiber.Map{
				"error": err.Error(),
			})
		},
	})

	// Add middleware
	app.Use(logger.New())
	app.Use(recover.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",                           // Consider restricting this in production
		AllowMethods: "GET,POST,PUT,DELETE,OPTIONS", // Added OPTIONS for preflight
		AllowHeaders: "Origin, Content-Type, Accept, Authorization",
	}))

	// Initialize repositories
	userRepo := repositories.NewUserRepository(dbClient)

	// Initialize handlers
	userHandler := handlers.NewUserHandler(userRepo)

	// --- Route Registration ---
	api := app.Group("/api") // Base group for API routes

	// Public User Routes (register, login)
	userHandler.RegisterRoutes(api) // This will now only register public routes

	// Protected User Routes (require JWT)
	authMiddleware := middleware.JWTMiddleware(jwtSecret)
	userProtected := api.Group("/users", authMiddleware) // Apply middleware here

	userProtected.Get("/", userHandler.GetAllUsers)
	userProtected.Get("/:id", userHandler.GetUser)
	userProtected.Put("/:id", userHandler.UpdateUser)
	userProtected.Delete("/:id", userHandler.DeleteUser)
	userProtected.Put("/:id/activate", userHandler.ActivateUser)
	userProtected.Put("/:id/deactivate", userHandler.DeactivateUser)
	userProtected.Put("/:id/password", userHandler.UpdatePassword)

	// Add a health check endpoint (public)
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"status": "ok",
			"time":   time.Now().Format(time.RFC3339),
		})
	})

	// Start server in a goroutine so we can listen for shutdown signal
	go func() {
		port := getEnv("PORT", "8080")
		log.Printf("Starting server on :%s", port)
		if err := app.Listen(":" + port); err != nil {
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
