package lunchmoney

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

const (
	BaseURL = "https://dev.lunchmoney.app"
)

type Client struct {
	HTTPClient *http.Client
}

func New(ctx context.Context) (*Client, error) {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr}).With().Caller().Logger()

	return &Client{}, nil
}

func (c *Client) Do(ctx context.Context, req *http.Request, resp interface{}) error {
	httpResponse, err := c.HTTPClient.Do(req)
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
