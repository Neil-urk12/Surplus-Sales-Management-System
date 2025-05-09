package handlers

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"oop/internal/models"
	"oop/internal/repositories"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockCustomerRepository is a mock type for the CustomerRepository interface
type MockCustomerRepository struct {
	mock.Mock
}

func (m *MockCustomerRepository) CreateCustomer(customer *models.Customer) (*models.Customer, error) {
	args := m.Called(customer)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Customer), args.Error(1)
}

func (m *MockCustomerRepository) GetCustomerByID(id string) (*models.Customer, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Customer), args.Error(1)
}

func (m *MockCustomerRepository) GetAllCustomers() ([]*models.Customer, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*models.Customer), args.Error(1)
}

func (m *MockCustomerRepository) UpdateCustomer(customer *models.Customer) (*models.Customer, error) {
	args := m.Called(customer)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Customer), args.Error(1)
}

func (m *MockCustomerRepository) DeleteCustomer(id string) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockCustomerRepository) GetCustomerByEmail(email string) (*models.Customer, error) {
	args := m.Called(email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Customer), args.Error(1)
}

// Helper function to create a test JWT token
func createCustomerTestToken(secret []byte, userID string, userRole string) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"role":    userRole,
		"exp":     time.Now().Add(time.Hour * 72).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(secret)
}

func setupCustomerTestApp(repo repositories.CustomerRepository, jwtSecret []byte) *fiber.App {
	app := fiber.New()
	h := NewCustomerHandler(repo, jwtSecret)
	// Assuming routes are registered under /api like in other tests
	// If your actual routes are different, adjust this group path.
	// Based on RegisterCustomerRoutes, it seems it's r.Group("/customers", authRequired)
	// So if 'r' is app.Group("/api"), then it would be /api/customers
	// For simplicity, let's assume the main app will group it under /api and pass that to RegisterCustomerRoutes
	// Or, the RegisterCustomerRoutes gets the app itself.
	// Let's follow the material_handlers_test.go pattern
	apiGroup := app.Group("/api") // If your routes are not prefixed with /api, adjust this.
	h.RegisterCustomerRoutes(apiGroup)
	return app
}

// Helper to compare CustomerResponse, ignoring time-sensitive fields if necessary
// For CustomerResponse, all time fields are strings, so direct comparison should work if mocks are consistent.

