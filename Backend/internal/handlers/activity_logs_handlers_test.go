package handlers

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"oop/internal/models"
	"strings"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockLogsRepository is a mock type for the LogsRepositoryInterface
type MockLogsRepository struct {
	mock.Mock
}

func (m *MockLogsRepository) Create(log *models.ActivityLog) error {
	args := m.Called(log)
	return args.Error(0)
}

func (m *MockLogsRepository) GetLogs(page, limit int) ([]models.ActivityLog, int64, error) {
	args := m.Called(page, limit)
	return args.Get(0).([]models.ActivityLog), args.Get(1).(int64), args.Error(2)
}

func (m *MockLogsRepository) GetBasedOnFilter(page, limit int, user, action, status string, startDate, endDate *time.Time) ([]models.ActivityLog, int64, error) {
	args := m.Called(page, limit, user, action, status, startDate, endDate)
	return args.Get(0).([]models.ActivityLog), args.Get(1).(int64), args.Error(2)
}

func setupAppAndHandler(mockRepo *MockLogsRepository) *fiber.App {
	app := fiber.New()
	h := NewActivityLogHandler(mockRepo)
	api := app.Group("/api")
	activityLogRoutes := api.Group("/activity-logs")
	activityLogRoutes.Get("/", h.GetActivityLogs)
	activityLogRoutes.Get("/filter", h.GetFilteredActivityLogs)
	activityLogRoutes.Post("/", h.CreateActivityLog)
	return app
}

func TestGetActivityLogsHandler(t *testing.T) {
	t.Run("SuccessfulRetrieval", func(t *testing.T) {
		mockRepo := new(MockLogsRepository)
		app := setupAppAndHandler(mockRepo)

		expectedLogs := []models.ActivityLog{{ID: "1", User: "test"}}
		mockRepo.On("GetLogs", 1, 10).Return(expectedLogs, int64(1), nil)

		req := httptest.NewRequest(http.MethodGet, "/api/activity-logs?page=1&limit=10", nil)
		resp, _ := app.Test(req)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		bodyBytes, _ := io.ReadAll(resp.Body)
		var result fiber.Map
		json.Unmarshal(bodyBytes, &result)

		assert.Equal(t, float64(1), result["total"])
		assert.Equal(t, float64(1), result["page"])
		assert.Equal(t, float64(1), result["last_page"])
		data, ok := result["data"].([]interface{})
		assert.True(t, ok)
		assert.Len(t, data, 1)
		if len(data) > 0 {
			firstLog, ok := data[0].(map[string]interface{})
			assert.True(t, ok)
			assert.Equal(t, expectedLogs[0].ID, firstLog["id"])
		}
		mockRepo.AssertExpectations(t)
	})

	t.Run("DefaultPagination", func(t *testing.T) {
		mockRepo := new(MockLogsRepository)
		app := setupAppAndHandler(mockRepo)

		mockRepo.On("GetLogs", 1, 10).Return([]models.ActivityLog{}, int64(0), nil)

		req := httptest.NewRequest(http.MethodGet, "/api/activity-logs", nil)
		resp, _ := app.Test(req)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)
		mockRepo.AssertExpectations(t)
	})

	t.Run("InvalidPaginationParams", func(t *testing.T) {
		mockRepo := new(MockLogsRepository)
		app := setupAppAndHandler(mockRepo)

		mockRepo.On("GetLogs", 1, 10).Return([]models.ActivityLog{}, int64(0), nil).Once()

		req := httptest.NewRequest(http.MethodGet, "/api/activity-logs?page=abc&limit=xyz", nil)
		resp, _ := app.Test(req)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)
		mockRepo.AssertExpectations(t)
	})

	t.Run("RepositoryError", func(t *testing.T) {
		mockRepo := new(MockLogsRepository)
		app := setupAppAndHandler(mockRepo)

		mockRepo.On("GetLogs", 1, 10).Return([]models.ActivityLog{}, int64(0), errors.New("db error"))

		req := httptest.NewRequest(http.MethodGet, "/api/activity-logs", nil)
		resp, _ := app.Test(req)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
		bodyBytes, _ := io.ReadAll(resp.Body)
		var result fiber.Map
		json.Unmarshal(bodyBytes, &result)
		assert.Equal(t, "Failed to retrieve activity logs", result["error"])
		mockRepo.AssertExpectations(t)
	})
}

