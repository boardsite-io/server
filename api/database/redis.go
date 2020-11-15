package database

import (
	"encoding/binary"

	"github.com/gomodule/redigo/redis"

	"boardsite/api/board"
)

// BoardDB Holds the connection to the DB
type BoardDB struct {
	SizeX, SizeY, NumBytes int
	BoardLen               int
	BoardKey               string
	Conn                   redis.Conn
}

// NewConnection Sets up redis DB connection with credentials
func NewConnection(sessionID string, x, y, numBytes int) (*BoardDB, error) {
	// TODO parse from config
	conn, err := redis.Dial("tcp", "localhost:6379")

	return &BoardDB{
		Conn:     conn,
		SizeX:    x,
		SizeY:    y,
		NumBytes: numBytes,
		BoardLen: numBytes * x * y,
		BoardKey: sessionID,
	}, err
}

// Close Closes connection to redis DB
func (db *BoardDB) Close() {
	db.Conn.Close()
}

// Reset Creates an new (empty) board
func (db *BoardDB) Reset() error {
	// empty slice of board size
	board := make([]byte, db.BoardLen)

	// default color 0xffffff
	for i := range board {
		board[i] = 0xff
	}

	_, err := db.Conn.Do("SET", db.BoardKey, board)
	return err
}

// Clear clears the board from Redis
func (db *BoardDB) Clear() error {
	_, err := db.Conn.Do("DEL", db.BoardKey)
	return err
}

// Set Stores board values to the database
func (db *BoardDB) Set(boardpos []board.Position) error {
	b := make([]byte, 4)

	for _, pos := range boardpos {
		// encode to byte slice
		binary.LittleEndian.PutUint32(b, pos.Value)

		// store only boardValBytes least significant bytes
		if pos.X < db.SizeX && pos.Y < db.SizeY {
			db.Conn.Send("SETRANGE", db.BoardKey, db.getDBIndex(pos.X, pos.Y), b[:db.NumBytes])
		}
	}

	if err := db.Conn.Flush(); err != nil {
		return err
	}

	return nil
}

func (db *BoardDB) getDBIndex(x, y int) int {
	return (db.SizeX*y + x) * db.NumBytes
}

// FetchAll Fetches all the values of the board from the DB
func (db *BoardDB) FetchAll() ([]board.Position, error) {
	// slice with max capacity
	boardpos := make([]board.Position, 0, db.BoardLen)

	reply, err := db.Conn.Do("GET", db.BoardKey)
	if err != nil {
		return nil, err
	} else if reply == nil {
		return nil, nil
	}

	data := reply.([]byte)

	for i := 0; i < db.BoardLen; i += db.NumBytes {
		// convert the board.NumBytes bytes to uint32
		var value uint32
		for j := 0; j < db.NumBytes; j++ {
			value |= uint32(data[i+j]) << (8 * j) // little endian
		}

		// only retrieve non-white (0xffffff) values
		if value != 0xffffff {
			boardpos = append(boardpos, board.Position{Value: value, X: i / db.NumBytes % db.SizeX, Y: i / db.NumBytes / db.SizeX})
		}
	}

	return boardpos, nil
}