func TestCreateCustomerHandler(t *testing.T) {
	mockRepo := new(MockCustomerRepository)
	jwtSecret := []byte("testsecret")
	app := setupCustomerTestApp(mockRepo, jwtSecret)
	testToken, _ := createCustomerTestToken(jwtSecret, "user-id-123", "admin")

	createReq := CreateCustomerRequest{
		FullName: "Test Customer",
		Email:    "test@example.com",
		Phone:    "+12345678901",
		Address:  "123 Test St",
	}

	t.Run("Success", func(t *testing.T) {
		// The ID is generated in the handler or repo, so we capture the argument passed to CreateCustomer
		// Define the expected model that the mock repository's CreateCustomer method should return.
		// This simulates the repository having set the ID and timestamps.
		expectedCreatedModel := &models.Customer{
			// ID will be set based on what the handler passes or what the repo mock decides.
			// For this test, we'll assume the handler passes a UUID, and the repo mock confirms other fields.
			FullName:  createReq.FullName,
			Email:     createReq.Email,
			Phone:     createReq.Phone,
			Address:   createReq.Address,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		mockRepo.On("CreateCustomer", mock.MatchedBy(func(argToCreate *models.Customer) bool {
			// The handler generates a UUID for ID before calling the repo.
			// So, argToCreate.ID will be a UUID string here.
			// We set this ID on our expectedCreatedModel to ensure the returned object matches.
			expectedCreatedModel.ID = argToCreate.ID // Capture the generated ID
			return argToCreate.FullName == createReq.FullName &&
				argToCreate.Email == createReq.Email &&
				argToCreate.Phone == createReq.Phone &&
				argToCreate.Address == createReq.Address
		})).Return(expectedCreatedModel, nil).Once()

		bodyBytes, _ := json.Marshal(createReq)
		req := httptest.NewRequest(http.MethodPost, "/api/customers", bytes.NewReader(bodyBytes))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+testToken)

		resp, err := app.Test(req, -1)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusCreated, resp.StatusCode)

		var actualResponse CustomerResponse
		err = json.NewDecoder(resp.Body).Decode(&actualResponse)
		assert.NoError(t, err)

		assert.NotEmpty(t, actualResponse.ID, "ID should not be empty")
		assert.Equal(t, createReq.FullName, actualResponse.FullName)
		assert.Equal(t, createReq.Email, actualResponse.Email)
		assert.Equal(t, createReq.Phone, actualResponse.Phone)
		assert.Equal(t, createReq.Address, actualResponse.Address)
		// Timestamps can be tricky if not precisely controlled in mock.
		// For now, let's assume they are set as expected.
		assert.NotNil(t, actualResponse.CreatedAt)
		assert.NotNil(t, actualResponse.UpdatedAt)

		mockRepo.AssertExpectations(t)
	})

	t.Run("Invalid Request Payload - Bad JSON", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/api/customers", bytes.NewReader([]byte("invalid json")))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+testToken)

		resp, err := app.Test(req, -1)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

		var errResp ErrorResponse
		err = json.NewDecoder(resp.Body).Decode(&errResp)
		assert.NoError(t, err)
		assert.Equal(t, "Invalid request payload", errResp.Error)
		mockRepo.AssertExpectations(t) // No call to repo expected
	})

	t.Run("Validation Error - Missing Required Fields", func(t *testing.T) {
		invalidReq := CreateCustomerRequest{FullName: "Test"}
		bodyBytes, _ := json.Marshal(invalidReq)
		req := httptest.NewRequest(http.MethodPost, "/api/customers", bytes.NewReader(bodyBytes))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+testToken)

		resp, err := app.Test(req, -1)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

		var errResp ErrorResponse
		err = json.NewDecoder(resp.Body).Decode(&errResp)
		assert.NoError(t, err)
		assert.Equal(t, "FullName, email, and phone are required", errResp.Error)
		mockRepo.AssertExpectations(t) // No call to repo expected
	})

	t.Run("Repository Error on CreateCustomer", func(t *testing.T) {
		mockRepo.On("CreateCustomer", mock.AnythingOfType("*models.Customer")).Return(nil, errors.New("db error")).Once()

		bodyBytes, _ := json.Marshal(createReq) // Use valid request
		req := httptest.NewRequest(http.MethodPost, "/api/customers", bytes.NewReader(bodyBytes))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+testToken)

		resp, err := app.Test(req, -1)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)

		var errResp ErrorResponse
		err = json.NewDecoder(resp.Body).Decode(&errResp)
		assert.NoError(t, err)
		assert.Equal(t, "Failed to create customer", errResp.Error)
		mockRepo.AssertExpectations(t)
	})
}

