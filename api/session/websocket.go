package session

import (
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
	fmt.Println(sessionID + "::" + conn.RemoteAddr().String() + " connected")
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
		ActiveSession[sessionID].SetInactive()
	}

	fmt.Println(sessionID + "::" + conn.RemoteAddr().String() + " disconnected")
	conn.WriteMessage(websocket.TextMessage, []byte("connection closed by host"))

	// close the websocket connection
	conn.Close()
}

func closeHandler(code int, text string) error {
	fmt.Printf("Connection closed %d: %s\n", code, text)
	return nil
}

func initBoard(sessionID string) (*database.BoardDB, []board.Position, error) {
	db, err := database.NewConnection(
		sessionID,
		ActiveSession[sessionID].SizeX,
		ActiveSession[sessionID].SizeY,
		ActiveSession[sessionID].NumBytes,
	)
	if err != nil {
		return nil, nil, err
	}

	data, err := db.FetchAll()
	if err != nil {
		return nil, nil, err
	} else if data == nil { // if board does not exist, create it
		db.Reset()
		data = []board.Position{}
	}

	return db, data, nil
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
	conn.WriteJSON(&boardData)

	for {
		var data []board.Position

		if err := conn.ReadJSON(&data); err != nil {
			break
		}

		fmt.Printf("Data Received: %v\n", data)

		// broadcast board values
		ActiveSession[sessionID].Board <- data

		// save to database
		ActiveSession[sessionID].DBCache <- data
	}
}
