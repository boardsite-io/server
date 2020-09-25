package board

const (
	SizeX    = 5
	SizeY    = 5
	NumBytes = 3
)

// Position Holds the position and value of pixels
type Position struct {
	Action string `json:"action"`
	Value  uint32 `json:"value"`
	X      int    `json:"x"`
	Y      int    `json:"y"`
}
