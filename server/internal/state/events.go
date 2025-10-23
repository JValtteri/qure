package state

import (
    "log"
    "fmt"
    "sync"
    "slices"
    "encoding/json"
    "github.com/JValtteri/qure/server/internal/utils"
    "github.com/JValtteri/qure/server/internal/crypt"
)

type Epoch = utils.Epoch

type Event struct {
    ID               crypt.ID
    Name             string;
    ShortDescription string;
    LongDescription  string;
    Draft            bool;
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
var events map[crypt.ID]Event = make(map[crypt.ID]Event)

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
