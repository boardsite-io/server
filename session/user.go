package session

import (
	"context"
	"errors"
	"fmt"

	gws "github.com/gorilla/websocket"
	gonanoid "github.com/matoous/go-nanoid/v2"

	apiErrors "github.com/heat1q/boardsite/api/errors"
	"github.com/heat1q/boardsite/api/types"
)

const maxNameLen = 32

// User declares some information about connected users.
type User struct {
	ID    string    `json:"id"`
	Alias string    `json:"alias"`
	Color string    `json:"color"`
	Conn  *gws.Conn `json:"-"`
}

// NewUser generate a new user struct based on
// the alias and color attribute
//
// Does some sanitize checks.
func (scb *controlBlock) NewUser(alias, color string) (*User, error) {
	if len(alias) > maxNameLen {
		alias = alias[:maxNameLen]
	}
	//TODO check if html color ?
	if len(color) != 7 {
		return nil, fmt.Errorf("incorrect html color")
	}

	id, err := gonanoid.New(16)
	if err != nil {
		return nil, err
	}
	user := &User{
		ID:    id,
		Alias: alias,
		Color: color,
	}
	// set user waiting
	err = scb.userReady(user)
	return user, err
}

// UserReady adds an user to the usersReady map.
func (scb *controlBlock) userReady(u *User) error {
	scb.muUsr.RLock()
	defer scb.muUsr.RUnlock()
	if scb.numUsers >= scb.maxUsers {
		return apiErrors.From(apiErrors.CodeMaxNumberOfUsersReached).Wrap(
			apiErrors.WithErrorf("maximum number of connected users reached"))
	}

	scb.muRdyUsr.Lock()
	scb.usersReady[u.ID] = u
	scb.muRdyUsr.Unlock()
	return nil
}

// GetUserReady returns the user with userID ready to join a session.
func (scb *controlBlock) GetUserReady(userID string) (*User, error) {
	scb.muRdyUsr.RLock()
	defer scb.muRdyUsr.RUnlock()
	u, ok := scb.usersReady[userID]
	if !ok {
		return nil, errors.New("ready user not found")
	}
	return u, nil
}

// UserConnect adds user from the userReady state to clients.
//
// Broadcast that user has connected to session.
func (scb *controlBlock) UserConnect(u *User) {
	scb.muUsr.Lock()
	scb.users[u.ID] = u
	scb.numUsers++
	numCl := scb.numUsers

	// the first user to connect needs to start the session
	if numCl == 1 {
		scb.Start()
	}

	scb.muUsr.Unlock()

	// broadcast that user has joined
	scb.broadcaster.Broadcast() <- types.Message{
		Type:    MessageTypeUserConnected,
		Content: u,
	}
}

// UserDisconnect removes user from clients.
//
// Broadcast that user has disconnected from session.
func (scb *controlBlock) UserDisconnect(ctx context.Context, userID string) {
	scb.muUsr.Lock()
	u := scb.users[userID]
	delete(scb.users, u.ID)
	scb.numUsers--
	numCl := scb.numUsers
	scb.muUsr.Unlock()

	// if session is empty after client disconnect
	// the session needs to be set to inactive
	if numCl == 0 {
		_ = scb.dispatcher.Close(ctx, scb.id)
		return
	}

	// broadcast that user has left
	scb.broadcaster.Broadcast() <- types.Message{
		Type:    MessageTypeUserDisconnected,
		Content: u,
	}
}

// IsUserConnected checks if the user with userID is an active client in the session.
func (scb *controlBlock) isUserConnected(userID string) bool {
	scb.muUsr.RLock()
	defer scb.muUsr.RUnlock()
	_, ok := scb.users[userID]
	return ok
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
