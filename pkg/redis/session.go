package redis

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/gomodule/redigo/redis"
)

type Stroke interface {
	Id() string
	PageId() string
	IsDeleted() bool
}

// getPageRankKey returns the Redis key for the pageRank of a session.
func getPageRankKey(sessionId string) string {
	return sessionId + ".rank"
}

// getStrokesKey returns the Redis key for the given pageId.
func getStrokesKey(sessionId, pageId string) string {
	return fmt.Sprintf("%s.%s.strokes", sessionId, pageId)
}

// getPageMetaKey returns the redis key for page meta data.
func getPageMetaKey(sessionId, pageId string) string {
	return fmt.Sprintf("%s.%s.meta", sessionId, pageId)
}

func (h *handler) UpdateStrokes(ctx context.Context, sessionId string, strokes ...Stroke) error {
	conn, err := h.pool.GetContext(ctx)
	if err != nil {
		return err
	}
	defer conn.Close()

	for _, s := range strokes {
		pid := getStrokesKey(sessionId, s.PageId())
		if s.IsDeleted() {
			if err := conn.Send("HDEL", pid, s.Id()); err != nil {
				return err
			}
		} else {
			bytes, err := json.Marshal(s)
			if err != nil {
				return err
			}
			if err := conn.Send("HMSET", pid, s.Id(), bytes); err != nil {
				return err
			}
		}

	}
	return conn.Flush()
}

func (h *handler) GetPageStrokes(ctx context.Context, sessionId, pageId string) ([][]byte, error) {
	pid := getStrokesKey(sessionId, pageId)
	keys, err := redis.Strings(h.Do(ctx, "HKEYS", pid))
	if err != nil {
		return nil, err
	}
	if len(keys) == 0 { // page is empty
		return [][]byte{}, nil
	}

	query := make([]any, 1, len(keys)+1)
	query[0] = pid
	for _, key := range keys {
		query = append(query, key)
	}

	return redis.ByteSlices(h.Do(ctx, "HMGET", query...))
}

func (h *handler) GetPageRank(ctx context.Context, sessionId string) ([]string, error) {
	pages, err := redis.Strings(h.Do(ctx, "ZRANGE", getPageRankKey(sessionId), 0, -1))
	if err != nil {
		return nil, err
	}
	return pages, nil
}

func (h *handler) GetPageMeta(ctx context.Context, sessionId, pageId string, meta any) error {
	resp, err := redis.Bytes(h.Do(ctx, "GET", getPageMetaKey(sessionId, pageId)))
	if err != nil {
		return err
	}
	if err := json.Unmarshal(resp, meta); err != nil {
		return err
	}
	return nil
}

func (h *handler) SetPageMeta(ctx context.Context, sessionId, pageId string, meta any) error {
	pMeta, err := json.Marshal(meta)
	if err != nil {
		return err
	}
	_, err = h.Do(ctx, "SET", getPageMetaKey(sessionId, pageId), pMeta)
	return err
}

func (h *handler) AddPage(ctx context.Context, sessionId, newpageId string, index int, meta any) error {
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
			"SET", getPageMetaKey(sessionId, newpageId), pMeta); err != nil {
			return err
		}
	}

	// get all pageIds
	pageRankKey := getPageRankKey(sessionId)
	pageRank, err := h.GetPageRank(ctx, sessionId)
	if err != nil {
		return err
	}
	if len(pageRank) > 0 {
		var score, diff, prevIndex int

		if index >= 0 && index < len(pageRank) { // add page in between
			// increment scores of proceeding pages
			for _, pid := range pageRank[index:] {
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
			prevIndex = len(pageRank) - 1
			diff = 1
		}

		// get score of preceding page
		score, err = redis.Int(conn.Do("ZSCORE", pageRankKey, pageRank[prevIndex]))
		if err != nil {
			return err
		}
		if _, err := conn.Do(
			"ZADD", pageRankKey, "NX", score+diff, newpageId); err != nil {
			return err
		}
	} else { // no pages exist yet
		if _, err := conn.Do("ZADD", pageRankKey, "NX", 0, newpageId); err != nil {
			return err
		}
	}
	return nil
}

func (h *handler) DeletePage(ctx context.Context, sessionId, pageId string) error {
	conn, err := h.pool.GetContext(ctx)
	if err != nil {
		return err
	}
	defer conn.Close()

	if err := conn.Send(
		"DEL",
		getStrokesKey(sessionId, pageId),
		getPageMetaKey(sessionId, pageId),
	); err != nil {
		return err
	}
	if err := conn.Send(
		"ZREM",
		getPageRankKey(sessionId),
		pageId,
	); err != nil {
		return err
	}
	return conn.Flush()
}

func (h *handler) ClearSession(ctx context.Context, sessionId string) error {
	pageRank, err := h.GetPageRank(ctx, sessionId)
	if err != nil {
		return err
	}

	if len(pageRank) == 0 { // nothing to do
		return nil
	}

	query := make([]any, 1, len(pageRank)*2+1)
	query[0] = getPageRankKey(sessionId)
	for _, pid := range pageRank {
		query = append(query, getStrokesKey(sessionId, pid), getPageMetaKey(sessionId, pid))
	}

	_, err = h.Do(ctx, "DEL", query...)
	return err
}

func (h *handler) ClearPage(ctx context.Context, sessionId, pageId string) error {
	_, err := h.Do(ctx, "DEL", getStrokesKey(sessionId, pageId))
	return err
}
