package unsocket

import (
	"fmt"
	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
	"time"
	"net/http"
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
	header	http.Header
	conn    *websocket.Conn
	send    chan []byte
	receive chan []byte
	// done channel acts as shutdown signal
	done chan struct{}
	// error channel allows to signal an error from go routines
	error chan struct{}
}

type wsClientConfig struct {
	url string
	header http.Header
}

func newWSClient(config *wsClientConfig) *wsClient {
	return &wsClient{
		url: config.url,
		header: config.header,
	}
}

func (c *wsClient) RunAndWait() error {
	conn, _, err := websocket.DefaultDialer.Dial(c.url, c.header)
	if err != nil {
		return fmt.Errorf("unable to establish websocket connection: %w", err)
	}

	c.conn = conn

	// initialize new done channel to act as shutdown signal
	c.done = make(chan struct{})
	c.error = make(chan struct{})

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
	c.conn.SetPongHandler(c.handlePong)

	for {
		_, message, err := c.conn.ReadMessage()

		c.conn.SetReadDeadline(time.Now().Add(pongWait))

		if err != nil {
			log.Errorf("error reading message: %v", err)
			close(c.error)
			close(c.done)
			return
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
