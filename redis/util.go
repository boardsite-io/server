package redis

import (
	"encoding/json"

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

// getPageMetaKey returns the redis key for page meta data.
func getPageMetaKey(sessionID, pageID string) string {
	return getPageKey(sessionID, pageID) + ".meta"
}

// ClearSession wipes the session from Redis.
//
// Removes all pages and the respective strokes on the pages
func ClearSession(sessionID string) error {
	conn := Pool.Get()
	defer conn.Close()

	pages, err := GetPages(sessionID)
	if err != nil {
		return err
	}

	if len(pages) == 0 { // nothing to do
		return nil
	}

	query := make([]interface{}, 1, len(pages)*2+1)
	query[0] = getPageRankKey(sessionID)
	for _, pid := range pages {
		query = append(query, getPageKey(sessionID, pid), getPageMetaKey(sessionID, pid))
	}
	if _, err := conn.Do("DEL", query...); err != nil {
		return err
	}
	return nil
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
		var err error
		if strokes[i].IsDeleted() {
			err = conn.Send("HDEL", pid, strokes[i].GetID())
		} else {
			if strokeStr, err := strokes[i].JSONStringify(); err == nil {
				err = conn.Send("HMSET", pid, strokes[i].GetID(), strokeStr)
			}
		}
		if err != nil {
			return err
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
func GetPages(sessionID string) ([]string, error) {
	conn := Pool.Get()
	defer conn.Close()

	pages, err := redis.Strings(
		conn.Do("ZRANGE", getPageRankKey(sessionID), 0, -1))
	if err != nil {
		return nil, err
	}
	return pages, nil
}

// GetPagesMeta returns a slice of all page meta data.
func GetPagesMeta(sessionID string, pageIDs []string) ([]*types.PageMeta, error) {
	conn := Pool.Get()
	defer conn.Close()

	metaRank := make([]*types.PageMeta, len(pageIDs))
	for i, pid := range pageIDs {
		var meta types.PageMeta
		if resp, err := redis.Bytes(conn.Do("GET", getPageMetaKey(sessionID, pid))); err == nil {
			if err := json.Unmarshal(resp, &meta); err != nil {
				return nil, err
			}
		} else {
			return nil, err
		}
		metaRank[i] = &meta
	}
	return metaRank, nil
}

// AddPage adds a page with pageID at position index.
//
// Other pages are moved and their score is reassigned
// when pages are added in between
func AddPage(sessionID, newPageID string,
	index int, meta *types.PageMeta) error {
	conn := Pool.Get()
	defer conn.Close()

	if meta != nil {
		pMeta, err := json.Marshal(meta)
		if err != nil {
			return err
		}
		if _, err := conn.Do(
			"SET", getPageMetaKey(sessionID, newPageID), pMeta); err != nil {
			return err
		}
	}

	// get all pageids
	pageRankKey := getPageRankKey(sessionID)
	pageIDs, err := GetPages(sessionID)
	if err != nil {
		return err
	}
	if len(pageIDs) > 0 {
		var score, diff, prevIndex int

		if index >= 0 && index < len(pageIDs) { // add page in between
			// increment scores of proceding pages
			for _, pid := range pageIDs[index:] {
				if err := conn.Send(
					"ZINCRBY", pageRankKey, 1, pid); err != nil {
					return err
				}
			}
			if err := conn.Flush(); err != nil {
				return err
			}
			prevIndex = index
			diff = -1
		} else { // append page at the end
			prevIndex = len(pageIDs) - 1
			diff = 1
		}

		// get score of preceding page
		score, err = redis.Int(conn.Do("ZSCORE", pageRankKey, pageIDs[prevIndex]))
		if err != nil {
			return err
		}
		if _, err := conn.Do(
			"ZADD", pageRankKey, "NX", score+diff, newPageID); err != nil {
			return err
		}
	} else { // no pages exist yet
		if _, err := conn.Do("ZADD", pageRankKey, "NX", 0, newPageID); err != nil {
			return err
		}
	}
	return nil
}

// DeletePage deletes a page and the respective strokes on the page
// and remove the PageID from the list.
func DeletePage(sessionID, pageID string) error {
	conn := Pool.Get()
	defer conn.Close()

	if err := conn.Send(
		"DEL",
		getPageKey(sessionID, pageID),
		getPageMetaKey(sessionID, pageID),
	); err != nil {
		return err
	}
	if err := conn.Send(
		"ZREM",
		getPageRankKey(sessionID),
		pageID,
	); err != nil {
		return err
	}
	return conn.Flush()
}

// ClearPage removes all strokes with given pageID.
func ClearPage(sessionID, pageID string) error {
	conn := Pool.Get()
	defer conn.Close()
	_, err := conn.Do("DEL", getPageKey(sessionID, pageID))
	return err
}
