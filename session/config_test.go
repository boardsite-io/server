package session_test

import (
	"testing"

	"github.com/heat1q/opt"
	"github.com/stretchr/testify/assert"

	"github.com/boardsite-io/server/api/config"
	"github.com/boardsite-io/server/session"
)

func TestConfig_Update(t *testing.T) {
	tests := []struct {
		name     string
		config   session.Config
		incoming session.ConfigRequest
		want     session.Config
		wantErr  bool
	}{
		{
			name:     "set readOnly",
			config:   session.Config{Session: config.Session{ReadOnly: false}},
			incoming: session.ConfigRequest{ReadOnly: opt.New[bool](true)},
			want:     session.Config{Session: config.Session{ReadOnly: true}},
		},
		{
			name:     "unset readOnly",
			config:   session.Config{Session: config.Session{ReadOnly: true}},
			incoming: session.ConfigRequest{ReadOnly: opt.New[bool](false)},
			want:     session.Config{Session: config.Session{ReadOnly: false}},
		},
		{
			name:     "set password",
			config:   session.Config{},
			incoming: session.ConfigRequest{Password: opt.New[string]("test1234")},
			want:     session.Config{Password: "test1234"},
		},
		{
			name:     "unset password",
			config:   session.Config{Password: "test1234"},
			incoming: session.ConfigRequest{Password: opt.New[string]("")},
			want:     session.Config{},
		},
		{
			name:     "set long password returns error",
			config:   session.Config{Password: "test1234"},
			incoming: session.ConfigRequest{Password: opt.New[string]("AAAABBBBCCCCDDDDAAAABBBBCCCCDDDDAAAABBBBCCCCDDDDAAAABBBBCCCCDDDDAAAA")},
			wantErr:  true,
		},
		{
			name:     "set maxUsers",
			config:   session.Config{},
			incoming: session.ConfigRequest{MaxUsers: opt.New[int](5)},
			want:     session.Config{Session: config.Session{MaxUsers: 5}},
		},
		{
			name:     "unset maxUsers",
			config:   session.Config{Session: config.Session{MaxUsers: 5}},
			incoming: session.ConfigRequest{MaxUsers: opt.New[int](10)},
			want:     session.Config{Session: config.Session{MaxUsers: 10}},
		},
		{
			name:     "set invalid maxUsers returs error",
			config:   session.Config{Session: config.Session{MaxUsers: 5}},
			incoming: session.ConfigRequest{MaxUsers: opt.New[int](-1)},
			wantErr:  true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := tt.config
			err := cfg.Update(&tt.incoming)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, cfg)
			}
		})
	}
}
