package redis

import (
	"context"
	"fmt"
	"time"

	"github.com/gomodule/redigo/redis"

	"github.com/heat1q/boardsite/api/types"
)

const (
	maxNumIdleConnections = 3
	maxIdleTimeoutSec     = 5
)

type Handler interface {
	// ClearSession wipes the session from Redis.
	//
	// Removes all pages and the respective strokes on the pages
	ClearSession(ctx context.Context, sessionID string) error
	// Update board strokes in Redis.
	//
	// Creates a JSON encoding for each slice entry which
	// is stored to the database.
	// Delete the stroke with given id if stroke type is set to delete.
	Update(ctx context.Context, sessionID string, strokes []*types.Stroke) error
	// FetchStrokesRaw Fetches all strokes of the specified page.
	//
	// Preserves the JSON encoding of Redis and returns an array of
	// a stringified stroke objects.
	FetchStrokesRaw(ctx context.Context, sessionID, pageID string) ([][]byte, error)
	// GetPages returns a list of all pageIDs for the current session.
	//
	// The PageIDs are maintained in a list in redis since the ordering is important
	GetPages(ctx context.Context, sessionID string) ([]string, error)
	// GetPagesMeta returns a slice of all page meta data.
	GetPagesMeta(ctx context.Context, sessionID string, pageIDs ...string) (map[string]*types.PageMeta, error)
	// UpdatePageMeta updates a page meta data.
	UpdatePageMeta(ctx context.Context, sessionID, pageID string, update *types.PageMeta) error
	// AddPage adds a page with pageID at position index.
	//
	// Other pages are moved and their score is reassigned
	// when pages are added in between
	AddPage(ctx context.Context, sessionID, newPageID string, index int, meta *types.PageMeta) error
	// DeletePage deletes a page and the respective strokes on the page and remove the PageID from the list.
	DeletePage(ctx context.Context, sessionID, pageID string) error
	// ClearPage removes all strokes with given pageID.
	ClearPage(ctx context.Context, sessionID, pageID string) error
	ClosePool() error
}

type handler struct {
	pool *redis.Pool
}

// New creates a new redis intance and initializes the Redis thread pool.
func New(host string, port uint16) (Handler, error) {
	redisHost := fmt.Sprintf("%s:%d", host, port)
	h := &handler{
		pool: newPool(redisHost),
	}
	if err := h.Ping(); err != nil {
		return nil, err
	}
	return h, nil
}

func newPool(redisHost string) *redis.Pool {
	return &redis.Pool{
		MaxIdle:     maxNumIdleConnections,
		IdleTimeout: maxIdleTimeoutSec * time.Second,
		Dial: func() (redis.Conn, error) {
			return redis.Dial("tcp", redisHost)
		},
	}
}

// ClosePool closes the Redis thread pool.
func (h *handler) ClosePool() error {
	return h.pool.Close()
}

// Ping pings the connection to Redis and returns an error
// if the connection cannot be established.
func (h *handler) Ping() error {
	conn := h.pool.Get()
	defer conn.Close()

	if _, err := conn.Do("PING"); err != nil {
		return fmt.Errorf("PING redis failed: %v", err)
	}
	return nil
}
