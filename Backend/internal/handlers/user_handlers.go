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

const (
	RoleAdmin = "admin"
	RoleStaff = "staff"
)

// UserAuthResponse is the response for successful user registration or login.
type UserAuthResponse struct {
	Message string       `json:"message"`
	User    *models.User `json:"user"`
	Token   string       `json:"token"`
}

// UserListResponse is the response for listing multiple users.
type UserListResponse struct {
	Users []*models.User `json:"users"`
}

// SingleUserResponse is the response for fetching a single user.
type SingleUserResponse struct {
	User *models.User `json:"user"`
}

// UserActionResponse is for actions like create or update that return a user and a message.
type UserActionResponse struct {
	Message string       `json:"message"`
	User    *models.User `json:"user"`
}

// MessageResponse is a generic response for actions that only return a message.
type MessageResponse struct {
	Message string `json:"message"`
}

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
// @Summary Register a new user
// @Description Creates a new user account.
// @Tags Users
// @Accept json
// @Produce json
// @Param user body models.UserCreateRequest true "User Registration Information"
// @Success 201 {object} UserAuthResponse "User registered successfully"
// @Failure 400 {object} ErrorResponse "Invalid request body or missing fields"
// @Failure 409 {object} ErrorResponse "Email already in use"
// @Failure 500 {object} ErrorResponse "Internal server error"
// @Router /users/register [post]
func (h *UserHandler) Register(c *fiber.Ctx) error {
	var input struct {
		Name     string `json:"fullName"`
		Email    string `json:"email"`
		Password string `json:"password"`
		Role     string `json:"role"`
	}

	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Error:      "Invalid request body",
			StatusCode: fiber.StatusBadRequest,
		})
	}

	// Validate input
	if input.Name == "" || input.Email == "" || input.Password == "" {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Error:      "Name, email, and password are required",
			StatusCode: fiber.StatusBadRequest,
		})
	}

	// Check if email already exists
	exists, err := h.userRepo.EmailExists(input.Email)
	if err != nil {
		log.Printf("Error checking email existence: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{
			Error:      "Internal server error",
			StatusCode: fiber.StatusInternalServerError,
		})
	}

	if exists {
		return c.Status(fiber.StatusConflict).JSON(ErrorResponse{
			Error:      "Email already in use",
			StatusCode: fiber.StatusConflict,
		})
	}

	// Generate a token for the new user
	token, err := generateToken()
	if err != nil {
		log.Printf("Error generating token: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{
			Error:      "Failed to generate authentication token",
			StatusCode: fiber.StatusInternalServerError,
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
		return c.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{
			Error:      "Failed to create user",
			StatusCode: fiber.StatusInternalServerError,
		})
	}

	// Don't return the password hash
	user.Password = ""

	return c.Status(fiber.StatusCreated).JSON(UserAuthResponse{
		Message: "User registered successfully",
		User:    user,
		Token:   token,
	})
}

