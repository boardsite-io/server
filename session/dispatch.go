package session

import (
	"errors"
	"log"
	"sync"

	gonanoid "github.com/matoous/go-nanoid/v2"

	"github.com/heat1q/boardsite/redis"
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

type sessionsDispatcher struct {
	mu            sync.RWMutex
	activeSession map[string]*ControlBlock
	cache         redis.Handler
}

func NewDispatcher(cache redis.Handler) Dispatcher {
	return &sessionsDispatcher{
		activeSession: make(map[string]*ControlBlock),
		cache:         cache,
	}
}

func (d *sessionsDispatcher) GetSCB(sessionID string) (*ControlBlock, error) {
	d.mu.RLock()
	defer d.mu.RUnlock()
	scb, ok := d.activeSession[sessionID]
	if !ok {
		return nil, errors.New("session not found")
	}
	return scb, nil
}

func (d *sessionsDispatcher) Create() (string, error) {
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

	scb := NewControlBlock(sid, d.cache, d)
	// assign to SessionControl struct
	d.mu.Lock()
	d.activeSession[scb.ID] = scb
	d.mu.Unlock()
	log.Printf("Create Session with ID: %s\n", scb.ID)

	return sid, nil
}

func (d *sessionsDispatcher) Close(sessionID string) error {
	scb, err := d.GetSCB(sessionID)
	if err != nil {
		return err
	}
	scb.Close()
	d.mu.Lock()
	delete(d.activeSession, sessionID)
	d.mu.Unlock()

	if err := scb.Attachments.Clear(); err != nil {
		log.Printf("cannot clear attachments for %s: %v\n", scb.ID, err)
	}

	log.Printf("Close session %s", scb.ID)

	return nil
}

func (d *sessionsDispatcher) IsValid(sessionID string) bool {
	_, err := d.GetSCB(sessionID)
	return err == nil
}
