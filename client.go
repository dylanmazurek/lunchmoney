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

	"github.com/dylanmazurek/lunchmoney/models"
	"github.com/google/go-querystring/query"
)

const (
	// BaseAPIURL is the base url we use for all API requests.
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
	HTTP *http.Client
	Base *url.URL
}

// NewClient creates a new client with the specified API key.
func NewClient(apikey string) (*Client, error) {
	base, err := url.Parse(BaseAPIURL)
	if err != nil {
		return nil, fmt.Errorf("invalid base URI: %w", err)
	}

	return &Client{
		HTTP: &http.Client{
			Transport: &addAuthHeaderTransport{T: http.DefaultTransport, Key: apikey},
		},
		Base: base,
	}, nil
}

type ResponseType interface {
	models.TransactionsResponse | models.CategoriesResponse | models.AssetsResponse
}

// Request makes a request using the client
func Request[T ResponseType](ctx context.Context, c *Client, reqOptions models.RequestOptions) (transResp *T, err error) {
	url, err := url.Parse(c.Base.String())
	if err != nil {
		return nil, fmt.Errorf("bad path: %w", err)
	}

	url.Path = reqOptions.Path
	vals, err := query.Values(reqOptions.QueryValues)
	if err != nil {
		return nil, fmt.Errorf("bad query values: %w", err)
	}

	url.RawQuery = vals.Encode()

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

	resp, err := c.HTTP.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request (%+v) failed: %w", req, err)
	}

	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(bodyBytes, &transResp)
	if err != nil {
		//var serverError models.ErrorResponse
		//serverError = unmarshalServerError(bodyBytes)

		return nil, errors.New("unable to unmarshal")
	}

	return transResp, nil
}
