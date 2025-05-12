// Package middleware provides request middleware functions.
package middleware

import (
	"errors"
	"log"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
)

// JWTMiddleware creates a new Fiber middleware handler for JWT authentication.
// It expects the JWT secret key as a byte slice.
func JWTMiddleware(secret []byte) fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Missing or malformed JWT",
			})
		}

		// Check for "Bearer " prefix
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Missing or malformed JWT (Bearer token required)",
			})
		}
		tokenString := parts[1]

		// Parse and validate the token
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// Validate the alg is what you expect:
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fiber.NewError(fiber.StatusUnauthorized, "unexpected signing method")
			}
			return secret, nil
		})

		if err != nil {
			log.Printf("JWT Error: %v", err)
			// Check specifically for expiration error using errors.Is
			if errors.Is(err, jwt.ErrTokenExpired) {
				return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Token has expired"})
			}
			// Handle other validation errors or parsing errors
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid or malformed JWT", // General error for other issues
			})
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			// Check expiration (redundant with Parse validation but explicit)
			if exp, ok := claims["exp"].(float64); ok {
				if time.Unix(int64(exp), 0).Before(time.Now()) {
					return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Token has expired"})
				}
			} else {
				return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid token claims (exp)"})
			}

			// Store user info in locals for downstream handlers
			c.Locals("user_id", claims["user_id"])
			c.Locals("email", claims["email"])
			c.Locals("role", claims["role"])
			return c.Next()
		}

		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid JWT",
		})
	}
}
