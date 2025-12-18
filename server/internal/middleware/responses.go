package middleware

import (
	"github.com/JValtteri/qure/server/internal/utils"
	"github.com/JValtteri/qure/server/internal/crypt"
	"github.com/JValtteri/qure/server/internal/state/model"
)


type Response interface {
	RegistrationResponse |
	Authentication |
	EventManipulationResponse |
	model.Event |
	Reservation |
	Reservations |
	SuccessResponse |
	[]model.Event
}

type RegistrationResponse struct {
	SessionKey	crypt.Key
	Error		string
}

type EventManipulationResponse struct {
	EventID	    crypt.ID
	Error		string
}

type Authentication struct {
	User			string
	Authenticated	bool
	IsAdmin			bool
	SessionKey		crypt.Key
	Error			string
}

type Reservation struct {
	Id			crypt.ID
	EventID		crypt.ID
	ClientID	crypt.ID
	Size		int				// Party size
	Confirmed	int				// Reserved size
	Timeslot	utils.Epoch
	Expiration	utils.Epoch
	Error		string
}

type SuccessResponse struct {
	Success		bool
	Error		string
}

type Reservations struct {
	Reservations	[]Reservation
}
