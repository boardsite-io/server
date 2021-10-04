package session

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/heat1q/boardsite/api/types"
)

// GetStrokes fetches all stroke data for specified page.
func (scb *ControlBlock) GetStrokes(ctx context.Context, pageID string) ([]types.Stroke, error) {
	strokesRaw, err := scb.cache.FetchStrokesRaw(ctx, scb.ID, pageID)
	if err != nil {
		return nil, errors.New("unable to fetch strokes")
	}

	strokes := make([]types.Stroke, len(strokesRaw))
	for i, s := range strokesRaw {
		if err := json.Unmarshal(s, &strokes[i]); err != nil {
			return nil, err
		}
	}
	return strokes, nil
}

// Receive is the entry point when a message is received in
// the session via the websocket.
func (scb *ControlBlock) Receive(ctx context.Context, msg *types.Message) error {
	if !scb.IsUserConnected(msg.Sender) {
		return errors.New("invalid sender userId")
	}

	var err error
	switch msg.Type {
	case types.MessageTypeStroke:
		err = scb.sanitizeStrokes(ctx, msg)

	case types.MessageTypeMouseMove:
		err = scb.mouseMove(msg)

	default:
		err = fmt.Errorf("message type not recognized: %s", msg.Type)
	}
	return err
}

// sanitizeStrokes parses the stroke content of the message.
//
// It further checks if the strokes have a valid pageId and userId.
func (scb *ControlBlock) sanitizeStrokes(ctx context.Context, msg *types.Message) error {
	var strokes []*types.Stroke
	if err := msg.UnmarshalContent(&strokes); err != nil {
		return err
	}

	validStrokes := make([]*types.Stroke, 0, len(strokes))
	pageIDs := scb.GetPagesSet(ctx)

	for _, stroke := range strokes {
		if _, ok := pageIDs[stroke.GetPageID()]; ok { // valid pageID
			if stroke.GetUserID() == msg.Sender { // valid userID
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
func (scb *ControlBlock) updateStrokes(userID string, strokes []*types.Stroke) {
	// broadcast changes
	scb.broadcast <- &types.Message{
		Type:    types.MessageTypeStroke,
		Sender:  userID,
		Content: strokes,
	}

	// save to database
	scb.dbUpdate <- strokes
}

// mouseMove broadcast mouse move events.
func (scb *ControlBlock) mouseMove(msg *types.Message) error {
	var mouseUpdate types.ContentMouseMove
	if err := msg.UnmarshalContent(&mouseUpdate); err != nil {
		return err
	}
	scb.broadcast <- &types.Message{
		Type:    types.MessageTypeMouseMove,
		Sender:  msg.Sender,
		Content: mouseUpdate,
	}
	return nil
}
