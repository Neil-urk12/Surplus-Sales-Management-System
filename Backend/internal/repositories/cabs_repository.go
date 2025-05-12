package repositories

import (
	"database/sql"
	"fmt"
	"log"
	"oop/internal/models"
	"oop/internal/config"
	"strings"
	"time"
)

// CabsRepository defines the interface for cab data operations.
type CabsRepository interface {
	GetCabs(filters map[string]interface{}) ([]models.MultiCab, error)
	GetCabByID(id int) (*models.MultiCab, error)
	AddCab(cab models.MultiCab) (*models.MultiCab, error)
	UpdateCab(id int, cab models.MultiCab) (*models.MultiCab, error)
	DeleteCab(id int) error
}

// cabsRepository is a database implementation of CabsRepository.
type cabsRepository struct {
	DB *sql.DB
}

// NewCabsRepository creates a new instance of the database repository.
func NewCabsRepository(db *sql.DB) CabsRepository {
	return &cabsRepository{DB: db}
}

// GetCabs retrieves a list of cabs, applying filters if provided.
func (r *cabsRepository) GetCabs(filters map[string]interface{}) ([]models.MultiCab, error) {
	query := `SELECT id, name, make, quantity, price, status, unit_color, image, created_at, updated_at FROM multicabs WHERE 1=1`
	var args []interface{}

	// Apply filters to query
	if makeFilter, ok := filters["make"].(string); ok && makeFilter != "" {
		query += " AND make = ?"
		args = append(args, makeFilter)
	}

	if colorFilter, ok := filters["unit_color"].(string); ok && colorFilter != "" {
		query += " AND unit_color = ?"
		args = append(args, colorFilter)
	}

	if statusFilter, ok := filters["status"].(string); ok && statusFilter != "" {
		query += " AND status = ?"
		args = append(args, statusFilter)
	}

	if searchFilter, ok := filters["search"].(string); ok && searchFilter != "" {
		query += " AND (name LIKE ? OR make LIKE ?)"
		searchTerm := "%" + searchFilter + "%"
		args = append(args, searchTerm, searchTerm)
	}

	query += " ORDER BY created_at DESC"

	rows, err := r.DB.Query(query, args...)
	if err != nil {
		log.Printf("Error querying cabs: %v\nQuery: %s\nArgs: %v", err, query, args)
		return nil, err
	}
	defer rows.Close()

	var cabs []models.MultiCab
	for rows.Next() {
		var cab models.MultiCab
		var createdAt, updatedAt time.Time
		var imageSQL sql.NullString

		if err := rows.Scan(
			&cab.ID,
			&cab.Name,
			&cab.Make,
			&cab.Quantity,
			&cab.Price,
			&cab.Status,
			&cab.UnitColor,
			&imageSQL,
			&createdAt,
			&updatedAt,
		); err != nil {
			log.Printf("Error scanning cab row: %v", err)
			return nil, err
		}

		cab.CreatedAt = createdAt
		cab.UpdatedAt = updatedAt
		
		// Handle NULL image values with default image
		if imageSQL.Valid && imageSQL.String != "" {
			cab.Image = imageSQL.String
		} else {
			cab.Image = config.DefaultImageURL
		}
		
		cabs = append(cabs, cab)
	}

	if err = rows.Err(); err != nil {
		log.Printf("Error iterating cab rows: %v", err)
		return nil, err
	}

	return cabs, nil
}

// GetCabByID retrieves a single cab by its ID.
func (r *cabsRepository) GetCabByID(id int) (*models.MultiCab, error) {
	query := `SELECT id, name, make, quantity, price, status, unit_color, image, created_at, updated_at FROM multicabs WHERE id = ?`
	row := r.DB.QueryRow(query, id)

	var cab models.MultiCab
	var createdAt, updatedAt time.Time
	var imageSQL sql.NullString

	if err := row.Scan(
		&cab.ID,
		&cab.Name,
		&cab.Make,
		&cab.Quantity,
		&cab.Price,
		&cab.Status,
		&cab.UnitColor,
		&imageSQL,
		&createdAt,
		&updatedAt,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("cab with ID %d not found", id)
		}
		log.Printf("Error scanning cab row by ID %d: %v", id, err)
		return nil, err
	}

	cab.CreatedAt = createdAt
	cab.UpdatedAt = updatedAt
	
	// Handle NULL image values with default image
	if imageSQL.Valid && imageSQL.String != "" {
		cab.Image = imageSQL.String
	} else {
		cab.Image = config.DefaultImageURL
	}
	
	return &cab, nil
}

