package database

import (
	"fmt"
	"strings"

	"github.com/gomodule/redigo/redis"

	"github.com/heat1q/boardsite/pkg/api"
)

// DatabaseUpdater Declares a set of functions used for Database updates.
// type DatabaseUpdater interface {
// 	Delete(id string) error
// 	Update(value []api.StrokeReader) error
// 	Close()
// 	Clear() error
// }

// RedisDB Holds the connection to the DB
type RedisDB struct {
	Conn     redis.Conn
	BoardKey string
}

// NewRedisConn Sets up redis DB connection with credentials
func NewRedisConn(sessionID string) (*RedisDB, error) {
	// TODO parse from config
	conn, err := redis.Dial("tcp", "arch-pi4:6379")

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

// Update updates board strokes in the database
//
// Creates a JSON encoding for each slice entry which
// is stored to the database
//
// Delete the stroke with given id, if type is set to
// "delete"
func (db *RedisDB) Update(stroke []api.Stroke) error {
	for i := range stroke {
		if stroke[i].IsDeleted() {
			db.Conn.Send("HDEL", db.BoardKey, stroke[i].GetID())
		} else {
			if strokeStr, err := stroke[i].JSONStringify(); err == nil {
				db.Conn.Send("HMSET", db.BoardKey, stroke[i].GetID(), strokeStr)
			}
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
