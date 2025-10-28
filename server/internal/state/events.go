package state

import (
    "log"
    "fmt"
    "encoding/json"
    "github.com/JValtteri/qure/server/internal/utils"
    "github.com/JValtteri/qure/server/internal/crypt"
)

type Epoch = utils.Epoch

func CreateEvent(eventJson []byte) (crypt.ID, error) {
    eventObj, err := eventFromJson(eventJson)
    if err != nil {
        return "0", err
    }
    id := setId(&eventObj)
    eventslock.Lock()
    defer eventslock.Unlock()
    _, duplicate := events[id]
    if duplicate {
        return "0", fmt.Errorf("duplicate Event: %v", id)
    }
    events[eventObj.ID] = eventObj
    return id, nil
}

func GetEvent(id crypt.ID) (Event, error) {
    eventslock.RLock()
    defer eventslock.RUnlock()
    event, ok := events[id]
    if !ok {
        return event, fmt.Errorf("event not found")
    }
    return event, nil
}

func GetEvents(isAdmin bool) []Event {
    var outEvents []Event
    eventslock.RLock()
    defer eventslock.RUnlock()
    for _, obj := range events {
        if !obj.Draft {
            outEvents = append(outEvents, obj)
        }
    }
    return outEvents
}

func RemoveEvent(id crypt.ID) bool {
    eventslock.Lock()
    defer eventslock.Unlock()
    _, ok := events[id]
    if !ok {
        return false
    }
    delete(events, id)
    return true
}

func ListEvents() {
    e := GetEvents(false)
    log.Println("Events: ", len(e))
    for _, obj := range e {
        log.Println("->   ", obj.Name, obj.ID)
    }
}

func eventFromJson(eventJson []byte) (Event, error) {
    var eventObj Event
    err := json.Unmarshal(eventJson, &eventObj)
    if err != nil {
        return eventObj, fmt.Errorf("JSON Unmarshal error: %v", err)
    }
    return eventObj, nil
}

func setId(event *Event) crypt.ID {
    currentTime := utils.EpochNow()
    currentSeconds := currentTime % 60
    uID := event.DtStart + currentSeconds
    newID := crypt.ID(fmt.Sprintf("%v", uID))
    event.ID = newID
    return newID
}


func MakeTestEvent(size int) crypt.ID {
    time := utils.Epoch(1100)
    timeslot := Timeslot{size: size}
    eventID, err := CreateEvent(eventJson)
    if err != nil {
        log.Fatalf("Unexpected error in creating event: %v", err)
    }
    event := events[eventID]
    event.append(timeslot, time)
    return eventID
}
