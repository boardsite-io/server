package session

import (
	"fmt"
	"sync"

	"github.com/gorilla/websocket"

	"github.com/heat1q/boardsite/pkg/api"
	"github.com/heat1q/boardsite/pkg/database"
)

var (
	// ActiveSession maps the session is to the SessionControl struct
	ActiveSession = make(map[string]*SessionControl)
)

// BroadcastData holds data to be broadcasted and the origin
type BroadcastData struct {
	Origin  string
	Content []byte
}

// SessionControl holds the information and channels for sessions
type SessionControl struct {
	SizeX int
	SizeY int

	ID string

	Broadcast chan *BroadcastData
	DBUpdate  chan []api.Stroke
	DBClear   chan struct{}
	Close     chan struct{}

	// Active Connections
	Clients    map[string]*websocket.Conn
	NumClients int
	Mu         sync.Mutex
}

// NewSessionControl creates and initializes a new SessionControl struct
func NewSessionControl(id string, x, y int) *SessionControl {
	scb := &SessionControl{
		ID:         id,
		SizeX:      x,
		SizeY:      y,
		Broadcast:  make(chan *BroadcastData),
		DBUpdate:   make(chan []api.Stroke),
		DBClear:    make(chan struct{}),
		Close:      make(chan struct{}),
		Clients:    make(map[string]*websocket.Conn),
		NumClients: 0,
	}

	// start goroutines for broadcasting and saving changes to board
	go scb.broadcast()
	go scb.updateDatabase()

	fmt.Printf("Create Session with ID: %s\n", id)

	return scb
}

// Broadcast Broadcasts board updates to all clients
func (scb *SessionControl) broadcast() {
	select {
	case data := <-scb.Broadcast:
		scb.Mu.Lock()
		for addr, clientConn := range scb.Clients { // Send to all connected clients
			// except the origin, i.e. the initiator of message
			if addr != data.Origin {
				clientConn.WriteMessage(websocket.TextMessage, data.Content) // ignore error
			}
		}
		scb.Mu.Unlock()
	case <-scb.Close:
		return
	}
}

// UpdateDatabase Updates database according to given Stroke values
func (scb *SessionControl) updateDatabase() {
	db, _ := database.NewRedisConn(scb.ID)
	defer db.Close()

	select {
	case board := <-scb.DBUpdate:
		db.Update(board)
	case <-scb.DBClear:
		db.Clear()
	case <-scb.Close:
		db.Clear()
		return
	}
}

// Clear clears the data in this session
func (scb *SessionControl) Clear() {
	scb.DBClear <- struct{}{}
	scb.Broadcast <- &BroadcastData{Content: []byte("[]")}
}

func closeSession(sessionID string) {
	ActiveSession[sessionID].Close <- struct{}{}
}
