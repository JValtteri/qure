package middleware

import (
	"github.com/JValtteri/qure/server/internal/crypt"
	"github.com/JValtteri/qure/server/internal/state/model"
)

func reservationToResponse(res model.Reservation) ReservationResponse {
	errorMsg := res.Error
	if errorMsg != "" {
		return ReservationResponse {Error: errorMsg}
	} else {
		errorMsg = ""
	}
	return ReservationResponse {
		Id:         res.Id,
		EventID:	getReservationEventId(&res),
		ClientID:	res.Client,
		Size:       res.Size,
		Confirmed:  res.Confirmed,
		Timeslot:   res.Timeslot,
		Expiration: res.Expiration,
		Error:      errorMsg,
		Session:	res.Session,
		Event: 		getReservationEvent(&res),
	}
}

func getReservationEventId(res *model.Reservation) crypt.ID {
	if res.Event == nil {
		return "nil"
	}
	return res.Event.ID
}

func getReservationEvent(res *model.Reservation) Event {
	if res.Event == nil {
		return Event{}
	}
	return Event{
		ID:			res.Event.ID,
		Name:		res.Event.Name,
		DtStart:	res.Event.DtStart,
		DtEnd:		res.Event.DtEnd,
	}
}
