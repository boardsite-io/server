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
	boardUpdate = make(chan []board.Position)
	clients     = make(map[string]*websocket.Conn)
	mu          sync.Mutex
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

func onClientDisconnect(addr string) {
	mu.Lock()
	delete(clients, addr)
	mu.Unlock()
	fmt.Println(addr + " disconnected")
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

// For development purpose
func checkOrigin(r *http.Request) bool {
	_ = r
	return true
}

func closeHandler(code int, text string) error {
	fmt.Printf("Connection closed %d: %s\n", code, text)
	return nil
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
	db, err := database.NewConnection()
	if err != nil {
		fmt.Println("Cannot connect database")
		return
	}
	// close when we are done
	defer db.Close()

	onClientConnect(conn)

	// send the data to client on connect
	boardData, err := db.FetchAll()
	if err != nil {
		fmt.Println("Cannot retrieve board from database")
		return
	} else if boardData == nil { // if board does not exist, create it
		db.Reset()
		conn.WriteJSON(&[]board.Position{})
	} else {
		conn.WriteJSON(&boardData)
	}

	for {
		var dataReceived []board.Position

		if err := conn.ReadJSON(&dataReceived); err != nil {
			break
		}

		fmt.Printf("Data Received: %v\n", dataReceived)

		// broadcast board values
		boardUpdate <- dataReceived

		// save to database
		db.Set(dataReceived)
	}

	conn.WriteMessage(websocket.TextMessage, []byte("connection closed by host"))
	//fmt.Printf("Connection to %s closed by server\n", conn.RemoteAddr())
	onClientDisconnect(conn.RemoteAddr().String())
}
