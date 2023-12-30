package lunchmoney

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/dylanmazurek/lunchmoney/models"
	"github.com/dylanmazurek/lunchmoney/util/constants"
	"github.com/dylanmazurek/lunchmoney/util/secretstore"
	"github.com/rs/zerolog/log"
)

const (
	SecretStoreShelfKey = "lunchmoney-client"
)

type AuthClient struct {
	httpClient *http.Client

	SecretStore *secretstore.SecretStore

	lunchmoneySecrets *models.Secrets
}

func NewAuthClient(ctx context.Context) (*AuthClient, error) {
	authClient := &AuthClient{
		httpClient: &http.Client{Transport: http.DefaultTransport},
	}

	return authClient, nil
}

type addAuthHeaderTransport struct {
	T      http.RoundTripper
	APIKey *string
}

func (adt *addAuthHeaderTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", *adt.APIKey))
	req.Header.Add("User-Agent", constants.Config.SourceUserAgent)

	return adt.T.RoundTrip(req)
}

func (c *AuthClient) InitTransportSession(ctx context.Context) (*http.Client, error) {
	secretStore, err := secretstore.New(SecretStoreShelfKey)
	if err != nil {
		return nil, err
	}

	c.SecretStore = secretStore

	shelfDataBytes, err := secretStore.GetShelfData()
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(shelfDataBytes, &c.lunchmoneySecrets)
	if err != nil {
		return nil, err
	}

	currentAPIKey := c.lunchmoneySecrets.APIKey
	if currentAPIKey == "" {
		log.Error().Msg("api key is not set")

		return nil, err
	}

	userData, err := c.getUserData(c.lunchmoneySecrets.APIKey)
	if err != nil {
		return nil, err
	}

	log.Info().Msgf("user data fetched for %s", userData.UserName)

	authTransport, err := c.createAuthTransport(ctx)

	return authTransport, err
}

func (c *AuthClient) SetSecrets(secrets models.Secrets) error {
	if !secrets.HasSecrets() {
		return errors.New("no secrets provided")
	}

	jsonSecrets, err := json.Marshal(secrets)
	if err != nil {
		return err
	}

	err = c.SecretStore.UpdateShelf(SecretStoreShelfKey, jsonSecrets)
	if err != nil {
		return err
	}

	log.Debug().Msg("secrets stored")

	return nil
}

func (c *AuthClient) getUserData(apiKey string) (*models.User, error) {
	path := fmt.Sprintf("%s/%s", constants.Config.APIBaseURL, constants.Path.Me)

	req, err := http.NewRequest("GET", path, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", apiKey))

	resp, err := http.DefaultClient.Do(req)

	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("invalid api key")
	}

	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var me models.User
	err = json.Unmarshal(bodyBytes, &me)
	if err != nil {
		return nil, err
	}

	return &me, nil
}

func (c *AuthClient) createAuthTransport(ctx context.Context) (*http.Client, error) {
	shelfDataBytes, err := c.SecretStore.GetShelfData()
	if err != nil {
		return nil, err
	}

	secrets := &models.Secrets{}
	err = json.Unmarshal(shelfDataBytes, &secrets)
	if err != nil {
		return nil, err
	}

	authClient := &http.Client{
		Transport: &addAuthHeaderTransport{
			T:      http.DefaultTransport,
			APIKey: &c.lunchmoneySecrets.APIKey,
		},
	}

	return authClient, nil
}
