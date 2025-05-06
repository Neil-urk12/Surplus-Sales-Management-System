package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"oop/internal/models"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockAccessoryRepository is a mock implementation of the AccessoryRepository interface
type MockAccessoryRepository struct {
	mock.Mock
}

// GetAll mocks the GetAll method
func (m *MockAccessoryRepository) GetAll(ctx context.Context) ([]models.Accessory, error) {
	args := m.Called(ctx)
	return args.Get(0).([]models.Accessory), args.Error(1)
}

// GetByID mocks the GetByID method
func (m *MockAccessoryRepository) GetByID(ctx context.Context, id int) (models.Accessory, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(models.Accessory), args.Error(1)
}

// Create mocks the Create method
func (m *MockAccessoryRepository) Create(ctx context.Context, input models.NewAccessoryInput) (models.Accessory, error) {
	args := m.Called(ctx, input)
	return args.Get(0).(models.Accessory), args.Error(1)
}

// Update mocks the Update method
func (m *MockAccessoryRepository) Update(ctx context.Context, id int, input models.UpdateAccessoryInput) (models.Accessory, error) {
	args := m.Called(ctx, id, input)
	return args.Get(0).(models.Accessory), args.Error(1)
}

// Delete mocks the Delete method
func (m *MockAccessoryRepository) Delete(ctx context.Context, id int) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

// Helper function to setup a test Fiber app with the accessories handlers
func setupTestApp(mockRepo *MockAccessoryRepository) *fiber.App {
	app := fiber.New()
	handler := NewAccessoriesHandler(mockRepo)

	// Register routes
	app.Get("/api/accessories", handler.GetAllAccessories)
	app.Get("/api/accessories/:id", handler.GetAccessoryByID)
	app.Post("/api/accessories", handler.CreateAccessory)
	app.Put("/api/accessories/:id", handler.UpdateAccessory)
	app.Delete("/api/accessories/:id", handler.DeleteAccessory)

	return app
}

// Helper function to make HTTP requests and parse the response
func makeRequest(app *fiber.App, method, url string, body interface{}) (*http.Response, []byte, error) {
	// Prepare request
	var reqBody io.Reader
	if body != nil {
		jsonBytes, err := json.Marshal(body)
		if err != nil {
			return nil, nil, err
		}
		reqBody = bytes.NewBuffer(jsonBytes)
	}

	// Create HTTP request
	req := httptest.NewRequest(method, url, reqBody)
	req.Header.Set("Content-Type", "application/json")

	// Perform request
	resp, err := app.Test(req)
	if err != nil {
		return nil, nil, err
	}

	// Read response body
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return resp, nil, err
	}

	return resp, respBody, nil
}

// TestGetAllAccessories tests the GetAllAccessories handler
func TestGetAllAccessories(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		// Create mock repository
		mockRepo := new(MockAccessoryRepository)

		// Setup test data
		now := time.Now()
		accessories := []models.Accessory{
			{
				ID:        1,
				Name:      "Test Accessory 1",
				Make:      models.MakeOEM,
				Quantity:  10,
				Price:     100.0,
				Status:    models.StatusInStock,
				UnitColor: models.ColorBlack,
				CreatedAt: now,
				UpdatedAt: now,
			},
			{
				ID:        2,
				Name:      "Test Accessory 2",
				Make:      models.MakeAftermarket,
				Quantity:  0,
				Price:     200.0,
				Status:    models.StatusOutOfStock,
				UnitColor: models.ColorSilver,
				CreatedAt: now,
				UpdatedAt: now,
			},
		}

		// Set expectations
		mockRepo.On("GetAll", mock.Anything).Return(accessories, nil)

		// Setup app and make request
		app := setupTestApp(mockRepo)
		resp, body, err := makeRequest(app, "GET", "/api/accessories", nil)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var response map[string]interface{}
		err = json.Unmarshal(body, &response)
		assert.NoError(t, err)

		data, ok := response["data"].([]interface{})
		assert.True(t, ok)
		assert.Len(t, data, 2)
		assert.Equal(t, float64(2), response["count"])

		// Verify mock
		mockRepo.AssertExpectations(t)
	})

	t.Run("Error", func(t *testing.T) {
		// Create mock repository
		mockRepo := new(MockAccessoryRepository)

		// Set expectations - simulate a database error
		mockRepo.On("GetAll", mock.Anything).Return([]models.Accessory{}, errors.New("database error"))

		// Setup app and make request
		app := setupTestApp(mockRepo)
		resp, body, err := makeRequest(app, "GET", "/api/accessories", nil)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)

		var response map[string]interface{}
		err = json.Unmarshal(body, &response)
		assert.NoError(t, err)
		assert.Contains(t, response["error"], "Failed to retrieve accessories")

		// Verify mock
		mockRepo.AssertExpectations(t)
	})
}

