package secretstore

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"reflect"
	"strings"

	"github.com/rs/zerolog/log"
)

type Client struct {
	StoreURL string
	Source   string

	accessKey string
}

func New(storeUrl string, source string, accessKey string) (*Client, error) {
	if storeUrl == "" || source == "" || accessKey == "" {
		return nil, errors.New("invalid parameters")
	}

	store := &Client{
		StoreURL:  storeUrl,
		Source:    source,
		accessKey: accessKey,
	}

	return store, nil
}

func (c *Client) PingServer() error {
	url := fmt.Sprintf("%s/secret/%s", c.StoreURL, c.Source)

	log.Info().Msg("pinging secret server")

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return err
	}

	req.Header.Add("x-access-key", c.accessKey)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return errors.New("secret not found")
	}

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("client: could not read response body: %s\n", err)
		os.Exit(1)
	}

	var pongResp map[string]interface{}
	err = json.Unmarshal(respBody, &pongResp)
	if err != nil {
		log.Error().Err(err)
	}

	log.Info().Msg("secret server replied")

	return nil
}

func (c *Client) LoadSecrets(s interface{}) (map[string]interface{}, error) {
	var err error
	structType := reflect.TypeOf(s)
	if structType.Kind() != reflect.Struct {
		return nil, errors.New("input param should be a struct")
	}

	structVal := reflect.ValueOf(s)
	fieldNum := structVal.NumField()

	secrets := make(map[string]interface{})
	for i := 0; i < fieldNum; i++ {
		field := structType.Field(i)

		fieldTag := field.Tag.Get("secret")
		fieldTagSplit := strings.Split(fieldTag, ",")

		tagValueLen := len(fieldTagSplit)
		fieldTagName := fieldTagSplit[0]
		if fieldTagName == "" {
			continue
		}

		omitEmpty := false
		if tagValueLen > 1 && fieldTagSplit[1] == "omitempty" {
			omitEmpty = true
		}

		secretVal, err := c.getSecret(fieldTagName)
		if err != nil {
			if err.Error() == "secret not found" && omitEmpty {
				continue
			}

			return nil, err
		}

		if secretVal == nil || *secretVal == "" {
			continue
		}

		secrets[fieldTagName] = *secretVal
	}

	return secrets, err
}

func (c *Client) SetSecret(secretName string, value string) error {
	url := fmt.Sprintf("%s/secret/%s/%s", c.StoreURL, c.Source, secretName)

	body := struct {
		Data string `json:"data"`
	}{
		Data: value,
	}

	jsonData, err := json.Marshal(body)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}

	req.Header.Add("x-access-key", c.accessKey)
	req.Header.Add("content-type", "application/json")

	_, err = http.DefaultClient.Do(req)
	if err != nil {
		return err
	}

	return nil
}

func (c *Client) getSecret(secretName string) (*string, error) {
	url := fmt.Sprintf("%s/secret/%s/%s", c.StoreURL, c.Source, secretName)

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("x-access-key", c.accessKey)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, errors.New("secret not found")
	}

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("client: could not read response body: %s\n", err)
		os.Exit(1)
	}

	var secretValue Secret
	err = json.Unmarshal(respBody, &secretValue)
	if err != nil {
		log.Error().Err(err)
	}

	return &secretValue.Data, nil
}
