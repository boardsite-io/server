package session

import (
	"encoding/json"
	"log"
	"math/rand"
	"time"

	gws "github.com/gorilla/websocket"

	"github.com/heat1q/boardsite/api/types"
	"github.com/heat1q/boardsite/database"
)

const (
	letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
)

var (
	// ActiveSession maps the session is to the SessionControl struct
	ActiveSession = make(map[string]*ControlBlock)
)

// GenerateID generates a unique sessionID.
func GenerateID() string {
	rand.Seed(time.Now().UnixNano())
	id := make([]byte, 6)
	// find available id
	for {
		for i := range id {
			id[i] = letters[rand.Intn(len(letters))]
		}

		if ActiveSession[string(id)] == nil {
			break
		}
	}

	return string(id)
}

// Create creates and initializes a new SessionControl struct
func Create() string {
	scb := NewControlBlock(GenerateID())

	// assign to SessionControl struct
	ActiveSession[scb.ID] = scb

	// start goroutines for broadcasting and saving changes to board
	go scb.broadcast()
	go scb.updateDatabase()

	log.Printf("Create Session with ID: %s\n", scb.ID)

	return scb.ID
}

// IsValid checks if session with sessionID exists.
func IsValid(sessionID string) bool {
	return ActiveSession[sessionID] != nil
}

// GetPages returns all pageIDs in order.
func GetPages(sessionID string) []string {
	db, err := database.NewRedisConn(sessionID)
	defer db.Close()

	if err == nil {
		return []string{}
	}

	return db.GetPages()
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

// AddPage adds a page with pageID to the session.
func AddPage(sessionID, pageID string, index int) {
	// TODO
}

// Close closes a session.
func Close(sessionID string) {
	ActiveSession[sessionID].close()
	delete(ActiveSession, sessionID)
}

// AddClient adds the client to the session
func AddClient(sessionID string, conn *gws.Conn) {
	ActiveSession[sessionID].Mu.Lock()

	// add current remote connections to clients
	ActiveSession[sessionID].NumClients++
	ActiveSession[sessionID].Clients[conn.RemoteAddr().String()] = conn

	ActiveSession[sessionID].Mu.Unlock()
}

// RemoveClient removes the client from the session
func RemoveClient(sessionID, remoteAddr string) {
	ActiveSession[sessionID].Mu.Lock()

	// remove current remote connection from clients
	ActiveSession[sessionID].NumClients--
	delete(ActiveSession[sessionID].Clients, remoteAddr)

	// if session is empty after client disconnect
	// the session needs to be set to inactive
	if ActiveSession[sessionID].NumClients == 0 {
		Close(sessionID)
	}

	ActiveSession[sessionID].Mu.Unlock()
}

// SendAllToClient schedules a request to send all data of a session to client.
func SendAllToClient(sessionID, remoteAddr string) {
	ActiveSession[sessionID].DBFetch <- remoteAddr
}

// Update updates the session data by scheduling a broadcast to all connected clients
// and a store request to Redis.
func Update(sessionID, remoteAddr string, strokes []*types.Stroke, strokesEncoded []byte) {
	// broadcast changes
	ActiveSession[sessionID].Broadcast <- &BroadcastData{
		Origin:  remoteAddr,
		Content: strokesEncoded,
	}

	// save to database
	ActiveSession[sessionID].DBUpdate <- strokes
}

// SanitizeReceived sanitizes websocket input data.
// Returns a slice with stroke structs and the respective sanitized JSON encoding.
func SanitizeReceived(sessionID, remoteAddr string, data []byte) ([]*types.Stroke, []byte, error) {
	var strokes = []types.Stroke{}
	if err := json.Unmarshal(data, &strokes); err != nil {
		return nil, nil, err
	}

	strokesSanitized := Sanitize(sessionID, strokes)

	// ignore error
	// since it is unlikely that marshalling fails
	strokesEncoded, _ := json.Marshal(&strokesSanitized)

	return strokesSanitized, strokesEncoded, nil
}

// Sanitize slice of strokes. Return a slice of sanitized strokes consisting
// of pointers to the strokes to prevent a hardcopy.
func Sanitize(sessionID string, strokes []types.Stroke) []*types.Stroke {
	strokesSanitized := make([]*types.Stroke, 0, len(strokes))
	pageIDs := GetPagesSet(sessionID)

	for i := range strokes {
		if _, ok := pageIDs[strokes[i].GetPageID()]; ok {
			strokesSanitized = append(strokesSanitized, &strokes[i])
		}
	}

	return strokesSanitized
}
