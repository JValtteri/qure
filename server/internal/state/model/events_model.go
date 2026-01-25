package model

import (
	"fmt"
	"slices"
	"sync"

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
	t.Reservations = addNElementsToList(partySize, res.Id, t.Reservations)
    t.Reserved = len(t.Reservations)
}

func (t *Timeslot)appendReservations(newReservations []crypt.ID) {
	t.Reservations = append(t.Reservations, newReservations...)
}

// Removes all instances of targetID form Reservations
func (t *Timeslot)purgeReservations(targetID crypt.ID) int {
	var count int = 0
	for index, value := range t.Reservations {
		if value == targetID {
			t.Reservations = slices.Delete(t.Reservations, index, index+1)
			count++
		}
	}
	return count
}

// Removes all instances of targetID form Queue
func (t *Timeslot)removeFromQueue(targetID crypt.ID) {
	var filtered = []crypt.ID{}
	for _, value := range t.Queue {
		if value != targetID {
			filtered = append(filtered, value)
		}
	}
	t.Reservations = filtered
}

// Removes N instances of targetID form Queue.
// Elements are removed in reverse order, from back of the queue
func (t *Timeslot)removeNfromQueue(targetID crypt.ID, count int) {
	var filtered = []crypt.ID{}
	for i := len(t.Queue) ; i >= 0 && count > 0 ; i-- {
		if t.Queue[i] != targetID {
			filtered = append(filtered, t.Queue[i])
		} else {
			count--
		}
	}
	slices.Reverse(filtered)
	t.Queue = filtered
}

func (t *Timeslot)addToQueue(number int, targetID crypt.ID) {
	addNElementsToList(number, targetID, t.Queue)
}

// Pops the first N elements of the Queue
func (t *Timeslot)popFromQueue(number int) ([]crypt.ID, error) {
	var err error = nil
	var queueSize = len(t.Queue)
	if number > queueSize {
		number = queueSize
		err = fmt.Errorf("over index")
	}
	var popped = t.Queue[0:number]
	t.Queue = t.Queue[number:]
	return popped, err
}

func addNElementsToList(number int, item crypt.ID, list []crypt.ID) []crypt.ID {
	items := slices.Repeat([]crypt.ID{item}, number)
	return append(list, items...)
}

func (t *Timeslot)countInQueue(targetID crypt.ID) int {
	var count int = 0
	for _, value := range t.Reservations {
		if value == targetID {
			count++
		}
	}
	return count
}
