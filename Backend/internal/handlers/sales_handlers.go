package handlers

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"oop/internal/middleware"
	"oop/internal/models"
	"oop/internal/repositories"

	"github.com/gofiber/fiber/v2"
)

// SaleRepository defines the interface for sales data operations
type SaleRepository interface {
	// GetAll retrieves all sales, with optional filtering
	GetAll(filters map[string]interface{}) ([]models.Sale, error)

	// GetByID retrieves a sale by its ID
	GetByID(id string) (*models.Sale, error)

	// GetSaleItems retrieves all items for a specific sale
	GetSaleItems(saleID string) ([]models.SaleItem, error)

	// GetCustomerSales retrieves all sales for a specific customer
	GetCustomerSales(customerID string) ([]models.Sale, error)

	// Create creates a new sale record
	Create(sale *models.Sale) (string, error)

	// CreateSaleItem adds a new item to a sale
	CreateSaleItem(item *models.SaleItem) (string, error)

	// Update updates an existing sale
	Update(sale *models.Sale) error

	// Delete deletes a sale and its associated items
	Delete(id string) error
}

// SaleHandlers holds the repository dependency and JWT secret
type SaleHandlers struct {
	Repo      SaleRepository
	CabRepo   interface{} // Generic interface for cab repository
	AccRepo   interface{} // Generic interface for accessory repository
	CustRepo  interface{} // Generic interface for customer repository
	jwtSecret []byte
}

// NewSaleHandlers creates a new instance of SaleHandlers
func NewSaleHandlers(
	repo SaleRepository,
	cabRepo interface{},
	accRepo interface{},
	custRepo interface{},
	jwtSecret []byte,
) *SaleHandlers {
	return &SaleHandlers{
		Repo:      repo,
		CabRepo:   cabRepo,
		AccRepo:   accRepo,
		CustRepo:  custRepo,
		jwtSecret: jwtSecret,
	}
}

// RegisterSaleRoutes sets up the routes for sale operations within the provided Fiber router
func (h *SaleHandlers) RegisterSaleRoutes(r fiber.Router) {
	// Define middleware - Use the imported middleware package and the injected jwtSecret
	authRequired := middleware.JWTMiddleware(h.jwtSecret)

	// Group routes under '/sales'
	salesGroup := r.Group("/sales", authRequired)

	// Sales endpoints
	salesGroup.Get("/", h.GetSalesHandler)              // GET /api/sales
	salesGroup.Get("/:id", h.GetSaleByIDHandler)        // GET /api/sales/{id}
	salesGroup.Get("/:id/items", h.GetSaleItemsHandler) // GET /api/sales/{id}/items
	salesGroup.Post("/", h.CreateSaleHandler)           // POST /api/sales
	salesGroup.Put("/:id", h.UpdateSaleHandler)         // PUT /api/sales/{id}
	salesGroup.Delete("/:id", h.DeleteSaleHandler)      // DELETE /api/sales/{id}

	// Customer sales endpoints
	r.Get("/customers/:id/sales", authRequired, h.GetCustomerSalesHandler) // GET /api/customers/{id}/sales

	// Cab sales endpoint
	r.Post("/cabs/:id/sell", authRequired, h.SellCabHandler) // POST /api/cabs/{id}/sell
}

// GetSalesHandler handles requests to retrieve all sales with optional filtering
// @Summary Get all sales
// @Description Retrieves a list of sales, with optional filtering.
// @Tags Sales
// @Produce json
// @Security ApiKeyAuth
// @Param customer_id query string false "Filter by customer ID"
// @Param sold_by query string false "Filter by seller ID"
// @Param date_from query string false "Filter by sale date (from)"
// @Param date_to query string false "Filter by sale date (to)"
// @Success 200 {array} models.Sale "Successfully retrieved list of sales"
// @Failure 500 {object} ErrorResponse "Failed to retrieve sales"
// @Router /sales [get]
func (h *SaleHandlers) GetSalesHandler(c *fiber.Ctx) error {
	// Extract query parameters for filtering
	filters := make(map[string]interface{})

	if customerID := c.Query("customer_id"); customerID != "" {
		filters["customer_id"] = customerID
	}

	if soldBy := c.Query("sold_by"); soldBy != "" {
		filters["sold_by"] = soldBy
	}

	if dateFrom := c.Query("date_from"); dateFrom != "" {
		filters["date_from"] = dateFrom
	}

	if dateTo := c.Query("date_to"); dateTo != "" {
		filters["date_to"] = dateTo
	}

	sales, err := h.Repo.GetAll(filters)
	if err != nil {
		log.Printf("Error getting sales: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":       "Failed to retrieve sales",
			"status_code": fiber.StatusInternalServerError,
		})
	}

	return c.Status(fiber.StatusOK).JSON(sales)
}

