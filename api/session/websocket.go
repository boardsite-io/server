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

func closeHandler(code int, text string) error  {
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

	fmt.Printf("%s connected.\n", conn.RemoteAddr())

	// send the data to client on connect
	boardData, err := db.FetchAll()
	if err != nil {
		fmt.Println("Cannot retrieve board from database")
		return
	}

	// if board does not exist, create it
	if boardData == nil {
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
		db.Set(dataReceived)
	}

	conn.WriteMessage(websocket.TextMessage, []byte("connection closed by host"))
	fmt.Printf("Connection to %s closed by server\n", conn.RemoteAddr())
}
