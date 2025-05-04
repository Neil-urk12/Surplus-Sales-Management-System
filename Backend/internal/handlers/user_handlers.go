package handlers

import (
	"crypto/rand"
	"encoding/base64"
	"log"
	"oop/internal/models"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

// UserRepository defines the interface for user repository operations
type UserRepository interface {
	Create(user *models.User) error
	GetByID(id string) (*models.User, error)
	GetByEmail(email string) (*models.User, error)
	Update(user *models.User) error
	UpdatePassword(userID string, newPassword string) error
	Delete(id string) error
	GetAll() ([]*models.User, error)
	VerifyPassword(email, password string) (*models.User, error)
	ActivateUser(id string) error
	DeactivateUser(id string) error
	EmailExists(email string) (bool, error)
}

// UserHandler handles HTTP requests related to users
type UserHandler struct {
	userRepo  UserRepository
	jwtSecret []byte
}

// NewUserHandler creates a new UserHandler instance
func NewUserHandler(userRepo UserRepository, jwtSecret []byte) *UserHandler {
	return &UserHandler{
		userRepo:  userRepo,
		jwtSecret: jwtSecret,
	}
}

// RegisterRoutes registers only the public user-related routes (register, login)
// Protected routes are registered in main.go with middleware.
func (h *UserHandler) RegisterRoutes(router fiber.Router) {
	// Base path is assumed to be '/api' from main.go
	userGroup := router.Group("/users")

	// Public routes
	userGroup.Post("/register", h.Register)
	userGroup.Post("/login", h.Login)

}

// generateToken creates a secure random token (used for registration maybe)
func generateToken() (string, error) {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}

// Register handles user registration
func (h *UserHandler) Register(c *fiber.Ctx) error {
	var input struct {
		Name     string `json:"fullName"`
		Email    string `json:"email"`
		Password string `json:"password"`
		Role     string `json:"role"`
	}

	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Validate input
	if input.Name == "" || input.Email == "" || input.Password == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Name, email, and password are required",
		})
	}

	// Check if email already exists
	exists, err := h.userRepo.EmailExists(input.Email)
	if err != nil {
		log.Printf("Error checking email existence: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Internal server error",
		})
	}

	if exists {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{
			"error": "Email already in use",
		})
	}

	// Generate a token for the new user
	token, err := generateToken()
	if err != nil {
		log.Printf("Error generating token: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to generate authentication token",
		})
	}

	// Create user
	user := &models.User{
		Id:       uuid.New().String(),
		FullName: input.Name, // Use FullName from input
		Email:    input.Email,
		Password: input.Password, // Will be hashed in the repository
		Role:     input.Role,
	}

	if err := h.userRepo.Create(user); err != nil {
		log.Printf("Error creating user: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create user",
		})
	}

	// Don't return the password hash
	user.Password = ""

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "User registered successfully",
		"user":    user,
		"token":   token,
	})
}

// Login handles user authentication and JWT generation
func (h *UserHandler) Login(c *fiber.Ctx) error {
	var input struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Validate input
	if input.Email == "" || input.Password == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Email and password are required",
		})
	}

	// First, get the user by email to check if they exist and if they're active
	user, err := h.userRepo.GetByEmail(input.Email)
	if err != nil {
		log.Printf("Login failed - user not found for email %s: %v", input.Email, err)
		// Return a generic error message to avoid revealing which part failed (email or password)
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid credentials",
		})
	}

	// Check if user is active before verifying password
	if !user.IsActive {
		log.Printf("Login attempt for inactive user: %s", input.Email)
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "Account is inactive",
		})
	}

	// Verify password
	user, err = h.userRepo.VerifyPassword(input.Email, input.Password)
	if err != nil {
		log.Printf("Login failed - invalid password for email %s: %v", input.Email, err)
		// Return a generic error message to avoid revealing which part failed
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid credentials",
		})
	}

	// Create the claims
	claims := jwt.MapClaims{
		"user_id": user.Id,
		"email":   user.Email,
		"role":    user.Role,
		"exp":     time.Now().Add(time.Hour * 72).Unix(), // Token expires in 72 hours
		"iat":     time.Now().Unix(),                     // Issued at
	}

	// Create token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Generate encoded token string
	tokenString, err := token.SignedString(h.jwtSecret) // Use the injected secret
	if err != nil {
		log.Printf("Error signing JWT token: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to generate authentication token",
		})
	}

	// Optionally update the token in the database (Consider if needed for session invalidation)
	// if err := h.userRepo.UpdateToken(user.Id, tokenString); err != nil {
	// 	log.Printf("Error updating JWT token in DB: %v", err)
	// 	// Decide if this should be a fatal error for login
	// 	return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
	// 		"error": "Failed to update session", // More generic error
	// 	})
	// }

	// Don't return the password hash
	user.Password = ""

	// Return user info and the JWT
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Login successful",
		"user":    user,
		"token":   tokenString, // Return the JWT string
	})
}

