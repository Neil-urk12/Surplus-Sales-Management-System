package models

import "time"

// UserCreateRequest defines the shape of the request body for user registration and admin/staff creation.
// It's used when a new user is being added to the system.
type UserCreateRequest struct {
	FullName string `json:"fullName" example:"John Doe"`
	Email    string `json:"email" example:"john.doe@example.com"`
	Password string `json:"password" example:"securepassword123"`
	Role     string `json:"role" example:"staff" enums:"staff,admin"`
}

// UserLoginRequest defines the shape of the request body for user login.
// It contains the credentials required for a user to authenticate.
type UserLoginRequest struct {
	Email    string `json:"email" example:"john.doe@example.com"`
	Password string `json:"password" example:"securepassword123"`
}

// UserResponse defines the shape of user data returned by the API.
// This model is used for responses to avoid exposing sensitive information like password hashes.
type UserResponse struct {
	Id        string    `json:"id" example:"xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx"`
	FullName  string    `json:"fullName" example:"John Doe"`
	Email     string    `json:"email" example:"john.doe@example.com"`
	Role      string    `json:"role" example:"staff"`
	IsActive  bool      `json:"isActive" example:"true"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

// UserUpdateRequest defines the shape of the request body for updating user information.
// Fields are optional (omitempty), and a pointer is used for boolean 'IsActive'
// to differentiate between explicitly setting it to false and not providing the field.
type UserUpdateRequest struct {
	FullName string `json:"fullName,omitempty" example:"Johnathan Doe"`
	Email    string `json:"email,omitempty" example:"johnathan.doe@example.com"`
	Role     string `json:"role,omitempty" example:"admin" enums:"staff,admin"`
	IsActive *bool  `json:"isActive,omitempty"` // Using pointer to distinguish between false and not provided
}

// UserPasswordUpdateRequest defines the shape of the request body for updating a user's password.
// It typically requires the new password. The current password might be handled by user verification
// if needed by the endpoint, but not usually part of this specific request model if ID is path param.
type UserPasswordUpdateRequest struct {
	NewPassword string `json:"newPassword" example:"newsecurepassword123"`
}
