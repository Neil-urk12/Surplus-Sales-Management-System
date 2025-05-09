package repositories_test

import (
	"database/sql"
	"errors"
	"fmt"
	"testing"
	"time"

	"oop/internal/models"
	"oop/internal/repositories"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// newMockCustomerRepo is a helper to set up the mock database and repository
func newMockCustomerRepo(t *testing.T) (repositories.CustomerRepository, sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	require.NoError(t, err, "Failed to create sqlmock")
	t.Cleanup(func() { db.Close() })
	repo := repositories.NewCustomerRepository(db)
	return repo, mock
}

func TestCreateCustomer(t *testing.T) {
	repo, mock := newMockCustomerRepo(t)

	tests := []struct {
		name          string
		customer      *models.Customer
		mockSetup     func(mock sqlmock.Sqlmock, customer *models.Customer)
		expectError   bool
		errorContains string
	}{
		{
			name: "Success - new ID generated",
			customer: &models.Customer{
				FullName:  "Test User",
				Email:     "test@example.com",
				Phone:     "1234567890",
				Address:   "123 Test St",
			},
			mockSetup: func(mock sqlmock.Sqlmock, customer *models.Customer) {
				mock.ExpectExec("INSERT INTO customers (id, full_name, email, phone, address, date_registered, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?)").
					WithArgs(sqlmock.AnyArg(), customer.FullName, customer.Email, customer.Phone, customer.Address, sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
					WillReturnResult(sqlmock.NewResult(1, 1))
			},
			expectError: false,
		},
		{
			name: "Success - with existing ID",
			customer: &models.Customer{
				ID:       uuid.New().String(),
				FullName: "Test User Existing ID",
				Email:     "testexisting@example.com",
				Phone:     "0987654321",
				Address:   "456 Test Ave",
			},
			mockSetup: func(mock sqlmock.Sqlmock, customer *models.Customer) {
				mock.ExpectExec("INSERT INTO customers (id, full_name, email, phone, address, date_registered, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?)").
					WithArgs(customer.ID, customer.FullName, customer.Email, customer.Phone, customer.Address, sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
					WillReturnResult(sqlmock.NewResult(1, 1))
			},
			expectError: false,
		},
		{
			name: "Error - database exec fails",
			customer: &models.Customer{
				FullName: "Error User",
				Email:    "error@example.com",
			},
			mockSetup: func(mock sqlmock.Sqlmock, customer *models.Customer) {
				mock.ExpectExec("INSERT INTO customers (id, full_name, email, phone, address, date_registered, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?)").
					WithArgs(sqlmock.AnyArg(), customer.FullName, customer.Email, customer.Phone, customer.Address, sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
					WillReturnError(errors.New("db error"))
			},
			expectError:   true,
			errorContains: "failed to create customer: db error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup(mock, tt.customer)
			createdCustomer, err := repo.CreateCustomer(tt.customer)

			if tt.expectError {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errorContains)
				assert.Nil(t, createdCustomer)
			} else {
				assert.NoError(t, err)
				require.NotNil(t, createdCustomer)
				if tt.customer.ID != "" {
					assert.Equal(t, tt.customer.ID, createdCustomer.ID)
				} else {
					assert.NotEmpty(t, createdCustomer.ID)
				}
				assert.Equal(t, tt.customer.FullName, createdCustomer.FullName)
				assert.False(t, createdCustomer.DateRegistered.IsZero())
				assert.False(t, createdCustomer.CreatedAt.IsZero())
				assert.False(t, createdCustomer.UpdatedAt.IsZero())
			}
			assert.NoError(t, mock.ExpectationsWereMet(), "Sqlmock expectations not met")
		})
	}
}

func TestGetCustomerByID(t *testing.T) {
	repo, mock := newMockCustomerRepo(t)
	customerID := uuid.New().String()

	tests := []struct {
		name           string
		mockSetup      func(mock sqlmock.Sqlmock)
		expectCustomer *models.Customer
		expectError    bool
		errorContains  string
	}{
		{
			name: "Success",
			mockSetup: func(mock sqlmock.Sqlmock) {
				query := `SELECT id, full_name, email, phone, address, date_registered, created_at, updated_at FROM customers WHERE id = ?`
				rows := sqlmock.NewRows([]string{"id", "full_name", "email", "phone", "address", "date_registered", "created_at", "updated_at"}).
					AddRow(customerID, "Test User", "get@example.com", "111", "Addr1", time.Now(), time.Now(), time.Now())
				mock.ExpectQuery(query).WithArgs(customerID).WillReturnRows(rows)
			},
			expectCustomer: &models.Customer{ID: customerID, FullName: "Test User", Email: "get@example.com", Phone: "111", Address: "Addr1"},
			expectError:    false,
		},
		{
			name: "Not Found",
			mockSetup: func(mock sqlmock.Sqlmock) {
				query := `SELECT id, full_name, email, phone, address, date_registered, created_at, updated_at FROM customers WHERE id = ?`
				mock.ExpectQuery(query).WithArgs(customerID).WillReturnError(sql.ErrNoRows)
			},
			expectError:   true,
			errorContains: fmt.Sprintf("customer with ID %s not found", customerID),
		},
		{
			name: "Scan Error",
			mockSetup: func(mock sqlmock.Sqlmock) {
				query := `SELECT id, full_name, email, phone, address, date_registered, created_at, updated_at FROM customers WHERE id = ?`
				rows := sqlmock.NewRows([]string{"id", "full_name"}).AddRow(customerID, "Test User") // Mismatched columns
				mock.ExpectQuery(query).WithArgs(customerID).WillReturnRows(rows)
			},
			expectError:   true,
			errorContains: fmt.Sprintf("failed to get customer by ID %s", customerID),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup(mock)
			customer, err := repo.GetCustomerByID(customerID)

			if tt.expectError {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errorContains)
				assert.Nil(t, customer)
			} else {
				assert.NoError(t, err)
				require.NotNil(t, customer)
				assert.Equal(t, tt.expectCustomer.ID, customer.ID)
				assert.Equal(t, tt.expectCustomer.FullName, customer.FullName)
				assert.False(t, customer.DateRegistered.IsZero())
				assert.False(t, customer.CreatedAt.IsZero())
				assert.False(t, customer.UpdatedAt.IsZero())
			}
			assert.NoError(t, mock.ExpectationsWereMet(), "Sqlmock expectations not met")
		})
	}
}

