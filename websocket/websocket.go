package websocket

import (
	"context"
	"log"

	gws "github.com/gorilla/websocket"

	"github.com/heat1q/boardsite/api/types"
	"github.com/heat1q/boardsite/session"
)

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
	_ = conn.WriteMessage(gws.TextMessage, []byte("connection closed by host"))
	_ = conn.Close()
}

// Subscribe subscribes to the websocket connection
func Subscribe(ctx context.Context, conn *gws.Conn, scb *session.ControlBlock, userID string) error {
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
			if errSanitize := scb.Receive(
				ctx,
				msg,
			); errSanitize != nil {
				log.Println(scb.ID+" :: Error Receive :: %v", err)
				continue // skip if data is corrupted
			}
		} else {
			break // socket closed
		}
	}
	return nil
}
