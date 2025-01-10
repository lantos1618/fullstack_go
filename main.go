package main

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"

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
		h.broadcast(msg)
	case *actor.Started:
		h.engine = c.Engine()
		log.Println("Hub actor started")
	case *actor.Stopped:
		log.Println("Hub actor stopped")
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
		go c.readPump(ctx.Engine())
	case *shared.WSMessage:
		c.writeMessage(msg)
	case *actor.Stopped:
		c.conn.Close()
	}
}

func (c *ClientActor) readPump(engine *actor.Engine) {
	defer c.conn.Close()

	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}

		var wsMsg shared.WSMessage
		if err := json.Unmarshal(message, &wsMsg); err != nil {
			log.Printf("error unmarshaling message: %v", err)
			continue
		}

		engine.Send(c.hub, &wsMsg)
	}
}

func (c *ClientActor) writeMessage(msg *shared.WSMessage) {
	data, err := json.Marshal(msg)
	if err != nil {
		log.Printf("error marshaling message: %v", err)
		return
	}

	if err := c.conn.WriteMessage(websocket.TextMessage, data); err != nil {
		log.Printf("error writing message: %v", err)
	}
}

func handleWS(e *actor.Engine, hub *actor.PID) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Println("upgrade error:", err)
			return
		}

		clientProducer := NewClientActor(conn, hub)
		pid := e.Spawn(clientProducer, "client", actor.WithInboxSize(100))
		log.Printf("New client connected: %s", pid.ID)
	}
}

func main() {
	config := actor.NewEngineConfig()
	engine, _ := actor.NewEngine(config)

	hubProducer := NewHubActor()
	hub := engine.Spawn(hubProducer, "hub", actor.WithInboxSize(1000))

	http.HandleFunc("/ws", handleWS(engine, hub))
	http.Handle("/", http.FileServer(http.Dir("frontend")))

	log.Println("Server starting on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
