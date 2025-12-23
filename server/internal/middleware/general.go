package middleware

import (
	"log"

	"github.com/JValtteri/qure/server/internal/state"
	"github.com/JValtteri/qure/server/internal/state/model"
)


func GetEvents(rq EventRequest) []model.Event {
	isAdmin := checkAdminStatus(rq)
	events := state.GetEvents(isAdmin)
	for i := range(events) {
		events[i].LongDescription = ""
	}
	return events
}

func GetEvent(eventRequest EventRequest) model.Event {
	isAdmin := checkAdminStatus(eventRequest)
	event, err := state.GetEvent(eventRequest.EventID, isAdmin)
	if err != nil {
		log.Printf("Error getting event: %v\n", err)
	}
	return event
}

func GetUserReservatoions(req UserReservationsRequest) []ReservationResponse {
	client, found := state.GetClientBySession(req.SessionKey)
	var response []ReservationResponse
	if !found {
		return response
	}
	reservations := client.GetReservations()
	for _, value := range(reservations) {
		response = append(response, reservationToResponse(*value))
	}
	return response
}

func checkAdminStatus(rq EventRequest) bool {
	var isAdmin = false
	client, err := state.ResumeSession(rq.SessionKey, rq.Fingerprint)
	if err == nil {
		isAdmin = client.IsAdmin()
	}
	return isAdmin
}
