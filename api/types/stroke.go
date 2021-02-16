package types

import (
	"encoding/json"

	gws "github.com/gorilla/websocket"
)

// User defines info about connected users.
type User struct {
	ID    string    `json:"id"`
	Alias string    `json:"alias"`
	Color string    `json:"color"`
	Conn  *gws.Conn `json:"-"`
}

// Style defines the stoke style.
type Style struct {
	Color string  `json:"color,omitempty"`
	Width float32 `json:"width,omitempty"`
}

// Stroke Holds the Stroke as the basic data type
// for all websocket communication.
type Stroke struct {
	ID     string    `json:"id,omitempty"`
	PageID string    `json:"pageId,omitempty"`
	UserID string    `json:"userId"`
	Type   int       `json:"type"`
	X      float64   `json:"x"`
	Y      float64   `json:"y"`
	Points []float64 `json:"points,omitempty"`
	Style  Style     `json:"style,omitempty"`

	// set for page updates
	PageRank []string `json:"pageRank,omitempty"`

	// pageIDs of pages to clear
	PageClear []string `json:"pageClear,omitempty"`

	// active users in session
	// required in the frontend to display all connected users
	ConnectedUsers map[string]*User `json:"connectedUsers,omitempty"`
}

// StrokeReader defines the set of common function
// to interact with strokes
type StrokeReader interface {
	JSONStringify() ([]byte, error)
	IsDeleted() bool
	GetID() string
	GetUserID() string
	GetPageID() string
}

// JSONStringify return the JSON encoding of Stroke
func (s *Stroke) JSONStringify() ([]byte, error) {
	return json.Marshal(s)
}

// IsDeleted verifies whether stroke is deleted or not
func (s *Stroke) IsDeleted() bool {
	return s.Type == 0
}

// GetID returns the id of the stroke
func (s *Stroke) GetID() string {
	return s.ID
}

// GetUserID returns the userid of the stroke
func (s *Stroke) GetUserID() string {
	return s.UserID
}

// GetPageID returns the page id of the stroke
func (s *Stroke) GetPageID() string {
	return s.PageID
}