// Login handles user authentication and JWT generation
// @Summary Log in an existing user
// @Description Authenticates a user and returns a JWT token.
// @Tags Users
// @Accept json
// @Produce json
// @Param credentials body models.UserLoginRequest true "User Login Credentials"
// @Success 200 {object} UserAuthResponse "Login successful"
// @Failure 400 {object} ErrorResponse "Invalid request body or missing fields"
// @Failure 401 {object} ErrorResponse "Invalid credentials"
// @Failure 403 {object} ErrorResponse "Account is inactive"
// @Failure 500 {object} ErrorResponse "Internal server error"
// @Router /users/login [post]
func (h *UserHandler) Login(c *fiber.Ctx) error {
	var input struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Error:      "Invalid request body",
			StatusCode: fiber.StatusBadRequest,
		})
	}

	// Validate input
	if input.Email == "" || input.Password == "" {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Error:      "Email and password are required",
			StatusCode: fiber.StatusBadRequest,
		})
	}

	// First, get the user by email to check if they exist and if they're active
	user, err := h.userRepo.GetByEmail(input.Email)
	if err != nil {
		log.Printf("Login failed - user not found for email %s: %v", input.Email, err)
		// Return a generic error message to avoid revealing which part failed (email or password)
		return c.Status(fiber.StatusUnauthorized).JSON(ErrorResponse{
			Error:      "Invalid credentials",
			StatusCode: fiber.StatusUnauthorized,
		})
	}

	// Check if user is active before verifying password
	if !user.IsActive {
		log.Printf("Login attempt for inactive user: %s", input.Email)
		return c.Status(fiber.StatusForbidden).JSON(ErrorResponse{
			Error:      "Account is inactive",
			StatusCode: fiber.StatusForbidden,
		})
	}

	// Verify password
	user, err = h.userRepo.VerifyPassword(input.Email, input.Password)
	if err != nil {
		log.Printf("Login failed - invalid password for email %s: %v", input.Email, err)
		// Return a generic error message to avoid revealing which part failed
		return c.Status(fiber.StatusUnauthorized).JSON(ErrorResponse{
			Error:      "Invalid credentials",
			StatusCode: fiber.StatusUnauthorized,
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
		return c.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{
			Error:      "Failed to generate authentication token",
			StatusCode: fiber.StatusInternalServerError,
		})
	}

	// Optionally update the token in the database (Consider if needed for session invalidation)
	// if err := h.userRepo.UpdateToken(user.Id, tokenString); err != nil {
	// 	log.Printf("Error updating JWT token in DB: %v", err)
	// 	// Decide if this should be a fatal error for login
	// 	return c.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{Error: "Failed to update session", StatusCode: fiber.StatusInternalServerError})
	// }

	// Don't return the password hash
	user.Password = ""

	// Return user info and the JWT
	return c.Status(fiber.StatusOK).JSON(UserAuthResponse{
		Message: "Login successful",
		User:    user,
		Token:   tokenString,
	})
}

// GetAllUsers returns a list of all users
// @Summary Get all users
// @Description Retrieves a list of all registered users. This is a protected route.
// @Tags Users
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} UserListResponse "Successfully retrieved list of users"
// @Failure 500 {object} ErrorResponse "Failed to retrieve users"
// @Router /users [get]
func (h *UserHandler) GetAllUsers(c *fiber.Ctx) error {
	// Access user info from middleware if needed:
	// userID := c.Locals("user_id")
	// userRole := c.Locals("role")
	// log.Printf("GetAllUsers called by user %s with role %s", userID, userRole)

	users, err := h.userRepo.GetAll()
	if err != nil {
		log.Printf("Error getting users: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{
			Error:      "Failed to retrieve users",
			StatusCode: fiber.StatusInternalServerError,
		})
	}

	// Remove password hashes from response
	for _, user := range users {
		user.Password = ""
	}

	return c.Status(fiber.StatusOK).JSON(UserListResponse{
		Users: users,
	})
}

// GetUser returns a single user by ID
// @Summary Get user by ID
// @Description Retrieves a single user by their ID. This is a protected route.
// @Tags Users
// @Produce json
// @Security ApiKeyAuth
// @Param id path string true "User ID"
// @Success 200 {object} SingleUserResponse "Successfully retrieved user"
// @Failure 403 {object} ErrorResponse "Permission denied"
// @Failure 404 {object} ErrorResponse "User not found"
// @Router /users/{id} [get]
func (h *UserHandler) GetUser(c *fiber.Ctx) error {
	id := c.Params("id")
	// Optional: Check if the requesting user (c.Locals("user_id")) is allowed to view this profile
	// requestedUserID := c.Locals("user_id")
	// requestedUserRole := c.Locals("role")
	// if requestedUserID != id && requestedUserRole != "admin" { // Example policy
	// 	return c.Status(fiber.StatusForbidden).JSON(ErrorResponse{Error: "Permission denied", StatusCode: fiber.StatusForbidden})
	// }

	user, err := h.userRepo.GetByID(id)
	if err != nil {
		log.Printf("Error getting user: %v", err)
		return c.Status(fiber.StatusNotFound).JSON(ErrorResponse{
			Error:      "User not found",
			StatusCode: fiber.StatusNotFound,
		})
	}

	// Don't return the password hash
	user.Password = ""

	return c.Status(fiber.StatusOK).JSON(SingleUserResponse{
		User: user,
	})
}

