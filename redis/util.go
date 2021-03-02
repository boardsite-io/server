package redis

import (
	"github.com/gomodule/redigo/redis"

	"github.com/heat1q/boardsite/api/types"
)

// getPageRankKey returns the Redis key for the pageRank of a session.
func getPageRankKey(sessionID string) string {
	return sessionID + ".rank"
}

// getPageKey returns the Redis key for the given PageID.
func getPageKey(sessionID, pageID string) string {
	return sessionID + "." + pageID
}

// ClearSession wipes the session from Redis.
//
// Removes all pages and the respective strokes on the pages
func ClearSession(sessionID string) {
	conn := Pool.Get()
	defer conn.Close()

	for _, pageID := range GetPages(sessionID) {
		conn.Send("DEL", getPageKey(sessionID, pageID))
	}
	conn.Send("DEL", getPageRankKey(sessionID))
	conn.Flush()
}

// Update board strokes in Redis.
//
// Creates a JSON encoding for each slice entry which
// is stored to the database.
// Delete the stroke with given id if stroke type is set to delete.
func Update(sessionID string, strokes []*types.Stroke) error {
	conn := Pool.Get()
	defer conn.Close()

	for i := range strokes {
		pid := getPageKey(sessionID, strokes[i].GetPageID())
		if strokes[i].IsDeleted() {
			conn.Send("HDEL", pid, strokes[i].GetID())
		} else {
			if strokeStr, err := strokes[i].JSONStringify(); err == nil {
				conn.Send("HMSET", pid, strokes[i].GetID(), strokeStr)
			}
		}
	}

	if err := conn.Flush(); err != nil {
		return err
	}

	return nil
}

// FetchStrokesRaw Fetches all strokes of the specified page.
//
// Preserves the JSON encoding of Redis and returns an array of
// a stringified stroke objects.
func FetchStrokesRaw(sessionID, pageID string) ([][]byte, error) {
	conn := Pool.Get()
	defer conn.Close()

	pid := getPageKey(sessionID, pageID)
	keys, err := redis.Strings(conn.Do("HKEYS", pid))
	if err != nil {
		return nil, err
	}
	if len(keys) == 0 { // page is empty
		return [][]byte{}, nil
	}

	query := make([]interface{}, 1, len(keys)+1)
	query[0] = pid
	for _, key := range keys {
		query = append(query, key)
	}

	strokes, errFetch := redis.ByteSlices(conn.Do("HMGET", query...))
	if errFetch != nil {
		return nil, errFetch
	}
	return strokes, nil
}

// GetPages returns a list of all pageIDs for the current session.
//
// The PageIDs are maintained in a list in redis since the ordering is important
func GetPages(sessionID string) []string {
	conn := Pool.Get()
	defer conn.Close()

	pages, err := redis.Strings(conn.Do("ZRANGE", getPageRankKey(sessionID), 0, -1))
	if err != nil {
		return []string{}
	}
	return pages
}

// AddPage adds a page with pageID at position index.
//
// Other pages are moved and their score is reassigned
// when pages are added in between
func AddPage(sessionID, newPageID string, index int) {
	conn := Pool.Get()
	defer conn.Close()

	// get all pageids
	pageRankKey := getPageRankKey(sessionID)
	pageIDs := GetPages(sessionID)
	if len(pageIDs) > 0 {
		var score, diff, prevIndex int

		if index >= 0 && index < len(pageIDs) { // add page in between
			// increment scores of proceding pages
			for _, pid := range pageIDs[index:] {
				conn.Send("ZINCRBY", pageRankKey, 1, pid)
			}
			conn.Flush() // ignore error
			prevIndex = index
			diff = -1
		} else { // append page at the end
			prevIndex = len(pageIDs) - 1
			diff = 1
		}

		// get score of preceding page
		score, _ = redis.Int(conn.Do("ZSCORE", pageRankKey, pageIDs[prevIndex]))
		conn.Do("ZADD", pageRankKey, "NX", score+diff, newPageID)
	} else { // no pages exist yet
		conn.Do("ZADD", pageRankKey, "NX", 0, newPageID)
	}
}

// DeletePage deletes a page and the respective strokes on the page
// and remove the PageID from the list.
func DeletePage(sessionID, pageID string) {
	conn := Pool.Get()
	defer conn.Close()

	conn.Do("DEL", getPageKey(sessionID, pageID))
	conn.Do("ZREM", getPageRankKey(sessionID), pageID)
}

// ClearPage removes all strokes with given pageID.
func ClearPage(sessionID, pageID string) {
	conn := Pool.Get()
	defer conn.Close()

	conn.Do("DEL", getPageKey(sessionID, pageID))
}
