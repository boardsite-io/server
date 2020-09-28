package board

const (
	SizeX    = 10
	SizeY    = 10
	NumBytes = 3
)

// Position Holds the position and value of pixels
type Position struct {
	Action string `json:"action"`
	Value  uint32 `json:"color"`
	X      int    `json:"x"`
	Y      int    `json:"y"`
}
