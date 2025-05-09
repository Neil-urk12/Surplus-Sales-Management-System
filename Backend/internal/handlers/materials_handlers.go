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
func (h *MaterialHandlers) GetMaterialsHandler(c *fiber.Ctx) error {
	// Extract query parameters for filtering
	searchTerm := c.Query("search")
	category := c.Query("category")
	supplier := c.Query("supplier")
	status := c.Query("status")

	materials, err := h.Repo.GetAll(searchTerm, category, supplier, status)
	if err != nil {
		log.Printf("Error getting materials: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to retrieve materials",
		})
	}

	return c.Status(fiber.StatusOK).JSON(materials)
}

// GetMaterialHandler handles requests to retrieve a single material by ID
func (h *MaterialHandlers) GetMaterialHandler(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid Material ID format",
		})
	}

	material, err := h.Repo.GetByID(id)
	if err != nil {
		log.Printf("Error getting material by ID %d: %v", id, err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to retrieve material",
		})
	}
	if material == nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Material not found",
		})
	}

	return c.Status(fiber.StatusOK).JSON(material)
}

// CreateMaterialHandler handles requests to create a new material
func (h *MaterialHandlers) CreateMaterialHandler(c *fiber.Ctx) error {
	var newMaterial models.Material
	if err := c.BodyParser(&newMaterial); err != nil {
		log.Printf("Error decoding create material request: %v", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request payload",
		})
	}

	// Basic validation (could be expanded)
	if newMaterial.Name == "" || newMaterial.Category == "" || newMaterial.Supplier == "" || newMaterial.Status == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Missing required material fields",
		})
	}

	id, err := h.Repo.Create(&newMaterial)
	if err != nil {
		log.Printf("Error creating material: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create material",
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
func (h *MaterialHandlers) UpdateMaterialHandler(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid Material ID format",
		})
	}

	var updatedMaterial models.Material
	if err := c.BodyParser(&updatedMaterial); err != nil {
		log.Printf("Error decoding update material request for ID %d: %v", id, err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request payload",
		})
	}

	// Ensure the ID from the path is used
	updatedMaterial.ID = id

	// Basic validation
	if updatedMaterial.Name == "" || updatedMaterial.Category == "" || updatedMaterial.Supplier == "" || updatedMaterial.Status == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Missing required material fields",
		})
	}

	err = h.Repo.Update(&updatedMaterial)
	if err != nil {
		log.Printf("Error updating material ID %d: %v", id, err)
		// Could check for specific errors like 'not found' if the repo layer provides them
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update material",
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
func (h *MaterialHandlers) DeleteMaterialHandler(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid Material ID format",
		})
	}

	err = h.Repo.Delete(id)
	if err != nil {
		log.Printf("Error deleting material ID %d: %v", id, err)
		// Could check for specific errors like 'not found'
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to delete material",
		})
	}

	return c.SendStatus(fiber.StatusNoContent) // Standard response for successful deletion
}

// GetPaginatedMaterialsHandler handles requests to retrieve paginated materials
func (h *MaterialHandlers) GetPaginatedMaterialsHandler(c *fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "10"))
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
