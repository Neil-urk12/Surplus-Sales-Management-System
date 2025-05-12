package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"oop/internal/models"
	"oop/internal/repositories"
	"strings"
	"oop/internal/config"

	// Added for Timestamp in ErrorResponse/SuccessResponse
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
	Success   bool        `json:"success"`
	Data      interface{} `json:"data,omitempty"` // Changed to interface{} to be more generic
	Message   string      `json:"message,omitempty"`
	Timestamp string      `json:"timestamp"`
}

// AccessoriesListResponse represents a response with a list of accessories
// This can be used directly in Swagger docs if GetAllAccessories returns this structure.
// However, the current GetAllAccessories returns fiber.Map{"data": accessories, "count": len(accessories)}
// For simplicity with current handler structure, we might define success directly in annotations.
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
// @Summary Get all accessories
// @Description Get a list of all accessories, with optional filtering.
// @Tags Accessories
// @Accept json
// @Produce json
// @Param make query string false "Filter by make"
// @Param status query string false "Filter by status"
// @Param unit_color query string false "Filter by unit color"
// @Param search query string false "General search term"
// @Success 200 {object} AccessoriesListResponse "Successfully retrieved list of accessories"
// @Failure 500 {object} ErrorResponse "Failed to retrieve accessories"
// @Router /accessories [get]
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
// @Summary Get accessory by ID
// @Description Get a single accessory by its ID.
// @Tags Accessories
// @Accept json
// @Produce json
// @Param id path int true "Accessory ID"
// @Success 200 {object} models.Accessory "Successfully retrieved accessory"
// @Failure 400 {object} ErrorResponse "Invalid ID format. ID must be an integer."
// @Failure 404 {object} ErrorResponse "Accessory not found"
// @Failure 500 {object} ErrorResponse "Failed to retrieve accessory"
// @Router /accessories/{id} [get]
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
// @Summary Create a new accessory
// @Description Add a new accessory to the inventory.
// @Tags Accessories
// @Accept json
// @Produce json
// @Param accessory_input body models.NewAccessoryInput true "Accessory object to create"
// @Success 201 {object} SuccessResponse "Accessory created successfully"
// @Failure 400 {object} ErrorResponse "Invalid JSON format or failed to parse request body"
// @Failure 422 {object} ErrorResponse "Missing required fields or validation error"
// @Failure 500 {object} ErrorResponse "Failed to create accessory or failed to retrieve details after creation"
// @Router /accessories [post]
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

	// Handle null or empty image with default image URL
	if input.Image == "null" || input.Image == "" {
		input.Image = config.DefaultImageURL
	}

	// Create accessory
	createdAccessoryID, err := h.Repo.Create(c.Context(), input)
	if err != nil {
		// Check for specific repository errors
		if strings.Contains(err.Error(), "cannot be empty") {
			return c.Status(http.StatusUnprocessableEntity).JSON(fiber.Map{"error": err.Error()})
		}
		// Handle other potential repository errors
		fmt.Printf("Error creating accessory: %v\n", err) // Replace with proper logging
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to create accessory"})
	}

	// Fetch the newly created accessory to get all its details
	newlyCreatedAccessory, err := h.Repo.GetByID(c.Context(), createdAccessoryID)
	if err != nil {
		fmt.Printf("Error fetching newly created accessory %d: %v\n", createdAccessoryID, err) // Replace with proper logging
		// For now, returning an error if fetching fails, as the client expects the full object.
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Accessory created but failed to retrieve details"})
	}

	log.Printf("Returning newly created accessory: %+v\n", newlyCreatedAccessory)

	// Return success response
	return c.Status(http.StatusCreated).JSON(fiber.Map{
		"success": true,
		"data":    newlyCreatedAccessory,
		"message": "Accessory created successfully",
	})
}

// UpdateAccessory updates an existing accessory
// @Summary Update an existing accessory
// @Description Update an existing accessory by its ID.
// @Tags Accessories
// @Accept json
// @Produce json
// @Param id path int true "Accessory ID"
// @Param accessory_update body models.UpdateAccessoryInput true "Accessory object with updated fields"
// @Success 200 {object} SuccessResponse "Accessory updated successfully"
// @Failure 400 {object} ErrorResponse "Invalid ID format or invalid JSON format/parsing error"
// @Failure 404 {object} ErrorResponse "Accessory not found for update"
// @Failure 500 {object} ErrorResponse "Failed to update accessory"
// @Router /accessories/{id} [put]
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

	// Handle null image values
	if input.Image != nil {
		if *input.Image == "null" || *input.Image == "" {
			defaultImage := config.DefaultImageURL
			input.Image = &defaultImage
		}
	}

	// Update accessory
	updatedAccessory, err := h.Repo.Update(c.Context(), id, input)
	if err != nil {
		// Check if the error is 'not found'
		if strings.Contains(strings.ToLower(err.Error()), "not found") {
			return c.Status(http.StatusNotFound).JSON(fiber.Map{"error": fmt.Sprintf("Accessory with ID %d not found for update", id)})
		}
		// Handle other potential repository errors
		fmt.Printf("Error updating accessory ID %d: %v\n", id, err) // Replace with proper logging
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to update accessory"})
	}

	// Return success response with the updated accessory data
	return c.Status(http.StatusOK).JSON(fiber.Map{
		"success": true,
		"data":    updatedAccessory, // Return the full accessory object
		"message": "Accessory updated successfully",
	})
}

// DeleteAccessory deletes an accessory
// @Summary Delete an accessory
// @Description Delete an accessory by its ID.
// @Tags Accessories
// @Accept json
// @Produce json
// @Param id path int true "Accessory ID"
// @Success 204 "Accessory deleted successfully (No Content)"
// @Failure 400 {object} ErrorResponse "Invalid ID format. ID must be an integer."
// @Failure 404 {object} ErrorResponse "Accessory not found for deletion"
// @Failure 500 {object} ErrorResponse "Failed to delete accessory"
// @Router /accessories/{id} [delete]
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
