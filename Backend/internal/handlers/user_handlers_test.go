package handlers

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"oop/internal/models"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockUserRepository is a mock implementation of the UserRepository
type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) Create(user *models.User) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *MockUserRepository) GetByID(id string) (*models.User, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserRepository) GetByEmail(email string) (*models.User, error) {
	args := m.Called(email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserRepository) Update(user *models.User) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *MockUserRepository) UpdatePassword(userID string, newPassword string) error {
	args := m.Called(userID, newPassword)
	return args.Error(0)
}

func (m *MockUserRepository) Delete(id string) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockUserRepository) GetAll() ([]*models.User, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*models.User), args.Error(1)
}

func (m *MockUserRepository) VerifyPassword(email, password string) (*models.User, error) {
	args := m.Called(email, password)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserRepository) ActivateUser(id string) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockUserRepository) DeactivateUser(id string) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockUserRepository) EmailExists(email string) (bool, error) {
	args := m.Called(email)
	return args.Bool(0), args.Error(1)
}

// Helper function to create a test app and handler
func setupTest() (*fiber.App, *UserHandler, *MockUserRepository) {
	app := fiber.New()
	mockRepo := new(MockUserRepository)
	handler := NewUserHandler(mockRepo, []byte("dummy_secret_for_test"))
	return app, handler, mockRepo
}

// Helper function to create a test user
func createTestUser() *models.User {
	return &models.User{
		Id:        uuid.New().String(),
		FullName:  "Test User",
		Email:     "test@example.com",
		Password:  "hashed_password",
		Role:      "user",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		IsActive:  true,
	}
}