// TestGetAccessoryByID tests the GetAccessoryByID handler
func TestGetAccessoryByID(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		// Create mock repository
		mockRepo := new(MockAccessoryRepository)

		// Setup test data
		now := time.Now()
		accessory := models.Accessory{
			ID:        1,
			Name:      "Test Accessory",
			Make:      models.MakeOEM,
			Quantity:  10,
			Price:     100.0,
			Status:    models.StatusInStock,
			UnitColor: models.ColorBlack,
			CreatedAt: now,
			UpdatedAt: now,
		}

		// Set expectations
		mockRepo.On("GetByID", mock.Anything, 1).Return(accessory, nil)

		// Setup app and make request
		app := setupTestApp(mockRepo)
		resp, body, err := makeRequest(app, "GET", "/api/accessories/1", nil)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var response models.Accessory
		err = json.Unmarshal(body, &response)
		assert.NoError(t, err)
		assert.Equal(t, accessory.ID, response.ID)
		assert.Equal(t, accessory.Name, response.Name)

		// Verify mock
		mockRepo.AssertExpectations(t)
	})

	t.Run("Not Found", func(t *testing.T) {
		// Create mock repository
		mockRepo := new(MockAccessoryRepository)

		// Set expectations
		mockRepo.On("GetByID", mock.Anything, 999).Return(models.Accessory{}, errors.New("accessory not found"))

		// Setup app and make request
		app := setupTestApp(mockRepo)
		resp, body, err := makeRequest(app, "GET", "/api/accessories/999", nil)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, http.StatusNotFound, resp.StatusCode)

		var response map[string]interface{}
		err = json.Unmarshal(body, &response)
		assert.NoError(t, err)
		assert.Contains(t, response["error"], "Accessory with ID 999 not found")

		// Verify mock
		mockRepo.AssertExpectations(t)
	})

	t.Run("Invalid ID", func(t *testing.T) {
		// Create mock repository
		mockRepo := new(MockAccessoryRepository)

		// Setup app and make request
		app := setupTestApp(mockRepo)
		resp, _, err := makeRequest(app, "GET", "/api/accessories/invalid", nil)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

		// No need to verify mock as it shouldn't be called
	})
}

// TestCreateAccessory tests the CreateAccessory handler
func TestCreateAccessory(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		// Create mock repository
		mockRepo := new(MockAccessoryRepository)

		// Setup test data
		now := time.Now()
		input := models.NewAccessoryInput{
			Name:      "New Accessory",
			Make:      models.MakeOEM,
			Quantity:  10,
			Price:     100.0,
			UnitColor: models.ColorBlack,
		}

		createdAccessory := models.Accessory{
			ID:        1,
			Name:      input.Name,
			Make:      input.Make,
			Quantity:  input.Quantity,
			Price:     input.Price,
			Status:    models.StatusInStock,
			UnitColor: input.UnitColor,
			CreatedAt: now,
			UpdatedAt: now,
		}

		// Set expectations
		mockRepo.On("Create", mock.Anything, input).Return(createdAccessory, nil)

		// Setup app and make request
		app := setupTestApp(mockRepo)
		resp, body, err := makeRequest(app, "POST", "/api/accessories", input)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, http.StatusCreated, resp.StatusCode)

		var response map[string]interface{}
		err = json.Unmarshal(body, &response)
		assert.NoError(t, err)
		assert.True(t, response["success"].(bool))
		assert.Equal(t, float64(1), response["id"])

		// Verify mock
		mockRepo.AssertExpectations(t)
	})

	t.Run("Validation Error", func(t *testing.T) {
		// Create mock repository
		mockRepo := new(MockAccessoryRepository)

		// Setup test data - missing required fields
		input := models.NewAccessoryInput{
			// Name is missing
			Make:      models.MakeOEM,
			Quantity:  10,
			Price:     100.0,
			UnitColor: models.ColorBlack,
		}

		// Setup app and make request
		app := setupTestApp(mockRepo)
		resp, body, err := makeRequest(app, "POST", "/api/accessories", input)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, http.StatusUnprocessableEntity, resp.StatusCode)

		var response map[string]interface{}
		err = json.Unmarshal(body, &response)
		assert.NoError(t, err)
		assert.Contains(t, response["error"], "Missing required fields")

		// No need to verify mock as it shouldn't be called
	})

	t.Run("Repository Error", func(t *testing.T) {
		// Create mock repository
		mockRepo := new(MockAccessoryRepository)

		// Setup test data
		input := models.NewAccessoryInput{
			Name:      "New Accessory",
			Make:      models.MakeOEM,
			Quantity:  10,
			Price:     100.0,
			UnitColor: models.ColorBlack,
		}

		// Set expectations - simulate a database error
		mockRepo.On("Create", mock.Anything, input).Return(models.Accessory{}, errors.New("database error"))

		// Setup app and make request
		app := setupTestApp(mockRepo)
		resp, body, err := makeRequest(app, "POST", "/api/accessories", input)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)

		var response map[string]interface{}
		err = json.Unmarshal(body, &response)
		assert.NoError(t, err)
		assert.Contains(t, response["error"], "Failed to create accessory")

		// Verify mock
		mockRepo.AssertExpectations(t)
	})
}

