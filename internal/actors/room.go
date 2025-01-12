package actors

import (
	"go-chat/shared/ws"
	"sync"

	"github.com/anthdm/hollywood/actor"
	"github.com/charmbracelet/log"
)

// RoomActor manages a group of connected clients
type RoomActor struct {
	clients map[string]*actor.PID
	mu      sync.RWMutex
}

// NewRoom creates a new room actor producer
func NewRoom() actor.Producer {
	return func() actor.Receiver {
		return &RoomActor{
			clients: make(map[string]*actor.PID),
		}
	}
}

// Receive implements actor.Receiver
func (r *RoomActor) Receive(ctx *actor.Context) {
	switch msg := ctx.Message().(type) {
	case *actor.Started:
		log.Info("RoomActor started")

	case *actor.Stopped:
		log.Info("RoomActor stopped")

	case *ClientJoined:
		r.mu.Lock()
		defer r.mu.Unlock()
		if r.clients == nil {
			r.clients = make(map[string]*actor.PID)
		}
		r.clients[msg.ClientPID.String()] = msg.ClientPID
		log.Info("client joined room", "pid", msg.ClientPID.String(), "total_clients", len(r.clients))

	case *ClientLeft:
		r.mu.Lock()
		defer r.mu.Unlock()
		if r.clients != nil {
			delete(r.clients, msg.ClientPID.String())
			log.Info("client left room", "pid", msg.ClientPID.String(), "total_clients", len(r.clients))
		}

	case *ws.Message:
		// Broadcast both text messages and typing indicators
		if msg.Type == ws.TypeMessage || msg.Type == ws.TypeTyping {
			r.mu.RLock()
			defer r.mu.RUnlock()

			if r.clients == nil {
				log.Error("clients map is nil")
				return
			}

			for pid, client := range r.clients {
				ctx.Engine().Send(client, msg)
				log.Debug("sent message to client", "pid", pid)
			}
		}
	}
}
