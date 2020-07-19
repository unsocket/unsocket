package unsocket

import (
	"fmt"
	"github.com/gorilla/websocket"
)

type wsClient struct {
	url  string
	conn *websocket.Conn
}

type wsClientConfig struct {
	url string
}

func newWSClient(config *wsClientConfig) *wsClient {
	return &wsClient{
		url: config.url,
	}
}

func (c *wsClient) RunAndWait() error {
	conn, _, err := websocket.DefaultDialer.Dial(c.url, nil)
	if err != nil {
		return fmt.Errorf("unable to establish websocket connection: %w", err)
	}

	c.conn = conn

	for {
		messageType, message, err := conn.ReadMessage()
		if err != nil {
			return fmt.Errorf("unable to read message: %w", err)
		} else if messageType == websocket.CloseNormalClosure {
			break
		}

		fmt.Printf("received message: %s", message)
	}

	return nil
}

func (c *wsClient) Stop() {
	_ = c.conn.Close()
}
