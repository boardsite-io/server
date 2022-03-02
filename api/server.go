package api

import (
	"context"
	"fmt"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"github.com/heat1q/boardsite/api/config"
	"github.com/heat1q/boardsite/api/github"
	"github.com/heat1q/boardsite/api/log"
	"github.com/heat1q/boardsite/api/metrics"
	apimw "github.com/heat1q/boardsite/api/middleware"
	"github.com/heat1q/boardsite/redis"
	"github.com/heat1q/boardsite/session"
)

type Server struct {
	cfg        *config.Configuration
	echo       *echo.Echo
	metrics    metrics.Handler
	session    session.Handler
	dispatcher session.Dispatcher
	github     github.Handler
	validator  github.Validator
}

func NewServer(cfg *config.Configuration) *Server {
	return &Server{cfg: cfg}
}

// Serve wraps the main application
func (s *Server) Serve(ctx context.Context) (func() error, func() error) {
	s.echo = echo.New()
	s.echo.HideBanner = true
	s.echo.HTTPErrorHandler = apimw.NewErrorHandler()

	// setup redis cache
	cache, err := redis.New(s.cfg.Cache.Host, s.cfg.Cache.Port)
	if err != nil {
		log.Global().Fatalf("redis pool: %v", err)
	}
	log.Global().Info("Redis connection pool initialized.")

	s.dispatcher = session.NewDispatcher(cache)

	// set up session dispatcher/handler
	s.session = session.NewHandler(s.cfg.Session, s.dispatcher)

	s.metrics = metrics.NewHandler(s.dispatcher)

	if s.cfg.Github.Enabled {
		githubClient := github.NewClient(&s.cfg.Github, cache)
		s.github = github.NewHandler(s.cfg, cache, githubClient)
		s.validator = github.NewValidator(&s.cfg.Github, cache, githubClient)
	}

	s.echo.Use(
		middleware.Recover(),
		middleware.Secure(),
		apimw.CORS(s.cfg.Server.AllowedOrigins),
		apimw.Monitoring(s.metrics))

	// set routes
	s.setRoutes()

	// configure CORS
	origins := strings.Split(s.cfg.Server.AllowedOrigins, ",")
	log.Global().Infof("CORS: allowed origins: %v", origins)

	return func() error {
			log.Global().Infof("Starting %s@%s listening on :%d\n", s.cfg.App.Name, s.cfg.App.Version, s.cfg.Server.Port)
			return s.echo.Start(fmt.Sprintf(":%d", s.cfg.Server.Port))
		}, func() error {
			_ = cache.ClosePool()
			return s.echo.Shutdown(ctx)
		}
}
