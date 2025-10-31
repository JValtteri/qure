package state

import (
    "testing"
    "github.com/JValtteri/qure/server/internal/crypt"
)


func setTimeslot(size int) Timeslot {
    return Timeslot{
        Size: size,
    }
}

func TestValidateBadReservation(t *testing.T) {
    ResetEvents()
    ResetClients()
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
    ResetEvents()
    ResetClients()
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
    event := EventFromJson(EventJson)
    eventID, err := CreateEvent(event)
    event.append(timeslot, time)
    if err != nil {
        t.Errorf("Unexpected error in creating event: %v", err)
    }
    res := MakeReservation(sessionKey, email, ip, size, eventID, time)
    if res.Error != "" {
        t.Errorf("Expected: %v, Got: %v\n", "", res.Error)
    }
    if res.Confirmed != size {
        t.Errorf("Expected: %v, Got: %v\n", size, res.Confirmed)
    }
    if reservations.byEmail[email].Id != reservations.byID[res.Id].Id {
        t.Fatalf("Reservations byEmail and byID do not agree.\n")
    }
    clientReservation := reservationsFor(clients.byEmail[email].Id)[0]
    if reservations.byEmail[email].Id != clientReservation.Id {
        t.Fatalf("Reservations byEmail and clientReservations do not agree.\n")
    }
}

func TestCreateReservationWithUnregistered(t *testing.T) {
    ResetEvents()
    ResetClients()
    email := "unregistered@example"
    ip := IP("0.0.0.0")
    time := Epoch(1100)
    size := 1
    timeslot := setTimeslot(1)
    event := EventFromJson(EventJson)
    eventID, err := CreateEvent(event)
    if err != nil {
        t.Fatalf("Unexpected error in creating event: %v", err)
    }
    event.append(timeslot, time)
    res := MakeReservation("0", email, ip, size, eventID, 1100)
    if res.Error != "" {
        t.Errorf("Expected: %v, Got: %v\n", nil, res.Error)
    }
    if res.Confirmed != size {
        t.Errorf("Expected: %v, Got: %v\n", size, res.Confirmed)
    }
}

func TestTooSmallReservation(t *testing.T) {
    ResetEvents()
    ResetClients()
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
    event := EventFromJson(EventJson)
    eventID, err := CreateEvent(event)
    event.append(timeslot, time)
    if err != nil {
        t.Errorf("Unexpected error in creating event: %v", err)
    }
    res := MakeReservation(sessionKey, email, ip, size, eventID, 1100)
    if res.Error != "" {
        t.Errorf("Expected: %v, Got: %v\n", nil, res.Error)
    }
    if res.Confirmed != slotSize {
        t.Errorf("Expected: %v, Got: %v\n", size, res.Confirmed)
    }
}

func TestInvalidReservation(t *testing.T) {
    ResetEvents()
    ResetClients()
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
    if res.Confirmed != 0 {
        t.Errorf("Expected: %v, Got: %v\n", size, res.Confirmed)
    }
}

func TestFullSlotsReservation(t *testing.T) {
    ResetEvents()
    role := "test"
    email := "full@example"
    ip := IP("0.0.0.0")
    time := Epoch(1100)
    size := 3
    temp := false
    slotSize := 3
    timeslot := setTimeslot(slotSize)
    event := EventFromJson(EventJson)
    eventID, err := CreateEvent(event)
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
    if res.Confirmed != 0 {
        t.Errorf("Expected: %v, Got: %v\n", 0, res.Confirmed)
    }
}

func TestGetReservations(t *testing.T) {
    ResetEvents()
    email := "getreservationsemail@example"
    ip := IP("0.0.0.1")
    time := Epoch(1100)
    size := 1
    timeslot := setTimeslot(1)
    event := EventFromJson(EventJson)
    eventID, _ := CreateEvent(event)
    event.append(timeslot, time)
    res := MakeReservation("0", email, ip, size, eventID, 1100)
    expected := 1
    clientID := res.Client.Id
    reservations := reservationsFor(clientID)
    if len(reservations) < expected {
        t.Fatalf("Expected: %v, Got: <%v\n", expected, expected)
    }
    if len(reservations) != expected {
        t.Errorf("Expected: %v, Got: %v\n", expected, len(reservations))
    }
    if reservations[0].Client.email != email  {
        t.Errorf("Expected: %v, Got: %v\n", email, reservations[0].Client.email)
    }
    if reservations[0].Client.email != email  {
        t.Errorf("Expected: %v, Got: %v\n", email, reservations[0].Client.email)
    }
}

func TestNoReservationsForNobody(t *testing.T) {
    ResetEvents()
    if res := reservationsFor(crypt.ID("no-id")); res != nil {
        t.Errorf("Expected: %v, Got: %v\n", nil, len(res))
    }
}
