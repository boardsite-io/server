package database

import (
	"fmt"
	"os"
	"strings"

	"github.com/gomodule/redigo/redis"

	"github.com/heat1q/boardsite/api/types"
)

// DatabaseUpdater Declares a set of functions used for Database updates.
// type DatabaseUpdater interface {
// 	Delete(id string) error
// 	Update(value []api.StrokeReader) error
// 	Close()
// 	Clear() error
// }
var (
	redisHost = fmt.Sprintf("%s:%s",
		os.Getenv("B_REDIS_HOST"),
		os.Getenv("B_REDIS_PORT"),
	)
)

// RedisDB Holds the connection to the DB
type RedisDB struct {
	Conn        redis.Conn
	SessionKey  string
	PageRankKey string
}

// NewRedisConn Sets up redis DB connection with credentials
func NewRedisConn(sessionID string) (*RedisDB, error) {
	conn, err := redis.Dial("tcp", redisHost)

	return &RedisDB{
		Conn:        conn,
		SessionKey:  sessionID,
		PageRankKey: sessionID + ".rank",
	}, err
}

// Close Closes connection to redis DB
func (db *RedisDB) Close() {
	db.Conn.Close()
}

// getPageKey return the Redis key for the given PageID.
func (db *RedisDB) getPageKey(pageID string) string {
	return db.SessionKey + "." + pageID
}

// Clear wipes the session from Redis.
//
// Removes all pages and the respective strokes on the pages
func (db *RedisDB) Clear() {
	for _, pageID := range db.GetPages() {
		db.Conn.Send("DEL", db.getPageKey(pageID))
	}
	db.Conn.Send("DEL", db.SessionKey)
	db.Conn.Send("DEL", db.PageRankKey)
	db.Conn.Flush()
}

// Update board strokes in Redis.
// Creates a JSON encoding for each slice entry which
// is stored to the database.
// Delete the stroke with given id if stroke type is set to delete.
func (db *RedisDB) Update(strokes []*types.Stroke) error {
	for i := range strokes {
		pid := db.getPageKey(strokes[i].GetPageID())
		if strokes[i].IsDeleted() {
			db.Conn.Send("HDEL", pid, strokes[i].GetID())
		} else {
			if strokeStr, err := strokes[i].JSONStringify(); err == nil {
				db.Conn.Send("HMSET", pid, strokes[i].GetID(), strokeStr)
			}
		}
	}

	if err := db.Conn.Flush(); err != nil {
		return err
	}

	return nil
}

// Delete a single stroke from Redis given the ID.
// func (db *RedisDB) Delete(strokeID string) error {
// 	_, err := db.Conn.Do("HDEL", db.SessionKey, strokeID)
// 	return err
// }

// FetchAll Fetches all strokes of the board from the DB
//
// Preserves the JSON encoding of DB
func (db *RedisDB) FetchAll() (string, error) {
	keys, err := redis.ByteSlices(db.Conn.Do("HKEYS", db.SessionKey))
	if err != nil {
		return "[]", err
	}

	// slice with capacity equal to num keys
	strokeStr := make([]string, 0, len(keys))

	for i := range keys {
		stroke, _ := redis.ByteSlices(db.Conn.Do("HMGET", db.SessionKey, keys[i]))
		strokeStr = append(strokeStr, string(stroke[0]))
	}

	return fmt.Sprintf("[%s]", strings.Join(strokeStr, ",")), nil
}

// GetPages returns a list of all pageIDs for the current session.
//
// The PageIDs are maintained in a list in redis since the ordering is important
func (db *RedisDB) GetPages() []string {
	pages, err := redis.Strings(db.Conn.Do("ZRANGE", db.PageRankKey, 0, -1))
	if err != nil {
		return []string{}
	}
	return pages
}

// AddPage adds a page with pageID at position index.
//
// Other pages are moved and their score is reassigned
// when pages are added in between
func (db *RedisDB) AddPage(newPageID string, index int) {
	// get all pageids
	pageIDs := db.GetPages()
	if len(pageIDs) > 0 {
		var score, diff, prevIndex int

		if index >= 0 && index < len(pageIDs) { // add page in between
			// increment scores of proceding pages
			for _, pid := range pageIDs[index:] {
				db.Conn.Send("ZINCRBY", db.PageRankKey, 1, pid)
			}
			db.Conn.Flush() // ignore error
			prevIndex = index
			diff = -1
		} else { // append page at the end
			prevIndex = len(pageIDs) - 1
			diff = 1
		}

		// get score of preceding page
		score, _ = redis.Int(db.Conn.Do("ZSCORE", db.PageRankKey, pageIDs[prevIndex]))
		db.Conn.Do("ZADD", db.PageRankKey, "NX", score+diff, newPageID)
	} else { // no pages exist yet
		db.Conn.Do("ZADD", db.PageRankKey, "NX", 0, newPageID)
	}
}

// DeletePage deletes a page and the respective strokes on the page
// and remove the PageID from the list.
func (db *RedisDB) DeletePage(pageID string) {
	db.Conn.Do("DEL", db.getPageKey(pageID))
	db.Conn.Do("ZREM", db.PageRankKey, pageID)
}

// ClearPage removes all strokes with given pageID.
func (db *RedisDB) ClearPage(pageID string) {
	db.Conn.Do("DEL", db.getPageKey(pageID))
}
