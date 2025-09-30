package main

import (
    "io"
    "log"
    "testing"
)

func TestCreateEvent(t *testing.T) {
    var input []byte = eventJson
    expect := "ok"
    id, got := CreateEvent(input)
    if got != nil {
        t.Errorf("Expected: %v, Got: %v\n", expect, got)
    }
    ok := RemoveEvent(id)
    if !ok {
        t.Errorf("Expected: %v, Got: %v\n", 0, ok)
    }
}

func TestDuplicateEvent(t *testing.T) {
    var input []byte = eventJson
    expect := "ok"
    id, got := CreateEvent(input)
    if got != nil {
        t.Errorf("Expected: %v, Got: %v\n", expect, got)
    }
    _, got = CreateEvent(input)
    if got == nil {
        t.Errorf("Expected: %v, Got: %v\n", "error", got)
    }
    ok := RemoveEvent(id)
    if !ok {
        t.Errorf("Expected: %v, Got: %v\n", 0, ok)
    }
}

func TestCreateBadEvent(t *testing.T) {
    log.SetOutput(io.Discard)
    var input []byte = badJson
    const expect ID = "0"
    got, _ := CreateEvent(input)
    if got != expect {
        t.Errorf("Expected: %v, Got: %v\n", expect, got)
    }
}

func TestRemoveNonexistent(t *testing.T) {
    input := ID("1010")
    const expect bool = false
    got := RemoveEvent(input)
    if got != expect {
        t.Errorf("Expected: %v, Got: %v\n", expect, got)
    }
}

func TestEventLifecycle(t *testing.T) {
    var input []byte = eventJson
    // Create Secondary target event
    _, err := CreateEvent(input)
    if err != nil {
        t.Errorf("Expected: %v, Got: %v\n", "ok", err)
    }
    // Create Primary target event
    eventObj, _ := eventFromJson(eventJson)
    expected := "Get event"
    id := ID("1337")
    eventObj.Name = expected
    eventObj.ID = id
    events[eventObj.ID] = eventObj
    if len(events) != 2 {
        t.Errorf("Expected: len=2, Got: len=%v\n", len(events))
    }
    log.SetOutput(io.Discard)
    ListEvents()
    got, err1 := GetEvent(id)
    if err1 != nil {
        t.Errorf("Expected: found, Got: %v\n", err1)
    }
    if got.Name != expected {
        t.Errorf("Expected: %v, Got: %v\n", expected, got.Name)
    }

    ok := RemoveEvent(id)
    if !ok {
        t.Errorf("Expected: remove() %v, Got: %v\n", true, ok)
    }
    got, err2 := GetEvent(id)
    if err2 == nil {
        t.Errorf("Expected: find nonexistant %v, Got: %v\n", "nil", got.Name)
    }
    expected = "Test event"
    for id, event := range events {
        if event.Name != expected {
            t.Errorf("Expected: %v, Got: %v\n", expected, got.Name)
        }
        ok = RemoveEvent(id)
        if !ok {
            t.Errorf("Expected: remove() %v, Got: %v\n", true, ok)
        }
    }
    log.SetOutput(io.Discard)
    ListEvents()
}
