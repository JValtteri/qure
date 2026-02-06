package state

import (
	"testing"

	"github.com/JValtteri/qure/server/internal/utils"
	"github.com/JValtteri/qure/server/internal/crypt"
	"github.com/JValtteri/qure/server/internal/state/model"
	"github.com/JValtteri/qure/server/internal/testjson"
)


func setTimeslot(size int) model.Timeslot {
	return model.Timeslot{
		Size: size,
	}
}

func TestValidateBadReservation(t *testing.T) {
	ResetEvents()
	ResetClients()
	timeslot := utils.Epoch(0)
	size := 1
	res, _ := newReservation(nil, nil, timeslot, size)
	err := res.Register(&reservations, &clients)
	t.Logf("%v\n", err)
	if err == nil {
		t.Errorf("Expected: %v, Got: %v\n", "error", err)
	}
}

func TestCreateReservationWithRegistered(t *testing.T) {
    ResetEvents()
    ResetClients()
    role := "test"
    email := "session@example.com"
    fingerprint := "0.0.0.0"
	time := utils.Epoch(1100)
    size := 1
    temp := false
    timeslot := setTimeslot(5)
    client, err := NewClient(role, email, crypt.Key("asdf"), temp)
    if err != nil {
        t.Fatalf("Error in creating client: %v", err)
    }
	sessionKey, _ := client.AddSession(role, email, temp, crypt.GenerateHash(fingerprint), &clients)
    event := EventFromJson(testjson.EventJson)
    eventID, err := CreateEvent(event)
	event.Append(timeslot, time)
    if err != nil {
        t.Errorf("Unexpected error in creating event: %v", err)
    }
    res := MakeReservation(sessionKey, email, fingerprint, crypt.GenerateHash(fingerprint), size, eventID, time, crypt.ID(""))
    if res.Error != "" {
        t.Errorf("Expected: %v, Got: %v\n", "", res.Error)
    }
    if res.Confirmed != size {
        t.Errorf("Expected: %v, Got: %v\n", size, res.Confirmed)
    }
	if reservations.ByEmail[email].Id != reservations.ByID[res.Id].Id {
		t.Fatalf("Reservations ByEmail and byID do not agree.\n")
	}
	clientReservation := reservationsFor(clients.ByEmail[email].Id)[0]
	if reservations.ByEmail[email].Id != clientReservation.Id {
		t.Fatalf("Reservations ByEmail and clientReservations do not agree.\n")
	}
}

func TestCreateReservationWithUnregistered(t *testing.T) {
    ResetEvents()
    ResetClients()
    email := "unregistered@example"
    fingerprint := "0.0.0.0"
	time := utils.Epoch(1100)
    size := 1
    timeslot := setTimeslot(1)
    event := EventFromJson(testjson.EventJson)
    eventID, err := CreateEvent(event)
    if err != nil {
        t.Fatalf("Unexpected error in creating event: %v", err)
    }
	event.Append(timeslot, time)
    res := MakeReservation("0", email, fingerprint, crypt.GenerateHash(fingerprint), size, eventID, 1100, crypt.ID(""))
    if res.Error != "" {
        t.Errorf("Expected: %v, Got: %v\n", nil, res.Error)
    }
    if res.Confirmed != size {
        t.Errorf("Expected: %v, Got: %v\n", size, res.Confirmed)
    }
}

func TestZeroReservation(t *testing.T) {
    ResetEvents()
    ResetClients()
    email := "zero@example"
    fingerprint := "0.0.0.0"
	time := utils.Epoch(1100)
    size := 0
    timeslot := setTimeslot(1)
    event := EventFromJson(testjson.EventJson)
    eventID, err := CreateEvent(event)
    if err != nil {
        t.Fatalf("Unexpected error in creating event: %v", err)
    }
	event.Append(timeslot, time)
    res := MakeReservation("0", email, fingerprint, crypt.GenerateHash(fingerprint), size, eventID, 1100, crypt.ID(""))
    if res.Error == "" {
        t.Errorf("Expected: %v, Got: %v\n", "error", res.Error)
    }
}

