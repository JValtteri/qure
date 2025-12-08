package middleware

import (
	"log"

	"github.com/JValtteri/qure/server/internal/state"
	"github.com/JValtteri/qure/server/internal/state/model"
)


func GetEvents(isAdmin bool) []model.Event {
	events := state.GetEvents(isAdmin)
	for i := range(events) {
		events[i].LongDescription = ""
	}
	return events
}

func GetEvent(eventRequest EventRequest) model.Event {
	event, err := state.GetEvent(eventRequest.EventID, eventRequest.IsAdmin)
	if err != nil {
		log.Printf("Error getting event: %v\n", err)
	}
	return event
}

func GetUserReservatoions(req UserReservationsRequest) Reservations {
	client, found := state.GetClientBySession(req.SessionKey)
	if !found {
		return Reservations{}
	}
	reservations := client.GetReservations()
	var response []Reservation
	for _, value := range(reservations) {
		response = append(response, reservationToResponse(*value))
	}
	return Reservations{Reservations: response}
}
