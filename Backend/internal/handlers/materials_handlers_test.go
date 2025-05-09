package handlers

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"oop/internal/models"
	"oop/internal/repositories"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockMaterialRepository is a mock type for the MaterialRepository interface
type MockMaterialRepository struct {
	mock.Mock
}

// GetPaginated implements repositories.MaterialRepository.
func (m *MockMaterialRepository) GetPaginated(page int, limit int, searchTerm string, category string, supplier string, status string) ([]models.Material, int64, error) {
	panic("unimplemented")
}

func (m *MockMaterialRepository) GetAll(searchTerm string, category string, supplier string, status string) ([]models.Material, error) {
	args := m.Called(searchTerm, category, supplier, status)
	return args.Get(0).([]models.Material), args.Error(1)
}

func (m *MockMaterialRepository) GetByID(id int) (*models.Material, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Material), args.Error(1)
}

func (m *MockMaterialRepository) Create(material *models.Material) (int, error) {
	args := m.Called(material)
	return args.Int(0), args.Error(1)
}

func (m *MockMaterialRepository) Update(material *models.Material) error {
	args := m.Called(material)
	return args.Error(0)
}

func (m *MockMaterialRepository) Delete(id int) error {
	args := m.Called(id)
	return args.Error(0)
}

// Helper function to create a test JWT token
func createTestToken(secret []byte, userID uint, userRole string) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"role":    userRole,
		"exp":     time.Now().Add(time.Hour * 72).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(secret)
}

func setupMaterialTestApp(repo repositories.MaterialRepository, jwtSecret []byte) *fiber.App {
	app := fiber.New()
	h := NewMaterialHandlers(repo, jwtSecret)
	api := app.Group("/api") // Match the main setup
	h.RegisterMaterialRoutes(api)
	return app
}

func TestGetMaterialsHandler(t *testing.T) {
	mockRepo := new(MockMaterialRepository)
	jwtSecret := []byte("test_secret")
	app := setupMaterialTestApp(mockRepo, jwtSecret)

	testToken, _ := createTestToken(jwtSecret, 1, "admin")

	now := time.Now()
	expectedMaterials := []models.Material{
		{ID: 1, Name: "Mat 1", Category: "C1", Supplier: "S1", Quantity: 10, Status: "Active", CreatedAt: now, UpdatedAt: now},
		{ID: 2, Name: "Mat 2", Category: "C2", Supplier: "S2", Quantity: 20, Status: "Inactive", CreatedAt: now, UpdatedAt: now},
	}

	t.Run("Success - No Filters", func(t *testing.T) {
		mockRepo.On("GetAll", "", "", "", "").Return(expectedMaterials, nil).Once()

		req := httptest.NewRequest(http.MethodGet, "/api/materials", nil)
		req.Header.Set("Authorization", "Bearer "+testToken)

		resp, err := app.Test(req, -1) // Use -1 for no timeout
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var actualMaterials []models.Material
		err = json.NewDecoder(resp.Body).Decode(&actualMaterials)
		assert.NoError(t, err)
		// Zero out time fields for comparison
		zeroTimeFieldsSlice(expectedMaterials)
		zeroTimeFieldsSlice(actualMaterials)
		assert.Equal(t, expectedMaterials, actualMaterials)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Success - With Filters", func(t *testing.T) {
		filteredMaterials := []models.Material{expectedMaterials[0]}
		mockRepo.On("GetAll", "search", "C1", "S1", "Active").Return(filteredMaterials, nil).Once()

		req := httptest.NewRequest(http.MethodGet, "/api/materials?search=search&category=C1&supplier=S1&status=Active", nil)
		req.Header.Set("Authorization", "Bearer "+testToken)

		resp, err := app.Test(req, -1)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var actualMaterials []models.Material
		err = json.NewDecoder(resp.Body).Decode(&actualMaterials)
		assert.NoError(t, err)
		// Zero out time fields for comparison
		zeroTimeFieldsSlice(filteredMaterials)
		zeroTimeFieldsSlice(actualMaterials)
		assert.Equal(t, filteredMaterials, actualMaterials)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Repository Error", func(t *testing.T) {
		mockRepo.On("GetAll", "", "", "", "").Return([]models.Material{}, errors.New("db error")).Once()

		req := httptest.NewRequest(http.MethodGet, "/api/materials", nil)
		req.Header.Set("Authorization", "Bearer "+testToken)

		resp, err := app.Test(req, -1)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
		mockRepo.AssertExpectations(t)
	})
}

