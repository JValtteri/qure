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
	_, err := r.confirmReservation(timeslot, reservations, clients) // Adds reservation to master data

	return err
}

func (r *Reservation) Amend(reservations *Reservations, clients *Clients) error {
	if err := r.checkBasicValidity(); err != nil {
		return fmt.Errorf("invalid reservation (event/client)")
	}
	oldReservation, err := r.getOldReservation(reservations)
	if err != nil {
		return fmt.Errorf("invalid reservation (event/client)")
	}

	Eventslock.Lock()
	defer Eventslock.Unlock()

	var timeslot = oldReservation.getTimeslot()
	var reserve, queue = r.calculateReservation(reservations)
	if reserve == 0 && queue == 0 {
		return fmt.Errorf("No change")
	}
	r.updateTimeslot(reserve, queue, &timeslot)
	r.Confirmed += reserve
	r.Event.Timeslots[r.Timeslot] = timeslot	// Update to Event
	reservations.update(*r, clients)			// Sync reservations
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

func (r *Reservation) calculateReservation(reservations *Reservations) (int, int) {
	var reserve	= 0
	var queue	= 0
	var size	= r.Size

	var oldReservation, _	= r.getOldReservation(reservations)
	var timeslot 		  	= oldReservation.getTimeslot()
	var freeSlots			= timeslot.hasFree()
	var inQueue				= timeslot.countInQueue(r.Id)
	var additionalSlots 	= r.Size - oldReservation.Confirmed - inQueue

	if size == oldReservation.Size {
		return reserve, queue
	}
	if size > oldReservation.Size {
		reserve = min(freeSlots, additionalSlots)
		queue = additionalSlots - reserve
	} else if additionalSlots < 0 {
		if additionalSlots + inQueue >= 0 {
			queue = additionalSlots
		} else {
			queue = -inQueue
			reserve = additionalSlots + inQueue
		}
	}
	return reserve, queue
}

func (r *Reservation) updateTimeslot(reserve int, queue int, timeslot *Timeslot) {
	if reserve < 0 {
		timeslot.purgeFromQueue(r.Id)
		timeslot.removeNfromReservations(r.Id, -reserve)
		promoteFromQueue(timeslot)
	} else if queue < 0 {
		timeslot.removeNfromQueue(r.Id, -queue)
	} else {
		timeslot.addToReservations(reserve, r.Id)
		timeslot.addToQueue(queue, r.Id)
	}
	timeslot.Reserved += reserve
}



func (r *Reservation) confirmReservation(timeslot Timeslot, reservations *Reservations, clients *Clients) (int, error) {
	freeSlots := timeslot.hasFree()
	confirmed := r.confirmSlots(freeSlots)
	r.updateEventTimeslot()
	err := reservations.update(*r, clients)
	return confirmed, err
}

// Sets confirmed slots:
// According to what is available vs.
// What was requested
func (r *Reservation) confirmSlots(freeSlots int) int {
	confirmed := min(freeSlots, r.Size)
	r.Confirmed = confirmed
	// We could require the user confirm they want the partial reservation
	//r.Expiration = utils.EpochNow() + c.CONFIG.MAX_PENDIG_RESERVATION_TIME
	return confirmed
}

// Adds reservation to Event's timeslot
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
