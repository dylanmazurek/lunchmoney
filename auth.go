package lunchmoney

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"runtime"
	"time"

	"github.com/hasura/go-graphql-client"
	"github.com/infamousjoeg/go-keyconfig"
	"github.com/lestrrat-go/jwx/jwt"
	"github.com/rs/zerolog/log"
)

const (
	SecretStoreKey = "lunchmoney-client"
)

type AuthClient struct {
	secretStore *models.SecretStore
}

func NewAuthClient(ctx context.Context) (*AuthClient, error) {
	authClient := &AuthClient{
		secretStore: &models.SecretStore{},
	}

	return authClient, nil
}

type addAuthHeaderTransport struct {
	T      http.RoundTripper
	APIKey *string
}

func (adt *addAuthHeaderTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	clientVersion := runtime.Version()

	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", adt.APIKey))
	req.Header.Add("User-Agent", fmt.Sprintf("github.com/dylanmazurek/lunchmoney/%s", clientVersion))

	return adt.T.RoundTrip(req)
}

func (c *AuthClient) InitTransportSession(ctx context.Context) (*http.Client, error) {
	err := keyconfig.GetConfig(SecretStoreKey, &c.secretStore)
	if err != nil || c.secretStore == nil || c.secretStore.Session == nil {
		log.Debug().Msg("session not found, logging in")

		err = c.login(ctx)

		if err != nil {
			log.Error().Msg("unable login")

			return nil, err
		}
	} else {
		log.Debug().Msg("session found, validating")
	}

	tokenPair := c.secretStore.Session.TokenPair

	authToken := tokenPair.AuthToken
	err = jwt.Validate(authToken)
	authTokenValid := (err == nil)

	if authTokenValid {
		authTokenExpiry := tokenPair.AuthToken.Expiration()
		log.Debug().Msgf("auth token valid for %.0f min", time.Until(authTokenExpiry).Minutes())

		authTransport, err := c.createAuthTransport(ctx)

		return authTransport, err
	}

	log.Debug().Msg("auth token not valid")

	refreshToken := tokenPair.RefreshToken
	err = jwt.Validate(refreshToken)
	refreshTokenValid := (err == nil)

	if refreshTokenValid {
		refreshTokenExpiry := tokenPair.AuthToken.Expiration()
		log.Debug().Msgf("refresh token valid for %.0f min, refreshing auth token", time.Until(refreshTokenExpiry).Minutes())

		err = c.refreshToken(ctx)
		if err != nil {
			return nil, err
		}

		authTransport, err := c.createAuthTransport(ctx)
		return authTransport, err
	}

	log.Info().Msg("refresh token not valid, logging in")

	err = c.login(ctx)
	if err != nil {
		return nil, err
	}

	authTransport, err := c.createAuthTransport(ctx)

	return authTransport, err
}

func (c *AuthClient) login(ctx context.Context) error {
	if c.secretStore.Email == "" || c.secretStore.Password == "" {
		email, password, err := c.getLoginDetails()
		if err != nil {
			return err
		}

		c.secretStore = &models.SecretStore{
			Email:    email,
			Password: password,
		}

		log.Debug().Msg("credentials stored")
	}

	var resp models.LoginResponse
	variables := map[string]interface{}{
		"input": models.LoginInput{
			Email:    c.secretStore.Email,
			Password: c.secretStore.Password,
		},
	}

	err := c.graphQLClient.Mutate(context.Background(), &resp, variables, graphql.OperationName("Login"))
	if err != nil {
		return err
	}

	auth := resp.LoginFunction.Auth

	c.secretStore.SetSession(auth.ID, auth.AuthToken, auth.RefreshToken)

	err = keyconfig.SetConfig(SecretStoreKey, c.secretStore)
	if err != nil {
		return err
	}

	log.Info().Msgf("login successful user id: %s", auth.ID)

	newTokenPair := c.secretStore.Session.TokenPair
	log.Debug().Msgf("auth token expires in %.0f min", time.Until(newTokenPair.AuthToken.Expiration()).Minutes())
	log.Debug().Msgf("refresh token expires in %.0f min", time.Until(newTokenPair.RefreshToken.Expiration()).Minutes())

	return nil
}

func (c *AuthClient) getLoginDetails() (string, string, error) {
	var email string
	flag.StringVar(&email, "email", "", "spaceship email address")

	var password string
	flag.StringVar(&password, "password", "", "spaceship password")

	flag.Parse()

	var err error
	for email == "" || err != nil {
		fmt.Print("email address: ")
		n, err := fmt.Scanf("%s\n", &email)

		if err != nil || n != 1 {
			fmt.Println("invalid username entered")
		}

		err = nil
	}

	for password == "" || err != nil {
		fmt.Print("password: ")
		n, err := fmt.Scanf("%s\n", &password)

		if err != nil || n != 1 {
			fmt.Println("invalid password entered")
		}

		err = nil
	}

	return email, password, err
}

func (c *AuthClient) refreshToken(ctx context.Context) error {
	_, rawRefreshToken, err := c.secretStore.GetRawTokenPair()
	if err != nil {
		return err
	}

	var resp models.RefreshTokenResponse
	variables := map[string]interface{}{
		"input": models.RefreshTokenInput{
			RefreshToken: rawRefreshToken,
		},
	}

	err = c.graphQLClient.Mutate(context.Background(), &resp, variables, graphql.OperationName("RefreshToken"))
	if err != nil {
		return err
	}

	auth := resp.RefreshTokenFunction.Auth
	c.secretStore.SetSession(auth.ID, auth.AuthToken, auth.RefreshToken)

	err = keyconfig.SetConfig(SecretStoreKey, c.secretStore)
	if err != nil {
		return err
	}

	log.Info().Msgf("auth token refresh successful user id: %s", auth.ID)

	newTokenPair := c.secretStore.Session.TokenPair
	log.Debug().Msgf("auth token expires in %.0f min", time.Until(newTokenPair.AuthToken.Expiration()).Minutes())
	log.Debug().Msgf("refresh token expires in %.0f min", time.Until(newTokenPair.RefreshToken.Expiration()).Minutes())

	return nil
}

func (c *AuthClient) createAuthTransport(ctx context.Context) (*http.Client, error) {
	rawAuthToken, _, err := c.secretStore.GetRawTokenPair()
	if err != nil {
		return nil, err
	}

	authClient := &http.Client{
		Transport: &addAuthHeaderTransport{
			T:            http.DefaultTransport,
			RawAuthToken: &rawAuthToken,
		},
	}

	return authClient, nil
}
