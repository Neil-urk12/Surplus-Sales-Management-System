package handlers

import (
	"encoding/json"
	"net/http"
	"net/url"
	"os"
)

// verifyResp models the JSON returned by Cloudflare
type verifyResp struct {
	Success    bool     `json:"success"`
	ErrorCodes []string `json:"error-codes"`
}

// VerifyTurnstile sends `token` to Cloudflare and returns whether it passed
func VerifyTurnstile(token string) (bool, error) {
	form := url.Values{}
	form.Set("secret", os.Getenv("TURNSTSTILE_SECRET_KEY"))
	form.Set("response", token)

	resp, err := http.PostForm(
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