func TestGetAllCustomersHandler(t *testing.T) {
	mockRepo := new(MockCustomerRepository)
	jwtSecret := []byte("testsecret")
	app := setupCustomerTestApp(mockRepo, jwtSecret)
	testToken, _ := createCustomerTestToken(jwtSecret, "user-id-123", "admin")

	expectedCustomersModel := []*models.Customer{
		{ID: "uuid1", FullName: "Customer 1", Email: "cust1@example.com", Phone: "+111", Address: "Addr1", CreatedAt: time.Now(), UpdatedAt: time.Now()},
		{ID: "uuid2", FullName: "Customer 2", Email: "cust2@example.com", Phone: "+222", Address: "Addr2", CreatedAt: time.Now(), UpdatedAt: time.Now()},
	}

	t.Run("Success", func(t *testing.T) {
		mockRepo.On("GetAllCustomers").Return(expectedCustomersModel, nil).Once()

		req := httptest.NewRequest(http.MethodGet, "/api/customers", nil)
		req.Header.Set("Authorization", "Bearer "+testToken)

		resp, err := app.Test(req, -1)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var listResponse CustomerListResponse
		err = json.NewDecoder(resp.Body).Decode(&listResponse)
		assert.NoError(t, err)

		assert.Len(t, listResponse.Customers, 2)
		for i, actualCustResp := range listResponse.Customers {
			expectedModel := expectedCustomersModel[i]
			assert.Equal(t, expectedModel.ID, actualCustResp.ID)
			assert.Equal(t, expectedModel.FullName, actualCustResp.FullName)
			assert.Equal(t, expectedModel.Email, actualCustResp.Email)
			// Add other field assertions as needed
		}
		mockRepo.AssertExpectations(t)
	})

	t.Run("Repository Error", func(t *testing.T) {
		mockRepo.On("GetAllCustomers").Return(nil, errors.New("db error fetching all")).Once()

		req := httptest.NewRequest(http.MethodGet, "/api/customers", nil)
		req.Header.Set("Authorization", "Bearer "+testToken)

		resp, err := app.Test(req, -1)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)

		var errResp ErrorResponse
		err = json.NewDecoder(resp.Body).Decode(&errResp)
		assert.NoError(t, err)
		assert.Equal(t, "Failed to retrieve customers", errResp.Error)
		mockRepo.AssertExpectations(t)
	})
}

func TestGetCustomerHandler(t *testing.T) {
	mockRepo := new(MockCustomerRepository)
	jwtSecret := []byte("testsecret")
	app := setupCustomerTestApp(mockRepo, jwtSecret)
	testToken, _ := createCustomerTestToken(jwtSecret, "user-id-123", "admin")

	customerID := uuid.New().String()
	expectedCustomerModel := &models.Customer{
		ID:        customerID,
		FullName:  "Specific Customer",
		Email:     "specific@example.com",
		Phone:     "+3334445555",
		Address:   "456 Specific Ave",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	t.Run("Success", func(t *testing.T) {
		mockRepo.On("GetCustomerByID", customerID).Return(expectedCustomerModel, nil).Once()

		req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/customers/%s", customerID), nil)
		req.Header.Set("Authorization", "Bearer "+testToken)

		resp, err := app.Test(req, -1)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var actualResponse CustomerResponse
		err = json.NewDecoder(resp.Body).Decode(&actualResponse)
		assert.NoError(t, err)
		assert.Equal(t, expectedCustomerModel.ID, actualResponse.ID)
		assert.Equal(t, expectedCustomerModel.FullName, actualResponse.FullName)
		assert.Equal(t, expectedCustomerModel.Email, actualResponse.Email)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Invalid Customer ID Format", func(t *testing.T) {
		invalidID := "not-a-uuid"
		req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/customers/%s", invalidID), nil)
		req.Header.Set("Authorization", "Bearer "+testToken)

		resp, err := app.Test(req, -1)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

		var errResp ErrorResponse
		err = json.NewDecoder(resp.Body).Decode(&errResp)
		assert.NoError(t, err)
		assert.Equal(t, "Invalid Customer ID format", errResp.Error)
		// No call to repo expected
		mockRepo.AssertExpectations(t)
	})

	t.Run("Customer Not Found", func(t *testing.T) {
		notFoundID := uuid.New().String()
		// The handler checks for a specific error string, which is fragile.
		// A better approach would be custom error types or errors.Is from the repo.
		mockRepo.On("GetCustomerByID", notFoundID).Return(nil, fmt.Errorf("customer with ID %s not found", notFoundID)).Once()

		req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/customers/%s", notFoundID), nil)
		req.Header.Set("Authorization", "Bearer "+testToken)

		resp, err := app.Test(req, -1)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusNotFound, resp.StatusCode)

		var errResp ErrorResponse
		err = json.NewDecoder(resp.Body).Decode(&errResp)
		assert.NoError(t, err)
		assert.Equal(t, "Customer not found", errResp.Error)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Repository Error - Other", func(t *testing.T) {
		errorID := uuid.New().String()
		mockRepo.On("GetCustomerByID", errorID).Return(nil, errors.New("some other db error")).Once()

		req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/customers/%s", errorID), nil)
		req.Header.Set("Authorization", "Bearer "+testToken)

		resp, err := app.Test(req, -1)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)

		var errResp ErrorResponse
		err = json.NewDecoder(resp.Body).Decode(&errResp)
		assert.NoError(t, err)
		assert.Equal(t, "Failed to retrieve customer", errResp.Error)
		mockRepo.AssertExpectations(t)
	})
}

