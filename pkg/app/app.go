package app

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"

	"github.com/heat1q/boardsite/pkg/session"
)

// Serve wraps the main application
func Serve(ctx context.Context, port int) (func() error, func() error) {
	router := mux.NewRouter()

	// api routes
	router.HandleFunc("/board/create", session.CreateBoard)
	router.HandleFunc("/board/{id}", session.HandleBoardRequest)

	handl := handlers.CORS(
		handlers.AllowedOrigins([]string{"*"}),
		handlers.AllowedMethods([]string{"GET", "HEAD", "POST", "PUT", "DELETE"}),
	)(router)
	serv := http.Server{Addr: fmt.Sprintf(":%d", port), Handler: handl}
	log.Printf("Starting on port %d\n", port)

	return serv.ListenAndServe, func() error {
		return serv.Shutdown(ctx)
	}
}
