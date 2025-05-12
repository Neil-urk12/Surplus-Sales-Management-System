package repositories

import (
	"database/sql/driver"
	"fmt"
	"oop/internal/models"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestCreateActivityLog(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := NewLogsRepository(db)

	now := time.Now()
	logEntry := &models.ActivityLog{
		User:           "testuser",
		Action:         "LOGIN",
		Details:        "User logged in",
		Status:         "SUCCESS",
		IsSystemAction: false,
		Timestamp:      now,
	}

	t.Run("SuccessfulCreate", func(t *testing.T) {
		mock.ExpectExec(regexp.QuoteMeta("INSERT INTO activity_logs (id, timestamp, user_id, action_type, details, status, is_system_action, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)")).
			WithArgs(sqlmock.AnyArg(), logEntry.Timestamp, logEntry.User, logEntry.Action, logEntry.Details, logEntry.Status, logEntry.IsSystemAction, sqlmock.AnyArg(), sqlmock.AnyArg()).
			WillReturnResult(sqlmock.NewResult(1, 1))

		err := repo.Create(logEntry)
		assert.NoError(t, err)
		assert.NotEmpty(t, logEntry.ID)
		assert.WithinDuration(t, now, logEntry.CreatedAt, time.Second)
		assert.WithinDuration(t, now, logEntry.UpdatedAt, time.Second)
	})

	t.Run("SuccessfulCreate_WithExistingID", func(t *testing.T) {
		existingID := "existing-uuid"
		logEntryWithID := &models.ActivityLog{
			ID:             existingID,
			User:           "testuser2",
			Action:         "UPDATE_PROFILE",
			Details:        "Profile updated",
			Status:         "SUCCESS",
			IsSystemAction: false,
			Timestamp:      now,
		}
		mock.ExpectExec(regexp.QuoteMeta("INSERT INTO activity_logs (id, timestamp, user_id, action_type, details, status, is_system_action, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)")).
			WithArgs(existingID, logEntryWithID.Timestamp, logEntryWithID.User, logEntryWithID.Action, logEntryWithID.Details, logEntryWithID.Status, logEntryWithID.IsSystemAction, sqlmock.AnyArg(), sqlmock.AnyArg()).
			WillReturnResult(sqlmock.NewResult(1, 1))

		err := repo.Create(logEntryWithID)
		assert.NoError(t, err)
		assert.Equal(t, existingID, logEntryWithID.ID)
	})

	t.Run("DatabaseError", func(t *testing.T) {
		mock.ExpectExec(regexp.QuoteMeta("INSERT INTO activity_logs (id, timestamp, user_id, action_type, details, status, is_system_action, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)")).
			WillReturnError(fmt.Errorf("db error"))

		err := repo.Create(logEntry)
		assert.Error(t, err)
	})

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetLogs(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := NewLogsRepository(db)
	now := time.Now()

	t.Run("SuccessfulGetLogs", func(t *testing.T) {
		expectedLogs := []models.ActivityLog{
			{ID: "uuid1", User: "user1", Action: "ACTION1", Timestamp: now, CreatedAt: now, UpdatedAt: now},
			{ID: "uuid2", User: "user2", Action: "ACTION2", Timestamp: now.Add(-time.Hour), CreatedAt: now.Add(-time.Hour), UpdatedAt: now.Add(-time.Hour)},
		}
		rows := sqlmock.NewRows([]string{"id", "timestamp", "user_id", "action_type", "details", "status", "is_system_action", "created_at", "updated_at"}).
			AddRow(expectedLogs[0].ID, expectedLogs[0].Timestamp, expectedLogs[0].User, expectedLogs[0].Action, "", "", false, expectedLogs[0].CreatedAt, expectedLogs[0].UpdatedAt).
			AddRow(expectedLogs[1].ID, expectedLogs[1].Timestamp, expectedLogs[1].User, expectedLogs[1].Action, "", "", false, expectedLogs[1].CreatedAt, expectedLogs[1].UpdatedAt)

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT COUNT(*) FROM activity_logs`)).
			WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(2))
		mock.ExpectQuery(regexp.QuoteMeta("SELECT id, timestamp, user_id, action_type, details, status, is_system_action, created_at, updated_at FROM activity_logs ORDER BY timestamp DESC LIMIT ? OFFSET ?")).
			WithArgs(10, 0).
			WillReturnRows(rows)

		logs, total, err := repo.GetLogs(1, 10)
		assert.NoError(t, err)
		assert.Equal(t, int64(2), total)
		assert.Len(t, logs, 2)
		assert.Equal(t, expectedLogs, logs)
	})

	t.Run("NoLogsFound", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT COUNT(*) FROM activity_logs`)).
			WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(0))
		mock.ExpectQuery(regexp.QuoteMeta("SELECT id, timestamp, user_id, action_type, details, status, is_system_action, created_at, updated_at FROM activity_logs ORDER BY timestamp DESC LIMIT ? OFFSET ?")).
			WithArgs(10, 0).
			WillReturnRows(sqlmock.NewRows([]string{"id", "timestamp", "user_id", "action_type", "details", "status", "is_system_action", "created_at", "updated_at"}))

		logs, total, err := repo.GetLogs(1, 10)
		assert.NoError(t, err)
		assert.Equal(t, int64(0), total)
		assert.Empty(t, logs)
	})

	t.Run("CountQueryError", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT COUNT(*) FROM activity_logs`)).
			WillReturnError(fmt.Errorf("count db error"))

		_, _, err := repo.GetLogs(1, 10)
		assert.Error(t, err)
	})

	t.Run("MainQueryError", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT COUNT(*) FROM activity_logs`)).
			WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1)) // Assume count is fine
		mock.ExpectQuery(regexp.QuoteMeta("SELECT id, timestamp, user_id, action_type, details, status, is_system_action, created_at, updated_at FROM activity_logs ORDER BY timestamp DESC LIMIT ? OFFSET ?")).
			WillReturnError(fmt.Errorf("main query db error"))

		_, _, err := repo.GetLogs(1, 10)
		assert.Error(t, err)
	})

	t.Run("ScanError", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"id", "timestamp"}). // Mismatching columns to cause scan error
									AddRow("uuid1", now)

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT COUNT(*) FROM activity_logs`)).
			WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))
		mock.ExpectQuery(regexp.QuoteMeta("SELECT id, timestamp, user_id, action_type, details, status, is_system_action, created_at, updated_at FROM activity_logs ORDER BY timestamp DESC LIMIT ? OFFSET ?")).
			WillReturnRows(rows)

		_, _, err := repo.GetLogs(1, 10)
		assert.Error(t, err)
	})

	assert.NoError(t, mock.ExpectationsWereMet())
}

// AnyTime can be used for time.Time values in mock.ExpectQuery arguments
type AnyTime struct{}

// Match satisfies sqlmock.Argument interface
func (a AnyTime) Match(v driver.Value) bool {
	_, ok := v.(time.Time)
	return ok
}

func TestGetBasedOnFilter(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := NewLogsRepository(db)
	now := time.Now()
	startDate := now.Add(-24 * time.Hour)
	endDate := now

	defaultLog := models.ActivityLog{ID: "uuid1", User: "test_user", Action: "TEST_ACTION", Status: "SUCCESS", Timestamp: now, CreatedAt: now, UpdatedAt: now}

	t.Run("SuccessfulGetBasedOnFilter_AllFilters", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"id", "timestamp", "user_id", "action_type", "details", "status", "is_system_action", "created_at", "updated_at"}).
			AddRow(defaultLog.ID, defaultLog.Timestamp, defaultLog.User, defaultLog.Action, defaultLog.Details, defaultLog.Status, defaultLog.IsSystemAction, defaultLog.CreatedAt, defaultLog.UpdatedAt)

		expectedCountQuery := "SELECT COUNT(*) FROM activity_logs WHERE 1=1 AND LOWER(user_id) LIKE LOWER(?) AND LOWER(action_type) LIKE LOWER(?) AND LOWER(status) LIKE LOWER(?) AND timestamp >= ? AND timestamp <= ?"
		expectedQuery := "SELECT id, timestamp, user_id, action_type, details, status, is_system_action, created_at, updated_at FROM activity_logs WHERE 1=1 AND LOWER(user_id) LIKE LOWER(?) AND LOWER(action_type) LIKE LOWER(?) AND LOWER(status) LIKE LOWER(?) AND timestamp >= ? AND timestamp <= ? ORDER BY timestamp DESC LIMIT ? OFFSET ?"

		mock.ExpectQuery(regexp.QuoteMeta(expectedCountQuery)).
			WithArgs("%test_user%", "%TEST_ACTION%", "%SUCCESS%", AnyTime{}, AnyTime{}).
			WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))
		mock.ExpectQuery(regexp.QuoteMeta(expectedQuery)).
			WithArgs("%test_user%", "%TEST_ACTION%", "%SUCCESS%", AnyTime{}, AnyTime{}, 10, 0).
			WillReturnRows(rows)

		logs, total, err := repo.GetBasedOnFilter(1, 10, "test_user", "TEST_ACTION", "SUCCESS", &startDate, &endDate)
		assert.NoError(t, err)
		assert.Equal(t, int64(1), total)
		assert.Len(t, logs, 1)
		assert.Equal(t, defaultLog, logs[0])
	})

	t.Run("SuccessfulGetBasedOnFilter_OnlyUser", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"id", "timestamp", "user_id", "action_type", "details", "status", "is_system_action", "created_at", "updated_at"}).
			AddRow(defaultLog.ID, defaultLog.Timestamp, defaultLog.User, defaultLog.Action, defaultLog.Details, defaultLog.Status, defaultLog.IsSystemAction, defaultLog.CreatedAt, defaultLog.UpdatedAt)

		expectedCountQuery := "SELECT COUNT(*) FROM activity_logs WHERE 1=1 AND LOWER(user_id) LIKE LOWER(?)"
		expectedQuery := "SELECT id, timestamp, user_id, action_type, details, status, is_system_action, created_at, updated_at FROM activity_logs WHERE 1=1 AND LOWER(user_id) LIKE LOWER(?) ORDER BY timestamp DESC LIMIT ? OFFSET ?"

		mock.ExpectQuery(regexp.QuoteMeta(expectedCountQuery)).
			WithArgs("%test_user%").
			WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))
		mock.ExpectQuery(regexp.QuoteMeta(expectedQuery)).
			WithArgs("%test_user%", 10, 0).
			WillReturnRows(rows)

		logs, total, err := repo.GetBasedOnFilter(1, 10, "test_user", "", "", nil, nil)
		assert.NoError(t, err)
		assert.Equal(t, int64(1), total)
		assert.Len(t, logs, 1)
		assert.Equal(t, defaultLog, logs[0])
	})

	t.Run("SuccessfulGetBasedOnFilter_OnlyDateRange", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"id", "timestamp", "user_id", "action_type", "details", "status", "is_system_action", "created_at", "updated_at"}).
			AddRow(defaultLog.ID, defaultLog.Timestamp, defaultLog.User, defaultLog.Action, defaultLog.Details, defaultLog.Status, defaultLog.IsSystemAction, defaultLog.CreatedAt, defaultLog.UpdatedAt)

		expectedCountQuery := "SELECT COUNT(*) FROM activity_logs WHERE 1=1 AND timestamp >= ? AND timestamp <= ?"
		expectedQuery := "SELECT id, timestamp, user_id, action_type, details, status, is_system_action, created_at, updated_at FROM activity_logs WHERE 1=1 AND timestamp >= ? AND timestamp <= ? ORDER BY timestamp DESC LIMIT ? OFFSET ?"

		mock.ExpectQuery(regexp.QuoteMeta(expectedCountQuery)).
			WithArgs(AnyTime{}, AnyTime{}).
			WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))
		mock.ExpectQuery(regexp.QuoteMeta(expectedQuery)).
			WithArgs(AnyTime{}, AnyTime{}, 10, 0).
			WillReturnRows(rows)

		logs, total, err := repo.GetBasedOnFilter(1, 10, "", "", "", &startDate, &endDate)
		assert.NoError(t, err)
		assert.Equal(t, int64(1), total)
		assert.Len(t, logs, 1)
		assert.Equal(t, defaultLog, logs[0])
	})

	t.Run("NoResultsFound", func(t *testing.T) {
		expectedCountQuery := "SELECT COUNT(*) FROM activity_logs WHERE 1=1 AND LOWER(user_id) LIKE LOWER(?)"
		expectedQuery := "SELECT id, timestamp, user_id, action_type, details, status, is_system_action, created_at, updated_at FROM activity_logs WHERE 1=1 AND LOWER(user_id) LIKE LOWER(?) ORDER BY timestamp DESC LIMIT ? OFFSET ?"

		mock.ExpectQuery(regexp.QuoteMeta(expectedCountQuery)).
			WithArgs("%nonexistent%").
			WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(0))
		mock.ExpectQuery(regexp.QuoteMeta(expectedQuery)).
			WithArgs("%nonexistent%", 10, 0).
			WillReturnRows(sqlmock.NewRows([]string{"id", "timestamp", "user_id", "action_type", "details", "status", "is_system_action", "created_at", "updated_at"}))

		logs, total, err := repo.GetBasedOnFilter(1, 10, "nonexistent", "", "", nil, nil)
		assert.NoError(t, err)
		assert.Equal(t, int64(0), total)
		assert.Empty(t, logs)
	})

	t.Run("CountQueryError_WithFilter", func(t *testing.T) {
		expectedCountQuery := "SELECT COUNT(*) FROM activity_logs WHERE 1=1 AND LOWER(user_id) LIKE LOWER(?)"
		mock.ExpectQuery(regexp.QuoteMeta(expectedCountQuery)).
			WithArgs("%test_user%").
			WillReturnError(fmt.Errorf("filter count db error"))

		_, _, err := repo.GetBasedOnFilter(1, 10, "test_user", "", "", nil, nil)
		assert.Error(t, err)
	})

	t.Run("MainQueryError_WithFilter", func(t *testing.T) {
		expectedCountQuery := "SELECT COUNT(*) FROM activity_logs WHERE 1=1 AND LOWER(user_id) LIKE LOWER(?)"
		expectedQuery := "SELECT id, timestamp, user_id, action_type, details, status, is_system_action, created_at, updated_at FROM activity_logs WHERE 1=1 AND LOWER(user_id) LIKE LOWER(?) ORDER BY timestamp DESC LIMIT ? OFFSET ?"

		mock.ExpectQuery(regexp.QuoteMeta(expectedCountQuery)).
			WithArgs("%test_user%").
			WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1)) // Assume count is fine
		mock.ExpectQuery(regexp.QuoteMeta(expectedQuery)).
			WithArgs("%test_user%", 10, 0).
			WillReturnError(fmt.Errorf("filter main query db error"))

		_, _, err := repo.GetBasedOnFilter(1, 10, "test_user", "", "", nil, nil)
		assert.Error(t, err)
	})

	t.Run("ScanError_WithFilter", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"id", "timestamp"}). // Mismatch columns
									AddRow("uuid1", now)

		expectedCountQuery := "SELECT COUNT(*) FROM activity_logs WHERE 1=1 AND LOWER(user_id) LIKE LOWER(?)"
		expectedQuery := "SELECT id, timestamp, user_id, action_type, details, status, is_system_action, created_at, updated_at FROM activity_logs WHERE 1=1 AND LOWER(user_id) LIKE LOWER(?) ORDER BY timestamp DESC LIMIT ? OFFSET ?"

		mock.ExpectQuery(regexp.QuoteMeta(expectedCountQuery)).
			WithArgs("%test_user%").
			WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))
		mock.ExpectQuery(regexp.QuoteMeta(expectedQuery)).
			WithArgs("%test_user%", 10, 0).
			WillReturnRows(rows)

		_, _, err := repo.GetBasedOnFilter(1, 10, "test_user", "", "", nil, nil)
		assert.Error(t, err)
	})

	assert.NoError(t, mock.ExpectationsWereMet())
}
