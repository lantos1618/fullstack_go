package main

import (
	"encoding/json"
	"fmt"
	"go-chat/internal/actors"
	"go-chat/shared/api"
	"go-chat/shared/ws"
	"net/http"
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

func setupWebSocket(engine *actor.Engine, roomPID *actor.PID) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Fatal(err)
		}

		pid := engine.Spawn(actors.NewClient(conn), string(actors.TypeClient))

		// Notify room about new client
		engine.Send(roomPID, &actors.ClientJoined{ClientPID: pid})

		log.Info(fmt.Sprintf("Client connected: %s", pid))

		defer func() {
			log.Info(fmt.Sprintf("Closing connection for %s", pid))
			// Notify room about client leaving
			engine.Send(roomPID, &actors.ClientLeft{ClientPID: pid})
			engine.Send(pid, ws.TypeClose)
			engine.Poison(pid)
			conn.Close()
		}()

		// Start reading messages
		for {
			var msg ws.Message
			err := conn.ReadJSON(&msg)
			if err != nil {
				log.Error("error reading message", "error", err)
				break
			}

			// Send message to room instead of directly to client
			engine.Send(roomPID, &msg)
		}
	}
}

func setupHealthCheck() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != api.RouteHealth.Method {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		response := api.HealthResponse{
			Status:    "ok",
			Timestamp: time.Now().Unix(),
		}
		json.NewEncoder(w).Encode(response)
	}
}

func main() {
	// Initialize actor system
	engine, err := actor.NewEngine(actor.NewEngineConfig())
	if err != nil {
		log.Fatal(err)
	}

	// Spawn the room actor
	roomPID := engine.Spawn(actors.NewRoom(), string(actors.TypeRoom))
	log.Info("Room actor started", "pid", roomPID)

	// Setup routes
	http.HandleFunc("/ws", setupWebSocket(engine, roomPID))
	http.HandleFunc(api.RouteHealth.Path, setupHealthCheck())
	http.Handle("/", http.FileServer(http.Dir("./dist")))

	// Start server
	log.Info("Server starting on port 8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
