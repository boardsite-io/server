package session

import "github.com/boardsite-io/server/redis"

// Style declares the stroke style.
type Style struct {
	Color   string  `json:"color"`
	Width   float64 `json:"width"`
	Opacity float64 `json:"opacity"`
}

// Stroke declares the structure of most stoke types.
type Stroke struct {
	Type   int       `json:"type"`
	ID     string    `json:"id,omitempty"`
	PageID string    `json:"pageId,omitempty"`
	UserID string    `json:"userId"`
	X      float64   `json:"x"`
	Y      float64   `json:"y"`
	ScaleX float64   `json:"scaleX,omitempty"`
	ScaleY float64   `json:"scaleY,omitempty"`
	Points []float64 `json:"points,omitempty"`
	Style  Style     `json:"style,omitempty"`
}

var _ redis.Stroke = (*Stroke)(nil)

// IsDeleted verifies whether stroke is deleted or not
func (s *Stroke) IsDeleted() bool {
	return s.Type == 0
}

// Id returns the id of the stroke
func (s *Stroke) Id() string {
	return s.ID
}

// UserId returns the userid of the stroke
func (s *Stroke) UserId() string {
	return s.UserID
}

// PageId returns the page id of the stroke
func (s *Stroke) PageId() string {
	return s.PageID
}
