package repositories

import (
	"database/sql"
	"fmt"
	"log"
	"oop/internal/models"
	"strings"
	"time"

	"github.com/google/uuid"
)

// LogsRepositoryInterface defines the interface for activity log database operations
type LogsRepositoryInterface interface {
	Create(log *models.ActivityLog) error
	GetLogs(page, limit int) ([]models.ActivityLog, int64, error)
	GetBasedOnFilter(page, limit int, user, action, status string, startDate, endDate *time.Time) ([]models.ActivityLog, int64, error)
}

// LogsRepository handles database operations related to users
type LogsRepository struct {
	dbClient *sql.DB // Changed from *DatabaseClient to *sql.DB assuming it's a standard SQL database client
}

// NewLogsRepository creates a new LogsRepository instance
func NewLogsRepository(dbClient *sql.DB) LogsRepositoryInterface {
	return &LogsRepository{
		dbClient: dbClient,
	}
}

func (r *LogsRepository) Create(logEntry *models.ActivityLog) error {
	// Generate a new UUID for the user if not provided
	if logEntry.ID == "" {
		logEntry.ID = uuid.New().String()
	}
	now := time.Now()
	logEntry.CreatedAt = now
	logEntry.UpdatedAt = now
	if logEntry.Timestamp.IsZero() {
		logEntry.Timestamp = now
	}

	query := `INSERT INTO activity_logs (id, timestamp, user_id, action_type, details, status, is_system_action, created_at, updated_at)
	          VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`
	_, err := r.dbClient.Exec(query, logEntry.ID, logEntry.Timestamp, logEntry.User, logEntry.Action, logEntry.Details, logEntry.Status, logEntry.IsSystemAction, logEntry.CreatedAt, logEntry.UpdatedAt)
	if err != nil {
		log.Printf("Error creating activity log: %v", err)
		return fmt.Errorf("could not create activity log: %w", err)
	}
	return nil
}

// GetLogs retrieves paginated activity logs.
func (r *LogsRepository) GetLogs(page, limit int) ([]models.ActivityLog, int64, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10 // Default limit
	}
	offset := (page - 1) * limit

	query := `SELECT id, timestamp, user_id, action_type, details, status, is_system_action, created_at, updated_at
	          FROM activity_logs ORDER BY timestamp DESC LIMIT ? OFFSET ?`
	countQuery := `SELECT COUNT(*) FROM activity_logs`

	var total int64
	err := r.dbClient.QueryRow(countQuery).Scan(&total)
	if err != nil {
		log.Printf("Error counting activity logs: %v", err)
		return nil, 0, fmt.Errorf("could not count activity logs: %w", err)
	}

	rows, err := r.dbClient.Query(query, limit, offset)
	if err != nil {
		log.Printf("Error querying activity logs: %v", err)
		return nil, 0, fmt.Errorf("could not query activity logs: %w", err)
	}
	defer rows.Close()

	logs := []models.ActivityLog{}
	for rows.Next() {
		var l models.ActivityLog
		if err := rows.Scan(&l.ID, &l.Timestamp, &l.User, &l.Action, &l.Details, &l.Status, &l.IsSystemAction, &l.CreatedAt, &l.UpdatedAt); err != nil {
			log.Printf("Error scanning activity log row: %v", err)
			return nil, 0, fmt.Errorf("could not scan activity log: %w", err)
		}
		logs = append(logs, l)
	}

	if err = rows.Err(); err != nil {
		log.Printf("Error iterating activity log rows: %v", err)
		return nil, 0, fmt.Errorf("error iterating activity log rows: %w", err)
	}

	return logs, total, nil
}

// GetBasedOnFilter retrieves paginated activity logs based on specified filters.
func (r *LogsRepository) GetBasedOnFilter(page, limit int, user, action, status string, startDate, endDate *time.Time) ([]models.ActivityLog, int64, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10 // Default limit
	}
	offset := (page - 1) * limit

	var queryBuilder strings.Builder
	var countQueryBuilder strings.Builder
	args := []interface{}{}
	paramIndex := 1

	queryBuilder.WriteString("SELECT id, timestamp, user_id, action_type, details, status, is_system_action, created_at, updated_at FROM activity_logs WHERE 1=1")
	countQueryBuilder.WriteString("SELECT COUNT(*) FROM activity_logs WHERE 1=1")

	addCondition := func(field, value string) {
		if value != "" {
			condition := fmt.Sprintf(" AND LOWER(%s) LIKE LOWER(?)", field)
			queryBuilder.WriteString(condition)
			countQueryBuilder.WriteString(condition)
			args = append(args, "%"+value+"%")
			paramIndex++
		}
	}

	addCondition("user_id", user)
	addCondition("action_type", action)
	addCondition("status", status)

	if startDate != nil {
		condition := " AND timestamp >= ?"
		queryBuilder.WriteString(condition)
		countQueryBuilder.WriteString(condition)
		args = append(args, *startDate)
		paramIndex++
	}

	if endDate != nil {
		// Use the endDate directly - the handler already adjusts it if needed
		condition := " AND timestamp <= ?"
		queryBuilder.WriteString(condition)
		countQueryBuilder.WriteString(condition)
		args = append(args, *endDate)
		paramIndex++
	}

	queryBuilder.WriteString(" ORDER BY timestamp DESC LIMIT ? OFFSET ?")

	finalArgs := make([]interface{}, len(args))
	copy(finalArgs, args)
	finalArgs = append(finalArgs, limit, offset)

	var total int64
	err := r.dbClient.QueryRow(countQueryBuilder.String(), args...).Scan(&total)
	if err != nil {
		log.Printf("Error counting filtered activity logs: %v\nQuery: %s\nArgs: %v", err, countQueryBuilder.String(), args)
		return nil, 0, fmt.Errorf("could not count filtered activity logs: %w", err)
	}

	rows, err := r.dbClient.Query(queryBuilder.String(), finalArgs...)
	if err != nil {
		log.Printf("Error querying filtered activity logs: %v\nQuery: %s\nArgs: %v", err, queryBuilder.String(), finalArgs)
		return nil, 0, fmt.Errorf("could not query filtered activity logs: %w", err)
	}
	defer rows.Close()

	logs := []models.ActivityLog{}
	for rows.Next() {
		var l models.ActivityLog
		if err := rows.Scan(&l.ID, &l.Timestamp, &l.User, &l.Action, &l.Details, &l.Status, &l.IsSystemAction, &l.CreatedAt, &l.UpdatedAt); err != nil {
			log.Printf("Error scanning filtered activity log row: %v", err)
			return nil, 0, fmt.Errorf("could not scan filtered activity log: %w", err)
		}
		logs = append(logs, l)
	}

	if err = rows.Err(); err != nil {
		log.Printf("Error iterating filtered activity log rows: %v", err)
		return nil, 0, fmt.Errorf("error iterating filtered activity log rows: %w", err)
	}

	return logs, total, nil
}
