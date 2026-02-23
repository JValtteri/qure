package model

import (
	"fmt"
	"slices"

	"github.com/JValtteri/qure/server/internal/crypt"
	"github.com/JValtteri/qure/server/internal/utils"
)

type Timeslot struct {
	Size         int
	Reserved     int
	Reservations []crypt.ID
	Queue        []crypt.ID
}

func (e *Event) Append(timeslot Timeslot, time utils.Epoch) {
	e.Timeslots[time] = timeslot
}

func (t *Timeslot) isFull() bool {
	return len(t.Reservations) == t.Size
}

func (t *Timeslot) hasFree() int {
	return t.Size - len(t.Reservations)
}

func (t *Timeslot) addToReservations(number int, targetID crypt.ID) {
	t.Reservations = addNElementsToList(number, targetID, t.Reservations)
}

func (t *Timeslot) appendToReservationsFromList(newReservations []crypt.ID) {
	t.Reservations = append(t.Reservations, newReservations...)
}

// Removes all instances of targetID form Reservations
func (t *Timeslot) purgeFromReservations(targetID crypt.ID) {
	var filtered []crypt.ID = FilterFrom(t.Reservations, targetID)
	t.Reservations = filtered
}

// Removes all instances of targetID form Queue
func (t *Timeslot) purgeFromQueue(targetID crypt.ID) {
	var filtered []crypt.ID = FilterFrom(t.Queue, targetID)
	t.Queue = filtered
}

// Removes N instances of targetID form Queue.
func (t *Timeslot) removeNfromReservations(targetID crypt.ID, count int) {
	var filtered []crypt.ID = filterNfrom(count, t.Reservations, targetID)
	t.Reservations = filtered
}

// Removes N instances of targetID form Queue.
//
// Elements are removed in REVERSE ORDER, from back of the queue
func (t *Timeslot) removeNfromQueue(targetID crypt.ID, count int) {
	var filtered = []crypt.ID{}
	for i := len(t.Queue) - 1; i >= 0 && count > 0; i-- {
		if t.Queue[i] != targetID {
			filtered = append(filtered, t.Queue[i])
		} else {
			count--
		}
	}
	slices.Reverse(filtered)
	t.Queue = filtered
}

func (t *Timeslot) addToQueue(number int, targetID crypt.ID) {
	t.Queue = addNElementsToList(number, targetID, t.Queue)
}

// Pops the first N elements of the Queue
func (t *Timeslot) popFromQueue(number int) ([]crypt.ID, error) {
	var popped, err = pop(&t.Queue, number)
	return popped, err
}

func (t *Timeslot) countInQueue(targetID crypt.ID) int {
	var count = t.countInList(targetID, &t.Queue)
	return count
}

func (t *Timeslot) countInList(targetID crypt.ID, list *[]crypt.ID) int {
	var count int = 0
	for _, value := range *list {
		if value == targetID {
			count++
		}
	}
	return count
}

func (t *Timeslot) QueueSize() int {
	return len(t.Queue)
}


// Filters from list all instances of targetID
func FilterFrom(list []crypt.ID, targetID crypt.ID) []crypt.ID {
	var filtered = []crypt.ID{}
	for _, value := range list {
		if value != targetID {
			filtered = append(filtered, value)
		}
	}
	return filtered
}

// Filters count of targetID from list
func filterNfrom(count int, list []crypt.ID, targetID crypt.ID) []crypt.ID {
	var filtered = []crypt.ID{}
	for _, value := range list {
		if count == 0 {
			filtered = append(filtered, value)
		} else if value != targetID {
			filtered = append(filtered, value)
		} else {
			count--
		}
	}
	return filtered
}

func pop(list *[]crypt.ID, number int) ([]crypt.ID, error) {
	var err error = nil
	var queueSize = len(*list)
	if number > queueSize {
		number = queueSize
		err = fmt.Errorf("over index")
	}
	var popped = (*list)[0:number]
	*list = (*list)[number:]
	return popped, err
}
