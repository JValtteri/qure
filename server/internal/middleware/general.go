package middleware

import (
	"log"

	"github.com/JValtteri/qure/server/internal/state"
	"github.com/JValtteri/qure/server/internal/state/model"
)


func GetEvents(rq EventRequest) []model.Event {
	isAdmin := checkStaffStatus(rq)
	events := state.GetEvents(isAdmin)
	for i := range(events) {
		events[i].LongDescription = ""
	}
	return events
}

func GetEvent(eventRequest EventRequest) model.Event {
	isAdmin := checkStaffStatus(eventRequest)
	event, err := state.GetEvent(eventRequest.EventID, isAdmin)
	if err != nil {
		log.Printf("Error getting event: %v\n", err)
	}
	return event
}

func GetUserReservatoions(rq UserReservationsRequest) []ReservationResponse {
	client, found := state.GetClientBySession(rq.SessionKey)
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

func checkStaffStatus(rq EventRequest) bool {
	var isStaff = false
	client, err := state.ResumeSession(rq.SessionKey, rq.Fingerprint)
	if err == nil {
		isStaff = client.IsStaff()
	}
	return isStaff
}
