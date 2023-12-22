package lunchmoney

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/dylanmazurek/lunchmoney/models"
	"github.com/infamousjoeg/go-keyconfig"
	"github.com/rs/zerolog/log"
)

const (
	SecretStoreKey = "lunchmoney-client"
)

type AuthClient struct {
	secretStore *models.SecretStore

	BaseURL string
}

func NewAuthClient(ctx context.Context, baseUrl string) (*AuthClient, error) {
	authClient := &AuthClient{
		secretStore: &models.SecretStore{},
		BaseURL:     baseUrl,
	}

	return authClient, nil
}

type addAuthHeaderTransport struct {
	T      http.RoundTripper
	APIKey *string
}

func (adt *addAuthHeaderTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", *adt.APIKey))
	req.Header.Add("User-Agent", "github.com/dylanmazurek/lunchmoney")

	return adt.T.RoundTrip(req)
}

func (c *AuthClient) InitTransportSession(ctx context.Context) (*http.Client, error) {
	newLogin := false
	err := keyconfig.GetConfig(SecretStoreKey, &c.secretStore)
	if err != nil || c.secretStore == nil || c.secretStore.APIKey == nil || *c.secretStore.APIKey == "" {
		log.Debug().Msg("api key not found, logging in")

		err = c.login(ctx)

		if err != nil {
			return nil, err
		}

		newLogin = true
	}

	if !newLogin {
		log.Debug().Msg("api key found, validating")

		_, err = c.getKeyUser()
		if err != nil {
			log.Debug().Msg("api key invalid, logging in")

			err = c.login(ctx)

			if err != nil {
				return nil, err
			}
		}
	}

	authTransport, err := c.createAuthTransport(ctx)

	return authTransport, err
}

func (c *AuthClient) login(ctx context.Context) error {
	var err error
	var user *models.User
	var apiKey string

	if c.secretStore == nil || c.secretStore.APIKey == nil {
		apiKey, user, err = c.getLoginDetails()
		if err != nil {
			return err
		}
	}

	c.secretStore = &models.SecretStore{
		UserID:    &user.UserID,
		AccountID: &user.AccountID,
		APIKey:    &apiKey,
	}

	err = keyconfig.SetConfig(SecretStoreKey, c.secretStore)
	if err != nil {
		return err
	}

	log.Info().Msgf("api key confirmed for %s", user.UserEmail)

	return nil
}

func (c *AuthClient) getKeyUser() (*models.User, error) {
	path := fmt.Sprintf("%s/%s", c.BaseURL, "me")

	req, err := http.NewRequest("GET", path, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", *c.secretStore.APIKey))

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

func (c *AuthClient) getLoginDetails() (string, *models.User, error) {
	var err error
	var keyUser *models.User

	var apiKey string
	for apiKey == "" || err != nil {
		fmt.Print("api key: ")
		n, err := fmt.Scanf("%s\n", &apiKey)

		if err != nil || n != 1 {
			fmt.Println("invalid api key entered")
		}

		c.secretStore.APIKey = &apiKey

		keyUser, err = c.getKeyUser()
		if err != nil {
			log.Info().Msg("unable to validate key")
		}
	}

	return apiKey, keyUser, nil
}

func (c *AuthClient) createAuthTransport(ctx context.Context) (*http.Client, error) {
	authClient := &http.Client{
		Transport: &addAuthHeaderTransport{
			T:      http.DefaultTransport,
			APIKey: c.secretStore.APIKey,
		},
	}

	return authClient, nil
}
