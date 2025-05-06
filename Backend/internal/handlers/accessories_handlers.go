package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"oop/internal/models"
	"oop/internal/repositories"
	"strings"

	"github.com/gofiber/fiber/v2"
)

// ErrorResponse represents an API error response
type ErrorResponse struct {
	Error      string `json:"error"`
	Message    string `json:"message,omitempty"`
	StatusCode int    `json:"statusCode"`
	Timestamp  string `json:"timestamp"`
}

// SuccessResponse represents a successful API response
type SuccessResponse struct {
	Success   bool   `json:"success"`
	ID        int    `json:"id,omitempty"`
	Message   string `json:"message,omitempty"`
	Timestamp string `json:"timestamp"`
}

// AccessoriesListResponse represents a response with a list of accessories
type AccessoriesListResponse struct {
	Data       []models.Accessory `json:"data"`
	Count      int                `json:"count"`
	Page       int                `json:"page,omitempty"`
	PageSize   int                `json:"pageSize,omitempty"`
	TotalPages int                `json:"totalPages,omitempty"`
}

// AccessoriesHandler handles accessory-related requests
type AccessoriesHandler struct {
	Repo repositories.AccessoryRepository
}

// NewAccessoriesHandler creates a new accessories handler
func NewAccessoriesHandler(repo repositories.AccessoryRepository) *AccessoriesHandler {
	return &AccessoriesHandler{Repo: repo}
}

// GetAllAccessories returns all accessories
func (h *AccessoriesHandler) GetAllAccessories(c *fiber.Ctx) error {
	// Extract query parameters for filtering
	filters := make(map[string]interface{})
	if makeFilter := c.Query("make"); makeFilter != "" {
		filters["make"] = makeFilter
	}
	if statusFilter := c.Query("status"); statusFilter != "" {
		filters["status"] = statusFilter
	}
	if colorFilter := c.Query("unit_color"); colorFilter != "" {
		filters["unit_color"] = colorFilter
	}
	if searchFilter := c.Query("search"); searchFilter != "" {
		filters["search"] = searchFilter
	}

	// Call repository to get accessories
	accessories, err := h.Repo.GetAll(c.Context())
	if err != nil {
		// Log the error internally
		fmt.Printf("Error fetching accessories: %v\n", err) // Replace with proper logging
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to retrieve accessories"})
	}

	// Return the list of accessories as JSON
	return c.Status(http.StatusOK).JSON(fiber.Map{
		"data":  accessories,
		"count": len(accessories),
	})
}

// GetAccessoryByID returns a specific accessory by ID
func (h *AccessoriesHandler) GetAccessoryByID(c *fiber.Ctx) error {
	// Parse ID from URL path parameter
	id, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Invalid ID format. ID must be an integer."})
	}

	// Call repository to get accessory by ID
	accessory, err := h.Repo.GetByID(c.Context(), id)
	if err != nil {
		// Check if the error is 'not found'
		if strings.Contains(strings.ToLower(err.Error()), "not found") {
			return c.Status(http.StatusNotFound).JSON(fiber.Map{"error": fmt.Sprintf("Accessory with ID %d not found", id)})
		}
		// Handle other potential repository errors
		fmt.Printf("Error fetching accessory by ID %d: %v\n", id, err) // Replace with proper logging
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to retrieve accessory"})
	}

	// Return the accessory as JSON
	return c.Status(http.StatusOK).JSON(accessory)
}

// CreateAccessory creates a new accessory
func (h *AccessoriesHandler) CreateAccessory(c *fiber.Ctx) error {
	var input models.NewAccessoryInput

	// Parse request body into NewAccessoryInput struct
	if err := c.BodyParser(&input); err != nil {
		// Check for specific JSON parsing errors
		if _, ok := err.(*json.SyntaxError); ok || err == fiber.ErrUnprocessableEntity {
			return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Invalid JSON format in request body"})
		}
		// Handle other potential BodyParser errors
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": fmt.Sprintf("Failed to parse request body: %v", err)})
	}

	// Perform basic validation
	if input.Name == "" || string(input.Make) == "" || string(input.UnitColor) == "" {
		return c.Status(http.StatusUnprocessableEntity).JSON(fiber.Map{
			"error": "Missing required fields. Required: name, make, unit_color",
		})
	}

	// Create accessory
	accessory, err := h.Repo.Create(c.Context(), input)
	if err != nil {
		// Check for specific repository errors
		if strings.Contains(err.Error(), "cannot be empty") {
			return c.Status(http.StatusUnprocessableEntity).JSON(fiber.Map{"error": err.Error()})
		}
		// Handle other potential repository errors
		fmt.Printf("Error creating accessory: %v\n", err) // Replace with proper logging
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to create accessory"})
	}

	// Return success response
	return c.Status(http.StatusCreated).JSON(fiber.Map{
		"success": true,
		"id":      accessory.ID,
		"message": "Accessory created successfully",
	})
}

// UpdateAccessory updates an existing accessory
func (h *AccessoriesHandler) UpdateAccessory(c *fiber.Ctx) error {
	// Parse ID from URL path parameter
	id, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Invalid ID format. ID must be an integer."})
	}

	var input models.UpdateAccessoryInput

	// Parse request body
	if err := c.BodyParser(&input); err != nil {
		// Check for specific JSON parsing errors
		if _, ok := err.(*json.SyntaxError); ok || err == fiber.ErrUnprocessableEntity {
			return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Invalid JSON format in request body"})
		}
		// Handle other potential BodyParser errors
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": fmt.Sprintf("Failed to parse request body: %v", err)})
	}

	// Update accessory
	accessory, err := h.Repo.Update(c.Context(), id, input)
	if err != nil {
		// Check if the error is 'not found'
		if strings.Contains(strings.ToLower(err.Error()), "not found") {
			return c.Status(http.StatusNotFound).JSON(fiber.Map{"error": fmt.Sprintf("Accessory with ID %d not found for update", id)})
		}
		// Handle other potential repository errors
		fmt.Printf("Error updating accessory ID %d: %v\n", id, err) // Replace with proper logging
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to update accessory"})
	}

	// Return success response
	return c.Status(http.StatusOK).JSON(fiber.Map{
		"success": true,
		"id":      accessory.ID,
		"message": "Accessory updated successfully",
	})
}

// DeleteAccessory deletes an accessory
func (h *AccessoriesHandler) DeleteAccessory(c *fiber.Ctx) error {
	// Parse ID from URL path parameter
	id, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Invalid ID format. ID must be an integer."})
	}

	// Delete accessory
	err = h.Repo.Delete(c.Context(), id)
	if err != nil {
		// Check if the error is 'not found'
		if strings.Contains(strings.ToLower(err.Error()), "not found") {
			return c.Status(http.StatusNotFound).JSON(fiber.Map{"error": fmt.Sprintf("Accessory with ID %d not found for deletion", id)})
		}
		// Handle other potential repository errors
		fmt.Printf("Error deleting accessory ID %d: %v\n", id, err) // Replace with proper logging
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to delete accessory"})
	}

	// Return No Content status for successful deletion
	return c.SendStatus(http.StatusNoContent)
}