func TestGetFilteredActivityLogsHandler(t *testing.T) {
	t.Run("SuccessfulRetrievalWithFilters", func(t *testing.T) {
		mockRepo := new(MockLogsRepository)
		app := setupAppAndHandler(mockRepo)

		expectedLogs := []models.ActivityLog{{ID: "1", User: "filter_user", Action: "LOGIN"}}
		startDate := time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)
		endDate := time.Date(2023, 1, 31, 0, 0, 0, 0, time.UTC)

		mockRepo.On("GetBasedOnFilter", 1, 5, "filter_user", "LOGIN", "SUCCESS", &startDate, &endDate).
			Return(expectedLogs, int64(1), nil)

		url := fmt.Sprintf("/api/activity-logs/filter?page=1&limit=5&user=filter_user&action=LOGIN&status=SUCCESS&startDate=%s&endDate=%s",
			startDate.Format("2006-01-02"), endDate.Format("2006-01-02"))
		req := httptest.NewRequest(http.MethodGet, url, nil)
		resp, _ := app.Test(req)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)
		bodyBytes, _ := io.ReadAll(resp.Body)
		var result fiber.Map
		json.Unmarshal(bodyBytes, &result)
		assert.Equal(t, float64(1), result["total"])
		data, _ := result["data"].([]interface{})
		assert.Len(t, data, 1)
		mockRepo.AssertExpectations(t)
	})

	t.Run("InvalidDateFormat_StartDate", func(t *testing.T) {
		mockRepo := new(MockLogsRepository)
		app := setupAppAndHandler(mockRepo)

		req := httptest.NewRequest(http.MethodGet, "/api/activity-logs/filter?startDate=invalid-date", nil)
		resp, _ := app.Test(req)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
		bodyBytes, _ := io.ReadAll(resp.Body)
		var result fiber.Map
		json.Unmarshal(bodyBytes, &result)
		assert.Contains(t, result["error"], "Invalid startDate format")
	})

	t.Run("InvalidDateFormat_EndDate", func(t *testing.T) {
		mockRepo := new(MockLogsRepository)
		app := setupAppAndHandler(mockRepo)

		req := httptest.NewRequest(http.MethodGet, "/api/activity-logs/filter?endDate=invalid-date", nil)
		resp, _ := app.Test(req)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
		bodyBytes, _ := io.ReadAll(resp.Body)
		var result fiber.Map
		json.Unmarshal(bodyBytes, &result)
		assert.Contains(t, result["error"], "Invalid endDate format")
	})

	t.Run("RepositoryError_WithFilter", func(t *testing.T) {
		mockRepo := new(MockLogsRepository)
		app := setupAppAndHandler(mockRepo)

		mockRepo.On("GetBasedOnFilter", 1, 10, "", "", "", (*time.Time)(nil), (*time.Time)(nil)).
			Return([]models.ActivityLog{}, int64(0), errors.New("db filter error"))

		req := httptest.NewRequest(http.MethodGet, "/api/activity-logs/filter", nil)
		resp, _ := app.Test(req)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
		bodyBytes, _ := io.ReadAll(resp.Body)
		var result fiber.Map
		json.Unmarshal(bodyBytes, &result)
		assert.Equal(t, "Failed to retrieve filtered activity logs", result["error"])
		mockRepo.AssertExpectations(t)
	})

	t.Run("SuccessfulRetrieval_NoFilters", func(t *testing.T) {
		mockRepo := new(MockLogsRepository)
		app := setupAppAndHandler(mockRepo)

		expectedLogs := []models.ActivityLog{{ID: "2", User: "another_user"}}
		mockRepo.On("GetBasedOnFilter", 1, 10, "", "", "", (*time.Time)(nil), (*time.Time)(nil)).
			Return(expectedLogs, int64(1), nil)

		req := httptest.NewRequest(http.MethodGet, "/api/activity-logs/filter?page=1&limit=10", nil)
		resp, _ := app.Test(req)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)
		bodyBytes, _ := io.ReadAll(resp.Body)
		var result fiber.Map
		json.Unmarshal(bodyBytes, &result)
		assert.Equal(t, float64(1), result["total"])
		data, _ := result["data"].([]interface{})
		assert.Len(t, data, 1)
		mockRepo.AssertExpectations(t)
	})
}

