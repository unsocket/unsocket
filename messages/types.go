package messages

import (
	"encoding/json"
	"fmt"
)

type Type string

const (
	ReadyType   Type = "ready"
	ConnectType Type = "connect"
	OpenType    Type = "open"
	TextType    Type = "text"
	ExcludeType Type = "exclude"
)

type Message struct {
	Type Type `json:"type"`

	// self holds a pointer to the outer struct
	self interface{}
}

type ReadyData struct{}
type Ready struct {
	Message
	*ReadyData
}

type ConnectData struct {
	URL string `json:"url"`
	Headers map[string]string `json:"headers"`
}
type Connect struct {
	Message
	*ConnectData
}

type OpenData struct{}
type Open struct {
	Message
	*OpenData
}

type TextData struct {
	Text string `json:"text"`
}
type Text struct {
	Message
	*TextData
}

type ExcludeData struct {
	Pattern string `json:"pattern"`
}
type Exclude struct {
	Message
	*ExcludeData
}

func NewReady(data *ReadyData) *Ready {
	message := Ready{Message: Message{Type: ReadyType}, ReadyData: data}
	message.Message.self = &message
	return &message
}

func NewConnect(data *ConnectData) *Connect {
	message := Connect{Message: Message{Type: ConnectType}, ConnectData: data}
	message.Message.self = &message
	return &message
}

func NewOpen(data *OpenData) *Open {
	message := Open{Message: Message{Type: OpenType}, OpenData: data}
	message.Message.self = &message
	return &message
}

func NewText(data *TextData) *Text {
	message := Text{Message: Message{Type: TextType}, TextData: data}
	message.Message.self = &message
	return &message
}

func NewExclude(data *ExcludeData) *Exclude {
	message := Exclude{Message: Message{Type: ExcludeType}, ExcludeData: data}
	message.Message.self = &message
	return &message
}

func (m *Message) Get() interface{} {
	return m.self
}

func (m *Message) MarshalJSON() ([]byte, error) {
	var data interface{}

	switch v := m.self.(type) {
	case *Ready:
		data = struct {
			Type Type `json:"type"`
			*ReadyData
		}{
			Type: ReadyType,
			ReadyData: v.ReadyData,
		}
	case *Connect:
		data = struct {
			Type Type `json:"type"`
			*ConnectData
		}{
			Type: ConnectType,
			ConnectData: v.ConnectData,
		}
	case *Open:
		data = struct {
			Type Type `json:"type"`
			*OpenData
		}{
			Type: OpenType,
			OpenData: v.OpenData,
		}
	case *Text:
		data = struct {
			Type Type `json:"type"`
			*TextData
		}{
			Type: TextType,
			TextData: v.TextData,
		}
	case *Exclude:
		data = struct {
			Type Type `json:"type"`
			*ExcludeData
		}{
			Type: ExcludeType,
			ExcludeData: v.ExcludeData,
		}
	default:
		return nil, fmt.Errorf("invalid type %T", m.Type)
	}

	return json.Marshal(data)
}

func (m *Message) UnmarshalJSON(b []byte) error {
	base := struct {
		Type Type `json:"type"`
	}{}

	if err := json.Unmarshal(b, &base); err != nil {
		return fmt.Errorf("unable to unmarshal: %w", err)
	}

	message := Message{Type: base.Type}
	var data interface{}

	switch base.Type {
	case ReadyType:
		readyData := &ReadyData{}
		*m = NewReady(readyData).Message
		data = readyData
	case ConnectType:
		connectData := &ConnectData{}
		*m = NewConnect(connectData).Message
		data = connectData
	case OpenType:
		openData := &OpenData{}
		*m = NewOpen(openData).Message
		data = openData
	case TextType:
		textData := &TextData{}
		*m = NewText(textData).Message
		data = textData
	case ExcludeType:
		excludeData := &ExcludeData{}
		*m = NewExclude(excludeData).Message
		data = excludeData
	default:
		return fmt.Errorf("invalid type %s", message.Type)
	}

	if err := json.Unmarshal(b, data); err != nil {
		return fmt.Errorf("unable to unmarshal %s type: %w", message.Type, err)
	}

	return nil
}
