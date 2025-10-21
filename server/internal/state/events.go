package state

import (
    "log"
    "fmt"
    "sync"
    "slices"
    "encoding/json"
    "github.com/JValtteri/qure/server/internal/utils"
)

type Epoch = utils.Epoch

type Event struct {
    ID               ID
    Name             string;
    ShortDescription string;
    LongDescription  string;
    DtStart          Epoch;
    DtEnd            Epoch;
    StaffSlots       int;
    Staff            int;
    Timeslots        map[Epoch]Timeslot
}

type Timeslot struct {
    size            int
    reservations    []*Reservation
    queue           []*Reservation
}

func (e *Event)append(timeslot Timeslot, time Epoch) {
    e.Timeslots[time] = timeslot
}

func (t *Timeslot)isFull() bool {
    return len(t.reservations) == t.size
}

func (t *Timeslot)guests() int {
    return len(t.reservations)
}

func (t *Timeslot)hasFree() int {
    return t.size - len(t.reservations)
}

func (t *Timeslot)append(res *Reservation) {
    partySize := res.confirmed
    reSlice := slices.Repeat([]*Reservation{res}, partySize)
    t.reservations = append(t.reservations, reSlice...)
}

var eventslock sync.RWMutex = sync.RWMutex{}
var events map[ID]Event = make(map[ID]Event)

func CreateEvent(eventJson []byte) (ID, error) {
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

func GetEvent(id ID) (Event, error) {
    eventslock.RLock()
    defer eventslock.RUnlock()
    event, ok := events[id]
    if !ok {
        return event, fmt.Errorf("event not found")
    }
    return event, nil
}

func RemoveEvent(id ID) bool {
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
    eventslock.RLock()
    defer eventslock.RUnlock()
    log.Println("Events: ", len(events))
    for key, obj := range events {
        log.Println("->   ", obj.Name, key)
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

func setId(event *Event) ID {
    currentTime := utils.EpochNow()
    currentSeconds := currentTime % 60
    uID := event.DtStart + currentSeconds
    newID := ID(fmt.Sprintf("%v", uID))
    event.ID = newID
    return newID
}
