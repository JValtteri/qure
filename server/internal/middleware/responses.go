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
	ReservationResponse |
	SuccessResponse |
	PasswordChangeResponse |
	[]model.Event |
	[]ReservationResponse
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

type ReservationResponse struct {
	Id			crypt.ID
	EventID		crypt.ID
	ClientID	crypt.ID
	Size		int				// Party size
	Confirmed	int				// Reserved size
	Timeslot	utils.Epoch
	Expiration	utils.Epoch
	Error		string
	Event 		Event
	Session		crypt.Key
}

type SuccessResponse struct {
	Success		bool
	Error		string
}

type PasswordChangeResponse struct {
	Success		bool
	SessionKey	crypt.Key
	Error		string
}


type Event struct {
	ID		crypt.ID
	Name	string
	DtStart	utils.Epoch
	DtEnd	utils.Epoch
}
