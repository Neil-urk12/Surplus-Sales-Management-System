package repositories

import (
	"database/sql"
	"fmt"
	"log"
	"oop/internal/models"
	"strings"
	"time"
)

// SalesRepository defines the interface for sales data operations
type SalesRepository interface {
	GetAllSales(filters map[string]interface{}) ([]models.Sale, error)
	GetSaleByID(id string) (*models.Sale, error)
	GetSalesByCustomerID(customerID string) ([]models.Sale, error)
	CreateSale(sale *models.Sale) (string, error)
	UpdateSale(sale *models.Sale) error
	DeleteSale(id string) error
	GetSaleItems(saleID string) ([]models.SaleItem, error)
	CreateSaleItem(item *models.SaleItem) (string, error)
	SellCab(cabID int, customerID string, quantity int, soldBy string, accessories []models.AccessoryForSale) (*models.Sale, error)
}

// salesRepository is a database implementation of SalesRepository
type salesRepository struct {
	DB *sql.DB
}

// NewSalesRepository creates a new instance of the sales repository
func NewSalesRepository(db *sql.DB) SalesRepository {
	return &salesRepository{DB: db}
}

// GetAllSales retrieves all sales records with optional filtering
func (r *salesRepository) GetAllSales(filters map[string]interface{}) ([]models.Sale, error) {
	query := `SELECT id, customer_id, sold_by, sale_date, total_price, created_at, updated_at FROM sales WHERE 1=1`
	args := []interface{}{}

	// Apply filters
	if customerID, ok := filters["customer_id"].(string); ok && customerID != "" {
		query += " AND customer_id = ?"
		args = append(args, customerID)
	}

	if soldBy, ok := filters["sold_by"].(string); ok && soldBy != "" {
		query += " AND sold_by = ?"
		args = append(args, soldBy)
	}

	if startDate, ok := filters["start_date"].(string); ok && startDate != "" {
		query += " AND sale_date >= ?"
		args = append(args, startDate)
	}

	if endDate, ok := filters["end_date"].(string); ok && endDate != "" {
		query += " AND sale_date <= ?"
		args = append(args, endDate)
	}

	query += " ORDER BY created_at DESC"

	rows, err := r.DB.Query(query, args...)
	if err != nil {
		log.Printf("Error querying sales: %v\nQuery: %s\nArgs: %v", err, query, args)
		return nil, err
	}
	defer rows.Close()

	var sales []models.Sale
	for rows.Next() {
		var sale models.Sale
		var createdAt, updatedAt time.Time

		if err := rows.Scan(
			&sale.ID,
			&sale.CustomerID,
			&sale.SoldBy,
			&sale.SaleDate,
			&sale.TotalPrice,
			&createdAt,
			&updatedAt,
		); err != nil {
			log.Printf("Error scanning sale row: %v", err)
			return nil, err
		}

		sale.CreatedAt = createdAt
		sale.UpdatedAt = updatedAt
		sales = append(sales, sale)
	}

	if err = rows.Err(); err != nil {
		log.Printf("Error iterating sale rows: %v", err)
		return nil, err
	}

	return sales, nil
}

// GetSaleByID retrieves a single sale by its ID
func (r *salesRepository) GetSaleByID(id string) (*models.Sale, error) {
	query := `SELECT id, customer_id, sold_by, sale_date, total_price, created_at, updated_at FROM sales WHERE id = ?`
	row := r.DB.QueryRow(query, id)

	var sale models.Sale
	var createdAt, updatedAt time.Time

	if err := row.Scan(
		&sale.ID,
		&sale.CustomerID,
		&sale.SoldBy,
		&sale.SaleDate,
		&sale.TotalPrice,
		&createdAt,
		&updatedAt,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("sale with ID %s not found", id)
		}
		log.Printf("Error scanning sale row by ID %s: %v", id, err)
		return nil, err
	}

	sale.CreatedAt = createdAt
	sale.UpdatedAt = updatedAt
	return &sale, nil
}

// GetSalesByCustomerID retrieves all sales for a specific customer
func (r *salesRepository) GetSalesByCustomerID(customerID string) ([]models.Sale, error) {
	filters := map[string]interface{}{
		"customer_id": customerID,
	}
	return r.GetAllSales(filters)
}

