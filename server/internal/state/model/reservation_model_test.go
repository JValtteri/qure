package model

import (
	"testing"

	"github.com/JValtteri/qure/server/internal/crypt"
	"github.com/JValtteri/qure/server/internal/utils"
)


func TestReservations(t *testing.T) {
	// Setup init objects
	reservations := makeTestReservations()
	clients := getTestClients()
	client := getTestClient()
	clients.ByEmail[client.Email] = &client
	clients.ByID[client.Id] = &client

	// Add Event
	event, slot := getTestEvent()
	slot.Size = 2
	time := utils.Epoch(200)
	event.Append(slot, time)

	// Test Reservation
	res := getTestReservation()
	res.Client = client.Id
	res.Event = &event
	err := res.Register(&reservations, &clients)
	if err != nil {
		t.Fatalf("Validating reservation failed: %s\n", err)
	}

	// Test Amending reservation
	var oldSize = res.Size
	var newSize = res.Size+1
	if oldSize+1 != newSize {
		t.Fatalf("Expected: %v, Got: %v\n", 2, newSize)
	}
	res.Size = newSize
	var oldReservations = len(res.getTimeslot().Reservations)
	var oldQueue = len(res.getTimeslot().Queue)

	//////////////////////////////

	var old = reservations.ByID[res.Id]
	var oldSlot = old.getTimeslot()
	t.Logf("%v", len(oldSlot.Queue))
	t.Logf("%v", len(oldSlot.Queue))

	///////////////////////////////

	err = res.Amend(&reservations, &clients)					// Amend /////  Here
	if err != nil {
		t.Fatalf("Amending reservation failed: %s\n", err)
	}
	var newReservations = len(res.getTimeslot().Reservations)
	var newQueue = len(res.getTimeslot().Queue)
	if oldReservations != 1 {
		t.Fatalf("Expected: %v, Got: %v\n", 1, oldReservations)
	}

	//r.Size - oldReservation.Confirmed - timeslot.countInQueue(r.Id)
	t.Logf("Size:%v Confirmed:%v oldQueue:%v newQueue:%v\n", res.Size, reservations.ByID[res.Event.ID].Confirmed, oldQueue, newQueue)

	if newReservations != oldReservations+1 {
		t.Fatalf("Expected: %v, Got: %v\n", oldReservations+1, newReservations)
	}
	if oldQueue != 0 {
		t.Fatalf("Expected: %v, Got: %v (oldQueue)\n", 0, oldQueue)
	}
	if newQueue != 0 {
		t.Fatalf("Expected: %v, Got: %v (oldQueue)\n", 0, newQueue)
	}
}


func getTestReservation() Reservation {
	idtype := crypt.ID("")
	id, _ := crypt.CreateHumanReadableKey(&idtype, 20)
	return Reservation{
		Id: 		id,
		Size: 		1,
		Timeslot:	utils.Epoch(200),
	}
}

func makeTestReservations() Reservations {
	return Reservations {
		ByID: 		make(map[crypt.ID]Reservation),
		ByEmail:	make(map[string]*Reservation),
	}
}
