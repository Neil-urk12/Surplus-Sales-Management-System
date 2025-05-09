package repositories

import (
	"context"
	"database/sql"
	"errors"
	"oop/internal/models"
)

// Default image to use when no image is provided
const DefaultImageURL = "https://encrypted-tbn0.gstatic.com/images?q=tbn:ANd9GcQU0N_pZ1FmfWhbnKjb-rlqcfOO65_PRLhvTg&s"

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
	Create(ctx context.Context, input models.NewAccessoryInput) (int, error)
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

	stmt, err := r.DB.PrepareContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	rows, err := stmt.QueryContext(ctx)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var accessories []models.Accessory
	for rows.Next() {
		var a models.Accessory
		var makeStr, colorStr, statusStr string
		var imageSQL sql.NullString

		err := rows.Scan(
			&a.ID,
			&a.Name,
			&makeStr,
			&a.Quantity,
			&a.Price,
			&statusStr,
			&colorStr,
			&imageSQL,
			&a.CreatedAt,
			&a.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		a.Make = models.AccessoryMake(makeStr)
		a.UnitColor = models.AccessoryColor(colorStr)
		a.Status = models.AccessoryStatus(statusStr)
		
		// Handle NULL image values with default image
		if imageSQL.Valid && imageSQL.String != "" {
			a.Image = imageSQL.String
		} else {
			a.Image = DefaultImageURL
		}

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
		WHERE id = ?
	`

	stmt, err := r.DB.PrepareContext(ctx, query)
	if err != nil {
		return models.Accessory{}, err
	}
	defer stmt.Close()

	var a models.Accessory
	var makeStr, colorStr, statusStr string
	var imageSQL sql.NullString

	err = stmt.QueryRowContext(ctx, id).Scan(
		&a.ID,
		&a.Name,
		&makeStr,
		&a.Quantity,
		&a.Price,
		&statusStr,
		&colorStr,
		&imageSQL,
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
	
	// Handle NULL image values with default image
	if imageSQL.Valid && imageSQL.String != "" {
		a.Image = imageSQL.String
	} else {
		a.Image = DefaultImageURL
	}

	return a, nil
}

// Create inserts a new accessory into the database and returns its ID.
func (r *AccessoryRepositoryImpl) Create(ctx context.Context, input models.NewAccessoryInput) (int, error) {
	status := determineStatus(input.Quantity)

	query := `
		INSERT INTO accessories (name, make, quantity, price, status, unit_color, image, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, NOW(), NOW())
	`

	stmt, err := r.DB.PrepareContext(ctx, query)
	if err != nil {
		return 0, err
	}
	defer stmt.Close()
	
	var imageValue interface{}
	if input.Image == "" || input.Image == DefaultImageURL {
		imageValue = nil // Use NULL if image is empty or default
	} else {
		imageValue = input.Image
	}

	res, err := stmt.ExecContext(
		ctx,
		input.Name,
		string(input.Make),
		input.Quantity,
		input.Price,
		string(status),
		string(input.UnitColor),
		imageValue,
	)

	if err != nil {
		return 0, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), nil
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
		// If image is empty string or "null", use the default image URL
		if *input.Image == "" || *input.Image == "null" {
			accessory.Image = DefaultImageURL
		} else {
		accessory.Image = *input.Image
		}
	}

	updateQuery := `
		UPDATE accessories
		SET name = ?, make = ?, quantity = ?, price = ?, status = ?, unit_color = ?, image = ?, updated_at = NOW()
		WHERE id = ?
	`

	stmt, err := r.DB.PrepareContext(ctx, updateQuery)
	if err != nil {
		return models.Accessory{}, err
	}
	defer stmt.Close()
	
	var imageValue interface{}
	if accessory.Image == "" || accessory.Image == DefaultImageURL {
		imageValue = nil // Use NULL if image is empty or default
	} else {
		imageValue = accessory.Image
	}

	_, err = stmt.ExecContext(
		ctx,
		accessory.Name,
		string(accessory.Make),
		accessory.Quantity,
		accessory.Price,
		string(accessory.Status),
		string(accessory.UnitColor),
		imageValue,
		id,
	)

	if err != nil {
		return models.Accessory{}, err
	}

	// Since RETURNING is not used, we might want to re-fetch the accessory to get the updated_at time set by NOW().
	// Or, we can assume the client will handle this or that accessory.UpdatedAt will be set by application logic if needed after this call.
	// For now, to keep it simple and avoid an extra DB call, we'll return the accessory object with its existing UpdatedAt value if it wasn't changed by NOW().
	// To get the DB-generated updated_at, we'd need another GetByID call here.
	// Let's call GetByID to ensure the returned accessory has the correct updated_at from the database.
	updatedAccessory, err := r.GetByID(ctx, id)
	if err != nil {
		// If fetching the updated accessory fails, we might return the partially updated one or a specific error.
		// For now, returning the error from GetByID.
		return models.Accessory{}, err
	}
	return updatedAccessory, nil
}

// Delete removes an accessory from the database
func (r *AccessoryRepositoryImpl) Delete(ctx context.Context, id int) error {
	query := `DELETE FROM accessories WHERE id = ?`

	stmt, err := r.DB.PrepareContext(ctx, query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	result, err := stmt.ExecContext(ctx, id)
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
