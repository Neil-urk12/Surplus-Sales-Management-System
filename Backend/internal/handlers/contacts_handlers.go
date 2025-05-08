package handlers

import (
	"log"
	"oop/internal/middleware"
	"oop/internal/models"
	"oop/internal/repositories"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// CustomerResponse defines the structure for a single customer response.
// It omits sensitive or unnecessary fields for client-side display.
type CustomerResponse struct {
	ID             string `json:"id"`
	Name           string `json:"name"`
	Email          string `json:"email"`
	Phone          string `json:"phone"`
	Address        string `json:"address,omitempty"`
	DateRegistered string `json:"dateRegistered"`
	CreatedAt      string `json:"createdAt"`
	UpdatedAt      string `json:"updatedAt"`
}

// CustomerListResponse defines the structure for a list of customers.
type CustomerListResponse struct {
	Customers []*CustomerResponse `json:"customers"`
}

// CreateCustomerRequest defines the expected payload for creating a new customer.
type CreateCustomerRequest struct {
	Name    string `json:"name" validate:"required,min=2,max=100"`
	Email   string `json:"email" validate:"required,email"`
	Phone   string `json:"phone" validate:"required,e164"` // e164 format for phone numbers
	Address string `json:"address,omitempty" validate:"max=255"`
}

// UpdateCustomerRequest defines the expected payload for updating an existing customer.
// All fields are optional for partial updates.
// Similar to CreateCustomerRequest but fields are pointers or have omitempty if not pointers and not always required.
// For simplicity, using direct fields and relying on handler logic to apply non-empty values.
type UpdateCustomerRequest struct {
	Name    string `json:"name,omitempty" validate:"omitempty,min=2,max=100"`
	Email   string `json:"email,omitempty" validate:"omitempty,email"`
	Phone   string `json:"phone,omitempty" validate:"omitempty,e164"`
	Address string `json:"address,omitempty" validate:"omitempty,max=255"`
}

// CustomerHandler holds the repository and JWT secret.
type CustomerHandler struct {
	Repo      repositories.CustomerRepository
	jwtSecret []byte
}

// NewCustomerHandler creates a new CustomerHandler instance.
func NewCustomerHandler(repo repositories.CustomerRepository, jwtSecret []byte) *CustomerHandler {
	return &CustomerHandler{
		Repo:      repo,
		jwtSecret: jwtSecret,
	}
}

// toCustomerResponse converts a models.Customer to a CustomerResponse.
func toCustomerResponse(customer *models.Customer) *CustomerResponse {
	if customer == nil {
		return nil
	}
	return &CustomerResponse{
		ID:             customer.ID,
		Name:           customer.Name,
		Email:          customer.Email,
		Phone:          customer.Phone,
		Address:        customer.Address,
		DateRegistered: customer.DateRegistered,
		CreatedAt:      customer.CreatedAt,
		UpdatedAt:      customer.UpdatedAt,
	}
}

// RegisterCustomerRoutes sets up the routes for customer operations.
func (h *CustomerHandler) RegisterCustomerRoutes(r fiber.Router) {
	authRequired := middleware.JWTMiddleware(h.jwtSecret)
	customerGroup := r.Group("/customers", authRequired)

	customerGroup.Post("/", h.CreateCustomer)
	customerGroup.Get("/", h.GetAllCustomers)
	customerGroup.Get("/:id", h.GetCustomer)
	customerGroup.Put("/:id", h.UpdateCustomer)
	customerGroup.Delete("/:id", h.DeleteCustomer)
}

// CreateCustomer handles the creation of a new customer.
// @Summary Create a new customer
// @Description Adds a new customer to the system.
// @Tags Customers
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param customer body CreateCustomerRequest true "Customer information"
// @Success 201 {object} CustomerResponse "Customer created successfully"
// @Failure 400 {object} ErrorResponse "Invalid request payload or validation error"
// @Failure 500 {object} ErrorResponse "Failed to create customer"
// @Router /customers [post]
func (h *CustomerHandler) CreateCustomer(c *fiber.Ctx) error {
	var req CreateCustomerRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{Error: "Invalid request payload", StatusCode: fiber.StatusBadRequest})
	}

	// TODO: Add validation for req struct using a library like go-playground/validator
	if req.Name == "" || req.Email == "" || req.Phone == "" {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{Error: "Name, email, and phone are required", StatusCode: fiber.StatusBadRequest})
	}

	customer := &models.Customer{
		ID:      uuid.New().String(), // Repository will also generate if empty, but good to have it here too.
		Name:    req.Name,
		Email:   req.Email,
		Phone:   req.Phone,
		Address: req.Address,
	}

	createdCustomer, err := h.Repo.CreateCustomer(customer)
	if err != nil {
		log.Printf("Error creating customer: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{Error: "Failed to create customer", StatusCode: fiber.StatusInternalServerError})
	}

	return c.Status(fiber.StatusCreated).JSON(toCustomerResponse(createdCustomer))
}

