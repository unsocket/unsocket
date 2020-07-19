package unsocket

import (
	"fmt"
	"github.com/go-resty/resty/v2"
)

type httpClient struct {
	client *resty.Client
	url    string
}

type httpClientConfig struct {
	url string
}

func newHTTPClient(config *httpClientConfig) *httpClient {
	return &httpClient{
		client: resty.New(),
		url:    config.url,
	}
}

type httpClientInitResponse struct {
	ws string
}

func (c *httpClient) init() (*httpClientInitResponse, error) {
	type request struct{}
	var response struct {
		Url string `json:"url"`
	}

	res, err := c.client.R().
		SetHeader("Content-Type", "application/json").
		SetBody(&request{}).
		SetResult(&response).
		Post(c.url)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}

	if res.IsError() {
		return nil, fmt.Errorf("%v", res.String())
	}

	return &httpClientInitResponse{
		ws: response.Url,
	}, nil
}
