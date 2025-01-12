package actors

import (
	"go-chat/shared/ws"
	"sync"
	"time"

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
		// Start periodic client cleanup
		go r.periodicCleanup(ctx)

	case *actor.Stopped:
		log.Info("RoomActor stopped", "total_clients", len(r.clients))

	case *ClientJoined:
		r.mu.Lock()
		r.clients[msg.ClientPID.String()] = msg.ClientPID
		clientCount := len(r.clients)
		r.mu.Unlock()

		log.Info("client joined room",
			"pid", msg.ClientPID.String(),
			"total_clients", clientCount)

		// Broadcast join message to all clients
		joinMsg := &ws.Message{
			Type: ws.TypeJoin,
			Payload: &ws.JoinMessage{
				From: msg.Username,
			},
		}
		r.broadcastMessage(ctx, joinMsg)

	case *ClientLeft:
		r.mu.Lock()
		if _, exists := r.clients[msg.ClientPID.String()]; exists {
			delete(r.clients, msg.ClientPID.String())
			clientCount := len(r.clients)
			r.mu.Unlock()

			log.Info("client left room",
				"pid", msg.ClientPID.String(),
				"total_clients", clientCount)

			// Broadcast leave message
			leaveMsg := &ws.Message{
				Type: ws.TypeLeave,
				Payload: &ws.LeaveMessage{
					From: msg.Username,
				},
			}
			r.broadcastMessage(ctx, leaveMsg)
		} else {
			r.mu.Unlock()
		}

	case *ws.Message:
		r.broadcastMessage(ctx, msg)
	}
}

func (r *RoomActor) broadcastMessage(ctx *actor.Context, msg *ws.Message) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	log.Debug("broadcasting message",
		"type", msg.Type,
		"total_clients", len(r.clients))

	for pid, client := range r.clients {
		ctx.Engine().Send(client, msg)
		log.Debug("sent message to client", "pid", pid, "type", msg.Type)
	}
}

func (r *RoomActor) periodicCleanup(ctx *actor.Context) {
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		r.mu.Lock()
		initialCount := len(r.clients)

		// Verify each client is still alive
		for pid, client := range r.clients {
			ctx.Engine().Send(client, &ws.Message{Type: ws.TypePing})
			log.Debug("sent ping to client", "pid", pid)
		}

		finalCount := len(r.clients)
		r.mu.Unlock()

		if initialCount != finalCount {
			log.Info("cleanup complete",
				"initial_clients", initialCount,
				"final_clients", finalCount,
				"removed", initialCount-finalCount)
		}
	}
}
