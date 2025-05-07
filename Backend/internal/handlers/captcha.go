package handlers

import (
	"encoding/json"
	"net/http"
	"net/url"
	"oop/internal/config"
	"time"
)

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
	return result.Success, nil
}
