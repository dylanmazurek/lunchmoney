package models

type SecretStore struct {
	UserID string `json:"email"`

	APIKey *string `json:"apiKey"`
}