func TestGetCustomerByEmail(t *testing.T) {
	repo, mock := newMockCustomerRepo(t)
	customerEmail := "testemail@example.com"

	tests := []struct {
		name           string
		mockSetup      func(mock sqlmock.Sqlmock)
		expectCustomer *models.Customer
		expectError    bool
		errorContains  string
	}{
		{
			name: "Success",
			mockSetup: func(mock sqlmock.Sqlmock) {
				query := `SELECT id, full_name, email, phone, address, date_registered, created_at, updated_at FROM customers WHERE email = ?`
				rows := sqlmock.NewRows([]string{"id", "full_name", "email", "phone", "address", "date_registered", "created_at", "updated_at"}).
					AddRow(uuid.New().String(), "Email User", customerEmail, "222", "Addr2", time.Now(), time.Now(), time.Now())
				mock.ExpectQuery(query).WithArgs(customerEmail).WillReturnRows(rows)
			},
			expectCustomer: &models.Customer{FullName: "Email User", Email: customerEmail, Phone: "222", Address: "Addr2"},
			expectError:    false,
		},
		{
			name: "Not Found",
			mockSetup: func(mock sqlmock.Sqlmock) {
				query := `SELECT id, full_name, email, phone, address, date_registered, created_at, updated_at FROM customers WHERE email = ?`
				mock.ExpectQuery(query).WithArgs(customerEmail).WillReturnError(sql.ErrNoRows)
			},
			expectError:   true,
			errorContains: fmt.Sprintf("customer with email %s not found", customerEmail),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup(mock)
			customer, err := repo.GetCustomerByEmail(customerEmail)

			if tt.expectError {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errorContains)
				assert.Nil(t, customer)
			} else {
				assert.NoError(t, err)
				require.NotNil(t, customer)
				assert.Equal(t, tt.expectCustomer.FullName, customer.FullName)
				assert.Equal(t, tt.expectCustomer.Email, customer.Email)
			}
			assert.NoError(t, mock.ExpectationsWereMet(), "Sqlmock expectations not met")
		})
	}
}

