package models

type SecretStore struct {
	UserID    *int `json:"userId"`
	AccountID *int `json:"accountId"`

	APIKey *string `json:"apiKey"`
}
