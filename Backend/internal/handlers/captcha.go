package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"oop/internal/config"
	"time"
)

// TurnstileError represents an error from Cloudflare Turnstile verification
type TurnstileError struct {
	ErrorCodes []string
}

func (e *TurnstileError) Error() string {
	return fmt.Sprintf("turnstile verification failed: %v", e.ErrorCodes)
}

// verifyResp models the JSON returned by Cloudflare
type verifyResp struct {
	Success    bool     `json:"success"`
	ErrorCodes []string `json:"error-codes"`
}

// VerifyTurnstile sends `token` to Cloudflare and returns whether it passed
func VerifyTurnstile(token string) (bool, error) {
	// Load Turnstile configuration
	turnstileConfig, err := config.LoadTurnstileConfig()
	if err != nil {
		return false, err
	}

	form := url.Values{}
	form.Set("secret", turnstileConfig.SecretKey)
	form.Set("response", token)

	// Create a client with a 10-second timeout
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	resp, err := client.PostForm(
		"https://challenges.cloudflare.com/turnstile/v0/siteverify",
		form,
	)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	var result verifyResp
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return false, err
	}

	if !result.Success {
		log.Printf("Turnstile verification failed with error codes: %v", result.ErrorCodes)
		return false, &TurnstileError{ErrorCodes: result.ErrorCodes}
	}

	return result.Success, nil
}