func TestUserHandler_Register_Success(t *testing.T) {
	// Setup
	app, handler, mockRepo := setupTest()

	// Setup route
	app.Post("/register", handler.Register)

	// Setup expectations
	mockRepo.On("EmailExists", "test@example.com").Return(false, nil)
	mockRepo.On("Create", mock.AnythingOfType("*models.User")).Return(nil)

	// Create request body
	reqBody := map[string]string{
		"fullName": "Test User",
		"email":    "test@example.com",
		"password": "password123",
		"role":     "user",
	}
	jsonBody, _ := json.Marshal(reqBody)

	// Create request
	req := httptest.NewRequest(http.MethodPost, "/register", bytes.NewReader(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	// Perform request
	resp, err := app.Test(req)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, http.StatusCreated, resp.StatusCode)

	// Verify response body
	var result map[string]interface{}
	error := json.NewDecoder(resp.Body).Decode(&result)
	assert.NoError(t, error)

	assert.Equal(t, "User registered successfully", result["message"])
	assert.NotNil(t, result["user"])

	// Verify expectations
	mockRepo.AssertExpectations(t)
}

func TestUserHandler_Register_EmailExists(t *testing.T) {
	// Setup
	app, handler, mockRepo := setupTest()

	// Setup route
	app.Post("/register", handler.Register)

	// Setup expectations
	mockRepo.On("EmailExists", "test@example.com").Return(true, nil)

	// Create request body
	reqBody := map[string]string{
		"fullName": "Test User",
		"email":    "test@example.com",
		"password": "password123",
	}
	jsonBody, _ := json.Marshal(reqBody)

	// Create request
	req := httptest.NewRequest(http.MethodPost, "/register", bytes.NewReader(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	// Perform request
	resp, err := app.Test(req)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, http.StatusConflict, resp.StatusCode)

	// Verify response body
	var result map[string]interface{}
	error := json.NewDecoder(resp.Body).Decode(&result)
	assert.NoError(t, error)

	assert.Equal(t, "Email already in use", result["error"])

	// Verify expectations
	mockRepo.AssertExpectations(t)
}

func TestUserHandler_Login_Success(t *testing.T) {
	// Setup
	app, handler, mockRepo := setupTest()

	// Setup route
	app.Post("/login", handler.Login)

	// Create test user
	user := createTestUser()

	// Setup expectations
	mockRepo.On("GetByEmail", "test@example.com").Return(user, nil)
	mockRepo.On("VerifyPassword", "test@example.com", "password123").Return(user, nil)

	// Create request body
	reqBody := map[string]string{
		"email":    "test@example.com",
		"password": "password123",
	}
	jsonBody, _ := json.Marshal(reqBody)

	// Create request
	req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewReader(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	// Perform request
	resp, err := app.Test(req)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	// Verify response body
	var result map[string]interface{}
	error := json.NewDecoder(resp.Body).Decode(&result)
	assert.NoError(t, error)

	assert.Equal(t, "Login successful", result["message"])
	assert.NotNil(t, result["user"])

	// Verify expectations
	mockRepo.AssertExpectations(t)
}

func TestUserHandler_Login_InvalidCredentials(t *testing.T) {
	// Setup
	app, handler, mockRepo := setupTest()

	// Setup route
	app.Post("/login", handler.Login)

	// Create test user
	user := createTestUser()

	// Setup expectations
	mockRepo.On("GetByEmail", "test@example.com").Return(user, nil)
	mockRepo.On("VerifyPassword", "test@example.com", "wrong_password").Return(nil, errors.New("invalid password"))

	// Create request body
	reqBody := map[string]string{
		"email":    "test@example.com",
		"password": "wrong_password",
	}
	jsonBody, _ := json.Marshal(reqBody)

	// Create request
	req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewReader(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	// Perform request
	resp, err := app.Test(req)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)

	// Verify response body
	var result map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&result)
	assert.NoError(t, err)

	assert.Equal(t, "Invalid credentials", result["error"])

	// Verify expectations
	mockRepo.AssertExpectations(t)
}

func TestUserHandler_Login_InactiveUser(t *testing.T) {
	// Setup
	app, handler, mockRepo := setupTest()

	// Setup route
	app.Post("/login", handler.Login)

	// Create inactive test user
	inactiveUser := createTestUser()
	inactiveUser.IsActive = false

	// Setup expectations
	mockRepo.On("GetByEmail", "test@example.com").Return(inactiveUser, nil)

	// Create request body
	reqBody := map[string]string{
		"email":    "test@example.com",
		"password": "password123",
	}
	jsonBody, _ := json.Marshal(reqBody)

	// Create request
	req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewReader(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	// Perform request
	resp, err := app.Test(req)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, http.StatusForbidden, resp.StatusCode)

	// Verify response body
	var result map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&result)
	assert.NoError(t, err)

	assert.Equal(t, "Account is inactive", result["error"])

	// Verify expectations
	mockRepo.AssertExpectations(t)
}

func TestUserHandler_GetAllUsers(t *testing.T) {
	// Setup
	app, handler, mockRepo := setupTest()

	// Setup route
	app.Get("/users", handler.GetAllUsers)

	// Create test users
	now := time.Now()
	users := []*models.User{
		{
			Id:        "user1",
			FullName:  "User One",
			Email:     "user1@example.com",
			Password:  "hashed_password1",
			Role:      "user",
			CreatedAt: now,
			UpdatedAt: now,
			IsActive:  true,
		},
		{
			Id:        "user2",
			FullName:  "User Two",
			Email:     "user2@example.com",
			Password:  "hashed_password2",
			Role:      "admin",
			CreatedAt: now,
			UpdatedAt: now,
			IsActive:  true,
		},
	}

	// Setup expectations
	mockRepo.On("GetAll").Return(users, nil)

	// Create request
	req := httptest.NewRequest(http.MethodGet, "/users", nil)

	// Perform request
	resp, err := app.Test(req)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	// Verify response body
	var result map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&result)
	assert.NoError(t, err)

	assert.NotNil(t, result["users"])

	// Verify expectations
	mockRepo.AssertExpectations(t)
}

func TestUserHandler_GetUser(t *testing.T) {
	// Setup
	app, handler, mockRepo := setupTest()

	// Setup route
	app.Get("/users/:id", handler.GetUser)

	// Create test user
	user := createTestUser()

	// Setup expectations
	mockRepo.On("GetByID", user.Id).Return(user, nil)

	// Create request
	req := httptest.NewRequest(http.MethodGet, "/users/"+user.Id, nil)

	// Perform request
	resp, err := app.Test(req)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	// Verify response body
	var result map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&result)
	assert.NoError(t, err)

	assert.NotNil(t, result["user"])

	// Verify expectations
	mockRepo.AssertExpectations(t)
}

