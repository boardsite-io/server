package main

import (
	"net/http"

	"github.com/gorilla/mux"

	"boardsite/api/session"
)

func main() {
	router := mux.NewRouter()

	// go session.DatabaseUpdater()
	// go session.Broadcaster()

	router.HandleFunc("/board/create", session.CreateBoard)
	router.HandleFunc("/board/{id}", session.ServeBoard)
	//http.HandleFunc("/api/board", session.ServeBoard)
	//http.HandleFunc("/board/create", board.Create)
	http.Handle("/", router)
	http.ListenAndServe(":8000", nil)
}
