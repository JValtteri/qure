package middleware

import (
	"github.com/JValtteri/qure/server/internal/utils"
	"github.com/JValtteri/qure/server/internal/crypt"
	"github.com/JValtteri/qure/server/internal/state/model"
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
	EventManipulationRequest |
	UserEventRequest |
	PasswordChangeRequest |
	RemovalRequest
}

type EventRequest struct {
	EventID		crypt.ID
	IsAdmin		bool
}

type LoginRequest struct {
	User		string
	Password	crypt.Key
	HashPrint	crypt.Hash
}

type EventLogin struct {
	EventID		crypt.Key
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
	EventID		crypt.ID
	Timeslot	utils.Epoch
}

type UserReservationsRequest struct {
	SessionKey	crypt.Key
}

type EventManipulationRequest struct {
	SessionKey	crypt.Key
	Fingerprint	string
	EventID		crypt.ID
	Event		model.Event
}

type UserEventRequest struct {
	SessionKey	crypt.Key
}

type PasswordChangeRequest struct {
	User		string
	SessionKey	crypt.Key
	Fingerprint	string
	HashPrint	crypt.Hash
	Password	crypt.Key
	NewPassword	crypt.Key
}

type RemovalRequest struct {
	User		string
	SessionKey	crypt.Key
	Fingerprint	string
	HashPrint	crypt.Hash
	Password	crypt.Key
}

type UniversalRequest struct {
	User			string
	Password		crypt.Key
	NewPassword		crypt.Key
	IsAdmin			bool
	EventID			crypt.ID
	Size			int
	Timeslot		utils.Epoch
	Event			model.Event
	Fingerprint		string		// This is sensed by server
	HashPrint		crypt.Hash	// Hashed Fingerprint
	SessionKey		crypt.Key	// This comes from the cookie
}


type Identity interface {
	SessionKey()	crypt.Key
	Fingerprint()	string
}
