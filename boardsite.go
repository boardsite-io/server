package main

import (
	"net/http"

	"boardsite/api/session"
)

func main() {
	go session.DatabaseUpdater()
	go session.Broadcaster()

	http.HandleFunc("/api/board", session.ServeBoard)
	http.ListenAndServe(":8000", nil)
}