func TestGetMaterialHandler(t *testing.T) {
	mockRepo := new(MockMaterialRepository)
	jwtSecret := []byte("test_secret")
	app := setupMaterialTestApp(mockRepo, jwtSecret)
	testToken, _ := createTestToken(jwtSecret, 1, "admin")

	now := time.Now()
	expectedMaterial := &models.Material{ID: 1, Name: "Mat 1", Category: "C1", Supplier: "S1", Quantity: 10, Status: "Active", CreatedAt: now, UpdatedAt: now}

	t.Run("Success", func(t *testing.T) {
		mockRepo.On("GetByID", 1).Return(expectedMaterial, nil).Once()

		req := httptest.NewRequest(http.MethodGet, "/api/materials/1", nil)
		req.Header.Set("Authorization", "Bearer "+testToken)

		resp, err := app.Test(req, -1)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var actualMaterial models.Material
		err = json.NewDecoder(resp.Body).Decode(&actualMaterial)
		assert.NoError(t, err)
		// Compare with the dereferenced expected object, zeroing time fields
		expected := *expectedMaterial
		zeroTimeFields(&expected)
		zeroTimeFields(&actualMaterial)
		assert.Equal(t, expected, actualMaterial)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Invalid ID Format", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/materials/abc", nil)
		req.Header.Set("Authorization", "Bearer "+testToken)

		resp, err := app.Test(req, -1)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Not Found", func(t *testing.T) {
		mockRepo.On("GetByID", 2).Return((*models.Material)(nil), nil).Once() // Repo returns nil, nil for not found

		req := httptest.NewRequest(http.MethodGet, "/api/materials/2", nil)
		req.Header.Set("Authorization", "Bearer "+testToken)

		resp, err := app.Test(req, -1)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusNotFound, resp.StatusCode)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Repository Error", func(t *testing.T) {
		mockRepo.On("GetByID", 3).Return((*models.Material)(nil), errors.New("db error")).Once()

		req := httptest.NewRequest(http.MethodGet, "/api/materials/3", nil)
		req.Header.Set("Authorization", "Bearer "+testToken)

		resp, err := app.Test(req, -1)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
		mockRepo.AssertExpectations(t)
	})
}

