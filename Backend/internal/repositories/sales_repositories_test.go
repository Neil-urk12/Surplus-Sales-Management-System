package repositories

import (
	"database/sql"
	"fmt"
	"regexp"
	"testing"
	"time"

	"oop/internal/models"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetAllSales_NoFilters(t *testing.T) {
	db, mock := NewMockDB(t)
	defer db.Close()
	repo := NewSalesRepository(db)

	now := time.Now()
	expected := []models.Sale{
		{ID: "s1", CustomerID: "cust1", SoldBy: "user1", SaleDate: "2025-05-09", TotalPrice: 100.0, CreatedAt: now, UpdatedAt: now},
		{ID: "s2", CustomerID: "cust2", SoldBy: "user2", SaleDate: "2025-05-08", TotalPrice: 200.0, CreatedAt: now, UpdatedAt: now},
	}

	rows := sqlmock.NewRows([]string{"id", "customer_id", "sold_by", "sale_date", "total_price", "created_at", "updated_at"}).
		AddRow(expected[0].ID, expected[0].CustomerID, expected[0].SoldBy, expected[0].SaleDate, expected[0].TotalPrice, expected[0].CreatedAt, expected[0].UpdatedAt).
		AddRow(expected[1].ID, expected[1].CustomerID, expected[1].SoldBy, expected[1].SaleDate, expected[1].TotalPrice, expected[1].CreatedAt, expected[1].UpdatedAt)

	query := "SELECT id, customer_id, sold_by, sale_date, total_price, created_at, updated_at FROM sales WHERE 1=1 ORDER BY created_at DESC"
	mock.ExpectQuery(regexp.QuoteMeta(query)).WillReturnRows(rows)

	sales, err := repo.GetAllSales(nil)
	require.NoError(t, err)
	assert.Equal(t, expected, sales)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetSaleByID_Exists(t *testing.T) {
	db, mock := NewMockDB(t)
	defer db.Close()
	repo := NewSalesRepository(db)

	now := time.Now()
	expected := &models.Sale{ID: "s1", CustomerID: "cust1", SoldBy: "user1", SaleDate: "2025-05-09", TotalPrice: 150.0, CreatedAt: now, UpdatedAt: now}

	rows := sqlmock.NewRows([]string{"id", "customer_id", "sold_by", "sale_date", "total_price", "created_at", "updated_at"}).
		AddRow(expected.ID, expected.CustomerID, expected.SoldBy, expected.SaleDate, expected.TotalPrice, expected.CreatedAt, expected.UpdatedAt)

	query := "SELECT id, customer_id, sold_by, sale_date, total_price, created_at, updated_at FROM sales WHERE id = ?"
	mock.ExpectQuery(regexp.QuoteMeta(query)).WithArgs(expected.ID).WillReturnRows(rows)

	sale, err := repo.GetSaleByID(expected.ID)
	require.NoError(t, err)
	require.NotNil(t, sale)
	assert.Equal(t, expected, sale)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetSaleByID_NotExists(t *testing.T) {
	db, mock := NewMockDB(t)
	defer db.Close()
	repo := NewSalesRepository(db)

	id := "notfound"
	query := "SELECT id, customer_id, sold_by, sale_date, total_price, created_at, updated_at FROM sales WHERE id = ?"
	mock.ExpectQuery(regexp.QuoteMeta(query)).WithArgs(id).WillReturnError(sql.ErrNoRows)

	sale, err := repo.GetSaleByID(id)
	require.Error(t, err)
	assert.Nil(t, sale)
	assert.Contains(t, err.Error(), fmt.Sprintf("sale with ID %s not found", id))
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestCreateSale_Success(t *testing.T) {
	db, mock := NewMockDB(t)
	defer db.Close()
	repo := NewSalesRepository(db)

	s := &models.Sale{ID: "testid", CustomerID: "cust1", SoldBy: "user1", SaleDate: "2025-05-10", TotalPrice: 75.5}
	query := "INSERT INTO sales (id, customer_id, sold_by, sale_date, total_price, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?)"
	mock.ExpectExec(regexp.QuoteMeta(query)).
		WithArgs(s.ID, s.CustomerID, s.SoldBy, s.SaleDate, s.TotalPrice, sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnResult(sqlmock.NewResult(0, 1))

	id, err := repo.CreateSale(s)
	require.NoError(t, err)
	assert.Equal(t, s.ID, id)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUpdateSale_Success(t *testing.T) {
	db, mock := NewMockDB(t)
	defer db.Close()
	repo := NewSalesRepository(db)

	s := &models.Sale{ID: "s1", CustomerID: "cust2", SoldBy: "user2", SaleDate: "2025-05-11", TotalPrice: 120.0}
	// Mock existing sale lookup
	getQuery := "SELECT id, customer_id, sold_by, sale_date, total_price, created_at, updated_at FROM sales WHERE id = ?"
	now := time.Now().Add(-time.Hour)
	rowsGet := sqlmock.NewRows([]string{"id", "customer_id", "sold_by", "sale_date", "total_price", "created_at", "updated_at"}).
		AddRow(s.ID, s.CustomerID, s.SoldBy, s.SaleDate, s.TotalPrice, now, now)
	mock.ExpectQuery(regexp.QuoteMeta(getQuery)).WithArgs(s.ID).WillReturnRows(rowsGet)

	// Mock update
	updateQuery := "UPDATE sales SET customer_id = ?, sold_by = ?, sale_date = ?, total_price = ?, updated_at = ? WHERE id = ?"
	mock.ExpectExec(regexp.QuoteMeta(updateQuery)).
		WithArgs(s.CustomerID, s.SoldBy, s.SaleDate, s.TotalPrice, sqlmock.AnyArg(), s.ID).
		WillReturnResult(sqlmock.NewResult(0, 1))

	err := repo.UpdateSale(s)
	require.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUpdateSale_NotExists(t *testing.T) {
	db, mock := NewMockDB(t)
	defer db.Close()
	repo := NewSalesRepository(db)

	s := &models.Sale{ID: "notexists"}
	getQuery := "SELECT id, customer_id, sold_by, sale_date, total_price, created_at, updated_at FROM sales WHERE id = ?"
	mock.ExpectQuery(regexp.QuoteMeta(getQuery)).WithArgs(s.ID).WillReturnError(sql.ErrNoRows)

	err := repo.UpdateSale(s)
	require.Error(t, err)
	assert.Contains(t, err.Error(), fmt.Sprintf("sale with ID %s not found for update", s.ID))
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestDeleteSale_Success(t *testing.T) {
	db, mock := NewMockDB(t)
	defer db.Close()
	repo := NewSalesRepository(db)

	id := "sdel"
	// Mock existing sale lookup
	getQuery := "SELECT id, customer_id, sold_by, sale_date, total_price, created_at, updated_at FROM sales WHERE id = ?"
	now := time.Now()
	rowsGet := sqlmock.NewRows([]string{"id", "customer_id", "sold_by", "sale_date", "total_price", "created_at", "updated_at"}).
		AddRow(id, "cust", "user", "2025-05-09", 50.0, now, now)
	mock.ExpectQuery(regexp.QuoteMeta(getQuery)).WithArgs(id).WillReturnRows(rowsGet)

	// Mock transaction
	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta("DELETE FROM sale_items WHERE sale_id = ?")).WithArgs(id).WillReturnResult(sqlmock.NewResult(0, 2))
	mock.ExpectExec(regexp.QuoteMeta("DELETE FROM sales WHERE id = ?")).WithArgs(id).WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectCommit()

	err := repo.DeleteSale(id)
	require.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestDeleteSale_NotExists(t *testing.T) {
	db, mock := NewMockDB(t)
	defer db.Close()
	repo := NewSalesRepository(db)

	id := "none"
	getQuery := "SELECT id, customer_id, sold_by, sale_date, total_price, created_at, updated_at FROM sales WHERE id = ?"
	mock.ExpectQuery(regexp.QuoteMeta(getQuery)).WithArgs(id).WillReturnError(sql.ErrNoRows)

	err := repo.DeleteSale(id)
	require.Error(t, err)
	assert.Contains(t, err.Error(), fmt.Sprintf("sale with ID %s not found for deletion", id))
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetSaleItems(t *testing.T) {
	db, mock := NewMockDB(t)
	defer db.Close()
	repo := NewSalesRepository(db)

	saleID := "sitems"
	now := time.Now()
	expected := []models.SaleItem{
		{ID: "i1", SaleID: saleID, ItemType: "cab", MultiCabID: "10", AccessoryID: "", MaterialID: "", Quantity: 1, UnitPrice: 100.0, Subtotal: 100.0, CreatedAt: now, UpdatedAt: now},
	}

	rows := sqlmock.NewRows([]string{"id", "sale_id", "item_type", "multi_cab_id", "accessory_id", "material_id", "quantity", "unit_price", "subtotal", "created_at", "updated_at"}).
		AddRow(expected[0].ID, expected[0].SaleID, expected[0].ItemType, expected[0].MultiCabID, expected[0].AccessoryID, expected[0].MaterialID, expected[0].Quantity, expected[0].UnitPrice, expected[0].Subtotal, expected[0].CreatedAt, expected[0].UpdatedAt)

	query := "SELECT id, sale_id, item_type, multi_cab_id, accessory_id, material_id, quantity, unit_price, subtotal, created_at, updated_at FROM sale_items WHERE sale_id = ?"
	mock.ExpectQuery(regexp.QuoteMeta(query)).WithArgs(saleID).WillReturnRows(rows)

	items, err := repo.GetSaleItems(saleID)
	require.NoError(t, err)
	assert.Equal(t, expected, items)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestCreateSaleItem(t *testing.T) {
	db, mock := NewMockDB(t)
	defer db.Close()
	repo := NewSalesRepository(db)

	item := &models.SaleItem{ID: "item1", SaleID: "sale1", ItemType: "cab", MultiCabID: "5", AccessoryID: "", MaterialID: "", Quantity: 2, UnitPrice: 300.0, Subtotal: 600.0}
	query := "INSERT INTO sale_items (id, sale_id, item_type, multi_cab_id, accessory_id, material_id, quantity, unit_price, subtotal, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)"
	mock.ExpectExec(regexp.QuoteMeta(query)).
		WithArgs(item.ID, item.SaleID, item.ItemType, item.MultiCabID, item.AccessoryID, item.MaterialID, item.Quantity, item.UnitPrice, item.Subtotal, sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnResult(sqlmock.NewResult(0, 1))

	id, err := repo.CreateSaleItem(item)
	require.NoError(t, err)
	assert.Equal(t, item.ID, id)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestSellCab_NoAccessories(t *testing.T) {
	db, mock := NewMockDB(t)
	defer db.Close()
	repo := NewSalesRepository(db)

	cabID := 10
	customer := "cust"
	user := "user"
	quantity := 2
	// Mock transaction and queries
	mock.ExpectBegin()
	// Mock cab lookup
	cabPrice := 500.0
	mock.ExpectQuery(regexp.QuoteMeta("SELECT id, name, price FROM multicabs WHERE id = ?")).WithArgs(cabID).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name", "price"}).AddRow(cabID, "Test", cabPrice))
	// Mock create sale
	mock.ExpectExec(regexp.QuoteMeta("INSERT INTO sales (id, customer_id, sold_by, sale_date, total_price, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?)")).WithArgs(sqlmock.AnyArg(), customer, user, sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).WillReturnResult(sqlmock.NewResult(0, 1))
	// Mock cab sale item insert
	mock.ExpectExec(regexp.QuoteMeta("INSERT INTO sale_items (id, sale_id, item_type, multi_cab_id, quantity, unit_price, subtotal, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)")).WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), "cab", cabID, quantity, cabPrice, cabPrice*float64(quantity), sqlmock.AnyArg(), sqlmock.AnyArg()).WillReturnResult(sqlmock.NewResult(0, 1))
	// Mock update cab inventory
	mock.ExpectExec(regexp.QuoteMeta("UPDATE multicabs SET quantity = quantity - ?, updated_at = ? WHERE id = ?")).WithArgs(quantity, sqlmock.AnyArg(), cabID).WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectCommit()

	sale, err := repo.SellCab(cabID, customer, quantity, user, nil)
	require.NoError(t, err)
	assert.Equal(t, customer, sale.CustomerID)
	assert.Equal(t, user, sale.SoldBy)
	assert.WithinDuration(t, time.Now(), sale.CreatedAt, 5*time.Second)
	assert.WithinDuration(t, time.Now(), sale.UpdatedAt, 5*time.Second)
	expectedTotal := cabPrice * float64(quantity)
	assert.Equal(t, expectedTotal, sale.TotalPrice)
	assert.NoError(t, mock.ExpectationsWereMet())
}
