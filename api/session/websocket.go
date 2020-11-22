package session

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/websocket"

	"boardsite/api/board"
	"boardsite/api/database"
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

	// if session is empty after client disconnect
	// the session needs to be set to inactive
	if ActiveSession[sessionID].NumClients == 0 {
		closeSession(sessionID)
	}

	fmt.Println(sessionID + " :: " + conn.RemoteAddr().String() + " disconnected")
	conn.WriteMessage(websocket.TextMessage, []byte("connection closed by host"))

	// close the websocket connection
	conn.Close()
}

func closeHandler(code int, text string) error {
	fmt.Printf("Connection closed %d: %s\n", code, text)
	return nil
}

func initBoard(sessionID string) (*database.RedisDB, string, error) {
	db, err := database.NewConnection(sessionID)
	if err != nil {
		return nil, "", err
	}

	data, err := db.FetchAll()
	return db, data, err
}

// InitWebsocket starts the websocket
func InitWebsocket(sessionID string, conn *websocket.Conn) {

	conn.SetCloseHandler(closeHandler)

	// connect the database
	db, boardData, err := initBoard(sessionID)
	if err != nil {
		fmt.Println("Cannot connect to database")
		return
	}
	// close when we are done
	defer db.Close()

	// send the data to client on connect
	conn.WriteMessage(websocket.TextMessage, []byte(boardData))

	for {
		var stroke []board.Stroke

		if _, data, err := conn.ReadMessage(); err == nil {
			// sanitize received data
			if e := json.Unmarshal(data, &stroke); e != nil {
				continue
			}
			fmt.Printf(sessionID+" :: Data Received from %s: %v\n", conn.RemoteAddr().String(), stroke)
		} else {
			break // socket closed
		}

		// broadcast board values
		ActiveSession[sessionID].Broadcast <- &board.BroadcastData{
			Origin:  conn.RemoteAddr().String(),
			Content: stroke,
		}

		// save to database
		ActiveSession[sessionID].DBCache <- stroke
	}
}
