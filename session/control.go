package session

import (
	"log"
	"sync"

	gws "github.com/gorilla/websocket"
	"github.com/heat1q/boardsite/api/types"
	"github.com/heat1q/boardsite/database"
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
	DBClear  chan struct{}
	DBFetch  chan string

	SignalClose chan struct{}

	// Active Connections
	Clients    map[string]*gws.Conn
	NumClients int
	Mu         sync.Mutex
}

// NewControlBlock creates a new Session ControlBlock with unique ID.
func NewControlBlock(sessionID string) *ControlBlock {
	return &ControlBlock{
		ID:          sessionID,
		Broadcast:   make(chan *BroadcastData),
		Echo:        make(chan *BroadcastData),
		DBUpdate:    make(chan []*types.Stroke),
		DBClear:     make(chan struct{}),
		DBFetch:     make(chan string),
		SignalClose: make(chan struct{}),
		Clients:     make(map[string]*gws.Conn),
		NumClients:  0,
	}
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
	db, err := database.NewRedisConn(scb.ID)
	defer db.Close()
	// close session if db connection fails
	if err != nil {
		scb.close()
		log.Fatal("Cannot connect to database")
		return
	}

	for {
		select {
		case board := <-scb.DBUpdate:
			db.Update(board)
		case origin := <-scb.DBFetch:
			data, _ := db.FetchAll()
			scb.Echo <- &BroadcastData{Origin: origin, Content: []byte(data)}
		case <-scb.DBClear:
			db.Clear()
		case <-scb.SignalClose:
			db.Clear()
			return
		}
	}
}

func (scb *ControlBlock) close() {
	scb.SignalClose <- struct{}{}
}

// Clear clears the data in this session
func (scb *ControlBlock) clear() {
	scb.DBClear <- struct{}{}
	scb.Broadcast <- &BroadcastData{Content: []byte("[]")}
}
