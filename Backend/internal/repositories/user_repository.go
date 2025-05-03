package repositories

import (
	"database/sql"
	"errors"
	"fmt"
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
	now := time.Now().Format(time.RFC3339)
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
		INSERT INTO users (id, name, email, password, role, created_at, updated_at, is_active, token)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
	`
	_, err = r.dbClient.DB.Exec(
		query,
		user.Id,
		user.Name,
		user.Email,
		string(hashedPassword),
		user.Role,
		user.CreatedAt,
		user.UpdatedAt,
		user.IsActive,
		user.Token,
	)

	if err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}

	return nil
}

// GetByID retrieves a user by their ID
func (r *UserRepository) GetByID(id string) (*models.User, error) {
	query := `
		SELECT id, name, email, password, role, created_at, updated_at, is_active, token
		FROM users
		WHERE id = ?
	`
	row := r.dbClient.DB.QueryRow(query, id)

	var user models.User
	err := row.Scan(
		&user.Id,
		&user.Name,
		&user.Email,
		&user.Password,
		&user.Role,
		&user.CreatedAt,
		&user.UpdatedAt,
		&user.IsActive,
		&user.Token,
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
		SELECT id, name, email, password, role, created_at, updated_at, is_active, token
		FROM users
		WHERE email = ?
	`
	row := r.dbClient.DB.QueryRow(query, email)

	var user models.User
	err := row.Scan(
		&user.Id,
		&user.Name,
		&user.Email,
		&user.Password,
		&user.Role,
		&user.CreatedAt,
		&user.UpdatedAt,
		&user.IsActive,
		&user.Token,
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
	user.UpdatedAt = time.Now().Format(time.RFC3339)

	query := `
		UPDATE users
		SET name = ?, email = ?, role = ?, updated_at = ?, is_active = ?, token = ?
		WHERE id = ?
	`
	result, err := r.dbClient.DB.Exec(
		query,
		user.Name,
		user.Email,
		user.Role,
		user.UpdatedAt,
		user.IsActive,
		user.Token,
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
	updatedAt := time.Now().Format(time.RFC3339)

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
		SELECT id, name, email, password, role, created_at, updated_at, is_active, token
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
			&user.Name,
			&user.Email,
			&user.Password,
			&user.Role,
			&user.CreatedAt,
			&user.UpdatedAt,
			&user.IsActive,
			&user.Token,
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
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return nil, fmt.Errorf("invalid password")
	}

	return user, nil
}

// ActivateUser sets a user's status to active
func (r *UserRepository) ActivateUser(id string) error {
	query := `
		UPDATE users
		SET is_active = true, updated_at = ?
		WHERE id = ?
	`
	updatedAt := time.Now().Format(time.RFC3339)
	result, err := r.dbClient.DB.Exec(query, updatedAt, id)
	if err != nil {
		return fmt.Errorf("failed to activate user: %w", err)
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

// DeactivateUser sets a user's status to inactive
func (r *UserRepository) DeactivateUser(id string) error {
	query := `
		UPDATE users
		SET is_active = false, updated_at = ?
		WHERE id = ?
	`
	updatedAt := time.Now().Format(time.RFC3339)
	result, err := r.dbClient.DB.Exec(query, updatedAt, id)
	if err != nil {
		return fmt.Errorf("failed to deactivate user: %w", err)
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

// EmailExists checks if an email already exists in the database
func (r *UserRepository) EmailExists(email string) (bool, error) {
	query := `SELECT COUNT(*) FROM users WHERE email = ?`
	var count int
	err := r.dbClient.DB.QueryRow(query, email).Scan(&count)
	if err != nil {
		return false, fmt.Errorf("failed to check email existence: %w", err)
	}
	return count > 0, nil
}

// UpdateToken updates a user's token
func (r *UserRepository) UpdateToken(userID string, token string) error {
	// Update the timestamp
	updatedAt := time.Now().Format(time.RFC3339)

	query := `
		UPDATE users
		SET token = ?, updated_at = ?
		WHERE id = ?
	`
	result, err := r.dbClient.DB.Exec(
		query,
		token,
		updatedAt,
		userID,
	)

	if err != nil {
		return fmt.Errorf("failed to update token: %w", err)
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