// UpdateUser updates a user's information
// @Summary Update user information
// @Description Updates a user's full name, email, role, or active status. This is a protected route (Admin/Staff only).
// @Tags Users
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path string true "User ID"
// @Param user_update body models.UserUpdateRequest true "User Update Information"
// @Success 200 {object} UserActionResponse "User updated successfully"
// @Failure 400 {object} ErrorResponse "Invalid request body"
// @Failure 403 {object} ErrorResponse "Permission denied"
// @Failure 404 {object} ErrorResponse "User not found"
// @Failure 500 {object} ErrorResponse "Failed to update user"
// @Router /users/{id} [put]
func (h *UserHandler) UpdateUser(c *fiber.Ctx) error {
	id := c.Params("id")
	// Check if the current user is admin or staff
	roleValue := c.Locals("role")
	requestUserRole, ok := roleValue.(string)
	if !ok || requestUserRole == "" {
		// Log the issue for debugging
		log.Printf("UpdateUser: 'role' not found or not a string in context locals for user ID %s. Value: %v", id, roleValue)
		// Return forbidden, as the role couldn't be determined or is invalid
		return c.Status(fiber.StatusForbidden).JSON(ErrorResponse{Error: "Permission denied: Unable to verify user role", StatusCode: fiber.StatusForbidden})
	}
	if requestUserRole != RoleAdmin && requestUserRole != RoleStaff { // Now check the validated role
		return c.Status(fiber.StatusForbidden).JSON(ErrorResponse{Error: "Permission denied", StatusCode: fiber.StatusForbidden})
	}

	// Get existing user
	existingUser, err := h.userRepo.GetByID(id)
	if err != nil {
		log.Printf("Error getting user: %v", err)
		return c.Status(fiber.StatusNotFound).JSON(ErrorResponse{
			Error:      "User not found",
			StatusCode: fiber.StatusNotFound,
		})
	}

	// Parse request body
	var input struct {
		FullName string `json:"fullName"`
		Email    string `json:"email"`
		Role     string `json:"role"`
		IsActive bool   `json:"isActive"`
	}

	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Error:      "Invalid request body",
			StatusCode: fiber.StatusBadRequest,
		})
	}
	log.Printf("UpdateUser received input: %+v", input) // Add logging here

	// Update fields if provided in the request body
	if input.FullName != "" {
		existingUser.FullName = input.FullName
	}
	if input.Email != "" {
		existingUser.Email = input.Email
	}
	if input.Role != "" {
		existingUser.Role = input.Role
	}

	// Update isActive status
	existingUser.IsActive = input.IsActive

	// Update timestamp
	existingUser.UpdatedAt = time.Now()

	// Save changes
	if err := h.userRepo.Update(existingUser); err != nil {
		log.Printf("Error updating user: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{
			Error:      "Failed to update user",
			StatusCode: fiber.StatusInternalServerError,
		})
	}

	// Don't return the password hash
	existingUser.Password = ""

	return c.Status(fiber.StatusOK).JSON(UserActionResponse{
		Message: "User updated successfully",
		User:    existingUser,
	})
}

