package model

import (
	"testing"

	"github.com/JValtteri/qure/server/internal/crypt"
	"github.com/JValtteri/qure/server/internal/utils"
)

func TestEventModel(t *testing.T) {
	event, slot := getTestEvent()
	time := utils.Epoch(100)
	event.Append(slot, time)
	theSlot := event.Timeslots[time]
	if theSlot.isFull() {
		t.Errorf("Empty timeslot should not be full\n")
	}
	if free := theSlot.hasFree() ; free != 1 {
		t.Errorf("expected %v, got %v\n", 1, free)
	}
	res := getTestReservation()
	theSlot.Reservations = append(theSlot.Reservations, res.Id)

	theSlot.append(&res)
	if full := theSlot.isFull() ; !full {
		t.Errorf("Timeslot should be full. Is full: %v\n", full)
	}
	if free := theSlot.hasFree() ; free != 0 {
		t.Errorf("expected %v, got %v\n", 0, free)
	}
}

func getTestEvent() (Event, Timeslot) {
	event := Event{
		ID:        crypt.ID("event"),
		Timeslots: make(map[utils.Epoch]Timeslot),
	}
	slot := Timeslot{
		Size:     1,
		Reserved: 0,
		Reservations:	[]crypt.ID{},
		Queue: 			[]crypt.ID{},
	}
	return event, slot
}
