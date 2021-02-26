package types

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUnmarshalMsgContent(t *testing.T) {
	tests := []struct {
		msg  string
		want Stroke
	}{
		{
			`{"id": 0, "type": "", "content": {"type":0,"id":"strokeid1","userId":"user1"}}`,
			Stroke{ID: "strokeid1", UserID: "user1", Type: 0},
		},
	}

	for _, test := range tests {
		var c Stroke
		assert.NoError(t, UnmarshalMsgContent([]byte(test.msg), &c))
		assert.Equal(t, test.want, c, "incorrect unmarshalling of message content")
	}
}
