package unsocket

import (
	"errors"
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/unsocket/unsocket/messages"
	"regexp"
	"strings"
)

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

	log.Info("sending READY to webhook")

	// send the ready message to webhook
	res, err := u.httpClient.request([]*messages.Message{
		&messages.NewReady(&messages.ReadyData{}).Message,
	})
	if err != nil {
		return fmt.Errorf("unable to send ready message: %w", err)
	}

	log.Infof("received %d messages to be processed", len(res.messages))

	// expect at least one connect message from webhook
	if len(res.messages) < 1 {
		return errors.New("must not get empty reply to ready message")
	}

	// look for the connect message
	connect, ok := res.messages[0].Get().(*messages.Connect)
	if !ok {
		return errors.New("must get connect message to ready message")
	}

	log.Infof("connecting to %s", connect.URL)

	headers := make(map[string][]string)

	for key, value := range connect.Headers {
		headers[key] = []string{value}
	}

	wsClient := newWSClient(&wsClientConfig{
		url: connect.URL,
		header: headers,
	})

	err = wsClient.RunAndWait()
	if err != nil {
		return fmt.Errorf("unable to run websocket client: %w", err)
	}

	var backlog []*messages.Message

	// append any excess messages from init to be handled
	backlog = append(backlog, res.messages[1:]...)

	// handle incoming messages
	for {
		for len(backlog) > 0 {
			m := backlog[0]
			backlog = backlog[1:]

			log.Infof("processing %s message", strings.ToUpper(string(m.Type)))

			switch message := m.Get().(type) {
			case *messages.Text:
				wsClient.send <- []byte(message.Text)
			case *messages.Exclude:
				// TODO: register exclude filter
			}
		}

		select {
		case <-wsClient.error:
			return errors.New("Websocket client exited with an error")
		case <-u.done:
			goto Escape
		case text := <-wsClient.receive:
			// TODO: apply registered exclude filters
			if matched, _ := regexp.Match(`heartbeat`, text); matched {
				log.Debug("excluding websocket message")
				continue
			}

			log.Info("processing websocket message")

			log.Info(string(text))

			res, err := u.httpClient.request([]*messages.Message{
				&messages.NewText(&messages.TextData{
					Text: string(text),
				}).Message,
			})
			if err != nil {
				log.Errorf("request failed with %s", err.Error())
				continue
			}

			if len(res.messages) > 0 {
				log.Infof("received %d messages to be processed", len(res.messages))

				backlog = append(backlog, res.messages...)
			}
		}
	}

Escape:
	wsClient.Stop()

	log.Info("stopped processing")

	return nil
}

func (u *Unsocket) Stop() error {
	// signal shutdown by closing done channel
	close(u.done)

	return nil
}
