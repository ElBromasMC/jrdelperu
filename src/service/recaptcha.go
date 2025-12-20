package service

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"time"
)

var (
	ErrRecaptchaFailed   = errors.New("reCAPTCHA verification failed")
	ErrRecaptchaNotFound = errors.New("reCAPTCHA token not provided")
	ErrRecaptchaLowScore = errors.New("reCAPTCHA score too low")
)

const (
	recaptchaVerifyURL = "https://www.google.com/recaptcha/api/siteverify"
	requestTimeout     = 10 * time.Second
	// Minimum score threshold (0.0 = bot, 1.0 = human)
	minScoreThreshold = 0.5
)

// RecaptchaService handles Google reCAPTCHA v3 verification
type RecaptchaService struct {
	siteKey   string
	secretKey string
	client    *http.Client
}

// recaptchaV3Response represents Google's reCAPTCHA v3 verification response
type recaptchaV3Response struct {
	Success     bool     `json:"success"`
	Score       float64  `json:"score"`       // v3: score between 0.0 and 1.0
	Action      string   `json:"action"`      // v3: action name
	ChallengeTS string   `json:"challenge_ts"`
	Hostname    string   `json:"hostname"`
	ErrorCodes  []string `json:"error-codes"`
}

// NewRecaptchaService creates a new reCAPTCHA service instance
func NewRecaptchaService(siteKey, secretKey string) *RecaptchaService {
	return &RecaptchaService{
		siteKey:   siteKey,
		secretKey: secretKey,
		client: &http.Client{
			Timeout: requestTimeout,
		},
	}
}

// GetSiteKey returns the site key for client-side widget
func (s *RecaptchaService) GetSiteKey() string {
	return s.siteKey
}

// IsConfigured returns true if the service has valid configuration
func (s *RecaptchaService) IsConfigured() bool {
	return s.siteKey != "" && s.secretKey != ""
}

// VerifyToken verifies the reCAPTCHA v3 response token with Google API
// For v3, it also checks the score against the threshold
func (s *RecaptchaService) VerifyToken(token string) (bool, error) {
	if token == "" {
		return false, ErrRecaptchaNotFound
	}

	if !s.IsConfigured() {
		// If not configured, skip verification (for development)
		return true, nil
	}

	// Prepare request data
	data := url.Values{
		"secret":   {s.secretKey},
		"response": {token},
	}

	// Send verification request to Google
	resp, err := s.client.PostForm(recaptchaVerifyURL, data)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	// Parse response
	var result recaptchaV3Response
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return false, err
	}

	if !result.Success {
		return false, ErrRecaptchaFailed
	}

	// For v3, check the score
	if result.Score < minScoreThreshold {
		return false, ErrRecaptchaLowScore
	}

	return true, nil
}
