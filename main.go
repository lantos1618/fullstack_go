package main

import (
	"encoding/json"
	"fmt"
	"go-chat/shared"
	"net/http"
	"sync"
	"time"

	"github.com/anthdm/hollywood/actor"
	"github.com/charmbracelet/log"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // Allow all origins for demo
	},
}

type ActorType string

const (
	ActorTypeClient ActorType = "client"
	ActorTypeRoom   ActorType = "room"
)

type ClientActor struct {
	conn *websocket.Conn
}

// Receive implements actor.Receiver.
func (c *ClientActor) Receive(ctx *actor.Context) {
	switch msg := ctx.Message().(type) {
	case *actor.Started:
		log.Info("ClientActor started")
	case *actor.Stopped:
		log.Info("ClientActor stopped")
	case *shared.WSMessage:
		switch msg.Type {
		case shared.TypeMessage, shared.TypeTyping:
			err := c.conn.WriteJSON(msg)
			if err != nil {
				log.Error("failed to write message", "error", err)
			}
		case shared.TypePing:
			pongMsg := &shared.WSMessage{Type: shared.TypePong}
			err := c.conn.WriteJSON(pongMsg)
			if err != nil {
				log.Error("failed to send pong", "error", err)
			}
		}
	}
}

func newClientActor(conn *websocket.Conn) actor.Producer {
	return func() actor.Receiver {
		return &ClientActor{conn: conn}
	}
}

type RoomActor struct {
	clients map[string]*actor.PID
	mu      sync.RWMutex
}

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

	case *shared.WSMessage:
		// Broadcast both text messages and typing indicators
		if msg.Type == shared.TypeMessage || msg.Type == shared.TypeTyping {
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

type ClientJoined struct {
	ClientPID *actor.PID
}

type ClientLeft struct {
	ClientPID *actor.PID
}

func newRoomActor() actor.Producer {
	return func() actor.Receiver {
		return &RoomActor{
			clients: make(map[string]*actor.PID),
		}
	}
}

func main() {
	engine, err := actor.NewEngine(actor.NewEngineConfig())
	if err != nil {
		log.Fatal(err)
	}

	// Spawn the room actor
	roomPID := engine.Spawn(newRoomActor(), string(ActorTypeRoom))
	log.Info("Room actor started", "pid", roomPID)

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Fatal(err)
		}

		pid := engine.Spawn(newClientActor(conn), string(ActorTypeClient))

		// Notify room about new client
		engine.Send(roomPID, &ClientJoined{ClientPID: pid})

		log.Info(fmt.Sprintf("Client connected: %s", pid))

		defer func() {
			log.Info(fmt.Sprintf("Closing connection for %s", pid))
			// Notify room about client leaving
			engine.Send(roomPID, &ClientLeft{ClientPID: pid})
			engine.Send(pid, shared.TypeClose)
			engine.Poison(pid)
			conn.Close()
		}()

		// Start reading messages
		for {
			var msg shared.WSMessage
			err := conn.ReadJSON(&msg)
			if err != nil {
				log.Error("error reading message", "error", err)
				break
			}

			// Send message to room instead of directly to client
			engine.Send(roomPID, &msg)
		}
	})

	// dist
	http.Handle("/", http.FileServer(http.Dir("./dist")))

	// Health check endpoint
	http.HandleFunc(shared.RouteHealth.Path, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != shared.RouteHealth.Method {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		response := shared.HealthResponse{
			Status:    "ok",
			Timestamp: time.Now().Unix(),
		}
		json.NewEncoder(w).Encode(response)
	})

	http.ListenAndServe(":8080", nil)
	log.Info("Server started on port 8080")
}
