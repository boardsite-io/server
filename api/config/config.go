package config

import (
	"github.com/spf13/viper"
)

const (
	name    = "boardsite-server"
	version = "0.1.0"

	serverHost        = "B_HOST"
	defaultServerHost = "localhost"
	serverPort        = "B_PORT"
	defaultServerPort = "8000"

	allowedOrigins = "B_CORS_ORIGINS"
	// comma separated list of allowed origins
	defaultOrigins = "*"

	cacheHost        = "B_REDIS_HOST"
	defaultCacheHost = "localhost"
	cachePort        = "B_REDIS_PORT"
	defaultCachePort = "6379"

	metricsRoute           = "B_METRICS_ROUTE"
	defaultMetricsRoute    = "/metrics"
	metricsUser            = "B_METRICS_USER"
	defaultMetricsUser     = "admin"
	metricsPassword        = "B_METRICS_USER"
	defaultMetricsPassword = "admin"

	// max number of users allowed in one session
	sessionMaxUsers = 10
)

type Configuration struct {
	App struct {
		Name    string
		Version string
	}

	Server struct {
		Host           string
		Port           uint16
		AllowedOrigins string
		Metrics        struct {
			Route    string
			User     string
			Password string
		}
	}

	Cache struct {
		Host string
		Port uint16
	}

	Session struct {
		MaxUsers int
	}
}

func New() (*Configuration, error) {
	cfg := &Configuration{}

	viper.AutomaticEnv()
	set("app.name", "", name)
	set("app.version", "", version)

	set("server.host", serverHost, defaultServerHost)
	set("server.port", serverPort, defaultServerPort)
	set("server.allowedOrigins", allowedOrigins, defaultOrigins)
	set("server.metrics.route", metricsRoute, defaultMetricsRoute)
	set("server.metrics.user", metricsUser, defaultMetricsUser)
	set("server.metrics.password", metricsPassword, defaultMetricsPassword)

	set("cache.host", cacheHost, defaultCacheHost)
	set("cache.port", cachePort, defaultCachePort)

	viper.Set("session.maxUsers", sessionMaxUsers)

	if err := viper.Unmarshal(cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}

func set(key string, envKey string, defaultVal interface{}) {
	viper.Set(key, defaultVal)
	if envKey != "" && viper.IsSet(envKey) {
		viper.Set(key, viper.GetString(envKey))
	}
}
