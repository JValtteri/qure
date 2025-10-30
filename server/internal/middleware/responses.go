package middleware

import (
	"github.com/JValtteri/qure/server/internal/crypt"
	"github.com/JValtteri/qure/server/internal/state"
)


type Response interface {
    RegistrationResponse |
    Authentication |
    EventCreationResponse |
    state.Event
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
