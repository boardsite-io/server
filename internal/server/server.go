package server

import (
	"context"
	"fmt"
	"strings"

	"github.com/labstack/echo/v4"
	echomw "github.com/labstack/echo/v4/middleware"

	"github.com/boardsite-io/server/internal/config"
	"github.com/boardsite-io/server/internal/metrics"
	apimw "github.com/boardsite-io/server/internal/middleware"
	"github.com/boardsite-io/server/internal/session"
	sessionHttp "github.com/boardsite-io/server/internal/session/http"
	"github.com/boardsite-io/server/pkg/log"
	libmw "github.com/boardsite-io/server/pkg/middleware"
	"github.com/boardsite-io/server/pkg/redis"
)

type Server struct {
	cfg        *config.Configuration
	echo       *echo.Echo
	metrics    metrics.Handler
	session    sessionHttp.Handler
	dispatcher session.Dispatcher
}

func New(cfg *config.Configuration) *Server {
	return &Server{cfg: cfg}
}

// Serve wraps the main application
func (s *Server) Serve(ctx context.Context) (func() error, func() error) {
	s.echo = echo.New()
	s.echo.HideBanner = true
	s.echo.HTTPErrorHandler = libmw.NewErrorHandler()

	// setup redis cache
	cache, err := redis.New(s.cfg.Cache.Host, s.cfg.Cache.Port)
	if err != nil {
		log.Global().Fatalf("redis pool: %v", err)
	}
	log.Global().Info("Redis connection pool initialized.")

	s.dispatcher = session.NewDispatcher(cache)

	// set up session dispatcher/handler
	s.session = sessionHttp.NewHandler(s.cfg.Session, s.dispatcher)

	s.echo.Use(
		echomw.Recover(),
		echomw.Secure(),
		libmw.CORS(s.cfg.Server.AllowedOrigins),
		libmw.Monitoring())

	if s.cfg.Metrics.Enabled {
		s.metrics = metrics.NewHandler(s.dispatcher)
		s.echo.Use(apimw.Metrics(s.metrics))
	}

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

func (s *Server) setupMetrics() {
	s.metrics = metrics.NewHandler(s.dispatcher)
}
