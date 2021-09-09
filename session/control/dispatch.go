package session

import (
	"errors"
	"log"
	"sync"

	gonanoid "github.com/matoous/go-nanoid/v2"
)

const (
	alphabet = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
)

type Dispatcher interface {
	// GetSCB returns the session control block for given sessionID.
	GetSCB(sessionID string) (*ControlBlock, error)
	// Create creates and initializes a new SessionControl struct
	Create() (string, error)
	// Close removes the SCB from the activesession map and closes the session.
	Close(sessionID string) error
	// IsValid checks if session with sessionID exists.
	IsValid(sessionID string) bool
}

type sessionDispatcher struct {
	mu            sync.RWMutex
	activeSession map[string]*ControlBlock
}

func NewDispatcher() Dispatcher {
	return &sessionDispatcher{
		activeSession: make(map[string]*ControlBlock),
	}
}

func (d *sessionDispatcher) GetSCB(sessionID string) (*ControlBlock, error) {
	d.mu.RLock()
	defer d.mu.RUnlock()
	scb, ok := d.activeSession[sessionID]
	if !ok {
		return nil, errors.New("session not found")
	}
	return scb, nil
}

func (d *sessionDispatcher) Create() (string, error) {
	var sid string
	for {
		id, err := gonanoid.Generate(alphabet, 8)
		if err != nil {
			return "", err
		}
		// ensure uniqueness of id
		if _, err := d.GetSCB(id); err != nil {
			sid = id
			break
		}
	}

	scb := NewControlBlock(sid)
	// assign to SessionControl struct
	d.mu.Lock()
	d.activeSession[scb.ID] = scb
	d.mu.Unlock()
	log.Printf("Create Session with ID: %s\n", scb.ID)

	return sid, nil
}

func (d *sessionDispatcher) Close(sessionID string) error {
	scb, err := d.GetSCB(sessionID)
	if err != nil {
		return err
	}
	scb.SignalClose <- struct{}{}
	d.mu.Lock()
	delete(d.activeSession, sessionID)
	d.mu.Unlock()

	if err := scb.Attachments.Clear(); err != nil {
		log.Printf("cannot clear attachments for %s: %v\n", scb.ID, err)
	}

	log.Printf("Close session %s", scb.ID)

	return nil
}

func (d *sessionDispatcher) IsValid(sessionID string) bool {
	_, err := d.GetSCB(sessionID)
	return err == nil
}
