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
	// Create a new mock database
	db, mock, err := sqlmock.New()
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
		Name:     "Test User",
		Email:    "test@example.com",
		Password: "password123",
		Role:     "user",
		Token:    "test-token",
	}

	// Set up the expected SQL query and result
	mock.ExpectExec("INSERT INTO users").
		WithArgs(
			user.Id,
			user.Name,
			user.Email,
			sqlmock.AnyArg(), // Password will be hashed
			user.Role,
			sqlmock.AnyArg(), // CreatedAt will be set by the function
			sqlmock.AnyArg(), // UpdatedAt will be set by the function
			true,             // IsActive is set to true by default
			user.Token,
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
	// Create a new mock database
	db, mock, err := sqlmock.New()
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
	expectedUser := &models.User{
		Id:        userID,
		Name:      "Test User",
		Email:     "test@example.com",
		Password:  "hashed_password",
		Role:      "user",
		CreatedAt: time.Now().Format(time.RFC3339),
		UpdatedAt: time.Now().Format(time.RFC3339),
		IsActive:  true,
		Token:     "test-token",
	}

	// Set up the expected SQL query and result
	rows := sqlmock.NewRows([]string{"id", "name", "email", "password", "role", "created_at", "updated_at", "is_active", "token"}).
		AddRow(expectedUser.Id, expectedUser.Name, expectedUser.Email, expectedUser.Password, expectedUser.Role, expectedUser.CreatedAt, expectedUser.UpdatedAt, expectedUser.IsActive, expectedUser.Token)

	mock.ExpectQuery("SELECT (.+) FROM users WHERE id = ?").
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
	// Create a new mock database
	db, mock, err := sqlmock.New()
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
	expectedUser := &models.User{
		Id:        "test-id",
		Name:      "Test User",
		Email:     email,
		Password:  "hashed_password",
		Role:      "user",
		CreatedAt: time.Now().Format(time.RFC3339),
		UpdatedAt: time.Now().Format(time.RFC3339),
		IsActive:  true,
		Token:     "test-token",
	}

	// Set up the expected SQL query and result
	rows := sqlmock.NewRows([]string{"id", "name", "email", "password", "role", "created_at", "updated_at", "is_active", "token"}).
		AddRow(expectedUser.Id, expectedUser.Name, expectedUser.Email, expectedUser.Password, expectedUser.Role, expectedUser.CreatedAt, expectedUser.UpdatedAt, expectedUser.IsActive, expectedUser.Token)

	mock.ExpectQuery("SELECT (.+) FROM users WHERE email = ?").
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
	// Create a new mock database
	db, mock, err := sqlmock.New()
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
		Name:     "Updated User",
		Email:    "updated@example.com",
		Role:     "admin",
		IsActive: true,
		Token:    "updated-token",
	}

	// Set up the expected SQL query and result
	mock.ExpectExec("UPDATE users SET name = \\?, email = \\?, role = \\?, updated_at = \\?, is_active = \\?, token = \\? WHERE id = \\?").
		WithArgs(user.Name, user.Email, user.Role, sqlmock.AnyArg(), user.IsActive, user.Token, user.Id).
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
	// Create a new mock database
	db, mock, err := sqlmock.New()
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
	mock.ExpectExec("UPDATE users SET password = \\?, updated_at = \\? WHERE id = \\?").
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), userID).
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
	// Create a new mock database
	db, mock, err := sqlmock.New()
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
	mock.ExpectExec("DELETE FROM users WHERE id = \\?").
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
	// Create a new mock database
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	// Create a new repository with the mock database
	repo := &UserRepository{
		dbClient: &DatabaseClient{DB: db},
	}

	// Set up test data
	expectedUsers := []*models.User{
		{
			Id:        "user1",
			Name:      "User One",
			Email:     "user1@example.com",
			Password:  "hashed_password1",
			Role:      "user",
			CreatedAt: time.Now().Format(time.RFC3339),
			UpdatedAt: time.Now().Format(time.RFC3339),
			IsActive:  true,
			Token:     "token1",
		},
		{
			Id:        "user2",
			Name:      "User Two",
			Email:     "user2@example.com",
			Password:  "hashed_password2",
			Role:      "admin",
			CreatedAt: time.Now().Format(time.RFC3339),
			UpdatedAt: time.Now().Format(time.RFC3339),
			IsActive:  true,
			Token:     "token2",
		},
	}

	// Set up the expected SQL query and result
	rows := sqlmock.NewRows([]string{"id", "name", "email", "password", "role", "created_at", "updated_at", "is_active", "token"})
	for _, user := range expectedUsers {
		rows.AddRow(user.Id, user.Name, user.Email, user.Password, user.Role, user.CreatedAt, user.UpdatedAt, user.IsActive, user.Token)
	}

	mock.ExpectQuery("SELECT (.+) FROM users ORDER BY created_at DESC").
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
	// Create a new mock database
	db, mock, err := sqlmock.New()
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
	// Generate a real bcrypt hash for "password123"
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		t.Fatalf("Failed to generate password hash: %v", err)
	}
	hashedPassword := string(hashedBytes)

	expectedUser := &models.User{
		Id:        "test-id",
		Name:      "Test User",
		Email:     email,
		Password:  hashedPassword,
		Role:      "user",
		CreatedAt: time.Now().Format(time.RFC3339),
		UpdatedAt: time.Now().Format(time.RFC3339),
		IsActive:  true,
		Token:     "test-token",
	}

	// Set up the expected SQL query and result
	rows := sqlmock.NewRows([]string{"id", "name", "email", "password", "role", "created_at", "updated_at", "is_active", "token"}).
		AddRow(expectedUser.Id, expectedUser.Name, expectedUser.Email, expectedUser.Password, expectedUser.Role, expectedUser.CreatedAt, expectedUser.UpdatedAt, expectedUser.IsActive, expectedUser.Token)

	mock.ExpectQuery("SELECT (.+) FROM users WHERE email = ?").
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

	// Test with inactive user
	db2, mock2, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db2.Close()

	// Create a new repository with the new mock database
	repo2 := &UserRepository{
		dbClient: &DatabaseClient{DB: db2},
	}

	inactiveUser := *expectedUser
	inactiveUser.IsActive = false

	rows2 := sqlmock.NewRows([]string{"id", "name", "email", "password", "role", "created_at", "updated_at", "is_active", "token"}).
		AddRow(inactiveUser.Id, inactiveUser.Name, inactiveUser.Email, inactiveUser.Password, inactiveUser.Role, inactiveUser.CreatedAt, inactiveUser.UpdatedAt, inactiveUser.IsActive, inactiveUser.Token)

	mock2.ExpectQuery("SELECT (.+) FROM users WHERE email = ?").
		WithArgs(email).
		WillReturnRows(rows2)

	// Call the function being tested
	_, err = repo2.VerifyPassword(email, password)

	// Assert that an error occurred
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "user account is inactive")
}

