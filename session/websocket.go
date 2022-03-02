package session

import (
	"context"
	"encoding/json"
	"net/http"

	gws "github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"

	"github.com/heat1q/boardsite/api/log"
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

// upgrade upgrade a connection to a websocket connection
func upgrade(c echo.Context, onConnectFn func(conn *gws.Conn) error) error {
	conn, err := upgrader.Upgrade(c.Response().Writer, c.Request(), nil)
	if err != nil {
		return err
	}
	return onConnectFn(conn)
}

func onClientConnect(ctx context.Context, scb Controller, userID string, conn *gws.Conn) error {
	if err := scb.UserConnect(userID, conn); err != nil {
		return err
	}
	log.Ctx(ctx).Infof("session %s :: %s (%s) connected", scb.ID(), userID, conn.RemoteAddr().String())
	return nil
}

func onClientDisconnect(ctx context.Context, scb Controller, userID string, conn *gws.Conn) {
	scb.UserDisconnect(ctx, userID)
	log.Ctx(ctx).Infof("session %s :: %s (%s) disconnected", scb.ID(), userID, conn.RemoteAddr().String())
	// _ = conn.WriteMessage(gws.TextMessage, []byte("connection closed by host"))
	_ = conn.Close()
}

// Subscribe subscribes to the websocket connection
func Subscribe(ctx context.Context, conn *gws.Conn, scb Controller, userID string) error {
	if err := onClientConnect(ctx, scb, userID, conn); err != nil {
		writeError(scb, userID, err)
		_ = conn.Close()
		return err
	}
	defer onClientDisconnect(ctx, scb, userID, conn)

	for {
		_, data, err := conn.ReadMessage()
		if err != nil {
			break // socket closed
		}

		msg, err := types.UnmarshalMessage(data)
		if err != nil {
			continue
		}

		// sanitize received data
		if err := scb.Receive(ctx, msg); err != nil {
			log.Ctx(ctx).Warnf("session %s :: error receive message from %s: %v", scb.ID(), msg.Sender, err)
			writeError(scb, userID, err)
		}
	}
	return nil
}

func writeError(scb Controller, userID string, content interface{}) {
	payload, _ := json.Marshal(types.Message{
		Type:    "error",
		Content: content,
	})
	scb.Broadcaster().Send() <- types.Message{
		Receiver: userID,
		Content:  payload,
	}
}
