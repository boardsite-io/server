package session

import (
	"context"
	"errors"
	"fmt"
	"sync"

	gws "github.com/gorilla/websocket"

	"github.com/google/uuid"

	"github.com/heat1q/boardsite/api/config"
	apiErrors "github.com/heat1q/boardsite/api/errors"
	"github.com/heat1q/boardsite/api/types"
	"github.com/heat1q/boardsite/attachment"
	"github.com/heat1q/boardsite/redis"
)

// TODO move to config
const maxUsers = 50

type CreateSessionResponse struct {
	Config `json:"config"`
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
	SetConfig(cfg Config) error

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
	// Receive handles data received in the session
	Receive(ctx context.Context, msg *types.Message) error
	// Attachments returns the session's attachment handler
	Attachments() attachment.Handler
	// Broadcaster returns the session's broadcaster
	Broadcaster() Broadcaster
	// NumUsers returns the number of active users in the session
	NumUsers() int
}

type Config struct {
	ID     string `json:"id"`
	Host   string `json:"host"`
	Secret string `json:"-"`

	config.Session
	Password *string `json:"password,omitempty"`
}

func NewConfig(sessionCfg config.Session) Config {
	return Config{
		Session: sessionCfg,
		Secret:  uuid.NewString(),
	}
}

func (c *Config) validate() error {
	if c.MaxUsers > maxUsers {
		return fmt.Errorf("invalid MaxUsers")
	}
	return nil
}

// controlBlock holds the information and channels for sessions
type controlBlock struct {
	cfg Config

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
	scb.broadcaster.Close() <- struct{}{}
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
	return scb.numUsers
}

func (scb *controlBlock) Config() Config {
	return scb.cfg
}

func (scb *controlBlock) SetConfig(incoming Config) error {
	if err := incoming.validate(); err != nil {
		return apiErrors.ErrBadRequest.Wrap(apiErrors.WithErrorf("validate config: %w", err))
	}
	// TODO replace only non-nil
	//if incoming.Host != "" {
	//	scb.cfg.Host = incoming.Host
	//}
	if incoming.MaxUsers > 0 {
		scb.cfg.MaxUsers = incoming.MaxUsers
	}
	if incoming.ReadOnly != nil {
		scb.cfg.ReadOnly = incoming.ReadOnly
	}
	if incoming.Password != nil {
		scb.cfg.Password = incoming.Password
	}

	scb.broadcaster.Broadcast() <- types.Message{
		Type:    MessageTypeSessionConfig,
		Content: CreateSessionResponse{Config: scb.cfg},
	}

	return nil
}
