package app

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"

	"github.com/heat1q/boardsite/api/routes"
)

// Serve wraps the main application
func Serve(ctx context.Context, port int) (func() error, func() error) {
	router := mux.NewRouter()

	// api routes
	router.HandleFunc("/b/create", routes.HandleCreateSession)
	router.HandleFunc("/b/{id}", routes.HandleSessionRequest)
	router.HandleFunc("/b/{id}/pages", routes.HandlePageRequest)
	router.HandleFunc("/b/{id}/pages/{pageId}", routes.HandlePageUpdate)

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
