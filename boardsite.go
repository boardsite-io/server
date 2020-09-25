package main

import (
	"net/http"

	"boardsite/api/session"
)

func main() {
	http.HandleFunc("/api/board", session.ServeBoard)
	http.ListenAndServe(":8000", nil)
}