// DeleteUser removes a user from the system
// @Summary Delete user
// @Description Deletes a user by their ID. This is a protected route.
// @Tags Users
// @Produce json
// @Security ApiKeyAuth
// @Param id path string true "User ID"
// @Success 200 {object} MessageResponse "User deleted successfully"
// @Failure 500 {object} ErrorResponse "Failed to delete user"
// @Router /users/{id} [delete]
func (h *UserHandler) DeleteUser(c *fiber.Ctx) error {
	id := c.Params("id")
	// Optional: Check permissions

	if err := h.userRepo.Delete(id); err != nil {
		log.Printf("Error deleting user: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{
			Error:      "Failed to delete user",
			StatusCode: fiber.StatusInternalServerError,
		})
	}

	return c.Status(fiber.StatusOK).JSON(MessageResponse{
		Message: "User deleted successfully",
	})
}

// ActivateUser activates a user account
// @Summary Activate user account
// @Description Activates a previously deactivated user account. This is a protected route.
// @Tags Users
// @Produce json
// @Security ApiKeyAuth
// @Param id path string true "User ID"
// @Success 200 {object} MessageResponse "User activated successfully"
// @Failure 500 {object} ErrorResponse "Failed to activate user"
// @Router /users/{id}/activate [put]
func (h *UserHandler) ActivateUser(c *fiber.Ctx) error {
	id := c.Params("id")

	if err := h.userRepo.ActivateUser(id); err != nil {
		log.Printf("Error activating user: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{
			Error:      "Failed to activate user",
			StatusCode: fiber.StatusInternalServerError,
		})
	}

	return c.Status(fiber.StatusOK).JSON(MessageResponse{
		Message: "User activated successfully",
	})
}

// CreateUser creates a new user (admin/staff only)
// @Summary Create user (Admin/Staff)
// @Description Allows Admin or Staff to create a new user account. This is a protected route.
// @Tags Users
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param user body models.UserCreateRequest true "User Creation Information"
// @Success 201 {object} UserActionResponse "User created successfully"
// @Failure 400 {object} ErrorResponse "Invalid request body or missing fields"
// @Failure 403 {object} ErrorResponse "Permission denied"
// @Failure 409 {object} ErrorResponse "Email already in use"
// @Failure 500 {object} ErrorResponse "Internal server error or failed to create user"
// @Router /users [post] // Note: This matches the route in main.go for creating users by admin/staff
func (h *UserHandler) CreateUser(c *fiber.Ctx) error {
	roleValue := c.Locals("role")
	requestUserRole, ok := roleValue.(string)

	if !ok || requestUserRole == "" {
		log.Printf("CreateUser: 'role' not found or not a string in context locals. Value: %v", roleValue)
		return c.Status(fiber.StatusForbidden).JSON(ErrorResponse{Error: "Permission denied: Unable to verify user role", StatusCode: fiber.StatusForbidden})
	}
	if requestUserRole != RoleAdmin && requestUserRole != RoleStaff {
		return c.Status(fiber.StatusForbidden).JSON(ErrorResponse{Error: "Permission denied", StatusCode: fiber.StatusForbidden})
	}

	// Parse request body
	var input struct {
		FullName string `json:"fullName"`
		Email    string `json:"email"`
		Password string `json:"password"`
		Role     string `json:"role"`
	}

	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Error:      "Invalid request body",
			StatusCode: fiber.StatusBadRequest,
		})
	}

	// Validate input
	if input.FullName == "" || input.Email == "" || input.Password == "" || input.Role == "" {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Error:      "Full name, email, password, and role are required",
			StatusCode: fiber.StatusBadRequest,
		})
	}

	// Check if email already exists
	exists, err := h.userRepo.EmailExists(input.Email)
	if err != nil {
		log.Printf("Error checking email existence: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{
			Error:      "Internal server error",
			StatusCode: fiber.StatusInternalServerError,
		})
	}

	if exists {
		return c.Status(fiber.StatusConflict).JSON(ErrorResponse{
			Error:      "Email already in use",
			StatusCode: fiber.StatusConflict,
		})
	}

	// Create user
	user := &models.User{
		Id:       uuid.New().String(),
		FullName: input.FullName,
		Email:    input.Email,
		Password: input.Password, // Will be hashed in the repository
		Role:     input.Role,
		IsActive: true,
	}

	if err := h.userRepo.Create(user); err != nil {
		log.Printf("Error creating user: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{
			Error:      "Failed to create user",
			StatusCode: fiber.StatusInternalServerError,
		})
	}

	// Don't return the password hash
	user.Password = ""

	return c.Status(fiber.StatusCreated).JSON(UserActionResponse{
		Message: "User created successfully",
		User:    user,
	})
}

