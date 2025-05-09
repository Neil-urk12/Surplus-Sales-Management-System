package handlers

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"oop/internal/models" // Assuming models are in this path
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockSaleRepository is a mock implementation of SaleRepository
type MockSaleRepository struct {
	mock.Mock
}

func (m *MockSaleRepository) GetAll(filters map[string]interface{}) ([]models.Sale, error) {
	args := m.Called(filters)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.Sale), args.Error(1)
}

func (m *MockSaleRepository) GetByID(id string) (*models.Sale, error) {
	args := m.Called(id)
	// Handle the case where Get(0) might be nil for a *models.Sale
	if ret := args.Get(0); ret != nil {
		return ret.(*models.Sale), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockSaleRepository) GetSaleItems(saleID string) ([]models.SaleItem, error) {
	args := m.Called(saleID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.SaleItem), args.Error(1)
}

func (m *MockSaleRepository) GetCustomerSales(customerID string) ([]models.Sale, error) {
	args := m.Called(customerID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.Sale), args.Error(1)
}

func (m *MockSaleRepository) Create(sale *models.Sale) (string, error) {
	args := m.Called(sale)
	return args.String(0), args.Error(1)
}

func (m *MockSaleRepository) CreateSaleItem(item *models.SaleItem) (string, error) {
	args := m.Called(item)
	return args.String(0), args.Error(1)
}

func (m *MockSaleRepository) Update(sale *models.Sale) error {
	args := m.Called(sale)
	return args.Error(0)
}

func (m *MockSaleRepository) Delete(id string) error {
	args := m.Called(id)
	return args.Error(0)
}

// Helper function to create a test Fiber app and SaleHandlers
// It also includes a mock middleware to simulate authentication
func setupSaleTestApp(mockRepo *MockSaleRepository, t *testing.T) (*fiber.App, *SaleHandlers) {
	app := fiber.New(fiber.Config{
		ErrorHandler: func(c *fiber.Ctx, err error) error { // Default error handler for Fiber
			code := fiber.StatusInternalServerError
			if e, ok := err.(*fiber.Error); ok {
				code = e.Code
			}
			return c.Status(code).JSON(fiber.Map{"error": err.Error()})
		},
	})
	jwtSecret := []byte("testsecret") // Dummy secret for tests
	handlers := NewSaleHandlers(mockRepo, nil, nil, nil, jwtSecret)

	// Mock middleware to set user_id for authenticated routes
	// This simulates what middleware.JWTMiddleware would do after successful auth
	authMiddleware := func(c *fiber.Ctx) error {
		c.Locals("user_id", "test_user_id") // Set a dummy user ID
		return c.Next()
	}

	// Register routes (mirroring RegisterSaleRoutes but with mock middleware)
	// Group routes under '/api/sales' (assuming an /api prefix for tests)
	salesGroup := app.Group("/api/sales", authMiddleware)
	salesGroup.Get("/", handlers.GetSalesHandler)
	salesGroup.Get("/:id", handlers.GetSaleByIDHandler)
	salesGroup.Get("/:id/items", handlers.GetSaleItemsHandler)
	salesGroup.Post("/", handlers.CreateSaleHandler)
	salesGroup.Put("/:id", handlers.UpdateSaleHandler)
	salesGroup.Delete("/:id", handlers.DeleteSaleHandler)

	app.Get("/api/customers/:id/sales", authMiddleware, handlers.GetCustomerSalesHandler)
	app.Post("/api/cabs/:id/sell", authMiddleware, handlers.SellCabHandler)

	return app, handlers
}

// TestGetSalesHandler
func TestGetSalesHandler(t *testing.T) {
	t.Parallel()
	mockRepo := new(MockSaleRepository)
	app, _ := setupSaleTestApp(mockRepo, t)

	t.Run("success - no filters", func(t *testing.T) {
		expectedSales := []models.Sale{{ID: "1", CustomerID: "cust1", SaleDate: "2023-01-01"}}
		mockRepo.On("GetAll", map[string]interface{}{}).Return(expectedSales, nil).Once()

		req := httptest.NewRequest(http.MethodGet, "/api/sales", nil)
		resp, err := app.Test(req, -1) // -1 for no timeout
		assert.NoError(t, err)

		assert.Equal(t, http.StatusOK, resp.StatusCode)
		var sales []models.Sale
		err = json.NewDecoder(resp.Body).Decode(&sales)
		assert.NoError(t, err)
		assert.Equal(t, expectedSales, sales)
		mockRepo.AssertExpectations(t)
	})

	t.Run("success - with filters", func(t *testing.T) {
		filters := map[string]interface{}{"customer_id": "cust1", "date_from": "2023-01-01"}
		expectedSales := []models.Sale{{ID: "1", CustomerID: "cust1", SaleDate: "2023-01-01"}}
		mockRepo.On("GetAll", filters).Return(expectedSales, nil).Once()

		req := httptest.NewRequest(http.MethodGet, "/api/sales?customer_id=cust1&date_from=2023-01-01", nil)
		resp, err := app.Test(req, -1)
		assert.NoError(t, err)

		assert.Equal(t, http.StatusOK, resp.StatusCode)
		var sales []models.Sale
		err = json.NewDecoder(resp.Body).Decode(&sales)
		assert.NoError(t, err)
		assert.Equal(t, expectedSales, sales)
		mockRepo.AssertExpectations(t)
	})

	t.Run("failure - repository error", func(t *testing.T) {
		mockRepo.On("GetAll", map[string]interface{}{}).Return(nil, errors.New("db error")).Once()

		req := httptest.NewRequest(http.MethodGet, "/api/sales", nil)
		resp, err := app.Test(req, -1)
		assert.NoError(t, err)

		assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
		var errResp map[string]interface{}
		err = json.NewDecoder(resp.Body).Decode(&errResp)
		assert.NoError(t, err)
		assert.Equal(t, "Failed to retrieve sales", errResp["error"])
		mockRepo.AssertExpectations(t)
	})
}

// TestGetSaleByIDHandler
func TestGetSaleByIDHandler(t *testing.T) {
	t.Parallel()
	mockRepo := new(MockSaleRepository)
	app, _ := setupSaleTestApp(mockRepo, t)

	t.Run("success", func(t *testing.T) {
		saleID := "sale123"
		expectedSale := &models.Sale{ID: saleID, CustomerID: "cust1", SaleDate: "2023-01-01"}
		mockRepo.On("GetByID", saleID).Return(expectedSale, nil).Once()

		req := httptest.NewRequest(http.MethodGet, "/api/sales/"+saleID, nil)
		resp, err := app.Test(req, -1)
		assert.NoError(t, err)

		assert.Equal(t, http.StatusOK, resp.StatusCode)
		var sale models.Sale
		err = json.NewDecoder(resp.Body).Decode(&sale)
		assert.NoError(t, err)
		assert.Equal(t, *expectedSale, sale)
		mockRepo.AssertExpectations(t)
	})

	t.Run("not found", func(t *testing.T) {
		saleID := "nonexistent"
		mockRepo.On("GetByID", saleID).Return(nil, nil).Once() // Repo returns nil, no error for not found

		req := httptest.NewRequest(http.MethodGet, "/api/sales/"+saleID, nil)
		resp, err := app.Test(req, -1)
		assert.NoError(t, err)

		assert.Equal(t, http.StatusNotFound, resp.StatusCode)
		var errResp map[string]interface{}
		err = json.NewDecoder(resp.Body).Decode(&errResp)
		assert.NoError(t, err)
		assert.Equal(t, "Sale not found", errResp["error"])
		mockRepo.AssertExpectations(t)
	})

	t.Run("repository error", func(t *testing.T) {
		saleID := "errorID"
		mockRepo.On("GetByID", saleID).Return(nil, errors.New("db error")).Once()

		req := httptest.NewRequest(http.MethodGet, "/api/sales/"+saleID, nil)
		resp, err := app.Test(req, -1)
		assert.NoError(t, err)

		assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
		var errResp map[string]interface{}
		err = json.NewDecoder(resp.Body).Decode(&errResp)
		assert.NoError(t, err)
		assert.Equal(t, "Failed to retrieve sale", errResp["error"])
		mockRepo.AssertExpectations(t)
	})
}

// TestGetSaleItemsHandler
func TestGetSaleItemsHandler(t *testing.T) {
	t.Parallel()
	mockRepo := new(MockSaleRepository)
	app, _ := setupSaleTestApp(mockRepo, t)
	saleID := "sale123"

	t.Run("success", func(t *testing.T) {
		existingSale := &models.Sale{ID: saleID, SaleDate: "2023-01-01"}
		expectedItems := []models.SaleItem{{ID: "item1", SaleID: saleID, MultiCabID: "prod1", Quantity: 1, UnitPrice: 10.0}}

		mockRepo.On("GetByID", saleID).Return(existingSale, nil).Once()
		mockRepo.On("GetSaleItems", saleID).Return(expectedItems, nil).Once()

		req := httptest.NewRequest(http.MethodGet, "/api/sales/"+saleID+"/items", nil)
		resp, err := app.Test(req, -1)
		assert.NoError(t, err)

		assert.Equal(t, http.StatusOK, resp.StatusCode)
		var items []models.SaleItem
		err = json.NewDecoder(resp.Body).Decode(&items)
		assert.NoError(t, err)
		assert.Equal(t, expectedItems, items)
		mockRepo.AssertExpectations(t)
	})

	t.Run("sale not found", func(t *testing.T) {
		mockRepo.On("GetByID", saleID).Return(nil, nil).Once() // Sale not found

		req := httptest.NewRequest(http.MethodGet, "/api/sales/"+saleID+"/items", nil)
		resp, err := app.Test(req, -1)
		assert.NoError(t, err)

		assert.Equal(t, http.StatusNotFound, resp.StatusCode)
		var errResp map[string]interface{}
		err = json.NewDecoder(resp.Body).Decode(&errResp)
		assert.NoError(t, err)
		assert.Equal(t, "Sale not found", errResp["error"])
		mockRepo.AssertExpectations(t)
	})

	t.Run("error getting sale for check", func(t *testing.T) {
		mockRepo.On("GetByID", saleID).Return(nil, errors.New("db error checking sale")).Once()

		req := httptest.NewRequest(http.MethodGet, "/api/sales/"+saleID+"/items", nil)
		resp, err := app.Test(req, -1)
		assert.NoError(t, err)

		assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
		var errResp map[string]interface{}
		err = json.NewDecoder(resp.Body).Decode(&errResp)
		assert.NoError(t, err)
		assert.Equal(t, "Failed to retrieve sale", errResp["error"])
		mockRepo.AssertExpectations(t)
	})

	t.Run("error getting sale items", func(t *testing.T) {
		existingSale := &models.Sale{ID: saleID, SaleDate: "2023-01-01"}
		mockRepo.On("GetByID", saleID).Return(existingSale, nil).Once()
		mockRepo.On("GetSaleItems", saleID).Return(nil, errors.New("db error getting items")).Once()

		req := httptest.NewRequest(http.MethodGet, "/api/sales/"+saleID+"/items", nil)
		resp, err := app.Test(req, -1)
		assert.NoError(t, err)

		assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
		var errResp map[string]interface{}
		err = json.NewDecoder(resp.Body).Decode(&errResp)
		assert.NoError(t, err)
		assert.Equal(t, "Failed to retrieve sale items", errResp["error"])
		mockRepo.AssertExpectations(t)
	})
}

// TestCreateSaleHandler
func TestCreateSaleHandler(t *testing.T) {
	t.Parallel()
	mockRepo := new(MockSaleRepository)
	app, _ := setupSaleTestApp(mockRepo, t)

	t.Run("success", func(t *testing.T) {
		saleInput := models.Sale{
			CustomerID: "cust1",
			SoldBy:     "user1",
			SaleDate:   "2023-01-01",
			TotalPrice: 100.0,
		}
		expectedSaleID := "newSaleID"

		mockRepo.On("Create", mock.MatchedBy(func(s *models.Sale) bool {
			return s.CustomerID == saleInput.CustomerID &&
				s.SoldBy == saleInput.SoldBy &&
				s.SaleDate == saleInput.SaleDate &&
				s.TotalPrice == saleInput.TotalPrice &&
				!s.CreatedAt.IsZero() &&
				!s.UpdatedAt.IsZero()
		})).Return(expectedSaleID, nil).Once()

		payload, _ := json.Marshal(saleInput)
		req := httptest.NewRequest(http.MethodPost, "/api/sales", bytes.NewBuffer(payload))
		req.Header.Set("Content-Type", "application/json")
		resp, err := app.Test(req, -1)
		assert.NoError(t, err)

		assert.Equal(t, http.StatusCreated, resp.StatusCode)
		var createdSale models.Sale
		err = json.NewDecoder(resp.Body).Decode(&createdSale)
		assert.NoError(t, err)
		assert.Equal(t, expectedSaleID, createdSale.ID)
		assert.Equal(t, saleInput.CustomerID, createdSale.CustomerID)
		assert.False(t, createdSale.CreatedAt.IsZero())
		assert.False(t, createdSale.UpdatedAt.IsZero())
		mockRepo.AssertExpectations(t)
	})

	t.Run("invalid payload - bad json", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/api/sales", strings.NewReader("{bad json"))
		req.Header.Set("Content-Type", "application/json")
		resp, err := app.Test(req, -1)
		assert.NoError(t, err)

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
		var errResp map[string]interface{}
		err = json.NewDecoder(resp.Body).Decode(&errResp)
		assert.NoError(t, err)
		assert.Equal(t, "Invalid request payload", errResp["error"])
	})

	t.Run("invalid payload - missing fields", func(t *testing.T) {
		saleInput := models.Sale{CustomerID: "cust1"} // Missing SoldBy, SaleDate
		payload, _ := json.Marshal(saleInput)
		req := httptest.NewRequest(http.MethodPost, "/api/sales", bytes.NewBuffer(payload))
		req.Header.Set("Content-Type", "application/json")
		resp, err := app.Test(req, -1)
		assert.NoError(t, err)

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
		var errResp map[string]interface{}
		err = json.NewDecoder(resp.Body).Decode(&errResp)
		assert.NoError(t, err)
		assert.Equal(t, "Missing required sale fields", errResp["error"])
	})

	t.Run("repository error", func(t *testing.T) {
		saleInput := models.Sale{
			CustomerID: "cust1",
			SoldBy:     "user1",
			SaleDate:   "2023-01-01",
			TotalPrice: 100.0,
		}
		mockRepo.On("Create", mock.AnythingOfType("*models.Sale")).Return("", errors.New("db error")).Once()

		payload, _ := json.Marshal(saleInput)
		req := httptest.NewRequest(http.MethodPost, "/api/sales", bytes.NewBuffer(payload))
		req.Header.Set("Content-Type", "application/json")
		resp, err := app.Test(req, -1)
		assert.NoError(t, err)

		assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
		var errResp map[string]interface{}
		err = json.NewDecoder(resp.Body).Decode(&errResp)
		assert.NoError(t, err)
		assert.Equal(t, "Failed to create sale", errResp["error"])
		mockRepo.AssertExpectations(t)
	})
}

// TestUpdateSaleHandler
func TestUpdateSaleHandler(t *testing.T) {
	t.Parallel()
	mockRepo := new(MockSaleRepository)
	app, _ := setupSaleTestApp(mockRepo, t)
	saleID := "saleToUpdate"
	originalCreatedAt := time.Now().Add(-time.Hour).Truncate(time.Second) // Truncate for comparison

	t.Run("success", func(t *testing.T) {
		existingSale := &models.Sale{ID: saleID, CustomerID: "oldCust", SoldBy: "oldUser", SaleDate: "2023-01-01", CreatedAt: originalCreatedAt, UpdatedAt: originalCreatedAt}
		updatePayload := models.Sale{
			CustomerID: "newCust",
			SoldBy:     "newUser",
			SaleDate:   "2023-02-02",
			TotalPrice: 200.0,
		}

		mockRepo.On("GetByID", saleID).Return(existingSale, nil).Once()
		mockRepo.On("Update", mock.MatchedBy(func(s *models.Sale) bool {
			return s.ID == saleID &&
				s.CustomerID == updatePayload.CustomerID &&
				s.SoldBy == updatePayload.SoldBy &&
				s.SaleDate == updatePayload.SaleDate &&
				s.TotalPrice == updatePayload.TotalPrice &&
				s.CreatedAt.Equal(originalCreatedAt) &&
				!s.UpdatedAt.Equal(originalCreatedAt) // UpdatedAt should be new
		})).Return(nil).Once()

		payload, _ := json.Marshal(updatePayload)
		req := httptest.NewRequest(http.MethodPut, "/api/sales/"+saleID, bytes.NewBuffer(payload))
		req.Header.Set("Content-Type", "application/json")
		resp, err := app.Test(req, -1)
		assert.NoError(t, err)

		assert.Equal(t, http.StatusOK, resp.StatusCode)
		var updatedSale models.Sale
		err = json.NewDecoder(resp.Body).Decode(&updatedSale)
		assert.NoError(t, err)
		assert.Equal(t, saleID, updatedSale.ID)
		assert.Equal(t, updatePayload.CustomerID, updatedSale.CustomerID)
		assert.True(t, updatedSale.CreatedAt.Equal(originalCreatedAt))
		assert.False(t, updatedSale.UpdatedAt.IsZero())
		assert.False(t, updatedSale.UpdatedAt.Equal(originalCreatedAt))
		mockRepo.AssertExpectations(t)
	})

	t.Run("sale not found", func(t *testing.T) {
		mockRepo.On("GetByID", saleID).Return(nil, nil).Once()

		payload, _ := json.Marshal(models.Sale{CustomerID: "cust", SoldBy: "user", SaleDate: "date"})
		req := httptest.NewRequest(http.MethodPut, "/api/sales/"+saleID, bytes.NewBuffer(payload))
		req.Header.Set("Content-Type", "application/json")
		resp, err := app.Test(req, -1)
		assert.NoError(t, err)

		assert.Equal(t, http.StatusNotFound, resp.StatusCode)
		mockRepo.AssertExpectations(t)
	})

	t.Run("error getting sale for check", func(t *testing.T) {
		mockRepo.On("GetByID", saleID).Return(nil, errors.New("db error")).Once()

		payload, _ := json.Marshal(models.Sale{CustomerID: "cust", SoldBy: "user", SaleDate: "date"})
		req := httptest.NewRequest(http.MethodPut, "/api/sales/"+saleID, bytes.NewBuffer(payload))
		req.Header.Set("Content-Type", "application/json")
		resp, err := app.Test(req, -1)
		assert.NoError(t, err)

		assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
		mockRepo.AssertExpectations(t)
	})

	t.Run("invalid payload - missing fields", func(t *testing.T) {
		existingSale := &models.Sale{ID: saleID, CustomerID: "oldCust", CreatedAt: originalCreatedAt}
		mockRepo.On("GetByID", saleID).Return(existingSale, nil).Once()

		updatePayload := models.Sale{CustomerID: "newCust"} // Missing SoldBy, SaleDate
		payload, _ := json.Marshal(updatePayload)
		req := httptest.NewRequest(http.MethodPut, "/api/sales/"+saleID, bytes.NewBuffer(payload))
		req.Header.Set("Content-Type", "application/json")
		resp, err := app.Test(req, -1)
		assert.NoError(t, err)

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
		var errResp map[string]interface{}
		err = json.NewDecoder(resp.Body).Decode(&errResp)
		assert.NoError(t, err)
		assert.Equal(t, "Missing required sale fields", errResp["error"])
		mockRepo.AssertExpectations(t)
	})

	t.Run("repository error on update", func(t *testing.T) {
		existingSale := &models.Sale{ID: saleID, CustomerID: "oldCust", CreatedAt: originalCreatedAt}
		updatePayload := models.Sale{
			CustomerID: "newCust",
			SoldBy:     "newUser",
			SaleDate:   "2023-02-02",
			TotalPrice: 200.0,
		}
		mockRepo.On("GetByID", saleID).Return(existingSale, nil).Once()
		mockRepo.On("Update", mock.AnythingOfType("*models.Sale")).Return(errors.New("db update error")).Once()

		payload, _ := json.Marshal(updatePayload)
		req := httptest.NewRequest(http.MethodPut, "/api/sales/"+saleID, bytes.NewBuffer(payload))
		req.Header.Set("Content-Type", "application/json")
		resp, err := app.Test(req, -1)
		assert.NoError(t, err)

		assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
		var errResp map[string]interface{}
		err = json.NewDecoder(resp.Body).Decode(&errResp)
		assert.NoError(t, err)
		assert.Equal(t, "Failed to update sale", errResp["error"])
		mockRepo.AssertExpectations(t)
	})
}

// TestDeleteSaleHandler
func TestDeleteSaleHandler(t *testing.T) {
	t.Parallel()
	mockRepo := new(MockSaleRepository)
	app, _ := setupSaleTestApp(mockRepo, t)
	saleID := "saleToDelete"

	t.Run("success", func(t *testing.T) {
		existingSale := &models.Sale{ID: saleID, SaleDate: "2023-01-01"}
		mockRepo.On("GetByID", saleID).Return(existingSale, nil).Once()
		mockRepo.On("Delete", saleID).Return(nil).Once()

		req := httptest.NewRequest(http.MethodDelete, "/api/sales/"+saleID, nil)
		resp, err := app.Test(req, -1)
		assert.NoError(t, err)

		assert.Equal(t, http.StatusNoContent, resp.StatusCode)
		mockRepo.AssertExpectations(t)
	})

	t.Run("sale not found", func(t *testing.T) {
		mockRepo.On("GetByID", saleID).Return(nil, nil).Once()

		req := httptest.NewRequest(http.MethodDelete, "/api/sales/"+saleID, nil)
		resp, err := app.Test(req, -1)
		assert.NoError(t, err)

		assert.Equal(t, http.StatusNotFound, resp.StatusCode)
		mockRepo.AssertExpectations(t)
	})

	t.Run("error getting sale for check", func(t *testing.T) {
		mockRepo.On("GetByID", saleID).Return(nil, errors.New("db error")).Once()

		req := httptest.NewRequest(http.MethodDelete, "/api/sales/"+saleID, nil)
		resp, err := app.Test(req, -1)
		assert.NoError(t, err)

		assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
		mockRepo.AssertExpectations(t)
	})

	t.Run("repository error on delete", func(t *testing.T) {
		existingSale := &models.Sale{ID: saleID, SaleDate: "2023-01-01"}
		mockRepo.On("GetByID", saleID).Return(existingSale, nil).Once()
		mockRepo.On("Delete", saleID).Return(errors.New("db delete error")).Once()

		req := httptest.NewRequest(http.MethodDelete, "/api/sales/"+saleID, nil)
		resp, err := app.Test(req, -1)
		assert.NoError(t, err)

		assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
		var errResp map[string]interface{}
		err = json.NewDecoder(resp.Body).Decode(&errResp)
		assert.NoError(t, err)
		assert.Equal(t, "Failed to delete sale", errResp["error"])
		mockRepo.AssertExpectations(t)
	})
}

// TestSellCabHandler
func TestSellCabHandler(t *testing.T) {
	t.Parallel()
	mockRepo := new(MockSaleRepository)
	app, _ := setupSaleTestApp(mockRepo, t)
	cabID := 123

	t.Run("success", func(t *testing.T) {
		salePayload := models.CabSalePayload{
			CustomerID: "cust123",
			Quantity:   1,
			Accessories: []models.AccessoryForSale{
				{ID: 1, Name: "acc1", Price: 10.0, Quantity: 1, UnitPrice: 10.0},
				{ID: 2, Name: "acc2", Price: 20.0, Quantity: 1, UnitPrice: 20.0},
			},
		}
		expectedSaleID := "newSaleFromCab"
		expectedSaleDate := time.Now().Format("2006-01-02")

		mockRepo.On("Create", mock.MatchedBy(func(s *models.Sale) bool {
			return s.CustomerID == salePayload.CustomerID &&
				s.SoldBy == "test_user_id" && // From mock middleware
				s.SaleDate == expectedSaleDate &&
				s.TotalPrice == 0 && // Current implementation sets totalPrice to 0
				!s.CreatedAt.IsZero() &&
				!s.UpdatedAt.IsZero()
		})).Return(expectedSaleID, nil).Once()

		// TODO: Add mocks for CabRepo.GetByID, AccRepo.GetByID, and Repo.CreateSaleItem
		// if/when those parts of SellCabHandler are fully implemented.

		payload, _ := json.Marshal(salePayload)
		req := httptest.NewRequest(http.MethodPost, "/api/cabs/"+strconv.Itoa(cabID)+"/sell", bytes.NewBuffer(payload))
		req.Header.Set("Content-Type", "application/json")
		resp, err := app.Test(req, -1)
		assert.NoError(t, err)

		assert.Equal(t, http.StatusCreated, resp.StatusCode)
		var respBody map[string]interface{}
		err = json.NewDecoder(resp.Body).Decode(&respBody)
		assert.NoError(t, err)

		assert.True(t, respBody["success"].(bool))
		assert.Equal(t, "Cab sold successfully", respBody["message"])
		assert.Equal(t, float64(cabID), respBody["cabId"])
		assert.Equal(t, salePayload.CustomerID, respBody["customerId"])
		assert.Equal(t, float64(salePayload.Quantity), respBody["quantity"])
		assert.Equal(t, expectedSaleID, respBody["saleId"])
		assert.Equal(t, expectedSaleDate, respBody["saleDate"])

		// Check accessories are returned in the response
		respAccessories, ok := respBody["accessories"].([]interface{})
		assert.True(t, ok)
		assert.Equal(t, len(salePayload.Accessories), len(respAccessories))

		// Verify each accessory has the expected structure
		for i, acc := range respAccessories {
			accMap, ok := acc.(map[string]interface{})
			assert.True(t, ok)
			assert.Equal(t, float64(salePayload.Accessories[i].ID), accMap["id"])
			assert.Equal(t, salePayload.Accessories[i].Name, accMap["name"])
		}

		mockRepo.AssertExpectations(t)
	})

	t.Run("invalid cab ID format", func(t *testing.T) {
		salePayload := models.CabSalePayload{CustomerID: "cust123", Quantity: 1}
		payload, _ := json.Marshal(salePayload)
		req := httptest.NewRequest(http.MethodPost, "/api/cabs/invalidID/sell", bytes.NewBuffer(payload))
		req.Header.Set("Content-Type", "application/json")
		resp, err := app.Test(req, -1)
		assert.NoError(t, err)

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
		var errResp map[string]interface{}
		err = json.NewDecoder(resp.Body).Decode(&errResp)
		assert.NoError(t, err)
		assert.Equal(t, "Invalid cab ID format", errResp["error"])
	})

	t.Run("invalid request payload - missing fields", func(t *testing.T) {
		salePayload := models.CabSalePayload{Quantity: 1} // Missing CustomerID
		payload, _ := json.Marshal(salePayload)
		req := httptest.NewRequest(http.MethodPost, "/api/cabs/"+strconv.Itoa(cabID)+"/sell", bytes.NewBuffer(payload))
		req.Header.Set("Content-Type", "application/json")
		resp, err := app.Test(req, -1)
		assert.NoError(t, err)

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
		var errResp map[string]interface{}
		err = json.NewDecoder(resp.Body).Decode(&errResp)
		assert.NoError(t, err)
		assert.Equal(t, "Missing required sale fields or invalid quantity", errResp["error"])
	})

	t.Run("invalid request payload - zero quantity", func(t *testing.T) {
		salePayload := models.CabSalePayload{CustomerID: "cust123", Quantity: 0}
		payload, _ := json.Marshal(salePayload)
		req := httptest.NewRequest(http.MethodPost, "/api/cabs/"+strconv.Itoa(cabID)+"/sell", bytes.NewBuffer(payload))
		req.Header.Set("Content-Type", "application/json")
		resp, err := app.Test(req, -1)
		assert.NoError(t, err)

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
		var errResp map[string]interface{}
		err = json.NewDecoder(resp.Body).Decode(&errResp)
		assert.NoError(t, err)
		assert.Equal(t, "Missing required sale fields or invalid quantity", errResp["error"])
	})

	t.Run("repository error on sale create", func(t *testing.T) {
		salePayload := models.CabSalePayload{CustomerID: "cust123", Quantity: 1}
		mockRepo.On("Create", mock.AnythingOfType("*models.Sale")).Return("", errors.New("db error creating sale")).Once()

		payload, _ := json.Marshal(salePayload)
		req := httptest.NewRequest(http.MethodPost, "/api/cabs/"+strconv.Itoa(cabID)+"/sell", bytes.NewBuffer(payload))
		req.Header.Set("Content-Type", "application/json")
		resp, err := app.Test(req, -1)
		assert.NoError(t, err)

		assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
		var errResp map[string]interface{}
		err = json.NewDecoder(resp.Body).Decode(&errResp)
		assert.NoError(t, err)
		assert.Equal(t, "Failed to create sale", errResp["error"])
		mockRepo.AssertExpectations(t)
	})
}

// TestGetCustomerSalesHandler
func TestGetCustomerSalesHandler(t *testing.T) {
	t.Parallel()
	mockRepo := new(MockSaleRepository)
	app, _ := setupSaleTestApp(mockRepo, t)
	customerID := "cust789"

	t.Run("success - sales found", func(t *testing.T) {
		expectedSales := []models.Sale{
			{ID: "sale1", CustomerID: customerID, SaleDate: "2023-01-01"},
			{ID: "sale2", CustomerID: customerID, SaleDate: "2023-01-02"},
		}
		mockRepo.On("GetCustomerSales", customerID).Return(expectedSales, nil).Once()

		req := httptest.NewRequest(http.MethodGet, "/api/customers/"+customerID+"/sales", nil)
		resp, err := app.Test(req, -1)
		assert.NoError(t, err)

		assert.Equal(t, http.StatusOK, resp.StatusCode)
		var sales []models.Sale
		err = json.NewDecoder(resp.Body).Decode(&sales)
		assert.NoError(t, err)
		assert.Equal(t, expectedSales, sales)
		mockRepo.AssertExpectations(t)
	})

	t.Run("success - no sales found", func(t *testing.T) {
		mockRepo.On("GetCustomerSales", customerID).Return([]models.Sale{}, nil).Once()

		req := httptest.NewRequest(http.MethodGet, "/api/customers/"+customerID+"/sales", nil)
		resp, err := app.Test(req, -1)
		assert.NoError(t, err)

		assert.Equal(t, http.StatusOK, resp.StatusCode)
		var sales []models.Sale
		err = json.NewDecoder(resp.Body).Decode(&sales)
		assert.NoError(t, err)
		assert.Empty(t, sales)
		mockRepo.AssertExpectations(t)
	})

	t.Run("repository error", func(t *testing.T) {
		mockRepo.On("GetCustomerSales", customerID).Return(nil, errors.New("db error")).Once()

		req := httptest.NewRequest(http.MethodGet, "/api/customers/"+customerID+"/sales", nil)
		resp, err := app.Test(req, -1)
		assert.NoError(t, err)

		assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
		var errResp map[string]interface{}
		err = json.NewDecoder(resp.Body).Decode(&errResp)
		assert.NoError(t, err)
		assert.Equal(t, "Failed to retrieve customer sales", errResp["error"])
		mockRepo.AssertExpectations(t)
	})
}
