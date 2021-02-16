package session

import (
	"encoding/json"
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
		[]string{},
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
		[]string{},
	)
}

// ClearPage clears all strokes on page with pageID and broadcasts
// the change to all connected clients.
func ClearPage(sessionID, pageID string) {
	//TODO handle error
	redis.ClearPage(sessionID, pageID)
	UpdatePages(
		sessionID,
		[]string{},
		[]string{pageID},
	)
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
	stroke := &types.Stroke{
		Type:           -1,
		ConnectedUsers: ActiveSession[sessionID].Clients,
	}

	UpdateStrokes(
		sessionID,
		"", // send to all clients
		[]*types.Stroke{stroke},
		[]*types.Stroke{},
	)
}

// UpdatePages broadcasts the current PageRank to all connected
// clients indicating an update in the pages (or ordering).
func UpdatePages(sessionID string, pageIDsToUpdate, pageIDsToClear []string) {
	stroke := &types.Stroke{
		Type:      -1, // non-zero, since it's no deletion
		PageRank:  pageIDsToUpdate,
		PageClear: pageIDsToClear,
	}

	UpdateStrokes(
		sessionID,
		"", // send to all clients
		[]*types.Stroke{stroke},
		[]*types.Stroke{},
	)
}

// UpdateStrokes updates the data in the session with sessionID.
//
// RemoteAddr indicates the initiator of the message, which is
// to be excluded in the broadcast.
// Strokes in the first slice are broadcasted to all connected
// clients. Stroke in the second slice (those with type >= 0)
// are updated in Redis.
func UpdateStrokes(sessionID, userID string, strokes, strokesDB []*types.Stroke) {
	// ignore error
	// since it is unlikely that marshalling fails
	strokesEncoded, _ := json.Marshal(&strokes)
	// broadcast changes
	ActiveSession[sessionID].Broadcast <- &BroadcastData{
		UserID:  userID,
		Content: strokesEncoded,
	}

	// save to database
	if len(strokesDB) > 0 {
		ActiveSession[sessionID].DBUpdate <- strokesDB
	}
}

// SanitizeAndRelay sanitizes websocket input data and returns an
// error if data is corrupted.
func SanitizeAndRelay(sessionID, remoteAddr string, data []byte) error {
	var strokes = []types.Stroke{}
	if err := json.Unmarshal(data, &strokes); err != nil {
		return err
	}

	strokesSanitized, strokesDB := SanitizeStrokes(sessionID, strokes)
	// update the session data
	UpdateStrokes(sessionID, remoteAddr, strokesSanitized, strokesDB)
	return nil
}

// SanitizeStrokes sanitizes a slice of strokes.
//
// It divides the input slice into two slices of pointer
// to strokes to prevent hardcopies. The first contains all
// sanitizes slices. The second contains only the stroke that
// need to be stored in the DB (i.e. type >= 0).
func SanitizeStrokes(sessionID string, strokes []types.Stroke) ([]*types.Stroke, []*types.Stroke) {
	strokesSanitized := make([]*types.Stroke, 0, len(strokes))
	strokesDB := make([]*types.Stroke, 0, len(strokes))
	pageIDs := GetPagesSet(sessionID)

	for i := range strokes {
		if _, ok := pageIDs[strokes[i].GetPageID()]; ok { // valid pageID
			if IsValidClient(sessionID, strokes[i].GetUserID()) { // valid userID
				strokesSanitized = append(strokesSanitized, &strokes[i])
				if strokes[i].Type >= 0 {
					strokesDB = append(strokesDB, &strokes[i])
				}
			}
		}
	}

	return strokesSanitized, strokesDB
}
