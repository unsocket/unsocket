package messages

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"log"
	"testing"
)

func TestMessage_MarshalJSON(t *testing.T) {
	t.Parallel()

	connectMessage := NewConnect(&ConnectData{
		URL: "wss://example.com",
	})

	b, err := json.Marshal(&connectMessage.Message)
	if err != nil {
		log.Fatalf("unable to serialize JSON: %v", err)
	}

	expected := "{\"type\":\"connect\",\"url\":\"wss://example.com\"}"

	assert.Equal(t, expected, string(b), "they should be equal")
}

func TestOpen_MarshalJSON(t *testing.T) {
	t.Parallel()

	connectMessage := Connect{
		Message: Message{Type: ConnectType},
		ConnectData: &ConnectData{
			URL: "wss://example.com",
		},
	}

	b, err := json.Marshal(connectMessage)
	if err != nil {
		log.Fatalf("unable to serialize JSON: %v", err)
	}

	expected := "{\"type\":\"connect\",\"url\":\"wss://example.com\"}"

	assert.Equal(t, expected, string(b), "they should be equal")
}

func TestOpenConnect_MarshalJSON(t *testing.T) {
	t.Parallel()

	openMessage := Open{
		Message: Message{Type: OpenType},
	}

	connectMessage := Connect{
		Message: Message{Type: ConnectType},
		ConnectData: &ConnectData{
			URL: "wss://example.com",
		},
	}

	message := []interface{}{
		openMessage,
		connectMessage,
	}

	b, err := json.Marshal(message)
	if err != nil {
		log.Fatalf("unable to serialize JSON: %v", err)
	}

	expected := "[{\"type\":\"open\"},{\"type\":\"connect\",\"url\":\"wss://example.com\"}]"

	assert.Equal(t, expected, string(b), "they should be equal")
}

func TestMessage_UnmarshalJSON(t *testing.T) {
	t.Parallel()

	payload := []byte(`[{
		"type": "open"
	}, {
		"type": "connect",
		"url": "wss://example.com"
	}]`)

	messages := []*Message{}

	err := json.Unmarshal(payload, &messages)
	if err != nil {
		log.Fatalf("unable to parse JSON: %v", err)
	}

	assert.Len(t, messages, 2, "they should have two messages")
	assert.Equal(t, OpenType, messages[0].Type, "they should be equal")
	assert.Equal(t, ConnectType, messages[1].Type, "they should be equal")

	if connect, ok := messages[1].Get().(*Connect); ok {
		assert.Equal(t, "wss://example.com", connect.URL, "they should be equal")
	} else {
		assert.IsType(t, &Connect{}, messages[1].Get(), "they should be correctly typed")
	}
}
