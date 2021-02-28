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

var (
	mu            sync.RWMutex
	activeSession = make(map[string]*ControlBlock)
)

func GetSCB(sessionID string) (*ControlBlock, error) {
	mu.RLock()
	defer mu.RUnlock()
	scb, ok := activeSession[sessionID]
	if !ok {
		return nil, errors.New("session not found")
	}
	return scb, nil
}

// IsValid checks if session with sessionID exists.
func IsValid(sessionID string) bool {
	_, err := GetSCB(sessionID)
	return err == nil
}

// Create creates and initializes a new SessionControl struct
func Create() (string, error) {
	var sid string
	for {
		id, err := gonanoid.Generate(alphabet, 8)
		if err != nil {
			return "", err
		}
		// ensure uniqueness of id
		if _, err := GetSCB(id); err != nil {
			sid = id
			break
		}
	}

	scb := NewControlBlock(sid)
	// assign to SessionControl struct
	mu.Lock()
	activeSession[scb.ID] = scb
	mu.Unlock()
	log.Printf("Create Session with ID: %s\n", scb.ID)

	return sid, nil
}

// Close removes the SCB from the activesession map and closes the session.
func Close(sessionID string) error {
	scb, err := GetSCB(sessionID)
	if err != nil {
		return err
	}
	scb.SignalClose <- struct{}{}
	mu.Lock()
	delete(activeSession, sessionID)
	mu.Unlock()

	return nil
}
