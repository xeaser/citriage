package httpclient

import (
	"net/http"
)

type ApiClient interface {
	Get(url string) (*http.Response, error)
}

type client struct {
	httpClient http.Client
}

func NewClient(httpClient http.Client) ApiClient {
	return &client{
		httpClient: httpClient,
	}
}

func (c *client) Get(url string) (*http.Response, error) {
	resp, err := c.httpClient.Get(url)
	if err != nil {
		return nil, err
	}

	return resp, nil
}