// AddCab adds a new cab to the repository.
func (r *cabsRepository) AddCab(cab models.MultiCab) (*models.MultiCab, error) {
	// Basic validation
	if cab.Name == "" || cab.Make == "" || cab.UnitColor == "" {
		return nil, fmt.Errorf("cab name, make, and color cannot be empty")
	}

	query := `INSERT INTO multicabs (name, make, quantity, price, status, unit_color, image, created_at, updated_at) 
              VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`

	now := time.Now()
	
	var imageValue interface{}
	if cab.Image == "" || cab.Image == config.DefaultImageURL {
		imageValue = nil // Use NULL if image is empty or default
	} else {
		imageValue = cab.Image
	}
	
	result, err := r.DB.Exec(
		query,
		cab.Name,
		cab.Make,
		cab.Quantity,
		cab.Price,
		cab.Status,
		cab.UnitColor,
		imageValue,
		now,
		now,
	)

	if err != nil {
		log.Printf("Error adding cab: %v", err)
		return nil, err
	}

	// Get the auto-generated ID
	id, err := result.LastInsertId()
	if err != nil {
		log.Printf("Error getting last insert ID for cab: %v", err)
		return nil, err
	}

	// Set the generated fields on the cab
	cab.ID = int(id)
	cab.CreatedAt = now
	cab.UpdatedAt = now
	
	// Ensure image is set to default if it was empty
	if cab.Image == "" {
		cab.Image = config.DefaultImageURL
	}

	return &cab, nil
}

// UpdateCab updates an existing cab in the repository.
func (r *cabsRepository) UpdateCab(id int, cab models.MultiCab) (*models.MultiCab, error) {
	// Check if the cab exists
	existingCab, err := r.GetCabByID(id)
	if err != nil {
		if strings.Contains(strings.ToLower(err.Error()), "not found") {
			return nil, fmt.Errorf("cab with ID %d not found for update", id)
		}
		return nil, err
	}

	// Prepare the update query
	query := `UPDATE multicabs 
              SET name = ?, make = ?, quantity = ?, price = ?, status = ?, unit_color = ?, image = ?, updated_at = ? 
              WHERE id = ?`

	now := time.Now()
	
	var imageValue interface{}
	if cab.Image == "" || cab.Image == config.DefaultImageURL {
		imageValue = nil // Use NULL if image is empty or default
	} else {
		imageValue = cab.Image
	}
	
	_, err = r.DB.Exec(
		query,
		cab.Name,
		cab.Make,
		cab.Quantity,
		cab.Price,
		cab.Status,
		cab.UnitColor,
		imageValue,
		now,
		id,
	)

	if err != nil {
		log.Printf("Error updating cab ID %d: %v", id, err)
		return nil, err
	}

	// Update the cab with the current data
	cab.ID = id
	cab.CreatedAt = existingCab.CreatedAt
	cab.UpdatedAt = now
	
	// Ensure image is set to default if it was empty
	if cab.Image == "" {
		cab.Image = config.DefaultImageURL
	}

	return &cab, nil
}

// DeleteCab removes a cab from the repository by its ID.
func (r *cabsRepository) DeleteCab(id int) error {
	// Check if the cab exists
	_, err := r.GetCabByID(id)
	if err != nil {
		if strings.Contains(strings.ToLower(err.Error()), "not found") {
			return fmt.Errorf("cab with ID %d not found for deletion", id)
		}
		return err
	}

	query := `DELETE FROM multicabs WHERE id = ?`
	_, err = r.DB.Exec(query, id)
	if err != nil {
		log.Printf("Error deleting cab ID %d: %v", id, err)
		return err
	}

	return nil
}
