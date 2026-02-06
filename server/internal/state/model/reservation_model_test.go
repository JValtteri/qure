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
	// 1=1+0
	res := getTestReservation()
	res.Client = client.Id
	res.Event = &event
	err := res.Register(&reservations, &clients)
	if err != nil {
		t.Fatalf("Validating reservation failed: %s\n", err)
	}

	// Test Amending reservation
	// 2=2+0
	firstAmendTest(t, &res, &reservations, &clients)

	// Test Amend again
	// Should flow to queue
	// 3=2+1
	secondAmendTest(t, &res, &reservations, &clients)

	// Test equavilant reservation
	amendNoChangeTest(t, &res, &reservations, &clients)

	// Test smaller reservation
	// Should remove one slot from queue
	// 2=2+0
	firstReduceAmendTest(t, &res, &reservations, &clients)

	// Test smaller reservation
	// Should remove one slot from reservations
	// 1=1+0
	secondReduceAmendTest(t, &res, &reservations, &clients)
}

func TestPromoteFromQueue(t *testing.T) {
	// TODO:
	t.Log("TestPromoteFromQueue not implemented!")
}


func firstAmendTest(t *testing.T, res *Reservation, reservations *Reservations, clients *Clients) {
	var newSize = 2
	res.Size = newSize
	var oldReservations = len(res.getTimeslot().Reservations)
	var oldQueue = len(res.getTimeslot().Queue)
	var err = res.Amend(reservations, clients) // Amend /////  Here
	if err != nil {
		t.Fatalf("Amending reservation failed: %s\n", err)
	}
	var newReservations = len(res.getTimeslot().Reservations)
	var newQueue = len(res.getTimeslot().Queue)
	if oldReservations != 1 {
		t.Fatalf("Expected: %v, Got: %v\n", 1, oldReservations)
	}
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

func secondAmendTest(t *testing.T, res *Reservation, reservations *Reservations, clients *Clients) {
	res.Size = 3
	var oldQueue = len(res.getTimeslot().Queue)
	var err = res.Amend(reservations, clients) // Amend /////  Here
	if err != nil {
		t.Fatalf("Amending reservation failed: %s\n", err)
	}
	var newReservations = len(res.getTimeslot().Reservations)
	var newQueue = len(res.getTimeslot().Queue)
	if newReservations != 2 {
		t.Fatalf("Expected: %v, Got: %v\n", 2, newReservations)
	}
	if oldQueue != 0 {
		t.Fatalf("Expected: %v, Got: %v (oldQueue)\n", 0, oldQueue)
	}
	if newQueue != 1 {
		t.Fatalf("Expected: %v, Got: %v (oldQueue)\n", 1, newQueue)
	}
}

func amendNoChangeTest(t *testing.T, res *Reservation, reservations *Reservations, clients *Clients) {
	var err = res.Amend(reservations, clients) // Amend /////  Here
	if err == nil {
		t.Fatalf("Ammend should fail when no change is made:\n")
	}
}

func firstReduceAmendTest(t *testing.T, res *Reservation, reservations *Reservations, clients *Clients) {
	res.Size = 2
	var oldQueue = len(res.getTimeslot().Queue)
	var err = res.Amend(reservations, clients) // Amend /////  Here
	if err != nil {
		t.Fatalf("Amending reservation failed: %s\n", err)
	}
	var newReservations = len(res.getTimeslot().Reservations)
	var newQueue = len(res.getTimeslot().Queue)
	if newReservations != 2 {
		t.Fatalf("Expected: %v, Got: %v\n", 2, newReservations)
	}
	if oldQueue != 1 {
		t.Fatalf("Expected: %v, Got: %v (oldQueue)\n", 0, oldQueue)
	}
	if newQueue != 0 {
		t.Fatalf("Expected: %v, Got: %v (oldQueue)\n", 0, newQueue)
	}
}

func secondReduceAmendTest(t *testing.T, res *Reservation, reservations *Reservations, clients *Clients) {
	res.Size = 1
	var oldQueue = len(res.getTimeslot().Queue)
	var err = res.Amend(reservations, clients) // Amend /////  Here
	if err != nil {
		t.Fatalf("Amending reservation failed: %s\n", err)
	}
	var newReservations = len(res.getTimeslot().Reservations)
	var newQueue = len(res.getTimeslot().Queue)
	if newReservations != 1 {
		t.Fatalf("Expected: %v, Got: %v\n", 1, newReservations)
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
