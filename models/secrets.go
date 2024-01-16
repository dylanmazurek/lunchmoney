package models

type Secrets struct {
	APIKey string `secret:"api_key"`

	UserID string `secret:"user_id,omitempty"`
}

func (s *Secrets) HasSecrets() bool {
	return s.APIKey != ""
}
