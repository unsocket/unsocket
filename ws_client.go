package unsocket

import (
	"fmt"
	"github.com/gorilla/websocket"
	"time"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 512
)

type wsClient struct {
	url     string
	conn    *websocket.Conn
	send    chan []byte
	receive chan []byte
	// done channel acts as shutdown signal
	done chan struct{}
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
	websocket.DefaultDialer.ReadBufferSize = 102400;
	conn, _, err := websocket.DefaultDialer.Dial(c.url, nil)
	if err != nil {
		return fmt.Errorf("unable to establish websocket connection: %w", err)
	}

	c.conn = conn

	// initialize new done channel to act as shutdown signal
	c.done = make(chan struct{})

	c.send = make(chan []byte)
	c.receive = make(chan []byte)

	go c.readPump()
	go c.writePump()

	// block until the done channel is closed
	// <-c.done

	return nil
}

func (c *wsClient) handlePong(message string) error {
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	return nil
}

func (c *wsClient) readPump() {
	// c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(c.handlePong)

	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			close(c.done)
		}

		if message != nil {
			c.receive <- message
		}
	}
}

func (c *wsClient) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer ticker.Stop()

	for {
		select {
		case <-c.done:
			c.conn.WriteMessage(websocket.CloseMessage, []byte{})
			_ = c.conn.Close()
			c.conn = nil
			return
		case message := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))

			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}

			w.Write(message)
			if err := w.Close(); err != nil {
				return
			}
		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

func (c *wsClient) Stop() {
	// signal shutdown by closing done channel
	close(c.done)
}
