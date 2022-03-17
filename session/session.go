package session

import (
	"context"
	"errors"
	"fmt"

	"github.com/heat1q/boardsite/api/types"
	"github.com/heat1q/boardsite/redis"
)

// Message type definitions.
const (
	MessageTypeSessionConfig    = "config"
	MessageTypeStroke           = "stroke"
	MessageTypeUserHost         = "userhost"
	MessageTypeUserConnected    = "userconn"
	MessageTypeUserSync         = "usersync"
	MessageTypeUserDisconnected = "userdisc"
	MessageTypeUserKick         = "userkick"
	MessageTypePageSync         = "pagesync"
	MessageTypeMouseMove        = "mmove"
)

// ContentMouseMove declares mouse move updates.
type ContentMouseMove struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}

// Receive is the entry point when a message is received in
// the session via the websocket.
func (scb *controlBlock) Receive(ctx context.Context, msg *types.Message) error {
	if !scb.isUserConnected(msg.Sender) {
		return errors.New("invalid sender userId")
	}

	var err error
	switch msg.Type {
	case MessageTypeStroke:
		err = scb.sanitizeStrokes(ctx, msg)

	case MessageTypeMouseMove:
		err = scb.mouseMove(msg)

	default:
		err = fmt.Errorf("message type not recognized: %s", msg.Type)
	}
	return err
}

// sanitizeStrokes parses the stroke content of the message.
//
// It further checks if the strokes have a valid pageId and userId.
func (scb *controlBlock) sanitizeStrokes(ctx context.Context, msg *types.Message) error {
	var strokes []*Stroke
	if err := msg.UnmarshalContent(&strokes); err != nil {
		return err
	}

	validStrokes := make([]redis.Stroke, 0, len(strokes))
	pageIDs := scb.getPagesSet(ctx)

	for _, stroke := range strokes {
		if _, ok := pageIDs[stroke.PageId()]; ok { // valid pageID
			if stroke.UserId() == msg.Sender { // valid userID
				validStrokes = append(validStrokes, stroke)
			}
		}
	}
	if len(validStrokes) > 0 {
		scb.updateStrokes(msg.Sender, validStrokes)
		return nil
	}
	return errors.New("strokes not validated")
}

// updateStrokes updates the strokes in the session with sessionID.
//
// userID indicates the initiator of the message, which is
// to be excluded in the broadcast. The strokes are scheduled for an
// update to Redis.
func (scb *controlBlock) updateStrokes(userID string, strokes []redis.Stroke) {
	// broadcast changes
	scb.broadcaster.Broadcast() <- types.Message{
		Type:    MessageTypeStroke,
		Sender:  userID,
		Content: strokes,
	}

	// save to database
	scb.broadcaster.Cache() <- strokes
}

// mouseMove broadcast mouse move events.
func (scb *controlBlock) mouseMove(msg *types.Message) error {
	var mouseUpdate ContentMouseMove
	if err := msg.UnmarshalContent(&mouseUpdate); err != nil {
		return err
	}
	scb.broadcaster.Broadcast() <- types.Message{
		Type:    MessageTypeMouseMove,
		Sender:  msg.Sender,
		Content: mouseUpdate,
	}
	return nil
}
