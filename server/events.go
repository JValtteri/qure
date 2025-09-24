package main

import (
    "log"
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

var events map[uint]Event = make(map[uint]Event)

func CreateEvent(eventJson []byte) bool {
    eventObj, err := eventFromJson(eventJson)
    if err != nil {
        return false
    }
    /* TODO check for duplicates
    if duplicate {
        log.Println("Duplicate Event")
        return false
    }
    */
    setId(&eventObj)
    events[eventObj.ID] = eventObj
    return true
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

func setId(event *Event) {
    currentTime := uint(time.Now().Unix())
    currentSeconds := currentTime % 60
    newID := uint(event.DtStart) + currentSeconds
    event.ID = newID
}

