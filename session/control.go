package session

import (
	"sync"

	gws "github.com/gorilla/websocket"
	"github.com/heat1q/boardsite/api/types"
	"github.com/heat1q/boardsite/redis"
)

// BroadcastData holds data to be broadcasted and the origin
type BroadcastData struct {
	Origin  string
	Content []byte
}

// ControlBlock holds the information and channels for sessions
type ControlBlock struct {
	ID string

	Broadcast chan *BroadcastData
	Echo      chan *BroadcastData

	DBUpdate chan []*types.Stroke

	SignalClose chan struct{}

	// Active Connections
	Clients    map[string]*gws.Conn
	NumClients int
	Mu         sync.Mutex
}

// NewControlBlock creates a new Session ControlBlock with unique ID.
func NewControlBlock(sessionID string) *ControlBlock {
	scb := &ControlBlock{
		ID:          sessionID,
		Broadcast:   make(chan *BroadcastData),
		Echo:        make(chan *BroadcastData),
		DBUpdate:    make(chan []*types.Stroke),
		SignalClose: make(chan struct{}),
		Clients:     make(map[string]*gws.Conn),
		NumClients:  0,
	}

	// start goroutines for broadcasting and saving changes to board
	go scb.broadcast()
	go scb.updateDatabase()

	return scb
}

// Broadcast Broadcasts board updates to all clients
func (scb *ControlBlock) broadcast() {
	for {
		select {
		case data := <-scb.Broadcast:
			scb.Mu.Lock()
			for addr, clientConn := range scb.Clients { // Send to all connected clients
				// except the origin, i.e. the initiator of message
				if addr != data.Origin {
					clientConn.WriteMessage(gws.TextMessage, data.Content) // ignore error
				}
			}
			scb.Mu.Unlock()
		case data := <-scb.Echo:
			// echo message back to origin
			scb.Mu.Lock()
			scb.Clients[data.Origin].WriteMessage(gws.TextMessage, data.Content)
			scb.Mu.Unlock()
		case <-scb.SignalClose:
			return
		}
	}
}

// UpdateDatabase Updates database according to given Stroke values
func (scb *ControlBlock) updateDatabase() {
	for {
		select {
		case strokes := <-scb.DBUpdate:
			redis.Update(scb.ID, strokes)
		case <-scb.SignalClose:
			redis.ClearSession(scb.ID)
			return
		}
	}
}

func (scb *ControlBlock) close() {
	scb.SignalClose <- struct{}{}
}
