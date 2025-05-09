package repositories

import (
	"database/sql"
	"regexp"
	"testing"
	"time"

	"oop/internal/models"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func NewMockDB(t *testing.T) (*sql.DB, sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherRegexp)) // Use regexp matching
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	return db, mock
}

func TestGetAllMaterials(t *testing.T) {
	db, mock := NewMockDB(t)
	defer db.Close()
	repo := NewMaterialRepository(db)

	now := time.Now()
	expectedMaterials := []models.Material{
		{ID: 1, Name: "Material 1", Category: "Cat A", Supplier: "Sup 1", Quantity: 10, Status: "Active", Image: "img1.jpg", CreatedAt: now, UpdatedAt: now},
		{ID: 2, Name: "Material 2", Category: "Cat B", Supplier: "Sup 2", Quantity: 20, Status: "Inactive", Image: "img2.jpg", CreatedAt: now, UpdatedAt: now},
	}

	t.Run("No Filters", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"id", "name", "category", "supplier", "quantity", "status", "image", "created_at", "updated_at"}).
			AddRow(expectedMaterials[0].ID, expectedMaterials[0].Name, expectedMaterials[0].Category, expectedMaterials[0].Supplier, expectedMaterials[0].Quantity, expectedMaterials[0].Status, expectedMaterials[0].Image, expectedMaterials[0].CreatedAt, expectedMaterials[0].UpdatedAt).
			AddRow(expectedMaterials[1].ID, expectedMaterials[1].Name, expectedMaterials[1].Category, expectedMaterials[1].Supplier, expectedMaterials[1].Quantity, expectedMaterials[1].Status, expectedMaterials[1].Image, expectedMaterials[1].CreatedAt, expectedMaterials[1].UpdatedAt)

		query := "SELECT id, name, category, supplier, quantity, status, image, created_at, updated_at FROM materials WHERE 1=1 ORDER BY created_at DESC"
		mock.ExpectQuery(regexp.QuoteMeta(query)).WillReturnRows(rows)

		materials, err := repo.GetAll("", "", "", "")
		assert.NoError(t, err)
		assert.Equal(t, expectedMaterials, materials)
	})

	t.Run("With Search Term", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"id", "name", "category", "supplier", "quantity", "status", "image", "created_at", "updated_at"}).
			AddRow(expectedMaterials[0].ID, expectedMaterials[0].Name, expectedMaterials[0].Category, expectedMaterials[0].Supplier, expectedMaterials[0].Quantity, expectedMaterials[0].Status, expectedMaterials[0].Image, expectedMaterials[0].CreatedAt, expectedMaterials[0].UpdatedAt)

		querySearch := "SELECT id, name, category, supplier, quantity, status, image, created_at, updated_at FROM materials WHERE 1=1 AND (name LIKE ? OR category LIKE ? OR supplier LIKE ? OR CAST(id AS CHAR) LIKE ?) ORDER BY created_at DESC"
		mock.ExpectQuery(regexp.QuoteMeta(querySearch)).WithArgs("%term%", "%term%", "%term%", "%term%").WillReturnRows(rows)

		materials, err := repo.GetAll("term", "", "", "")
		assert.NoError(t, err)
		assert.Equal(t, []models.Material{expectedMaterials[0]}, materials)
	})

	t.Run("With Category", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"id", "name", "category", "supplier", "quantity", "status", "image", "created_at", "updated_at"}).
			AddRow(expectedMaterials[0].ID, expectedMaterials[0].Name, expectedMaterials[0].Category, expectedMaterials[0].Supplier, expectedMaterials[0].Quantity, expectedMaterials[0].Status, expectedMaterials[0].Image, expectedMaterials[0].CreatedAt, expectedMaterials[0].UpdatedAt)

		queryCategory := "SELECT id, name, category, supplier, quantity, status, image, created_at, updated_at FROM materials WHERE 1=1 AND category = ? ORDER BY created_at DESC"
		mock.ExpectQuery(regexp.QuoteMeta(queryCategory)).WithArgs("Cat A").WillReturnRows(rows)

		materials, err := repo.GetAll("", "Cat A", "", "")
		assert.NoError(t, err)
		assert.Equal(t, []models.Material{expectedMaterials[0]}, materials)
	})

	t.Run("With Supplier", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"id", "name", "category", "supplier", "quantity", "status", "image", "created_at", "updated_at"}).
			AddRow(expectedMaterials[0].ID, expectedMaterials[0].Name, expectedMaterials[0].Category, expectedMaterials[0].Supplier, expectedMaterials[0].Quantity, expectedMaterials[0].Status, expectedMaterials[0].Image, expectedMaterials[0].CreatedAt, expectedMaterials[0].UpdatedAt)

		querySupplier := "SELECT id, name, category, supplier, quantity, status, image, created_at, updated_at FROM materials WHERE 1=1 AND supplier = ? ORDER BY created_at DESC"
		mock.ExpectQuery(regexp.QuoteMeta(querySupplier)).WithArgs("Sup 1").WillReturnRows(rows)

		materials, err := repo.GetAll("", "", "Sup 1", "")
		assert.NoError(t, err)
		assert.Equal(t, []models.Material{expectedMaterials[0]}, materials)
	})

	t.Run("With Status", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"id", "name", "category", "supplier", "quantity", "status", "image", "created_at", "updated_at"}).
			AddRow(expectedMaterials[0].ID, expectedMaterials[0].Name, expectedMaterials[0].Category, expectedMaterials[0].Supplier, expectedMaterials[0].Quantity, expectedMaterials[0].Status, expectedMaterials[0].Image, expectedMaterials[0].CreatedAt, expectedMaterials[0].UpdatedAt)

		queryStatus := "SELECT id, name, category, supplier, quantity, status, image, created_at, updated_at FROM materials WHERE 1=1 AND status = ? ORDER BY created_at DESC"
		mock.ExpectQuery(regexp.QuoteMeta(queryStatus)).WithArgs("Active").WillReturnRows(rows)

		materials, err := repo.GetAll("", "", "", "Active")
		assert.NoError(t, err)
		assert.Equal(t, []models.Material{expectedMaterials[0]}, materials)
	})

	t.Run("All Filters", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"id", "name", "category", "supplier", "quantity", "status", "image", "created_at", "updated_at"}).
			AddRow(expectedMaterials[0].ID, expectedMaterials[0].Name, expectedMaterials[0].Category, expectedMaterials[0].Supplier, expectedMaterials[0].Quantity, expectedMaterials[0].Status, expectedMaterials[0].Image, expectedMaterials[0].CreatedAt, expectedMaterials[0].UpdatedAt)

		queryAll := "SELECT id, name, category, supplier, quantity, status, image, created_at, updated_at FROM materials WHERE 1=1 AND (name LIKE ? OR category LIKE ? OR supplier LIKE ? OR CAST(id AS CHAR) LIKE ?) AND category = ? AND supplier = ? AND status = ? ORDER BY created_at DESC"
		mock.ExpectQuery(regexp.QuoteMeta(queryAll)).WithArgs("%term%", "%term%", "%term%", "%term%", "Cat A", "Sup 1", "Active").WillReturnRows(rows)

		materials, err := repo.GetAll("term", "Cat A", "Sup 1", "Active")
		assert.NoError(t, err)
		assert.Equal(t, []models.Material{expectedMaterials[0]}, materials)
	})

	t.Run("Query Error", func(t *testing.T) {
		query := "SELECT id, name, category, supplier, quantity, status, image, created_at, updated_at FROM materials WHERE 1=1 ORDER BY created_at DESC"
		mock.ExpectQuery(regexp.QuoteMeta(query)).WillReturnError(sql.ErrConnDone)

		materials, err := repo.GetAll("", "", "", "")
		assert.ErrorIs(t, err, sql.ErrConnDone)
		assert.Nil(t, materials)
	})

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetMaterialByID(t *testing.T) {
	db, mock := NewMockDB(t)
	defer db.Close()
	repo := NewMaterialRepository(db)

	now := time.Now()
	expectedMaterial := &models.Material{ID: 1, Name: "Material 1", Category: "Cat A", Supplier: "Sup 1", Quantity: 10, Status: "Active", Image: "img1.jpg", CreatedAt: now, UpdatedAt: now}

	rows := sqlmock.NewRows([]string{"id", "name", "category", "supplier", "quantity", "status", "image", "created_at", "updated_at"}).
		AddRow(expectedMaterial.ID, expectedMaterial.Name, expectedMaterial.Category, expectedMaterial.Supplier, expectedMaterial.Quantity, expectedMaterial.Status, expectedMaterial.Image, expectedMaterial.CreatedAt, expectedMaterial.UpdatedAt)

	query := "SELECT id, name, category, supplier, quantity, status, image, created_at, updated_at FROM materials WHERE id = ?"

	// Test case 1: Found
	mock.ExpectQuery(regexp.QuoteMeta(query)).WithArgs(1).WillReturnRows(rows)
	material, err := repo.GetByID(1)
	assert.NoError(t, err)
	assert.Equal(t, expectedMaterial, material)

	// Test case 2: Not Found
	mock.ExpectQuery(regexp.QuoteMeta(query)).WithArgs(2).WillReturnError(sql.ErrNoRows)
	material, err = repo.GetByID(2)
	assert.NoError(t, err) // Expecting nil error for not found as per repo logic
	assert.Nil(t, material)

	// Test case 3: Query Error
	mock.ExpectQuery(regexp.QuoteMeta(query)).WithArgs(3).WillReturnError(sql.ErrConnDone)
	material, err = repo.GetByID(3)
	assert.ErrorIs(t, err, sql.ErrConnDone) // Be specific about the error type
	assert.Nil(t, material)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestCreateMaterial(t *testing.T) {
	db, mock := NewMockDB(t)
	defer db.Close()
	repo := NewMaterialRepository(db)

	newMaterial := &models.Material{
		Name:     "New Material",
		Category: "Cat C",
		Supplier: "Sup 3",
		Quantity: 5,
		Status:   "Pending",
		Image:    "new.jpg",
	}

	query := "INSERT INTO materials (name, category, supplier, quantity, status, image, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?)"

	// Test case 1: Success
	mock.ExpectExec(regexp.QuoteMeta(query)).
		WithArgs(newMaterial.Name, newMaterial.Category, newMaterial.Supplier, newMaterial.Quantity, newMaterial.Status, newMaterial.Image, sqlmock.AnyArg(), sqlmock.AnyArg()). // Use AnyArg for time.Now()
		WillReturnResult(sqlmock.NewResult(1, 1))                                                                                                                                // Insert ID 1, 1 row affected

	id, err := repo.Create(newMaterial)
	assert.NoError(t, err)
	assert.Equal(t, 1, id)

	// Test case 2: Exec Error
	mock.ExpectExec(regexp.QuoteMeta(query)).
		WithArgs(newMaterial.Name, newMaterial.Category, newMaterial.Supplier, newMaterial.Quantity, newMaterial.Status, newMaterial.Image, sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnError(sql.ErrConnDone)

	id, err = repo.Create(newMaterial)
	assert.ErrorIs(t, err, sql.ErrConnDone) // Be specific about the error type
	assert.Equal(t, 0, id)

	// Test case 3: Result Error (e.g., driver doesn't support LastInsertId)
	expectedErr := sql.ErrNoRows // Simulate the error returned by mock when getting LastInsertId fails
	mock.ExpectExec(regexp.QuoteMeta(query)).
		WithArgs(newMaterial.Name, newMaterial.Category, newMaterial.Supplier, newMaterial.Quantity, newMaterial.Status, newMaterial.Image, sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnResult(sqlmock.NewErrorResult(expectedErr)) // Simulate error getting LastInsertId

	id, err = repo.Create(newMaterial)
	assert.ErrorIs(t, err, expectedErr) // Check for the specific result error
	assert.Equal(t, 0, id)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUpdateMaterial(t *testing.T) {
	db, mock := NewMockDB(t)
	defer db.Close()
	repo := NewMaterialRepository(db)

	updatedMaterial := &models.Material{
		ID:       1,
		Name:     "Updated Material",
		Category: "Cat A Updated",
		Supplier: "Sup 1 Updated",
		Quantity: 15,
		Status:   "Inactive",
		Image:    "updated.jpg",
	}

	query := "UPDATE materials SET name = ?, category = ?, supplier = ?, quantity = ?, status = ?, image = ?, updated_at = ? WHERE id = ?"

	// Test case 1: Success
	mock.ExpectExec(regexp.QuoteMeta(query)).
		WithArgs(updatedMaterial.Name, updatedMaterial.Category, updatedMaterial.Supplier, updatedMaterial.Quantity, updatedMaterial.Status, updatedMaterial.Image, sqlmock.AnyArg(), updatedMaterial.ID). // Use AnyArg for time.Now()
		WillReturnResult(sqlmock.NewResult(0, 1))                                                                                                                                                          // 1 row affected

	err := repo.Update(updatedMaterial)
	assert.NoError(t, err)

	// Test case 2: Exec Error
	mock.ExpectExec(regexp.QuoteMeta(query)).
		WithArgs(updatedMaterial.Name, updatedMaterial.Category, updatedMaterial.Supplier, updatedMaterial.Quantity, updatedMaterial.Status, updatedMaterial.Image, sqlmock.AnyArg(), updatedMaterial.ID).
		WillReturnError(sql.ErrConnDone)

	err = repo.Update(updatedMaterial)
	assert.ErrorIs(t, err, sql.ErrConnDone) // Be specific about the error type

	// Test case 3: No Rows Affected (could indicate material not found, though current repo doesn't check)
	mock.ExpectExec(regexp.QuoteMeta(query)).
		WithArgs(updatedMaterial.Name, updatedMaterial.Category, updatedMaterial.Supplier, updatedMaterial.Quantity, updatedMaterial.Status, updatedMaterial.Image, sqlmock.AnyArg(), updatedMaterial.ID).
		WillReturnResult(sqlmock.NewResult(0, 0)) // 0 rows affected

	err = repo.Update(updatedMaterial)
	assert.NoError(t, err) // The repository currently doesn't return an error if no rows are affected

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestDeleteMaterial(t *testing.T) {
	db, mock := NewMockDB(t)
	defer db.Close()
	repo := NewMaterialRepository(db)

	materialID := 1
	query := "DELETE FROM materials WHERE id = ?"

	// Test case 1: Success
	mock.ExpectExec(regexp.QuoteMeta(query)).WithArgs(materialID).WillReturnResult(sqlmock.NewResult(0, 1)) // 1 row affected

	err := repo.Delete(materialID)
	assert.NoError(t, err)

	// Test case 2: Exec Error
	mock.ExpectExec(regexp.QuoteMeta(query)).WithArgs(materialID).WillReturnError(sql.ErrConnDone)

	err = repo.Delete(materialID)
	assert.ErrorIs(t, err, sql.ErrConnDone) // Be specific about the error type

	// Test case 3: No Rows Affected (could indicate material not found)
	mock.ExpectExec(regexp.QuoteMeta(query)).WithArgs(materialID).WillReturnResult(sqlmock.NewResult(0, 0)) // 0 rows affected

	err = repo.Delete(materialID)
	assert.NoError(t, err) // Repository doesn't currently return error for 0 rows affected

	assert.NoError(t, mock.ExpectationsWereMet())
}
