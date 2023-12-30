package models

type Secrets struct {
	APIKey string `json:"apiKey"`
}

func (s *Secrets) HasSecrets() bool {
	return s.APIKey != ""
}