// GetSaleByIDHandler handles requests to retrieve a single sale by ID
// @Summary Get sale by ID
// @Description Retrieves a single sale by its ID.
// @Tags Sales
// @Produce json
// @Security ApiKeyAuth
// @Param id path string true "Sale ID"
// @Success 200 {object} models.Sale "Successfully retrieved sale"
// @Failure 404 {object} ErrorResponse "Sale not found"
// @Failure 500 {object} ErrorResponse "Failed to retrieve sale"
// @Router /sales/{id} [get]
func (h *SaleHandlers) GetSaleByIDHandler(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":       "Sale ID is required",
			"status_code": fiber.StatusBadRequest,
		})
	}

	sale, err := h.Repo.GetByID(id)
	if err != nil {
		log.Printf("Error getting sale by ID %s: %v", id, err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":       "Failed to retrieve sale",
			"status_code": fiber.StatusInternalServerError,
		})
	}

	if sale == nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error":       "Sale not found",
			"status_code": fiber.StatusNotFound,
		})
	}

	return c.Status(fiber.StatusOK).JSON(sale)
}

// GetSaleItemsHandler handles requests to retrieve all items for a specific sale
// @Summary Get sale items
// @Description Retrieves all items for a specific sale.
// @Tags Sales
// @Produce json
// @Security ApiKeyAuth
// @Param id path string true "Sale ID"
// @Success 200 {array} models.SaleItem "Successfully retrieved sale items"
// @Failure 404 {object} ErrorResponse "Sale not found"
// @Failure 500 {object} ErrorResponse "Failed to retrieve sale items"
// @Router /sales/{id}/items [get]
func (h *SaleHandlers) GetSaleItemsHandler(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":       "Sale ID is required",
			"status_code": fiber.StatusBadRequest,
		})
	}

	// First check if the sale exists
	sale, err := h.Repo.GetByID(id)
	if err != nil {
		log.Printf("Error checking sale existence by ID %s: %v", id, err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":       "Failed to retrieve sale",
			"status_code": fiber.StatusInternalServerError,
		})
	}

	if sale == nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error":       "Sale not found",
			"status_code": fiber.StatusNotFound,
		})
	}

	// Get the sale items
	items, err := h.Repo.GetSaleItems(id)
	if err != nil {
		log.Printf("Error getting items for sale ID %s: %v", id, err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":       "Failed to retrieve sale items",
			"status_code": fiber.StatusInternalServerError,
		})
	}

	return c.Status(fiber.StatusOK).JSON(items)
}

// CreateSaleHandler handles requests to create a new sale
// @Summary Create a new sale
// @Description Creates a new sale record.
// @Tags Sales
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param sale body models.Sale true "Sale object to create"
// @Success 201 {object} models.Sale "Sale created successfully"
// @Failure 400 {object} ErrorResponse "Invalid request payload or missing required fields"
// @Failure 500 {object} ErrorResponse "Failed to create sale"
// @Router /sales [post]
func (h *SaleHandlers) CreateSaleHandler(c *fiber.Ctx) error {
	var newSale models.Sale
	if err := c.BodyParser(&newSale); err != nil {
		log.Printf("Error decoding create sale request: %v", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":       "Invalid request payload",
			"status_code": fiber.StatusBadRequest,
		})
	}

	// Basic validation
	if newSale.CustomerID == "" || newSale.SoldBy == "" || newSale.SaleDate == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":       "Missing required sale fields",
			"status_code": fiber.StatusBadRequest,
		})
	}

	// Set timestamps
	newSale.CreatedAt = time.Now()
	newSale.UpdatedAt = time.Now()

	// Create the sale
	saleID, err := h.Repo.Create(&newSale)
	if err != nil {
		log.Printf("Error creating sale: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":       "Failed to create sale",
			"status_code": fiber.StatusInternalServerError,
		})
	}

	// Set the ID in the response
	newSale.ID = saleID

	return c.Status(fiber.StatusCreated).JSON(newSale)
}

