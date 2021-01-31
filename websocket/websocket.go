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
func UpgradeProtocol(w http.ResponseWriter, r *http.Request, sessionID string) error {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err == nil {
		initSocket(sessionID, conn)
	}
	return err
}

// For development purpose
func checkOrigin(r *http.Request) bool {
	_ = r
	return true
}

func onClientConnect(sessionID string, conn *gws.Conn) {
	session.AddClient(sessionID, conn)
	log.Println(sessionID + " :: " + conn.RemoteAddr().String() + " connected")
}

func onClientDisconnect(sessionID string, conn *gws.Conn) {
	session.RemoveClient(sessionID, conn.RemoteAddr().String())
	log.Println(sessionID + " :: " + conn.RemoteAddr().String() + " disconnected")
	conn.WriteMessage(gws.TextMessage, []byte("connection closed by host"))
	// close the websocket connection
	conn.Close()
}

// Init starts the websocket
func initSocket(sessionID string, conn *gws.Conn) {
	onClientConnect(sessionID, conn)
	defer onClientDisconnect(sessionID, conn)

	// send all the data to client on connect
	session.SendAllToClient(sessionID, conn.RemoteAddr().String())

	for {
		if _, data, err := conn.ReadMessage(); err == nil {
			// sanitize received data
			if strokes, strokesEncoded, e := session.SanitizeReceived(
				sessionID,
				conn.RemoteAddr().String(),
				data,
			); e == nil {
				// update the session data
				session.Update(sessionID, conn.RemoteAddr().String(), strokes, strokesEncoded)

				log.Printf(sessionID+" :: Data Received from %s: %d stroke(s)\n",
					conn.RemoteAddr().String(),
					len(strokes),
				)
			} else {
				continue // skip if data is corrupted
			}
		} else {
			break // socket closed
		}
	}
}
