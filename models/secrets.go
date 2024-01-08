package models

type Secrets struct {
	APIKey string `json:"apiKey"`

	UserID string `json:"userId"`
}

func (s *Secrets) HasSecrets() bool {
	return s.APIKey != ""
}
