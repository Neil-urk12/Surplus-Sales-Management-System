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
	secretKey := os.Getenv("TURNSTILE_SECRET_KEY")
	if secretKey == "" {
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
		return TurnstileConfig{}, fmt.Errorf("TURNSTILE_SECRET_KEY appears to be invalid (too short)")
	}

	return TurnstileConfig{
		SecretKey: secretKey,
	}, nil
}
