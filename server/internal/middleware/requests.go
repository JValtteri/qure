package middleware

import (
	"github.com/JValtteri/qure/server/internal/crypt"
	"github.com/JValtteri/qure/server/internal/state"
)

type Request interface {
	EventRequest |
	LoginRequest |
	EventLogin	 |
	AuthenticateRequest |
	RegisterRequest |
	ReserveRequest |
	UserReservationsRequest |
	UniversalRequest
}

type EventRequest struct {
	EventID		crypt.ID
}

type LoginRequest struct {
	User		string
	Password	crypt.Key
	Ip			state.IP
}

type EventLogin struct {
	EventID		crypt.Key
	Ip			state.IP
}

type AuthenticateRequest struct {
	SessionKey	crypt.Key	// This should come from the cookie
	Ip			state.IP	// This should be sensed by server
}

type RegisterRequest struct {
	User		string
	Password	crypt.Key
	Ip			state.IP
}

type ReserveRequest struct {
	SessionKey	crypt.Key	//
	Email		string
	Ip			state.IP	//
	Size		int
	EventId		state.ID
	Timeslot	state.Epoch
}

type UserReservationsRequest struct {
	SessionKey	crypt.Key
}

type UniversalRequest struct {
	Email		string
	Password	crypt.Key
	EventID		crypt.ID
	Size		int
	EventId		state.ID
	Imeslot		state.Epoch
	Ip			state.IP	// This should be sensed by server
	SessionKey	crypt.Key	// This should come from the cookie
}
