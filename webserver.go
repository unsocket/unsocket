package unsocket

import (
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"

	log "github.com/sirupsen/logrus"
	"github.com/unsocket/unsocket/messages"
)

type webserver struct {
	srv    	*http.Server
	error chan struct{}
	receive chan messages.Text
}

func newWebserver() *webserver {
	return &webserver{
		srv: &http.Server{
			Addr: ":3009",
		},
		error: make(chan struct{}),
	}
}

func (c *webserver) RunAndWait() error {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		var bodyBytes []byte
		var err error
		if r.Body != nil {
			bodyBytes, err = ioutil.ReadAll(r.Body)
			if err != nil {
				fmt.Printf("Body reading error: %v", err)
				return
			}
			defer r.Body.Close()
		}

	    m := messages.Text{}
	    if err := m.UnmarshalJSON(bodyBytes); err == nil {
		    c.receive <- m
		} else{
			log.Errorf("Error during unmarshal: %v", err)
			http.Error(w, "unsupported message type", 422)
		}
    })

	c.receive = make(chan messages.Text)

    log.Printf("starting to listen on %v", c.srv.Addr)

    go c.listen()

    return nil
}

func (c *webserver) listen() {
	err := c.srv.ListenAndServe()

	if errors.Is(err, http.ErrServerClosed) {
		log.Printf("webserver closed")
	} else if err != nil {
		log.Printf("error starting web server: %s\n", err)
		c.error <- struct{}{}
	}
}

func (c *webserver) Stop() {
	if err := c.srv.Shutdown(context.Background()); err != nil {
		log.Printf("HTTP server Shutdown: %v", err)
	}
}
