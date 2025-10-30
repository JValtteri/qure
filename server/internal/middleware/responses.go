package middleware

import (
	"github.com/JValtteri/qure/server/internal/crypt"
	"github.com/JValtteri/qure/server/internal/state"
)


type Response interface {
	RegistrationResponse |
	Authentication |
	EventCreationResponse |
	state.Event |
	Reservation
}

type RegistrationResponse struct {
	SessionKey	crypt.Key
	Error		string
}

type EventCreationResponse struct {
	EventID	    crypt.ID
}

type Authentication struct {
    Authenticated bool
    IsAdmin       bool
    SessionKey    crypt.Key
    Error         string
}

type Reservation struct {
	Id			crypt.ID
	EventID		crypt.ID
	ClientID	crypt.ID
	Size		int				// Party size
	Confirmed	int				// Reserved size
	Timeslot	state.Epoch
	Expiration	state.Epoch
	Error		string
}
