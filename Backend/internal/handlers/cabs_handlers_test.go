package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"oop/internal/models"
	"oop/internal/repositories"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// MockCabsRepository is a mock implementation of CabsRepository for testing handlers.
// It allows setting expectations on function calls.
type MockCabsRepository struct {
	GetCabsFn    func(filters map[string]interface{}) ([]models.MultiCab, error)
	GetCabByIDFn func(id int) (*models.MultiCab, error)
	AddCabFn     func(cab models.MultiCab) (*models.MultiCab, error)
	UpdateCabFn  func(id int, cab models.MultiCab) (*models.MultiCab, error)
	DeleteCabFn  func(id int) error
}

// Implement the CabsRepository interface for the mock
func (m *MockCabsRepository) GetCabs(filters map[string]interface{}) ([]models.MultiCab, error) {
	if m.GetCabsFn != nil {
		return m.GetCabsFn(filters)
	}
	return nil, fmt.Errorf("mock GetCabsFn not implemented")
}

func (m *MockCabsRepository) GetCabByID(id int) (*models.MultiCab, error) {
	if m.GetCabByIDFn != nil {
		return m.GetCabByIDFn(id)
	}
	return nil, fmt.Errorf("mock GetCabByIDFn not implemented")
}

func (m *MockCabsRepository) AddCab(cab models.MultiCab) (*models.MultiCab, error) {
	if m.AddCabFn != nil {
		return m.AddCabFn(cab)
	}
	return nil, fmt.Errorf("mock AddCabFn not implemented")
}

func (m *MockCabsRepository) UpdateCab(id int, cab models.MultiCab) (*models.MultiCab, error) {
	if m.UpdateCabFn != nil {
		return m.UpdateCabFn(id, cab)
	}
	return nil, fmt.Errorf("mock UpdateCabFn not implemented")
}

func (m *MockCabsRepository) DeleteCab(id int) error {
	if m.DeleteCabFn != nil {
		return m.DeleteCabFn(id)
	}
	return fmt.Errorf("mock DeleteCabFn not implemented")
}

// Helper to setup Fiber app with handlers using a provided (mock) repository
func setupAppWithMockRepo(repo repositories.CabsRepository) *fiber.App {
	h := NewCabsHandlers(repo)
	app := fiber.New()

	app.Get("/api/v1/cabs", h.GetCabs)
	app.Get("/api/v1/cabs/:id", h.GetCabByID)
	app.Post("/api/v1/cabs", h.AddCab)
	app.Put("/api/v1/cabs/:id", h.UpdateCab)
	app.Delete("/api/v1/cabs/:id", h.DeleteCab)

	return app
}

func TestGetCabs_Handler_NoFilters(t *testing.T) {
	mockRepo := &MockCabsRepository{}
	app := setupAppWithMockRepo(mockRepo)

	expectedCabs := []models.MultiCab{
		{ID: 1, Name: "Cab 1"},
		{ID: 2, Name: "Cab 2"},
	}

	mockRepo.GetCabsFn = func(filters map[string]interface{}) ([]models.MultiCab, error) {
		assert.Empty(t, filters, "Expected empty filters map for no-filter request")
		return expectedCabs, nil
	}

	req := httptest.NewRequest(http.MethodGet, "/api/v1/cabs", nil)
	resp, err := app.Test(req, -1)

	require.NoError(t, err, "Request execution failed")
	assert.Equal(t, http.StatusOK, resp.StatusCode, "Expected status OK")

	var cabs []models.MultiCab
	err = json.NewDecoder(resp.Body).Decode(&cabs)
	require.NoError(t, err, "Failed to decode response body")
	assert.Equal(t, expectedCabs, cabs, "Expected cabs returned from mock")
	resp.Body.Close()
}

