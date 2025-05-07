package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
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
func (m *MockAccessoryRepository) Create(ctx context.Context, input models.NewAccessoryInput) (int, error) {
	args := m.Called(ctx, input)
	return args.Int(0), args.Error(1)
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
			Name:      "New Test Accessory",
			Make:      models.MakeCustom,
			Quantity:  5,
			Price:     150.0,
			UnitColor: models.ColorCustom,
			Image:     "new_image.png",
		}

		newAccessoryID := 123 // ID returned by Create
		expectedAccessory := models.Accessory{
			ID:        newAccessoryID,
			Name:      input.Name,
			Make:      input.Make,
			Quantity:  input.Quantity,
			Price:     input.Price,
			Status:    models.StatusInStock,
			UnitColor: input.UnitColor,
			Image:     input.Image,
			CreatedAt: now,
			UpdatedAt: now,
		}

		// Set expectations
		mockRepo.On("Create", mock.Anything, input).Return(newAccessoryID, nil)
		mockRepo.On("GetByID", mock.Anything, newAccessoryID).Return(expectedAccessory, nil)

		// Setup app and make request
		app := setupTestApp(mockRepo)
		resp, body, err := makeRequest(app, "POST", "/api/accessories", input)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, http.StatusCreated, resp.StatusCode)

		var response map[string]interface{}
		err = json.Unmarshal(body, &response)
		assert.NoError(t, err)

		assert.True(t, response["success"].(bool), "Expected success to be true")
		assert.Equal(t, "Accessory created successfully", response["message"])

		// Check for the 'data' field
		data, ok := response["data"].(map[string]interface{})
		assert.True(t, ok, "Expected 'data' field in the response")

		// Assertions for the accessory data
		assert.Equal(t, float64(expectedAccessory.ID), data["id"].(float64))
		assert.Equal(t, expectedAccessory.Name, data["name"].(string))
		assert.Equal(t, string(expectedAccessory.Make), data["make"].(string))
		assert.Equal(t, float64(expectedAccessory.Quantity), data["quantity"].(float64))
		assert.Equal(t, expectedAccessory.Price, data["price"].(float64))
		assert.Equal(t, string(expectedAccessory.Status), data["status"].(string))
		assert.Equal(t, string(expectedAccessory.UnitColor), data["unit_color"].(string))
		assert.Equal(t, expectedAccessory.Image, data["image"].(string))

		// Verify mock
		mockRepo.AssertExpectations(t)
	})

	t.Run("Missing Required Fields", func(t *testing.T) {
		// Create mock repository
		mockRepo := new(MockAccessoryRepository)

		// Test data with missing fields
		input := models.NewAccessoryInput{
			Name: "Test Accessory",
			// Make and UnitColor are missing
		}

		// No expectation for mockRepo.Create as it shouldn't be called

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

		// Verify mock (ensure Create was not called)
		mockRepo.AssertNotCalled(t, "Create", mock.Anything, mock.Anything)
	})

	t.Run("Invalid JSON", func(t *testing.T) {
		mockRepo := new(MockAccessoryRepository)
		app := setupTestApp(mockRepo)

		// Create HTTP request with invalid JSON
		req := httptest.NewRequest("POST", "/api/accessories", bytes.NewBufferString("invalid-json"))
		req.Header.Set("Content-Type", "application/json")

		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

		var responseBody map[string]interface{}
		_ = json.NewDecoder(resp.Body).Decode(&responseBody)
		assert.Contains(t, responseBody["error"], "Invalid JSON format")
		mockRepo.AssertNotCalled(t, "Create", mock.Anything, mock.Anything)
	})

	t.Run("Repository Error", func(t *testing.T) {
		mockRepo := new(MockAccessoryRepository)
		input := models.NewAccessoryInput{
			Name:      "Error Accessory",
			Make:      models.MakeOEM,
			Quantity:  1,
			Price:     10.0,
			UnitColor: models.ColorBlack,
		}
		mockRepo.On("Create", mock.Anything, input).Return(0, errors.New("database create error"))

		app := setupTestApp(mockRepo)
		resp, body, err := makeRequest(app, "POST", "/api/accessories", input)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)

		var responseBody map[string]interface{}
		err = json.Unmarshal(body, &responseBody)
		assert.NoError(t, err)
		assert.Contains(t, responseBody["error"], "Failed to create accessory")
		mockRepo.AssertExpectations(t)
	})

	t.Run("Error Fetching Created Accessory", func(t *testing.T) {
		mockRepo := new(MockAccessoryRepository)
		input := models.NewAccessoryInput{
			Name:      "Test Accessory",
			Make:      models.MakeOEM,
			Quantity:  1,
			Price:     10.0,
			UnitColor: models.ColorBlack,
		}

		createdID := 123
		mockRepo.On("Create", mock.Anything, input).Return(createdID, nil)
		mockRepo.On("GetByID", mock.Anything, createdID).Return(models.Accessory{}, errors.New("error fetching accessory"))

		app := setupTestApp(mockRepo)
		resp, body, err := makeRequest(app, "POST", "/api/accessories", input)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)

		var responseBody map[string]interface{}
		err = json.Unmarshal(body, &responseBody)
		assert.NoError(t, err)
		assert.Contains(t, responseBody["error"], "Accessory created but failed to retrieve details")
		mockRepo.AssertExpectations(t)
	})
}

