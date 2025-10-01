package main

import (
    "testing"
)

func resetEvents() {
    events = make(map[ID]Event)
}

func setTimeslot(size int) Timeslot {
    return Timeslot{
        size: size,
    }
}

func TestCreateReservationWithRegistered(t *testing.T) {
    resetEvents()
    role := "test"
    email := "session@example.com"
    ip := IP("0.0.0.0")
    time := Epoch(1100)
    size := 1
    temp := false
    timeslot := setTimeslot(5)
    sessionKey, _ := AddSession(role, email, temp, ip)
    eventID, err := CreateEvent(eventJson)
    event := events[eventID]
    event.append(timeslot, time)
    if err != nil {
        t.Errorf("Unexpected error in creating event: %v", err)
    }
    res, err := MakeReservation(sessionKey, email, ip, size, eventID, time)
    if err != nil {
        t.Errorf("Expected: %v, Got: %v\n", nil, err)
    }
    if res.confirmed != size {
        t.Errorf("Expected: %v, Got: %v\n", size, res.confirmed)
    }
}

func TestCreateReservationWithUnregistered(t *testing.T) {
    resetEvents()
    email := "unregistered@example"
    ip := IP("0.0.0.0")
    time := Epoch(1100)
    size := 1
    timeslot := setTimeslot(1)
    eventID, err := CreateEvent(eventJson)
    event := events[eventID]
    event.append(timeslot, time)
    if err != nil {
        t.Errorf("Unexpected error in creating event: %v", err)
    }
    res, err := MakeReservation("0", email, ip, size, eventID, 1100)
    if err != nil {
        t.Errorf("Expected: %v, Got: %v\n", nil, err)
    }
    if res.confirmed != size {
        t.Errorf("Expected: %v, Got: %v\n", size, res.confirmed)
    }
}

func TestTooSmallReservation(t *testing.T) {
    resetEvents()
    role := "test"
    email := "session@example"
    ip := IP("0.0.0.0")
    time := Epoch(1100)
    size := 3
    temp := false
    slotSize := 2
    timeslot := setTimeslot(slotSize)
    sessionKey, _ := AddSession(role, email, temp, ip)
    eventID, err := CreateEvent(eventJson)
    event := events[eventID]
    event.append(timeslot, time)
    if err != nil {
        t.Errorf("Unexpected error in creating event: %v", err)
    }
    res, err := MakeReservation(sessionKey, email, ip, size, eventID, 1100)
    if err != nil {
        t.Errorf("Expected: %v, Got: %v\n", nil, err)
    }
    if res.confirmed != slotSize {
        t.Errorf("Expected: %v, Got: %v\n", size, res.confirmed)
    }
}

func TestInvalidReservation(t *testing.T) {
    resetEvents()
    role := "test"
    email := "session@example.com"
    ip := IP("0.0.0.0")
    size := 1
    temp := false
    eventID := ID("none")
    timeslot := Epoch(1)
    key, err := AddSession(role, email, temp, ip)
    if err != nil {
        t.Errorf("Unexpected error in creating event: %v", err)
    }
    res, err := MakeReservation(key, email, ip, size, eventID, timeslot)
    if err == nil {
        t.Errorf("Expected: %v, Got: %v\n", "error", err)
    }
    if res.confirmed != 0 {
        t.Errorf("Expected: %v, Got: %v\n", size, res.confirmed)
    }
}

func TestFullSlotsReservation(t *testing.T) {
    resetEvents()
    email := "session@example"
    ip := IP("0.0.0.0")
    time := Epoch(1100)
    size := 3
    slotSize := 3
    timeslot := setTimeslot(slotSize)
    eventID, err := CreateEvent(eventJson)
    event := events[eventID]
    event.append(timeslot, time)
    if err != nil {
        t.Errorf("Unexpected error in creating event: %v", err)
    }
    _, _      = MakeReservation("0", email, ip, size, eventID, 1100)
    res, err := MakeReservation("0", email, ip, size, eventID, 1100)
    if err == nil {
        t.Errorf("Expected: %v, Got: %v\n", "error", err)
    }
    if res.confirmed != 0 {
        t.Errorf("Expected: %v, Got: %v\n", 0, res.confirmed)
    }
}
