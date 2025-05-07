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
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

var jwtSecret = []byte(getEnv("JWT_SECRET", "7b324dbe6535a315ef300b1c79bc35324c8c2fb1495195176364d8e12379f22c469bea496c1336c85724175d83d566d2271df2f49a14a09d5fddecc85962435a17196950ad693461e38c6645148a15726626c35e0ead273673a8f4a98547d015c89dfdf0d6bc5332ad9cbc1a180363d881ab320ef0f825c8cc83286aea871562f442c71c05e44298b1e83f43e1e7a57a101718bdd58489c05978317afd1feaa7fa2d2898f6bd41e09823172c87d676b025c488eb849e71adc2408cbe36ca9814bd38091dd14b6d790c68e4210350f7bdd365fbfa903fe59a9744f70021943b2c4f65101695ad1c1d25468bb4589fefcaa863fdb57b8321b2a5cbfa4cb6ac275a"))

// Helper function to get environment variable or default value with validation and sanitization
// Supported environment variables:
// - JWT_SECRET: Secret key for JWT token generation and validation
// - PORT: Port on which the server will listen
// - ALLOWED_ORIGINS: Comma-separated list of allowed origins for CORS (e.g., "http://localhost:9000,https://example.com")
func getEnv(key, fallback string) string {
	// Validate key parameter
	if key == "" {
		log.Printf("Error: Empty environment variable key provided")
		return fallback
	}

	// Get environment variable value
	value, ok := os.LookupEnv(key)
	if !ok {
		// TODO: Convert to log.Fatalf in production
		log.Printf("Warning: Environment variable %s not set. Using default.", key)
		return fallback
	}

	// Validate and sanitize value based on the environment variable type
	switch key {
	case "JWT_SECRET":
		// JWT_SECRET should be at least 32 characters long for security
		if len(value) < 32 {
			log.Printf("Warning: JWT_SECRET is too short (< 32 chars). Using default for security.")
			return fallback
		}
		// No need to sanitize JWT_SECRET as it's used as-is for cryptographic purposes

	case "PORT":
		// Validate port number (should be numeric and within valid range 1-65535)
		port := 0
		_, err := fmt.Sscanf(value, "%d", &port)
		if err != nil || port < 1 || port > 65535 {
			log.Printf("Warning: Invalid PORT value '%s'. Must be a number between 1-65535. Using default.", value)
			return fallback
		}
		// Return the validated port as a string
		return fmt.Sprintf("%d", port)

	case "ALLOWED_ORIGINS":
		// Validate and sanitize CORS origins
		if value == "*" {
			log.Printf("Warning: ALLOWED_ORIGINS set to '*' which allows all origins. This is not recommended for production.")
		} else if value == "" {
			log.Printf("Warning: ALLOWED_ORIGINS is empty. Using default.")
			return fallback
		}
		// Basic validation of origins format (could be enhanced further)
		// This simple check ensures each origin starts with http:// or https://
		origins := strings.Split(value, ",")
		validOrigins := make([]string, 0, len(origins))
		for _, origin := range origins {
			origin = strings.TrimSpace(origin)
			if origin == "" {
				continue
			}
			if !strings.HasPrefix(origin, "http://") && !strings.HasPrefix(origin, "https://") {
				log.Printf("Warning: Origin '%s' doesn't start with http:// or https://. Skipping.", origin)
				continue
			}
			validOrigins = append(validOrigins, origin)
		}
		if len(validOrigins) == 0 {
			log.Printf("Warning: No valid origins found in ALLOWED_ORIGINS. Using default.")
			return fallback
		}
		return strings.Join(validOrigins, ",")
	}

	// For other environment variables, perform basic sanitization
	// Remove leading/trailing whitespace
	value = strings.TrimSpace(value)
	if value == "" {
		log.Printf("Warning: Environment variable %s is empty after sanitization. Using default.", key)
		return fallback
	}

	return value
}

func main() {
	// Load db config
	dbClient, err := initDatabase()
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	// Load and validate Turnstile config
	if _, err := config.LoadTurnstileConfig(); err != nil {
		log.Fatalf("Failed to load Turnstile configuration: %v", err)
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

	// Get allowed origins from environment variable or use default for development
	allowedOrigins := getEnv("ALLOWED_ORIGINS", "http://localhost:9000,http://localhost:8080")

	app.Use(cors.New(cors.Config{
		AllowOrigins: allowedOrigins,                // Restricted to specific origins from environment variable
		AllowMethods: "GET,POST,PUT,DELETE,OPTIONS", // Added OPTIONS for preflight
		AllowHeaders: "Origin, Content-Type, Accept, Authorization",
	}))

	// Initialize repositories
	userRepo := repositories.NewUserRepository(dbClient)
	materialRepo := repositories.NewMaterialRepository(dbClient.DB)

	// Initialize cabs repository directly with DB
	cabsRepo := repositories.NewCabsRepository(dbClient.DB)

	// Initialize handlers
	userHandler := handlers.NewUserHandler(userRepo, jwtSecret)
	materialHandler := handlers.NewMaterialHandlers(materialRepo, jwtSecret)
	// Initialize cabs handler
	cabsHandler := handlers.NewCabsHandlers(cabsRepo)

	// --- Route Registration ---
	api := app.Group("/api") // Base group for API routes
	app.Post("/submit", func(c *fiber.Ctx) error {
		// 1) grab the Turnstile token from the client
		token := c.FormValue("cf-turnstile-response")
		if token == "" {
			log.Printf("Error: Missing Turnstile token in request")
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "captcha token missing",
			})
		}

		// Log token for debugging (truncate for security)
		tokenLength := len(token)
		truncatedToken := ""
		if tokenLength > 10 {
			truncatedToken = token[:5] + "..." + token[tokenLength-5:]
		} else {
			truncatedToken = token
		}
		log.Printf("Debug: Processing Turnstile token: %s", truncatedToken)

		// 2) verify with Cloudflare
		ok, err := handlers.VerifyTurnstile(token)
		if err != nil {
			log.Printf("Error: Turnstile verification failed: %v, token: %s", err, truncatedToken)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		if !ok {
			log.Printf("Error: Invalid Turnstile captcha with token: %s", truncatedToken)
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error": "invalid captcha",
			})
		}

		// 3) CAPTCHA passed — now invoke your normal form logic.
		//    e.g., parse the rest of the body and save to DB:
		//
		//    var payload YourPayloadType
		//    if err := c.BodyParser(&payload); err != nil { … }
		//    // do your business logic here…
		//
		log.Printf("Success: Turnstile verification passed for token: %s", truncatedToken)
		return c.JSON(fiber.Map{"status": "ok"})
	})

	// Public User Routes (register, login)
	userHandler.RegisterRoutes(api) // This will now only register public routes
	materialHandler.RegisterMaterialRoutes(api)

	// Register Cabs routes
	api.Get("/cabs", cabsHandler.GetCabs)          // GET /api/cabs
	api.Get("/cabs/:id", cabsHandler.GetCabByID)   // GET /api/cabs/:id
	api.Post("/cabs", cabsHandler.AddCab)          // POST /api/cabs
	api.Put("/cabs/:id", cabsHandler.UpdateCab)    // PUT /api/cabs/:id
	api.Delete("/cabs/:id", cabsHandler.DeleteCab) // DELETE /api/cabs/:id

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
	userProtected.Post("/", userHandler.CreateUser)

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
