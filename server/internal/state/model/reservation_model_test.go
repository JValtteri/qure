package model

import (
	"testing"

	"github.com/JValtteri/qure/server/internal/crypt"
	"github.com/JValtteri/qure/server/internal/utils"
)


func TestReservations(t *testing.T) {
	reservations := makeTestReservations()
	res := getTestReservation()
	clients := getTestClients()
	client := getTestClient()
	clients.ByEmail[client.Email] = &client
	clients.ByID[client.Id] = &client
	event, slot := getTestEvent()
	time := utils.Epoch(200)
	event.Append(slot, time)
	res.Client = client.Id
	res.Event = &event
	err := res.Validate(&reservations, &clients)
	if err != nil {
		t.Fatalf("Validating reservation failed: %s\n", err)
	}
}


func getTestReservation() Reservation {
	return Reservation{
		Id: crypt.ID("test reservation"),
		Size: 1,
		Timeslot:	utils.Epoch(200),
	}
}

func makeTestReservations() Reservations {
	return Reservations {
		ByID: 		make(map[crypt.ID]Reservation),
		ByEmail:	make(map[string]*Reservation),
	}
}
