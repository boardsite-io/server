package types

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUnmarshalMsg(t *testing.T) {
	tests := []struct {
		msg  string
		want *Message
		err  error
	}{
		{"", nil, assert.AnError},
		{`{{"content": {}`, nil, assert.AnError},
		{`{"content": {"strokeType": "0}}`, nil, assert.AnError},
		{`{}`, &Message{}, nil},
		{`{"type":"sometype"}`, &Message{Type: "sometype"}, nil},
		{
			`{"type":"sometype","sender":"heat","content":{}}`,
			&Message{Type: "sometype", Sender: "heat"},
			nil,
		},
	}

	for _, test := range tests {
		m, err := UnmarshalMessage([]byte(test.msg))
		if test.err != nil {
			assert.Error(t, err)
			continue
		} else {
			assert.NoError(t, err)
		}
		assert.Equal(t, test.want.Type, m.Type, "incorrect marshalling of message type")
		assert.Equal(t, test.want.Sender, m.Sender, "incorrect marshalling of message sender")
	}
}

func TestUnmarshalMsgContent(t *testing.T) {
	tests := []struct {
		msg  string
		want interface{}
		err  error
	}{
		{"", nil, assert.AnError},
		{"{}", nil, assert.AnError},
		{`{{"content": {}`, nil, assert.AnError},
		{`{"content": {"strokeType": "0}}`, nil, assert.AnError},
		{
			`{"type":"sometype","sender":"heat","content":{"strokeType":0,"strokeId":"strokeid1","userId":"user1"}}`,
			Stroke{ID: "strokeid1", UserID: "user1", Type: 0},
			nil,
		},
		{
			`{"type":"sometype","sender":"heat","content":{"id":"user1","alias":"user1","color":"#ff00ff"}}`,
			User{ID: "user1", Alias: "user1", Color: "#ff00ff"},
			nil,
		},
	}

	for _, test := range tests {
		c := reflect.ValueOf(test.want)
		w := reflect.ValueOf(test.want)
		if test.err != nil {
			assert.Error(t, UnmarshalMsgContent([]byte(test.msg), &c))
		} else {
			assert.NoError(t, UnmarshalMsgContent([]byte(test.msg), &c))
		}
		assert.Equal(t, w, c, "incorrect unmarshalling of message content")
	}
}