func TestUserHandler_GetUser_NotFound(t *testing.T) {
	// Setup
	app, handler, mockRepo := setupTest()

	// Setup route
	app.Get("/users/:id", handler.GetUser)

	// Setup expectations
	mockRepo.On("GetByID", "non-existent-id").Return(nil, errors.New("user not found"))

	// Create request
	req := httptest.NewRequest(http.MethodGet, "/users/non-existent-id", nil)

	// Perform request
	resp, err := app.Test(req)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, http.StatusNotFound, resp.StatusCode)

	// Verify response body
	var result map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&result)
	assert.NoError(t, err)

	assert.Equal(t, "User not found", result["error"])

	// Verify expectations
	mockRepo.AssertExpectations(t)
}

func TestUserHandler_UpdateUser(t *testing.T) {
	// Setup
	app, handler, mockRepo := setupTest()

	// Set a role in context locals to simulate authenticated admin/staff
	app.Use(func(c *fiber.Ctx) error {
		c.Locals("role", "admin") // or "staff"
		return c.Next()
	})

	// Setup route
	app.Put("/users/:id", handler.UpdateUser)

	// Create test user
	user := createTestUser()

	// Setup expectations
	mockRepo.On("GetByID", user.Id).Return(user, nil)
	// mockRepo.On("EmailExists", "updated@example.com").Return(false, nil) // Removed: Handler doesn't check email existence
	mockRepo.On("Update", mock.MatchedBy(func(u *models.User) bool {
		// Check if the updated fields are correct
		return u.Id == user.Id && u.Role == "staff" && u.IsActive == false
	})).Return(nil)

	// Create request body - only include fields the handler uses
	reqBody := map[string]interface{}{ // Use interface{} for boolean
		// "name":  "Updated User", // Removed: Handler doesn't update name
		// "email": "updated@example.com", // Removed: Handler doesn't update email
		"role":     "staff",
		"isActive": false,
	}
	jsonBody, _ := json.Marshal(reqBody)

	// Create request
	req := httptest.NewRequest(http.MethodPut, "/users/"+user.Id, bytes.NewReader(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	// Perform request
	resp, err := app.Test(req)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	// Verify response body
	var result map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&result)
	assert.NoError(t, err)

	assert.Equal(t, "User updated successfully", result["message"])
	// Further check the returned user if necessary
	userData, ok := result["user"].(map[string]interface{})
	assert.True(t, ok)
	assert.Equal(t, "staff", userData["role"])
	assert.Equal(t, false, userData["isActive"])

	// Verify expectations
	mockRepo.AssertExpectations(t)
}

func TestUserHandler_UpdateUser_InactiveUser(t *testing.T) {
	// Setup
	app, handler, mockRepo := setupTest()

	// Set a role in context locals to simulate authenticated admin/staff
	app.Use(func(c *fiber.Ctx) error {
		c.Locals("role", "admin") // or "staff"
		return c.Next()
	})

	// Setup route
	app.Put("/users/:id", handler.UpdateUser)

	// Create inactive test user
	inactiveUser := createTestUser()
	inactiveUser.IsActive = false

	// Setup expectations
	mockRepo.On("GetByID", inactiveUser.Id).Return(&models.User{
		Id:        inactiveUser.Id,
		FullName:  inactiveUser.FullName,
		Email:     inactiveUser.Email,
		Password:  inactiveUser.Password, // Will be hashed in reality
		Role:      inactiveUser.Role,
		CreatedAt: inactiveUser.CreatedAt,
		UpdatedAt: inactiveUser.UpdatedAt,
		IsActive:  false, // Explicitly false for this test case
	}, nil)
	// Add the Update expectation, as admin/staff *can* update role/isActive for inactive users
	mockRepo.On("Update", mock.MatchedBy(func(u *models.User) bool {
		return u.Id == inactiveUser.Id && u.Role == "user" && u.IsActive == true // Example: activating the user
	})).Return(nil)

	// Create request body - only include fields the handler uses
	reqBody := map[string]interface{}{ // Use interface{} for boolean
		// "name":  "Updated User", // Removed: Handler doesn't update name
		// "email": "updated@example.com", // Removed: Handler doesn't update email
		"role":     "user",
		"isActive": true, // Try to activate the user
	}
	jsonBody, _ := json.Marshal(reqBody)

	// Create request
	req := httptest.NewRequest(http.MethodPut, "/users/"+inactiveUser.Id, bytes.NewReader(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	// Perform request
	resp, err := app.Test(req)

	// Assert
	assert.NoError(t, err)
	// Expect 200 OK because admin/staff can update isActive/role even for inactive users
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	// Verify response body
	var result map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&result)
	assert.NoError(t, err)

	// Expect success message
	assert.Equal(t, "User updated successfully", result["message"])
	userData, ok := result["user"].(map[string]interface{})
	assert.True(t, ok)
	assert.Equal(t, "user", userData["role"])
	assert.Equal(t, true, userData["isActive"])
	// assert.Equal(t, "Account is inactive", result["error"]) // Removed: Update should succeed

	// Verify expectations
	mockRepo.AssertExpectations(t)
}

func TestUserHandler_DeleteUser(t *testing.T) {
	// Setup
	app, handler, mockRepo := setupTest()

	// Setup route
	app.Delete("/users/:id", handler.DeleteUser)

	// Setup expectations
	mockRepo.On("Delete", "test-id").Return(nil)

	// Create request
	req := httptest.NewRequest(http.MethodDelete, "/users/test-id", nil)

	// Perform request
	resp, err := app.Test(req)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	// Verify response body
	var result map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&result)
	assert.NoError(t, err)

	assert.Equal(t, "User deleted successfully", result["message"])

	// Verify expectations
	mockRepo.AssertExpectations(t)
}

func TestUserHandler_ActivateUser(t *testing.T) {
	// Setup
	app, handler, mockRepo := setupTest()

	// Setup route
	app.Put("/users/:id/activate", handler.ActivateUser)

	// Setup expectations
	mockRepo.On("ActivateUser", "test-id").Return(nil)

	// Create request
	req := httptest.NewRequest(http.MethodPut, "/users/test-id/activate", nil)

	// Perform request
	resp, err := app.Test(req)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	// Verify response body
	var result map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&result)
	assert.NoError(t, err)

	assert.Equal(t, "User activated successfully", result["message"])

	// Verify expectations
	mockRepo.AssertExpectations(t)
}

