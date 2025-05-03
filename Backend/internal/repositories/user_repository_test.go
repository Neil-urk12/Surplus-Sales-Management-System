package repositories

import (
	"database/sql"
	"oop/internal/models"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
)

func TestUserRepository_Create(t *testing.T) {
	// Create a new mock database with exact query matching
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	// Create a new repository with the mock database
	repo := &UserRepository{
		dbClient: &DatabaseClient{DB: db},
	}

	// Create a test user
	user := &models.User{
		Id:       "test-id",
		FullName: "Test User",
		Email:    "test@example.com",
		Password: "password123",
		Role:     "user",
	}

	// Set up the expected SQL query and result
	mock.ExpectExec("INSERT INTO users (id, full_name, email, password_hash, role, created_at, updated_at, is_active) VALUES (?, ?, ?, ?, ?, ?, ?, ?)").
		WithArgs(
			user.Id,
			user.FullName,
			user.Email,
			sqlmock.AnyArg(), // Password will be hashed
			user.Role,
			sqlmock.AnyArg(), // CreatedAt will be set by the function
			sqlmock.AnyArg(), // UpdatedAt will be set by the function
			true,             // IsActive is set to true by default
		).
		WillReturnResult(sqlmock.NewResult(1, 1))

	// Call the function being tested
	err = repo.Create(user)

	// Assert that no errors occurred
	assert.NoError(t, err)

	// Verify that all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestUserRepository_GetByID(t *testing.T) {
	// Create a new mock database with exact query matching
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	// Create a new repository with the mock database
	repo := &UserRepository{
		dbClient: &DatabaseClient{DB: db},
	}

	// Set up test data
	userID := "test-id"
	now := time.Now()
	expectedUser := &models.User{
		Id:        userID,
		FullName:  "Test User",
		Email:     "test@example.com",
		Password:  "hashed_password",
		Role:      "user",
		CreatedAt: now,
		UpdatedAt: now,
		IsActive:  true,
	}

	// Set up the expected SQL query and result
	rows := sqlmock.NewRows([]string{"id", "full_name", "email", "password", "role", "created_at", "updated_at", "is_active"}).
		AddRow(expectedUser.Id, expectedUser.FullName, expectedUser.Email, expectedUser.Password, expectedUser.Role, expectedUser.CreatedAt, expectedUser.UpdatedAt, expectedUser.IsActive)

	mock.ExpectQuery("SELECT id, full_name, email, password_hash, role, created_at, updated_at, is_active FROM users WHERE id = ?").
		WithArgs(userID).
		WillReturnRows(rows)

	// Call the function being tested
	user, err := repo.GetByID(userID)

	// Assert that no errors occurred and the user matches the expected user
	assert.NoError(t, err)
	assert.Equal(t, expectedUser, user)

	// Verify that all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestUserRepository_GetByEmail(t *testing.T) {
	// Create a new mock database with exact query matching
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	// Create a new repository with the mock database
	repo := &UserRepository{
		dbClient: &DatabaseClient{DB: db},
	}

	// Set up test data
	email := "test@example.com"
	now := time.Now()
	expectedUser := &models.User{
		Id:        "test-id",
		FullName:  "Test User",
		Email:     email,
		Password:  "hashed_password",
		Role:      "user",
		CreatedAt: now,
		UpdatedAt: now,
		IsActive:  true,
	}

	// Set up the expected SQL query and result
	rows := sqlmock.NewRows([]string{"id", "full_name", "email", "password", "role", "created_at", "updated_at", "is_active"}).
		AddRow(expectedUser.Id, expectedUser.FullName, expectedUser.Email, expectedUser.Password, expectedUser.Role, expectedUser.CreatedAt, expectedUser.UpdatedAt, expectedUser.IsActive)

	mock.ExpectQuery("SELECT id, full_name, email, password_hash, role, created_at, updated_at, is_active FROM users WHERE email = ?").
		WithArgs(email).
		WillReturnRows(rows)

	// Call the function being tested
	user, err := repo.GetByEmail(email)

	// Assert that no errors occurred and the user matches the expected user
	assert.NoError(t, err)
	assert.Equal(t, expectedUser, user)

	// Verify that all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestUserRepository_Update(t *testing.T) {
	// Create a new mock database with exact query matching
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	// Create a new repository with the mock database
	repo := &UserRepository{
		dbClient: &DatabaseClient{DB: db},
	}

	// Create a test user
	user := &models.User{
		Id:       "test-id",
		FullName: "Updated User",
		Email:    "updated@example.com",
		Role:     "admin",
		IsActive: true,
	}

	// Set up the expected SQL query and result
	mock.ExpectExec("UPDATE users SET full_name = ?, email = ?, role = ?, updated_at = ?, is_active = ? WHERE id = ?").
		WithArgs(user.FullName, user.Email, user.Role, sqlmock.AnyArg(), user.IsActive, user.Id).
		WillReturnResult(sqlmock.NewResult(0, 1))

	// Call the function being tested
	err = repo.Update(user)

	// Assert that no errors occurred
	assert.NoError(t, err)

	// Verify that all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestUserRepository_UpdatePassword(t *testing.T) {
	// Create a new mock database with exact query matching
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	// Create a new repository with the mock database
	repo := &UserRepository{
		dbClient: &DatabaseClient{DB: db},
	}

	// Set up test data
	userID := "test-id"
	newPassword := "new_password"

	// Set up the expected SQL query and result
	mock.ExpectExec("UPDATE users SET password_hash = ?, updated_at = ? WHERE id = ?").
		WithArgs(
			sqlmock.AnyArg(), // Hashed password
			sqlmock.AnyArg(), // UpdatedAt timestamp
			userID,
		).
		WillReturnResult(sqlmock.NewResult(0, 1))

	// Call the function being tested
	err = repo.UpdatePassword(userID, newPassword)

	// Assert that no errors occurred
	assert.NoError(t, err)

	// Verify that all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestUserRepository_Delete(t *testing.T) {
	// Create a new mock database with exact query matching
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	// Create a new repository with the mock database
	repo := &UserRepository{
		dbClient: &DatabaseClient{DB: db},
	}

	// Set up test data
	userID := "test-id"

	// Set up the expected SQL query and result
	mock.ExpectExec("DELETE FROM users WHERE id = ?").
		WithArgs(userID).
		WillReturnResult(sqlmock.NewResult(0, 1))

	// Call the function being tested
	err = repo.Delete(userID)

	// Assert that no errors occurred
	assert.NoError(t, err)

	// Verify that all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestUserRepository_GetAll(t *testing.T) {
	// Create a new mock database with exact query matching
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	// Create a new repository with the mock database
	repo := &UserRepository{
		dbClient: &DatabaseClient{DB: db},
	}

	// Set up test data
	now := time.Now()
	user1 := &models.User{
		Id:        "id1",
		FullName:  "User One",
		Email:     "one@example.com",
		Password:  "pass1",
		Role:      "user",
		CreatedAt: now,
		UpdatedAt: now,
		IsActive:  true,
	}
	user2 := &models.User{
		Id:        "id2",
		FullName:  "User Two",
		Email:     "two@example.com",
		Password:  "pass2",
		Role:      "admin",
		CreatedAt: now,
		UpdatedAt: now,
		IsActive:  false,
	}
	expectedUsers := []*models.User{user1, user2}

	// Set up the expected SQL query and result
	rows := sqlmock.NewRows([]string{"id", "full_name", "email", "password", "role", "created_at", "updated_at", "is_active"}).
		AddRow(user1.Id, user1.FullName, user1.Email, user1.Password, user1.Role, user1.CreatedAt, user1.UpdatedAt, user1.IsActive).
		AddRow(user2.Id, user2.FullName, user2.Email, user2.Password, user2.Role, user2.CreatedAt, user2.UpdatedAt, user2.IsActive)

	mock.ExpectQuery("SELECT id, full_name, email, password_hash, role, created_at, updated_at, is_active FROM users ORDER BY created_at DESC").
		WillReturnRows(rows)

	// Call the function being tested
	users, err := repo.GetAll()

	// Assert that no errors occurred and the users match the expected users
	assert.NoError(t, err)
	assert.Equal(t, expectedUsers, users)

	// Verify that all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestUserRepository_VerifyPassword(t *testing.T) {
	// Create a new mock database with exact query matching
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	// Create a new repository with the mock database
	repo := &UserRepository{
		dbClient: &DatabaseClient{DB: db},
	}

	// Set up test data
	email := "test@example.com"
	password := "password123"
	now := time.Now()
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		t.Fatalf("Failed to generate password hash: %v", err)
	}
	hashedPassword := string(hashedBytes)
	expectedUser := &models.User{
		Id:        "test-id",
		FullName:  "Test User",
		Email:     email,
		Password:  hashedPassword,
		Role:      "user",
		CreatedAt: now,
		UpdatedAt: now,
		IsActive:  true,
	}

	// Set up the expected SQL query and result
	rows := sqlmock.NewRows([]string{"id", "full_name", "email", "password", "role", "created_at", "updated_at", "is_active"}).
		AddRow(expectedUser.Id, expectedUser.FullName, expectedUser.Email, expectedUser.Password, expectedUser.Role, expectedUser.CreatedAt, expectedUser.UpdatedAt, expectedUser.IsActive)

	mock.ExpectQuery("SELECT id, full_name, email, password_hash, role, created_at, updated_at, is_active FROM users WHERE email = ?").
		WithArgs(email).
		WillReturnRows(rows)

	// Call the function being tested
	user, err := repo.VerifyPassword(email, password)

	// Assert that no errors occurred and the user matches the expected user
	assert.NoError(t, err)
	assert.Equal(t, expectedUser.Id, user.Id)
	assert.Equal(t, expectedUser.Email, user.Email)

	// Verify that all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}

	// Test with incorrect password
	db2, mock2, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db2.Close()

	// Create a new repository with the new mock database
	repo2 := &UserRepository{
		dbClient: &DatabaseClient{DB: db2},
	}

	rows2 := sqlmock.NewRows([]string{"id", "full_name", "email", "password", "role", "created_at", "updated_at", "is_active"}).
		AddRow(expectedUser.Id, expectedUser.FullName, expectedUser.Email, expectedUser.Password, expectedUser.Role, expectedUser.CreatedAt, expectedUser.UpdatedAt, expectedUser.IsActive)

	mock2.ExpectQuery("SELECT id, full_name, email, password_hash, role, created_at, updated_at, is_active FROM users WHERE email = ?").
		WithArgs(email).
		WillReturnRows(rows2)

	// Call the function being tested with wrong password
	_, err = repo2.VerifyPassword(email, "wrong_password")

	// Assert that an error occurred
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid password")
}

