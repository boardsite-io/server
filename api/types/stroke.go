package types

import (
	"encoding/json"
)

// Stroke Holds the Stroke and value of pixels
type Stroke struct {
	ID     string    `json:"id"`
	PageID string    `json:"pageId"`
	Type   int       `json:"type"`
	X      float64   `json:"x"`
	Y      float64   `json:"y"`
	Points []float64 `json:"points"`
	Style  struct {
		Color string  `json:"color"`
		Width float64 `json:"width"`
	} `json:"style"`
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
