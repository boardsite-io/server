package session

import (
	"fmt"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"

	"boardsite/api/board"
	"boardsite/api/database"
)

var (
	boardUpdate    = make(chan []board.Position)
	databaseUpdate = make(chan []board.Position)
	clients        = make(map[string]*websocket.Conn)
	mu             sync.Mutex
)

type errorStatus struct {
	Error string `json:"error"`
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     checkOrigin,
}

func onClientConnect(conn *websocket.Conn) {
	mu.Lock()
	clients[conn.RemoteAddr().String()] = conn
	mu.Unlock()
	fmt.Println(conn.RemoteAddr().String() + " connected")
}

func onClientDisconnect(conn *websocket.Conn) {
	mu.Lock()
	delete(clients, conn.RemoteAddr().String())
	mu.Unlock()
	fmt.Println(conn.RemoteAddr().String() + " disconnected")
	conn.WriteMessage(websocket.TextMessage, []byte("connection closed by host"))
}

// Broadcaster Broadcasts board updates to all clients
func Broadcaster() {
	for {
		msg := <-boardUpdate

		mu.Lock()
		for _, clientConn := range clients { // Send to all connected clients
			clientConn.WriteJSON(&msg) // ignore error
		}
		mu.Unlock()
	}
}

// DatabaseUpdater Updates database according to given position values
func DatabaseUpdater() {
	db, err := database.NewConnection()
	if err != nil {
		fmt.Println("Cannot connect database")
		return
	}
	defer db.Close()

	for {
		board := <-databaseUpdate

		if board[0].Action == "clear" {
			db.Reset()
			continue
		}

		db.Set(board)
	}
}

// For development purpose
func checkOrigin(r *http.Request) bool {
	_ = r
	return true
}

func closeHandler(code int, text string) error {
	fmt.Printf("Connection closed %d: %s\n", code, text)
	return nil
}

func initBoard() (*database.BoardDB, []board.Position, error) {
	db, err := database.NewConnection()
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

// ServeBoard starts the websocket
func ServeBoard(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	defer conn.Close()

	conn.SetCloseHandler(closeHandler)

	// connect the database
	db, boardData, err := initBoard()
	if err != nil {
		fmt.Println("Cannot connect to database")
		return
	}
	// close when we are done
	defer db.Close()

	// send the data to client on connect
	conn.WriteJSON(&boardData)

	// on client connect/disconnect
	onClientConnect(conn)
	defer onClientDisconnect(conn)

	for {
		var data []board.Position

		if err := conn.ReadJSON(&data); err != nil {
			break
		}

		fmt.Printf("Data Received: %v\n", data)

		// broadcast board values
		boardUpdate <- data

		// save to database
		databaseUpdate <- data
	}
}
