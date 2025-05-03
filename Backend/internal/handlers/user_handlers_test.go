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
	handler := NewUserHandler(mockRepo)
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
		"name":     "Test User",
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
		"name":     "Test User",
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

	// Setup expectations
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

	// Setup route
	app.Put("/users/:id", handler.UpdateUser)

	// Create test user
	user := createTestUser()

	// Setup expectations
	mockRepo.On("GetByID", user.Id).Return(user, nil)
	mockRepo.On("EmailExists", "updated@example.com").Return(false, nil)
	mockRepo.On("Update", mock.AnythingOfType("*models.User")).Return(nil)

	// Create request body
	reqBody := map[string]string{
		"name":  "Updated User",
		"email": "updated@example.com",
		"role":  "admin",
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
	assert.NotNil(t, result["user"])

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

	// Setup route
	app.Put("/users/:id/password", handler.UpdatePassword)

	// Create test user
	user := createTestUser()

	// Setup expectations
	mockRepo.On("GetByID", user.Id).Return(user, nil)
	mockRepo.On("VerifyPassword", user.Email, "current_password").Return(user, nil)
	mockRepo.On("UpdatePassword", user.Id, "new_password").Return(nil)

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

func TestUserHandler_UpdatePassword_IncorrectCurrentPassword(t *testing.T) {
	// Setup
	app, handler, mockRepo := setupTest()

	// Setup route
	app.Put("/users/:id/password", handler.UpdatePassword)

	// Create test user
	user := createTestUser()

	// Setup expectations
	mockRepo.On("GetByID", user.Id).Return(user, nil)
	mockRepo.On("VerifyPassword", user.Email, "wrong_password").Return(nil, errors.New("invalid password"))

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
