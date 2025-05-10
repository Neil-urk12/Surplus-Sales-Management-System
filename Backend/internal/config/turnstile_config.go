package config

import (
	"fmt"
	"os"
	"strings"
)

// TurnstileConfig holds the configuration for Cloudflare Turnstile
type TurnstileConfig struct {
	SecretKey string
}

// LoadTurnstileConfig loads and validates the Turnstile configuration
func LoadTurnstileConfig() (TurnstileConfig, error) {
	// Check if environment variable exists
	secretKey, exists := os.LookupEnv("TURNSTILE_SECRET_KEY")
	if !exists {
		return TurnstileConfig{}, fmt.Errorf("TURNSTILE_SECRET_KEY environment variable is not set")
	}

	// Remove any whitespace
	secretKey = strings.TrimSpace(secretKey)
	if secretKey == "" {
		return TurnstileConfig{}, fmt.Errorf("TURNSTILE_SECRET_KEY environment variable is empty")
	}

	// Cloudflare Turnstile secret keys are typically 32 characters long
	// This is a basic validation - adjust if Cloudflare has different requirements
	if len(secretKey) < 32 {
		return TurnstileConfig{}, fmt.Errorf("TURNSTILE_SECRET_KEY appears to be invalid (too short, expected at least 32 characters)")
	}

	// Additional validation: Check for common invalid characters
	invalidChars := []string{" ", "\t", "\n", "\r"}
	for _, char := range invalidChars {
		if strings.Contains(secretKey, char) {
			return TurnstileConfig{}, fmt.Errorf("TURNSTILE_SECRET_KEY contains invalid characters (whitespace or newlines)")
		}
	}

	// Additional validation: Check if the key appears to be a valid format
	// Cloudflare Turnstile keys typically contain only alphanumeric characters and some special characters
	validChars := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789-_"
	for _, char := range secretKey {
		if !strings.ContainsRune(validChars, char) {
			return TurnstileConfig{}, fmt.Errorf("TURNSTILE_SECRET_KEY contains invalid characters (only alphanumeric, hyphen, and underscore are allowed)")
		}
	}

	return TurnstileConfig{
		SecretKey: secretKey,
	}, nil
}
