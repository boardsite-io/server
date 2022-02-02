package session

import (
	"context"
	"errors"
	"sync"

	"github.com/heat1q/boardsite/api/types"
	"github.com/heat1q/boardsite/attachment"
	"github.com/heat1q/boardsite/redis"
)

const defaultMaxUsers = 10

type CreateSessionResponse struct {
	SessionId string `json:"sessionId"`
}

//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 -generate
//counterfeiter:generate . Controller
type Controller interface {
	// ID returns the session id
	ID() string

	// GetPageRank returns the current page rank of a session
	GetPageRank(ctx context.Context) ([]string, error)
	// GetPage returns a page from the session
	GetPage(ctx context.Context, pageId string, withStrokes bool) (*Page, error)
	// AddPages adds pages to the session
	AddPages(ctx context.Context, pageRequest PageRequest) error
	// UpdatePages perform an operation on the given pages.
	// Operations include: clear, delete and update meta data
	UpdatePages(ctx context.Context, pageRequest PageRequest, operation string) error
	// GetPageSync returns the page rank and all pages from the session (optionally with all strokes)
	GetPageSync(ctx context.Context, pageIds []string, withStrokes bool) (*PageSync, error)
	// SyncSession synchronizes the session with the given page rank and pages
	SyncSession(ctx context.Context, sync PageSync) error
	// IsValidPage checks if the given page ids are valid pages
	IsValidPage(ctx context.Context, pageID ...string) bool

	// NewUser creates a new ready user for the session
	NewUser(alias string, color string) (*User, error)
	// GetUserReady returns the current ready user for the given user id
	GetUserReady(userID string) (*User, error)
	// UserConnect connects a ready user to the session
	UserConnect(u *User)
	// UserDisconnect disconnects a user from the session
	UserDisconnect(ctx context.Context, userID string)
	// GetUsers returns all active users in the session
	GetUsers() map[string]*User

	// Close closes a session
	Close()
	// Receive handles data received in the session
	Receive(ctx context.Context, msg *types.Message) error
	// Attachments returns the session's attachment handler
	Attachments() attachment.Handler
	// NumUsers returns the number of active users in the session
	NumUsers() int
}

// controlBlock holds the information and channels for sessions
type controlBlock struct {
	id string

	attachments attachment.Handler
	dispatcher  Dispatcher
	broadcaster Broadcaster

	cache redis.Handler

	muRdyUsr sync.RWMutex
	// users that have previously been created via POST
	// and have not yet joined the session
	usersReady map[string]*User

	muUsr sync.RWMutex
	// Active Client users that are in the session
	// and have an intact WS connection
	users    map[string]*User
	maxUsers int
	numUsers int
}

var _ Controller = (*controlBlock)(nil)

type ControlBlockOption = func(scb *controlBlock)

// WithMaxUsers sets the maximum numbers of users allowed in a session
// This functional argument is passed to NewControlBlock.
func WithMaxUsers(maxUsers int) ControlBlockOption {
	return func(scb *controlBlock) {
		scb.maxUsers = maxUsers
	}
}

// WithCache sets the redis.Handler
// This functional argument is passed to NewControlBlock.
func WithCache(cache redis.Handler) ControlBlockOption {
	return func(scb *controlBlock) {
		scb.cache = cache
	}
}

// WithDispatcher sets the Dispatcher
// This functional argument is passed to NewControlBlock.
func WithDispatcher(dispatcher Dispatcher) ControlBlockOption {
	return func(scb *controlBlock) {
		scb.dispatcher = dispatcher
	}
}

// WithAttachments sets the attachment.Handler
// This functional argument is passed to NewControlBlock.
func WithAttachments(attachments attachment.Handler) ControlBlockOption {
	return func(scb *controlBlock) {
		scb.attachments = attachments
	}
}

// WithBroadcaster sets the Broadcaster
// This functional argument is passed to NewControlBlock.
func WithBroadcaster(broadcaster Broadcaster) ControlBlockOption {
	return func(scb *controlBlock) {
		scb.broadcaster = broadcaster
	}
}

// NewControlBlock creates a new Session controlBlock with unique ID.
func NewControlBlock(sessionId string, options ...ControlBlockOption) (*controlBlock, error) {
	scb := &controlBlock{
		id:         sessionId,
		usersReady: make(map[string]*User),
		users:      make(map[string]*User),
		maxUsers:   defaultMaxUsers,
	}

	for _, o := range options {
		o(scb)
	}

	if scb.cache == nil || scb.dispatcher == nil {
		return nil, errors.New("some of the required handlers are not set")
	}

	if scb.attachments == nil {
		scb.attachments = attachment.NewLocalHandler(sessionId)
	}

	if scb.broadcaster == nil {
		scb.broadcaster = NewBroadcaster(scb.cache)
	}

	return scb, nil
}

// Start starts a session when the first user has joined.
// It binds the broadcaster and starts its goroutines.
func (scb *controlBlock) Start() {
	scb.broadcaster.Bind(scb)
}

// Close sends a close signal
func (scb *controlBlock) Close() {
	scb.broadcaster.Close() <- struct{}{}
}

func (scb *controlBlock) ID() string {
	return scb.id
}

func (scb *controlBlock) Attachments() attachment.Handler {
	return scb.attachments
}

func (scb *controlBlock) NumUsers() int {
	return scb.numUsers
}