func TestUserHandler_DeactivateUser(t *testing.T) {
	// Setup
	app, handler, mockRepo := setupTest()

	// Setup route
	app.Put("/users/:id/deactivate", handler.DeactivateUser)

	// Setup expectations
	mockRepo.On("DeactivateUser", "test-id").Return(nil)

	// Create request
	req := httptest.NewRequest(http.MethodPut, "/users/test-id/deactivate", nil)

	// Perform request
	resp, err := app.Test(req)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	// Verify response body
	var result map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&result)
	assert.NoError(t, err)

	assert.Equal(t, "User deactivated successfully", result["message"])

	// Verify expectations
	mockRepo.AssertExpectations(t)
}

func TestUserHandler_UpdatePassword(t *testing.T) {
	// Setup
	app, handler, mockRepo := setupTest()

	// Create test user
	user := createTestUser()

	// Setup middleware to simulate authenticated user
	app.Use(func(c *fiber.Ctx) error {
		// For this test, simulate the user updating their own password
		if c.Path() == "/users/"+user.Id+"/password" {
			c.Locals("user_id", user.Id)
			c.Locals("role", user.Role) // Use the role from the test user
		}
		return c.Next()
	})

	// Setup route
	app.Put("/users/:id/password", handler.UpdatePassword)

	// Setup expectations
	mockRepo.On("GetByID", user.Id).Return(user, nil).Once()
	mockRepo.On("VerifyPassword", user.Email, "current_password").Return(user, nil).Once()
	mockRepo.On("UpdatePassword", user.Id, "new_password").Return(nil).Once()

	// Create request body
	reqBody := map[string]string{
		"current_password": "current_password",
		"new_password":     "new_password",
	}
	jsonBody, _ := json.Marshal(reqBody)

	// Create request
	req := httptest.NewRequest(http.MethodPut, "/users/"+user.Id+"/password", bytes.NewReader(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	// Perform request
	resp, err := app.Test(req)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	// Verify response body
	var result map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&result)
	assert.NoError(t, err)

	assert.Equal(t, "Password updated successfully", result["message"])

	// Verify expectations
	mockRepo.AssertExpectations(t)
}

func TestUserHandler_UpdatePassword_InactiveUser(t *testing.T) {
	// Setup
	app, handler, mockRepo := setupTest()

	// Create inactive test user
	inactiveUser := createTestUser()
	inactiveUser.IsActive = false

	// Setup middleware to simulate authenticated user
	app.Use(func(c *fiber.Ctx) error {
		// For this test, simulate the user attempting to update their own password
		if c.Path() == "/users/"+inactiveUser.Id+"/password" {
			c.Locals("user_id", inactiveUser.Id)
			c.Locals("role", inactiveUser.Role) // Use the role from the test user
		}
		return c.Next()
	})

	// Setup route
	app.Put("/users/:id/password", handler.UpdatePassword)

	// Setup expectations
	// Only GetByID should be called before the active check
	mockRepo.On("GetByID", inactiveUser.Id).Return(&models.User{
		Id:        inactiveUser.Id,
		FullName:  inactiveUser.FullName,
		Email:     inactiveUser.Email,
		Password:  inactiveUser.Password, // Will be hashed in reality
		Role:      inactiveUser.Role,
		CreatedAt: inactiveUser.CreatedAt,
		UpdatedAt: inactiveUser.UpdatedAt,
		IsActive:  false, // Explicitly false for this test case
	}, nil).Once()

	// Create request body
	reqBody := map[string]string{
		"current_password": "current_password",
		"new_password":     "new_password",
	}
	jsonBody, _ := json.Marshal(reqBody)

	// Create request
	req := httptest.NewRequest(http.MethodPut, "/users/"+inactiveUser.Id+"/password", bytes.NewReader(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	// Perform request
	resp, err := app.Test(req)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, http.StatusForbidden, resp.StatusCode) // Expect Forbidden due to inactive user

	// Verify response body
	var result map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&result)
	assert.NoError(t, err)

	assert.Equal(t, "Account is inactive", result["error"]) // Expect inactive account error

	// Verify expectations
	mockRepo.AssertExpectations(t)
}

func TestUserHandler_UpdatePassword_IncorrectCurrentPassword(t *testing.T) {
	// Setup
	app, handler, mockRepo := setupTest()

	// Create test user (defaults to active)
	user := createTestUser()

	// Setup middleware to simulate authenticated user
	app.Use(func(c *fiber.Ctx) error {
		if c.Path() == "/users/"+user.Id+"/password" {
			c.Locals("user_id", user.Id)
			c.Locals("role", user.Role)
		}
		return c.Next()
	})

	// Setup route
	app.Put("/users/:id/password", handler.UpdatePassword)

	// Create test user
	// user := createTestUser() // Moved user creation up

	// Setup expectations
	mockRepo.On("GetByID", user.Id).Return(user, nil).Once()
	mockRepo.On("VerifyPassword", user.Email, "wrong_password").Return(nil, errors.New("invalid password")).Once()

	// Create request body
	reqBody := map[string]string{
		"current_password": "wrong_password",
		"new_password":     "new_password",
	}
	jsonBody, _ := json.Marshal(reqBody)

	// Create request
	req := httptest.NewRequest(http.MethodPut, "/users/"+user.Id+"/password", bytes.NewReader(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	// Perform request
	resp, err := app.Test(req)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)

	// Verify response body
	var result map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&result)
	assert.NoError(t, err)

	assert.Equal(t, "Current password is incorrect", result["error"])

	// Verify expectations
	mockRepo.AssertExpectations(t)
}
