package secretstore

type Secret struct {
	Source string `json:"source"`
	Key    string `json:"key"`

	Data string `json:"data"`
}
