package unsocket

import (
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/unsocket/unsocket/messages"
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

type httpClientResponse struct {
	messages []*messages.Message
}

func (c *httpClient) request(msgs []*messages.Message) (*httpClientResponse, error) {
	type request struct {
		Messages []*messages.Message `json:"messages"`
	}
	var response struct {
		Messages []*messages.Message `json:"messages"`
	}

	res, err := c.client.R().
		SetHeader("Content-Type", "application/json").
		SetBody(&request{
			Messages: msgs,
		}).
		SetResult(&response).
		Post(c.url)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}

	if res.IsError() {
		return nil, fmt.Errorf("%v", res.String())
	}

	return &httpClientResponse{
		messages: response.Messages,
	}, nil
}
