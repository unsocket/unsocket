package unsocket

import "fmt"

type Config struct {
	WebhookURL string
}

type Unsocket struct {
	// done channel acts as shutdown signal
	done chan struct{}

	// httpClient manages all webhook communication
	httpClient *httpClient
}

func NewUnsocket(config *Config) (*Unsocket, error) {
	httpClient := newHTTPClient(&httpClientConfig{
		url: config.WebhookURL,
	})

	return &Unsocket{
		httpClient: httpClient,
	}, nil
}

func (u *Unsocket) RunAndWait() error {
	// initialize new done channel to act as shutdown signal
	u.done = make(chan struct{})

	res, err := u.httpClient.init()
	if err != nil {
		return fmt.Errorf("unable to init %w", err)
	}

	fmt.Printf("connect to %s", res.ws)

	wsClient := newWSClient(&wsClientConfig{
		url: res.ws,
	})

	err = wsClient.RunAndWait()
	if err != nil {
		return fmt.Errorf("unable to run websocket client: %w", err)
	}

	// block until the done channel is closed
	<-u.done

	return nil
}

func (u *Unsocket) Stop() error {
	// signal shutdown by closing done channel
	close(u.done)

	return nil
}
