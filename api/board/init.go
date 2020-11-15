package board

import (
	"fmt"
	"sync"

	"github.com/gorilla/websocket"
)

// Position Holds the position and value of pixels
type Position struct {
	Action string `json:"action"`
	Value  uint32 `json:"color"`
	X      int    `json:"x"`
	Y      int    `json:"y"`
}

// SetupForm Form to setup a new board
type SetupForm struct {
	X int `json:"x"`
	Y int `json:"y"`
}

// SessionControl holds the information and channels for sessions
type SessionControl struct {
	SizeX    int
	SizeY    int
	NumBytes int

	ID string

	Board    chan []Position
	DBCache  chan []Position
	DB       DatabaseUpdater
	IsActive bool
	// Active Connections
	Clients    map[string]*websocket.Conn
	NumClients int
	Mu         sync.Mutex
}

// DatabaseUpdater Declares a set of functions used for Database updates.
type DatabaseUpdater interface {
	Set(value []Position) error
	Close()
	Reset() error
	Clear() error
}

// NewSessionControl creates and initializes a new SessionControl struct
func NewSessionControl(id string, x, y, numBytes int, db DatabaseUpdater) *SessionControl {
	scb := &SessionControl{
		ID:         id,
		SizeX:      x,
		SizeY:      y,
		NumBytes:   numBytes,
		IsActive:   true,
		Board:      make(chan []Position),
		DBCache:    make(chan []Position),
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
		msg := <-scb.Board

		scb.Mu.Lock()
		for _, clientConn := range scb.Clients { // Send to all connected clients
			clientConn.WriteJSON(&msg) // ignore error
		}
		scb.Mu.Unlock()
	}
}

// UpdateDatabase Updates database according to given position values
func (scb *SessionControl) updateDatabase() {
	for scb.IsActive {
		board := <-scb.DBCache

		if board[0].Action == "clear" {
			scb.DB.Reset()
			continue
		}

		scb.DB.Set(board)
	}
	scb.DB.Close()
}

// SetInactive set the current session inactive
func (scb *SessionControl) SetInactive() {
	scb.IsActive = false

	// clear database
	scb.DB.Clear()
}
