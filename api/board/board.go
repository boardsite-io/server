package board

import (
    "fmt"
    "net/http"

    "github.com/gorilla/websocket"
)

type boardPos struct {
    Action string `json:"action"`
    Value  uint8  `json:"value"`
    X      int    `json:"x"`
    Y      int    `json:"y"`
}

type errorStatus struct {
    Error string `json:"error"`
}

var upgrader = websocket.Upgrader{
    ReadBufferSize:  1024,
    WriteBufferSize: 1024,
    CheckOrigin: checkOrigin,
}

// For development purpose
func checkOrigin(r *http.Request) bool {
    _ = r
    return true
}

// ServeBoard starts the websocket
func ServeBoard(w http.ResponseWriter, r *http.Request) {
    conn, err := upgrader.Upgrade(w, r, nil)
    if err != nil {
        fmt.Println(err.Error())
        return
    }

    fmt.Printf("%s connected.\n", conn.RemoteAddr())

    // send the data to client on connect

    for {
        var data []boardPos

        if err := conn.ReadJSON(&data); err != nil {
            errMsg := errorStatus{Error: "JSON unmarshaling failed"}
            conn.WriteJSON(&errMsg)
            continue
        }

        fmt.Println(data)
    }
}