func TestTooSmallReservation(t *testing.T) {
    ResetEvents()
    ResetClients()
    role := "test"
    email := "session@example"
    fingerprint := "0.0.0.0"
	time := utils.Epoch(1100)
    size := 3
    temp := false
    slotSize := 4
    timeslot := setTimeslot(slotSize)
    client, err := NewClient(role, email, crypt.Key("asdf"), temp)
	sessionKey, _ := client.AddSession(role, email, temp, crypt.GenerateHash(fingerprint), &clients)
    event := EventFromJson(testjson.EventJson)
    eventID, err := CreateEvent(event)
	event.Append(timeslot, time)
    if err != nil {
        t.Errorf("Unexpected error in creating event: %v", err)
    }
    res := MakeReservation(sessionKey, email, fingerprint, crypt.GenerateHash(fingerprint), size, eventID, 1100, crypt.ID(""))
    if res.Error != "" {
        t.Errorf("Expected: %v, Got: %v\n", nil, res.Error)
    }
    if res.Confirmed != size {
        t.Errorf("Expected: %v, Got: %v\n", size, res.Confirmed)
    }
    res = MakeReservation(sessionKey, email, fingerprint, crypt.GenerateHash(fingerprint), size, eventID, 1100, crypt.ID(""))
    if res.Error != "" {
        t.Errorf("Expected: %v, Got: %v\n", nil, res.Error)
    }
    if res.Confirmed != (slotSize - size) {
        t.Errorf("Expected: %v, Got: %v\n", (slotSize - size), res.Confirmed)
    }
}

func TestInvalidReservation(t *testing.T) {
    ResetEvents()
    ResetClients()
    role := "test"
    email := "session@example.com"
    fingerprint := "0.0.0.0"
    size := 1
    temp := false
    eventID := crypt.ID("none")
	timeslot := utils.Epoch(1)
    client, err := NewClient(role, email, crypt.Key("asdf"), temp)
	key, _ := client.AddSession(role, email, temp, crypt.GenerateHash(fingerprint), &clients)
    if err != nil {
        t.Errorf("Unexpected error in creating event: %v", err)
    }
    res := MakeReservation(key, email, fingerprint, crypt.GenerateHash(fingerprint), size, eventID, timeslot, crypt.ID(""))
    if res.Error == "<nil>" {
        t.Errorf("Expected: %v, Got: %v\n", "error", res.Error)
    }
    if res.Confirmed != 0 {
        t.Errorf("Expected: %v, Got: %v\n", size, res.Confirmed)
    }
}

func TestFullSlotsReservation(t *testing.T) {
    ResetEvents()
    role := "test"
    email := "full@example"
    fingerprint := "0.0.0.0"
	time := utils.Epoch(1100)
    size := 3
    temp := false
    slotSize := 3
    timeslot := setTimeslot(slotSize)
    event := EventFromJson(testjson.EventJson)
    eventID, err := CreateEvent(event)
	event.Append(timeslot, time)
    client, err := NewClient(role, email, crypt.Key("asdf"), temp)
	key, _ := client.AddSession(role, email, temp, crypt.GenerateHash(fingerprint), &clients)
    if err != nil {
        t.Errorf("Unexpected error in creating event: %v", err)
    }
    _      = MakeReservation(key, email, fingerprint, crypt.GenerateHash(fingerprint), size, eventID, 1100, crypt.ID(""))
    res   := MakeReservation(key, email, fingerprint, crypt.GenerateHash(fingerprint), size, eventID, 1100, crypt.ID(""))
    if res.Error == "<nil>" {
        t.Errorf("Expected: %v, Got: %v\n", "error", res.Error)
    }
    if res.Confirmed != 0 {
        t.Errorf("Expected: %v, Got: %v\n", 0, res.Confirmed)
    }
}

func TestGetReservations(t *testing.T) {
    ResetEvents()
    email := "getreservationsemail@example"
    fingerprint := "0.0.0.1"
	time := utils.Epoch(1100)
    size := 1
    timeslot := setTimeslot(1)
    event := EventFromJson(testjson.EventJson)
    eventID, _ := CreateEvent(event)
	event.Append(timeslot, time)
    res := MakeReservation("0", email, fingerprint, crypt.GenerateHash(fingerprint), size, eventID, 1100, crypt.ID(""))
    expected := 1
	clientID := res.Client
	client := clients.ByID[res.Client]
    reservations := reservationsFor(clientID)
    if len(reservations) < expected {
        t.Fatalf("Expected: %v, Got: <%v\n", expected, expected)
    }
    if len(reservations) != expected {
        t.Errorf("Expected: %v, Got: %v\n", expected, len(reservations))
    }
	if reservations[0].Client != client.Id  {
		t.Errorf("Expected: %v, Got: %v\n", reservations[0].Client, client.Id)
	}
}

func TestNoReservationsForNobody(t *testing.T) {
    ResetEvents()
    if res := reservationsFor(crypt.ID("no-id")); res != nil {
        t.Errorf("Expected: %v, Got: %v\n", nil, len(res))
    }
}
