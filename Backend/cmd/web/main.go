// Package main implements the web server for the Cortes Surplus Inventory Management System.
package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"time"

	"oop/internal/config"
	"oop/internal/handlers"
	"oop/internal/middleware"
	"oop/internal/repositories"

	_ "oop/docs" // load API docs generated by Swag CLI

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	swagger "github.com/gofiber/swagger" // swagger handler
	"github.com/joho/godotenv"
)

var jwtSecret = []byte(getEnv("JWT_SECRET", "7b324dbe6535a315ef300b1c79bc35324c8c2fb1495195176364d8e12379f22c469bea496c1336c85724175d83d566d2271df2f49a14a09d5fddecc85962435a17196950ad693461e38c6645148a15726626c35e0ead273673a8f4a98547d015c89dfdf0d6bc5332ad9cbc1a180363d881ab320ef0f825c8cc83286aea871562f442c71c05e44298b1e83f43e1e7a57a101718bdd58489c05978317afd1feaa7fa2d2898f6bd41e09823172c87d676b025c488eb849e71adc2408cbe36ca9814bd38091dd14b6d790c68e4210350f7bdd365fbfa903fe59a9744f70021943b2c4f65101695ad1c1d25468bb4589fefcaa863fdb57b8321b2a5cbfa4cb6ac275a"))

// @title Cortes Surplus Inventory Management API
// @version 1.0
// @description This is the API for the Cortes Surplus Inventory Management System.
// @termsOfService http://swagger.io/terms/
// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io
// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html
// @host localhost:8080
// @BasePath /api
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization

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

// turnstileMiddleware creates a middleware that verifies Cloudflare Turnstile tokens
func turnstileMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
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

		log.Printf("Success: Turnstile verification passed for token: %s", truncatedToken)
		return c.Next()
	}
}

func main() {
	// Load env variables
	err := godotenv.Load()

	if err != nil {
		log.Println("Failed to load environment variables... Retrying...")
		return
	}

	if os.Getenv("FRONTEND_URL") == "" {
		log.Fatalf("FRONTEND_URL is not set. Please set it to a valid frontend URL (e.g., http://localhost:9000).")
	}

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
	app.Use(cors.New(cors.Config{
		AllowOrigins:     os.Getenv("FRONTEND_URL"),     // Restricted to specific origins from environment variable
		AllowMethods:     "GET,POST,PUT,DELETE,OPTIONS", // Added OPTIONS for preflight
		AllowCredentials: true,
		AllowHeaders:     "Origin, Content-Type, Accept, Authorization",
	}))

	// Initialize repositories
	userRepo := repositories.NewUserRepository(dbClient)
	materialRepo := repositories.NewMaterialRepository(dbClient.DB)
	accessoryRepo := repositories.NewAccessoryRepository(dbClient.DB)
	customerRepo := repositories.NewCustomerRepository(dbClient.DB)

	// Initialize cabs repository directly with DB
	cabsRepo := repositories.NewCabsRepository(dbClient.DB)

	// Initialize sales repository
	saleRepo := repositories.NewSalesRepository(dbClient.DB)

	// Initialize logs repository
	logsRepo := repositories.NewLogsRepository(dbClient.DB) // Assuming dbClient.DB is the *sql.DB instance

	// Initialize handlers
	userHandler := handlers.NewUserHandler(userRepo, jwtSecret)
	materialHandler := handlers.NewMaterialHandlers(materialRepo, jwtSecret)
	customerHandler := handlers.NewCustomerHandler(customerRepo, jwtSecret)
	// Initialize cabs handler
	cabsHandler := handlers.NewCabsHandlers(cabsRepo)
	accessoryHandler := handlers.NewAccessoriesHandler(accessoryRepo)
	// Initialize sales handler
	saleHandler := handlers.NewSaleHandlers(saleRepo, cabsRepo, accessoryRepo, customerRepo, jwtSecret)
	// Initialize activity log handler
	activityLogHandler := handlers.NewActivityLogHandler(logsRepo)

	// --- Route Registration ---
	api := app.Group("/api") // Base group for API routes

	// Swagger docs route
	api.Get("/swagger/*", swagger.HandlerDefault) // get /api/swagger/*

	// @Summary Submit Turnstile Captcha
	// @Description Verifies a Cloudflare Turnstile token.
	// @Tags Captcha
	// @Accept x-www-form-urlencoded
	// @Produce json
	// @Param cf-turnstile-response formData string true "Cloudflare Turnstile Token"
	// @Success 200 {object} fiber.Map{"status=ok"}
	// @Failure 400 {object} fiber.Map{"error=captcha token missing"}
	// @Failure 403 {object} fiber.Map{"error=invalid captcha"}
	// @Failure 500 {object} fiber.Map{"error=verification failed"}
	// @Router /submit [post]
	app.Post("/submit", turnstileMiddleware(), func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"status": "ok"})
	})

	// Public User Routes (register, login)
	userHandler.RegisterRoutes(api) // This will now only register public routes
	materialHandler.RegisterMaterialRoutes(api)
	customerHandler.RegisterCustomerRoutes(api)

	// Register Cabs routes - Detailed Swagger annotations are in cabs_handlers.go
	api.Get("/cabs", cabsHandler.GetCabs)          // GET /api/cabs
	api.Get("/cabs/:id", cabsHandler.GetCabByID)   // GET /api/cabs/:id
	api.Post("/cabs", cabsHandler.AddCab)          // POST /api/cabs
	api.Put("/cabs/:id", cabsHandler.UpdateCab)    // PUT /api/cabs/:id
	api.Delete("/cabs/:id", cabsHandler.DeleteCab) // DELETE /api/cabs/:id

	// Register Accessories routes - Detailed Swagger annotations are in accessories_handlers.go
	api.Get("/accessories", accessoryHandler.GetAllAccessories)      // GET /api/accessories
	api.Get("/accessories/:id", accessoryHandler.GetAccessoryByID)   // GET /api/accessories/:id
	api.Post("/accessories", accessoryHandler.CreateAccessory)       // POST /api/accessories
	api.Put("/accessories/:id", accessoryHandler.UpdateAccessory)    // PUT /api/accessories/:id
	api.Delete("/accessories/:id", accessoryHandler.DeleteAccessory) // DELETE /api/accessories/:id

	// Register Sale routes - Detailed Swagger annotations are in sales_handlers.go
	saleHandler.RegisterSaleRoutes(api)

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

	// Protected Activity Log Routes (require JWT)
	activityLogProtected := api.Group("/activity-logs", authMiddleware)
	activityLogProtected.Get("/", activityLogHandler.GetActivityLogs)
	activityLogProtected.Get("/filter", activityLogHandler.GetFilteredActivityLogs)
	activityLogProtected.Post("/", activityLogHandler.CreateActivityLog)

	// Add a health check endpoint (public)
	// @Summary Health Check
	// @Description Checks if the server is running
	// @Tags Health
	// @Accept json
	// @Produce json
	// @Success 200 {object} fiber.Map
	// @Router /health [get]
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
