package repositories

import (
	"database/sql"
	"errors"
	"fmt"
	"oop/internal/models"
	"time"

	"github.com/google/uuid"
)

// CustomerRepository defines the interface for customer data operations.
type CustomerRepository interface {
	CreateCustomer(customer *models.Customer) (*models.Customer, error)
	GetCustomerByID(id string) (*models.Customer, error)
	GetAllCustomers() ([]*models.Customer, error)
	UpdateCustomer(customer *models.Customer) (*models.Customer, error)
	DeleteCustomer(id string) error
	GetCustomerByEmail(email string) (*models.Customer, error)
}

// customerRepository implements the CustomerRepository interface.
type customerRepository struct {
	DB *sql.DB
}

// NewCustomerRepository creates a new instance of customerRepository.
func NewCustomerRepository(db *sql.DB) CustomerRepository {
	return &customerRepository{DB: db}
}

// CreateCustomer adds a new customer to the database.
func (r *customerRepository) CreateCustomer(customer *models.Customer) (*models.Customer, error) {
	if customer.ID == "" {
		customer.ID = uuid.New().String()
	}

	now := time.Now()
	customer.DateRegistered = now.Format(time.RFC3339)
	customer.CreatedAt = now.Format(time.RFC3339)
	customer.UpdatedAt = now.Format(time.RFC3339)

	query := `
		INSERT INTO customers (id, name, email, phone, address, date_registered, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)
	`
	_, err := r.DB.Exec(
		query,
		customer.ID,
		customer.Name,
		customer.Email,
		customer.Phone,
		customer.Address,
		customer.DateRegistered,
		customer.CreatedAt,
		customer.UpdatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to create customer: %w", err)
	}

	return customer, nil
}

// GetCustomerByID retrieves a customer by their ID.
func (r *customerRepository) GetCustomerByID(id string) (*models.Customer, error) {
	query := `
		SELECT id, name, email, phone, address, date_registered, created_at, updated_at
		FROM customers
		WHERE id = ?
	`
	row := r.DB.QueryRow(query, id)

	var customer models.Customer
	err := row.Scan(
		&customer.ID,
		&customer.Name,
		&customer.Email,
		&customer.Phone,
		&customer.Address,
		&customer.DateRegistered,
		&customer.CreatedAt,
		&customer.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("customer with ID %s not found", id)
		}
		return nil, fmt.Errorf("failed to get customer by ID %s: %w", id, err)
	}

	return &customer, nil
}

// GetCustomerByEmail retrieves a customer by their email.
func (r *customerRepository) GetCustomerByEmail(email string) (*models.Customer, error) {
	query := `
		SELECT id, name, email, phone, address, date_registered, created_at, updated_at
		FROM customers
		WHERE email = ?
	`
	row := r.DB.QueryRow(query, email)

	var customer models.Customer
	err := row.Scan(
		&customer.ID,
		&customer.Name,
		&customer.Email,
		&customer.Phone,
		&customer.Address,
		&customer.DateRegistered,
		&customer.CreatedAt,
		&customer.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("customer with email %s not found", email)
		}
		return nil, fmt.Errorf("failed to get customer by email %s: %w", email, err)
	}

	return &customer, nil
}

// GetAllCustomers retrieves all customers from the database.
func (r *customerRepository) GetAllCustomers() ([]*models.Customer, error) {
	query := `
		SELECT id, name, email, phone, address, date_registered, created_at, updated_at
		FROM customers
		ORDER BY created_at DESC
	`
	rows, err := r.DB.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to query customers: %w", err)
	}
	defer rows.Close()

	var customers []*models.Customer
	for rows.Next() {
		var customer models.Customer
		err := rows.Scan(
			&customer.ID,
			&customer.Name,
			&customer.Email,
			&customer.Phone,
			&customer.Address,
			&customer.DateRegistered,
			&customer.CreatedAt,
			&customer.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan customer row: %w", err)
		}
		customers = append(customers, &customer)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating customer rows: %w", err)
	}

	return customers, nil
}

// UpdateCustomer updates an existing customer in the database.
func (r *customerRepository) UpdateCustomer(customer *models.Customer) (*models.Customer, error) {
	customer.UpdatedAt = time.Now().Format(time.RFC3339)

	query := `
		UPDATE customers
		SET name = ?, email = ?, phone = ?, address = ?, updated_at = ?
		WHERE id = ?
	`
	result, err := r.DB.Exec(
		query,
		customer.Name,
		customer.Email,
		customer.Phone,
		customer.Address,
		customer.UpdatedAt,
		customer.ID,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to update customer with ID %s: %w", customer.ID, err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return nil, fmt.Errorf("failed to get rows affected for customer ID %s: %w", customer.ID, err)
	}

	if rowsAffected == 0 {
		return nil, fmt.Errorf("customer with ID %s not found for update", customer.ID)
	}

	return customer, nil
}

// DeleteCustomer removes a customer from the database by their ID.
func (r *customerRepository) DeleteCustomer(id string) error {
	query := `DELETE FROM customers WHERE id = ?`
	result, err := r.DB.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete customer with ID %s: %w", id, err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected for customer ID %s: %w", id, err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("customer with ID %s not found for deletion", id)
	}

	return nil
}