func TestGetCabs_Handler_WithFilters(t *testing.T) {
	mockRepo := &MockCabsRepository{}
	app := setupAppWithMockRepo(mockRepo)

	// Test filter by make
	t.Run("Filter by Make", func(t *testing.T) {
		expectedCabsMake := []models.MultiCab{{ID: 3, Name: "Porsche 1", Make: "Porsche"}}
		mockRepo.GetCabsFn = func(filters map[string]interface{}) ([]models.MultiCab, error) {
			assert.Equal(t, map[string]interface{}{"make": "Porsche"}, filters, "Expected 'make' filter")
			return expectedCabsMake, nil
		}
		reqMake := httptest.NewRequest(http.MethodGet, "/api/v1/cabs?make=Porsche", nil)
		respMake, _ := app.Test(reqMake, -1)
		require.Equal(t, http.StatusOK, respMake.StatusCode)
		var cabsMake []models.MultiCab
		err := json.NewDecoder(respMake.Body).Decode(&cabsMake)
		require.NoError(t, err)
		assert.Equal(t, expectedCabsMake, cabsMake)
		respMake.Body.Close()
	})

	// Test filter by status
	t.Run("Filter by Status", func(t *testing.T) {
		expectedCabsStatus := []models.MultiCab{{ID: 4, Name: "Cab 4", Status: "Available"}}
		mockRepo.GetCabsFn = func(filters map[string]interface{}) ([]models.MultiCab, error) {
			assert.Equal(t, map[string]interface{}{"status": "Available"}, filters, "Expected 'status' filter")
			return expectedCabsStatus, nil
		}
		reqStatus := httptest.NewRequest(http.MethodGet, "/api/v1/cabs?status=Available", nil)
		respStatus, _ := app.Test(reqStatus, -1)
		require.Equal(t, http.StatusOK, respStatus.StatusCode)
		var cabsStatus []models.MultiCab
		err := json.NewDecoder(respStatus.Body).Decode(&cabsStatus)
		require.NoError(t, err)
		assert.Equal(t, expectedCabsStatus, cabsStatus)
		respStatus.Body.Close()
	})

	// Test filter by search term
	t.Run("Filter by Search", func(t *testing.T) {
		expectedCabsSearch := []models.MultiCab{{ID: 6, Name: "Navara"}}
		mockRepo.GetCabsFn = func(filters map[string]interface{}) ([]models.MultiCab, error) {
			assert.Equal(t, map[string]interface{}{"search": "Navara"}, filters, "Expected 'search' filter")
			return expectedCabsSearch, nil
		}
		reqSearch := httptest.NewRequest(http.MethodGet, "/api/v1/cabs?search=Navara", nil)
		respSearch, _ := app.Test(reqSearch, -1)
		require.Equal(t, http.StatusOK, respSearch.StatusCode)
		var cabsSearch []models.MultiCab
		err := json.NewDecoder(respSearch.Body).Decode(&cabsSearch)
		require.NoError(t, err)
		assert.Equal(t, expectedCabsSearch, cabsSearch)
		respSearch.Body.Close()
	})

	// Test combined filters
	t.Run("Combined Filters", func(t *testing.T) {
		expectedCabsCombined := []models.MultiCab{{ID: 5, Name: "Porsche 2", Make: "Porsche", Status: "In Stock"}}
		mockRepo.GetCabsFn = func(filters map[string]interface{}) ([]models.MultiCab, error) {
			expectedFilters := map[string]interface{}{"make": "Porsche", "status": "In Stock"}
			assert.Equal(t, expectedFilters, filters, "Expected combined filters")
			return expectedCabsCombined, nil
		}
		reqCombined := httptest.NewRequest(http.MethodGet, "/api/v1/cabs?make=Porsche&status=In%20Stock", nil)
		respCombined, _ := app.Test(reqCombined, -1)
		require.Equal(t, http.StatusOK, respCombined.StatusCode)
		var cabsCombined []models.MultiCab
		err := json.NewDecoder(respCombined.Body).Decode(&cabsCombined)
		require.NoError(t, err)
		assert.Equal(t, expectedCabsCombined, cabsCombined)
		respCombined.Body.Close()
	})

	// Test repository error
	t.Run("Repository Error", func(t *testing.T) {
		mockRepo.GetCabsFn = func(filters map[string]interface{}) ([]models.MultiCab, error) {
			return nil, fmt.Errorf("internal database error")
		}
		reqErr := httptest.NewRequest(http.MethodGet, "/api/v1/cabs", nil)
		respErr, _ := app.Test(reqErr, -1)
		assert.Equal(t, http.StatusInternalServerError, respErr.StatusCode)
		// Assert error message in body if needed
		var body map[string]interface{}
		json.NewDecoder(respErr.Body).Decode(&body)
		assert.Contains(t, body["error"], "Failed to retrieve cabs")
		respErr.Body.Close()
	})
}