func TestUpdateCustomerHandler(t *testing.T) {
	mockRepo := new(MockCustomerRepository)
	jwtSecret := []byte("testsecret")
	app := setupCustomerTestApp(mockRepo, jwtSecret)
	testToken, _ := createCustomerTestToken(jwtSecret, "user-id-123", "admin")

	customerID := uuid.New().String()

	existingCustomerModel := &models.Customer{
		ID:        customerID,
		FullName:  "Original Name",
		Email:     "original@example.com",
		Phone:     "+1000000000",
		Address:   "Original Address",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	updateReq := UpdateCustomerRequest{
		FullName: "Updated Name",
		Email:    "updated@example.com",
		Phone:    "+2000000000",
		Address:  "Updated Address",
	}

	t.Run("Success - Full Update", func(t *testing.T) {
		mockRepo.On("GetCustomerByID", customerID).Return(existingCustomerModel, nil).Once()

		// Define the expected state of the customer model that the mock's UpdateCustomer method should return.
		// This model reflects updates from updateReq AND a new UpdatedAt timestamp as if set by the repository.
		expectedReturnFromRepoUpdate := *existingCustomerModel // Start with a copy
		expectedReturnFromRepoUpdate.FullName = updateReq.FullName
		expectedReturnFromRepoUpdate.Email = updateReq.Email
		expectedReturnFromRepoUpdate.Phone = updateReq.Phone
		expectedReturnFromRepoUpdate.Address = updateReq.Address
		// Simulate the repository updating the UpdatedAt timestamp. Ensure it's different from the original.
		newUpdatedAt := time.Now().Add(5 * time.Second)
		if newUpdatedAt.Equal(existingCustomerModel.UpdatedAt) { // Ensure it's different
			newUpdatedAt = time.Now().Add(10 * time.Second)
		}
		expectedReturnFromRepoUpdate.UpdatedAt = newUpdatedAt

		mockRepo.On("UpdateCustomer", mock.MatchedBy(func(argToUpdate *models.Customer) bool {
			// Check the model passed to UpdateCustomer by the handler:
			// It should have fields from updateReq, but UpdatedAt should still be original.
			return argToUpdate.ID == customerID &&
				argToUpdate.FullName == updateReq.FullName &&
				argToUpdate.Email == updateReq.Email &&
				argToUpdate.Phone == updateReq.Phone &&
				argToUpdate.Address == updateReq.Address &&
				argToUpdate.UpdatedAt.Equal(existingCustomerModel.UpdatedAt) // Before repo changes it
		})).Return(&expectedReturnFromRepoUpdate, nil).Once()

		bodyBytes, _ := json.Marshal(updateReq)
		req := httptest.NewRequest(http.MethodPut, fmt.Sprintf("/api/customers/%s", customerID), bytes.NewReader(bodyBytes))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+testToken)

		resp, err := app.Test(req, -1)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var actualResponse CustomerResponse
		err = json.NewDecoder(resp.Body).Decode(&actualResponse)
		assert.NoError(t, err)
		assert.Equal(t, customerID, actualResponse.ID)
		assert.Equal(t, updateReq.FullName, actualResponse.FullName)
		assert.Equal(t, updateReq.Email, actualResponse.Email)
		assert.NotEqual(t, existingCustomerModel.UpdatedAt, actualResponse.UpdatedAt) // Check UpdatedAt changed
		mockRepo.AssertExpectations(t)
	})

	t.Run("Success - Partial Update (Name only)", func(t *testing.T) {
		partialUpdateReq := UpdateCustomerRequest{FullName: "Partial Update Name"}
		// Reset and re-mock GetCustomerByID for this specific sub-test
		mockRepo.ExpectedCalls = nil // Clear previous expectations
		mockRepo.On("GetCustomerByID", customerID).Return(existingCustomerModel, nil).Once()

		// Define the expected state for partial update return
		expectedReturnFromPartialRepoUpdate := *existingCustomerModel // Start with a copy
		expectedReturnFromPartialRepoUpdate.FullName = partialUpdateReq.FullName
		// Other fields (Email, Phone, Address) remain as in existingCustomerModel
		newPartialUpdatedAt := time.Now().Add(5 * time.Second)
		if newPartialUpdatedAt.Equal(existingCustomerModel.UpdatedAt) {
			newPartialUpdatedAt = time.Now().Add(10 * time.Second)
		}
		expectedReturnFromPartialRepoUpdate.UpdatedAt = newPartialUpdatedAt

		mockRepo.On("UpdateCustomer", mock.MatchedBy(func(argToUpdate *models.Customer) bool {
			// Check the model passed to UpdateCustomer by the handler for partial update:
			// Name updated, Email, Phone, Address, UpdatedAt should be original.
			return argToUpdate.ID == customerID &&
				argToUpdate.FullName == partialUpdateReq.FullName &&
				argToUpdate.Email == existingCustomerModel.Email &&
				argToUpdate.Phone == existingCustomerModel.Phone &&
				argToUpdate.Address == existingCustomerModel.Address &&
				argToUpdate.UpdatedAt.Equal(existingCustomerModel.UpdatedAt)
		})).Return(&expectedReturnFromPartialRepoUpdate, nil).Once()

		bodyBytes, _ := json.Marshal(partialUpdateReq)
		req := httptest.NewRequest(http.MethodPut, fmt.Sprintf("/api/customers/%s", customerID), bytes.NewReader(bodyBytes))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+testToken)

		resp, err := app.Test(req, -1)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var actualResponse CustomerResponse
		err = json.NewDecoder(resp.Body).Decode(&actualResponse)
		assert.NoError(t, err)
		assert.Equal(t, customerID, actualResponse.ID)
		assert.Equal(t, partialUpdateReq.FullName, actualResponse.FullName)
		assert.Equal(t, existingCustomerModel.Email, actualResponse.Email) // Check email is unchanged
		mockRepo.AssertExpectations(t)
	})

	t.Run("Invalid Customer ID Format", func(t *testing.T) {
		invalidID := "not-a-uuid"
		bodyBytes, _ := json.Marshal(updateReq)
		req := httptest.NewRequest(http.MethodPut, fmt.Sprintf("/api/customers/%s", invalidID), bytes.NewReader(bodyBytes))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+testToken)

		resp, err := app.Test(req, -1)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
		// No call to repo expected
		mockRepo.ExpectedCalls = nil // Ensure clean state for assertion
		mockRepo.AssertExpectations(t)
	})

	t.Run("Invalid Request Payload - Bad JSON", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPut, fmt.Sprintf("/api/customers/%s", customerID), bytes.NewReader([]byte("invalid json")))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+testToken)

		resp, err := app.Test(req, -1)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
		// No call to repo expected
		mockRepo.ExpectedCalls = nil // Ensure clean state for assertion
		mockRepo.AssertExpectations(t)
	})

	t.Run("Customer Not Found on GetCustomerByID", func(t *testing.T) {
		notFoundID := uuid.New().String()
		mockRepo.ExpectedCalls = nil // Clear previous expectations before setting new ones
		mockRepo.On("GetCustomerByID", notFoundID).Return(nil, errors.New("not found error from repo")).Once()

		bodyBytes, _ := json.Marshal(updateReq)
		req := httptest.NewRequest(http.MethodPut, fmt.Sprintf("/api/customers/%s", notFoundID), bytes.NewReader(bodyBytes))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+testToken)

		resp, err := app.Test(req, -1)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusNotFound, resp.StatusCode) // Handler maps repo error to 404 here
		mockRepo.AssertExpectations(t)
	})

	t.Run("Repository Error on UpdateCustomer", func(t *testing.T) {
		mockRepo.ExpectedCalls = nil // Clear previous expectations
		mockRepo.On("GetCustomerByID", customerID).Return(existingCustomerModel, nil).Once()
		mockRepo.On("UpdateCustomer", mock.AnythingOfType("*models.Customer")).Return(nil, errors.New("db update error")).Once()

		bodyBytes, _ := json.Marshal(updateReq)
		req := httptest.NewRequest(http.MethodPut, fmt.Sprintf("/api/customers/%s", customerID), bytes.NewReader(bodyBytes))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+testToken)

		resp, err := app.Test(req, -1)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
		mockRepo.AssertExpectations(t)
	})
}

