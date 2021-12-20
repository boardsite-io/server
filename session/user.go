package session

import (
	"context"
	"errors"
	"fmt"

	gonanoid "github.com/matoous/go-nanoid/v2"

	"github.com/heat1q/boardsite/api/types"
)

// NewUser generate a new user struct based on
// the alias and color attribute
//
// Does some sanitize checks.
func (scb *ControlBlock) NewUser(alias, color string) (*types.User, error) {
	if len(alias) > 24 {
		alias = alias[:24]
	}
	//TODO check if html color ?
	if len(color) != 7 {
		return nil, fmt.Errorf("incorrect html color")
	}

	id, err := gonanoid.New(16)
	if err != nil {
		return nil, err
	}
	user := &types.User{
		ID:    id,
		Alias: alias,
		Color: color,
	}
	// set user waiting
	err = scb.UserReady(user)
	return user, err
}

// UserReady adds an user to the usersReady map.
func (scb *ControlBlock) UserReady(u *types.User) error {
	scb.muUsr.RLock()
	defer scb.muUsr.RUnlock()
	if scb.numUsers >= scb.maxUsers {
		return errors.New("maximum number of connected users in session has been reached")
	}

	scb.muRdyUsr.Lock()
	scb.usersReady[u.ID] = u
	scb.muRdyUsr.Unlock()
	return nil
}

// GetUserReady returns the user with userID ready to join a session.
func (scb *ControlBlock) GetUserReady(userID string) (*types.User, error) {
	scb.muRdyUsr.RLock()
	defer scb.muRdyUsr.RUnlock()
	u, ok := scb.usersReady[userID]
	if !ok {
		return nil, errors.New("ready user not found")
	}
	return u, nil
}

// IsUserReady checks if the user with userID is ready to join a session.
func (scb *ControlBlock) IsUserReady(userID string) bool {
	_, err := scb.GetUserReady(userID)
	return err == nil
}

// UserConnect adds user from the userReady state to clients.
//
// Broadcast that user has connected to session.
func (scb *ControlBlock) UserConnect(u *types.User) {
	scb.muUsr.Lock()
	scb.users[u.ID] = u
	scb.numUsers++
	scb.muUsr.Unlock()

	// broadcast that user has joined
	scb.broadcast <- &types.Message{
		Type:    types.MessageTypeUserConnected,
		Content: u,
	}
}

// UserDisconnect removes user from clients.
//
// Broadcast that user has disconnected from session.
func (scb *ControlBlock) UserDisconnect(ctx context.Context, userID string) {
	scb.muUsr.Lock()
	u := scb.users[userID]
	delete(scb.users, u.ID)
	scb.numUsers--
	numCl := scb.numUsers
	scb.muUsr.Unlock()

	// if session is empty after client disconnect
	// the session needs to be set to inactive
	if numCl == 0 {
		_ = scb.Dispatcher.Close(ctx, scb.ID)
		return
	}

	// broadcast that user has left
	scb.broadcast <- &types.Message{
		Type:    types.MessageTypeUserDisconnected,
		Content: u,
	}
}

// IsUserConnected checks if the user with userID is an active client in the session.
func (scb *ControlBlock) IsUserConnected(userID string) bool {
	scb.muUsr.RLock()
	defer scb.muUsr.RUnlock()
	_, ok := scb.users[userID]
	return ok
}

// GetUsers returns all active users/clients in the session.
func (scb *ControlBlock) GetUsers() map[string]*types.User {
	users := make(map[string]*types.User)
	scb.muUsr.RLock()
	for id, u := range scb.users {
		users[id] = u
	}
	scb.muUsr.RUnlock()
	return users
}
