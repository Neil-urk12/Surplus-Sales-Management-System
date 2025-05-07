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

// CabsHandlers struct holds dependencies specifically for cab-related handlers.
type CabsHandlers struct {
	Repo repositories.CabsRepository
}

// NewCabsHandlers creates a new CabsHandlers struct.
func NewCabsHandlers(repo repositories.CabsRepository) *CabsHandlers {
	return &CabsHandlers{Repo: repo}
}

// GetCabs handles requests to retrieve a list of cabs, applying filters.
// @Summary Get all cabs
// @Description Get a list of all cabs, with optional filtering by make, status, unit color, or a general search term.
// @Tags Cabs
// @Accept json
// @Produce json
// @Param make query string false "Filter by make (e.g., Toyota)"
// @Param status query string false "Filter by status (e.g., Available, Maintenance)"
// @Param unit_color query string false "Filter by unit color (e.g., Red)"
// @Param search query string false "General search term for various fields"
// @Success 200 {array} models.MultiCab "Successfully retrieved list of cabs"
// @Failure 500 {object} ErrorResponse "Failed to retrieve cabs"
// @Router /cabs [get]
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
		return c.Status(http.StatusInternalServerError).JSON(ErrorResponse{
			Error:      "Failed to retrieve cabs",
			StatusCode: http.StatusInternalServerError,
		})
	}

	// Return the list of cabs as JSON
	return c.Status(http.StatusOK).JSON(cabs)
}

// GetCabByID handles requests to retrieve a single cab by its ID.
// @Summary Get cab by ID
// @Description Get a single cab by its ID.
// @Tags Cabs
// @Accept json
// @Produce json
// @Param id path int true "Cab ID"
// @Success 200 {object} models.MultiCab "Successfully retrieved cab"
// @Failure 400 {object} ErrorResponse "Invalid ID format. ID must be an integer."
// @Failure 404 {object} ErrorResponse "Cab not found"
// @Failure 500 {object} ErrorResponse "Failed to retrieve cab"
// @Router /cabs/{id} [get]
func (h *CabsHandlers) GetCabByID(c *fiber.Ctx) error {
	// Parse ID from URL path parameter
	id, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(ErrorResponse{
			Error:      "Invalid ID format. ID must be an integer.",
			StatusCode: http.StatusBadRequest,
		})
	}

	// Call repository to get cab by ID
	cab, err := h.Repo.GetCabByID(id)
	if err != nil {
		// Check if the error is 'not found'
		if strings.Contains(strings.ToLower(err.Error()), "not found") {
			return c.Status(http.StatusNotFound).JSON(ErrorResponse{
				Error:      fmt.Sprintf("Cab with ID %d not found", id),
				StatusCode: http.StatusNotFound,
			})
		}
		// Handle other potential repository errors
		fmt.Printf("Error fetching cab by ID %d: %v\n", id, err) // Replace with proper logging
		return c.Status(http.StatusInternalServerError).JSON(ErrorResponse{
			Error:      "Failed to retrieve cab",
			StatusCode: http.StatusInternalServerError,
		})
	}

	// Return the cab as JSON
	return c.Status(http.StatusOK).JSON(cab)
}

// AddCab handles requests to add a new cab to the inventory.
// @Summary Add a new cab
// @Description Add a new cab to the inventory.
// @Tags Cabs
// @Accept json
// @Produce json
// @Param cab body models.MultiCab true "Cab object to add. ID is auto-generated and should be omitted."
// @Success 201 {object} models.MultiCab "Cab added successfully"
// @Failure 400 {object} ErrorResponse "Invalid JSON format or failed to parse request body"
// @Failure 422 {object} ErrorResponse "Missing required fields or validation error"
// @Failure 500 {object} ErrorResponse "Failed to add new cab"
// @Router /cabs [post]
func (h *CabsHandlers) AddCab(c *fiber.Ctx) error {
	var cab models.MultiCab

	// Parse request body into MultiCab struct
	if err := c.BodyParser(&cab); err != nil {
		// Check for specific JSON parsing errors
		if _, ok := err.(*json.SyntaxError); ok || err == fiber.ErrUnprocessableEntity {
			return c.Status(http.StatusBadRequest).JSON(ErrorResponse{
				Error:      "Invalid JSON format in request body",
				StatusCode: http.StatusBadRequest,
			})
		}
		// Handle other potential BodyParser errors
		return c.Status(http.StatusBadRequest).JSON(ErrorResponse{
			Error:      fmt.Sprintf("Failed to parse request body: %v", err),
			StatusCode: http.StatusBadRequest,
		})
	}

	// Perform basic validation (consider using a validation library for complex cases)
	if cab.Name == "" || cab.Make == "" || cab.UnitColor == "" || cab.Status == "" {
		return c.Status(http.StatusUnprocessableEntity).JSON(ErrorResponse{
			Error:      "Missing required fields. Required: name, make, unit_color, status",
			StatusCode: http.StatusUnprocessableEntity,
		})
	}
	// ID should not be provided by the client for Add operations
	cab.ID = 0

	// Call repository to add the new cab
	addedCab, err := h.Repo.AddCab(cab)
	if err != nil {
		// Check for specific repository errors (e.g., validation error from repo)
		if strings.Contains(err.Error(), "cannot be empty") { // Example check
			return c.Status(http.StatusUnprocessableEntity).JSON(ErrorResponse{
				Error:      err.Error(),
				StatusCode: http.StatusUnprocessableEntity,
			})
		}
		// Handle other potential repository errors
		fmt.Printf("Error adding cab: %v\n", err) // Replace with proper logging
		return c.Status(http.StatusInternalServerError).JSON(ErrorResponse{
			Error:      "Failed to add new cab",
			StatusCode: http.StatusInternalServerError,
		})
	}

	// Return the newly added cab with generated ID and timestamps
	return c.Status(http.StatusCreated).JSON(addedCab)
}

