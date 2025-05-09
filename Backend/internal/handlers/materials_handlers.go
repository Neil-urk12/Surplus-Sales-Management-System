package handlers

import (
	"log"
	"strconv"

	"oop/internal/middleware"
	"oop/internal/models"
	"oop/internal/repositories"

	"github.com/gofiber/fiber/v2"
)

// MaterialHandlers holds the repository dependency and JWT secret
type MaterialHandlers struct {
	Repo      repositories.MaterialRepository
	jwtSecret []byte
}

// NewMaterialHandlers creates a new instance of MaterialHandlers
func NewMaterialHandlers(repo repositories.MaterialRepository, jwtSecret []byte) *MaterialHandlers {
	return &MaterialHandlers{
		Repo:      repo,
		jwtSecret: jwtSecret,
	}
}

// RegisterMaterialRoutes sets up the routes for material operations within the provided Fiber router
func (h *MaterialHandlers) RegisterMaterialRoutes(r fiber.Router) {
	// Define middleware - Use the imported middleware package and the injected jwtSecret
	authRequired := middleware.JWTMiddleware(h.jwtSecret)

	// Group routes under '/materials'
	materialsGroup := r.Group("/materials", authRequired)

	materialsGroup.Get("/", h.GetMaterialsHandler)                   // GET /api/materials?params...
	materialsGroup.Get("/paginated", h.GetPaginatedMaterialsHandler) // GET /api/materials/paginated?page=1&limit=10
	materialsGroup.Get("/:id", h.GetMaterialHandler)                 // GET /api/materials/{id}
	materialsGroup.Post("/", h.CreateMaterialHandler)                // POST /api/materials
	materialsGroup.Put("/:id", h.UpdateMaterialHandler)              // PUT /api/materials/{id}
	materialsGroup.Delete("/:id", h.DeleteMaterialHandler)           // DELETE /api/materials/{id}
}

// GetMaterialsHandler handles requests to retrieve multiple materials with filtering
// @Summary Get all materials
// @Description Retrieves a list of materials, with optional filtering.
// @Tags Materials
// @Produce json
// @Security ApiKeyAuth
// @Param search query string false "Search term for material name or description"
// @Param category query string false "Filter by category"
// @Param supplier query string false "Filter by supplier"
// @Param status query string false "Filter by status (e.g., In Stock, Low Stock)"
// @Success 200 {array} models.Material "Successfully retrieved list of materials"
// @Failure 500 {object} ErrorResponse "Failed to retrieve materials"
// @Router /materials [get]
func (h *MaterialHandlers) GetMaterialsHandler(c *fiber.Ctx) error {
	// Extract query parameters for filtering
	searchTerm := c.Query("search")
	category := c.Query("category")
	supplier := c.Query("supplier")
	status := c.Query("status")

	materials, err := h.Repo.GetAll(searchTerm, category, supplier, status)
	if err != nil {
		log.Printf("Error getting materials: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{
			Error:      "Failed to retrieve materials",
			StatusCode: fiber.StatusInternalServerError,
		})
	}

	return c.Status(fiber.StatusOK).JSON(materials)
}

// GetMaterialHandler handles requests to retrieve a single material by ID
// @Summary Get material by ID
// @Description Retrieves a single material by its ID.
// @Tags Materials
// @Produce json
// @Security ApiKeyAuth
// @Param id path int true "Material ID"
// @Success 200 {object} models.Material "Successfully retrieved material"
// @Failure 400 {object} ErrorResponse "Invalid Material ID format"
// @Failure 404 {object} ErrorResponse "Material not found"
// @Failure 500 {object} ErrorResponse "Failed to retrieve material"
// @Router /materials/{id} [get]
func (h *MaterialHandlers) GetMaterialHandler(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Error:      "Invalid Material ID format",
			StatusCode: fiber.StatusBadRequest,
		})
	}

	material, err := h.Repo.GetByID(id)
	if err != nil {
		log.Printf("Error getting material by ID %d: %v", id, err)
		return c.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{
			Error:      "Failed to retrieve material",
			StatusCode: fiber.StatusInternalServerError,
		})
	}
	if material == nil {
		return c.Status(fiber.StatusNotFound).JSON(ErrorResponse{
			Error:      "Material not found",
			StatusCode: fiber.StatusNotFound,
		})
	}

	return c.Status(fiber.StatusOK).JSON(material)
}

// CreateMaterialHandler handles requests to create a new material
// @Summary Create a new material
// @Description Adds a new material to the inventory.
// @Tags Materials
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param material body models.Material true "Material object to create"
// @Success 201 {object} models.Material "Material created successfully"
// @Failure 400 {object} ErrorResponse "Invalid request payload or missing required fields"
// @Failure 500 {object} ErrorResponse "Failed to create material"
// @Router /materials [post]
func (h *MaterialHandlers) CreateMaterialHandler(c *fiber.Ctx) error {
	var newMaterial models.Material
	if err := c.BodyParser(&newMaterial); err != nil {
		log.Printf("Error decoding create material request: %v", err)
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Error:      "Invalid request payload",
			StatusCode: fiber.StatusBadRequest,
		})
	}

	// Basic validation (could be expanded)
	if newMaterial.Name == "" || newMaterial.Category == "" || newMaterial.Supplier == "" || newMaterial.Status == "" {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Error:      "Missing required material fields",
			StatusCode: fiber.StatusBadRequest,
		})
	}

	id, err := h.Repo.Create(&newMaterial)
	if err != nil {
		log.Printf("Error creating material: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{
			Error:      "Failed to create material",
			StatusCode: fiber.StatusInternalServerError,
		})
	}

	// Optionally fetch the full created object to get timestamps
	createdMaterial, err := h.Repo.GetByID(id)
	if err != nil || createdMaterial == nil {
		log.Printf("Error fetching created material ID %d after creation: %v", id, err)
		// Respond with the ID even if fetch fails, as creation succeeded
		// You might prefer to return the original input + ID here
		newMaterial.ID = id
		return c.Status(fiber.StatusCreated).JSON(newMaterial) // Return input + ID
	}

	return c.Status(fiber.StatusCreated).JSON(createdMaterial)
}

