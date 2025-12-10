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
	Size			int
	Reserved		int
	Reservations	[]crypt.ID
	Queue			[]crypt.ID
}

func (e *Event)Append(timeslot Timeslot, time utils.Epoch) {
    e.Timeslots[time] = timeslot
}

func (t *Timeslot)isFull() bool {
    return len(t.Reservations) == t.Size
}

func (t *Timeslot)hasFree() int {
    return t.Size - len(t.Reservations)
}

func (t *Timeslot)append(res *Reservation) {
    partySize := res.Confirmed
	reSlice := slices.Repeat([]crypt.ID{res.Id}, partySize)
    t.Reservations = append(t.Reservations, reSlice...)
    t.Reserved = len(t.Reservations)
}
