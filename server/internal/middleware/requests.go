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
	Fingerprint	string
	HashPrint	crypt.Hash
}

type EventLogin struct {
	EventID		crypt.Key
	Fingerprint	string
	HashPrint	crypt.Hash
}

type AuthenticateRequest struct {
	SessionKey	crypt.Key
	Fingerprint	string
}

type RegisterRequest struct {
	User		string
	Password	crypt.Key
	HashPrint	crypt.Hash
}

type ReserveRequest struct {
	SessionKey	crypt.Key
	User		string
	Fingerprint	string
	HashPrint	crypt.Hash
	Size		int
	EventID		state.ID
	Timeslot	state.Epoch
}

type UserReservationsRequest struct {
	SessionKey	crypt.Key
}

type EventCreationRequest struct {
	SessionKey	crypt.Key
	Fingerprint	string
	Event		state.Event
}

type UserEventRequest struct {
	SessionKey	crypt.Key
}

type UniversalRequest struct {
	User			string
	Password		crypt.Key
	IsAdmin			bool
	EventID			crypt.ID
	Size			int
	Timeslot		state.Epoch
	Event			state.Event
	Fingerprint		string		// This is sensed by server
	HashPrint		crypt.Hash	// Hashed fingerprint
	SessionKey		crypt.Key	// This comes from the cookie
}