// UpdateMaterialHandler handles requests to update an existing material
// @Summary Update an existing material
// @Description Updates an existing material by its ID.
// @Tags Materials
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path int true "Material ID"
// @Param material body models.Material true "Material object with updated fields"
// @Success 200 {object} models.Material "Material updated successfully"
// @Success 204 "Material updated, but fetch failed (No Content)"
// @Failure 400 {object} ErrorResponse "Invalid Material ID format or invalid request payload or missing required fields"
// @Failure 500 {object} ErrorResponse "Failed to update material"
// @Router /materials/{id} [put]
func (h *MaterialHandlers) UpdateMaterialHandler(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Error:      "Invalid Material ID format",
			StatusCode: fiber.StatusBadRequest,
		})
	}

	var updatedMaterial models.Material
	if err := c.BodyParser(&updatedMaterial); err != nil {
		log.Printf("Error decoding update material request for ID %d: %v", id, err)
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Error:      "Invalid request payload",
			StatusCode: fiber.StatusBadRequest,
		})
	}

	// Ensure the ID from the path is used
	updatedMaterial.ID = id

	// Basic validation
	if updatedMaterial.Name == "" || updatedMaterial.Category == "" || updatedMaterial.Supplier == "" || updatedMaterial.Status == "" {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Error:      "Missing required material fields",
			StatusCode: fiber.StatusBadRequest,
		})
	}

	err = h.Repo.Update(&updatedMaterial)
	if err != nil {
		log.Printf("Error updating material ID %d: %v", id, err)
		// Could check for specific errors like 'not found' if the repo layer provides them
		return c.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{
			Error:      "Failed to update material",
			StatusCode: fiber.StatusInternalServerError,
		})
	}

	// Fetch the updated object to return the latest state including UpdatedAt
	finalMaterial, err := h.Repo.GetByID(id)
	if err != nil || finalMaterial == nil {
		log.Printf("Error fetching updated material ID %d after update: %v", id, err)
		// Update succeeded, but fetch failed. Return No Content.
		return c.SendStatus(fiber.StatusNoContent)
	}

	return c.Status(fiber.StatusOK).JSON(finalMaterial)
}

// DeleteMaterialHandler handles requests to delete a material
// @Summary Delete a material
// @Description Deletes a material by its ID.
// @Tags Materials
// @Produce json
// @Security ApiKeyAuth
// @Param id path int true "Material ID"
// @Success 204 "Material deleted successfully (No Content)"
// @Failure 400 {object} ErrorResponse "Invalid Material ID format"
// @Failure 500 {object} ErrorResponse "Failed to delete material"
// @Router /materials/{id} [delete]
func (h *MaterialHandlers) DeleteMaterialHandler(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Error:      "Invalid Material ID format",
			StatusCode: fiber.StatusBadRequest,
		})
	}

	err = h.Repo.Delete(id)
	if err != nil {
		log.Printf("Error deleting material ID %d: %v", id, err)
		// Could check for specific errors like 'not found'
		return c.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{
			Error:      "Failed to delete material",
			StatusCode: fiber.StatusInternalServerError,
		})
	}

	return c.SendStatus(fiber.StatusNoContent) // Standard response for successful deletion
}

// Maximum number of items that can be requested per page
const maxPageLimit = 100

// GetPaginatedMaterialsHandler handles requests to retrieve paginated materials
func (h *MaterialHandlers) GetPaginatedMaterialsHandler(c *fiber.Ctx) error {
	pageStr := c.Query("page", "1")
	limitStr := c.Query("limit", "10")

	page, err := strconv.Atoi(pageStr)
	if err != nil {
		log.Printf("Invalid page parameter: %s, using default value 1", pageStr)
		page = 1
	}

	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		log.Printf("Invalid limit parameter: %s, using default value 10", limitStr)
		limit = 10
	}

	// Ensure page is at least 1
	if page < 1 {
		log.Printf("Page parameter less than 1: %d, using default value 1", page)
		page = 1
	}

	// Ensure limit is between 1 and maxPageLimit
	if limit < 1 {
		log.Printf("Limit parameter less than 1: %d, using default value 10", limit)
		limit = 10 // Default to 10 if invalid
	} else if limit > maxPageLimit {
		log.Printf("Limit parameter exceeds maximum allowed value: %d, using maximum value %d", limit, maxPageLimit)
		limit = maxPageLimit
	}

	searchTerm := c.Query("search")
	category := c.Query("category")
	supplier := c.Query("supplier")
	status := c.Query("status")

	materials, total, err := h.Repo.GetPaginated(page, limit, searchTerm, category, supplier, status)
	if err != nil {
		log.Printf("Error getting paginated materials: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to retrieve materials",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"materials":  materials,
		"total":      total,
		"page":       page,
		"limit":      limit,
		"totalPages": (total + int64(limit) - 1) / int64(limit),
	})
}
