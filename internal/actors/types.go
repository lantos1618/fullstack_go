package actors

import (
	"github.com/anthdm/hollywood/actor"
)

// Type represents different types of actors
type Type string

const (
	TypeClient Type = "client"
	TypeRoom   Type = "room"
)

// ClientJoined is sent when a new client joins the room
type ClientJoined struct {
	ClientPID *actor.PID
	Username  string
}

// ClientLeft is sent when a client leaves the room
type ClientLeft struct {
	ClientPID *actor.PID
	Username  string
}
