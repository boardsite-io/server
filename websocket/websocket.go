package websocket

import (
	"log"
	"net/http"

	"github.com/heat1q/boardsite/api/types"
	"github.com/heat1q/boardsite/session"

	gws "github.com/gorilla/websocket"
)

var upgrader = gws.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     checkOrigin,
}

// UpgradeProtocol to websocket protocol
func UpgradeProtocol(
	w http.ResponseWriter,
	r *http.Request,
	scb *session.ControlBlock,
	userID string,
) error {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return err
	}
	return initSocket(scb, userID, conn)
}

// For development purpose
func checkOrigin(r *http.Request) bool {
	_ = r
	return true
}

func onClientConnect(scb *session.ControlBlock, userID string, conn *gws.Conn) error {
	u, err := scb.GetUserReady(userID)
	if err != nil {
		return err
	}
	u.Conn = conn      // set the current ws connection
	scb.UserConnect(u) // already checked if user is ready at this point
	log.Println(scb.ID + " :: " + conn.RemoteAddr().String() + " connected")
	return nil
}

func onClientDisconnect(scb *session.ControlBlock, userID string, conn *gws.Conn) {
	scb.UserDisconnect(userID)
	log.Println(scb.ID + " :: " + conn.RemoteAddr().String() + " disconnected")
	conn.WriteMessage(gws.TextMessage, []byte("connection closed by host"))
	// close the websocket connection
	conn.Close()
}

// initSocket starts the websocket
func initSocket(scb *session.ControlBlock, userID string, conn *gws.Conn) error {
	if err := onClientConnect(scb, userID, conn); err != nil {
		return err
	}
	defer onClientDisconnect(scb, userID, conn)

	for {
		if _, data, err := conn.ReadMessage(); err == nil {
			msg, errMsg := types.UnmarshalMessage(data)
			if errMsg != nil {
				continue
			}

			// sanitize received data
			if errSanitize := session.Receive(
				scb,
				msg,
			); errSanitize != nil {
				log.Println(scb.ID+" :: Error Receive :: %v", err)
				continue // skip if data is corrupted
			}

			log.Printf(scb.ID+" :: Data Received from %s\n",
				conn.RemoteAddr().String(),
			)
		} else {
			break // socket closed
		}
	}
	return nil
}