// TestUpdateAccessory tests the UpdateAccessory handler
func TestUpdateAccessory(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		// Create mock repository
		mockRepo := new(MockAccessoryRepository)

		// Setup test data
		accessoryID := 1
		now := time.Now()
		updateInput := models.UpdateAccessoryInput{
			Name:     ptrToString("Updated Test Accessory"),
			Quantity: ptrToInt(15),
			Price:    ptrToFloat64(175.0),
		}
		updatedAccessory := models.Accessory{
			ID:        accessoryID,
			Name:      *updateInput.Name,
			Make:      models.MakeOEM, // Assuming make is not updated or fetched from existing
			Quantity:  *updateInput.Quantity,
			Price:     *updateInput.Price,
			Status:    models.StatusAvailable, // Or determined by quantity
			UnitColor: models.ColorBlack,      // Assuming color is not updated
			Image:     "original_image.png",
			CreatedAt: now.Add(-time.Hour), // Original creation time
			UpdatedAt: now,                 // Updated time
		}

		// Set expectations
		mockRepo.On("Update", mock.Anything, accessoryID, updateInput).Return(updatedAccessory, nil)

		// Setup app and make request
		app := setupTestApp(mockRepo)
		resp, body, err := makeRequest(app, "PUT", "/api/accessories/"+fmt.Sprintf("%d", accessoryID), updateInput)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var response map[string]interface{}
		err = json.Unmarshal(body, &response)
		assert.NoError(t, err)

		assert.True(t, response["success"].(bool), "Expected success to be true")
		assert.Equal(t, "Accessory updated successfully", response["message"])

		// Check for the 'data' field
		data, ok := response["data"].(map[string]interface{})
		assert.True(t, ok, "Expected 'data' field in the response")

		// Assertions for the accessory data
		assert.Equal(t, float64(updatedAccessory.ID), data["id"].(float64))
		assert.Equal(t, updatedAccessory.Name, data["name"].(string))
		assert.Equal(t, float64(updatedAccessory.Quantity), data["quantity"].(float64))
		assert.Equal(t, updatedAccessory.Price, data["price"].(float64))
		assert.Equal(t, string(updatedAccessory.Status), data["status"].(string))

		// Verify mock
		mockRepo.AssertExpectations(t)
	})

	t.Run("Accessory Not Found", func(t *testing.T) {
		// Create mock repository
		mockRepo := new(MockAccessoryRepository)
		accessoryID := 999 // Non-existent ID
		updateInput := models.UpdateAccessoryInput{Name: ptrToString("No Such Accessory")}

		// Set expectations
		mockRepo.On("Update", mock.Anything, accessoryID, updateInput).Return(models.Accessory{}, errors.New("accessory with ID 999 not found for update"))

		// Setup app and make request
		app := setupTestApp(mockRepo)
		resp, body, err := makeRequest(app, "PUT", "/api/accessories/"+fmt.Sprintf("%d", accessoryID), updateInput)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, http.StatusNotFound, resp.StatusCode)

		var response map[string]interface{}
		err = json.Unmarshal(body, &response)
		assert.NoError(t, err)
		assert.Contains(t, response["error"], "Accessory with ID 999 not found for update")

		// Verify mock
		mockRepo.AssertExpectations(t)
	})

	t.Run("Invalid ID Format", func(t *testing.T) {
		mockRepo := new(MockAccessoryRepository)
		app := setupTestApp(mockRepo)
		updateInput := models.UpdateAccessoryInput{Name: ptrToString("Test Update")}

		resp, body, err := makeRequest(app, "PUT", "/api/accessories/invalid-id", updateInput)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

		var responseBody map[string]interface{}
		err = json.Unmarshal(body, &responseBody)
		assert.NoError(t, err)
		assert.Contains(t, responseBody["error"], "Invalid ID format")
		mockRepo.AssertNotCalled(t, "Update", mock.Anything, mock.Anything, mock.Anything)
	})

	t.Run("Invalid JSON for Update", func(t *testing.T) {
		mockRepo := new(MockAccessoryRepository)
		app := setupTestApp(mockRepo)

		req := httptest.NewRequest("PUT", "/api/accessories/1", bytes.NewBufferString("invalid-json"))
		req.Header.Set("Content-Type", "application/json")

		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

		var responseBody map[string]interface{}
		_ = json.NewDecoder(resp.Body).Decode(&responseBody)
		assert.Contains(t, responseBody["error"], "Invalid JSON format")
		mockRepo.AssertNotCalled(t, "Update", mock.Anything, mock.Anything, mock.Anything)
	})

	t.Run("Repository Update Error", func(t *testing.T) {
		mockRepo := new(MockAccessoryRepository)
		accessoryID := 1
		updateInput := models.UpdateAccessoryInput{Name: ptrToString("Error Update")}
		mockRepo.On("Update", mock.Anything, accessoryID, updateInput).Return(models.Accessory{}, errors.New("database update error"))

		app := setupTestApp(mockRepo)
		resp, body, err := makeRequest(app, "PUT", "/api/accessories/"+fmt.Sprintf("%d", accessoryID), updateInput)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)

		var responseBody map[string]interface{}
		err = json.Unmarshal(body, &responseBody)
		assert.NoError(t, err)
		assert.Contains(t, responseBody["error"], "Failed to update accessory")
		mockRepo.AssertExpectations(t)
	})
}

// Helper functions for pointers (if not already defined elsewhere)
func ptrToString(s string) *string {
	return &s
}

func ptrToInt(i int) *int {
	return &i
}

func ptrToFloat64(f float64) *float64 {
	return &f
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
