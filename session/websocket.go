package session

import (
	"context"
	"log"
	"net/http"

	gws "github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"

	"github.com/heat1q/boardsite/api/types"
)

var upgrader = gws.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	// already checked by CORS middleware
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// upgrade upgrade a connection to a websocker connection
func upgrade(c echo.Context, onConnectFn func(conn *gws.Conn) error) error {
	conn, err := upgrader.Upgrade(c.Response().Writer, c.Request(), nil)
	if err != nil {
		return err
	}
	return onConnectFn(conn)
}

func onClientConnect(scb *ControlBlock, userID string, conn *gws.Conn) error {
	u, err := scb.GetUserReady(userID)
	if err != nil {
		return err
	}
	u.Conn = conn      // set the current ws connection
	scb.UserConnect(u) // already checked if user is ready at this point
	log.Println(scb.ID + " :: " + conn.RemoteAddr().String() + " connected")
	return nil
}

func onClientDisconnect(scb *ControlBlock, userID string, conn *gws.Conn) {
	scb.UserDisconnect(userID)
	log.Println(scb.ID + " :: " + conn.RemoteAddr().String() + " disconnected")
	_ = conn.WriteMessage(gws.TextMessage, []byte("connection closed by host"))
	_ = conn.Close()
}

// Subscribe subscribes to the websocket connection
func Subscribe(ctx context.Context, conn *gws.Conn, scb *ControlBlock, userID string) error {
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
