package session_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/boardsite-io/server/internal/config"
	"github.com/boardsite-io/server/internal/session"
	"github.com/boardsite-io/server/internal/session/sessionfakes"
	"github.com/boardsite-io/server/pkg/redis/redisfakes"
)

func Test_controlBlock_NewUser(t *testing.T) {
	var scb session.Controller
	dispatcher := &sessionfakes.FakeDispatcher{}
	cache := &redisfakes.FakeHandler{}
	broadcaster := &sessionfakes.FakeBroadcaster{}
	broadcaster.BroadcastReturns(make(chan session.Message, 999))

	tests := []struct {
		name    string
		cfg     session.Config
		userReq session.UserRequest
		want    session.User
		wantErr bool
	}{
		{
			name: "happy path",
			cfg: session.Config{
				Password: "password",
				Session:  config.Session{MaxUsers: 10},
			},
			userReq: session.UserRequest{
				Password: "password",
				User:     session.User{Alias: "potato", Color: "#00ff00"},
			},
			want: session.User{Alias: "potato", Color: "#00ff00"},
		},
		{
			name: "without password",
			cfg: session.Config{
				Session: config.Session{MaxUsers: 10},
			},
			userReq: session.UserRequest{
				Password: "password",
				User:     session.User{Alias: "potato", Color: "#00ff00"},
			},
			want: session.User{Alias: "potato", Color: "#00ff00"},
		},
		{
			name: "wrong password",
			cfg: session.Config{
				Password: "password",
				Session:  config.Session{MaxUsers: 10},
			},
			userReq: session.UserRequest{
				Password: "test",
				User:     session.User{Alias: "potato", Color: "#00ff00"},
			},
			wantErr: true,
		},
		{
			name: "max number of users reached",
			cfg: session.Config{
				Session: config.Session{MaxUsers: 0},
			},
			userReq: session.UserRequest{
				User: session.User{Alias: "potato", Color: "#00ff00"},
			},
			wantErr: true,
		},
		{
			name: "invalid user alias",
			cfg: session.Config{
				Session: config.Session{MaxUsers: 10},
			},
			userReq: session.UserRequest{
				User: session.User{Alias: "", Color: "#00ff00"},
			},
			wantErr: true,
		},
		{
			name: "invalid user color",
			cfg: session.Config{
				Session: config.Session{MaxUsers: 10},
			},
			userReq: session.UserRequest{
				User: session.User{Alias: "potato", Color: "#xxff00"},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			scb, _ = session.NewControlBlock(tt.cfg,
				session.WithDispatcher(dispatcher),
				session.WithCache(cache),
				session.WithBroadcaster(broadcaster))
			got, err := scb.NewUser(tt.userReq)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.NotEmpty(t, got.ID)
				assert.Equal(t, tt.want.Alias, got.Alias)
				assert.Equal(t, tt.want.Color, got.Color)
			}
		})
	}
}
