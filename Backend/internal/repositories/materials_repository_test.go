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
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("With Search Term", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"id", "name", "category", "supplier", "quantity", "status", "image", "created_at", "updated_at"}).
			AddRow(expectedMaterials[0].ID, expectedMaterials[0].Name, expectedMaterials[0].Category, expectedMaterials[0].Supplier, expectedMaterials[0].Quantity, expectedMaterials[0].Status, expectedMaterials[0].Image, expectedMaterials[0].CreatedAt, expectedMaterials[0].UpdatedAt)

		querySearch := "SELECT id, name, category, supplier, quantity, status, image, created_at, updated_at FROM materials WHERE 1=1 AND name ILIKE ? ORDER BY created_at DESC"
		mock.ExpectQuery(regexp.QuoteMeta(querySearch)).WithArgs("%term%").WillReturnRows(rows)

		materials, err := repo.GetAll("term", "", "", "")
		assert.NoError(t, err)
		assert.Equal(t, []models.Material{expectedMaterials[0]}, materials)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("With Category", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"id", "name", "category", "supplier", "quantity", "status", "image", "created_at", "updated_at"}).
			AddRow(expectedMaterials[0].ID, expectedMaterials[0].Name, expectedMaterials[0].Category, expectedMaterials[0].Supplier, expectedMaterials[0].Quantity, expectedMaterials[0].Status, expectedMaterials[0].Image, expectedMaterials[0].CreatedAt, expectedMaterials[0].UpdatedAt)

		queryCategory := "SELECT id, name, category, supplier, quantity, status, image, created_at, updated_at FROM materials WHERE 1=1 AND category = ? ORDER BY created_at DESC"
		mock.ExpectQuery(regexp.QuoteMeta(queryCategory)).WithArgs("Cat A").WillReturnRows(rows)

		materials, err := repo.GetAll("", "Cat A", "", "")
		assert.NoError(t, err)
		assert.Equal(t, []models.Material{expectedMaterials[0]}, materials)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("With Supplier", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"id", "name", "category", "supplier", "quantity", "status", "image", "created_at", "updated_at"}).
			AddRow(expectedMaterials[0].ID, expectedMaterials[0].Name, expectedMaterials[0].Category, expectedMaterials[0].Supplier, expectedMaterials[0].Quantity, expectedMaterials[0].Status, expectedMaterials[0].Image, expectedMaterials[0].CreatedAt, expectedMaterials[0].UpdatedAt)

		querySupplier := "SELECT id, name, category, supplier, quantity, status, image, created_at, updated_at FROM materials WHERE 1=1 AND supplier = ? ORDER BY created_at DESC"
		mock.ExpectQuery(regexp.QuoteMeta(querySupplier)).WithArgs("Sup 1").WillReturnRows(rows)

		materials, err := repo.GetAll("", "", "Sup 1", "")
		assert.NoError(t, err)
		assert.Equal(t, []models.Material{expectedMaterials[0]}, materials)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("With Status", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"id", "name", "category", "supplier", "quantity", "status", "image", "created_at", "updated_at"}).
			AddRow(expectedMaterials[0].ID, expectedMaterials[0].Name, expectedMaterials[0].Category, expectedMaterials[0].Supplier, expectedMaterials[0].Quantity, expectedMaterials[0].Status, expectedMaterials[0].Image, expectedMaterials[0].CreatedAt, expectedMaterials[0].UpdatedAt)

		queryStatus := "SELECT id, name, category, supplier, quantity, status, image, created_at, updated_at FROM materials WHERE 1=1 AND status = ? ORDER BY created_at DESC"
		mock.ExpectQuery(regexp.QuoteMeta(queryStatus)).WithArgs("Active").WillReturnRows(rows)

		materials, err := repo.GetAll("", "", "", "Active")
		assert.NoError(t, err)
		assert.Equal(t, []models.Material{expectedMaterials[0]}, materials)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("All Filters", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"id", "name", "category", "supplier", "quantity", "status", "image", "created_at", "updated_at"}).
			AddRow(expectedMaterials[0].ID, expectedMaterials[0].Name, expectedMaterials[0].Category, expectedMaterials[0].Supplier, expectedMaterials[0].Quantity, expectedMaterials[0].Status, expectedMaterials[0].Image, expectedMaterials[0].CreatedAt, expectedMaterials[0].UpdatedAt)

		queryAll := "SELECT id, name, category, supplier, quantity, status, image, created_at, updated_at FROM materials WHERE 1=1 AND name ILIKE ? AND category = ? AND supplier = ? AND status = ? ORDER BY created_at DESC"
		mock.ExpectQuery(regexp.QuoteMeta(queryAll)).WithArgs("%term%", "Cat A", "Sup 1", "Active").WillReturnRows(rows)

		materials, err := repo.GetAll("term", "Cat A", "Sup 1", "Active")
		assert.NoError(t, err)
		assert.Equal(t, []models.Material{expectedMaterials[0]}, materials)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Query Error", func(t *testing.T) {
		query := "SELECT id, name, category, supplier, quantity, status, image, created_at, updated_at FROM materials WHERE 1=1 ORDER BY created_at DESC"
		mock.ExpectQuery(regexp.QuoteMeta(query)).WillReturnError(sql.ErrConnDone)

		materials, err := repo.GetAll("", "", "", "")
		assert.ErrorIs(t, err, sql.ErrConnDone) 
		assert.Nil(t, materials)
	})
}

