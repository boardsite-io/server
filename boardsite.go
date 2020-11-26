package main

import (
	"net/http"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"

	"boardsite/api/session"
)

func main() {
	router := mux.NewRouter()

	router.HandleFunc("/board/create", session.CreateBoard)
	router.HandleFunc("/board/{id}", session.HandleBoardRequest)
	http.ListenAndServe(":8000",
		handlers.CORS(
			handlers.AllowedOrigins([]string{"*"}),
			handlers.AllowedMethods([]string{"GET", "HEAD", "POST", "PUT", "DELETE"}),
		)(router))
}
