package session

import (
	"errors"
	"fmt"
	"log"

	gws "github.com/gorilla/websocket"
	gonanoid "github.com/matoous/go-nanoid/v2"

	"github.com/heat1q/boardsite/api/types"
	"github.com/heat1q/boardsite/redis"
)

const (
	alphabet = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
)

var (
	// ActiveSession maps the session is to the SessionControl struct
	ActiveSession = make(map[string]*ControlBlock)
)

// Create creates and initializes a new SessionControl struct
func Create() string {
	id, _ := gonanoid.Generate(alphabet, 8)
	scb := NewControlBlock(id)

	// assign to SessionControl struct
	ActiveSession[scb.ID] = scb
	log.Printf("Create Session with ID: %s\n", scb.ID)

	return scb.ID
}

// IsValid checks if session with sessionID exists.
func IsValid(sessionID string) bool {
	return ActiveSession[sessionID] != nil
}

// GetStrokes fetches all stroke data for specified page
// as json stringified array of stroke objects.
func GetStrokes(sessionID, pageID string) string {
	return redis.FetchStrokes(sessionID, pageID)
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
func AddPage(sessionID, pageID string, index int) {
	//TODO handle error
	redis.AddPage(sessionID, pageID, index)
	UpdatePages(
		sessionID,
		redis.GetPages(sessionID),
	)
}

// DeletePage deletes a page with pageID and broadcasts
// the change to all connected clients.
func DeletePage(sessionID, pageID string) {
	//TODO handle error
	redis.DeletePage(sessionID, pageID)
	UpdatePages(
		sessionID,
		redis.GetPages(sessionID),
	)
}

// ClearPage clears all strokes on page with pageID and broadcasts
// the change to all connected clients.
func ClearPage(sessionID string, pageIDs ...string) {
	//TODO handle error
	for _, pid := range pageIDs {
		redis.ClearPage(sessionID, pid)
	}
	ActiveSession[sessionID].Broadcast <- &types.Message{
		Type:    types.MessageTypePageClear,
		Sender:  "", // send to all clients
		Content: pageIDs,
	}
}

// Close closes a session.
func Close(sessionID string) {
	ActiveSession[sessionID].close()
	delete(ActiveSession, sessionID)
}

// NewUser generate a new user struct based on
// the alias and color attribute
//
// Does some sanitize checks.
func NewUser(sessionID, alias, color string) (*types.User, error) {
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
	ActiveSession[sessionID].UserReady[id] = user
	return user, err
}

// IsReadyUser checks if the user with userID is ready to join a session.
func IsReadyUser(sessionID, userID string) bool {
	_, ok := ActiveSession[sessionID].UserReady[userID]
	return ok
}

// IsValidClient checks if the user with userID is an active client in the session.
func IsValidClient(sessionID, userID string) bool {
	_, ok := ActiveSession[sessionID].Clients[userID]
	return ok
}

// AddClient adds the client to the session
// and generates a unique userID.
func AddClient(sessionID, userID string, conn *gws.Conn) {
	ActiveSession[sessionID].Mu.Lock()
	ActiveSession[sessionID].NumClients++
	// add current userid connections to clients
	user := ActiveSession[sessionID].UserReady[userID]
	user.Conn = conn
	ActiveSession[sessionID].Clients[userID] = user

	// user joined, remove from waiting list
	delete(ActiveSession[sessionID].UserReady, userID)

	// broadcast that user has joined
	UpdateConnectedUsers(sessionID)

	ActiveSession[sessionID].Mu.Unlock()
}

// RemoveClient removes the client from the session
func RemoveClient(sessionID, userID string) {
	ActiveSession[sessionID].Mu.Lock()
	// remove current remote connection from clients
	ActiveSession[sessionID].NumClients--
	delete(ActiveSession[sessionID].Clients, userID)
	ActiveSession[sessionID].Mu.Unlock()

	// if session is empty after client disconnect
	// the session needs to be set to inactive
	if ActiveSession[sessionID].NumClients == 0 {
		Close(sessionID)
		return
	}

	// broadcast that user has left
	UpdateConnectedUsers(sessionID)
}

// UpdateConnectedUsers broadcasts the current set of active
// connected users/clients in the session.
//
// The client may use this metadata to update some information
// about the session.
func UpdateConnectedUsers(sessionID string) {
	ActiveSession[sessionID].Broadcast <- &types.Message{
		Type:    types.MessageTypeUserSync,
		Sender:  "",
		Content: ActiveSession[sessionID].Clients,
	}
}

// UpdatePages broadcasts the current PageRank to all connected
// clients indicating an update in the pages (or ordering).
func UpdatePages(sessionID string, pageIDsToUpdate []string) {
	ActiveSession[sessionID].Broadcast <- &types.Message{
		Type:    types.MessageTypePageSync,
		Sender:  "", // send to all clients
		Content: pageIDsToUpdate,
	}
}

// Receive is the entry point when a message is received in
// the session via the websocket.
func Receive(sessionID string, msg *types.Message) error {
	if !IsValidClient(sessionID, msg.Sender) {
		return errors.New("invalid sender userId")
	}

	switch msg.Type {
	case types.MessageTypeStroke:
		return SanitizeStrokes(sessionID, msg)
	default:
		return fmt.Errorf("message type not recognized: %s", msg.Type)
	}
}

// SanitizeStrokes parses the stroke content of the message.
//
// It further checks if the strokes have a valid pageId and userId.
func SanitizeStrokes(sessionID string, msg *types.Message) error {
	var strokes []*types.Stroke
	if err := msg.UnmarshalContent(&strokes); err != nil {
		return err
	}

	validStrokes := make([]*types.Stroke, 0, len(strokes))
	pageIDs := GetPagesSet(sessionID)

	for _, stroke := range strokes {
		if _, ok := pageIDs[stroke.GetPageID()]; ok { // valid pageID
			if stroke.GetUserID() == msg.Sender { // valid userID
				validStrokes = append(validStrokes, stroke)
			}
		}
	}
	UpdateStrokes(sessionID, msg.Sender, validStrokes)
	return nil
}

// UpdateStrokes updates the strokes in the session with sessionID.
//
// userID indicates the initiator of the message, which is
// to be excluded in the broadcast. The strokes are scheduled for an
// update to Redis.
func UpdateStrokes(sessionID, userID string, strokes []*types.Stroke) {
	// broadcast changes
	ActiveSession[sessionID].Broadcast <- &types.Message{
		Type:    types.MessageTypeStroke,
		Sender:  userID,
		Content: strokes,
	}

	// save to database
	ActiveSession[sessionID].DBUpdate <- strokes
}
