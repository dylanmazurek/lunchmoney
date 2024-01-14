package secretstore

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/rs/zerolog/log"
)

type SecretStore struct {
	StoreURL string
	Source   string

	accessKey string
}

func (s *SecretStore) New(storeUrl string, source string, accessKey string) (*SecretStore, error) {
	if storeUrl == "" || source == "" || accessKey == "" {
		return nil, errors.New("invalid parameters")
	}

	store := &SecretStore{
		StoreURL:  storeUrl,
		Source:    source,
		accessKey: accessKey,
	}

	return store, nil
}

func (s *SecretStore) GetSecret(secretName string) (*string, error) {
	resp, err := http.Get(s.StoreURL + "/secret/" + s.Source + "/" + secretName + "?accessKey=" + s.accessKey)

	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, errors.New("secret not found")
	}

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("client: could not read response body: %s\n", err)
		os.Exit(1)
	}

	var secretValue Secret
	err = json.Unmarshal(respBody, &secretValue)
	if err != nil {
		log.Error.Err(err)
	}

	value := string(respBody)
	return &value, nil
}
