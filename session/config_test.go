package session_test

import (
	"testing"

	"github.com/samber/lo"

	"github.com/heat1q/boardsite/api/config"

	"github.com/stretchr/testify/assert"

	"github.com/heat1q/boardsite/session"
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
			incoming: session.ConfigRequest{ReadOnly: lo.ToPtr[bool](true)},
			want:     session.Config{Session: config.Session{ReadOnly: true}},
		},
		{
			name:     "unset readOnly",
			config:   session.Config{Session: config.Session{ReadOnly: true}},
			incoming: session.ConfigRequest{ReadOnly: lo.ToPtr[bool](false)},
			want:     session.Config{Session: config.Session{ReadOnly: false}},
		},
		{
			name:     "set password",
			config:   session.Config{},
			incoming: session.ConfigRequest{Password: lo.ToPtr[string]("test1234")},
			want:     session.Config{Password: "test1234"},
		},
		{
			name:     "unset password",
			config:   session.Config{Password: "test1234"},
			incoming: session.ConfigRequest{Password: lo.ToPtr[string]("")},
			want:     session.Config{},
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
