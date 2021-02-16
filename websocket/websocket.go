package websocket

import (
	"log"
	"net/http"

	"github.com/heat1q/boardsite/session"

	gws "github.com/gorilla/websocket"
)

type errorStatus struct {
	Error string `json:"error"`
}

var upgrader = gws.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     checkOrigin,
}

// UpgradeProtocol to websocket protocol
func UpgradeProtocol(
	w http.ResponseWriter,
	r *http.Request,
	sessionID, userID string,
) error {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err == nil {
		initSocket(sessionID, userID, conn)
	}
	return err
}

// For development purpose
func checkOrigin(r *http.Request) bool {
	_ = r
	return true
}

func onClientConnect(sessionID, userID string, conn *gws.Conn) {
	session.AddClient(sessionID, userID, conn)
	log.Println(sessionID + " :: " + conn.RemoteAddr().String() + " connected")
}

func onClientDisconnect(sessionID, userID string, conn *gws.Conn) {
	session.RemoveClient(sessionID, userID)
	log.Println(sessionID + " :: " + conn.RemoteAddr().String() + " disconnected")
	conn.WriteMessage(gws.TextMessage, []byte("connection closed by host"))
	// close the websocket connection
	conn.Close()
}

// Init starts the websocket
func initSocket(sessionID, userID string, conn *gws.Conn) {
	onClientConnect(sessionID, userID, conn)
	defer onClientDisconnect(sessionID, userID, conn)

	for {
		if _, data, err := conn.ReadMessage(); err == nil {
			// sanitize received data
			if errSanitize := session.SanitizeAndRelay(
				sessionID,
				conn.RemoteAddr().String(),
				data,
			); errSanitize != nil {
				continue // skip if data is corrupted
			}

			log.Printf(sessionID+" :: Data Received from %s\n",
				conn.RemoteAddr().String(),
			)
		} else {
			break // socket closed
		}
	}
}
