package state

import (
	"io"
	"log"
	"testing"

	"github.com/JValtteri/qure/server/internal/crypt"
	"github.com/JValtteri/qure/server/internal/testjson"
	"github.com/JValtteri/qure/server/internal/utils"
)

func TestCreateEvent(t *testing.T) {
    var input []byte = testjson.EventJson
    event := EventFromJson(input)
    expect := "ok"
    id, got := CreateEvent(event)
    if got != nil {
        t.Errorf("Expected: %v, Got: %v\n", expect, got)
    }
    if len(id) < 10 {
        t.Errorf("Expected: %v, Got: %v\n", ">=10", len(id))
    }
    ok := RemoveEvent(id)
    if !ok {
        t.Errorf("Expected: %v, Got: %v\n", 0, ok)
    }
}

func TestEditEvent(t *testing.T) {
	ResetEvents()
	var input []byte = testjson.EventJson
	event := EventFromJson(input)
	id, _ := CreateEvent(event)
	event = events[id]
	event.DtEnd = utils.Epoch(1337)
	_, err := EditEvent(event)
	if err != nil {
		t.Errorf("Got error editing event: %v\n", err)
	}
	if leet := events[id].DtEnd; leet != 1337 {
		t.Errorf("Expected: %v, Got: %v\n", 1337, leet)
	}
	event.ID = crypt.ID("wrong")
	if _, err = EditEvent(event); err == nil {
		t.Errorf("Expected error for invalid event ID: %v\n", err)
	}
}

func TestDuplicateEvent(t *testing.T) {
	ResetEvents()
    var input []byte = testjson.EventJson
    event := EventFromJson(input)
    expect := "ok"
    id, got := CreateEvent(event)
    if got != nil {
        t.Errorf("Expected: %v, Got: %v\n", expect, got)
    }
    _, got = CreateEvent(event)
    if got == nil {
        t.Errorf("Expected: %v, Got: %v\n", "error", got)
    }
    ok := RemoveEvent(id)
    if !ok {
        t.Errorf("Expected: %v, Got: %v\n", 0, ok)
    }
}

func TestRemoveNonexistent(t *testing.T) {
    input := crypt.ID("1010")
    const expect bool = false
    got := RemoveEvent(input)
    if got != expect {
        t.Errorf("Expected: %v, Got: %v\n", expect, got)
    }
}

func TestEventLifecycle(t *testing.T) {
    ResetEvents()
    var input []byte = testjson.EventJson
    event := EventFromJson(input)
    // Create Secondary target event
    _, err := CreateEvent(event)
    if err != nil {
        t.Errorf("Expected: %v, Got: %v\n", "ok", err)
    }
    // Create Primary target event
    eventObj := EventFromJson(input)
    expected := "Get event"
    id := crypt.ID("1337")
    eventObj.Name = expected
    eventObj.ID = id
    events[eventObj.ID] = eventObj
    if len(events) != 2 {
        t.Errorf("Expected: len=2, Got: len=%v\n", len(events))
    }
    log.SetOutput(io.Discard)
    ListEvents()
    got, err1 := GetEvent(id, true)
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
    got, err2 := GetEvent(id, true)
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