func TestGetAllCustomers(t *testing.T) {
	repo, mock := newMockCustomerRepo(t)

	tests := []struct {
		name          string
		mockSetup     func(mock sqlmock.Sqlmock)
		expectCount   int
		expectError   bool
		errorContains string
	}{
		{
			name: "Success - multiple customers",
			mockSetup: func(mock sqlmock.Sqlmock) {
				query := `SELECT id, full_name, email, phone, address, date_registered, created_at, updated_at FROM customers ORDER BY created_at DESC`
				rows := sqlmock.NewRows([]string{"id", "full_name", "email", "phone", "address", "date_registered", "created_at", "updated_at"}).
					AddRow(uuid.New().String(), "User 1", "u1@example.com", "", "", time.Now(), time.Now(), time.Now()).
					AddRow(uuid.New().String(), "User 2", "u2@example.com", "", "", time.Now(), time.Now(), time.Now())
				mock.ExpectQuery(query).WillReturnRows(rows)
			},
			expectCount: 2,
			expectError: false,
		},
		{
			name: "Success - no customers",
			mockSetup: func(mock sqlmock.Sqlmock) {
				query := `SELECT id, full_name, email, phone, address, date_registered, created_at, updated_at FROM customers ORDER BY created_at DESC`
				rows := sqlmock.NewRows([]string{"id", "full_name", "email", "phone", "address", "date_registered", "created_at", "updated_at"})
				mock.ExpectQuery(query).WillReturnRows(rows)
			},
			expectCount: 0,
			expectError: false,
		},
		{
			name: "Error - query fails",
			mockSetup: func(mock sqlmock.Sqlmock) {
				query := `SELECT id, full_name, email, phone, address, date_registered, created_at, updated_at FROM customers ORDER BY created_at DESC`
				mock.ExpectQuery(query).WillReturnError(errors.New("db query error"))
			},
			expectError:   true,
			errorContains: "failed to query customers: db query error",
		},
		{
			name: "Error - scan fails",
			mockSetup: func(mock sqlmock.Sqlmock) {
				query := `SELECT id, full_name, email, phone, address, date_registered, created_at, updated_at FROM customers ORDER BY created_at DESC`
				rows := sqlmock.NewRows([]string{"id", "full_name"}).AddRow(uuid.New().String(), "User 1") // Mismatched columns
				mock.ExpectQuery(query).WillReturnRows(rows)
			},
			expectError:   true,
			errorContains: "failed to scan customer row",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup(mock)
			customers, err := repo.GetAllCustomers()

			if tt.expectError {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errorContains)
				assert.Nil(t, customers)
			} else {
				assert.NoError(t, err)
				assert.Len(t, customers, tt.expectCount)
			}
			assert.NoError(t, mock.ExpectationsWereMet(), "Sqlmock expectations not met")
		})
	}
}