// GetAllUsers returns a list of all users
func (h *UserHandler) GetAllUsers(c *fiber.Ctx) error {
	// Access user info from middleware if needed:
	// userID := c.Locals("user_id")
	// userRole := c.Locals("role")
	// log.Printf("GetAllUsers called by user %s with role %s", userID, userRole)

	users, err := h.userRepo.GetAll()
	if err != nil {
		log.Printf("Error getting users: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to retrieve users",
		})
	}

	// Remove password hashes from response
	for _, user := range users {
		user.Password = ""
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"users": users,
	})
}

// GetUser returns a single user by ID
func (h *UserHandler) GetUser(c *fiber.Ctx) error {
	id := c.Params("id")
	// Optional: Check if the requesting user (c.Locals("user_id")) is allowed to view this profile
	// requestedUserID := c.Locals("user_id")
	// requestedUserRole := c.Locals("role")
	// if requestedUserID != id && requestedUserRole != "admin" { // Example policy
	// 	return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "Permission denied"})
	// }

	user, err := h.userRepo.GetByID(id)
	if err != nil {
		log.Printf("Error getting user: %v", err)
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "User not found",
		})
	}

	// Don't return the password hash
	user.Password = ""

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"user": user,
	})
}

// UpdateUser updates a user's information
func (h *UserHandler) UpdateUser(c *fiber.Ctx) error {
	id := c.Params("id")
	// Optional: Check permissions similar to GetUser
	// requestedUserID := c.Locals("user_id")
	// requestedUserRole := c.Locals("role")
	// if requestedUserID != id && requestedUserRole != "admin" {
	// 	return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "Permission denied"})
	// }

	// Get existing user
	existingUser, err := h.userRepo.GetByID(id)
	if err != nil {
		log.Printf("Error getting user: %v", err)
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "User not found",
		})
	}

	// Check if user is active
	if !existingUser.IsActive {
		log.Printf("Update attempt for inactive user: %s", existingUser.Email)
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "Account is inactive",
		})
	}

	// Parse request body
	var input struct {
		Name  string `json:"name"`
		Email string `json:"email"`
		Role  string `json:"role"`
	}

	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Update user fields if provided
	if input.Name != "" {
		existingUser.FullName = input.Name // Update FullName
	}
	if input.Email != "" && input.Email != existingUser.Email {
		// Check if new email already exists
		exists, err := h.userRepo.EmailExists(input.Email)
		if err != nil {
			log.Printf("Error checking email existence: %v", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Internal server error",
			})
		}

		if exists {
			return c.Status(fiber.StatusConflict).JSON(fiber.Map{
				"error": "Email already in use",
			})
		}

		existingUser.Email = input.Email
	}
	if input.Role != "" {
		existingUser.Role = input.Role
	}

	// Update timestamp
	existingUser.UpdatedAt = time.Now() // Use time.Time directly

	// Save changes
	if err := h.userRepo.Update(existingUser); err != nil {
		log.Printf("Error updating user: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update user",
		})
	}

	// Don't return the password hash
	existingUser.Password = ""

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "User updated successfully",
		"user":    existingUser,
	})
}

// DeleteUser removes a user from the system
func (h *UserHandler) DeleteUser(c *fiber.Ctx) error {
	id := c.Params("id")
	// Optional: Check permissions

	if err := h.userRepo.Delete(id); err != nil {
		log.Printf("Error deleting user: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to delete user",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "User deleted successfully",
	})
}

// ActivateUser activates a user account
func (h *UserHandler) ActivateUser(c *fiber.Ctx) error {
	id := c.Params("id")
	// Optional: Check permissions

	if err := h.userRepo.ActivateUser(id); err != nil {
		log.Printf("Error activating user: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to activate user",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "User activated successfully",
	})
}

// DeactivateUser deactivates a user account
func (h *UserHandler) DeactivateUser(c *fiber.Ctx) error {
	id := c.Params("id")
	// Optional: Check permissions

	if err := h.userRepo.DeactivateUser(id); err != nil {
		log.Printf("Error deactivating user: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to deactivate user",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "User deactivated successfully",
	})
}

// UpdatePassword updates a user's password
func (h *UserHandler) UpdatePassword(c *fiber.Ctx) error {
	id := c.Params("id")
	// Optional: Check permissions

	var input struct {
		CurrentPassword string `json:"current_password"`
		NewPassword     string `json:"new_password"`
	}

	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Validate input
	if input.CurrentPassword == "" || input.NewPassword == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Current password and new password are required",
		})
	}

	// Get the user
	user, err := h.userRepo.GetByID(id)
	if err != nil {
		log.Printf("Error getting user: %v", err)
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "User not found",
		})
	}

	// Check if user is active
	if !user.IsActive {
		log.Printf("Password update attempt for inactive user: %s", user.Email)
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "Account is inactive",
		})
	}

	// Verify current password
	_, err = h.userRepo.VerifyPassword(user.Email, input.CurrentPassword)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Current password is incorrect",
		})
	}

	// Update password
	if err := h.userRepo.UpdatePassword(id, input.NewPassword); err != nil {
		log.Printf("Error updating password: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update password",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Password updated successfully",
	})
}
