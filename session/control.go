package session

import (
	"log"
	"sync"

	"github.com/heat1q/boardsite/api/types"
	"github.com/heat1q/boardsite/redis"
)

// ControlBlock holds the information and channels for sessions
type ControlBlock struct {
	ID string

	Broadcast chan *types.Message
	Echo      chan *types.Message

	DBUpdate chan []*types.Stroke

	SignalClose chan struct{}

	// users that have previously been created via POST
	// and have not yet joined the session
	UserReady map[string]*types.User

	// Active Client users that are in the session
	// and have an intact WS connection
	Clients    map[string]*types.User
	NumClients int
	Mu         sync.Mutex
}

// NewControlBlock creates a new Session ControlBlock with unique ID.
func NewControlBlock(sessionID string) *ControlBlock {
	scb := &ControlBlock{
		ID:          sessionID,
		Broadcast:   make(chan *types.Message),
		Echo:        make(chan *types.Message),
		DBUpdate:    make(chan []*types.Stroke),
		SignalClose: make(chan struct{}),
		UserReady:   make(map[string]*types.User),
		Clients:     make(map[string]*types.User),
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
			for userID, user := range scb.Clients { // Send to all connected clients
				// except the origin, i.e. the initiator of message
				if userID != data.Sender {
					if err := user.Conn.WriteJSON(data); err != nil {
						log.Printf("%s :: cannot broadcast to %s: %v",
							scb.ID, user.Conn.RemoteAddr(), err)
						continue
					}
				}
			}
			scb.Mu.Unlock()
		case data := <-scb.Echo:
			// echo message back to origin
			scb.Mu.Lock()
			if err := scb.Clients[data.Sender].Conn.WriteJSON(data); err != nil {
				// log.Println("")
				continue
			}
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
