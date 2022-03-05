package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	readOnly := false
	want := &Configuration{}
	want.App.Name = "boardsite-server"
	want.App.Version = "1.0.0"
	want.Server.BaseURL = "http://localhost"
	want.Server.Port = 8000
	want.Server.AllowedOrigins = "*"
	want.Server.Metrics.Enabled = true
	want.Server.Metrics.Route = "/metrics"
	want.Server.Metrics.User = "admin"
	want.Server.Metrics.Password = "admin"
	want.Server.RPM = 10
	want.Cache.Host = "localhost"
	want.Cache.Port = 6379
	want.Session.MaxUsers = 4
	want.Session.ReadOnly = &readOnly
	want.Github.ClientId = "client-Id"
	want.Github.ClientSecret = "client-Secret"
	want.Github.RedirectURI = "http://localhost:3000"
	want.Github.Scope = []string{"user:email"}
	want.Github.Emails = []string{"potato@boardsite.io"}
	want.Github.WhitelistedEmails = map[string]struct{}{
		"potato@boardsite.io": {},
	}

	got, err := New("./../../config.yaml")

	assert.NoError(t, err)
	assert.Equal(t, want, got)
}
