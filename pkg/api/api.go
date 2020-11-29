package api

import (
	"encoding/json"
)

// Stroke Holds the Stroke and value of pixels
type Stroke struct {
	ID        string    `json:"id"`
	Type      string    `json:"type"`
	Color     string    `json:"color"`
	LineWidth float64   `json:"line_width"`
	Position  []float64 `json:"position"`
}

// StrokeReader defines the set of common function
// to interact with strokes
type StrokeReader interface {
	JSONStringify() ([]byte, error)
	IsDeleted() bool
	GetID() string
}

// SetupForm Form to setup a new board
type SetupForm struct {
	X int `json:"x"`
	Y int `json:"y"`
}

type CreateBoardResponse struct {
	ID string `json:"id"`
}

type BoardRequest struct {
	Action string `json:"action"`
}

// JSONStringify return the JSON encoding of Stroke
func (s *Stroke) JSONStringify() ([]byte, error) {
	return json.Marshal(s)
}

// IsDeleted verifies whether stroke is deleted or not
func (s *Stroke) IsDeleted() bool {
	return s.Type == "delete"
}

// GetID returns the id of the stroke
func (s *Stroke) GetID() string {
	return s.ID
}
