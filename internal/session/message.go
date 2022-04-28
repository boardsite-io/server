package session

import (
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
)

// Message declares the generic message envelope
// of any API JSON encoded message.
type Message struct {
	Type     string `json:"type"`
	Sender   string `json:"sender,omitempty"`
	Receiver string `json:"-"`
	Content  any    `json:"content,omitempty"`
}

// NewMessage creates a new Message with any JSON encodable content,
// a message type and an optional sender.
func NewMessage(content any, msgType string, sender ...string) *Message {
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

// DecodeMsgContent is a shorthand wrapper to directly decode
// the content of generic API JSON messages.
func DecodeMsgContent(r io.Reader, v any) error {
	data, err := ioutil.ReadAll(r)
	if err != nil {
		return err
	}
	return UnmarshalMsgContent(data, v)
}

// UnmarshalMsgContent is a shorthand wrapper to directly unmarshal
// the content of generic API JSON messages.
func UnmarshalMsgContent(data []byte, v any) error {
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
func (m *Message) UnmarshalContent(v any) error {
	c, ok := (m.Content.(*json.RawMessage))
	if !ok {
		return errors.New("cannot unmarshal content")
	}
	if err := json.Unmarshal(*c, v); err != nil {
		return err
	}
	return nil
}
