package session

import (
	"context"

	"github.com/heat1q/boardsite/api/log"
	"github.com/heat1q/boardsite/api/types"
	"github.com/heat1q/boardsite/redis"
)

//counterfeiter:generate . Broadcaster
type Broadcaster interface {
	// Bind binds the broadcaster to a session
	Bind(scb Controller) Broadcaster
	// Broadcast returns a channel for messages to be broadcasted
	Broadcast() chan<- types.Message
	// Send returns a channel for messages to sent to a specific client
	Send() chan<- types.Message
	// Cache returns a channel for strokes to be stored in the cache
	Cache() chan<- []redis.Stroke
	// Close returns a channel for closing the broadcaster and clean up all goroutines
	Close() chan<- struct{}
}

type broadcaster struct {
	scb   Controller
	cache redis.Handler

	broadcast   chan types.Message
	send        chan types.Message
	cacheUpdate chan []redis.Stroke
	close       chan struct{}
}

// NewBroadcaster creates a new Broadcaster for a given session
func NewBroadcaster(cache redis.Handler) Broadcaster {
	return &broadcaster{
		cache:       cache,
		broadcast:   make(chan types.Message),
		send:        make(chan types.Message),
		cacheUpdate: make(chan []redis.Stroke),
		close:       make(chan struct{}),
	}
}

func (b *broadcaster) Bind(scb Controller) Broadcaster {
	if b.scb != nil {
		return nil
	}
	b.scb = scb
	// start goroutines for broadcasting and saving changes to board
	go b.broadcastLoop()
	go b.cacheUpdateLoop()

	return b
}

func (b *broadcaster) Broadcast() chan<- types.Message {
	return b.broadcast
}

func (b *broadcaster) Send() chan<- types.Message {
	return b.send
}

func (b *broadcaster) Cache() chan<- []redis.Stroke {
	return b.cacheUpdate
}

func (b *broadcaster) Close() chan<- struct{} {
	return b.close
}

func (b *broadcaster) getUsers() map[string]*User {
	if b.scb == nil {
		log.Global().Warnf("broadcaster is not bound to any session")
		return map[string]*User{}
	}
	return b.scb.GetUsers()
}

// broadcastLoop Broadcasts board updates to all clients
func (b *broadcaster) broadcastLoop() {
	for {
		select {
		case data := <-b.broadcast:
			users := b.getUsers()
			for userID, user := range users { // Send to all connected clients
				// except the origin, i.e. the initiator of message
				if userID != data.Sender {
					if err := user.Conn.WriteJSON(data); err != nil {
						log.Global().Warnf("cannot broadcast to %s: %v",
							user.Conn.RemoteAddr(), err)
					}
				}
			}
			break
		case data := <-b.send:
			users := b.getUsers()
			u, ok := users[data.Receiver]
			if !ok {
				log.Global().Warnf("broadcastLoop: unknown receiver %v", data.Receiver)
				continue
			}
			if err := u.Conn.WriteJSON(data); err != nil {
				log.Global().Warnf("broadcastLoop: %v", err)
			}
		case <-b.close:
			return
		}
	}
}

// dbUpdateLoop updates database according to given Stroke values
func (b *broadcaster) cacheUpdateLoop() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	for {
		select {
		case strokes := <-b.cacheUpdate:
			if err := b.cache.UpdateStrokes(ctx, b.scb.ID(), strokes...); err != nil {
				log.Global().Warnf("error in dbUpdateLoop: %v", err)
			}
		case <-b.close:
			_ = b.cache.ClearSession(ctx, b.scb.ID())
			return
		}
	}
}
