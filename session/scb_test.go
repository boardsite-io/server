package session_test

import (
	"encoding/json"
	"testing"

	"github.com/heat1q/boardsite/api/types"

	"github.com/heat1q/boardsite/redis/redisfakes"

	"github.com/stretchr/testify/require"

	"github.com/heat1q/boardsite/session/sessionfakes"

	"github.com/heat1q/boardsite/api/config"
	"github.com/stretchr/testify/assert"

	"github.com/heat1q/boardsite/session"
)

func Test_controlBlock_SetConfig(t *testing.T) {
	fakeDispatcher := &sessionfakes.FakeDispatcher{}
	fakeCache := &redisfakes.FakeHandler{}
	fakeBroadcaster := &sessionfakes.FakeBroadcaster{}
	fakeBroadcaster.BroadcastReturns(make(chan types.Message, 100))
	readonly := true
	wantReadonly := false
	password := "test1234"
	wantPassword := "newpassword"
	cfg := session.Config{
		ID:     "1234",
		Host:   "beef",
		Secret: "potato",
		Session: config.Session{
			MaxUsers: 10,
			ReadOnly: &readonly,
		},
		Password: &password,
	}
	var scb session.Controller
	var err error
	scb, err = session.NewControlBlock(
		cfg,
		session.WithDispatcher(fakeDispatcher),
		session.WithCache(fakeCache),
		session.WithBroadcaster(fakeBroadcaster))
	require.NoError(t, err)
	tests := []struct {
		name string
		set  string
		want session.Config
	}{
		{
			name: "",
			set:  `{"readOnly":false,"maxUsers":42,"id":"test","host":"test","secret":"test","password":"newpassword"}`,
			want: session.Config{
				ID:     "1234",
				Host:   "beef",
				Secret: "potato",
				Session: config.Session{
					MaxUsers: 42,
					ReadOnly: &wantReadonly,
				},
				Password: &wantPassword,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var cfg session.Config
			err := json.Unmarshal([]byte(tt.set), &cfg)
			assert.NoError(t, err)
			err = scb.SetConfig(cfg)
			assert.NoError(t, err)
			assert.Equal(t, tt.want, scb.Config())
		})
	}
}