// GetAllCustomers handles retrieving all customers.
// @Summary Get all customers
// @Description Retrieves a list of all customers.
// @Tags Customers
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} CustomerListResponse "Successfully retrieved list of customers"
// @Failure 500 {object} ErrorResponse "Failed to retrieve customers"
// @Router /customers [get]
func (h *CustomerHandler) GetAllCustomers(c *fiber.Ctx) error {
	customers, err := h.Repo.GetAllCustomers()
	if err != nil {
		log.Printf("Error getting all customers: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{Error: "Failed to retrieve customers", StatusCode: fiber.StatusInternalServerError})
	}

	customerResponses := make([]*CustomerResponse, len(customers))
	for i, cust := range customers {
		customerResponses[i] = toCustomerResponse(cust)
	}

	return c.Status(fiber.StatusOK).JSON(CustomerListResponse{Customers: customerResponses})
}

// GetCustomer handles retrieving a single customer by ID.
// @Summary Get customer by ID
// @Description Retrieves a single customer by their ID.
// @Tags Customers
// @Produce json
// @Security ApiKeyAuth
// @Param id path string true "Customer ID (UUID format)"
// @Success 200 {object} CustomerResponse "Successfully retrieved customer"
// @Failure 400 {object} ErrorResponse "Invalid Customer ID format"
// @Failure 404 {object} ErrorResponse "Customer not found"
// @Failure 500 {object} ErrorResponse "Failed to retrieve customer"
// @Router /customers/{id} [get]
func (h *CustomerHandler) GetCustomer(c *fiber.Ctx) error {
	id := c.Params("id")
	if _, err := uuid.Parse(id); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{Error: "Invalid Customer ID format", StatusCode: fiber.StatusBadRequest})
	}

	customer, err := h.Repo.GetCustomerByID(id)
	if err != nil {
		// Differentiate between not found and other errors if repo returns specific errors
		log.Printf("Error getting customer by ID %s: %v", id, err)
		// For now, assume any error from repo.GetCustomerByID for a non-existent ID might be caught by specific error string check
		if err.Error() == "customer with ID "+id+" not found" { // This is fragile; better to use custom error types or errors.Is
			return c.Status(fiber.StatusNotFound).JSON(ErrorResponse{Error: "Customer not found", StatusCode: fiber.StatusNotFound})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{Error: "Failed to retrieve customer", StatusCode: fiber.StatusInternalServerError})
	}

	return c.Status(fiber.StatusOK).JSON(toCustomerResponse(customer))
}

// UpdateCustomer handles updating an existing customer.
// @Summary Update an existing customer
// @Description Updates an existing customer by their ID.
// @Tags Customers
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path string true "Customer ID (UUID format)"
// @Param customer body UpdateCustomerRequest true "Customer information to update"
// @Success 200 {object} CustomerResponse "Customer updated successfully"
// @Failure 400 {object} ErrorResponse "Invalid Customer ID format or invalid request payload"
// @Failure 404 {object} ErrorResponse "Customer not found"
// @Failure 500 {object} ErrorResponse "Failed to update customer"
// @Router /customers/{id} [put]
func (h *CustomerHandler) UpdateCustomer(c *fiber.Ctx) error {
	id := c.Params("id")
	if _, err := uuid.Parse(id); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{Error: "Invalid Customer ID format", StatusCode: fiber.StatusBadRequest})
	}

	var req UpdateCustomerRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{Error: "Invalid request payload", StatusCode: fiber.StatusBadRequest})
	}

	// TODO: Add validation for req struct

	existingCustomer, err := h.Repo.GetCustomerByID(id)
	if err != nil {
		log.Printf("Error finding customer %s for update: %v", id, err)
		// Differentiate error types if possible
		return c.Status(fiber.StatusNotFound).JSON(ErrorResponse{Error: "Customer not found", StatusCode: fiber.StatusNotFound})
	}

	// Apply updates from request if fields are provided
	if req.Name != "" {
		existingCustomer.Name = req.Name
	}
	if req.Email != "" {
		existingCustomer.Email = req.Email
	}
	if req.Phone != "" {
		existingCustomer.Phone = req.Phone
	}
	if req.Address != "" {
		existingCustomer.Address = req.Address
	}
	// Note: ID, DateRegistered, CreatedAt should not be changed here. UpdatedAt is handled by the repo.

	updatedCustomer, err := h.Repo.UpdateCustomer(existingCustomer)
	if err != nil {
		log.Printf("Error updating customer %s: %v", id, err)
		return c.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{Error: "Failed to update customer", StatusCode: fiber.StatusInternalServerError})
	}

	return c.Status(fiber.StatusOK).JSON(toCustomerResponse(updatedCustomer))
}

// DeleteCustomer handles deleting a customer by ID.
// @Summary Delete a customer
// @Description Deletes a customer by their ID.
// @Tags Customers
// @Produce json
// @Security ApiKeyAuth
// @Param id path string true "Customer ID (UUID format)"
// @Success 204 "Customer deleted successfully (No Content)"
// @Failure 400 {object} ErrorResponse "Invalid Customer ID format"
// @Failure 404 {object} ErrorResponse "Customer not found"
// @Failure 500 {object} ErrorResponse "Failed to delete customer"
// @Router /customers/{id} [delete]
func (h *CustomerHandler) DeleteCustomer(c *fiber.Ctx) error {
	id := c.Params("id")
	if _, err := uuid.Parse(id); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{Error: "Invalid Customer ID format", StatusCode: fiber.StatusBadRequest})
	}

	err := h.Repo.DeleteCustomer(id)
	if err != nil {
		log.Printf("Error deleting customer %s: %v", id, err)
		// Check if the error indicates "not found"
		if err.Error() == "customer with ID "+id+" not found for deletion" { // Fragile check
			return c.Status(fiber.StatusNotFound).JSON(ErrorResponse{Error: "Customer not found", StatusCode: fiber.StatusNotFound})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{Error: "Failed to delete customer", StatusCode: fiber.StatusInternalServerError})
	}

	return c.SendStatus(fiber.StatusNoContent)
}