// UpdateSaleHandler handles requests to update an existing sale
// @Summary Update an existing sale
// @Description Updates an existing sale by its ID.
// @Tags Sales
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path string true "Sale ID"
// @Param sale body models.Sale true "Sale object with updated fields"
// @Success 200 {object} models.Sale "Sale updated successfully"
// @Failure 400 {object} ErrorResponse "Invalid request payload or missing required fields"
// @Failure 404 {object} ErrorResponse "Sale not found"
// @Failure 500 {object} ErrorResponse "Failed to update sale"
// @Router /sales/{id} [put]
func (h *SaleHandlers) UpdateSaleHandler(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":       "Sale ID is required",
			"status_code": fiber.StatusBadRequest,
		})
	}

	// Check if the sale exists
	existingSale, err := h.Repo.GetByID(id)
	if err != nil {
		log.Printf("Error checking sale existence by ID %s: %v", id, err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":       "Failed to retrieve sale",
			"status_code": fiber.StatusInternalServerError,
		})
	}

	if existingSale == nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error":       "Sale not found",
			"status_code": fiber.StatusNotFound,
		})
	}

	// Parse the updated sale
	var updatedSale models.Sale
	if err := c.BodyParser(&updatedSale); err != nil {
		log.Printf("Error decoding update sale request: %v", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":       "Invalid request payload",
			"status_code": fiber.StatusBadRequest,
		})
	}

	// Ensure the ID from the path is used
	updatedSale.ID = id

	// Basic validation
	if updatedSale.CustomerID == "" || updatedSale.SoldBy == "" || updatedSale.SaleDate == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":       "Missing required sale fields",
			"status_code": fiber.StatusBadRequest,
		})
	}

	// Update timestamp
	updatedSale.UpdatedAt = time.Now()

	// Preserve the creation timestamp
	updatedSale.CreatedAt = existingSale.CreatedAt

	// Update the sale
	err = h.Repo.Update(&updatedSale)
	if err != nil {
		log.Printf("Error updating sale: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":       "Failed to update sale",
			"status_code": fiber.StatusInternalServerError,
		})
	}

	return c.Status(fiber.StatusOK).JSON(updatedSale)
}

// DeleteSaleHandler handles requests to delete a sale
// @Summary Delete a sale
// @Description Deletes a sale by its ID.
// @Tags Sales
// @Produce json
// @Security ApiKeyAuth
// @Param id path string true "Sale ID"
// @Success 204 "Sale deleted successfully (No Content)"
// @Failure 404 {object} ErrorResponse "Sale not found"
// @Failure 500 {object} ErrorResponse "Failed to delete sale"
// @Router /sales/{id} [delete]
func (h *SaleHandlers) DeleteSaleHandler(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":       "Sale ID is required",
			"status_code": fiber.StatusBadRequest,
		})
	}

	// Check if the sale exists
	existingSale, err := h.Repo.GetByID(id)
	if err != nil {
		log.Printf("Error checking sale existence by ID %s: %v", id, err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":       "Failed to retrieve sale",
			"status_code": fiber.StatusInternalServerError,
		})
	}

	if existingSale == nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error":       "Sale not found",
			"status_code": fiber.StatusNotFound,
		})
	}

	// Delete the sale
	err = h.Repo.Delete(id)
	if err != nil {
		log.Printf("Error deleting sale: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":       "Failed to delete sale",
			"status_code": fiber.StatusInternalServerError,
		})
	}

	return c.Status(fiber.StatusNoContent).Send(nil)
}

