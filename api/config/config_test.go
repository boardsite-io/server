package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	want := &Configuration{}
	want.App.Name = "boardsite-server"
	want.App.Version = "1.0.0"
	want.Server.Host = "localhost"
	want.Server.Port = 8000
	want.Server.AllowedOrigins = "*"
	want.Server.Metrics.Enabled = true
	want.Server.Metrics.Route = "/metrics"
	want.Server.Metrics.User = "admin"
	want.Server.Metrics.Password = "admin"
	want.Cache.Host = "localhost"
	want.Cache.Port = 6379
	want.Session.MaxUsers = 10
	want.Session.RPM = 10

	got, err := New("./../../config.yaml")

	assert.NoError(t, err)
	assert.Equal(t, want, got)
}
