package model

import (
	"fmt"
	"log"

	"github.com/JValtteri/qure/server/internal/crypt"
	"github.com/JValtteri/qure/server/internal/utils"
)

type Reservation struct {
	Id         crypt.ID
	Client     crypt.ID
	Size       int // Party size
	Confirmed  int // Reserved size
	Event      *Event
	Timeslot   utils.Epoch
	Expiration utils.Epoch
	Session    crypt.Key // Not stored, but sent as part of a response. Needed when a session is created simultaneously
	Error      string
}

// Propagets reservation OR get an error for not enough room
// implements partial reservations
func (r *Reservation) Register(reservations *Reservations, clients *Clients) error {
	if err := r.checkBasicValidity(); err != nil {
		return fmt.Errorf("invalid reservation (event/client)")
	}
	Eventslock.Lock()
	defer Eventslock.Unlock()

	timeslot := r.getTimeslot()
	if timeslot.isFull() {
		return fmt.Errorf("slot full")
	}
	err := r.propagateReservation(timeslot, reservations, clients) // Adds reservation to master data

	return err
}

func (r *Reservation) Amend(reservations *Reservations, clients *Clients) error {
	if err := r.checkBasicValidity(); err != nil {
		return fmt.Errorf("invalid reservation (event/client)")
	}

	// Target event is valid check
	oldReservation, err := r.getOldReservation(reservations)
	if err != nil {
		return fmt.Errorf("invalid reservation (event/client)")
	}

	Eventslock.Lock()
	defer Eventslock.Unlock()

	oldTimeslot := oldReservation.getTimeslot()
	var additionalSlots = r.Size - oldReservation.Confirmed - oldTimeslot.countInQueue(r.Id)

	if r.Size == oldReservation.Size {
		return fmt.Errorf("no change was made")
	} else if r.Size > oldReservation.Size {
		if !oldTimeslot.isFull() {
			freeSlots := oldTimeslot.hasFree()
			addToRes := min(freeSlots, additionalSlots)
			// TODO: Add to Reservations ////////////////////////////////////////////////////////////////////////////

			additionalQueueSlots := additionalSlots - addToRes
			oldTimeslot.addToQueue(additionalQueueSlots, r.Id) // add additionalSlots to queue

			// UpdateReservationDataEverywhere
		} else {

		}
	} else if r.Size < oldReservation.Confirmed {
		oldTimeslot.purgeFromReservations(r.Id)                        // Remove old reservation
		oldTimeslot.purgeFromQueue(r.Id)                          // Clear queue
		r.propagateReservation(oldTimeslot, reservations, clients) // Validate new reservation
		promoteFromQueue(&oldTimeslot)
	} else { // Remove items from queue
		oldTimeslot.removeNfromQueue(r.Id, -additionalSlots) // remove additionalSlots from queue
		promoteFromQueue(&oldTimeslot)
	}


	/*
	// Promote from queue if possible
	promote := promoteFromQueue(&oldTimeslot)
	r.Confirmed = oldReservation.Confirmed + promote
	*/

	// Update Reserved count
	oldTimeslot.Reserved = len(oldTimeslot.Reservations)

	// Save timeslot
	r.Event.Timeslots[r.Timeslot] = oldTimeslot

	// Save reservation
	reservations.append(*r, clients)

	return err
}

func (r *Reservation) checkBasicValidity() error {
	if r.Event == nil || r.Client == "" {
		return fmt.Errorf("invalid reservation (event/client)")
	}
	return nil
}

func (r *Reservation) getOldReservation(reservations *Reservations) (Reservation, error) {
	reservations.Lock()
	defer reservations.Unlock()

	target, ok := reservations.ByID[r.Id]
	if !ok {
		return target, fmt.Errorf("invalid reservation id")
	}
	return target, nil
}

func (r *Reservation) getTimeslot() Timeslot {
	return r.Event.Timeslots[r.Timeslot]
}

func (r *Reservation) propagateReservation(timeslot Timeslot, reservations *Reservations, clients *Clients) error {
	freeSlots := timeslot.hasFree()
	r.confirmSlots(freeSlots)
	r.updateEventTimeslot()
	err := reservations.append(*r, clients)
	return err
}

// Sets confirmed slots:
// According to what is available vs.
// What was requested
func (r *Reservation) confirmSlots(freeSlots int) {
	if freeSlots >= r.Size {
		r.Confirmed = r.Size
	} else {
		r.Confirmed = freeSlots
		// We could require the user confirm they want the partial reservation
		//r.Expiration = utils.EpochNow() + c.CONFIG.MAX_PENDIG_RESERVATION_TIME
	}
}

// Adds reservation to event and returns the updated timeslot
func (r *Reservation) updateEventTimeslot() {
	partySize := r.Confirmed
	timeslot := r.getTimeslot()						// Gets timeslot
	timeslot.addToReservations(partySize, r.Id)
	timeslot.Reserved = len(timeslot.Reservations)
	r.Event.Timeslots[r.Timeslot] = timeslot		// Updates timeslot
}

// Progress the Queue; Top off reserve slots with reservations from queue
func promoteFromQueue(oldTimeslot *Timeslot) int {
	var promote = 0
	freeSlots := oldTimeslot.hasFree()
	if freeSlots > 0 {
		queueSize := oldTimeslot.QueueSize()
		promote = min(freeSlots, queueSize)
	}
	if promote > 0 {
		var fromQueue, err = oldTimeslot.popFromQueue(promote)
		// TODO: UpdateReservations
		if err != nil {
			log.Println("error: over indexed request from queue")
		}
		oldTimeslot.appendToReservationsFromList(fromQueue)
	}
	return promote
}
