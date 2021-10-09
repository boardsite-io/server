package session

import (
	"context"
	"log"
	"sync"

	"github.com/heat1q/boardsite/api/types"
	"github.com/heat1q/boardsite/attachment"
	"github.com/heat1q/boardsite/redis"
)

// ControlBlock holds the information and channels for sessions
type ControlBlock struct {
	ID string

	Attachments attachment.Handler
	Dispatcher  Dispatcher

	broadcast chan *types.Message
	echo      chan *types.Message

	cache    redis.Handler
	dbUpdate chan []*types.Stroke

	signalClose chan struct{}

	muRdyUsr sync.RWMutex
	// users that have previously been created via POST
	// and have not yet joined the session
	usersReady map[string]*types.User

	muUsr sync.RWMutex
	// Active Client users that are in the session
	// and have an intact WS connection
	users    map[string]*types.User
	maxUsers int
	numUsers int
}

// NewControlBlock creates a new Session ControlBlock with unique ID.
func NewControlBlock(sessionID string, cache redis.Handler, dispatcher Dispatcher, maxUsers int) *ControlBlock {
	scb := &ControlBlock{
		ID:          sessionID,
		Attachments: attachment.NewLocalHandler(sessionID),
		Dispatcher:  dispatcher,
		broadcast:   make(chan *types.Message),
		echo:        make(chan *types.Message),
		cache:       cache,
		dbUpdate:    make(chan []*types.Stroke),
		signalClose: make(chan struct{}),
		usersReady:  make(map[string]*types.User),
		users:       make(map[string]*types.User),
		maxUsers:    maxUsers,
	}

	// start goroutines for broadcasting and saving changes to board
	go scb.broadcastLoop()
	go scb.dbUpdateLoop()

	return scb
}

// Close sends a close signal
func (scb *ControlBlock) Close() {
	scb.signalClose <- struct{}{}
}

// broadcastLoop Broadcasts board updates to all clients
func (scb *ControlBlock) broadcastLoop() {
	for {
		select {
		case data := <-scb.broadcast:
			scb.muUsr.RLock()
			for userID, user := range scb.users { // Send to all connected clients
				// except the origin, i.e. the initiator of message
				if userID != data.Sender {
					if err := user.Conn.WriteJSON(data); err != nil {
						log.Printf("%s :: cannot broadcast to %s: %v",
							scb.ID, user.Conn.RemoteAddr(), err)
					}
				}
			}
			scb.muUsr.RUnlock()
		case data := <-scb.echo:
			// echo message back to origin
			scb.muUsr.RLock()
			if err := scb.users[data.Sender].Conn.WriteJSON(data); err != nil {
				log.Printf("error in broadcastLoop: %v", err)
			}
			scb.muUsr.RUnlock()
		case <-scb.signalClose:
			return
		}
	}
}

// dbUpdateLoop updates database according to given Stroke values
func (scb *ControlBlock) dbUpdateLoop() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	for {
		select {
		case strokes := <-scb.dbUpdate:
			if err := scb.cache.Update(ctx, scb.ID, strokes); err != nil {
				log.Printf("error in dbUpdateLoop: %v", err)
			}
		case <-scb.signalClose:
			_ = scb.cache.ClearSession(ctx, scb.ID)
			return
		}
	}
}
