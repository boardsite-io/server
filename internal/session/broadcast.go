package session

import (
	"context"
	"errors"
	"fmt"
	"runtime"

	gws "github.com/gorilla/websocket"

	"github.com/boardsite-io/server/pkg/log"
	"github.com/boardsite-io/server/pkg/redis"
)

var ErrBroadcasterClosed = errors.New("broadcaster: closed")

//counterfeiter:generate . Broadcaster
type Broadcaster interface {
	// Bind binds the broadcaster to a session
	Bind(scb Controller) Broadcaster
	// Broadcast returns a channel for messages to be broadcasted
	Broadcast() chan<- Message
	// Send returns a channel for messages to sent to a specific client
	Send() chan<- Message
	// Control returns a channel for close messages to sent to a specific client
	Control() chan<- Message
	// Cache returns a channel for strokes to be stored in the cache
	Cache() chan<- []redis.Stroke
	// Close the broadcaster and cleans up all goroutines
	Close()
}

type broadcaster struct {
	scb   Controller
	cache redis.Handler

	broadcast   chan Message
	send        chan Message
	control     chan Message
	cacheUpdate chan []redis.Stroke
	close       chan struct{}
}

// NewBroadcaster creates a new Broadcaster for a given session
func NewBroadcaster(cache redis.Handler) Broadcaster {
	return &broadcaster{
		cache:       cache,
		broadcast:   make(chan Message),
		send:        make(chan Message),
		control:     make(chan Message),
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

func (b *broadcaster) Broadcast() chan<- Message {
	return b.broadcast
}

func (b *broadcaster) Send() chan<- Message {
	return b.send
}

func (b *broadcaster) Control() chan<- Message {
	return b.control
}

func (b *broadcaster) Cache() chan<- []redis.Stroke {
	return b.cacheUpdate
}

func (b *broadcaster) Close() {
	close(b.close)
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
		users := b.getUsers()
		err := b.broadcastToUser(users)
		if errors.Is(err, ErrBroadcasterClosed) {
			return
		}
		if err != nil {
			log.Global().Warnf("broadcastLoop: %v", err)
		}
	}
}

func (b *broadcaster) broadcastToUser(users map[string]*User) error {
	defer func() {
		if r := recover(); r != nil {
			stack := make([]byte, 2<<12)
			length := runtime.Stack(stack, true)
			log.Global().Errorf("[PANIC RECOVER] %v %s", r, stack[:length])
		}
	}()

	select {
	case data := <-b.broadcast:
		for userID, user := range users { // Send to all connected clients
			// except the origin, i.e. the initiator of message
			if userID != data.Sender {
				if err := user.Conn.WriteJSON(data); err != nil {
					log.Global().Warnf("cannot broadcast to %s: %v",
						user.Conn.RemoteAddr(), err)
				}
			}
		}
	case data := <-b.send:
		u, ok := users[data.Receiver]
		if !ok {
			return fmt.Errorf("send: unkown receiver: %v", data.Receiver)
		}
		if err := u.Conn.WriteJSON(data); err != nil {
			return fmt.Errorf("send: writeJSON: %w", err)
		}
	case data := <-b.control:
		u, ok := users[data.Receiver]
		if !ok {
			return fmt.Errorf("control: unkown receiver: %v", data.Receiver)
		}
		msg := gws.FormatCloseMessage(gws.CloseNormalClosure, fmt.Sprintf("%v", data.Content))
		_ = u.Conn.WriteMessage(gws.CloseMessage, msg)
	case <-b.close:
		return ErrBroadcasterClosed
	}
	return nil
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
