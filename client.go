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
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

const (
	BaseURL = "https://dev.lunchmoney.app/v1"
)

type Client struct {
	HTTPClient *http.Client
}

func New(ctx context.Context) (*Client, error) {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr}).With().Caller().Logger()

	return &Client{}, nil
}

func (c *Client) InitClient(ctx context.Context) error {
	authClient, err := NewAuthClient(ctx, BaseURL)
	if err != nil {
		return err
	}

	authTransport, err := authClient.InitTransportSession(ctx)
	if err != nil {
		return err
	}

	c.HTTPClient = authTransport

	return nil
}

func (c *Client) NewRequest(ctx context.Context, method string, path string, body io.Reader, params *url.Values) (*models.Request, error) {
	urlString := fmt.Sprintf("%s/%s", BaseURL, path)
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
		req.Header.Set("Content-Type", "application/json")
	}

	request := &models.Request{
		HTTPRequest: req,
	}

	return request, nil
}

func (c *Client) Do(ctx context.Context, req *models.Request, resp any) error {
	httpResponse, err := c.HTTPClient.Do(req.HTTPRequest)
	if err != nil {
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
