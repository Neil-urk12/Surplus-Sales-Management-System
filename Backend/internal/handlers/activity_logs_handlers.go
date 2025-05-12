package handlers

import (
	"fmt"
	"math"
	"net/http"
	"oop/internal/models"
	"oop/internal/repositories"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
)

// ActivityLogHandler handles HTTP requests for activity logs.
type ActivityLogHandler struct {
	repo repositories.LogsRepositoryInterface
}

// NewActivityLogHandler creates a new instance of ActivityLogHandler.
func NewActivityLogHandler(repo repositories.LogsRepositoryInterface) *ActivityLogHandler {
	return &ActivityLogHandler{repo: repo}
}

// GetActivityLogs godoc
// @Summary Get paginated activity logs
// @Description Retrieves a paginated list of activity logs, ordered by timestamp descending.
// @Tags ActivityLogs
// @Accept json
// @Produce json
// @Param page query int false "Page number for pagination." default(1)
// @Param limit query int false "Number of logs per page." default(10)
// @Success 200 {object} map[string]interface{} "{\"data\":[]models.ActivityLog, \"total\":int64, \"page\":int, \"last_page\":float64}"
// @Failure 400 {object} map[string]string "{\"error\": \"Invalid query parameter(s)\"}"
// @Failure 500 {object} map[string]string "{\"error\": \"Failed to retrieve activity logs\"}"
// @Router /api/v1/activity-logs [get]
func (h *ActivityLogHandler) GetActivityLogs(c *fiber.Ctx) error {
	page, err := strconv.Atoi(c.Query("page", "1"))
	if err != nil || page < 1 {
		page = 1
	}

	limit, err := strconv.Atoi(c.Query("limit", "10"))
	if err != nil || limit < 1 {
		limit = 10
	}

	logs, total, err := h.repo.GetLogs(page, limit)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to retrieve activity logs",
		})
	}

	return c.Status(http.StatusOK).JSON(fiber.Map{
		"data":      logs,
		"total":     total,
		"page":      page,
		"last_page": math.Ceil(float64(total) / float64(limit)),
	})
}

// GetFilteredActivityLogs godoc
// @Summary Get filtered and paginated activity logs
// @Description Retrieves a list of activity logs based on specified filters and pagination, ordered by timestamp descending.
// @Tags ActivityLogs
// @Accept json
// @Produce json
// @Param page query int false "Page number for pagination." default(1)
// @Param limit query int false "Number of logs per page." default(10)
// @Param user query string false "Filter logs by the user who performed the action (case-insensitive, partial match)."
// @Param action query string false "Filter logs by the action performed (case-insensitive, partial match)."
// @Param status query string false "Filter logs by the status of the action (case-insensitive, partial match)."
// @Param startDate query string false "Filter logs from this date (YYYY-MM-DD). Includes the entire day."
// @Param endDate query string false "Filter logs up to this date (YYYY-MM-DD). Includes the entire day."
// @Success 200 {object} map[string]interface{} "{\"data\":[]models.ActivityLog, \"total\":int64, \"page\":int, \"last_page\":float64}"
// @Failure 400 {object} map[string]string "{\"error\": \"Invalid query parameter(s), e.g., invalid date format\"}"
// @Failure 500 {object} map[string]string "{\"error\": \"Failed to retrieve filtered activity logs\"}"
// @Router /api/v1/activity-logs/filter [get]
func (h *ActivityLogHandler) GetFilteredActivityLogs(c *fiber.Ctx) error {
	page, err := strconv.Atoi(c.Query("page", "1"))
	if err != nil || page < 1 {
		page = 1
	}

	limit, err := strconv.Atoi(c.Query("limit", "10"))
	if err != nil || limit < 1 {
		limit = 10
	}

	user := c.Query("user")
	action := c.Query("action")
	status := c.Query("status")

	var startDate, endDate *time.Time
	startDateStr := c.Query("startDate")
	if startDateStr != "" {
		// Try parsing ISO format first (RFC3339)
		parsedDate, err := time.Parse(time.RFC3339, startDateStr)
		if err != nil {
			// Fall back to YYYY-MM-DD format if ISO parsing fails
			parsedDate, err = time.Parse("2006-01-02", startDateStr)
			if err != nil {
				return c.Status(http.StatusBadRequest).JSON(fiber.Map{
					"error": fmt.Sprintf("Invalid startDate format: %s. Use ISO format or YYYY-MM-DD.", startDateStr),
				})
			}
		}
		startDate = &parsedDate
	}

	endDateStr := c.Query("endDate")
	if endDateStr != "" {
		// Try parsing ISO format first (RFC3339)
		parsedDate, err := time.Parse(time.RFC3339, endDateStr)
		if err != nil {
			// Fall back to YYYY-MM-DD format if ISO parsing fails
			parsedDate, err = time.Parse("2006-01-02", endDateStr)
			if err != nil {
				return c.Status(http.StatusBadRequest).JSON(fiber.Map{
					"error": fmt.Sprintf("Invalid endDate format: %s. Use ISO format or YYYY-MM-DD.", endDateStr),
				})
			}
			// For YYYY-MM-DD format, adjust to end of day
			parsedDate = time.Date(parsedDate.Year(), parsedDate.Month(), parsedDate.Day(), 23, 59, 59, 999999999, parsedDate.Location())
		}
		endDate = &parsedDate
	}

	logs, total, err := h.repo.GetBasedOnFilter(page, limit, user, action, status, startDate, endDate)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to retrieve filtered activity logs",
		})
	}

	return c.Status(http.StatusOK).JSON(fiber.Map{
		"data":      logs,
		"total":     total,
		"page":      page,
		"last_page": math.Ceil(float64(total) / float64(limit)),
	})
}

// CreateActivityLog godoc
// @Summary Create a new activity log
// @Description Adds a new activity log entry to the system.
// @Tags ActivityLogs
// @Accept json
// @Produce json
// @Param log body models.ActivityLog true "Activity Log data to create. ID, Timestamp, CreatedAt, UpdatedAt are auto-generated."
// @Success 201 {object} models.ActivityLog "Successfully created activity log"
// @Failure 400 {object} map[string]string "{"error": "Invalid request payload or missing required fields"}"
// @Failure 500 {object} map[string]string "{"error": "Failed to create activity log"}"
// @Router /api/activity-logs [post]
func (h *ActivityLogHandler) CreateActivityLog(c *fiber.Ctx) error {
	logEntry := new(models.ActivityLog)

	if err := c.BodyParser(logEntry); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request payload: " + err.Error(),
		})
	}

	// Basic validation (can be expanded based on model requirements)
	if logEntry.User == "" || logEntry.Action == "" || logEntry.Status == "" {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Missing required fields: user, action, status",
		})
	}

	// ID, Timestamp, CreatedAt, UpdatedAt will be set by the repository Create method
	// or by database defaults if applicable.
	// We can ensure Timestamp is set if not provided, or let the repo handle it.
	if logEntry.Timestamp.IsZero() {
		logEntry.Timestamp = time.Now()
	}

	if err := h.repo.Create(logEntry); err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create activity log: " + err.Error(),
		})
	}

	return c.Status(http.StatusCreated).JSON(logEntry)
}
