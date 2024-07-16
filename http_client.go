package unsocket

import (
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/unsocket/unsocket/messages"
)

type httpClient struct {
	client *resty.Client
	url    			string
	webhookSecret	string
}

type httpClientConfig struct {
	url				string
	webhookSecret	string
}

func newHTTPClient(config *httpClientConfig) *httpClient {
	return &httpClient{
		client: resty.New(),
		url:    config.url,
		webhookSecret: config.webhookSecret,
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

	headers := map[string]string{ "Content-Type":  "application/json" }

	if len(c.webhookSecret) > 0 {
		headers["Authorization"] = "Bearer " + c.webhookSecret
	}

	res, err := c.client.R().
		SetHeaders(headers).
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
