package main

import (
	"fmt"
	"go-chat/shared"
	"net/http"

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
		case shared.TypeMessage:
			if textMsg, ok := msg.Payload.(map[string]interface{}); ok {
				log.Info("received text message", "text", textMsg["text"], "from", textMsg["from"])
				// Echo the message back to the client for now
				err := c.conn.WriteJSON(msg)
				if err != nil {
					log.Error("failed to write message", "error", err)
				}
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

func main() {

	engine, err := actor.NewEngine(actor.NewEngineConfig())
	if err != nil {
		log.Fatal(err)
	}

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Fatal(err)
		}

		pid := engine.Spawn(newClientActor(conn), string(ActorTypeClient))

		log.Info(fmt.Sprintf("Client connected: %s", pid))

		defer func() {
			log.Info(fmt.Sprintf("Closing connection for %s", pid))
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

			engine.Send(pid, &msg)
		}
	})

	// dist
	http.Handle("/", http.FileServer(http.Dir("./dist")))

	http.ListenAndServe(":8080", nil)
	log.Info("Server started on port 8080")
}