func TestCreateMaterialHandler(t *testing.T) {
	mockRepo := new(MockMaterialRepository)
	jwtSecret := []byte("test_secret")
	app := setupMaterialTestApp(mockRepo, jwtSecret)
	testToken, _ := createTestToken(jwtSecret, 1, "admin")

	now := time.Now()
	newMaterialInput := models.Material{Name: "New Mat", Category: "CNew", Supplier: "SNew", Quantity: 5, Status: "Pending", Image: "new.jpg"}
	createdMaterial := models.Material{ID: 1, Name: "New Mat", Category: "CNew", Supplier: "SNew", Quantity: 5, Status: "Pending", Image: "new.jpg", CreatedAt: now, UpdatedAt: now}

	t.Run("Success", func(t *testing.T) {
		// Mock Create returning the new ID
		mockRepo.On("Create", mock.AnythingOfType("*models.Material")).Return(1, nil).Once().Run(func(args mock.Arguments) {
			marg := args.Get(0).(*models.Material)
			assert.Equal(t, newMaterialInput.Name, marg.Name) // Check input passed to repo
		})
		// Mock GetByID called after Create to return the full object
		mockRepo.On("GetByID", 1).Return(&createdMaterial, nil).Once()

		bodyBytes, _ := json.Marshal(newMaterialInput)
		req := httptest.NewRequest(http.MethodPost, "/api/materials", bytes.NewReader(bodyBytes))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+testToken)

		resp, err := app.Test(req, -1)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusCreated, resp.StatusCode)

		var actualMaterial models.Material
		err = json.NewDecoder(resp.Body).Decode(&actualMaterial)
		assert.NoError(t, err)
		// Zero out time fields for comparison
		zeroTimeFields(&createdMaterial)
		zeroTimeFields(&actualMaterial)
		assert.Equal(t, createdMaterial, actualMaterial)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Invalid Payload", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/api/materials", bytes.NewReader([]byte(`{invalid json`)))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+testToken)

		resp, err := app.Test(req, -1)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Validation Error - Missing Fields", func(t *testing.T) {
		invalidInput := models.Material{Name: "Only Name"} // Missing other required fields
		bodyBytes, _ := json.Marshal(invalidInput)
		req := httptest.NewRequest(http.MethodPost, "/api/materials", bytes.NewReader(bodyBytes))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+testToken)

		resp, err := app.Test(req, -1)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Repository Create Error", func(t *testing.T) {
		mockRepo.On("Create", mock.AnythingOfType("*models.Material")).Return(0, errors.New("db create error")).Once()

		bodyBytes, _ := json.Marshal(newMaterialInput)
		req := httptest.NewRequest(http.MethodPost, "/api/materials", bytes.NewReader(bodyBytes))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+testToken)

		resp, err := app.Test(req, -1)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Repository GetByID Error After Create", func(t *testing.T) {
		// Mock Create returning the new ID
		mockRepo.On("Create", mock.AnythingOfType("*models.Material")).Return(1, nil).Once()
		// Mock GetByID failing after successful creation
		mockRepo.On("GetByID", 1).Return((*models.Material)(nil), errors.New("db get error")).Once()

		bodyBytes, _ := json.Marshal(newMaterialInput)
		req := httptest.NewRequest(http.MethodPost, "/api/materials", bytes.NewReader(bodyBytes))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+testToken)

		resp, err := app.Test(req, -1)
		assert.NoError(t, err)
		// Handler returns input + ID if GetByID fails
		assert.Equal(t, http.StatusCreated, resp.StatusCode)
		var actualMaterial models.Material
		err = json.NewDecoder(resp.Body).Decode(&actualMaterial)
		assert.NoError(t, err)
		// Expect the input material with the ID set
		expectedResponse := newMaterialInput
		expectedResponse.ID = 1
		// Zero out time fields for comparison
		zeroTimeFields(&expectedResponse) // Assuming newMaterialInput doesn't have time fields set
		zeroTimeFields(&actualMaterial)
		assert.Equal(t, expectedResponse, actualMaterial)
		mockRepo.AssertExpectations(t)
	})
}