func TestCreateActivityLogHandler(t *testing.T) {
	t.Run("SuccessfulCreation", func(t *testing.T) {
		mockRepo := new(MockLogsRepository)
		app := setupAppAndHandler(mockRepo)

		logInput := models.ActivityLog{
			User:           "test_user_create",
			Action:         "CREATE_ITEM",
			Details:        "Item X was created",
			Status:         "successful",
			IsSystemAction: false,
		}
		mockRepo.On("Create", mock.MatchedBy(func(log *models.ActivityLog) bool {
			return log.User == logInput.User && log.Action == logInput.Action
		})).Return(nil).Run(func(args mock.Arguments) {
			arg := args.Get(0).(*models.ActivityLog)
			arg.ID = "new_log_id"
			arg.Timestamp = time.Now()
		})

		bodyBytes, _ := json.Marshal(logInput)
		req := httptest.NewRequest(http.MethodPost, "/api/activity-logs", bytes.NewReader(bodyBytes))
		req.Header.Set("Content-Type", "application/json")

		resp, _ := app.Test(req)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusCreated, resp.StatusCode)

		var createdLog models.ActivityLog
		json.NewDecoder(resp.Body).Decode(&createdLog)

		assert.Equal(t, logInput.User, createdLog.User)
		assert.Equal(t, logInput.Action, createdLog.Action)
		assert.Equal(t, "new_log_id", createdLog.ID)
		assert.False(t, createdLog.Timestamp.IsZero())
		mockRepo.AssertExpectations(t)
	})

	t.Run("InvalidRequestBody_BadJSON", func(t *testing.T) {
		mockRepo := new(MockLogsRepository)
		app := setupAppAndHandler(mockRepo)

		req := httptest.NewRequest(http.MethodPost, "/api/activity-logs", strings.NewReader("not a json string"))
		req.Header.Set("Content-Type", "application/json")

		resp, _ := app.Test(req)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
		var bodyBytesRead []byte
		bodyBytesRead, _ = io.ReadAll(resp.Body)
		var result fiber.Map
		json.Unmarshal(bodyBytesRead, &result)
		assert.Contains(t, result["error"], "Invalid request payload")
	})

	t.Run("MissingRequiredFields", func(t *testing.T) {
		mockRepo := new(MockLogsRepository)
		app := setupAppAndHandler(mockRepo)

		logInput := models.ActivityLog{
			User: "test_user",
		}
		bodyBytes, _ := json.Marshal(logInput)
		req := httptest.NewRequest(http.MethodPost, "/api/activity-logs", bytes.NewReader(bodyBytes))
		req.Header.Set("Content-Type", "application/json")

		resp, _ := app.Test(req)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
		var bodyBytesRead []byte
		bodyBytesRead, _ = io.ReadAll(resp.Body)
		var result fiber.Map
		json.Unmarshal(bodyBytesRead, &result)
		assert.Equal(t, "Missing required fields: user, action, status", result["error"])
	})

	t.Run("RepositoryErrorOnCreate", func(t *testing.T) {
		mockRepo := new(MockLogsRepository)
		app := setupAppAndHandler(mockRepo)

		logInput := models.ActivityLog{
			User:   "test_user_fail",
			Action: "CREATE_FAIL",
			Status: "attempted",
		}
		mockRepo.On("Create", mock.AnythingOfType("*models.ActivityLog")).Return(errors.New("database insert error"))

		bodyBytes, _ := json.Marshal(logInput)
		req := httptest.NewRequest(http.MethodPost, "/api/activity-logs", bytes.NewReader(bodyBytes))
		req.Header.Set("Content-Type", "application/json")

		resp, _ := app.Test(req)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
		var bodyBytesRead []byte
		bodyBytesRead, _ = io.ReadAll(resp.Body)
		var result fiber.Map
		json.Unmarshal(bodyBytesRead, &result)
		assert.Contains(t, result["error"], "Failed to create activity log: database insert error")
		mockRepo.AssertExpectations(t)
	})
}
