package repositories

import (
	"database/sql"
	"log"
	"time"

	"oop/internal/models"
)

// MaterialRepository defines the interface for material database operations
type MaterialRepository interface {
	GetAll(searchTerm string, category string, supplier string, status string) ([]models.Material, error)
	GetByID(id int) (*models.Material, error)
	Create(material *models.Material) (int, error)
	Update(material *models.Material) error
	Delete(id int) error
}

// materialRepository implements the MaterialRepository interface
type materialRepository struct {
	DB *sql.DB
}

// NewMaterialRepository creates a new instance of materialRepository
func NewMaterialRepository(db *sql.DB) MaterialRepository {
	return &materialRepository{DB: db}
}

// GetAll retrieves all materials from the database, with optional filtering
func (r *materialRepository) GetAll(searchTerm string, category string, supplier string, status string) ([]models.Material, error) {
	query := `SELECT id, name, category, supplier, quantity, status, image, created_at, updated_at FROM materials WHERE 1=1`
	args := []interface{}{}

	if searchTerm != "" {
		query += " AND name ILIKE ?"
		args = append(args, "%"+searchTerm+"%")
	}
	if category != "" {
		query += " AND category = ?"
		args = append(args, category)
	}
	if supplier != "" {
		query += " AND supplier = ?"
		args = append(args, supplier)
	}
	if status != "" {
		query += " AND status = ?"
		args = append(args, status)
	}

	query += " ORDER BY created_at DESC"

	rows, err := r.DB.Query(query, args...)
	if err != nil {
		log.Printf("Error querying materials: %v\nQuery: %s\nArgs: %v", err, query, args)
		return nil, err
	}
	defer rows.Close()

	materials := []models.Material{}
	for rows.Next() {
		var m models.Material
		if err := rows.Scan(&m.ID, &m.Name, &m.Category, &m.Supplier, &m.Quantity, &m.Status, &m.Image, &m.CreatedAt, &m.UpdatedAt); err != nil {
			log.Printf("Error scanning material row: %v", err)
			return nil, err
		}
		materials = append(materials, m)
	}

	if err = rows.Err(); err != nil {
		log.Printf("Error iterating material rows: %v", err)
		return nil, err
	}

	return materials, nil
}

// GetByID retrieves a single material by its ID
func (r *materialRepository) GetByID(id int) (*models.Material, error) {
	query := `SELECT id, name, category, supplier, quantity, status, image, created_at, updated_at FROM materials WHERE id = ?`
	row := r.DB.QueryRow(query, id)

	var m models.Material
	if err := row.Scan(&m.ID, &m.Name, &m.Category, &m.Supplier, &m.Quantity, &m.Status, &m.Image, &m.CreatedAt, &m.UpdatedAt); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // Return nil if no material found
		}
		log.Printf("Error scanning material row by ID %d: %v", id, err)
		return nil, err
	}
	return &m, nil
}

// Create inserts a new material into the database
func (r *materialRepository) Create(material *models.Material) (int, error) {
	query := `INSERT INTO materials (name, category, supplier, quantity, status, image, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?)`
	now := time.Now()
	res, err := r.DB.Exec(query, material.Name, material.Category, material.Supplier, material.Quantity, material.Status, material.Image, now, now)
	if err != nil {
		log.Printf("Error creating material: %v", err)
		return 0, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		log.Printf("Error getting last insert ID for material: %v", err)
		return 0, err
	}

	return int(id), nil
}

// Update modifies an existing material in the database
func (r *materialRepository) Update(material *models.Material) error {
	query := `UPDATE materials SET name = ?, category = ?, supplier = ?, quantity = ?, status = ?, image = ?, updated_at = ? WHERE id = ?`
	now := time.Now()
	_, err := r.DB.Exec(query, material.Name, material.Category, material.Supplier, material.Quantity, material.Status, material.Image, now, material.ID)
	if err != nil {
		log.Printf("Error updating material ID %d: %v", material.ID, err)
		return err
	}
	return nil
}

// Delete removes a material from the database by its ID
func (r *materialRepository) Delete(id int) error {
	query := `DELETE FROM materials WHERE id = ?`
	_, err := r.DB.Exec(query, id)
	if err != nil {
		log.Printf("Error deleting material ID %d: %v", id, err)
		return err
	}
	return nil
}
