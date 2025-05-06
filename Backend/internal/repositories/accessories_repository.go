package repositories

import (
	"context"
	"database/sql"
	"errors"
	"oop/internal/models"
)

// Helper function to determine status based on quantity
func determineStatus(quantity int) models.AccessoryStatus {
	switch {
	case quantity == 0:
		return models.StatusOutOfStock
	case quantity <= 2:
		return models.StatusLowStock
	case quantity <= 5:
		return models.StatusInStock
	default:
		return models.StatusAvailable
	}
}

// AccessoryRepository defines methods for working with accessories
type AccessoryRepository interface {
	GetAll(ctx context.Context) ([]models.Accessory, error)
	GetByID(ctx context.Context, id int) (models.Accessory, error)
	Create(ctx context.Context, input models.NewAccessoryInput) (models.Accessory, error)
	Update(ctx context.Context, id int, input models.UpdateAccessoryInput) (models.Accessory, error)
	Delete(ctx context.Context, id int) error
}

// AccessoryRepositoryImpl is a SQL implementation of AccessoryRepository
type AccessoryRepositoryImpl struct {
	DB *sql.DB
}

// NewAccessoryRepository creates a new accessory repository
func NewAccessoryRepository(db *sql.DB) AccessoryRepository {
	return &AccessoryRepositoryImpl{
		DB: db,
	}
}

// GetAll retrieves all accessories from the database
func (r *AccessoryRepositoryImpl) GetAll(ctx context.Context) ([]models.Accessory, error) {
	query := `
		SELECT id, name, make, quantity, price, status, unit_color, image, created_at, updated_at
		FROM accessories
		ORDER BY id ASC
	`

	rows, err := r.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var accessories []models.Accessory
	for rows.Next() {
		var a models.Accessory
		var makeStr, colorStr, statusStr string

		err := rows.Scan(
			&a.ID,
			&a.Name,
			&makeStr,
			&a.Quantity,
			&a.Price,
			&statusStr,
			&colorStr,
			&a.Image,
			&a.CreatedAt,
			&a.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		a.Make = models.AccessoryMake(makeStr)
		a.UnitColor = models.AccessoryColor(colorStr)
		a.Status = models.AccessoryStatus(statusStr)

		accessories = append(accessories, a)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return accessories, nil
}

// GetByID retrieves an accessory by its ID
func (r *AccessoryRepositoryImpl) GetByID(ctx context.Context, id int) (models.Accessory, error) {
	query := `
		SELECT id, name, make, quantity, price, status, unit_color, image, created_at, updated_at
		FROM accessories
		WHERE id = $1
	`

	var a models.Accessory
	var makeStr, colorStr, statusStr string

	err := r.DB.QueryRowContext(ctx, query, id).Scan(
		&a.ID,
		&a.Name,
		&makeStr,
		&a.Quantity,
		&a.Price,
		&statusStr,
		&colorStr,
		&a.Image,
		&a.CreatedAt,
		&a.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.Accessory{}, errors.New("accessory not found")
		}
		return models.Accessory{}, err
	}

	a.Make = models.AccessoryMake(makeStr)
	a.UnitColor = models.AccessoryColor(colorStr)
	a.Status = models.AccessoryStatus(statusStr)

	return a, nil
}

// Create inserts a new accessory into the database
func (r *AccessoryRepositoryImpl) Create(ctx context.Context, input models.NewAccessoryInput) (models.Accessory, error) {
	status := determineStatus(input.Quantity)

	query := `
		INSERT INTO accessories (name, make, quantity, price, status, unit_color, image, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, NOW(), NOW())
		RETURNING id, created_at, updated_at
	`

	var accessory models.Accessory
	accessory.Name = input.Name
	accessory.Make = input.Make
	accessory.Quantity = input.Quantity
	accessory.Price = input.Price
	accessory.Status = status
	accessory.UnitColor = input.UnitColor
	accessory.Image = input.Image

	err := r.DB.QueryRowContext(
		ctx,
		query,
		accessory.Name,
		string(accessory.Make),
		accessory.Quantity,
		accessory.Price,
		string(accessory.Status),
		string(accessory.UnitColor),
		accessory.Image,
	).Scan(&accessory.ID, &accessory.CreatedAt, &accessory.UpdatedAt)

	if err != nil {
		return models.Accessory{}, err
	}

	return accessory, nil
}

// Update modifies an existing accessory in the database
func (r *AccessoryRepositoryImpl) Update(ctx context.Context, id int, input models.UpdateAccessoryInput) (models.Accessory, error) {
	// First, get the current accessory
	accessory, err := r.GetByID(ctx, id)
	if err != nil {
		return models.Accessory{}, err
	}

	// Apply updates if provided
	if input.Name != nil {
		accessory.Name = *input.Name
	}
	if input.Make != nil {
		accessory.Make = *input.Make
	}
	if input.Quantity != nil {
		accessory.Quantity = *input.Quantity
		accessory.Status = determineStatus(*input.Quantity)
	}
	if input.Price != nil {
		accessory.Price = *input.Price
	}
	if input.UnitColor != nil {
		accessory.UnitColor = *input.UnitColor
	}
	if input.Image != nil {
		accessory.Image = *input.Image
	}

	query := `
		UPDATE accessories
		SET name = $1, make = $2, quantity = $3, price = $4, status = $5, unit_color = $6, image = $7, updated_at = NOW()
		WHERE id = $8
		RETURNING updated_at
	`

	err = r.DB.QueryRowContext(
		ctx,
		query,
		accessory.Name,
		string(accessory.Make),
		accessory.Quantity,
		accessory.Price,
		string(accessory.Status),
		string(accessory.UnitColor),
		accessory.Image,
		id,
	).Scan(&accessory.UpdatedAt)

	if err != nil {
		return models.Accessory{}, err
	}

	return accessory, nil
}

// Delete removes an accessory from the database
func (r *AccessoryRepositoryImpl) Delete(ctx context.Context, id int) error {
	query := `DELETE FROM accessories WHERE id = $1`

	result, err := r.DB.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New("accessory not found")
	}

	return nil
}
