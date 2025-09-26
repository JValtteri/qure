package main

import (
    "log"
    "fmt"
    "time"
    "errors"
    "encoding/json"
)

type Epoch uint

type Event struct {
    ID               uint
    Name             string;
    ShortDescription string;
    LongDescription  string;
    DtStart          Epoch;
    DtEnd            Epoch;
    StaffSlots       int;
    Staff            int;
    GuestSlots       int;
    Guests           int;
}

// TODO: Change to thread safe datastructure
var events map[uint]Event = make(map[uint]Event)


func CreateEvent(eventJson []byte) (uint, error) {
    eventObj, err := eventFromJson(eventJson)
    if err != nil {
        return 0, err
    }
    id := setId(&eventObj)
    _, duplicate := events[id]
    if duplicate {
        return 0, errors.New(fmt.Sprintf("Error: Duplicate Event: %v", id))
    }
    events[eventObj.ID] = eventObj
    return id, nil
}

func GetEvent(id uint) (Event, error) {
    event, ok := events[id]
    if !ok {
        return event, errors.New("Error: Event not found")
    }
    return event, nil
}

func RemoveEvent(id uint) bool {
    _, ok := events[id]
    if !ok {
        return false
    }
    delete(events, id)
    return true
}

func ListEvents() {
    log.Println("Events: ", len(events))
    for key, obj := range events {
        log.Println("->   ", obj.Name, key)
    }
}

func eventFromJson(eventJson []byte) (Event, error) {
    var eventObj Event
    err := json.Unmarshal(eventJson, &eventObj)
    if err != nil {
        log.Println("Weather JSON Unmarshal error:", err)
        return eventObj, errors.New("JSON Unmarshal error")
    }
    return eventObj, nil
}

func setId(event *Event) uint {
    currentTime := uint(time.Now().Unix())
    currentSeconds := currentTime % 60
    newID := uint(event.DtStart) + currentSeconds
    event.ID = newID
    return newID
}

