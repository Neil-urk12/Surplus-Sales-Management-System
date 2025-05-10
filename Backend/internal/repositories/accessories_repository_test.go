package repositories

import (
	"context"
	"database/sql"
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
	mock.ExpectPrepare(regexp.QuoteMeta(`
		SELECT id, name, make, quantity, price, status, unit_color, image, created_at, updated_at
		FROM accessories
		ORDER BY id ASC
	`)).ExpectQuery().WillReturnRows(rows)

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
	assert.Equal(t, models.MakeOEM, accessories[0].Make)
	assert.Equal(t, 10, accessories[0].Quantity)
	assert.Equal(t, 5000.0, accessories[0].Price)
	assert.Equal(t, models.StatusInStock, accessories[0].Status)
	assert.Equal(t, models.ColorBlack, accessories[0].UnitColor)

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

		mock.ExpectPrepare(regexp.QuoteMeta(`
			SELECT id, name, make, quantity, price, status, unit_color, image, created_at, updated_at
			FROM accessories
			WHERE id = ?
		`)).ExpectQuery().WithArgs(1).WillReturnRows(rows)

		repo := NewAccessoryRepository(db)
		accessory, err := repo.GetByID(context.Background(), 1)

		assert.NoError(t, err)
		assert.Equal(t, 1, accessory.ID)
		assert.Equal(t, "Steering Wheel", accessory.Name)
		assert.Equal(t, models.MakeOEM, accessory.Make)
		assert.Equal(t, 10, accessory.Quantity)
		assert.Equal(t, 5000.0, accessory.Price)
		assert.Equal(t, models.StatusInStock, accessory.Status)
		assert.Equal(t, models.ColorBlack, accessory.UnitColor)
	})

	// Test case 2: Accessory does not exist
	t.Run("Accessory does not exist", func(t *testing.T) {
		mock.ExpectPrepare(regexp.QuoteMeta(`
			SELECT id, name, make, quantity, price, status, unit_color, image, created_at, updated_at
			FROM accessories
			WHERE id = ?
		`)).ExpectQuery().WithArgs(99).WillReturnError(sql.ErrNoRows)

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
		Make:      models.MakeAftermarket,
		Quantity:  5,
		Price:     8500.0,
		UnitColor: models.ColorWhite,
		Image:     "image3.jpg",
	}

	// Setup expected query and result
	mock.ExpectPrepare(regexp.QuoteMeta(`
		INSERT INTO accessories (name, make, quantity, price, status, unit_color, image, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, NOW(), NOW())
	`)).ExpectExec().WithArgs(
		input.Name,
		string(input.Make),
		input.Quantity,
		input.Price,
		string(models.StatusInStock), // Status is calculated based on quantity
		string(input.UnitColor),
		input.Image,
	).WillReturnResult(sqlmock.NewResult(3, 1)) // ID 3, 1 row affected

	// Create repository with mock DB
	repo := NewAccessoryRepository(db)

	// Execute method
	id, err := repo.Create(context.Background(), input)

	// Assert no errors
	assert.NoError(t, err)

	// Assert result is the expected ID
	assert.Equal(t, 3, id)

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
	mock.ExpectPrepare(regexp.QuoteMeta(`
		SELECT id, name, make, quantity, price, status, unit_color, image, created_at, updated_at
		FROM accessories
		WHERE id = ?
	`)).ExpectQuery().WithArgs(id).WillReturnRows(
		sqlmock.NewRows([]string{"id", "name", "make", "quantity", "price", "status", "unit_color", "image", "created_at", "updated_at"}).
			AddRow(id, "Steering Wheel", "OEM", 10, 5000.0, "In Stock", "Black", "image1.jpg", now, now),
	)

	// Setup mock for update query
	mock.ExpectPrepare(regexp.QuoteMeta(`
		UPDATE accessories
		SET name = ?, make = ?, quantity = ?, price = ?, status = ?, unit_color = ?, image = ?, updated_at = NOW()
		WHERE id = ?
	`)).ExpectExec().WithArgs(
		name,                          // updated name
		string(models.MakeOEM),        // unchanged make
		quantity,                      // updated quantity
		price,                         // updated price
		string(models.StatusLowStock), // status updated based on quantity
		string(models.ColorBlack),     // unchanged color
		"image1.jpg",                  // unchanged image
		id,
	).WillReturnResult(sqlmock.NewResult(0, 1)) // No new ID, 1 row affected

	// Setup mock for GetByID again (to fetch the updated accessory)
	mock.ExpectPrepare(regexp.QuoteMeta(`
		SELECT id, name, make, quantity, price, status, unit_color, image, created_at, updated_at
		FROM accessories
		WHERE id = ?
	`)).ExpectQuery().WithArgs(id).WillReturnRows(
		sqlmock.NewRows([]string{"id", "name", "make", "quantity", "price", "status", "unit_color", "image", "created_at", "updated_at"}).
			AddRow(id, name, "OEM", quantity, price, "Low Stock", "Black", "image1.jpg", now, now),
	)

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

	// Test success case
	t.Run("Success", func(t *testing.T) {
		id := 1

		mock.ExpectPrepare(regexp.QuoteMeta(`
			DELETE FROM accessories
			WHERE id = ?
		`)).ExpectExec().WithArgs(id).WillReturnResult(sqlmock.NewResult(0, 1))

		repo := NewAccessoryRepository(db)
		err := repo.Delete(context.Background(), id)

		assert.NoError(t, err)
	})

	// Test not found case
	t.Run("Not Found", func(t *testing.T) {
		id := 999

		mock.ExpectPrepare(regexp.QuoteMeta(`
			DELETE FROM accessories
			WHERE id = ?
		`)).ExpectExec().WithArgs(id).WillReturnResult(sqlmock.NewResult(0, 0))

		repo := NewAccessoryRepository(db)
		err := repo.Delete(context.Background(), id)

		assert.Error(t, err)
		assert.Equal(t, "accessory not found", err.Error())
	})

	// Verify all expectations were met
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestDetermineStatus(t *testing.T) {
	testCases := []struct {
		name     string
		quantity int
		expected models.AccessoryStatus
	}{
		{
			name:     "Out of stock",
			quantity: 0,
			expected: models.StatusOutOfStock,
		},
		{
			name:     "Low stock - 1",
			quantity: 1,
			expected: models.StatusLowStock,
		},
		{
			name:     "Low stock - 2",
			quantity: 2,
			expected: models.StatusLowStock,
		},
		{
			name:     "In stock - 3",
			quantity: 3,
			expected: models.StatusInStock,
		},
		{
			name:     "In stock - 5",
			quantity: 5,
			expected: models.StatusInStock,
		},
		{
			name:     "Available - 6",
			quantity: 6,
			expected: models.StatusAvailable,
		},
		{
			name:     "Available - high quantity",
			quantity: 100,
			expected: models.StatusAvailable,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := determineStatus(tc.quantity)
			assert.Equal(t, tc.expected, result)
		})
	}
}