// UpdateCab handles requests to update an existing cab.
// @Summary Update an existing cab
// @Description Update an existing cab by its ID. The ID in the path is authoritative.
// @Tags Cabs
// @Accept json
// @Produce json
// @Param id path int true "Cab ID"
// @Param cab_update body models.MultiCab true "Cab object with updated fields. ID in body is ignored."
// @Success 200 {object} models.MultiCab "Cab updated successfully"
// @Failure 400 {object} ErrorResponse "Invalid ID format or invalid JSON format/parsing error"
// @Failure 404 {object} ErrorResponse "Cab not found for update"
// @Failure 500 {object} ErrorResponse "Failed to update cab"
// @Router /cabs/{id} [put]
func (h *CabsHandlers) UpdateCab(c *fiber.Ctx) error {
	// Parse ID from URL path parameter
	id, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(ErrorResponse{
			Error:      "Invalid ID format. ID must be an integer.",
			StatusCode: http.StatusBadRequest,
		})
	}

	var updatedCabData models.MultiCab
	// Parse request body into MultiCab struct
	if err := c.BodyParser(&updatedCabData); err != nil {
		// Check for specific JSON parsing errors
		if _, ok := err.(*json.SyntaxError); ok || err == fiber.ErrUnprocessableEntity {
			return c.Status(http.StatusBadRequest).JSON(ErrorResponse{
				Error:      "Invalid JSON format in request body",
				StatusCode: http.StatusBadRequest,
			})
		}
		// Handle other potential BodyParser errors
		return c.Status(http.StatusBadRequest).JSON(ErrorResponse{
			Error:      fmt.Sprintf("Failed to parse request body: %v", err),
			StatusCode: http.StatusBadRequest,
		})
	}

	// ID from payload is usually ignored in PUT, path parameter 'id' is authoritative
	updatedCabData.ID = id

	// Call repository to update the cab
	resultCab, err := h.Repo.UpdateCab(id, updatedCabData)
	if err != nil {
		// Check if the error is 'not found'
		if strings.Contains(strings.ToLower(err.Error()), "not found") {
			return c.Status(http.StatusNotFound).JSON(ErrorResponse{
				Error:      fmt.Sprintf("Cab with ID %d not found for update", id),
				StatusCode: http.StatusNotFound,
			})
		}
		// Handle other potential repository errors
		fmt.Printf("Error updating cab ID %d: %v\n", id, err) // Replace with proper logging
		return c.Status(http.StatusInternalServerError).JSON(ErrorResponse{
			Error:      "Failed to update cab",
			StatusCode: http.StatusInternalServerError,
		})
	}

	// Return the updated cab data
	return c.Status(http.StatusOK).JSON(resultCab)
}

// DeleteCab handles requests to delete a cab by its ID.
// @Summary Delete a cab
// @Description Delete a cab by its ID.
// @Tags Cabs
// @Accept json
// @Produce json
// @Param id path int true "Cab ID"
// @Success 204 "Cab deleted successfully (No Content)"
// @Failure 400 {object} ErrorResponse "Invalid ID format. ID must be an integer."
// @Failure 404 {object} ErrorResponse "Cab not found for deletion"
// @Failure 500 {object} ErrorResponse "Failed to delete cab"
// @Router /cabs/{id} [delete]
func (h *CabsHandlers) DeleteCab(c *fiber.Ctx) error {
	// Parse ID from URL path parameter
	id, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(ErrorResponse{
			Error:      "Invalid ID format. ID must be an integer.",
			StatusCode: http.StatusBadRequest,
		})
	}

	// Call repository to delete the cab
	err = h.Repo.DeleteCab(id)
	if err != nil {
		// Check if the error is 'not found'
		if strings.Contains(strings.ToLower(err.Error()), "not found") {
			return c.Status(http.StatusNotFound).JSON(ErrorResponse{
				Error:      fmt.Sprintf("Cab with ID %d not found for deletion", id),
				StatusCode: http.StatusNotFound,
			})
		}
		// Handle other potential repository errors
		fmt.Printf("Error deleting cab ID %d: %v\n", id, err) // Replace with proper logging
		return c.Status(http.StatusInternalServerError).JSON(ErrorResponse{
			Error:      "Failed to delete cab",
			StatusCode: http.StatusInternalServerError,
		})
	}

	// Return No Content status for successful deletion
	return c.SendStatus(http.StatusNoContent)
}
