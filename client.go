package lunchmoney

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"

	"github.com/dylanmazurek/lunchmoney/models"
	"github.com/dylanmazurek/lunchmoney/util/constants"
	"github.com/dylanmazurek/lunchmoney/util/secretstore"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type Client struct {
	ctx context.Context

	Client *http.Client
	State  string

	UserID *string
}

func New(ctx context.Context) (*Client, error) {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr}).With().Caller().Logger()

	return &Client{
		ctx:   ctx,
		State: constants.ClientState.New,
	}, nil
}

func (c *Client) InitClient(ssc *secretstore.Client) error {
	newAuthClient, err := NewAuthClient(c.ctx, ssc)
	if err != nil {
		c.State = constants.ClientState.Error
		return err
	}

	var secrets models.Secrets
	secretsMap, err := ssc.LoadSecrets(secrets)
	if err != nil {
		log.Error().Err(err).Msg("failed to get secrets")
	}

	newAuthClient.secrets = models.Secrets{
		APIKey: secretsMap["api_key"].(string),
	}

	authTransport, err := newAuthClient.InitTransportSession()
	if err != nil {
		return err
	}

	c.Client = authTransport
	c.UserID = &newAuthClient.secrets.UserID

	return nil
}

func (c *Client) NewRequest(method string, path string, body io.Reader, params *url.Values) (*models.Request, error) {
	urlString := fmt.Sprintf("%s%s", constants.Config.APIBaseURL, path)
	requestUrl, err := url.Parse(urlString)
	if err != nil {
		return nil, err
	}

	if params != nil {
		requestUrl.RawQuery = params.Encode()
	}

	req, err := http.NewRequest(method, requestUrl.String(), body)
	if err != nil {
		return nil, err
	}

	if body != nil {
		req.Header.Add("Content-Type", "application/json")
	}

	request := &models.Request{
		HTTPRequest: req,
	}

	return request, nil
}

func (c *Client) Do(req *models.Request, resp interface{}) error {
	httpResponse, err := c.Client.Do(req.HTTPRequest)
	if err != nil {
		if httpResponse != nil {
			log.Debug().Msgf("http response error: %s", httpResponse.Status)
		}

		return err
	}
	defer httpResponse.Body.Close()

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
