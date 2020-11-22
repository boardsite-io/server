package database

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/gomodule/redigo/redis"

	"boardsite/api/board"
)

// RedisDB Holds the connection to the DB
type RedisDB struct {
	Conn     redis.Conn
	BoardKey string
}

type strokeObj struct {
	ID string `json:"id"`
}

// NewConnection Sets up redis DB connection with credentials
func NewConnection(sessionID string) (*RedisDB, error) {
	// TODO parse from config
	conn, err := redis.Dial("tcp", "localhost:6379")

	return &RedisDB{
		Conn:     conn,
		BoardKey: sessionID,
	}, err
}

// Close Closes connection to redis DB
func (db *RedisDB) Close() {
	db.Conn.Close()
}

// Clear clears the board from Redis
func (db *RedisDB) Clear() error {
	_, err := db.Conn.Do("DEL", db.BoardKey)
	return err
}

// Set Stores board strokes to the database
//
// Creates a JSON encoding for each slice entry which
// is stored to the database
func (db *RedisDB) Set(stroke []board.Stroke) error {
	for i := range stroke {
		if strokeStr, err := json.Marshal(&stroke[i]); err == nil {
			db.Conn.Send("HMSET", db.BoardKey, stroke[i].ID, strokeStr)
		} else {
			fmt.Println(err)
		}
	}

	if err := db.Conn.Flush(); err != nil {
		return err
	}

	return nil
}

// Delete deletes stroke from database by ID
func (db *RedisDB) Delete(strokeID string) error {
	_, err := db.Conn.Do("HDEL", db.BoardKey, strokeID)
	return err
}

// FetchAll Fetches all strokes of the board from the DB
//
// Preserves the JSON encoding of DB
func (db *RedisDB) FetchAll() (string, error) {
	keys, err := redis.ByteSlices(db.Conn.Do("HKEYS", db.BoardKey))
	if err != nil {
		return "[]", err
	}

	// slice with capacity equal to num keys
	strokeStr := make([]string, 0, len(keys))

	for i := range keys {
		stroke, _ := redis.ByteSlices(db.Conn.Do("HMGET", db.BoardKey, keys[i]))
		strokeStr = append(strokeStr, string(stroke[0]))
	}

	return fmt.Sprintf("[%s]", strings.Join(strokeStr, ",")), nil
}
