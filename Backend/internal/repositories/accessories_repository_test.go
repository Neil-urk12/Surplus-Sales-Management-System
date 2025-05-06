package repositories

import (
	"context"
	"database/sql"
	"errors"
	"oop/internal/models"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestAccessoryRepository_GetAll(t *testing.T) {
	// Create a new SQL mock
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create mock: %v", err)
	}
	defer db.Close()

	// Create columns for the mock result
	columns := []string{"id", "name", "make", "quantity", "price", "status", "unit_color", "image", "created_at", "updated_at"}

	// Create expected time values
	now := time.Now()

	// Create expected rows
	rows := sqlmock.NewRows(columns).
		AddRow(1, "Steering Wheel", "OEM", 10, 5000.0, "In Stock", "Black", "image1.jpg", now, now).
		AddRow(2, "Sport Seats", "Aftermarket", 0, 12000.0, "Out of Stock", "Silver", "image2.jpg", now, now)

	// Set up expected query and result
	mock.ExpectQuery(regexp.QuoteMeta(`
		SELECT id, name, make, quantity, price, status, unit_color, image, created_at, updated_at
		FROM accessories
		ORDER BY id ASC
	`)).WillReturnRows(rows)

	// Create repository with mock DB
	repo := NewAccessoryRepository(db)

	// Execute the method
	accessories, err := repo.GetAll(context.Background())

	// Assert no errors occurred
	assert.NoError(t, err)

	// Assert we got the expected results
	assert.Len(t, accessories, 2)
	assert.Equal(t, 1, accessories[0].ID)
	assert.Equal(t, "Steering Wheel", accessories[0].Name)
	assert.Equal(t, models.AccessoryMake("OEM"), accessories[0].Make)
	assert.Equal(t, 10, accessories[0].Quantity)
	assert.Equal(t, 5000.0, accessories[0].Price)
	assert.Equal(t, models.AccessoryStatus("In Stock"), accessories[0].Status)
	assert.Equal(t, models.AccessoryColor("Black"), accessories[0].UnitColor)

	// Verify all expectations were met
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestAccessoryRepository_GetByID(t *testing.T) {
	// Create a new SQL mock
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create mock: %v", err)
	}
	defer db.Close()

	// Create columns for the mock result
	columns := []string{"id", "name", "make", "quantity", "price", "status", "unit_color", "image", "created_at", "updated_at"}

	// Create expected time values
	now := time.Now()

	// Test case 1: Accessory exists
	t.Run("Accessory exists", func(t *testing.T) {
		rows := sqlmock.NewRows(columns).
			AddRow(1, "Steering Wheel", "OEM", 10, 5000.0, "In Stock", "Black", "image1.jpg", now, now)

		mock.ExpectQuery(regexp.QuoteMeta(`
			SELECT id, name, make, quantity, price, status, unit_color, image, created_at, updated_at
			FROM accessories
			WHERE id = $1
		`)).WithArgs(1).WillReturnRows(rows)

		repo := NewAccessoryRepository(db)
		accessory, err := repo.GetByID(context.Background(), 1)

		assert.NoError(t, err)
		assert.Equal(t, 1, accessory.ID)
		assert.Equal(t, "Steering Wheel", accessory.Name)
		assert.Equal(t, models.AccessoryMake("OEM"), accessory.Make)
		assert.Equal(t, 10, accessory.Quantity)
		assert.Equal(t, 5000.0, accessory.Price)
		assert.Equal(t, models.AccessoryStatus("In Stock"), accessory.Status)
		assert.Equal(t, models.AccessoryColor("Black"), accessory.UnitColor)
	})

	// Test case 2: Accessory does not exist
	t.Run("Accessory does not exist", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta(`
			SELECT id, name, make, quantity, price, status, unit_color, image, created_at, updated_at
			FROM accessories
			WHERE id = $1
		`)).WithArgs(99).WillReturnError(sql.ErrNoRows)

		repo := NewAccessoryRepository(db)
		_, err := repo.GetByID(context.Background(), 99)

		assert.Error(t, err)
		assert.Equal(t, "accessory not found", err.Error())
	})

	// Verify all expectations were met
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestAccessoryRepository_Create(t *testing.T) {
	// Create a new SQL mock
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create mock: %v", err)
	}
	defer db.Close()

	// Create test input
	input := models.NewAccessoryInput{
		Name:      "LED Headlights",
		Make:      models.AccessoryMake("Aftermarket"),
		Quantity:  5,
		Price:     8500.0,
		UnitColor: models.AccessoryColor("White"),
		Image:     "image3.jpg",
	}

	// Create expected time values
	now := time.Now()

	// Setup expected query and result
	mock.ExpectQuery(regexp.QuoteMeta(`
		INSERT INTO accessories (name, make, quantity, price, status, unit_color, image, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, NOW(), NOW())
		RETURNING id, created_at, updated_at
	`)).WithArgs(
		input.Name,
		string(input.Make),
		input.Quantity,
		input.Price,
		string(models.StatusInStock), // Status is calculated based on quantity
		string(input.UnitColor),
		input.Image,
	).WillReturnRows(sqlmock.NewRows([]string{"id", "created_at", "updated_at"}).AddRow(3, now, now))

	// Create repository with mock DB
	repo := NewAccessoryRepository(db)

	// Execute method
	result, err := repo.Create(context.Background(), input)

	// Assert no errors
	assert.NoError(t, err)

	// Assert result
	assert.Equal(t, 3, result.ID)
	assert.Equal(t, input.Name, result.Name)
	assert.Equal(t, input.Make, result.Make)
	assert.Equal(t, input.Quantity, result.Quantity)
	assert.Equal(t, input.Price, result.Price)
	assert.Equal(t, models.StatusInStock, result.Status)
	assert.Equal(t, input.UnitColor, result.UnitColor)
	assert.Equal(t, input.Image, result.Image)

	// Verify all expectations were met
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestAccessoryRepository_Update(t *testing.T) {
	// Create a new SQL mock
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create mock: %v", err)
	}
	defer db.Close()

	// Create test data
	id := 1
	name := "Updated Steering Wheel"
	quantity := 2
	price := 5500.0

	// Create expected time values
	now := time.Now()

	// Create update input
	input := models.UpdateAccessoryInput{
		Name:     &name,
		Quantity: &quantity,
		Price:    &price,
	}

	// Setup mock for GetByID (first step in the update process)
	mock.ExpectQuery(regexp.QuoteMeta(`
		SELECT id, name, make, quantity, price, status, unit_color, image, created_at, updated_at
		FROM accessories
		WHERE id = $1
	`)).WithArgs(id).WillReturnRows(
		sqlmock.NewRows([]string{"id", "name", "make", "quantity", "price", "status", "unit_color", "image", "created_at", "updated_at"}).
			AddRow(id, "Steering Wheel", "OEM", 10, 5000.0, "In Stock", "Black", "image1.jpg", now, now),
	)

	// Setup mock for update query
	mock.ExpectQuery(regexp.QuoteMeta(`
		UPDATE accessories
		SET name = $1, make = $2, quantity = $3, price = $4, status = $5, unit_color = $6, image = $7, updated_at = NOW()
		WHERE id = $8
		RETURNING updated_at
	`)).WithArgs(
		name,                          // updated name
		string(models.MakeOEM),        // unchanged make
		quantity,                      // updated quantity
		price,                         // updated price
		string(models.StatusLowStock), // status updated based on quantity
		string(models.ColorBlack),     // unchanged color
		"image1.jpg",                  // unchanged image
		id,
	).WillReturnRows(sqlmock.NewRows([]string{"updated_at"}).AddRow(now))

	// Create repository with mock DB
	repo := NewAccessoryRepository(db)

	// Execute method
	result, err := repo.Update(context.Background(), id, input)

	// Assert no errors
	assert.NoError(t, err)

	// Assert result has expected updates
	assert.Equal(t, id, result.ID)
	assert.Equal(t, name, result.Name)
	assert.Equal(t, models.MakeOEM, result.Make) // Unchanged
	assert.Equal(t, quantity, result.Quantity)
	assert.Equal(t, price, result.Price)
	assert.Equal(t, models.StatusLowStock, result.Status) // Changed based on quantity
	assert.Equal(t, models.ColorBlack, result.UnitColor)  // Unchanged

	// Verify all expectations were met
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestAccessoryRepository_Delete(t *testing.T) {
	// Create a new SQL mock
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create mock: %v", err)
	}
	defer db.Close()

	// Test case 1: Successful deletion
	t.Run("Successful deletion", func(t *testing.T) {
		// Setup mock for delete query
		mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM accessories WHERE id = $1`)).
			WithArgs(1).
			WillReturnResult(sqlmock.NewResult(0, 1)) // 1 row affected

		// Create repository with mock DB
		repo := NewAccessoryRepository(db)

		// Execute method
		err := repo.Delete(context.Background(), 1)

		// Assert no errors
		assert.NoError(t, err)
	})

	// Test case 2: Accessory not found
	t.Run("Accessory not found", func(t *testing.T) {
		// Setup mock for delete query
		mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM accessories WHERE id = $1`)).
			WithArgs(99).
			WillReturnResult(sqlmock.NewResult(0, 0)) // 0 rows affected

		// Create repository with mock DB
		repo := NewAccessoryRepository(db)

		// Execute method
		err := repo.Delete(context.Background(), 99)

		// Assert error is returned
		assert.Error(t, err)
		assert.Equal(t, "accessory not found", err.Error())
	})

	// Test case 3: Database error
	t.Run("Database error", func(t *testing.T) {
		// Setup mock for delete query with error
		dbError := errors.New("database connection lost")
		mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM accessories WHERE id = $1`)).
			WithArgs(1).
			WillReturnError(dbError)

		// Create repository with mock DB
		repo := NewAccessoryRepository(db)

		// Execute method
		err := repo.Delete(context.Background(), 1)

		// Assert error is returned
		assert.Error(t, err)
		assert.Equal(t, dbError, err)
	})

	// Verify all expectations were met
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestDetermineStatus(t *testing.T) {
	// Test different quantity thresholds
	testCases := []struct {
		name     string
		quantity int
		expected models.AccessoryStatus
	}{
		{"Zero quantity", 0, models.StatusOutOfStock},
		{"Low stock threshold 1", 1, models.StatusLowStock},
		{"Low stock threshold 2", 2, models.StatusLowStock},
		{"In stock threshold 3", 3, models.StatusInStock},
		{"In stock threshold 5", 5, models.StatusInStock},
		{"Available threshold 6", 6, models.StatusAvailable},
		{"High quantity", 100, models.StatusAvailable},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			status := determineStatus(tc.quantity)
			assert.Equal(t, tc.expected, status)
		})
	}
}
