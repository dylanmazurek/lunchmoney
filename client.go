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
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type Client struct {
	ctx context.Context

	HTTPClient *http.Client
	State      string

	UserID *string
}

func New(ctx context.Context) (*Client, error) {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr}).With().Caller().Logger()

	return &Client{
		ctx:   ctx,
		State: constants.ClientState.New,
	}, nil
}

func (c *Client) InitClient(newCredentials *models.Secrets) error {
	newAuthClient, err := NewAuthClient(c.ctx)
	if err != nil {
		c.State = constants.ClientState.Error
		return err
	}

	if newCredentials.APIKey != "" {
		err := newAuthClient.SetSecrets(*newCredentials)
		if err != nil {
			return err
		}
	}

	authTransport, err := newAuthClient.InitTransportSession()
	if err != nil {
		return err
	}

	c.HTTPClient = authTransport
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
	httpResponse, err := c.HTTPClient.Do(req.HTTPRequest)
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
