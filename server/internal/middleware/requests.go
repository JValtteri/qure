package middleware

import (
	"github.com/JValtteri/qure/server/internal/crypt"
	"github.com/JValtteri/qure/server/internal/state"
)

type Request interface {
	UniversalRequest |
	EventRequest |
	LoginRequest |
	EventLogin	 |
	AuthenticateRequest |
	RegisterRequest |
	ReserveRequest |
	UserReservationsRequest |
	EventCreationRequest |
	UserEventRequest
}

type EventRequest struct {
	EventID		crypt.ID
	IsAdmin		bool
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
	SessionKey	crypt.Key
	Ip			state.IP
}

type RegisterRequest struct {
	User		string
	Password	crypt.Key
	Ip			state.IP
}

type ReserveRequest struct {
	SessionKey	crypt.Key
	User		string
	Ip			state.IP
	Size		int
	EventID		state.ID
	Timeslot	state.Epoch
}

type UserReservationsRequest struct {
	SessionKey	crypt.Key
}

type EventCreationRequest struct {
	SessionKey	crypt.Key
	Ip			state.IP
	Event		state.Event
}

type UserEventRequest struct {
	SessionKey	crypt.Key
}

type UniversalRequest struct {
	User		string
	Password	crypt.Key
	IsAdmin		bool
	EventID		crypt.ID
	Size		int
	Timeslot	state.Epoch
	Event		state.Event
	Ip			state.IP	// This is sensed by server
	SessionKey	crypt.Key	// This comes from the cookie
}
