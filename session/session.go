package session

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/heat1q/boardsite/api/types"
	"github.com/heat1q/boardsite/redis"
)

// Message type definitions.
const (
	MessageTypeStroke           = "stroke"
	MessageTypeUserConnected    = "userconn"
	MessageTypeUserDisconnected = "userdisc"
	MessageTypePageSync         = "pagesync"
	MessageTypeMouseMove        = "mmove"
)

// ContentMouseMove declares mouse move updates.
type ContentMouseMove struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}

// Style declares the stroke style.
type Style struct {
	Color   string  `json:"color"`
	Width   float64 `json:"width"`
	Opacity float64 `json:"opacity"`
}

// Stroke declares the structure of most stoke types.
type Stroke struct {
	Type   int       `json:"type"`
	ID     string    `json:"id,omitempty"`
	PageID string    `json:"pageId,omitempty"`
	UserID string    `json:"userId"`
	X      float64   `json:"x"`
	Y      float64   `json:"y"`
	ScaleX float64   `json:"scaleX,omitempty"`
	ScaleY float64   `json:"scaleY,omitempty"`
	Points []float64 `json:"points,omitempty"`
	Style  Style     `json:"style,omitempty"`
}

var _ redis.Stroke = (*Stroke)(nil)

// IsDeleted verifies whether stroke is deleted or not
func (s *Stroke) IsDeleted() bool {
	return s.Type == 0
}

// Id returns the id of the stroke
func (s *Stroke) Id() string {
	return s.ID
}

// UserId returns the userid of the stroke
func (s *Stroke) UserId() string {
	return s.UserID
}

// PageId returns the page id of the stroke
func (s *Stroke) PageId() string {
	return s.PageID
}

// GetStrokes fetches all stroke data for specified page.
func (scb *controlBlock) GetStrokes(ctx context.Context, pageID string) ([]*Stroke, error) {
	strokesRaw, err := scb.cache.GetPageStrokes(ctx, scb.id, pageID)
	if err != nil {
		return nil, errors.New("unable to fetch strokes")
	}

	strokes := make([]*Stroke, len(strokesRaw))
	for i, s := range strokesRaw {
		var stroke Stroke
		if err := json.Unmarshal(s, &stroke); err != nil {
			return nil, err
		}
		strokes[i] = &stroke
	}
	return strokes, nil
}

// Receive is the entry point when a message is received in
// the session via the websocket.
func (scb *controlBlock) Receive(ctx context.Context, msg *types.Message) error {
	if !scb.IsUserConnected(msg.Sender) {
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
	scb.broadcast <- &types.Message{
		Type:    MessageTypeStroke,
		Sender:  userID,
		Content: strokes,
	}

	// save to database
	scb.dbUpdate <- strokes
}

// mouseMove broadcast mouse move events.
func (scb *controlBlock) mouseMove(msg *types.Message) error {
	var mouseUpdate ContentMouseMove
	if err := msg.UnmarshalContent(&mouseUpdate); err != nil {
		return err
	}
	scb.broadcast <- &types.Message{
		Type:    MessageTypeMouseMove,
		Sender:  msg.Sender,
		Content: mouseUpdate,
	}
	return nil
}
