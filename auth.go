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

type AuthClient struct {
	Ctx        context.Context
	httpClient *http.Client

	secretStoreClient *secretstore.Client

	secrets models.Secrets
}

func NewAuthClient(ctx context.Context, ssc *secretstore.Client) (*AuthClient, error) {
	authClient := &AuthClient{
		Ctx:        ctx,
		httpClient: &http.Client{Transport: http.DefaultTransport},

		secretStoreClient: ssc,
		secrets:           models.Secrets{},
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

func (c *AuthClient) InitTransportSession() (*http.Client, error) {
	currentAPIKey := c.secrets.APIKey
	if currentAPIKey == "" {
		err := fmt.Errorf("api key is not set")

		return nil, err
	}

	user, err := c.getUserData(c.secrets.APIKey)
	if err != nil {
		err := fmt.Errorf("api key not valid")

		return nil, err
	}

	log.Info().
		Str("username", user.UserName).
		Msgf("user data fetched")

	authTransport, err := c.createAuthTransport()

	return authTransport, err
}

func (c *AuthClient) SetSecrets(secrets models.Secrets) error {
	if !secrets.HasSecrets() {
		return errors.New("no secrets provided")
	}

	sessionBytes, err := json.Marshal(secrets.UserID)
	if err != nil {
		return err
	}

	err = c.secretStoreClient.SetSecret("session", string(sessionBytes))
	if err != nil {
		return err
	}

	log.Debug().Msg("secrets stored")

	return nil
}

func (c *AuthClient) getUserData(apiKey string) (*models.User, error) {
	path := fmt.Sprintf("%s%s", constants.Config.APIBaseURL, constants.Path.Me)

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

func (c *AuthClient) createAuthTransport() (*http.Client, error) {
	authClient := &http.Client{
		Transport: &addAuthHeaderTransport{
			T:      http.DefaultTransport,
			APIKey: &c.secrets.APIKey,
		},
	}

	return authClient, nil
}
