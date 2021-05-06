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
	//router.Use(contentTypeMiddleware)

	routes.Set(router)

	handl := handlers.CORS(
		handlers.AllowedOrigins([]string{"*"}),
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
	)(router)
	handl = handlers.ContentTypeHandler(
		handl,
		"text/plain",
		"application/json",
		"image/*",
		"multipart/form-data",
	)

	serv := http.Server{Addr: fmt.Sprintf(":%d", port), Handler: handl}
	log.Printf("Starting on port %d\n", port)

	return serv.ListenAndServe, func() error {
		return serv.Shutdown(ctx)
	}
}

// func contentTypeMiddleware(next http.Handler) http.Handler {
// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		w.Header().Add("Content-Type", "application/json")
// 		next.ServeHTTP(w, r)
// 	})
// }