func TestGetCabByID_Handler_Exists(t *testing.T) {
	mockRepo := &MockCabsRepository{}
	app := setupAppWithMockRepo(mockRepo)
	id := 1
	expectedCab := &models.MultiCab{ID: id, Name: "RX‑7"}

	mockRepo.GetCabByIDFn = func(reqID int) (*models.MultiCab, error) {
		assert.Equal(t, id, reqID, "Expected correct ID passed to repo")
		return expectedCab, nil
	}

	req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/v1/cabs/%d", id), nil)
	resp, err := app.Test(req, -1)

	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var cab models.MultiCab
	err = json.NewDecoder(resp.Body).Decode(&cab)
	require.NoError(t, err)
	assert.Equal(t, *expectedCab, cab)
	resp.Body.Close()
}

func TestGetCabByID_Handler_NotExists(t *testing.T) {
	mockRepo := &MockCabsRepository{}
	app := setupAppWithMockRepo(mockRepo)
	id := 99 // Non-existent ID

	mockRepo.GetCabByIDFn = func(reqID int) (*models.MultiCab, error) {
		assert.Equal(t, id, reqID, "Expected correct ID passed to repo")
		return nil, fmt.Errorf("cab with ID %d not found", reqID)
	}

	req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/v1/cabs/%d", id), nil)
	resp, err := app.Test(req, -1)

	require.NoError(t, err)
	assert.Equal(t, http.StatusNotFound, resp.StatusCode, "Expected Not Found status")
	// Assert error message in body
	var body map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&body)
	assert.Contains(t, body["error"], fmt.Sprintf("Cab with ID %d not found", id))
	resp.Body.Close()
}

func TestGetCabByID_Handler_InvalidID(t *testing.T) {
	// This test does not involve the repository, only Fiber's parameter parsing.
	app := setupAppWithMockRepo(&MockCabsRepository{}) // Pass a dummy mock

	req := httptest.NewRequest(http.MethodGet, "/api/v1/cabs/invalid", nil)
	resp, err := app.Test(req, -1)

	require.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode, "Expected Bad Request for invalid ID format")
	var body map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&body)
	assert.Contains(t, body["error"], "Invalid ID format")
	resp.Body.Close()
}

func TestAddCab_Handler_Success(t *testing.T) {
	mockRepo := &MockCabsRepository{}
	app := setupAppWithMockRepo(mockRepo)

	newCabData := models.MultiCab{
		Name:      "Test Add Car",
		Make:      "Test Add Make",
		Quantity:  5,
		Price:     250000,
		Status:    "Available", // Required by handler validation
		UnitColor: "Blue",      // Required by handler validation
		Image:     "add.png",
	}
	expectedAddedCab := newCabData
	expectedAddedCab.ID = 8 // Simulate repo assigning ID
	now := time.Now()
	expectedAddedCab.CreatedAt = now // Simulate repo assigning timestamp
	expectedAddedCab.UpdatedAt = now // Simulate repo assigning timestamp

	mockRepo.AddCabFn = func(cab models.MultiCab) (*models.MultiCab, error) {
		assert.Equal(t, newCabData.Name, cab.Name)
		assert.Equal(t, newCabData.Make, cab.Make)
		assert.Equal(t, 0, cab.ID, "ID should be 0 before repo adds it")
		expectedAddedCab.CreatedAt = cab.CreatedAt // Use the timestamp generated by repo mock if needed
		expectedAddedCab.UpdatedAt = cab.UpdatedAt
		return &expectedAddedCab, nil
	}

	bodyBytes, _ := json.Marshal(newCabData)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/cabs", bytes.NewReader(bodyBytes))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req, -1)

	require.NoError(t, err)
	assert.Equal(t, http.StatusCreated, resp.StatusCode, "Expected status Created")

	var addedCabResp models.MultiCab
	err = json.NewDecoder(resp.Body).Decode(&addedCabResp)
	require.NoError(t, err)
	assert.Equal(t, expectedAddedCab.ID, addedCabResp.ID)
	assert.Equal(t, expectedAddedCab.Name, addedCabResp.Name)
	assert.Equal(t, expectedAddedCab.Make, addedCabResp.Make)
	// Timestamps are tricky with mocks, check if needed or use AnyArg
	resp.Body.Close()
}

