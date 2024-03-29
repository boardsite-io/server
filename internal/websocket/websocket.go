package websocket

import (
	"context"
	"fmt"
	"net/http"

	gws "github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"

	"github.com/boardsite-io/server/internal/session"
	"github.com/boardsite-io/server/pkg/log"
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
func upgrade(c echo.Context) (*gws.Conn, error) {
	return upgrader.Upgrade(c.Response().Writer, c.Request(), nil)
}

func onClientConnect(ctx context.Context, scb session.Controller, userID string, conn *gws.Conn) error {
	if err := scb.UserConnect(userID, conn); err != nil {
		return err
	}
	log.Ctx(ctx).Infof("session %s :: %s (%s) connected", scb.ID(), userID, conn.RemoteAddr().String())
	return nil
}

func onClientDisconnect(ctx context.Context, scb session.Controller, userID string, conn *gws.Conn) {
	scb.UserDisconnect(ctx, userID)
	log.Ctx(ctx).Infof("session %s :: %s (%s) disconnected", scb.ID(), userID, conn.RemoteAddr().String())
}

// Subscribe subscribes to the websocket connection
func Subscribe(c echo.Context, scb session.Controller, userID string) error {
	ctx := c.Request().Context()
	conn, err := upgrade(c)
	if err != nil {
		return err
	}

	if err := onClientConnect(ctx, scb, userID, conn); err != nil {
		_ = conn.WriteMessage(gws.CloseMessage, gws.FormatCloseMessage(gws.CloseNormalClosure, fmt.Sprintf("%v", err)))
		_ = conn.Close()
		return err
	}
	defer onClientDisconnect(ctx, scb, userID, conn)

	for {
		_, data, err := conn.ReadMessage()
		if err != nil {
			break // socket closed
		}

		msg, err := session.UnmarshalMessage(data)
		if err != nil {
			continue
		}

		// sanitize received data
		if err := scb.Receive(ctx, msg, userID); err != nil {
			log.Ctx(ctx).Warnf("session %s :: error receive message from %s: %v", scb.ID(), msg.Sender, err)
			scb.Broadcaster().Send() <- session.Message{
				Type:     "error",
				Receiver: userID,
				Content:  err.Error(),
			}
		}
	}
	return nil
}
