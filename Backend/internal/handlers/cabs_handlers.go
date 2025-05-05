package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"oop/internal/models"
	"oop/internal/repository"
	"strings"

	"github.com/gofiber/fiber/v2"
)

// CabsHandlers struct holds dependencies specifically for cab-related handlers.
type CabsHandlers struct {
	Repo repository.CabsRepository
}

// NewCabsHandlers creates a new CabsHandlers struct.
func NewCabsHandlers(repo repository.CabsRepository) *CabsHandlers {
	return &CabsHandlers{Repo: repo}
}

// GetCabs handles requests to retrieve a list of cabs, applying filters.
func (h *CabsHandlers) GetCabs(c *fiber.Ctx) error {
	// Extract query parameters for filtering
	filters := make(map[string]interface{})
	if makeFilter := c.Query("make"); makeFilter != "" {
		filters["make"] = makeFilter
	}
	if statusFilter := c.Query("status"); statusFilter != "" {
		filters["status"] = statusFilter
	}
	if unitColorFilter := c.Query("unit_color"); unitColorFilter != "" {
		filters["unit_color"] = unitColorFilter
	}
	if searchFilter := c.Query("search"); searchFilter != "" {
		filters["search"] = searchFilter
	}

	// Call repository to get cabs with filters
	cabs, err := h.Repo.GetCabs(filters)
	if err != nil {
		// Log the error internally
		fmt.Printf("Error fetching cabs: %v\n", err) // Replace with proper logging
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to retrieve cabs"})
	}

	// Return the list of cabs as JSON
	return c.Status(http.StatusOK).JSON(cabs)
}

// GetCabByID handles requests to retrieve a single cab by its ID.
func (h *CabsHandlers) GetCabByID(c *fiber.Ctx) error {
	// Parse ID from URL path parameter
	id, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Invalid ID format. ID must be an integer."})
	}

	// Call repository to get cab by ID
	cab, err := h.Repo.GetCabByID(id)
	if err != nil {
		// Check if the error is 'not found'
		if strings.Contains(strings.ToLower(err.Error()), "not found") {
			return c.Status(http.StatusNotFound).JSON(fiber.Map{"error": fmt.Sprintf("Cab with ID %d not found", id)})
		}
		// Handle other potential repository errors
		fmt.Printf("Error fetching cab by ID %d: %v\n", id, err) // Replace with proper logging
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to retrieve cab"})
	}

	// Return the cab as JSON
	return c.Status(http.StatusOK).JSON(cab)
}

// AddCab handles requests to add a new cab to the inventory.
func (h *CabsHandlers) AddCab(c *fiber.Ctx) error {
	var cab models.MultiCab

	// Parse request body into MultiCab struct
	if err := c.BodyParser(&cab); err != nil {
		// Check for specific JSON parsing errors
		if _, ok := err.(*json.SyntaxError); ok || err == fiber.ErrUnprocessableEntity {
			return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Invalid JSON format in request body"})
		}
		// Handle other potential BodyParser errors
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": fmt.Sprintf("Failed to parse request body: %v", err)})
	}

	// Perform basic validation (consider using a validation library for complex cases)
	if cab.Name == "" || cab.Make == "" || cab.UnitColor == "" || cab.Status == "" {
		return c.Status(http.StatusUnprocessableEntity).JSON(fiber.Map{
			"error": "Missing required fields. Required: name, make, unit_color, status",
		})
	}
	// ID should not be provided by the client for Add operations
	cab.ID = 0

	// Call repository to add the new cab
	addedCab, err := h.Repo.AddCab(cab)
	if err != nil {
		// Check for specific repository errors (e.g., validation error from repo)
		if strings.Contains(err.Error(), "cannot be empty") { // Example check
			return c.Status(http.StatusUnprocessableEntity).JSON(fiber.Map{"error": err.Error()})
		}
		// Handle other potential repository errors
		fmt.Printf("Error adding cab: %v\n", err) // Replace with proper logging
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to add new cab"})
	}

	// Return the newly added cab with generated ID and timestamps
	return c.Status(http.StatusCreated).JSON(addedCab)
}

// UpdateCab handles requests to update an existing cab.
func (h *CabsHandlers) UpdateCab(c *fiber.Ctx) error {
	// Parse ID from URL path parameter
	id, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Invalid ID format. ID must be an integer."})
	}

	var updatedCabData models.MultiCab
	// Parse request body into MultiCab struct
	if err := c.BodyParser(&updatedCabData); err != nil {
		// Check for specific JSON parsing errors
		if _, ok := err.(*json.SyntaxError); ok || err == fiber.ErrUnprocessableEntity {
			return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Invalid JSON format in request body"})
		}
		// Handle other potential BodyParser errors
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": fmt.Sprintf("Failed to parse request body: %v", err)})
	}

	// ID from payload is usually ignored in PUT, path parameter 'id' is authoritative
	updatedCabData.ID = id

	// Call repository to update the cab
	resultCab, err := h.Repo.UpdateCab(id, updatedCabData)
	if err != nil {
		// Check if the error is 'not found'
		if strings.Contains(strings.ToLower(err.Error()), "not found") {
			return c.Status(http.StatusNotFound).JSON(fiber.Map{"error": fmt.Sprintf("Cab with ID %d not found for update", id)})
		}
		// Handle other potential repository errors
		fmt.Printf("Error updating cab ID %d: %v\n", id, err) // Replace with proper logging
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to update cab"})
	}

	// Return the updated cab data
	return c.Status(http.StatusOK).JSON(resultCab)
}

// DeleteCab handles requests to delete a cab by its ID.
func (h *CabsHandlers) DeleteCab(c *fiber.Ctx) error {
	// Parse ID from URL path parameter
	id, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Invalid ID format. ID must be an integer."})
	}

	// Call repository to delete the cab
	err = h.Repo.DeleteCab(id)
	if err != nil {
		// Check if the error is 'not found'
		if strings.Contains(strings.ToLower(err.Error()), "not found") {
			return c.Status(http.StatusNotFound).JSON(fiber.Map{"error": fmt.Sprintf("Cab with ID %d not found for deletion", id)})
		}
		// Handle other potential repository errors
		fmt.Printf("Error deleting cab ID %d: %v\n", id, err) // Replace with proper logging
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to delete cab"})
	}

	// Return No Content status for successful deletion
	return c.SendStatus(http.StatusNoContent)
}