func TestAddCab_Handler_InvalidInput(t *testing.T) {
	// These tests check handler-level validation (parsing, required fields)
	// before repository interaction.
	app := setupAppWithMockRepo(&MockCabsRepository{}) // Dummy mock needed

	// Malformed JSON
	t.Run("Invalid JSON", func(t *testing.T) {
		invalidJSON := `{"name": "Test"`
		reqInvalidJSON := httptest.NewRequest(http.MethodPost, "/api/v1/cabs", bytes.NewReader([]byte(invalidJSON)))
		reqInvalidJSON.Header.Set("Content-Type", "application/json")
		respInvalidJSON, _ := app.Test(reqInvalidJSON, -1)
		assert.Equal(t, http.StatusBadRequest, respInvalidJSON.StatusCode)
		var body map[string]interface{}
		json.NewDecoder(respInvalidJSON.Body).Decode(&body)
		assert.Contains(t, body["error"], "Invalid JSON format")
		respInvalidJSON.Body.Close()
	})

	// Missing required fields (as defined in handler)
	t.Run("Missing Required Fields", func(t *testing.T) {
		missingFieldsData := `{"quantity": 5}`
		reqMissingFields := httptest.NewRequest(http.MethodPost, "/api/v1/cabs", bytes.NewReader([]byte(missingFieldsData)))
		reqMissingFields.Header.Set("Content-Type", "application/json")
		respMissingFields, _ := app.Test(reqMissingFields, -1)
		assert.Equal(t, http.StatusUnprocessableEntity, respMissingFields.StatusCode)
		var body map[string]interface{}
		json.NewDecoder(respMissingFields.Body).Decode(&body)
		assert.Contains(t, body["error"], "Missing required fields")
		respMissingFields.Body.Close()
	})

	// Repository validation error (e.g., from repo.AddCab)
	t.Run("Repository Validation Error", func(t *testing.T) {
		mockRepo := &MockCabsRepository{}
		app := setupAppWithMockRepo(mockRepo)
		validPayload := models.MultiCab{Name: "Valid", Make: "Valid", UnitColor: "Valid", Status: "Valid"}
		bodyBytes, _ := json.Marshal(validPayload)

		mockRepo.AddCabFn = func(cab models.MultiCab) (*models.MultiCab, error) {
			return nil, fmt.Errorf("repo validation: cannot be empty")
		}

		req := httptest.NewRequest(http.MethodPost, "/api/v1/cabs", bytes.NewReader(bodyBytes))
		req.Header.Set("Content-Type", "application/json")
		resp, _ := app.Test(req, -1)

		assert.Equal(t, http.StatusUnprocessableEntity, resp.StatusCode)
		var body map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&body)
		assert.Contains(t, body["error"], "cannot be empty")
		resp.Body.Close()
	})
}

func TestUpdateCab_Handler_Success(t *testing.T) {
	mockRepo := &MockCabsRepository{}
	app := setupAppWithMockRepo(mockRepo)
	idToUpdate := 1

	updatePayload := models.MultiCab{
		Name:      "RX-7 Updated Handler",
		Make:      "Mazda",
		Quantity:  5,
		Price:     7100000,
		Status:    "Low Stock",
		UnitColor: "Yellow",
		Image:     "rx7_updated.jpg",
	}
	expectedUpdatedCab := updatePayload
	expectedUpdatedCab.ID = idToUpdate
	now := time.Now()
	expectedUpdatedCab.UpdatedAt = now // Simulate repo updating timestamp
	// Assuming CreatedAt is preserved by repo logic and not needed in mock return here

	mockRepo.UpdateCabFn = func(id int, cab models.MultiCab) (*models.MultiCab, error) {
		assert.Equal(t, idToUpdate, id)
		assert.Equal(t, updatePayload.Name, cab.Name)
		assert.Equal(t, updatePayload.Make, cab.Make)
		assert.Equal(t, idToUpdate, cab.ID)          // Ensure ID from payload is ignored by repo call
		expectedUpdatedCab.UpdatedAt = cab.UpdatedAt // Match timestamp
		return &expectedUpdatedCab, nil
	}

	bodyBytes, _ := json.Marshal(updatePayload)
	req := httptest.NewRequest(http.MethodPut, fmt.Sprintf("/api/v1/cabs/%d", idToUpdate), bytes.NewReader(bodyBytes))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req, -1)

	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode, "Expected status OK")

	var updatedCabResp models.MultiCab
	err = json.NewDecoder(resp.Body).Decode(&updatedCabResp)
	require.NoError(t, err)
	assert.Equal(t, expectedUpdatedCab.ID, updatedCabResp.ID)
	assert.Equal(t, expectedUpdatedCab.Name, updatedCabResp.Name)
	assert.Equal(t, expectedUpdatedCab.Status, updatedCabResp.Status)
	// Timestamps are tricky, assert UpdatedAt if stable
	resp.Body.Close()
}

