package board

import (
	"fmt"
	"sync"

	"github.com/gorilla/websocket"
)

// Stroke Holds the Stroke and value of pixels
type Stroke struct {
	ID        string    `json:"id"`
	Type      string    `json:"type"`
	Color     string    `json:"color"`
	LineWidth float64   `json:"line_width"`
	Position  []float64 `json:"position"`
}

// BroadcastData holds data to be broadcasted and the origin
type BroadcastData struct {
	Origin  string
	Content []byte
}

// SetupForm Form to setup a new board
type SetupForm struct {
	X int `json:"x"`
	Y int `json:"y"`
}

// SessionControl holds the information and channels for sessions
type SessionControl struct {
	SizeX int
	SizeY int

	ID string

	Broadcast chan *BroadcastData
	DBCache   chan []Stroke
	DB        DatabaseUpdater
	IsActive  bool
	// Active Connections
	Clients    map[string]*websocket.Conn
	NumClients int
	Mu         sync.Mutex
}

// DatabaseUpdater Declares a set of functions used for Database updates.
type DatabaseUpdater interface {
	Delete(id string) error
	Set(value []Stroke) error
	Close()
	Clear() error
}

// NewSessionControl creates and initializes a new SessionControl struct
func NewSessionControl(id string, x, y int, db DatabaseUpdater) *SessionControl {
	scb := &SessionControl{
		ID:         id,
		SizeX:      x,
		SizeY:      y,
		IsActive:   true,
		Broadcast:  make(chan *BroadcastData),
		DBCache:    make(chan []Stroke),
		Clients:    make(map[string]*websocket.Conn),
		NumClients: 0,
		DB:         db,
	}

	// start goroutines for broadcasting and saving changes to board
	go scb.broadcast()
	go scb.updateDatabase()

	fmt.Printf("Create Session with ID: %s\n", id)

	return scb
}

// Broadcast Broadcasts board updates to all clients
func (scb *SessionControl) broadcast() {
	for scb.IsActive {
		data := <-scb.Broadcast

		scb.Mu.Lock()
		for addr, clientConn := range scb.Clients { // Send to all connected clients
			// except the origin, i.e. the initiator of message
			if addr != data.Origin {
				clientConn.WriteMessage(websocket.TextMessage, data.Content) // ignore error
			}
		}
		scb.Mu.Unlock()
	}
}

// UpdateDatabase Updates database according to given Stroke values
func (scb *SessionControl) updateDatabase() {
	for scb.IsActive {
		board := <-scb.DBCache

		if board[0].Type == "clear" {
			scb.DB.Clear()
			continue
		}

		scb.DB.Set(board)
	}
	scb.DB.Close()
}
