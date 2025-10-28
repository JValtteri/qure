package state

import (
    "testing"
    "github.com/JValtteri/qure/server/internal/crypt"
)


func resetEvents() {
    events = make(map[crypt.ID]Event)
}

func setTimeslot(size int) Timeslot {
    return Timeslot{
        size: size,
    }
}

func TestValidateBadReservation(t *testing.T) {
    resetEvents()
    resetClients()
    timeslot := Epoch(0)
    size := 1
    res, _ := newReservation(nil, nil, timeslot, size)
    err := res.validate()
    t.Logf("%v\n", err)
    if err == nil {
        t.Errorf("Expected: %v, Got: %v\n", "error", err)
    }
}

func TestCreateReservationWithRegistered(t *testing.T) {
    resetEvents()
    resetClients()
    role := "test"
    email := "session@example.com"
    ip := IP("0.0.0.0")
    time := Epoch(1100)
    size := 1
    temp := false
    timeslot := setTimeslot(5)
    client, err := NewClient(role, email, crypt.Key("asdf"), temp)
    if err != nil {
        t.Fatalf("Error in creating client: %v", err)
    }
    sessionKey, _ := client.AddSession(role, email, temp, ip)
    eventID, err := CreateEvent(eventJson)
    event := events[eventID]
    event.append(timeslot, time)
    if err != nil {
        t.Errorf("Unexpected error in creating event: %v", err)
    }
    res := MakeReservation(sessionKey, email, ip, size, eventID, time)
    if res.Error != "<nil>" {
        t.Errorf("Expected: %v, Got: %v\n", "", res.Error)
    }
    if res.confirmed != size {
        t.Errorf("Expected: %v, Got: %v\n", size, res.confirmed)
    }
    if reservations.byEmail[email].id != reservations.byID[res.id].id {
        t.Fatalf("Reservations byEmail and byID do not agree.\n")
    }
    clientReservation := reservationsFor(clients.byEmail[email].id)[0]
    if reservations.byEmail[email].id != clientReservation.id {
        t.Fatalf("Reservations byEmail and clientReservations do not agree.\n")
    }
}

func TestCreateReservationWithUnregistered(t *testing.T) {
    resetEvents()
    resetClients()
    email := "unregistered@example"
    ip := IP("0.0.0.0")
    time := Epoch(1100)
    size := 1
    timeslot := setTimeslot(1)
    eventID, err := CreateEvent(eventJson)
    if err != nil {
        t.Fatalf("Unexpected error in creating event: %v", err)
    }
    event := events[eventID]
    event.append(timeslot, time)
    res := MakeReservation("0", email, ip, size, eventID, 1100)
    if res.Error != "<nil>" {
        t.Errorf("Expected: %v, Got: %v\n", nil, res.Error)
    }
    if res.confirmed != size {
        t.Errorf("Expected: %v, Got: %v\n", size, res.confirmed)
    }
}

func TestTooSmallReservation(t *testing.T) {
    resetEvents()
    resetClients()
    role := "test"
    email := "session@example"
    ip := IP("0.0.0.0")
    time := Epoch(1100)
    size := 3
    temp := false
    slotSize := 2
    timeslot := setTimeslot(slotSize)
    client, err := NewClient(role, email, crypt.Key("asdf"), temp)
    sessionKey, _ := client.AddSession(role, email, temp, ip)
    eventID, err := CreateEvent(eventJson)
    event := events[eventID]
    event.append(timeslot, time)
    if err != nil {
        t.Errorf("Unexpected error in creating event: %v", err)
    }
    res := MakeReservation(sessionKey, email, ip, size, eventID, 1100)
    if res.Error != "<nil>" {
        t.Errorf("Expected: %v, Got: %v\n", nil, res.Error)
    }
    if res.confirmed != slotSize {
        t.Errorf("Expected: %v, Got: %v\n", size, res.confirmed)
    }
}

func TestInvalidReservation(t *testing.T) {
    resetEvents()
    resetClients()
    role := "test"
    email := "session@example.com"
    ip := IP("0.0.0.0")
    size := 1
    temp := false
    eventID := crypt.ID("none")
    timeslot := Epoch(1)
    client, err := NewClient(role, email, crypt.Key("asdf"), temp)
    key, _ := client.AddSession(role, email, temp, ip)
    if err != nil {
        t.Errorf("Unexpected error in creating event: %v", err)
    }
    res := MakeReservation(key, email, ip, size, eventID, timeslot)
    if res.Error == "<nil>" {
        t.Errorf("Expected: %v, Got: %v\n", "error", res.Error)
    }
    if res.confirmed != 0 {
        t.Errorf("Expected: %v, Got: %v\n", size, res.confirmed)
    }
}

func TestFullSlotsReservation(t *testing.T) {
    resetEvents()
    role := "test"
    email := "full@example"
    ip := IP("0.0.0.0")
    time := Epoch(1100)
    size := 3
    temp := false
    slotSize := 3
    timeslot := setTimeslot(slotSize)
    eventID, err := CreateEvent(eventJson)
    event := events[eventID]
    event.append(timeslot, time)
    client, err := NewClient(role, email, crypt.Key("asdf"), temp)
    key, _ := client.AddSession(role, email, temp, ip)
    if err != nil {
        t.Errorf("Unexpected error in creating event: %v", err)
    }
    _      = MakeReservation(key, email, ip, size, eventID, 1100)
    res   := MakeReservation(key, email, ip, size, eventID, 1100)
    if res.Error == "<nil>" {
        t.Errorf("Expected: %v, Got: %v\n", "error", res.Error)
    }
    if res.confirmed != 0 {
        t.Errorf("Expected: %v, Got: %v\n", 0, res.confirmed)
    }
}

func TestGetReservations(t *testing.T) {
    resetEvents()
    email := "getreservationsemail@example"
    ip := IP("0.0.0.1")
    time := Epoch(1100)
    size := 1
    timeslot := setTimeslot(1)
    eventID, _ := CreateEvent(eventJson)
    event := events[eventID]
    event.append(timeslot, time)
    res := MakeReservation("0", email, ip, size, eventID, 1100)
    expected := 1
    clientID := res.client.id
    reservations := reservationsFor(clientID)
    if len(reservations) < expected {
        t.Fatalf("Expected: %v, Got: <%v\n", expected, expected)
    }
    if len(reservations) != expected {
        t.Errorf("Expected: %v, Got: %v\n", expected, len(reservations))
    }
    if reservations[0].client.email != email  {
        t.Errorf("Expected: %v, Got: %v\n", email, reservations[0].client.email)
    }
    if reservations[0].client.email != email  {
        t.Errorf("Expected: %v, Got: %v\n", email, reservations[0].client.email)
    }
}