// DeactivateUser deactivates a user account
// @Summary Deactivate user account
// @Description Deactivates an active user account. This is a protected route.
// @Tags Users
// @Produce json
// @Security ApiKeyAuth
// @Param id path string true "User ID"
// @Success 200 {object} MessageResponse "User deactivated successfully"
// @Failure 500 {object} ErrorResponse "Failed to deactivate user"
// @Router /users/{id}/deactivate [put]
func (h *UserHandler) DeactivateUser(c *fiber.Ctx) error {
	id := c.Params("id")
	// Optional: Check permissions

	if err := h.userRepo.DeactivateUser(id); err != nil {
		log.Printf("Error deactivating user: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{
			Error:      "Failed to deactivate user",
			StatusCode: fiber.StatusInternalServerError,
		})
	}

	return c.Status(fiber.StatusOK).JSON(MessageResponse{
		Message: "User deactivated successfully",
	})
}

// UpdatePassword updates a user's password
// @Summary Update user password
// @Description Updates the password for a given user. This is a protected route.
// @Tags Users
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path string true "User ID"
// @Param password_update body models.UserPasswordUpdateRequest true "Password Update Information"
// @Success 200 {object} MessageResponse "Password updated successfully"
// @Failure 400 {object} ErrorResponse "Invalid request body or missing fields"
// @Failure 403 {object} ErrorResponse "Permission denied"
// @Failure 404 {object} ErrorResponse "User not found"
// @Failure 500 {object} ErrorResponse "Failed to update password"
// @Router /users/{id}/password [put]
func (h *UserHandler) UpdatePassword(c *fiber.Ctx) error {
	id := c.Params("id")
	// Check if the current user is admin or the user themselves
	requestUserID := c.Locals("user_id").(string)
	requestUserRole := c.Locals("role").(string)

	if requestUserID != id && requestUserRole != RoleAdmin {
		return c.Status(fiber.StatusForbidden).JSON(ErrorResponse{Error: "Permission denied", StatusCode: fiber.StatusForbidden})
	}

	var input struct {
		CurrentPassword string `json:"current_password"`
		NewPassword     string `json:"new_password"`
	}

	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Error:      "Invalid request body",
			StatusCode: fiber.StatusBadRequest,
		})
	}

	if input.NewPassword == "" {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Error:      "New password is required",
			StatusCode: fiber.StatusBadRequest,
		})
	}

	// Check if user exists before trying to update password
	user, err := h.userRepo.GetByID(id)
	if err != nil {
		log.Printf("UpdatePassword: User not found with ID %s: %v", id, err)
		return c.Status(fiber.StatusNotFound).JSON(ErrorResponse{
			Error:      "User not found",
			StatusCode: fiber.StatusNotFound,
		})
	}

	// Verify current password
	_, err = h.userRepo.VerifyPassword(user.Email, input.CurrentPassword)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(ErrorResponse{
			Error:      "Current password is incorrect",
			StatusCode: fiber.StatusUnauthorized,
		})
	}

	// Update password
	if err := h.userRepo.UpdatePassword(id, input.NewPassword); err != nil {
		log.Printf("Error updating password for user %s: %v", id, err)
		return c.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{
			Error:      "Failed to update password",
			StatusCode: fiber.StatusInternalServerError,
		})
	}

	return c.Status(fiber.StatusOK).JSON(MessageResponse{
		Message: "Password updated successfully",
	})
}