// TestUpdateAccessory tests the UpdateAccessory handler
func TestUpdateAccessory(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		// Create mock repository
		mockRepo := new(MockAccessoryRepository)

		// Setup test data
		now := time.Now()
		name := "Updated Accessory"
		price := 150.0
		input := models.UpdateAccessoryInput{
			Name:  &name,
			Price: &price,
		}

		updatedAccessory := models.Accessory{
			ID:        1,
			Name:      name,
			Make:      models.MakeOEM,
			Quantity:  10,
			Price:     price,
			Status:    models.StatusInStock,
			UnitColor: models.ColorBlack,
			CreatedAt: now,
			UpdatedAt: now,
		}

		// Set expectations
		mockRepo.On("Update", mock.Anything, 1, input).Return(updatedAccessory, nil)

		// Setup app and make request
		app := setupTestApp(mockRepo)
		resp, body, err := makeRequest(app, "PUT", "/api/accessories/1", input)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var response map[string]interface{}
		err = json.Unmarshal(body, &response)
		assert.NoError(t, err)
		assert.True(t, response["success"].(bool))
		assert.Equal(t, float64(1), response["id"])

		// Verify mock
		mockRepo.AssertExpectations(t)
	})

	t.Run("Not Found", func(t *testing.T) {
		// Create mock repository
		mockRepo := new(MockAccessoryRepository)

		// Setup test data
		name := "Updated Accessory"
		price := 150.0
		input := models.UpdateAccessoryInput{
			Name:  &name,
			Price: &price,
		}

		// Set expectations
		mockRepo.On("Update", mock.Anything, 999, input).Return(models.Accessory{}, errors.New("accessory not found"))

		// Setup app and make request
		app := setupTestApp(mockRepo)
		resp, body, err := makeRequest(app, "PUT", "/api/accessories/999", input)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, http.StatusNotFound, resp.StatusCode)

		var response map[string]interface{}
		err = json.Unmarshal(body, &response)
		assert.NoError(t, err)
		assert.Contains(t, response["error"], "Accessory with ID 999 not found")

		// Verify mock
		mockRepo.AssertExpectations(t)
	})

	t.Run("Invalid ID", func(t *testing.T) {
		// Create mock repository
		mockRepo := new(MockAccessoryRepository)

		// Setup app and make request
		app := setupTestApp(mockRepo)
		resp, _, err := makeRequest(app, "PUT", "/api/accessories/invalid", models.UpdateAccessoryInput{})

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

		// No need to verify mock as it shouldn't be called
	})
}

// TestDeleteAccessory tests the DeleteAccessory handler
func TestDeleteAccessory(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		// Create mock repository
		mockRepo := new(MockAccessoryRepository)

		// Set expectations
		mockRepo.On("Delete", mock.Anything, 1).Return(nil)

		// Setup app and make request
		app := setupTestApp(mockRepo)
		resp, _, err := makeRequest(app, "DELETE", "/api/accessories/1", nil)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, http.StatusNoContent, resp.StatusCode)

		// Verify mock
		mockRepo.AssertExpectations(t)
	})

	t.Run("Not Found", func(t *testing.T) {
		// Create mock repository
		mockRepo := new(MockAccessoryRepository)

		// Set expectations
		mockRepo.On("Delete", mock.Anything, 999).Return(errors.New("accessory not found"))

		// Setup app and make request
		app := setupTestApp(mockRepo)
		resp, body, err := makeRequest(app, "DELETE", "/api/accessories/999", nil)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, http.StatusNotFound, resp.StatusCode)

		var response map[string]interface{}
		err = json.Unmarshal(body, &response)
		assert.NoError(t, err)
		assert.Contains(t, response["error"], "Accessory with ID 999 not found")

		// Verify mock
		mockRepo.AssertExpectations(t)
	})

	t.Run("Invalid ID", func(t *testing.T) {
		// Create mock repository
		mockRepo := new(MockAccessoryRepository)

		// Setup app and make request
		app := setupTestApp(mockRepo)
		resp, _, err := makeRequest(app, "DELETE", "/api/accessories/invalid", nil)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

		// No need to verify mock as it shouldn't be called
	})
}