func TestUserRepository_ActivateUser(t *testing.T) {
	// Create a new mock database
	db, mock, err := sqlmock.New()
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
	mock.ExpectExec("UPDATE users SET is_active = true, updated_at = \\? WHERE id = \\?").
		WithArgs(sqlmock.AnyArg(), userID).
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
	// Create a new mock database
	db, mock, err := sqlmock.New()
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
	mock.ExpectExec("UPDATE users SET is_active = false, updated_at = \\? WHERE id = \\?").
		WithArgs(sqlmock.AnyArg(), userID).
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
	// Create a new mock database
	db, mock, err := sqlmock.New()
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
	mock.ExpectQuery("SELECT COUNT\\(\\*\\) FROM users WHERE email = \\?").
		WithArgs(email).
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))

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
	db2, mock2, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db2.Close()

	// Create a new repository with the new mock database
	repo2 := &UserRepository{
		dbClient: &DatabaseClient{DB: db2},
	}

	mock2.ExpectQuery("SELECT COUNT\\(\\*\\) FROM users WHERE email = \\?").
		WithArgs(email).
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(0))

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

func TestUserRepository_UpdateToken(t *testing.T) {
	// Create a new mock database
	db, mock, err := sqlmock.New()
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
	token := "new-token"

	// Set up the expected SQL query and result
	mock.ExpectExec("UPDATE users SET token = \\?, updated_at = \\? WHERE id = \\?").
		WithArgs(token, sqlmock.AnyArg(), userID).
		WillReturnResult(sqlmock.NewResult(0, 1))

	// Call the function being tested
	err = repo.UpdateToken(userID, token)

	// Assert that no errors occurred
	assert.NoError(t, err)

	// Verify that all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}

	// Test when user does not exist
	db2, mock2, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db2.Close()

	// Create a new repository with the new mock database
	repo2 := &UserRepository{
		dbClient: &DatabaseClient{DB: db2},
	}

	mock2.ExpectExec("UPDATE users SET token = \\?, updated_at = \\? WHERE id = \\?").
		WithArgs(token, sqlmock.AnyArg(), userID).
		WillReturnResult(sqlmock.NewResult(0, 0))

	// Call the function being tested
	err = repo2.UpdateToken(userID, token)

	// Assert that an error occurred
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "user not found")

	// Verify that all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestUserRepository_GetByID_NotFound(t *testing.T) {
	// Create a new mock database
	db, mock, err := sqlmock.New()
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

	// Set up the expected SQL query and result
	mock.ExpectQuery("SELECT (.+) FROM users WHERE id = ?").
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
