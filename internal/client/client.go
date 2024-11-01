package devhub

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"
)

const HostURL string = "http://localhost:4000"

type Client struct {
	HostURL    string
	HTTPClient *http.Client
	ApiKey     string
}

func NewClient(host, apiKey *string) (*Client, error) {
	c := Client{
		HTTPClient: &http.Client{Timeout: 10 * time.Second},
		HostURL:    HostURL,
	}

	if host != nil {
		c.HostURL = *host
	}

	if apiKey == nil {
		return &c, nil
	}

	c.ApiKey = *apiKey

	return &c, nil
}

func (c *Client) doRequest(req *http.Request) ([]byte, error) {
	apiKey := c.ApiKey

	req.Header.Set("x-api-key", apiKey)
	req.Header.Set("content-type", "application/json")

	res, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	if res.StatusCode == http.StatusNotFound {
		return nil, errors.New("not found")
	}

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("status: %d, body: %s", res.StatusCode, body)
	}

	return body, err
}