// CreateSale inserts a new sale record into the database
func (r *salesRepository) CreateSale(sale *models.Sale) (string, error) {
	query := `INSERT INTO sales (id, customer_id, sold_by, sale_date, total_price, created_at, updated_at) 
			VALUES (?, ?, ?, ?, ?, ?, ?)`

	// Generate a UUID if not provided
	if sale.ID == "" {
		// In a real implementation, you would generate a UUID here
		// For simplicity, we'll use a timestamp-based ID
		sale.ID = fmt.Sprintf("sale_%d", time.Now().UnixNano())
	}

	now := time.Now()
	sale.CreatedAt = now
	sale.UpdatedAt = now

	_, err := r.DB.Exec(
		query,
		sale.ID,
		sale.CustomerID,
		sale.SoldBy,
		sale.SaleDate,
		sale.TotalPrice,
		sale.CreatedAt,
		sale.UpdatedAt,
	)

	if err != nil {
		log.Printf("Error creating sale: %v", err)
		return "", err
	}

	return sale.ID, nil
}

// UpdateSale modifies an existing sale record in the database
func (r *salesRepository) UpdateSale(sale *models.Sale) error {
	// Check if the sale exists
	_, err := r.GetSaleByID(sale.ID)
	if err != nil {
		if strings.Contains(strings.ToLower(err.Error()), "not found") {
			return fmt.Errorf("sale with ID %s not found for update", sale.ID)
		}
		return err
	}

	query := `UPDATE sales 
			SET customer_id = ?, sold_by = ?, sale_date = ?, total_price = ?, updated_at = ? 
			WHERE id = ?`

	now := time.Now()
	sale.UpdatedAt = now

	_, err = r.DB.Exec(
		query,
		sale.CustomerID,
		sale.SoldBy,
		sale.SaleDate,
		sale.TotalPrice,
		now,
		sale.ID,
	)

	if err != nil {
		log.Printf("Error updating sale ID %s: %v", sale.ID, err)
		return err
	}

	return nil
}

// DeleteSale removes a sale record from the database
func (r *salesRepository) DeleteSale(id string) error {
	// Check if the sale exists
	_, err := r.GetSaleByID(id)
	if err != nil {
		if strings.Contains(strings.ToLower(err.Error()), "not found") {
			return fmt.Errorf("sale with ID %s not found for deletion", id)
		}
		return err
	}

	// Start a transaction to delete the sale and its items
	tx, err := r.DB.Begin()
	if err != nil {
		log.Printf("Error starting transaction for sale deletion: %v", err)
		return err
	}

	// Delete sale items first (foreign key constraint)
	_, err = tx.Exec("DELETE FROM sale_items WHERE sale_id = ?", id)
	if err != nil {
		tx.Rollback()
		log.Printf("Error deleting sale items for sale ID %s: %v", id, err)
		return err
	}

	// Delete the sale record
	_, err = tx.Exec("DELETE FROM sales WHERE id = ?", id)
	if err != nil {
		tx.Rollback()
		log.Printf("Error deleting sale ID %s: %v", id, err)
		return err
	}

	// Commit the transaction
	if err = tx.Commit(); err != nil {
		log.Printf("Error committing transaction for sale deletion: %v", err)
		return err
	}

	return nil
}

// GetSaleItems retrieves all items for a specific sale
func (r *salesRepository) GetSaleItems(saleID string) ([]models.SaleItem, error) {
	query := `SELECT id, sale_id, item_type, multi_cab_id, accessory_id, material_id, quantity, unit_price, subtotal, created_at, updated_at 
			FROM sale_items WHERE sale_id = ?`

	rows, err := r.DB.Query(query, saleID)
	if err != nil {
		log.Printf("Error querying sale items for sale ID %s: %v", saleID, err)
		return nil, err
	}
	defer rows.Close()

	var items []models.SaleItem
	for rows.Next() {
		var item models.SaleItem
		var createdAt, updatedAt time.Time

		if err := rows.Scan(
			&item.ID,
			&item.SaleID,
			&item.ItemType,
			&item.MultiCabID,
			&item.AccessoryID,
			&item.MaterialID,
			&item.Quantity,
			&item.UnitPrice,
			&item.Subtotal,
			&createdAt,
			&updatedAt,
		); err != nil {
			log.Printf("Error scanning sale item row: %v", err)
			return nil, err
		}

		item.CreatedAt = createdAt
		item.UpdatedAt = updatedAt
		items = append(items, item)
	}

	if err = rows.Err(); err != nil {
		log.Printf("Error iterating sale item rows: %v", err)
		return nil, err
	}

	return items, nil
}

