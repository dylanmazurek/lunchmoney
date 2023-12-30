package lunchmoney

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"os"

	"github.com/dylanmazurek/lunchmoney/models"
	"github.com/dylanmazurek/lunchmoney/util/constants"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type Client struct {
	Client *http.Client
	State  string
}

func New(ctx context.Context) (*Client, error) {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr}).With().Caller().Logger()

	return &Client{
		Client: nil,
		State:  constants.ClientState.New,
	}, nil
}

func (c *Client) InitClient(ctx context.Context, newCredentials *models.Secrets) error {
	authClient, err := NewAuthClient(ctx)
	if err != nil {
		c.State = constants.ClientState.Error
		return err
	}

	authTransport, err := authClient.InitTransportSession(ctx)
	if err != nil {
		c.State = constants.ClientState.Error
		return err
	}

	if authClient.lunchmoneySecrets == nil || authClient.lunchmoneySecrets.APIKey == "" {
		if newCredentials == nil {
			log.Warn().Msg("api key is not set, set API_KEY env")
			return nil
		} else {
			log.Info().Msg("API_KEY env set, saving to store")
			err := authClient.SetSecrets(*newCredentials)
			if err != nil {
				return err
			}
		}
	}

	c.Client = authTransport

	c.State = constants.ClientState.Authenticated

	return nil
}

func (c *Client) Do(ctx context.Context, req *http.Request, resp interface{}, query url.Values) error {
	httpResponse, err := c.Client.Do(req)
	if err != nil {
		if httpResponse.StatusCode == http.StatusUnauthorized {
			log.Debug().Msg("http response unauthorized, logging in")

			c.InitClient(ctx, nil)
		}

		return err
	}
	defer httpResponse.Body.Close()

	req.URL.RawQuery = query.Encode()

	bodyBytes, err := io.ReadAll(httpResponse.Body)
	if err != nil {
		return err
	}

	err = json.Unmarshal(bodyBytes, &resp)
	if err != nil {
		return err
	}

	return nil
}
