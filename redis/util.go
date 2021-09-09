package redis

import (
	"context"
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

func (h *handler) ClearSession(ctx context.Context, sessionID string) error {
	conn, err := h.pool.GetContext(ctx)
	if err != nil {
		return err
	}
	defer conn.Close()

	pages, err := h.GetPages(ctx, sessionID)
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

func (h *handler) Update(ctx context.Context, sessionID string, strokes []*types.Stroke) error {
	conn, err := h.pool.GetContext(ctx)
	if err != nil {
		return err
	}
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

func (h *handler) FetchStrokesRaw(ctx context.Context, sessionID, pageID string) ([][]byte, error) {
	conn, err := h.pool.GetContext(ctx)
	if err != nil {
		return nil, err
	}
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

func (h *handler) GetPages(ctx context.Context, sessionID string) ([]string, error) {
	conn, err := h.pool.GetContext(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	pages, err := redis.Strings(
		conn.Do("ZRANGE", getPageRankKey(sessionID), 0, -1))
	if err != nil {
		return nil, err
	}
	return pages, nil
}

func (h *handler) GetPagesMeta(ctx context.Context, sessionID string, pageIDs ...string) (map[string]*types.PageMeta, error) {
	conn, err := h.pool.GetContext(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	metaPages := make(map[string]*types.PageMeta)
	for _, pid := range pageIDs {
		var meta types.PageMeta
		if resp, err := redis.Bytes(conn.Do("GET", getPageMetaKey(sessionID, pid))); err == nil {
			if err := json.Unmarshal(resp, &meta); err != nil {
				return nil, err
			}
		} else {
			return nil, err
		}
		metaPages[pid] = &meta
	}
	return metaPages, nil
}

func (h *handler) UpdatePageMeta(ctx context.Context, sessionID, pageID string, update *types.PageMeta) error {
	conn, err := h.pool.GetContext(ctx)
	if err != nil {
		return err
	}
	defer conn.Close()

	var meta types.PageMeta
	resp, err := redis.Bytes(
		conn.Do("GET", getPageMetaKey(sessionID, pageID)))
	if err != nil {
		return err
	}
	if err := json.Unmarshal(resp, &meta); err != nil {
		return err
	}
	tmp, err := json.Marshal(update)
	if err != nil {
		return err
	}
	// update only non-zero entries
	if err := json.Unmarshal(tmp, &meta); err != nil {
		return err
	}
	// store result
	pMeta, err := json.Marshal(meta)
	if err != nil {
		return err
	}
	if _, err := conn.Do(
		"SET", getPageMetaKey(sessionID, pageID), pMeta); err != nil {
		return err
	}
	return nil
}

func (h *handler) AddPage(ctx context.Context, sessionID, newPageID string, index int, meta *types.PageMeta) error {
	conn, err := h.pool.GetContext(ctx)
	if err != nil {
		return err
	}
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
	pageIDs, err := h.GetPages(ctx, sessionID)
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

func (h *handler) DeletePage(ctx context.Context, sessionID, pageID string) error {
	conn, err := h.pool.GetContext(ctx)
	if err != nil {
		return err
	}
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

func (h *handler) ClearPage(ctx context.Context, sessionID, pageID string) error {
	conn, err := h.pool.GetContext(ctx)
	if err != nil {
		return err
	}
	defer conn.Close()
	if _, err := conn.Do("DEL", getPageKey(sessionID, pageID)); err != nil {
		return err
	}
	return nil
}