func TestUpdateMaterialHandler(t *testing.T) {
	mockRepo := new(MockMaterialRepository)
	jwtSecret := []byte("test_secret")
	app := setupMaterialTestApp(mockRepo, jwtSecret)
	testToken, _ := createTestToken(jwtSecret, 1, "admin")

	now := time.Now()
	updateID := 1
	updateInput := models.Material{Name: "Updated Mat", Category: "CUpdated", Supplier: "SUpdated", Quantity: 15, Status: "Active", Image: "updated.jpg"}
	updatedMaterial := models.Material{ID: updateID, Name: "Updated Mat", Category: "CUpdated", Supplier: "SUpdated", Quantity: 15, Status: "Active", Image: "updated.jpg", CreatedAt: now.Add(-time.Hour), UpdatedAt: now}

	t.Run("Success", func(t *testing.T) {
		// Mock Update succeeding
		mockRepo.On("Update", mock.MatchedBy(func(m *models.Material) bool { return m.ID == updateID && m.Name == updateInput.Name })).Return(nil).Once()
		// Mock GetByID called after Update to return the full object
		mockRepo.On("GetByID", updateID).Return(&updatedMaterial, nil).Once()

		bodyBytes, _ := json.Marshal(updateInput)
		req := httptest.NewRequest(http.MethodPut, fmt.Sprintf("/api/materials/%d", updateID), bytes.NewReader(bodyBytes))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+testToken)

		resp, err := app.Test(req, -1)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var actualMaterial models.Material
		err = json.NewDecoder(resp.Body).Decode(&actualMaterial)
		assert.NoError(t, err)
		// Ensure the exact object is compared, zeroing time fields
		zeroTimeFields(&updatedMaterial)
		zeroTimeFields(&actualMaterial)
		assert.Equal(t, updatedMaterial, actualMaterial)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Invalid ID Format", func(t *testing.T) {
		bodyBytes, _ := json.Marshal(updateInput)
		req := httptest.NewRequest(http.MethodPut, "/api/materials/abc", bytes.NewReader(bodyBytes))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+testToken)

		resp, err := app.Test(req, -1)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Invalid Payload", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPut, fmt.Sprintf("/api/materials/%d", updateID), bytes.NewReader([]byte(`{invalid`)))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+testToken)

		resp, err := app.Test(req, -1)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Validation Error - Missing Fields", func(t *testing.T) {
		invalidInput := models.Material{Name: "Only Name"}
		bodyBytes, _ := json.Marshal(invalidInput)
		req := httptest.NewRequest(http.MethodPut, fmt.Sprintf("/api/materials/%d", updateID), bytes.NewReader(bodyBytes))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+testToken)

		resp, err := app.Test(req, -1)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Repository Update Error", func(t *testing.T) {
		mockRepo.On("Update", mock.MatchedBy(func(m *models.Material) bool { return m.ID == updateID })).Return(errors.New("db update error")).Once()

		bodyBytes, _ := json.Marshal(updateInput)
		req := httptest.NewRequest(http.MethodPut, fmt.Sprintf("/api/materials/%d", updateID), bytes.NewReader(bodyBytes))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+testToken)

		resp, err := app.Test(req, -1)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Repository GetByID Error After Update", func(t *testing.T) {
		// Mock Update succeeding
		mockRepo.On("Update", mock.MatchedBy(func(m *models.Material) bool { return m.ID == updateID })).Return(nil).Once()
		// Mock GetByID failing after successful update
		mockRepo.On("GetByID", updateID).Return((*models.Material)(nil), errors.New("db get error")).Once()

		bodyBytes, _ := json.Marshal(updateInput)
		req := httptest.NewRequest(http.MethodPut, fmt.Sprintf("/api/materials/%d", updateID), bytes.NewReader(bodyBytes))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+testToken)

		resp, err := app.Test(req, -1)
		assert.NoError(t, err)
		// Handler returns No Content if GetByID fails after update
		assert.Equal(t, http.StatusNoContent, resp.StatusCode)
		mockRepo.AssertExpectations(t)
	})
}

func TestDeleteMaterialHandler(t *testing.T) {
	mockRepo := new(MockMaterialRepository)
	jwtSecret := []byte("test_secret")
	app := setupMaterialTestApp(mockRepo, jwtSecret)
	testToken, _ := createTestToken(jwtSecret, 1, "admin")

	deleteID := 1

	t.Run("Success", func(t *testing.T) {
		mockRepo.On("Delete", deleteID).Return(nil).Once()

		req := httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/api/materials/%d", deleteID), nil)
		req.Header.Set("Authorization", "Bearer "+testToken)

		resp, err := app.Test(req, -1)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusNoContent, resp.StatusCode)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Invalid ID Format", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodDelete, "/api/materials/abc", nil)
		req.Header.Set("Authorization", "Bearer "+testToken)

		resp, err := app.Test(req, -1)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Repository Error", func(t *testing.T) {
		mockRepo.On("Delete", deleteID).Return(errors.New("db delete error")).Once()

		req := httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/api/materials/%d", deleteID), nil)
		req.Header.Set("Authorization", "Bearer "+testToken)

		resp, err := app.Test(req, -1)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
		mockRepo.AssertExpectations(t)
	})
}

// Helper function to zero out time fields in a single Material struct
func zeroTimeFields(m *models.Material) {
	if m != nil {
		m.CreatedAt = time.Time{}
		m.UpdatedAt = time.Time{}
	}
}

// Helper function to zero out time fields in a slice of Material structs
func zeroTimeFieldsSlice(materials []models.Material) {
	for i := range materials {
		zeroTimeFields(&materials[i])
	}
}