// SellCabHandler handles requests to sell a cab with optional accessories
// @Summary Sell a cab
// @Description Sells a cab with optional accessories.
// @Tags Sales
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path int true "Cab ID"
// @Param sale body models.CabSalePayload true "Sale details"
// @Success 201 {object} models.CabSale "Cab sold successfully"
// @Failure 400 {object} ErrorResponse "Invalid request payload or missing required fields"
// @Failure 404 {object} ErrorResponse "Cab not found"
// @Failure 500 {object} ErrorResponse "Failed to process sale"
// @Router /cabs/{id}/sell [post]
func (h *SaleHandlers) SellCabHandler(c *fiber.Ctx) error {
	// Parse the cab ID from the URL
	cabIDStr := c.Params("id")
	cabID, err := strconv.Atoi(cabIDStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":       "Invalid cab ID format",
			"status_code": fiber.StatusBadRequest,
		})
	}

	// Parse the sale payload
	var salePayload models.CabSalePayload
	if err := c.BodyParser(&salePayload); err != nil {
		log.Printf("Error decoding cab sale request: %v", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":       "Invalid request payload",
			"status_code": fiber.StatusBadRequest,
		})
	}

	// Basic validation
	if salePayload.CustomerID == "" || salePayload.Quantity <= 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":       "Missing required sale fields or invalid quantity",
			"status_code": fiber.StatusBadRequest,
		})
	}

	// Get user ID from JWT token for the SoldBy field
	userID := c.Locals("user_id")
	if userID == nil {
		userID = "system" // Fallback if user ID is not available
	}

	// Create a new sale record
	newSale := models.Sale{
		CustomerID: salePayload.CustomerID,
		SoldBy:     fmt.Sprintf("%v", userID),
		SaleDate:   time.Now().Format("2006-01-02"),
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	// Calculate total price (cab price * quantity + accessories price)
	var totalPrice float64

	// Get cab price from repository and calculate total
	cabRepo, ok := h.CabRepo.(repositories.CabsRepository)
	if !ok {
		log.Println("Cab repository not initialized correctly")
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":       "Internal server error: Cab repository not available",
			"status_code": fiber.StatusInternalServerError,
		})
	}

	cab, err := cabRepo.GetCabByID(cabID)
	if err != nil {
		log.Printf("Error getting cab by ID %d: %v", cabID, err)
		// Check if the error is due to the cab not being found
		if strings.Contains(err.Error(), "not found") {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error":       fmt.Sprintf("Cab with ID %d not found", cabID),
				"status_code": fiber.StatusNotFound,
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":       "Failed to retrieve cab details",
			"status_code": fiber.StatusInternalServerError,
		})
	}

	// Calculate cab item price
	cabItemPrice := cab.Price * float64(salePayload.Quantity)
	totalPrice += cabItemPrice

	// Calculate accessories total price and create accessory sale items
	accRepo, ok := h.AccRepo.(repositories.AccessoryRepository)
	if !ok {
		log.Println("Accessory repository not initialized correctly")
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":       "Internal server error: Accessory repository not available",
			"status_code": fiber.StatusInternalServerError,
		})
	}

	var accessorySaleItems []models.SaleItem // To store successfully created accessory sale items for response

	for _, accessoryForSale := range salePayload.Accessories {
		accessory, err := accRepo.GetByID(c.Context(), accessoryForSale.ID)
		if err != nil {
			log.Printf("Error getting accessory by ID %d: %v. Skipping accessory.", accessoryForSale.ID, err)
			// Continue to the next accessory if not found or other error
			continue
		}

		// Calculate accessory item price
		accessoryItemPrice := accessory.Price * float64(accessoryForSale.Quantity)
		totalPrice += accessoryItemPrice

		// Create sale item for the accessory
		accessorySaleItem := models.SaleItem{
			SaleID:      "", // Will be set after the main sale is created
			ItemType:    "accessory",
			AccessoryID: strconv.Itoa(accessory.ID), // Convert int to string for SaleItem struct
			Quantity:    accessoryForSale.Quantity,
			UnitPrice:   accessory.Price,
			Subtotal:    accessoryItemPrice,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		}
		accessorySaleItems = append(accessorySaleItems, accessorySaleItem)
	}

	// Set the total price
	newSale.TotalPrice = totalPrice

	// Create the main sale record
	saleID, err := h.Repo.Create(&newSale)
	if err != nil {
		log.Printf("Error creating sale: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":       "Failed to create sale",
			"status_code": fiber.StatusInternalServerError,
		})
	}

	// Create sale item for the cab
	cabSaleItem := models.SaleItem{
		SaleID:     saleID,
		ItemType:   "cab",
		MultiCabID: strconv.Itoa(cab.ID), // Convert int to string for SaleItem struct
		Quantity:   salePayload.Quantity,
		UnitPrice:  cab.Price,
		Subtotal:   cabItemPrice,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	_, err = h.Repo.CreateSaleItem(&cabSaleItem)
	if err != nil {
		log.Printf("Error creating cab sale item for sale ID %s: %v", saleID, err)
		// Decide how to handle this error - potentially delete the main sale and other items?
		// For now, just log and continue, but this might leave inconsistent data.
		// A transaction would be better here.
	}

	// Create sale items for the accessories
	for _, accessorySaleItem := range accessorySaleItems {
		accessorySaleItem.SaleID = saleID // Set the sale ID
		_, err := h.Repo.CreateSaleItem(&accessorySaleItem)
		if err != nil {
			log.Printf("Error creating accessory sale item for sale ID %s and accessory ID %s: %v", saleID, accessorySaleItem.AccessoryID, err)
			// Log the error and continue. Again, a transaction would be better.
		}
	}

	// Prepare the accessories list for the response, including details from the fetched accessories
	responseAccessories := []map[string]interface{}{}
	for _, accessoryForSale := range salePayload.Accessories {
		accessory, err := accRepo.GetByID(c.Context(), accessoryForSale.ID)
		if err != nil {
			// If we couldn't fetch details for the sale item creation, we also can't for the response.
			// Log and skip this accessory in the response as well.
			log.Printf("Error getting accessory by ID %d for response: %v. Skipping.", accessoryForSale.ID, err)
			continue
		}
		responseAccessories = append(responseAccessories, map[string]interface{}{
			"id":        accessory.ID,
			"name":      accessory.Name,
			"price":     accessory.Price,
			"quantity":  accessoryForSale.Quantity,
			"unitPrice": accessory.Price, // Assuming unit price is the same as the accessory's price
		})
	}

	// Return the sale details
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"success":     true,
		"message":     "Cab sold successfully",
		"cabId":       cabID,
		"customerId":  salePayload.CustomerID,
		"quantity":    salePayload.Quantity,
		"accessories": responseAccessories, // Use the prepared accessories list
		"totalPrice":  totalPrice,
		"saleDate":    newSale.SaleDate,
		"saleId":      saleID,
	})
}

