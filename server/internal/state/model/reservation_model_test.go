package model

import (
	"testing"

	"github.com/JValtteri/qure/server/internal/crypt"
	"github.com/JValtteri/qure/server/internal/utils"
)

func TestCreateAndAmmendReservations(t *testing.T) {
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
	// 1 = 1 confirmed + 0 in queue
	res := getTestReservation()
	res.Client = client.Id
	res.Event = &event
	err := res.Register(&reservations, &clients)
	if err != nil {
		t.Fatalf("Validating reservation failed: %s\n", err)
	}

	// Test Amending reservation
	// 2 = 2 confirmed + 0 in queue
	firstAmendTest(t, &res, &reservations, &clients)

	// Test Amend again
	// Should flow to queue
	// 3 = 2 confirmed + 1 in queue
	secondAmendTest(t, &res, &reservations, &clients)

	// Test equavilant reservation
	amendNoChangeTest(t, &res, &reservations, &clients)

	// Test smaller reservation
	// Should remove one slot from queue
	// 2 = 2 confirmed + 0 in queue
	firstReduceAmendTest(t, &res, &reservations, &clients)

	// Test smaller reservation
	// Should remove one slot from reservations
	// 1 = 1 confirmed + 0 in queue
	secondReduceAmendTest(t, &res, &reservations, &clients)
}

func TestInitialReservationQueueFunction(t *testing.T) {
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
	// 3 = 2 confirmed + 1 in queue
	res := getTestReservation()
	res.Size = 3
	res.Client = client.Id
	res.Event = &event
	err := res.Register(&reservations, &clients)
	if err != nil {
		t.Fatalf("Validating reservation failed: %s\n", err)
	}
	var newReservations = len(res.getTimeslot().Reservations)
	var newQueue = len(res.getTimeslot().Queue)
	var expectedReservations = 2
	var expectedQueue = 1
	if newReservations != expectedReservations {
		t.Fatalf("Expected: %v, Got: %v\n", expectedReservations, newReservations)
	}
	if newQueue != expectedQueue {
		t.Fatalf("Expected: %v, Got: %v (oldQueue)\n", expectedQueue, newQueue)
	}
}

func TestFullSlot(t *testing.T) {
	// Setup init objects
	reservations := makeTestReservations()
	clients := getTestClients()
	client1 := getTestClient()
	clients.ByEmail[client1.Email] = &client1
	clients.ByID[client1.Id] = &client1

	client2 := getTestClient()
	client2.Id = crypt.ID("22")
	clients.ByEmail[client2.Email] = &client2
	clients.ByID[client2.Id] = &client2

	// Add Event
	event, slot := getTestEvent()
	slot.Size = 1
	time := utils.Epoch(200)
	event.Append(slot, time)

	// Test Reservation
	// 3 = 2 confirmed + 1 in queue
	res1 := getTestReservation()
	res1.Size = 1
	res1.Client = client1.Id
	res1.Event = &event
	err := res1.Register(&reservations, &clients)
	if err != nil {
		t.Fatalf("Validating reservation failed: %s\n", err)
	}
	res2 := getTestReservation()
	res2.Size = 1
	res2.Client = client2.Id
	res2.Event = &event
	err = res2.Register(&reservations, &clients)
	if err == nil {
		t.Fatalf("Expected slot full: %s\n", err)
	}
	var newReservations = len(res2.getTimeslot().Reservations)
	var newQueue = len(res2.getTimeslot().Queue)
	var expectedQueue = 0
	if newReservations != 1 {
		t.Fatalf("Expected: %v, Got: %v\n", 1, newReservations)
	}
	if newQueue != expectedQueue {
		t.Fatalf("Expected: %v, Got: %v (oldQueue)\n", expectedQueue, newQueue)
	}
}

func TestPromoteFromQueue(t *testing.T) {
	// TODO:
	t.Log("TestPromoteFromQueue not implemented!")
}

/* Sub-Tests */

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
	var expectedQueue = 0
	if oldQueue != expectedQueue {
		t.Fatalf("Expected: %v, Got: %v (oldQueue)\n", expectedQueue, oldQueue)
	}
	if newQueue != expectedQueue {
		t.Fatalf("Expected: %v, Got: %v (oldQueue)\n", expectedQueue, newQueue)
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
	var expectedReservations = 2
	var expectedOldQueue = 0
	var expectedNewQueue = 1
	if newReservations != expectedReservations {
		t.Fatalf("Expected: %v, Got: %v\n", expectedReservations, newReservations)
	}
	if oldQueue != expectedOldQueue {
		t.Fatalf("Expected: %v, Got: %v (oldQueue)\n", expectedOldQueue, oldQueue)
	}
	if newQueue != expectedNewQueue {
		t.Fatalf("Expected: %v, Got: %v (oldQueue)\n", expectedNewQueue, newQueue)
	}
}

func amendNoChangeTest(t *testing.T, res *Reservation, reservations *Reservations, clients *Clients) {
	var err = res.Amend(reservations, clients)
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
	var expectedReservations = 2
	var expectedOldQueue = 1
	var expectedNewQueue = 0
	if newReservations != expectedReservations {
		t.Fatalf("Expected: %v, Got: %v\n", expectedReservations, newReservations)
	}
	if oldQueue != expectedOldQueue {
		t.Fatalf("Expected: %v, Got: %v (oldQueue)\n", expectedOldQueue, oldQueue)
	}
	if newQueue != expectedNewQueue {
		t.Fatalf("Expected: %v, Got: %v (oldQueue)\n", expectedNewQueue, newQueue)
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
	var expectedReservations = 1
	var expectedOldQueue = 0
	var expectedNewQueue = 0
	if newReservations != expectedReservations {
		t.Fatalf("Expected: %v, Got: %v\n", expectedReservations, newReservations)
	}
	if oldQueue != expectedOldQueue {
		t.Fatalf("Expected: %v, Got: %v (oldQueue)\n", expectedOldQueue, oldQueue)
	}
	if newQueue != expectedNewQueue {
		t.Fatalf("Expected: %v, Got: %v (oldQueue)\n", expectedNewQueue, newQueue)
	}
}

/* Helper Functions */

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