// CreateSaleItem inserts a new sale item into the database
func (r *salesRepository) CreateSaleItem(item *models.SaleItem) (string, error) {
	query := `INSERT INTO sale_items (id, sale_id, item_type, multi_cab_id, accessory_id, material_id, quantity, unit_price, subtotal, created_at, updated_at) 
			VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

	// Generate a UUID if not provided
	if item.ID == "" {
		// In a real implementation, you would generate a UUID here
		// For simplicity, we'll use a timestamp-based ID
		item.ID = fmt.Sprintf("item_%d", time.Now().UnixNano())
	}

	now := time.Now()

	_, err := r.DB.Exec(
		query,
		item.ID,
		item.SaleID,
		item.ItemType,
		item.MultiCabID,
		item.AccessoryID,
		item.MaterialID,
		item.Quantity,
		item.UnitPrice,
		item.Subtotal,
		now,
		now,
	)

	if err != nil {
		log.Printf("Error creating sale item: %v", err)
		return "", err
	}

	return item.ID, nil
}

// SellCab handles the complete process of selling a cab with optional accessories
func (r *salesRepository) SellCab(cabID int, customerID string, quantity int, soldBy string, accessories []models.AccessoryForSale) (*models.Sale, error) {
	// Start a transaction
	tx, err := r.DB.Begin()
	if err != nil {
		log.Printf("Error starting transaction for cab sale: %v", err)
		return nil, err
	}

	// Get the cab details
	query := `SELECT id, name, price FROM multicabs WHERE id = ?`
	var cab struct {
		ID    int
		Name  string
		Price float64
	}

	err = tx.QueryRow(query, cabID).Scan(&cab.ID, &cab.Name, &cab.Price)
	if err != nil {
		tx.Rollback()
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("cab with ID %d not found", cabID)
		}
		log.Printf("Error getting cab details for ID %d: %v", cabID, err)
		return nil, err
	}

	// Calculate the total price
	cabTotal := cab.Price * float64(quantity)
	accessoriesTotal := 0.0
	for _, acc := range accessories {
		accessoriesTotal += acc.Price * float64(acc.Quantity)
	}
	totalPrice := cabTotal + accessoriesTotal

	// Create the sale record
	saleID := fmt.Sprintf("sale_%d", time.Now().UnixNano())
	saleDate := time.Now().Format("2006-01-02")

	_, err = tx.Exec(
		`INSERT INTO sales (id, customer_id, sold_by, sale_date, total_price, created_at, updated_at) 
		VALUES (?, ?, ?, ?, ?, ?, ?)`,
		saleID,
		customerID,
		soldBy,
		saleDate,
		totalPrice,
		time.Now(),
		time.Now(),
	)

	if err != nil {
		tx.Rollback()
		log.Printf("Error creating sale record: %v", err)
		return nil, err
	}

	// Add the cab as a sale item
	_, err = tx.Exec(
		`INSERT INTO sale_items (id, sale_id, item_type, multi_cab_id, quantity, unit_price, subtotal, created_at, updated_at) 
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		fmt.Sprintf("item_%d_cab", time.Now().UnixNano()),
		saleID,
		"cab",
		cabID,
		quantity,
		cab.Price,
		cabTotal,
		time.Now(),
		time.Now(),
	)

	if err != nil {
		tx.Rollback()
		log.Printf("Error creating cab sale item: %v", err)
		return nil, err
	}

	// Add each accessory as a sale item
	for _, acc := range accessories {
		_, err = tx.Exec(
			`INSERT INTO sale_items (id, sale_id, item_type, accessory_id, quantity, unit_price, subtotal, created_at, updated_at) 
			VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`,
			fmt.Sprintf("item_%d_acc_%d", time.Now().UnixNano(), acc.ID),
			saleID,
			"accessory",
			acc.ID,
			acc.Quantity,
			acc.Price,
			acc.Price*float64(acc.Quantity),
			time.Now(),
			time.Now(),
		)

		if err != nil {
			tx.Rollback()
			log.Printf("Error creating accessory sale item: %v", err)
			return nil, err
		}

		// Update the accessory inventory
		_, err = tx.Exec(
			"UPDATE accessories SET quantity = quantity - ?, updated_at = ? WHERE id = ?",
			acc.Quantity,
			time.Now(),
			acc.ID,
		)

		if err != nil {
			tx.Rollback()
			log.Printf("Error updating accessory inventory: %v", err)
			return nil, err
		}
	}

	// Update the cab inventory
	_, err = tx.Exec(
		"UPDATE multicabs SET quantity = quantity - ?, updated_at = ? WHERE id = ?",
		quantity,
		time.Now(),
		cabID,
	)

	if err != nil {
		tx.Rollback()
		log.Printf("Error updating cab inventory: %v", err)
		return nil, err
	}

	// Commit the transaction
	if err = tx.Commit(); err != nil {
		log.Printf("Error committing transaction for cab sale: %v", err)
		return nil, err
	}

	// Return the created sale
	sale := &models.Sale{
		ID:         saleID,
		CustomerID: customerID,
		SoldBy:     soldBy,
		SaleDate:   saleDate,
		TotalPrice: totalPrice,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	return sale, nil
}