func TestUserRepository_ActivateUser(t *testing.T) {
	// Create a new mock database with exact query matching
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	// Create a new repository with the mock database
	repo := &UserRepository{
		dbClient: &DatabaseClient{DB: db},
	}

	// Set up test data
	userID := "test-id"

	// Set up the expected SQL query and result
	mock.ExpectExec("UPDATE users SET is_active = ?, updated_at = ? WHERE id = ?").
		WithArgs(true, sqlmock.AnyArg(), userID).
		WillReturnResult(sqlmock.NewResult(0, 1))

	// Call the function being tested
	err = repo.ActivateUser(userID)

	// Assert that no errors occurred
	assert.NoError(t, err)

	// Verify that all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestUserRepository_DeactivateUser(t *testing.T) {
	// Create a new mock database with exact query matching
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	// Create a new repository with the mock database
	repo := &UserRepository{
		dbClient: &DatabaseClient{DB: db},
	}

	// Set up test data
	userID := "test-id"

	// Set up the expected SQL query and result
	mock.ExpectExec("UPDATE users SET is_active = ?, updated_at = ? WHERE id = ?").
		WithArgs(false, sqlmock.AnyArg(), userID).
		WillReturnResult(sqlmock.NewResult(0, 1))

	// Call the function being tested
	err = repo.DeactivateUser(userID)

	// Assert that no errors occurred
	assert.NoError(t, err)

	// Verify that all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestUserRepository_EmailExists(t *testing.T) {
	// Create a new mock database with exact query matching
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	// Create a new repository with the mock database
	repo := &UserRepository{
		dbClient: &DatabaseClient{DB: db},
	}

	// Set up test data
	email := "test@example.com"

	// Test when email exists
	mock.ExpectQuery("SELECT EXISTS(SELECT 1 FROM users WHERE email = ?)").
		WithArgs(email).
		WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(true))

	// Call the function being tested
	exists, err := repo.EmailExists(email)

	// Assert that no errors occurred and the email exists
	assert.NoError(t, err)
	assert.True(t, exists)

	// Verify that all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}

	// Test when email does not exist
	db2, mock2, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db2.Close()

	// Create a new repository with the new mock database
	repo2 := &UserRepository{
		dbClient: &DatabaseClient{DB: db2},
	}

	mock2.ExpectQuery("SELECT EXISTS(SELECT 1 FROM users WHERE email = ?)").
		WithArgs(email).
		WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(false))

	// Call the function being tested
	exists, err = repo2.EmailExists(email)

	// Assert that no errors occurred and the email does not exist
	assert.NoError(t, err)
	assert.False(t, exists)

	// Verify that all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestUserRepository_GetByID_NotFound(t *testing.T) {
	// Create a new mock database with exact query matching
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	// Create a new repository with the mock database
	repo := &UserRepository{
		dbClient: &DatabaseClient{DB: db},
	}

	// Set up test data
	userID := "non-existent-id"

	// Set up the expected SQL query to return no rows
	mock.ExpectQuery("SELECT id, full_name, email, password_hash, role, created_at, updated_at, is_active FROM users WHERE id = ?").
		WithArgs(userID).
		WillReturnError(sql.ErrNoRows)

	// Call the function being tested
	user, err := repo.GetByID(userID)

	// Assert that an error occurred and the user is nil
	assert.Error(t, err)
	assert.Nil(t, user)
	assert.Contains(t, err.Error(), "user not found")

	// Verify that all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
