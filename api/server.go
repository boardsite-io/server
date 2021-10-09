package api

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"

	"github.com/heat1q/boardsite/api/config"
	"github.com/heat1q/boardsite/redis"
	"github.com/heat1q/boardsite/session"
)

type Server struct {
	cfg        *config.Configuration
	router     *mux.Router
	dispatcher session.Dispatcher
}

func NewServer() (*Server, error) {
	cfg, err := config.New()
	if err != nil {
		return nil, err
	}

	return &Server{
		cfg:    cfg,
		router: mux.NewRouter(),
	}, nil
}

// Serve wraps the main application
func (s *Server) Serve(ctx context.Context) (func() error, func() error) {
	// setup redis cache
	redisHandler, err := redis.New(s.cfg.Cache.Host, s.cfg.Cache.Port)
	if err != nil {
		log.Fatal(err.Error())
	}
	log.Println("Redis connection pool initialized.")

	// set up session dipatcher/handler
	s.dispatcher = session.NewDispatcher(redisHandler)

	// set routes
	s.setRoutes()

	// configure CORS
	handl := handlers.CORS(
		handlers.AllowedOrigins(
			[]string{
				"https://boardsite.io",  // production
				"http://localhost:3000", // testing
			},
		),
		handlers.AllowedHeaders(
			[]string{
				"Content-Type",
			},
		),
		handlers.AllowedMethods(
			[]string{
				"GET",
				"HEAD",
				"POST",
				"PUT",
				"DELETE",
			},
		),
	)(s.router)
	handl = handlers.ContentTypeHandler(
		handl,
		"text/plain",
		"application/json",
		"image/*",
		"multipart/form-data",
	)

	serv := http.Server{Addr: fmt.Sprintf(":%d", s.cfg.Server.Port), Handler: handl}
	log.Printf("Starting %s@%s listening on :%d\n", s.cfg.App.Name, s.cfg.App.Version, s.cfg.Server.Port)

	return serv.ListenAndServe, func() error {
		redisHandler.ClosePool()
		return serv.Shutdown(ctx)
	}
}
