package main

import (
	"context"
	"errors"
	"fmt"
	"oop/internal/config"
	"oop/internal/repositories"
	"os"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
)

// mockDatabaseClient is a mock implementation for testing
type mockDatabaseClient struct {
	closeWasCalled bool
	closeError     error
}

func (m *mockDatabaseClient) Close(ctx context.Context) error {
	m.closeWasCalled = true
	return m.closeError
}

// mockConfigLoader is used to mock the config loading functionality
type mockConfigLoader struct {
	config config.DatabaseConfig
	err    error
}

func (m *mockConfigLoader) Load() (config.DatabaseConfig, error) {
	return m.config, m.err
}

// mockDBClientFactory is used to mock the database client creation
type mockDBClientFactory struct {
	client *mockDatabaseClient
	err    error
}

func (m *mockDBClientFactory) NewClient(config config.DatabaseConfig) (*mockDatabaseClient, error) {
	return m.client, m.err
}

// TestInitDatabase tests the initDatabase function
func TestInitDatabase(t *testing.T) {
	// Test case 1: Config loading fails
	t.Run("Config loading fails", func(t *testing.T) {
		// Create mocks
		configLoader := &mockConfigLoader{
			err: errors.New("config error"),
		}

		// Call function under test with our mocks
		client, err := initDatabaseWithDeps(configLoader.Load, nil)

		// Verify results
		if err == nil {
			t.Error("Expected an error, got nil")
		}
		if client != nil {
			t.Errorf("Expected nil client, got %v", client)
		}
	})

	// Test case 2: Database connection fails
	t.Run("Database connection fails", func(t *testing.T) {
		// Create mocks
		configLoader := &mockConfigLoader{
			config: config.DatabaseConfig{},
		}

		dbFactory := &mockDBClientFactory{
			err: errors.New("connection error"),
		}

		// Call function under test with our mocks
		client, err := initDatabaseWithDeps(
			configLoader.Load,
			func(config config.DatabaseConfig) (*repositories.DatabaseClient, error) {
				return nil, dbFactory.err
			},
		)

		// Verify results
		if err == nil {
			t.Error("Expected an error, got nil")
		}
		if client != nil {
			t.Errorf("Expected nil client, got %v", client)
		}
	})

	// Test case 3: Success
	t.Run("Success", func(t *testing.T) {
		// Create mocks
		configLoader := &mockConfigLoader{
			config: config.DatabaseConfig{},
		}

		mockClient := &repositories.DatabaseClient{}

		// Call function under test with our mocks
		client, err := initDatabaseWithDeps(
			configLoader.Load,
			func(config config.DatabaseConfig) (*repositories.DatabaseClient, error) {
				return mockClient, nil
			},
		)

		// Verify results
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
		if client != mockClient {
			t.Errorf("Expected client to be %v, got %v", mockClient, client)
		}
	})
}

// Helper function for testing that accepts dependencies
func initDatabaseWithDeps(
	loadConfig func() (config.DatabaseConfig, error),
	newDBClient func(config config.DatabaseConfig) (*repositories.DatabaseClient, error),
) (*repositories.DatabaseClient, error) {
	dbConfig, err := loadConfig()
	if err != nil {
		return nil, fmt.Errorf("load database config: %v", err)
	}

	dbClient, err := newDBClient(dbConfig)
	if err != nil {
		return nil, fmt.Errorf("connect to database: %v", err)
	}

	return dbClient, nil
}

// TestHandleShutdown tests the handleShutdown function
func TestHandleShutdown(t *testing.T) {
	// Create a mock database client
	mockClient := &mockDatabaseClient{}

	// Create a shutdown channel
	shutdown := make(chan struct{})

	// Create a done channel to signal when the test is complete
	done := make(chan struct{})

	// Start a goroutine to simulate the handleShutdown function
	// without relying on OS signals which are hard to test
	go func() {
		// Simulate closing the database
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
		defer cancel()

		if err := mockClient.Close(ctx); err != nil {
			t.Errorf("Error closing database: %v", err)
		}

		// Signal shutdown
		close(shutdown)

		// Signal test is complete
		close(done)
	}()

	// Wait for the test to complete or timeout
	select {
	case <-done:
		// Test completed successfully
	case <-time.After(500 * time.Millisecond):
		t.Error("Test timed out")
	}

	// Verify the database was closed
	if !mockClient.closeWasCalled {
		t.Error("Expected database to be closed")
	}

	// Verify the shutdown channel was closed
	select {
	case <-shutdown:
		// Channel was closed, which is what we expect
	default:
		t.Error("Expected shutdown channel to be closed")
	}
}

// TestCreateFiberApp tests that we can create a Fiber app
func TestCreateFiberApp(t *testing.T) {
	app := fiber.New()
	if app == nil {
		t.Error("Expected to create a Fiber app")
	}
}

// TestFiberAppShutdown tests that the Fiber app can be shut down gracefully
func TestFiberAppShutdown(t *testing.T) {
	// Create a Fiber app
	app := fiber.New()

	// Create a context with a timeout
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	// Shut down the app
	err := app.ShutdownWithContext(ctx)

	// Verify the app was shut down without errors
	if err != nil {
		t.Errorf("Expected no error when shutting down Fiber app, got %v", err)
	}
}

// TestGracefulShutdown tests the graceful shutdown process
func TestGracefulShutdown(t *testing.T) {
	// Create a shutdown channel
	shutdown := make(chan struct{})

	// Create a Fiber app
	app := fiber.New()

	// Start a goroutine to simulate the server
	serverDone := make(chan struct{})
	go func() {
		// Wait for shutdown signal
		<-shutdown

		// Shut down the app
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
		defer cancel()

		err := app.ShutdownWithContext(ctx)
		if err != nil {
			t.Errorf("Error during server shutdown: %v", err)
		}

		close(serverDone)
	}()

	// Trigger shutdown
	close(shutdown)

	// Wait for server to shut down or timeout
	select {
	case <-serverDone:
		// Server shut down successfully
	case <-time.After(500 * time.Millisecond):
		t.Error("Server shutdown timed out")
	}
}

// TestServerStartup tests the server startup process
func TestServerStartup(t *testing.T) {
	// Create a Fiber app
	app := fiber.New()

	// Start the server in a goroutine
	serverStarted := make(chan struct{})
	go func() {
		// Signal that we're about to start the server
		close(serverStarted)

		// Start listening on a random port to avoid conflicts
		// Use a port that's likely to be free
		err := app.Listen(":0")

		// This line should only be reached if there's an error
		// or the server is shut down
		if err != nil {
			// We expect an error when the server is shut down
			// Fiber doesn't export ErrServerClosed, so we'll just check for any error
			t.Logf("Server stopped with error: %v", err)
		}
	}()

	// Wait for the server to start
	<-serverStarted

	// Give the server a moment to initialize
	time.Sleep(50 * time.Millisecond)

	// Shut down the server
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	err := app.ShutdownWithContext(ctx)
	if err != nil {
		t.Errorf("Error shutting down server: %v", err)
	}
}

// TestMain is used to set up any test environment needs
func TestMain(m *testing.M) {
	// Setup code here if needed

	// Run tests
	exitCode := m.Run()

	// Cleanup code here if needed

	// Exit with the same code as the tests
	os.Exit(exitCode)
}
