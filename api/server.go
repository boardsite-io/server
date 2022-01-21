package api

import (
	"context"
	"fmt"
	"strings"

	"github.com/labstack/echo-contrib/prometheus"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"github.com/heat1q/boardsite/api/config"
	"github.com/heat1q/boardsite/api/log"
	apimw "github.com/heat1q/boardsite/api/middleware"
	"github.com/heat1q/boardsite/redis"
	"github.com/heat1q/boardsite/session"
)

type Server struct {
	cfg        *config.Configuration
	echo       *echo.Echo
	prom       *prometheus.Prometheus
	session    session.Handler
	dispatcher session.Dispatcher
}

func NewServer() (*Server, error) {
	cfg, err := config.New()
	if err != nil {
		return nil, err
	}

	return &Server{
		cfg: cfg,
	}, nil
}

// Serve wraps the main application
func (s *Server) Serve(ctx context.Context) (func() error, func() error) {
	s.echo = echo.New()
	s.echo.HideBanner = true
	s.echo.HTTPErrorHandler = apimw.NewErrorHandler()

	s.prom = prometheus.NewPrometheus("echo", nil)
	s.echo.Use(
		middleware.Recover(),
		middleware.Secure(),
		apimw.CORS(s.cfg.Server.AllowedOrigins),
		apimw.Monitoring(s.prom))

	// setup redis cache
	redisHandler, err := redis.New(s.cfg.Cache.Host, s.cfg.Cache.Port)
	if err != nil {
		s.echo.Logger.Fatalf("redis pool: %v", err)
	}
	log.Global().Info("Redis connection pool initialized.")

	s.dispatcher = session.NewDispatcher(redisHandler)

	// set up session dispatcher/handler
	s.session = session.NewHandler(s.cfg, s.dispatcher)

	// set routes
	s.setRoutes()

	// configure CORS
	origins := strings.Split(s.cfg.Server.AllowedOrigins, ",")
	log.Global().Infof("CORS: allowed origins: %v", origins)

	return func() error {
			log.Global().Infof("Starting %s@%s listening on :%d\n", s.cfg.App.Name, s.cfg.App.Version, s.cfg.Server.Port)
			return s.echo.Start(fmt.Sprintf(":%d", s.cfg.Server.Port))
		}, func() error {
			_ = redisHandler.ClosePool()
			return s.echo.Shutdown(ctx)
		}
}
