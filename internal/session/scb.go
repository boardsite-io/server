package session

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/google/uuid"
	gws "github.com/gorilla/websocket"

	"github.com/boardsite-io/server/internal/attachment"
	"github.com/boardsite-io/server/internal/config"
	"github.com/boardsite-io/server/pkg/redis"
)

type CreateSessionRequest struct {
	ConfigRequest *ConfigRequest `json:"config,omitempty"`
}

type CreateSessionResponse struct {
	Config Config `json:"config"`
}

type GetConfigResponse struct {
	Users  map[string]*User `json:"users"`
	Config `json:"config"`
}

//counterfeiter:generate . Controller
type Controller interface {
	// ID returns the session id
	ID() string
	// Config returns the session config
	Config() Config
	// SetConfig sets the session config
	SetConfig(cfg *ConfigRequest) error

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
	NewUser(userReq UserRequest) (*User, error)
	// UpdateUser updates a user alias or color
	UpdateUser(user User, userReq UserRequest) error
	// UserCanJoin check if a user can join the session
	UserCanJoin(userID string) error
	// UserConnect connects a ready user to the session
	UserConnect(userID string, conn *gws.Conn) error
	// UserDisconnect disconnects a user from the session
	UserDisconnect(ctx context.Context, userID string)
	// KickUser removes a user from the session
	KickUser(userID string) error
	// GetUsers returns all active users in the session
	GetUsers() map[string]*User

	// Close closes a session
	Close()
	// CloseAfter closes a session after a specified timeout and executes fn
	CloseAfter(t time.Duration, fn func())
	// Receive handles data received in the session
	Receive(ctx context.Context, msg *Message, userID string) error
	// Attachments returns the session's attachment handler
	Attachments() attachment.Handler
	// Broadcaster returns the session's broadcaster
	Broadcaster() Broadcaster
	// NumUsers returns the number of active users in the session
	NumUsers() int
	// Allow checks whether a user is allowed to modify the session
	Allow(userID string) bool
}

func NewConfig(sessionCfg config.Session) Config {
	return Config{
		Session: sessionCfg,
		Secret:  uuid.NewString(),
	}
}

// controlBlock holds the information and channels for sessions
type controlBlock struct {
	cfg Config

	attachments attachment.Handler
	dispatcher  Dispatcher
	broadcaster Broadcaster

	// close timer
	timer *time.Timer

	cache redis.Handler

	muRdyUsr sync.RWMutex
	// users that have previously been created via POST
	// and have not yet joined the session
	usersReady map[string]*User

	muUsr sync.RWMutex
	// Active Client users that are in the session
	// and have an intact WS connection
	users    map[string]*User
	numUsers int
}

var _ Controller = (*controlBlock)(nil)

type ControlBlockOption = func(scb *controlBlock)

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
func NewControlBlock(cfg Config, options ...ControlBlockOption) (*controlBlock, error) {
	scb := &controlBlock{
		cfg:        cfg,
		usersReady: make(map[string]*User),
		users:      make(map[string]*User),
	}

	for _, o := range options {
		o(scb)
	}

	if scb.cache == nil || scb.dispatcher == nil {
		return nil, errors.New("some of the required handlers are not set")
	}

	if scb.attachments == nil {
		scb.attachments = attachment.NewLocalHandler(cfg.ID)
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
	scb.broadcaster.Close()
}

func (scb *controlBlock) CloseAfter(t time.Duration, fn func()) {
	if scb.timer == nil {
		scb.timer = time.AfterFunc(t, func() {
			if scb.NumUsers() == 0 {
				scb.Close()
				fn()
			}
		})
		return
	}
	scb.timer.Reset(t)
}

func (scb *controlBlock) ID() string {
	return scb.cfg.ID
}

func (scb *controlBlock) Attachments() attachment.Handler {
	return scb.attachments
}

func (scb *controlBlock) Broadcaster() Broadcaster {
	return scb.broadcaster
}

func (scb *controlBlock) NumUsers() int {
	scb.muUsr.RLock()
	defer scb.muUsr.RUnlock()
	return scb.numUsers
}

func (scb *controlBlock) Config() Config {
	return scb.cfg
}

func (scb *controlBlock) SetConfig(incoming *ConfigRequest) error {
	if err := scb.cfg.Update(incoming); err != nil {
		return err
	}
	scb.broadcaster.Broadcast() <- Message{
		Type:    MessageTypeSessionConfig,
		Content: CreateSessionResponse{Config: scb.cfg},
	}
	return nil
}

func (scb *controlBlock) Allow(userID string) bool {
	if scb.Config().ReadOnly && userID != scb.Config().Host {
		return false
	}
	return true
}
