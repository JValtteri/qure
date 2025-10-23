package state

import (
    "sync"
    "slices"
    "github.com/JValtteri/qure/server/internal/crypt"
)


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