// GetCustomerSalesHandler handles requests to retrieve all sales for a specific customer
// @Summary Get customer sales
// @Description Retrieves all sales for a specific customer.
// @Tags Sales
// @Produce json
// @Security ApiKeyAuth
// @Param id path string true "Customer ID"
// @Success 200 {array} models.Sale "Successfully retrieved customer sales"
// @Failure 404 {object} ErrorResponse "Customer not found"
// @Failure 500 {object} ErrorResponse "Failed to retrieve customer sales"
// @Router /customers/{id}/sales [get]
func (h *SaleHandlers) GetCustomerSalesHandler(c *fiber.Ctx) error {
	customerID := c.Params("id")
	if customerID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":       "Customer ID is required",
			"status_code": fiber.StatusBadRequest,
		})
	}

	// Get the customer sales
	sales, err := h.Repo.GetCustomerSales(customerID)
	if err != nil {
		log.Printf("Error getting sales for customer ID %s: %v", customerID, err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":       "Failed to retrieve customer sales",
			"status_code": fiber.StatusInternalServerError,
		})
	}

	// If no sales found, return an empty array instead of 404
	if len(sales) == 0 {
		return c.Status(fiber.StatusOK).JSON([]models.Sale{})
	}

	return c.Status(fiber.StatusOK).JSON(sales)
}
