package session

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/websocket"

	"github.com/heat1q/boardsite/pkg/api"
	"github.com/heat1q/boardsite/pkg/database"
)

type errorStatus struct {
	Error string `json:"error"`
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     checkOrigin,
}

// For development purpose
func checkOrigin(r *http.Request) bool {
	_ = r
	return true
}

func onClientConnect(sessionID string, conn *websocket.Conn) {
	ActiveSession[sessionID].Mu.Lock()

	// add current remote connections to clients
	ActiveSession[sessionID].NumClients++
	ActiveSession[sessionID].Clients[conn.RemoteAddr().String()] = conn

	ActiveSession[sessionID].Mu.Unlock()
	fmt.Println(sessionID + " :: " + conn.RemoteAddr().String() + " connected")
}

func onClientDisconnect(sessionID string, conn *websocket.Conn) {
	ActiveSession[sessionID].Mu.Lock()

	// remove current remote connection from clients
	ActiveSession[sessionID].NumClients--
	delete(ActiveSession[sessionID].Clients, conn.RemoteAddr().String())

	ActiveSession[sessionID].Mu.Unlock()

	fmt.Println(sessionID + " :: " + conn.RemoteAddr().String() + " disconnected")
	conn.WriteMessage(websocket.TextMessage, []byte("connection closed by host"))

	// if session is empty after client disconnect
	// the session needs to be set to inactive
	if ActiveSession[sessionID].NumClients == 0 {
		closeSession(sessionID)
	}

	// close the websocket connection
	conn.Close()
}

// InitWebsocket starts the websocket
func InitWebsocket(sessionID string, conn *websocket.Conn, db *database.RedisDB) {
	boardData, _ := db.FetchAll()

	// send the data to client on connect
	conn.WriteMessage(websocket.TextMessage, []byte(boardData))

	for {
		var stroke []api.Stroke

		if _, data, err := conn.ReadMessage(); err == nil {
			// sanitize received data
			if e := json.Unmarshal(data, &stroke); e != nil {
				continue
			}
			fmt.Printf(sessionID+" :: Data Received from %s: %d stroke(s)\n",
				conn.RemoteAddr().String(),
				len(stroke),
			)
		} else {
			break // socket closed
		}

		if strokeContent, err := json.Marshal(&stroke); err == nil {
			// broadcast board values
			ActiveSession[sessionID].Broadcast <- &BroadcastData{
				Origin:  conn.RemoteAddr().String(),
				Content: strokeContent,
			}

			// save to database
			ActiveSession[sessionID].DBUpdate <- stroke
		}
	}
}
