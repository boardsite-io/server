package session_test

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/boardsite-io/server/internal/config"
	"github.com/boardsite-io/server/internal/session"
	"github.com/boardsite-io/server/internal/session/sessionfakes"
	"github.com/boardsite-io/server/pkg/redis/redisfakes"
	"github.com/boardsite-io/server/pkg/types"
)

func Test_controlBlock_SetConfig(t *testing.T) {
	fakeDispatcher := &sessionfakes.FakeDispatcher{}
	fakeCache := &redisfakes.FakeHandler{}
	fakeBroadcaster := &sessionfakes.FakeBroadcaster{}
	fakeBroadcaster.BroadcastReturns(make(chan types.Message, 100))
	cfg := session.Config{
		ID:     "1234",
		Host:   "beef",
		Secret: "potato",
		Session: config.Session{
			MaxUsers: 10,
			ReadOnly: true,
		},
		Password: "test1234",
	}
	var scb session.Controller
	var err error
	require.NoError(t, err)
	tests := []struct {
		name string
		set  string
		want session.Config
	}{
		{
			name: "new config without overwriting constant values",
			set:  `{"readOnly":false,"maxUsers":42,"id":"test","host":"test","secret":"test","password":"newpassword"}`,
			want: session.Config{
				ID:     "1234",
				Host:   "beef",
				Secret: "potato",
				Session: config.Session{
					MaxUsers: 42,
					ReadOnly: false,
				},
				Password: "newpassword",
			},
		},
		{
			name: "remove password",
			set:  `{"password":""}`,
			want: session.Config{
				ID:     "1234",
				Host:   "beef",
				Secret: "potato",
				Session: config.Session{
					MaxUsers: 10,
					ReadOnly: true,
				},
				Password: "",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			scb, err = session.NewControlBlock(
				cfg,
				session.WithDispatcher(fakeDispatcher),
				session.WithCache(fakeCache),
				session.WithBroadcaster(fakeBroadcaster))
			var cfg session.ConfigRequest
			err := json.Unmarshal([]byte(tt.set), &cfg)
			assert.NoError(t, err)
			err = scb.SetConfig(&cfg)
			assert.NoError(t, err)
			assert.Equal(t, tt.want, scb.Config())
		})
	}
}
