package actors

import (
	"go-chat/shared/ws"

	"github.com/anthdm/hollywood/actor"
	"github.com/charmbracelet/log"
	"github.com/gorilla/websocket"
)

// ClientActor handles individual WebSocket client connections
type ClientActor struct {
	conn *websocket.Conn
}

// NewClient creates a new client actor producer
func NewClient(conn *websocket.Conn) actor.Producer {
	return func() actor.Receiver {
		return &ClientActor{conn: conn}
	}
}

// Receive implements actor.Receiver
func (c *ClientActor) Receive(ctx *actor.Context) {
	switch msg := ctx.Message().(type) {
	case *actor.Started:
		log.Info("ClientActor started")
	case *actor.Stopped:
		log.Info("ClientActor stopped")
	case *ws.Message:
		switch msg.Type {
		case ws.TypeMessage, ws.TypeTyping:
			err := c.conn.WriteJSON(msg)
			if err != nil {
				log.Error("failed to write message", "error", err)
			}
		case ws.TypePing:
			pongMsg := &ws.Message{Type: ws.TypePong}
			err := c.conn.WriteJSON(pongMsg)
			if err != nil {
				log.Error("failed to send pong", "error", err)
			}
		}
	}
}
