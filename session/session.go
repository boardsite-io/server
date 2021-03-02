package session

import (
	"encoding/json"
	"errors"
	"fmt"

	gonanoid "github.com/matoous/go-nanoid/v2"

	"github.com/heat1q/boardsite/api/types"
	"github.com/heat1q/boardsite/redis"
)

// GetStrokes fetches all stroke data for specified page.
func GetStrokes(sessionID, pageID string) ([]types.Stroke, error) {
	strokesRaw, err := redis.FetchStrokesRaw(sessionID, pageID)
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

// GetPages returns all pageIDs in order.
func GetPages(sessionID string) []string {
	return redis.GetPages(sessionID)
}

// GetPagesSet returns all pageIDs in a map for fast verification.
func GetPagesSet(sessionID string) map[string]struct{} {
	pageIDs := GetPages(sessionID)
	pageIDSet := make(map[string]struct{})

	for _, pid := range pageIDs {
		pageIDSet[pid] = struct{}{}
	}

	return pageIDSet
}

// IsValidPage checks if a pageID is valid, i.e. the page exists.
func IsValidPage(sessionID, pageID string) bool {
	_, ok := GetPagesSet(sessionID)[pageID]
	return ok
}

// AddPage adds a page with pageID to the session and broadcasts
// the change to all connected clients.
func AddPage(scb *ControlBlock, pageID string, index int) {
	//TODO handle error
	redis.AddPage(scb.ID, pageID, index)
	UpdatePages(
		scb,
		redis.GetPages(scb.ID),
	)
}

// DeletePage deletes a page with pageID and broadcasts
// the change to all connected clients.
func DeletePage(scb *ControlBlock, pageID string) {
	//TODO handle error
	redis.DeletePage(scb.ID, pageID)
	UpdatePages(
		scb,
		redis.GetPages(scb.ID),
	)
}

// ClearPage clears all strokes on page with pageID and broadcasts
// the change to all connected clients.
func ClearPage(scb *ControlBlock, pageIDs ...string) {
	//TODO handle error
	for _, pid := range pageIDs {
		redis.ClearPage(scb.ID, pid)
	}
	scb.Broadcast <- &types.Message{
		Type:    types.MessageTypePageClear,
		Sender:  "", // send to all clients
		Content: pageIDs,
	}
}

// UpdatePages broadcasts the current PageRank to all connected
// clients indicating an update in the pages (or ordering).
func UpdatePages(scb *ControlBlock, pageIDsToUpdate []string) {
	scb.Broadcast <- &types.Message{
		Type:    types.MessageTypePageSync,
		Sender:  "", // send to all clients
		Content: pageIDsToUpdate,
	}
}

// NewUser generate a new user struct based on
// the alias and color attribute
//
// Does some sanitize checks.
func NewUser(scb *ControlBlock, alias, color string) (*types.User, error) {
	if len(alias) > 24 {
		alias = alias[:24]
	}
	//TODO check if html color ?
	if len(color) != 7 {
		return nil, fmt.Errorf("incorrect html color")
	}

	id, err := gonanoid.New(16)
	if err != nil {
		return nil, err
	}
	user := &types.User{
		ID:    id,
		Alias: alias,
		Color: color,
	}
	// set user waiting
	scb.UserReady(user)
	return user, err
}

// Receive is the entry point when a message is received in
// the session via the websocket.
func Receive(scb *ControlBlock, msg *types.Message) error {
	if !scb.IsUserConnected(msg.Sender) {
		return errors.New("invalid sender userId")
	}

	switch msg.Type {
	case types.MessageTypeStroke:
		return SanitizeStrokes(scb, msg)
	default:
		return fmt.Errorf("message type not recognized: %s", msg.Type)
	}
}

// SanitizeStrokes parses the stroke content of the message.
//
// It further checks if the strokes have a valid pageId and userId.
func SanitizeStrokes(scb *ControlBlock, msg *types.Message) error {
	var strokes []*types.Stroke
	if err := msg.UnmarshalContent(&strokes); err != nil {
		return err
	}

	validStrokes := make([]*types.Stroke, 0, len(strokes))
	pageIDs := GetPagesSet(scb.ID)

	for _, stroke := range strokes {
		if _, ok := pageIDs[stroke.GetPageID()]; ok { // valid pageID
			if stroke.GetUserID() == msg.Sender { // valid userID
				validStrokes = append(validStrokes, stroke)
			}
		}
	}
	UpdateStrokes(scb, msg.Sender, validStrokes)
	return nil
}

// UpdateStrokes updates the strokes in the session with sessionID.
//
// userID indicates the initiator of the message, which is
// to be excluded in the broadcast. The strokes are scheduled for an
// update to Redis.
func UpdateStrokes(scb *ControlBlock, userID string, strokes []*types.Stroke) {
	// broadcast changes
	scb.Broadcast <- &types.Message{
		Type:    types.MessageTypeStroke,
		Sender:  userID,
		Content: strokes,
	}

	// save to database
	scb.DBUpdate <- strokes
}
