package repository

import (
	"database/sql"
	"fmt"
	"oop/internal/models"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// NewMockDB creates a new mock database connection for testing.
func NewMockDB(t *testing.T) (*sql.DB, sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherRegexp)) // Use regexp matching
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	return db, mock
}

// --- Tests ---

func TestGetCabs_NoFilters(t *testing.T) {
	db, mock := NewMockDB(t)
	defer db.Close()
	repo := NewCabsRepository(db)

	now := time.Now()
	expectedCabs := []models.MultiCab{
		{ID: 1, Name: "RX‑7", Make: "Mazda", Quantity: 4, Price: 7000000, Status: "In Stock", UnitColor: "Blue", Image: "rx7.jpg", CreatedAt: now, UpdatedAt: now},
		{ID: 2, Name: "911 GT3", Make: "Porsche", Quantity: 2, Price: 15000000, Status: "In Stock", UnitColor: "White", Image: "911.jpg", CreatedAt: now, UpdatedAt: now},
	}

	rows := sqlmock.NewRows([]string{"id", "name", "make", "quantity", "price", "status", "unit_color", "image", "created_at", "updated_at"}).
		AddRow(expectedCabs[0].ID, expectedCabs[0].Name, expectedCabs[0].Make, expectedCabs[0].Quantity, expectedCabs[0].Price, expectedCabs[0].Status, expectedCabs[0].UnitColor, expectedCabs[0].Image, expectedCabs[0].CreatedAt, expectedCabs[0].UpdatedAt).
		AddRow(expectedCabs[1].ID, expectedCabs[1].Name, expectedCabs[1].Make, expectedCabs[1].Quantity, expectedCabs[1].Price, expectedCabs[1].Status, expectedCabs[1].UnitColor, expectedCabs[1].Image, expectedCabs[1].CreatedAt, expectedCabs[1].UpdatedAt)

	// Base query without filters
	query := "SELECT id, name, make, quantity, price, status, unit_color, image, created_at, updated_at FROM multicabs WHERE 1=1 ORDER BY created_at DESC"
	mock.ExpectQuery(regexp.QuoteMeta(query)).WillReturnRows(rows)

	cabs, err := repo.GetCabs(nil)
	require.NoError(t, err)
	assert.Equal(t, expectedCabs, cabs)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetCabs_WithFilters(t *testing.T) {
	db, mock := NewMockDB(t)
	defer db.Close()
	repo := NewCabsRepository(db)

	now := time.Now()
	// Sample data for filtering
	cabPorsche911 := models.MultiCab{ID: 2, Name: "911 GT3", Make: "Porsche", Quantity: 2, Price: 15000000, Status: "In Stock", UnitColor: "White", Image: "911.jpg", CreatedAt: now, UpdatedAt: now}
	cabPorscheCayenne := models.MultiCab{ID: 7, Name: "Cayenne", Make: "Porsche", Quantity: 1, Price: 12000000, Status: "Sold Out", UnitColor: "Silver", Image: "cayenne.jpg", CreatedAt: now, UpdatedAt: now}
	cabMustang := models.MultiCab{ID: 3, Name: "Mustang", Make: "Ford", Quantity: 5, Price: 5500000, Status: "Available", UnitColor: "Red", Image: "mustang.jpg", CreatedAt: now, UpdatedAt: now}
	cabRX7 := models.MultiCab{ID: 1, Name: "RX‑7", Make: "Mazda", Quantity: 4, Price: 7000000, Status: "In Stock", UnitColor: "Blue", Image: "rx7.jpg", CreatedAt: now, UpdatedAt: now}

	cols := []string{"id", "name", "make", "quantity", "price", "status", "unit_color", "image", "created_at", "updated_at"}

	// Filter by Make
	t.Run("Filter by Make", func(t *testing.T) {
		rowsMake := sqlmock.NewRows(cols).
			AddRow(cabPorsche911.ID, cabPorsche911.Name, cabPorsche911.Make, cabPorsche911.Quantity, cabPorsche911.Price, cabPorsche911.Status, cabPorsche911.UnitColor, cabPorsche911.Image, cabPorsche911.CreatedAt, cabPorsche911.UpdatedAt).
			AddRow(cabPorscheCayenne.ID, cabPorscheCayenne.Name, cabPorscheCayenne.Make, cabPorscheCayenne.Quantity, cabPorscheCayenne.Price, cabPorscheCayenne.Status, cabPorscheCayenne.UnitColor, cabPorscheCayenne.Image, cabPorscheCayenne.CreatedAt, cabPorscheCayenne.UpdatedAt)

		queryMake := "SELECT id, name, make, quantity, price, status, unit_color, image, created_at, updated_at FROM multicabs WHERE 1=1 AND make = \\? ORDER BY created_at DESC"
		mock.ExpectQuery(queryMake).WithArgs("Porsche").WillReturnRows(rowsMake)

		filtersMake := map[string]interface{}{"make": "Porsche"}
		cabsMake, errMake := repo.GetCabs(filtersMake)
		require.NoError(t, errMake)
		assert.Len(t, cabsMake, 2)
		assert.Equal(t, []models.MultiCab{cabPorsche911, cabPorscheCayenne}, cabsMake)
	})

	// Filter by Status
	t.Run("Filter by Status", func(t *testing.T) {
		rowsStatus := sqlmock.NewRows(cols).
			AddRow(cabMustang.ID, cabMustang.Name, cabMustang.Make, cabMustang.Quantity, cabMustang.Price, cabMustang.Status, cabMustang.UnitColor, cabMustang.Image, cabMustang.CreatedAt, cabMustang.UpdatedAt)

		queryStatus := "SELECT id, name, make, quantity, price, status, unit_color, image, created_at, updated_at FROM multicabs WHERE 1=1 AND status = \\? ORDER BY created_at DESC"
		mock.ExpectQuery(queryStatus).WithArgs("Available").WillReturnRows(rowsStatus)

		filtersStatus := map[string]interface{}{"status": "Available"}
		cabsStatus, errStatus := repo.GetCabs(filtersStatus)
		require.NoError(t, errStatus)
		assert.Len(t, cabsStatus, 1)
		assert.Equal(t, []models.MultiCab{cabMustang}, cabsStatus)
	})

	// Filter by Search (Name)
	t.Run("Filter by Search Name", func(t *testing.T) {
		rowsSearchName := sqlmock.NewRows(cols).
			AddRow(cabRX7.ID, cabRX7.Name, cabRX7.Make, cabRX7.Quantity, cabRX7.Price, cabRX7.Status, cabRX7.UnitColor, cabRX7.Image, cabRX7.CreatedAt, cabRX7.UpdatedAt)

		// Note: The actual query uses LIKE which needs % in args
		// Also need to escape parenthesis in regex for the OR condition
		querySearchName := "SELECT id, name, make, quantity, price, status, unit_color, image, created_at, updated_at FROM multicabs WHERE 1=1 AND \\(name LIKE \\? OR make LIKE \\?\\) ORDER BY created_at DESC"
		searchTerm := "%RX%"
		mock.ExpectQuery(querySearchName).WithArgs(searchTerm, searchTerm).WillReturnRows(rowsSearchName)

		filtersSearchName := map[string]interface{}{"search": "RX"}
		cabsSearchName, errSearchName := repo.GetCabs(filtersSearchName)
		require.NoError(t, errSearchName)
		assert.Len(t, cabsSearchName, 1)
		assert.Equal(t, []models.MultiCab{cabRX7}, cabsSearchName)
	})

	// Filter by Search (Make)
	t.Run("Filter by Search Make", func(t *testing.T) {
		rowsSearchMake := sqlmock.NewRows(cols).
			AddRow(cabMustang.ID, cabMustang.Name, cabMustang.Make, cabMustang.Quantity, cabMustang.Price, cabMustang.Status, cabMustang.UnitColor, cabMustang.Image, cabMustang.CreatedAt, cabMustang.UpdatedAt)

		// Note: The actual query uses LIKE which needs % in args
		// Also need to escape parenthesis in regex for the OR condition
		querySearchMake := "SELECT id, name, make, quantity, price, status, unit_color, image, created_at, updated_at FROM multicabs WHERE 1=1 AND \\(name LIKE \\? OR make LIKE \\?\\) ORDER BY created_at DESC"
		searchTerm := "%ford%"
		mock.ExpectQuery(querySearchMake).WithArgs(searchTerm, searchTerm).WillReturnRows(rowsSearchMake)

		filtersSearchMake := map[string]interface{}{"search": "ford"}
		cabsSearchMake, errSearchMake := repo.GetCabs(filtersSearchMake)
		require.NoError(t, errSearchMake)
		assert.Len(t, cabsSearchMake, 1)
		assert.Equal(t, []models.MultiCab{cabMustang}, cabsSearchMake)
	})

	// Combined Filters
	t.Run("Combined Filters", func(t *testing.T) {
		rowsCombined := sqlmock.NewRows(cols).
			AddRow(cabPorsche911.ID, cabPorsche911.Name, cabPorsche911.Make, cabPorsche911.Quantity, cabPorsche911.Price, cabPorsche911.Status, cabPorsche911.UnitColor, cabPorsche911.Image, cabPorsche911.CreatedAt, cabPorsche911.UpdatedAt)

		queryCombined := "SELECT id, name, make, quantity, price, status, unit_color, image, created_at, updated_at FROM multicabs WHERE 1=1 AND make = \\? AND status = \\? ORDER BY created_at DESC"
		mock.ExpectQuery(queryCombined).WithArgs("Porsche", "In Stock").WillReturnRows(rowsCombined)

		filtersCombined := map[string]interface{}{"make": "Porsche", "status": "In Stock"}
		cabsCombined, errCombined := repo.GetCabs(filtersCombined)
		require.NoError(t, errCombined)
		assert.Len(t, cabsCombined, 1)
		assert.Equal(t, []models.MultiCab{cabPorsche911}, cabsCombined)
	})

	// No results
	t.Run("No Results", func(t *testing.T) {
		rowsNone := sqlmock.NewRows(cols) // No rows added

		queryNone := "SELECT id, name, make, quantity, price, status, unit_color, image, created_at, updated_at FROM multicabs WHERE 1=1 AND make = \\? ORDER BY created_at DESC"
		mock.ExpectQuery(queryNone).WithArgs("Ferrari").WillReturnRows(rowsNone)

		filtersNone := map[string]interface{}{"make": "Ferrari"}
		cabsNone, errNone := repo.GetCabs(filtersNone)
		require.NoError(t, errNone)
		assert.Len(t, cabsNone, 0)
	})

	// Query Error
	t.Run("Query Error", func(t *testing.T) {
		queryErr := "SELECT id, name, make, quantity, price, status, unit_color, image, created_at, updated_at FROM multicabs WHERE 1=1 AND make = \\? ORDER BY created_at DESC"
		mock.ExpectQuery(queryErr).WithArgs("ErrorCase").WillReturnError(sql.ErrConnDone)

		filtersErr := map[string]interface{}{"make": "ErrorCase"}
		cabsErr, err := repo.GetCabs(filtersErr)
		assert.ErrorIs(t, err, sql.ErrConnDone)
		assert.Nil(t, cabsErr)
	})

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetCabByID_Exists(t *testing.T) {
	db, mock := NewMockDB(t)
	defer db.Close()
	repo := NewCabsRepository(db)

	now := time.Now()
	expectedCab := &models.MultiCab{ID: 1, Name: "RX‑7", Make: "Mazda", Quantity: 4, Price: 7000000, Status: "In Stock", UnitColor: "Blue", Image: "rx7.jpg", CreatedAt: now, UpdatedAt: now}

	rows := sqlmock.NewRows([]string{"id", "name", "make", "quantity", "price", "status", "unit_color", "image", "created_at", "updated_at"}).
		AddRow(expectedCab.ID, expectedCab.Name, expectedCab.Make, expectedCab.Quantity, expectedCab.Price, expectedCab.Status, expectedCab.UnitColor, expectedCab.Image, expectedCab.CreatedAt, expectedCab.UpdatedAt)

	query := "SELECT id, name, make, quantity, price, status, unit_color, image, created_at, updated_at FROM multicabs WHERE id = \\?"
	mock.ExpectQuery(query).WithArgs(expectedCab.ID).WillReturnRows(rows)

	id := 1
	cab, err := repo.GetCabByID(id)
	require.NoError(t, err)
	require.NotNil(t, cab)
	assert.Equal(t, expectedCab, cab)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetCabByID_NotExists(t *testing.T) {
	db, mock := NewMockDB(t)
	defer db.Close()
	repo := NewCabsRepository(db)

	query := "SELECT id, name, make, quantity, price, status, unit_color, image, created_at, updated_at FROM multicabs WHERE id = \\?"
	IDNotFound := 99
	mock.ExpectQuery(query).WithArgs(IDNotFound).WillReturnError(sql.ErrNoRows)

	id := 99
	cab, err := repo.GetCabByID(id)
	require.Error(t, err, "Should return error for non-existent ID")
	assert.Nil(t, cab, "Cab should be nil for non-existent ID")
	// Check error message specific to the DB implementation
	assert.Contains(t, err.Error(), fmt.Sprintf("cab with ID %d not found", id), "Error message should indicate 'not found'")
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestAddCab_Success(t *testing.T) {
	db, mock := NewMockDB(t)
	defer db.Close()
	repo := NewCabsRepository(db)

	newCabData := models.MultiCab{
		Name:      "Test Car",
		Make:      "Test Make",
		Quantity:  10,
		Price:     500000,
		Status:    "Available",
		UnitColor: "Green",
		Image:     "test.jpg",
	}

	insertQuery := "INSERT INTO multicabs \\(name, make, quantity, price, status, unit_color, image, created_at, updated_at\\) VALUES \\(\\?, \\?, \\?, \\?, \\?, \\?, \\?, \\?, \\?\\)"
	mock.ExpectExec(insertQuery).
		WithArgs(newCabData.Name, newCabData.Make, newCabData.Quantity, newCabData.Price, newCabData.Status, newCabData.UnitColor, newCabData.Image, sqlmock.AnyArg(), sqlmock.AnyArg()). // Use AnyArg for timestamps
		WillReturnResult(sqlmock.NewResult(8, 1))                                                                                                                                         // Expecting ID 8, 1 row affected

	addedCab, err := repo.AddCab(newCabData)
	require.NoError(t, err)
	require.NotNil(t, addedCab)

	// ID is generated by the database, check if it's non-zero (or > 0 depending on strategy)
	assert.Equal(t, 8, addedCab.ID) // Check against the LastInsertId returned by mock
	// Assert other fields match
	assert.Equal(t, newCabData.Name, addedCab.Name)
	assert.Equal(t, newCabData.Make, addedCab.Make)
	assert.Equal(t, newCabData.Quantity, addedCab.Quantity)
	assert.Equal(t, newCabData.Price, addedCab.Price)
	assert.Equal(t, newCabData.Status, addedCab.Status)
	assert.Equal(t, newCabData.UnitColor, addedCab.UnitColor)
	assert.Equal(t, newCabData.Image, addedCab.Image)
	// Check timestamps are recent (within a reasonable delta)
	assert.WithinDuration(t, time.Now(), addedCab.CreatedAt, 5*time.Second, "CreatedAt should be recent")
	assert.WithinDuration(t, time.Now(), addedCab.UpdatedAt, 5*time.Second, "UpdatedAt should be recent")

	assert.NoError(t, mock.ExpectationsWereMet())

	// Note: Verification by retrieval would require a separate test or mock setup within this test.
}

func TestAddCab_ValidationError(t *testing.T) {
	db, mock := NewMockDB(t)
	defer db.Close()
	repo := NewCabsRepository(db)

	// Test case for validation defined in the DB repository's AddCab method
	invalidCabData := models.MultiCab{
		Name:      "", // Empty name should trigger validation
		Make:      "Test Make",
		UnitColor: "Red", // Assuming UnitColor is required by DB repo validation
	}
	addedCab, err := repo.AddCab(invalidCabData)
	require.Error(t, err, "Should return error for invalid data")
	assert.Nil(t, addedCab)
	// Check for the specific validation error message from the repository
	assert.Contains(t, err.Error(), "cab name, make, and color cannot be empty")

	// Verify no DB expectations were made or missed
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUpdateCab_Success(t *testing.T) {
	db, mock := NewMockDB(t)
	defer db.Close()
	repo := NewCabsRepository(db)

	idToUpdate := 1
	now := time.Now()
	// We need to mock the initial GetCabByID call within UpdateCab
	originalCab := &models.MultiCab{ID: idToUpdate, Name: "RX‑7", Make: "Mazda", Quantity: 4, Price: 7000000, Status: "In Stock", UnitColor: "Blue", Image: "rx7.jpg", CreatedAt: now.Add(-time.Hour), UpdatedAt: now.Add(-time.Hour)}
	getByIDQuery := "SELECT id, name, make, quantity, price, status, unit_color, image, created_at, updated_at FROM multicabs WHERE id = \\?"
	rowsGet := sqlmock.NewRows([]string{"id", "name", "make", "quantity", "price", "status", "unit_color", "image", "created_at", "updated_at"}).
		AddRow(originalCab.ID, originalCab.Name, originalCab.Make, originalCab.Quantity, originalCab.Price, originalCab.Status, originalCab.UnitColor, originalCab.Image, originalCab.CreatedAt, originalCab.UpdatedAt)
	mock.ExpectQuery(getByIDQuery).WithArgs(idToUpdate).WillReturnRows(rowsGet)

	updateData := models.MultiCab{
		Name:      "RX-7 Updated",
		Make:      "Mazda", // Keep Make same
		Quantity:  originalCab.Quantity + 5,
		Price:     originalCab.Price + 100000,
		Status:    "Low Stock",
		UnitColor: "Red", // Change color
		Image:     "updated.jpg",
		// ID and CreatedAt should be ignored by UpdateCab logic
	}

	// Mock the UPDATE execution
	updateQuery := "UPDATE multicabs SET name = \\?, make = \\?, quantity = \\?, price = \\?, status = \\?, unit_color = \\?, image = \\?, updated_at = \\? WHERE id = \\?"
	mock.ExpectExec(updateQuery).
		WithArgs(updateData.Name, updateData.Make, updateData.Quantity, updateData.Price, updateData.Status, updateData.UnitColor, updateData.Image, sqlmock.AnyArg(), idToUpdate).
		WillReturnResult(sqlmock.NewResult(0, 1)) // 1 row affected

	updatedCab, err := repo.UpdateCab(idToUpdate, updateData)
	require.NoError(t, err)
	require.NotNil(t, updatedCab)

	// Assertions for updated fields
	assert.Equal(t, idToUpdate, updatedCab.ID)
	assert.Equal(t, updateData.Name, updatedCab.Name)
	assert.Equal(t, updateData.Make, updatedCab.Make) // Should remain "Mazda"
	assert.Equal(t, updateData.Quantity, updatedCab.Quantity)
	assert.Equal(t, updateData.Price, updatedCab.Price)
	assert.Equal(t, updateData.Status, updatedCab.Status)
	assert.Equal(t, updateData.UnitColor, updatedCab.UnitColor)
	assert.Equal(t, updateData.Image, updatedCab.Image)
	// CreatedAt should not change
	assert.Equal(t, originalCab.CreatedAt, updatedCab.CreatedAt, "CreatedAt should match original")
	// UpdatedAt should be newer than the original
	assert.True(t, updatedCab.UpdatedAt.After(originalCab.UpdatedAt), "UpdatedAt should be newer")
	// Check UpdatedAt is recent
	assert.WithinDuration(t, time.Now(), updatedCab.UpdatedAt, 5*time.Second, "UpdatedAt should be recent")

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUpdateCab_NotExists(t *testing.T) {
	db, mock := NewMockDB(t)
	defer db.Close()
	repo := NewCabsRepository(db)

	id := 99
	// Mock the GetCabByID call which should return not found
	getByIDQuery := "SELECT id, name, make, quantity, price, status, unit_color, image, created_at, updated_at FROM multicabs WHERE id = \\?"
	mock.ExpectQuery(getByIDQuery).WithArgs(id).WillReturnError(sql.ErrNoRows)

	updateData := models.MultiCab{Name: "Does not matter"}
	updatedCab, err := repo.UpdateCab(id, updateData)
	require.Error(t, err, "Should return error for non-existent ID")
	assert.Nil(t, updatedCab)
	// Check error message specific to the DB implementation's update check
	assert.Contains(t, err.Error(), fmt.Sprintf("cab with ID %d not found for update", id))
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestDeleteCab_Success(t *testing.T) {
	db, mock := NewMockDB(t)
	defer db.Close()
	repo := NewCabsRepository(db)

	idToDelete := 1

	// Mock the GetCabByID check before deletion
	getByIDQuery := "SELECT id, name, make, quantity, price, status, unit_color, image, created_at, updated_at FROM multicabs WHERE id = \\\\?"
	// Return data for all columns expected by GetCabByID's Scan
	cols := []string{"id", "name", "make", "quantity", "price", "status", "unit_color", "image", "created_at", "updated_at"}
	rowsGet := sqlmock.NewRows(cols).
		AddRow(idToDelete, "Dummy Name", "Dummy Make", 0, 0.0, "Dummy Status", "Dummy Color", "dummy.jpg", time.Now(), time.Now()) // Provide dummy values
	mock.ExpectQuery(getByIDQuery).WithArgs(idToDelete).WillReturnRows(rowsGet)

	// Mock the DELETE execution
	deleteQuery := "DELETE FROM multicabs WHERE id = \\\\?"
	mock.ExpectExec(deleteQuery).WithArgs(idToDelete).WillReturnResult(sqlmock.NewResult(0, 1)) // 1 row affected

	// Perform deletion
	errDelete := repo.DeleteCab(idToDelete)
	require.NoError(t, errDelete, "Should successfully delete the cab")

	assert.NoError(t, mock.ExpectationsWereMet())

	// Verification of deletion would require separate mock expectations if needed.
}

func TestDeleteCab_NotExists(t *testing.T) {
	db, mock := NewMockDB(t)
	defer db.Close()
	repo := NewCabsRepository(db)

	id := 99
	// Mock the GetCabByID check which should return not found
	getByIDQuery := "SELECT id, name, make, quantity, price, status, unit_color, image, created_at, updated_at FROM multicabs WHERE id = \\?"
	mock.ExpectQuery(getByIDQuery).WithArgs(id).WillReturnError(sql.ErrNoRows)

	// DELETE query should not be executed if GetByID fails

	err := repo.DeleteCab(id)
	require.Error(t, err, "Should return error for non-existent ID")
	// Check error message specific to the DB implementation's delete check
	assert.Contains(t, err.Error(), fmt.Sprintf("cab with ID %d not found for deletion", id))

	assert.NoError(t, mock.ExpectationsWereMet())
}
