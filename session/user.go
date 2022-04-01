package session

import (
	"context"
	"errors"
	"regexp"

	"github.com/google/uuid"
	gws "github.com/gorilla/websocket"

	apiErrors "github.com/heat1q/boardsite/api/errors"
	"github.com/heat1q/boardsite/api/types"
)

var (
	ErrReadyUser      = errors.New("ready user not found")
	ErrUserConnected  = errors.New("user already connected")
	ErrMaxUserReached = errors.New("maximum number of connected users reached")
)

var (
	htmlColor = regexp.MustCompile("^#[a-fA-F0-9]{6}$")
	aliasExp  = regexp.MustCompile("^[a-zA-Z0-9-_]{4,32}$")
)

// User declares some information about connected users.
type User struct {
	ID    string    `json:"id"`
	Alias string    `json:"alias"`
	Color string    `json:"color"`
	Conn  *gws.Conn `json:"-"`
}

func (u *User) validate() error {
	if !aliasExp.MatchString(u.Alias) {
		return apiErrors.From(apiErrors.BadUsername)
	}
	if !htmlColor.MatchString(u.Color) {
		return apiErrors.ErrBadRequest.Wrap(apiErrors.WithErrorf("incorrect html color"))
	}
	return nil
}

type UserRequest struct {
	Password string `json:"password"`
	User     `json:"user"`
}

type userHostContent struct {
	Secret string `json:"secret"`
}

// NewUser generate a new user struct based on
// the alias and color attribute
//
// Does some sanitize checks.
func (scb *controlBlock) NewUser(userReq UserRequest) (*User, error) {
	user := &User{
		ID:    uuid.NewString(),
		Alias: userReq.Alias,
		Color: userReq.Color,
	}

	if scb.Config().Password != "" && userReq.Password != scb.Config().Password {
		return nil, apiErrors.From(apiErrors.WrongPassword)
	}
	if err := user.validate(); err != nil {
		return nil, err
	}

	// set user waiting
	err := scb.userReady(user)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (scb *controlBlock) UpdateUser(user User, userReq UserRequest) error {
	if userReq.Alias == "" {
		userReq.Alias = user.Alias
	}
	if userReq.Color == "" {
		userReq.Color = user.Color
	}

	if err := userReq.validate(); err != nil {
		return err
	}

	scb.muUsr.Lock()
	scb.users[user.ID].Alias = userReq.Alias
	scb.users[user.ID].Color = userReq.Color
	scb.muUsr.Unlock()

	scb.broadcaster.Broadcast() <- types.Message{
		Type:    MessageTypeUserSync,
		Content: scb.GetUsers(),
	}

	return nil
}

// UserReady adds an user to the usersReady map.
func (scb *controlBlock) userReady(u *User) error {
	scb.muUsr.RLock()
	defer scb.muUsr.RUnlock()
	if scb.numUsers >= scb.cfg.MaxUsers {
		return apiErrors.From(apiErrors.MaxNumberOfUsersReached).Wrap(
			apiErrors.WithError(ErrMaxUserReached))
	}

	if scb.numUsers == 0 {
		scb.cfg.Host = u.ID
	}

	scb.muRdyUsr.Lock()
	scb.usersReady[u.ID] = u
	scb.muRdyUsr.Unlock()
	return nil
}

// GetUserReady returns the user with userID ready to join a session.
func (scb *controlBlock) getUserReady(userID string) (*User, error) {
	scb.muRdyUsr.RLock()
	defer scb.muRdyUsr.RUnlock()
	u, ok := scb.usersReady[userID]
	if !ok {
		return nil, ErrReadyUser
	}
	return u, nil
}

func (scb *controlBlock) UserCanJoin(userID string) error {
	scb.muRdyUsr.RLock()
	defer scb.muRdyUsr.RUnlock()
	if _, ok := scb.usersReady[userID]; !ok {
		return ErrReadyUser
	}
	scb.muUsr.RLock()
	defer scb.muUsr.RUnlock()
	if _, ok := scb.users[userID]; ok {
		return ErrUserConnected
	}
	if scb.numUsers >= scb.cfg.MaxUsers {
		return apiErrors.From(apiErrors.MaxNumberOfUsersReached).Wrap(
			apiErrors.WithError(ErrMaxUserReached))
	}

	return nil
}

// UserConnect adds user from the userReady state to clients.
//
// Broadcast that user has connected to session.
func (scb *controlBlock) UserConnect(userID string, conn *gws.Conn) error {
	u, err := scb.getUserReady(userID)
	if err != nil {
		return apiErrors.ErrBadRequest.Wrap(apiErrors.WithError(err))
	}
	u.Conn = conn

	scb.muUsr.Lock()
	if _, ok := scb.users[u.ID]; ok {
		scb.muUsr.Unlock()
		return apiErrors.ErrBadRequest.Wrap(apiErrors.WithErrorf("user already connected"))
	}
	scb.users[u.ID] = u
	scb.numUsers++
	numCl := scb.numUsers
	scb.muUsr.Unlock()

	// the first user to connect needs to start the session
	if numCl == 1 {
		scb.Start()
	}

	// broadcast that user has joined
	scb.broadcaster.Broadcast() <- types.Message{
		Type:    MessageTypeUserConnected,
		Content: u,
	}

	if scb.isHost(u) {
		scb.broadcaster.Send() <- types.Message{
			Type:     MessageTypeUserHost,
			Receiver: u.ID,
			Content:  userHostContent{Secret: scb.cfg.Secret},
		}
	}

	return nil
}

// UserDisconnect removes user from clients.
//
// Broadcast that user has disconnected from session.
func (scb *controlBlock) UserDisconnect(ctx context.Context, userID string) {
	scb.muUsr.Lock()
	u, ok := scb.users[userID]
	if ok {
		delete(scb.users, u.ID)
		scb.numUsers--
	}
	numCl := scb.numUsers
	scb.muUsr.Unlock()

	// if session is empty after client disconnect
	// the session needs to be set to inactive
	if numCl == 0 {
		_ = scb.dispatcher.Close(ctx, scb.cfg.ID)
		return
	}

	// broadcast that user has left
	scb.broadcaster.Broadcast() <- types.Message{
		Type:    MessageTypeUserDisconnected,
		Content: u,
	}

	scb.broadcaster.Control() <- types.Message{
		Receiver: userID,
		Content:  "Closed by server",
	}
}

func (scb *controlBlock) KickUser(userID string) error {
	if _, ok := scb.GetUsers()[userID]; !ok {
		return apiErrors.ErrBadRequest.Wrap(apiErrors.WithErrorf("user not found"))
	}

	// delete user from registered ones to prevent rejoin with same config
	scb.muRdyUsr.Lock()
	delete(scb.usersReady, userID)
	scb.muRdyUsr.Unlock()

	scb.broadcaster.Send() <- types.Message{
		Type:     MessageTypeUserKick,
		Receiver: userID,
	}
	scb.broadcaster.Control() <- types.Message{
		Receiver: userID,
		Content:  "Kicked by host",
	}
	return nil
}

// GetUsers returns all active users/clients in the session.
func (scb *controlBlock) GetUsers() map[string]*User {
	users := make(map[string]*User)
	scb.muUsr.RLock()
	for id, u := range scb.users {
		users[id] = u
	}
	scb.muUsr.RUnlock()
	return users
}

func (scb *controlBlock) isHost(u *User) bool {
	return u.ID == scb.cfg.Host
}
