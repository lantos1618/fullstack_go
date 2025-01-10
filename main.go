package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"go-chat/shared"

	"github.com/anthdm/hollywood/actor"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // Allow all origins for demo
	},
}

type ClientActor struct {
	conn *websocket.Conn
	hub  *actor.PID
	id   string
}

type HubActor struct {
	clients map[*actor.PID]bool
	mu      sync.RWMutex
	engine  *actor.Engine
}

func NewHubActor() actor.Producer {
	return func() actor.Receiver {
		return &HubActor{
			clients: make(map[*actor.PID]bool),
		}
	}
}

func (h *HubActor) Receive(c *actor.Context) {
	switch msg := c.Message().(type) {
	case *shared.WSMessage:
		log.Printf("[Hub] Broadcasting message: Type=%s, From=%v", msg.Type, msg.Payload)
		h.broadcast(msg)
	case *actor.Started:
		h.engine = c.Engine()
		log.Println("[Hub] Actor started")
	case *actor.Stopped:
		log.Println("[Hub] Actor stopped")
	case *actor.PID:
		// Client disconnected
		h.mu.Lock()
		delete(h.clients, msg)
		h.mu.Unlock()
		log.Printf("[Hub] Client disconnected: %s", msg.ID)
		log.Printf("[Hub] Active clients: %d", len(h.clients))
	}
}

func (h *HubActor) broadcast(msg *shared.WSMessage) {
	h.mu.RLock()
	defer h.mu.RUnlock()
	for client := range h.clients {
		h.engine.Send(client, msg)
	}
}

func NewClientActor(conn *websocket.Conn, hub *actor.PID) actor.Producer {
	return func() actor.Receiver {
		return &ClientActor{
			conn: conn,
			hub:  hub,
		}
	}
}

func (c *ClientActor) Receive(ctx *actor.Context) {
	switch msg := ctx.Message().(type) {
	case *actor.Started:
		c.id = ctx.PID().ID
		log.Printf("[Client %s] Connected", c.id)
		go c.readPump(ctx.Engine())
	case *shared.WSMessage:
		log.Printf("[Client %s] Sending message: %+v", c.id, msg)
		c.writeMessage(msg)
	case *actor.Stopped:
		log.Printf("[Client %s] Disconnected", c.id)
		c.conn.Close()
		// Notify hub about disconnection
		ctx.Engine().Send(c.hub, ctx.PID())
	}
}

func (c *ClientActor) readPump(engine *actor.Engine) {
	defer c.conn.Close()

	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("[Client %s] Error: %v", c.id, err)
			}
			break
		}

		var wsMsg shared.WSMessage
		if err := json.Unmarshal(message, &wsMsg); err != nil {
			log.Printf("[Client %s] Error unmarshaling message: %v", c.id, err)
			continue
		}

		log.Printf("[Client %s] Received message: %+v", c.id, wsMsg)
		engine.Send(c.hub, &wsMsg)
	}
}

func (c *ClientActor) writeMessage(msg *shared.WSMessage) {
	data, err := json.Marshal(msg)
	if err != nil {
		log.Printf("[Client %s] Error marshaling message: %v", c.id, err)
		return
	}

	if err := c.conn.WriteMessage(websocket.TextMessage, data); err != nil {
		log.Printf("[Client %s] Error writing message: %v", c.id, err)
	}
}

func handleWS(e *actor.Engine, hub *actor.PID) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Println("[Server] Upgrade error:", err)
			return
		}

		clientProducer := NewClientActor(conn, hub)
		pid := e.Spawn(clientProducer, fmt.Sprintf("client-%d", time.Now().UnixNano()), actor.WithInboxSize(100))
		log.Printf("[Server] New client connected: %s", pid.ID)
	}
}

func main() {
	log.SetFlags(log.Ltime | log.Lmicroseconds)
	log.Println("[Server] Starting chat server...")

	config := actor.NewEngineConfig()
	engine, _ := actor.NewEngine(config)

	hubProducer := NewHubActor()
	hub := engine.Spawn(hubProducer, "hub", actor.WithInboxSize(1000))
	log.Printf("[Server] Hub actor started with ID: %s", hub.ID)

	http.HandleFunc("/ws", handleWS(engine, hub))
	http.Handle("/", http.FileServer(http.Dir("dist")))

	log.Println("[Server] Listening on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
