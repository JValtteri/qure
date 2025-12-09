package model

import (
    "sync"
    "slices"
    "github.com/JValtteri/qure/server/internal/crypt"
    "github.com/JValtteri/qure/server/internal/utils"
)


var Eventslock sync.RWMutex = sync.RWMutex{}

type Event struct {
	ID					crypt.ID
	Name				string;
	ShortDescription	string;
	LongDescription		string;
	Draft				bool;
	DtStart				utils.Epoch;
	DtEnd				utils.Epoch;
	StaffSlots			int;
	Staff				int;
	Timeslots			map[utils.Epoch]Timeslot
}

type Timeslot struct {
    Size            int
    Reserved        int
    reservations    []*Reservation
    queue           []*Reservation
}

func (e *Event)Append(timeslot Timeslot, time utils.Epoch) {
    e.Timeslots[time] = timeslot
}

func (t *Timeslot)isFull() bool {
    return len(t.reservations) == t.Size
}

func (t *Timeslot)hasFree() int {
    return t.Size - len(t.reservations)
}

func (t *Timeslot)append(res *Reservation) {
    partySize := res.Confirmed
    reSlice := slices.Repeat([]*Reservation{res}, partySize)
    t.reservations = append(t.reservations, reSlice...)
    t.Reserved = len(t.reservations)
}
