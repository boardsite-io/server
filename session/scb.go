package session

import (
	"errors"
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

	muRdyUsr sync.RWMutex
	// users that have previously been created via POST
	// and have not yet joined the session
	usersReady map[string]*types.User

	muUsr sync.RWMutex
	// Active Client users that are in the session
	// and have an intact WS connection
	users    map[string]*types.User
	numUsers int
}

// NewControlBlock creates a new Session ControlBlock with unique ID.
func NewControlBlock(sessionID string) *ControlBlock {
	scb := &ControlBlock{
		ID:          sessionID,
		Broadcast:   make(chan *types.Message),
		Echo:        make(chan *types.Message),
		DBUpdate:    make(chan []*types.Stroke),
		SignalClose: make(chan struct{}),
		usersReady:  make(map[string]*types.User),
		users:       make(map[string]*types.User),
	}

	// start goroutines for broadcasting and saving changes to board
	go scb.broadcast()
	go scb.updateDatabase()

	return scb
}

func (scb *ControlBlock) UserReady(u *types.User) {
	scb.muRdyUsr.Lock()
	scb.usersReady[u.ID] = u
	scb.muRdyUsr.Unlock()
}

// GetUserReady returns the user with userID ready to join a session.
func (scb *ControlBlock) GetUserReady(userID string) (*types.User, error) {
	scb.muRdyUsr.RLock()
	defer scb.muRdyUsr.RUnlock()
	u, ok := scb.usersReady[userID]
	if !ok {
		return nil, errors.New("ready user not found")
	}
	return u, nil
}

// IsUserReady checks if the user with userID is ready to join a session.
func (scb *ControlBlock) IsUserReady(userID string) bool {
	_, err := scb.GetUserReady(userID)
	return err == nil
}

// UserConnect adds user from the userReady state to clients.
//
// Broadcast that user has connected to session.
func (scb *ControlBlock) UserConnect(userID string) {
	scb.muRdyUsr.Lock()
	u := scb.usersReady[userID]
	delete(scb.usersReady, userID)
	scb.muRdyUsr.Unlock()

	scb.muUsr.Lock()
	scb.users[userID] = u
	scb.numUsers++
	scb.muUsr.Unlock()

	// broadcast that user has joined
	scb.Broadcast <- &types.Message{
		Type:    types.MessageTypeUserConnected,
		Content: u,
	}
}

// UserDisconnect removes user from clients.
//
// Broadcast that user has disconnected from session.
func (scb *ControlBlock) UserDisconnect(userID string) {
	scb.muUsr.Lock()
	u := scb.users[userID]
	delete(scb.users, u.ID)
	scb.numUsers--
	numCl := scb.numUsers
	scb.muUsr.Unlock()

	// if session is empty after client disconnect
	// the session needs to be set to inactive
	if numCl == 0 {
		Close(scb.ID)
		return
	}

	// broadcast that user has left
	scb.Broadcast <- &types.Message{
		Type:    types.MessageTypeUserDisconnected,
		Content: u,
	}
}

// IsUserConnected checks if the user with userID is an active client in the session.
func (scb *ControlBlock) IsUserConnected(userID string) bool {
	scb.muUsr.RLock()
	defer scb.muUsr.RUnlock()
	_, ok := scb.users[userID]
	return ok
}

// Broadcast Broadcasts board updates to all clients
func (scb *ControlBlock) broadcast() {
	for {
		select {
		case data := <-scb.Broadcast:
			scb.muUsr.RLock()
			for userID, user := range scb.users { // Send to all connected clients
				// except the origin, i.e. the initiator of message
				if userID != data.Sender {
					if err := user.Conn.WriteJSON(data); err != nil {
						log.Printf("%s :: cannot broadcast to %s: %v",
							scb.ID, user.Conn.RemoteAddr(), err)
						continue
					}
				}
			}
			scb.muUsr.RUnlock()
		case data := <-scb.Echo:
			// echo message back to origin
			scb.muUsr.RLock()
			if err := scb.users[data.Sender].Conn.WriteJSON(data); err != nil {
				// log.Println("")
				continue
			}
			scb.muUsr.RUnlock()
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
