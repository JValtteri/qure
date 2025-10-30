package middleware

import (
	"log"

	"github.com/JValtteri/qure/server/internal/crypt"
	"github.com/JValtteri/qure/server/internal/state"
)


func GetEvents(isAdmin bool) []state.Event {
	events := state.GetEvents(isAdmin)
	for i := range(events) {
		events[i].LongDescription = ""
	}
	return events
}

func GetEvent(eventRequest EventRequest) state.Event {
	event, err := state.GetEvent(eventRequest.EventID, eventRequest.IsAdmin)
	if err != nil {
		log.Printf("Error getting event: %v\n", err)
	}
	return event
}

func GetUserReservatoions(sessionKey crypt.Key) []*state.Reservation {
	client, found := state.GetClientBySession(sessionKey)
	if !found {
		return nil
	}
	res := client.GetReservations()
	return res
}

