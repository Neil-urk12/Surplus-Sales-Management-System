package repositories

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"oop/internal/models"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

// UserRepository handles database operations related to users
type UserRepository struct {
	dbClient *DatabaseClient
}

// NewUserRepository creates a new UserRepository instance
func NewUserRepository(dbClient *DatabaseClient) *UserRepository {
	return &UserRepository{
		dbClient: dbClient,
	}
}

// Create adds a new user to the database
func (r *UserRepository) Create(user *models.User) error {
	// Generate a new UUID for the user if not provided
	if user.Id == "" {
		user.Id = uuid.New().String()
	}

	// Hash the password before storing
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}

	// Set timestamps
	now := time.Now()
	user.CreatedAt = now
	user.UpdatedAt = now

	// Set default role if not provided
	if user.Role == "" {
		user.Role = "user"
	}

	// Set active by default
	user.IsActive = true

	// Insert the user into the database
	query := `
		INSERT INTO users (id, full_name, email, password_hash, role, created_at, updated_at, is_active)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)
	`
	_, err = r.dbClient.DB.Exec(
		query,
		user.Id,
		user.FullName,
		user.Email,
		string(hashedPassword),
		user.Role,
		user.CreatedAt,
		user.UpdatedAt,
		user.IsActive,
	)

	if err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}

	return nil
}

// GetByID retrieves a user by their ID
func (r *UserRepository) GetByID(id string) (*models.User, error) {
	query := `
		SELECT id, full_name, email, password, role, created_at, updated_at, is_active
		FROM users
		WHERE id = ?
	`
	row := r.dbClient.DB.QueryRow(query, id)

	var user models.User
	err := row.Scan(
		&user.Id,
		&user.FullName,
		&user.Email,
		&user.Password,
		&user.Role,
		&user.CreatedAt,
		&user.UpdatedAt,
		&user.IsActive,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("user not found: %w", err)
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return &user, nil
}

// GetByEmail retrieves a user by their email address
func (r *UserRepository) GetByEmail(email string) (*models.User, error) {
	query := `
		SELECT id, full_name, email, password_hash, role, created_at, updated_at, is_active
		FROM users
		WHERE email = ?
	`
	row := r.dbClient.DB.QueryRow(query, email)

	var user models.User
	err := row.Scan(
		&user.Id,
		&user.FullName,
		&user.Email,
		&user.Password,
		&user.Role,
		&user.CreatedAt,
		&user.UpdatedAt,
		&user.IsActive,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("user not found: %w", err)
		}
		return nil, fmt.Errorf("failed to get user by email: %w", err)
	}

	return &user, nil
}

// Update updates an existing user in the database
func (r *UserRepository) Update(user *models.User) error {
	// Update the timestamp
	user.UpdatedAt = time.Now()

	query := `
		UPDATE users
		SET full_name = ?, email = ?, role = ?, updated_at = ?, is_active = ?
		WHERE id = ?
	`
	result, err := r.dbClient.DB.Exec(
		query,
		user.FullName,
		user.Email,
		user.Role,
		user.UpdatedAt,
		user.IsActive,
		user.Id,
	)

	if err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("user not found")
	}

	return nil
}

// UpdatePassword updates a user's password
func (r *UserRepository) UpdatePassword(userID string, newPassword string) error {
	// Hash the new password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}

	// Update the timestamp
	updatedAt := time.Now()

	query := `
		UPDATE users
		SET password = ?, updated_at = ?
		WHERE id = ?
	`
	result, err := r.dbClient.DB.Exec(
		query,
		string(hashedPassword),
		updatedAt,
		userID,
	)

	if err != nil {
		return fmt.Errorf("failed to update password: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("user not found")
	}

	return nil
}

// Delete removes a user from the database
func (r *UserRepository) Delete(id string) error {
	query := `DELETE FROM users WHERE id = ?`
	result, err := r.dbClient.DB.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("user not found")
	}

	return nil
}

// GetAll retrieves all users from the database
func (r *UserRepository) GetAll() ([]*models.User, error) {
	query := `
		SELECT id, full_name, email, password, role, created_at, updated_at, is_active
		FROM users
		ORDER BY created_at DESC
	`
	rows, err := r.dbClient.DB.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to query users: %w", err)
	}
	defer rows.Close()

	var users []*models.User
	for rows.Next() {
		var user models.User
		err := rows.Scan(
			&user.Id,
			&user.FullName,
			&user.Email,
			&user.Password,
			&user.Role,
			&user.CreatedAt,
			&user.UpdatedAt,
			&user.IsActive,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan user row: %w", err)
		}
		users = append(users, &user)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating user rows: %w", err)
	}

	return users, nil
}

// VerifyPassword checks if the provided password matches the stored hash
func (r *UserRepository) VerifyPassword(email, password string) (*models.User, error) {
	user, err := r.GetByEmail(email)
	if err != nil {
		return nil, err
	}

	// Check if the user is active
	if !user.IsActive {
		return nil, fmt.Errorf("user account is inactive")
	}

	// Compare the provided password with the stored hash
	log.Println(user.Password)
	log.Println(password)
	log.Println(bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost))
	log.Println(bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)))
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return nil, fmt.Errorf("invalid password")
	}

	return user, nil
}

// ActivateUser sets a user's status to active
func (r *UserRepository) ActivateUser(id string) error {
	query := `UPDATE users SET is_active = ?, updated_at = ? WHERE id = ?`
	_, err := r.dbClient.DB.Exec(query, true, time.Now(), id)
	if err != nil {
		return fmt.Errorf("failed to activate user: %w", err)
	}

	return nil
}

// DeactivateUser sets a user's status to inactive
func (r *UserRepository) DeactivateUser(id string) error {
	query := `UPDATE users SET is_active = ?, updated_at = ? WHERE id = ?`
	_, err := r.dbClient.DB.Exec(query, false, time.Now(), id)
	if err != nil {
		return fmt.Errorf("failed to deactivate user: %w", err)
	}

	return nil
}

// EmailExists checks if an email already exists in the database
func (r *UserRepository) EmailExists(email string) (bool, error) {
	var exists bool
	query := `SELECT EXISTS(SELECT 1 FROM users WHERE email = ?)`
	err := r.dbClient.DB.QueryRow(query, email).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("failed to check if email exists: %w", err)
	}
	return exists, nil
}
