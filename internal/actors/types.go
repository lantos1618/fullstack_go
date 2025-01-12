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

// ClientJoined represents a client joining message
type ClientJoined struct {
	ClientPID *actor.PID
}

// ClientLeft represents a client leaving message
type ClientLeft struct {
	ClientPID *actor.PID
}