func TestUpdateCustomer(t *testing.T) {
	repo, mock := newMockCustomerRepo(t)
	customerToUpdate := &models.Customer{
		ID:       uuid.New().String(),
		FullName: "Updated Name",
		Email:    "updated@example.com",
		Phone:    "333",
		Address:  "Updated Addr",
	}

	tests := []struct {
		name          string
		mockSetup     func(mock sqlmock.Sqlmock)
		expectError   bool
		errorContains string
	}{
		{
			name: "Success",
			mockSetup: func(mock sqlmock.Sqlmock) {
				query := `UPDATE customers SET full_name = ?, email = ?, phone = ?, address = ?, updated_at = ? WHERE id = ?`
				mock.ExpectExec(query).
					WithArgs(customerToUpdate.FullName, customerToUpdate.Email, customerToUpdate.Phone, customerToUpdate.Address, sqlmock.AnyArg(), customerToUpdate.ID).
					WillReturnResult(sqlmock.NewResult(0, 1))
			},
			expectError: false,
		},
		{
			name: "Not Found",
			mockSetup: func(mock sqlmock.Sqlmock) {
				query := `UPDATE customers SET full_name = ?, email = ?, phone = ?, address = ?, updated_at = ? WHERE id = ?`
				mock.ExpectExec(query).
					WithArgs(customerToUpdate.FullName, customerToUpdate.Email, customerToUpdate.Phone, customerToUpdate.Address, sqlmock.AnyArg(), customerToUpdate.ID).
					WillReturnResult(sqlmock.NewResult(0, 0)) // 0 rows affected
			},
			expectError:   true,
			errorContains: fmt.Sprintf("customer with ID %s not found for update", customerToUpdate.ID),
		},
		{
			name: "DB Error",
			mockSetup: func(mock sqlmock.Sqlmock) {
				query := `UPDATE customers SET full_name = ?, email = ?, phone = ?, address = ?, updated_at = ? WHERE id = ?`
				mock.ExpectExec(query).
					WithArgs(customerToUpdate.FullName, customerToUpdate.Email, customerToUpdate.Phone, customerToUpdate.Address, sqlmock.AnyArg(), customerToUpdate.ID).
					WillReturnError(errors.New("db update error"))
			},
			expectError:   true,
			errorContains: fmt.Sprintf("failed to update customer with ID %s: db update error", customerToUpdate.ID),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup(mock)
			updatedCustomer, err := repo.UpdateCustomer(customerToUpdate)

			if tt.expectError {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errorContains)
				assert.Nil(t, updatedCustomer)
			} else {
				assert.NoError(t, err)
				require.NotNil(t, updatedCustomer)
				assert.Equal(t, customerToUpdate.ID, updatedCustomer.ID)
				assert.Equal(t, customerToUpdate.FullName, updatedCustomer.FullName)
				assert.False(t, updatedCustomer.UpdatedAt.IsZero())
			}
			assert.NoError(t, mock.ExpectationsWereMet(), "Sqlmock expectations not met")
		})
	}
}

func TestDeleteCustomer(t *testing.T) {
	repo, mock := newMockCustomerRepo(t)
	customerID := uuid.New().String()

	tests := []struct {
		name          string
		mockSetup     func(mock sqlmock.Sqlmock)
		expectError   bool
		errorContains string
	}{
		{
			name: "Success",
			mockSetup: func(mock sqlmock.Sqlmock) {
				query := `DELETE FROM customers WHERE id = ?`
				mock.ExpectExec(query).WithArgs(customerID).WillReturnResult(sqlmock.NewResult(0, 1))
			},
			expectError: false,
		},
		{
			name: "Not Found",
			mockSetup: func(mock sqlmock.Sqlmock) {
				query := `DELETE FROM customers WHERE id = ?`
				mock.ExpectExec(query).WithArgs(customerID).WillReturnResult(sqlmock.NewResult(0, 0)) // 0 rows affected
			},
			expectError:   true,
			errorContains: fmt.Sprintf("customer with ID %s not found for deletion", customerID),
		},
		{
			name: "DB Error",
			mockSetup: func(mock sqlmock.Sqlmock) {
				query := `DELETE FROM customers WHERE id = ?`
				mock.ExpectExec(query).WithArgs(customerID).WillReturnError(errors.New("db delete error"))
			},
			expectError:   true,
			errorContains: fmt.Sprintf("failed to delete customer with ID %s: db delete error", customerID),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup(mock)
			err := repo.DeleteCustomer(customerID)

			if tt.expectError {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errorContains)
			} else {
				assert.NoError(t, err)
			}
			assert.NoError(t, mock.ExpectationsWereMet(), "Sqlmock expectations not met")
		})
	}
}
