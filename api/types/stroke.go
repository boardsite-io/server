package types

import (
	"encoding/json"
)

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
	Type   int       `json:"type"`
	X      float64   `json:"x"`
	Y      float64   `json:"y"`
	Points []float64 `json:"points,omitempty"`
	Style  Style     `json:"style,omitempty"`

	// set for page updates
	PageRank []string `json:"pageRank,omitempty"`

	// pageIDs of pages to clear
	PageClear []string `json:"pageClear,omitempty"`
}

// StrokeReader defines the set of common function
// to interact with strokes
type StrokeReader interface {
	JSONStringify() ([]byte, error)
	IsDeleted() bool
	GetID() string
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

// GetPageID returns the page id of the stroke
func (s *Stroke) GetPageID() string {
	return s.PageID
}
