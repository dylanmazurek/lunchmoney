package lunchmoney

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/dylanmazurek/lunchmoney/models"
	"github.com/rs/zerolog"
)

const (
	BaseAPIURL = "https://dev.lunchmoney.app/"
)

type addAuthHeaderTransport struct {
	T   http.RoundTripper
	Key string
}

func (adt *addAuthHeaderTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	if adt.Key == "" {
		return nil, fmt.Errorf("no key provided")
	}

	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", adt.Key))
	req.Header.Add("User-Agent", "github.com/dylanmazurek/lunchmoney/0.0.0")

	return adt.T.RoundTrip(req)
}

// Client holds our base configuration for our LunchMoney client.
type Client struct {
	HTTP    *http.Client
	BaseURL *url.URL

	Logger *zerolog.Logger

	TransactionCache models.TransactionCache
}

// NewClient creates a new client with the specified API key.
func NewClient(ctx context.Context, apikey string) (*Client, error) {
	baseUrl, err := url.Parse(BaseAPIURL)
	if err != nil {
		return nil, fmt.Errorf("invalid base url: %w", err)
	}

	return &Client{
		HTTP: &http.Client{
			Transport: &addAuthHeaderTransport{T: http.DefaultTransport, Key: apikey},
		},
		BaseURL: baseUrl,

		Logger: zerolog.Ctx(ctx),
	}, nil
}

// Request makes a request using the client
func Request(ctx context.Context, c *Client, reqOptions models.Request) (resp *models.Response, err error) {
	url, err := url.Parse(c.BaseURL.String())
	if err != nil {
		return nil, fmt.Errorf("bad path: %w", err)
	}

	url.Path = reqOptions.Path
	url.RawQuery = reqOptions.QueryValues.Encode()

	req, err := http.NewRequest(reqOptions.Method, url.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("could not create new request: %w", err)
	}

	req.Close = true

	if reqOptions.ReqBody != nil {
		jsonReqBody, err := json.Marshal(reqOptions.ReqBody)
		if err != nil {
			return nil, fmt.Errorf("could not marshal request body: %w", err)
		}

		req, _ = http.NewRequest(reqOptions.Method, url.String(), bytes.NewBuffer(jsonReqBody))
		req.Header.Set("Content-Type", "application/json")
	}

	httpResp, err := c.HTTP.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request (%+v) failed: %w", req, err)
	}

	defer httpResp.Body.Close()

	bodyBytes, err := io.ReadAll(httpResp.Body)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(bodyBytes, &resp)
	if err != nil {
		return nil, errors.New("unable to unmarshal")
	}

	if resp.Errors != nil {
		err = errors.New(strings.Join(*resp.Errors, ", "))
	}

	return resp, err
}