func TestGetMaterialByID(t *testing.T) {
	db, mock := NewMockDB(t)
	defer db.Close()
	repo := NewMaterialRepository(db)

	now := time.Now()
	expectedMaterial := &models.Material{ID: 1, Name: "Material 1", Category: "Cat A", Supplier: "Sup 1", Quantity: 10, Status: "Active", Image: "img1.jpg", CreatedAt: now, UpdatedAt: now}

	query := "SELECT id, name, category, supplier, quantity, status, image, created_at, updated_at FROM materials WHERE id = ?"

	t.Run("Found", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"id", "name", "category", "supplier", "quantity", "status", "image", "created_at", "updated_at"}).
			AddRow(expectedMaterial.ID, expectedMaterial.Name, expectedMaterial.Category, expectedMaterial.Supplier, expectedMaterial.Quantity, expectedMaterial.Status, expectedMaterial.Image, expectedMaterial.CreatedAt, expectedMaterial.UpdatedAt)
		mock.ExpectQuery(regexp.QuoteMeta(query)).WithArgs(1).WillReturnRows(rows)
		material, err := repo.GetByID(1)
		assert.NoError(t, err)
		assert.Equal(t, expectedMaterial, material)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Not Found", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta(query)).WithArgs(2).WillReturnError(sql.ErrNoRows)
		material, err := repo.GetByID(2)
		assert.NoError(t, err) // GetByID itself doesn't return sql.ErrNoRows directly, it returns nil material
		assert.Nil(t, material)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Query Error", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta(query)).WithArgs(3).WillReturnError(sql.ErrConnDone)
		material, err := repo.GetByID(3)
		assert.ErrorIs(t, err, sql.ErrConnDone)
		assert.Nil(t, material)
		// No mock.ExpectationsWereMet() here
	})
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

	t.Run("Success", func(t *testing.T) {
		mock.ExpectExec(regexp.QuoteMeta(query)).
			WithArgs(newMaterial.Name, newMaterial.Category, newMaterial.Supplier, newMaterial.Quantity, newMaterial.Status, newMaterial.Image, sqlmock.AnyArg(), sqlmock.AnyArg()).
			WillReturnResult(sqlmock.NewResult(1, 1))

		id, err := repo.Create(newMaterial)
		assert.NoError(t, err)
		assert.Equal(t, 1, id)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Exec Error", func(t *testing.T) {
		mock.ExpectExec(regexp.QuoteMeta(query)).
			WithArgs(newMaterial.Name, newMaterial.Category, newMaterial.Supplier, newMaterial.Quantity, newMaterial.Status, newMaterial.Image, sqlmock.AnyArg(), sqlmock.AnyArg()).
			WillReturnError(sql.ErrConnDone)

		id, err := repo.Create(newMaterial)
		assert.ErrorIs(t, err, sql.ErrConnDone)
		assert.Equal(t, 0, id)
		// No mock.ExpectationsWereMet() here
	})

	t.Run("Result Error", func(t *testing.T) {
		expectedErr := sql.ErrNoRows // Simulate driver not supporting LastInsertId or other result errors
		mock.ExpectExec(regexp.QuoteMeta(query)).
			WithArgs(newMaterial.Name, newMaterial.Category, newMaterial.Supplier, newMaterial.Quantity, newMaterial.Status, newMaterial.Image, sqlmock.AnyArg(), sqlmock.AnyArg()).
			WillReturnResult(sqlmock.NewErrorResult(expectedErr))

		id, err := repo.Create(newMaterial)
		assert.ErrorIs(t, err, expectedErr)
		assert.Equal(t, 0, id)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
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

	t.Run("Success", func(t *testing.T) {
		mock.ExpectExec(regexp.QuoteMeta(query)).
			WithArgs(updatedMaterial.Name, updatedMaterial.Category, updatedMaterial.Supplier, updatedMaterial.Quantity, updatedMaterial.Status, updatedMaterial.Image, sqlmock.AnyArg(), updatedMaterial.ID).
			WillReturnResult(sqlmock.NewResult(0, 1)) // 1 row affected

		err := repo.Update(updatedMaterial)
		assert.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Exec Error", func(t *testing.T) {
		mock.ExpectExec(regexp.QuoteMeta(query)).
			WithArgs(updatedMaterial.Name, updatedMaterial.Category, updatedMaterial.Supplier, updatedMaterial.Quantity, updatedMaterial.Status, updatedMaterial.Image, sqlmock.AnyArg(), updatedMaterial.ID).
			WillReturnError(sql.ErrConnDone)

		err := repo.Update(updatedMaterial)
		assert.ErrorIs(t, err, sql.ErrConnDone)
		// No mock.ExpectationsWereMet() here
	})

	t.Run("No Rows Affected", func(t *testing.T) {
		mock.ExpectExec(regexp.QuoteMeta(query)).
			WithArgs(updatedMaterial.Name, updatedMaterial.Category, updatedMaterial.Supplier, updatedMaterial.Quantity, updatedMaterial.Status, updatedMaterial.Image, sqlmock.AnyArg(), updatedMaterial.ID).
			WillReturnResult(sqlmock.NewResult(0, 0)) // 0 rows affected

		err := repo.Update(updatedMaterial)
		assert.NoError(t, err) // Update itself doesn't error on 0 rows affected, usually
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestDeleteMaterial(t *testing.T) {
	db, mock := NewMockDB(t)
	defer db.Close()
	repo := NewMaterialRepository(db)

	materialID := 1
	query := "DELETE FROM materials WHERE id = ?"

	t.Run("Success", func(t *testing.T) {
		mock.ExpectExec(regexp.QuoteMeta(query)).WithArgs(materialID).WillReturnResult(sqlmock.NewResult(0, 1)) // 1 row affected

		err := repo.Delete(materialID)
		assert.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Exec Error", func(t *testing.T) {
		mock.ExpectExec(regexp.QuoteMeta(query)).WithArgs(materialID).WillReturnError(sql.ErrConnDone)

		err := repo.Delete(materialID)
		assert.ErrorIs(t, err, sql.ErrConnDone)
		// No mock.ExpectationsWereMet() here
	})

	t.Run("No Rows Affected", func(t *testing.T) {
		mock.ExpectExec(regexp.QuoteMeta(query)).WithArgs(materialID).WillReturnResult(sqlmock.NewResult(0, 0)) // 0 rows affected

		err := repo.Delete(materialID)
		assert.NoError(t, err) // Delete itself doesn't error on 0 rows affected
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}
