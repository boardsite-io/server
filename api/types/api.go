package types

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"

	gws "github.com/gorilla/websocket"
)

// Message type definitions.
const (
	MessageTypeStroke           = "stroke"
	MessageTypeUserConnected    = "userconn"
	MessageTypeUserDisconnected = "userdisc"
	MessageTypePageSync         = "pagesync"
	MessageTypePageClear        = "pageclear"
	MessageTypeMouseMove        = "mmove"
)

// Message declares the generic message envelope
// of any API JSON encoded message.
type Message struct {
	Type    string      `json:"type"`
	Sender  string      `json:"sender,omitempty"`
	Content interface{} `json:"content,omitempty"`
	Error   string      `json:"error,omitempty"`
}

// User declares some information about connected users.
type User struct {
	ID    string    `json:"id"`
	Alias string    `json:"alias"`
	Color string    `json:"color"`
	Conn  *gws.Conn `json:"-"`
}

// PageMeta declares some page meta data.
type PageMeta struct {
	Background string `json:"background"`
}

// ContentPageRequest declares the message content for page requests.
type ContentPageRequest struct {
	PageID   string `json:"pageId"`
	Index    int    `json:"index"`
	PageMeta `json:"meta"`
}

// ContentPageSync message content for page sync.
type ContentPageSync struct {
	PageRank []string    `json:"pageRank"`
	Meta     []*PageMeta `json:"meta"`
}

// ContentMouseMove declares mouse move updates.
type ContentMouseMove struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}

// type ContentUserSync map[string]*User
// type ContentPageSync []string
// type ContentPageClear []string

// NewMessage creates a new Message with any JSON encodable content,
// a message type and an optional sender.
func NewMessage(content interface{}, msgType string, sender ...string) *Message {
	var s string
	if len(sender) == 1 {
		s = sender[0]
	}
	return &Message{
		Type:    msgType,
		Sender:  s,
		Content: content,
	}
}

// NewErrorMessage creates a new Message with the error field formatted
// accoring to the error.
func NewErrorMessage(err error) *Message {
	return &Message{Error: fmt.Sprintf("%v", err)}
}

// DecodeMsgContent is a shorthand wrapper to directly decode
// the content of generic API JSON messages.
func DecodeMsgContent(r io.Reader, v interface{}) error {
	data, err := ioutil.ReadAll(r)
	if err != nil {
		return err
	}
	return UnmarshalMsgContent(data, v)
}

// UnmarshalMsgContent is a shorthand wrapper to directly unmarshal
// the content of generic API JSON messages.
func UnmarshalMsgContent(data []byte, v interface{}) error {
	m, err := UnmarshalMessage(data)
	if err != nil {
		return err
	}
	if err := m.UnmarshalContent(v); err != nil {
		return err
	}
	return nil
}

// UnmarshalMessage parses the JSON-encoded message and stores the result
// in the Message struct. The content field not parsed.
func UnmarshalMessage(data []byte) (*Message, error) {
	var content json.RawMessage
	msg := Message{Content: &content}
	if err := json.Unmarshal(data, &msg); err != nil {
		return nil, err
	}
	return &msg, nil
}

// UnmarshalContent parses the JSON-encoded content of a Message and
// stores the result in the value pointed to by v.
func (m *Message) UnmarshalContent(v interface{}) error {
	c, ok := (m.Content.(*json.RawMessage))
	if !ok {
		return errors.New("cannot unmarshal content")
	}
	if err := json.Unmarshal(*c, v); err != nil {
		return err
	}
	return nil
}