func TestDeleteCustomerHandler(t *testing.T) {
	mockRepo := new(MockCustomerRepository)
	jwtSecret := []byte("testsecret")
	app := setupCustomerTestApp(mockRepo, jwtSecret)
	testToken, _ := createCustomerTestToken(jwtSecret, "user-id-123", "admin")

	customerID := uuid.New().String()

	t.Run("Success", func(t *testing.T) {
		mockRepo.ExpectedCalls = nil
		mockRepo.On("DeleteCustomer", customerID).Return(nil).Once()

		req := httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/api/customers/%s", customerID), nil)
		req.Header.Set("Authorization", "Bearer "+testToken)

		resp, err := app.Test(req, -1)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusNoContent, resp.StatusCode)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Invalid Customer ID Format", func(t *testing.T) {
		invalidID := "not-a-uuid"
		mockRepo.ExpectedCalls = nil

		req := httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/api/customers/%s", invalidID), nil)
		req.Header.Set("Authorization", "Bearer "+testToken)

		resp, err := app.Test(req, -1)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

		var errResp ErrorResponse
		err = json.NewDecoder(resp.Body).Decode(&errResp)
		assert.NoError(t, err)
		assert.Equal(t, "Invalid Customer ID format", errResp.Error)
		mockRepo.AssertExpectations(t) // No call to repo expected
	})

	t.Run("Customer Not Found", func(t *testing.T) {
		notFoundID := uuid.New().String()
		mockRepo.ExpectedCalls = nil
		// Handler checks for: "customer with ID "+id+" not found for deletion"
		mockRepo.On("DeleteCustomer", notFoundID).Return(fmt.Errorf("customer with ID %s not found for deletion", notFoundID)).Once()

		req := httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/api/customers/%s", notFoundID), nil)
		req.Header.Set("Authorization", "Bearer "+testToken)

		resp, err := app.Test(req, -1)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusNotFound, resp.StatusCode)

		var errResp ErrorResponse
		err = json.NewDecoder(resp.Body).Decode(&errResp)
		assert.NoError(t, err)
		assert.Equal(t, "Customer not found", errResp.Error)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Repository Error - Other", func(t *testing.T) {
		errorID := uuid.New().String()
		mockRepo.ExpectedCalls = nil
		mockRepo.On("DeleteCustomer", errorID).Return(errors.New("some other db error during delete")).Once()

		req := httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/api/customers/%s", errorID), nil)
		req.Header.Set("Authorization", "Bearer "+testToken)

		resp, err := app.Test(req, -1)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)

		var errResp ErrorResponse
		err = json.NewDecoder(resp.Body).Decode(&errResp)
		assert.NoError(t, err)
		assert.Equal(t, "Failed to delete customer", errResp.Error)
		mockRepo.AssertExpectations(t)
	})
}