func TestUpdateCab_Handler_NotExists(t *testing.T) {
	mockRepo := &MockCabsRepository{}
	app := setupAppWithMockRepo(mockRepo)
	id := 99 // Non-existent ID
	updatePayload := models.MultiCab{Name: "Doesn't Matter"}

	mockRepo.UpdateCabFn = func(reqID int, cab models.MultiCab) (*models.MultiCab, error) {
		assert.Equal(t, id, reqID)
		return nil, fmt.Errorf("cab with ID %d not found for update", reqID)
	}

	bodyBytes, _ := json.Marshal(updatePayload)
	req := httptest.NewRequest(http.MethodPut, fmt.Sprintf("/api/v1/cabs/%d", id), bytes.NewReader(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(req, -1)

	require.NoError(t, err)
	assert.Equal(t, http.StatusNotFound, resp.StatusCode, "Expected Not Found")
	var body map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&body)
	assert.Contains(t, body["error"], fmt.Sprintf("Cab with ID %d not found for update", id))
	resp.Body.Close()
}

func TestUpdateCab_Handler_InvalidID(t *testing.T) {
	app := setupAppWithMockRepo(&MockCabsRepository{}) // Dummy mock
	updatePayload := `{"name": "Doesn't Matter"}`

	req := httptest.NewRequest(http.MethodPut, "/api/v1/cabs/invalid", bytes.NewReader([]byte(updatePayload)))
	req.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(req, -1)

	require.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode, "Expected Bad Request")
	resp.Body.Close()
}

func TestUpdateCab_Handler_InvalidInput(t *testing.T) {
	app := setupAppWithMockRepo(&MockCabsRepository{}) // Dummy mock
	id := 1
	invalidJSON := `{"name": "Test"`

	req := httptest.NewRequest(http.MethodPut, fmt.Sprintf("/api/v1/cabs/%d", id), bytes.NewReader([]byte(invalidJSON)))
	req.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(req, -1)

	require.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode, "Expected Bad Request")
	resp.Body.Close()
}

func TestDeleteCab_Handler_Success(t *testing.T) {
	mockRepo := &MockCabsRepository{}
	app := setupAppWithMockRepo(mockRepo)
	idToDelete := 1
	deleteCalled := false

	mockRepo.DeleteCabFn = func(reqID int) error {
		assert.Equal(t, idToDelete, reqID)
		deleteCalled = true
		return nil // Success
	}

	req := httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/api/v1/cabs/%d", idToDelete), nil)
	resp, err := app.Test(req, -1)

	require.NoError(t, err)
	assert.Equal(t, http.StatusNoContent, resp.StatusCode, "Expected No Content")
	assert.True(t, deleteCalled, "Expected repo DeleteCab to be called")
	resp.Body.Close()
}

func TestDeleteCab_Handler_NotExists(t *testing.T) {
	mockRepo := &MockCabsRepository{}
	app := setupAppWithMockRepo(mockRepo)
	id := 99 // Non-existent ID
	deleteCalled := false

	mockRepo.DeleteCabFn = func(reqID int) error {
		assert.Equal(t, id, reqID)
		deleteCalled = true
		return fmt.Errorf("cab with ID %d not found for deletion", reqID)
	}

	req := httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/api/v1/cabs/%d", id), nil)
	resp, err := app.Test(req, -1)

	require.NoError(t, err)
	assert.Equal(t, http.StatusNotFound, resp.StatusCode, "Expected Not Found")
	assert.True(t, deleteCalled, "Expected repo DeleteCab to be called")
	var body map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&body)
	assert.Contains(t, body["error"], fmt.Sprintf("Cab with ID %d not found for deletion", id))
	resp.Body.Close()
}

func TestDeleteCab_Handler_InvalidID(t *testing.T) {
	app := setupAppWithMockRepo(&MockCabsRepository{}) // Dummy mock

	req := httptest.NewRequest(http.MethodDelete, "/api/v1/cabs/invalid", nil)
	resp, err := app.Test(req, -1)

	require.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode, "Expected Bad Request")
	resp.Body.Close()
}
