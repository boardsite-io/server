package main

import (
    "net/http"

    "boardsite/api/board"
)

func main() {
    http.HandleFunc("/api/board", board.ServeBoard)
    http.ListenAndServe(":8000", nil)
}
