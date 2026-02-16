package model

import (
	"fmt"
	"log"

	"github.com/JValtteri/qure/server/internal/crypt"
	"github.com/JValtteri/qure/server/internal/utils"
)

type Reservation struct {
	Id			crypt.ID
	Client		crypt.ID
	Size		int			// Party size
	Confirmed	int			// Reserved size
	Event		*Event
	Timeslot	utils.Epoch
	Expiration	utils.Epoch
	Session		crypt.Key	// Not stored, but sent as part of a response. Needed when a session is created simultaneously
	Error		string
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
	reserve, queue := calculateNewReservation(timeslot, r)
	r.updateTimeslotReservationsAndQueue(reserve, queue, &timeslot, reservations, clients)
	r.Confirmed = reserve
	err := reservations.update(*r, clients)
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
	var reserve, queue = r.calculateAmendedReservation(reservations)
	if reserve == 0 && queue == 0 {
		return fmt.Errorf("No change")
	}
	r.updateTimeslotReservationsAndQueue(reserve, queue, &timeslot, reservations, clients)
	r.Confirmed += reserve
	reservations.update(*r, clients)			// Sync reservations
	return err
}

func (r *Reservation) Cancel(reservations *Reservations, clients *Clients) error {
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
	if r.Size != 0 {
		return fmt.Errorf("No change")
	}
	// - Remove reservation
	//   - From event
	//   - From user
	//   - From reservations maps
	// - Advance queue
	r.removeFromReservationsAndQueue(&timeslot, reservations, clients)
	r.Confirmed = 0
	reservations.update(*r, clients)			// Sync reservations
	return err
}

func (r *Reservation) checkBasicValidity() error {
	if r.Event == nil || r.Client == "" {
		return fmt.Errorf("invalid reservation (event/client)")
	}
	return nil
}

func calculateNewReservation(timeslot Timeslot, r *Reservation) (int, int) {
	var freeSlots = timeslot.hasFree()
	var reserve = min(freeSlots, r.Size)
	var queue = r.Size - reserve
	return reserve, queue
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

func (r *Reservation) calculateAmendedReservation(reservations *Reservations) (int, int) {
	var reserve	= 0
	var queue	= 0
	var size	= r.Size

	var oldReservation, _	= r.getOldReservation(reservations)
	var timeslot 			= oldReservation.getTimeslot()
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

func (r *Reservation) updateTimeslotReservationsAndQueue(
	reserve int,			queue int,
	timeslot *Timeslot,		reservations *Reservations,		clients *Clients,
) {
	if reserve < 0 {
		timeslot.purgeFromQueue(r.Id)
		timeslot.removeNfromReservations(r.Id, -reserve)
		promoteFromQueue(timeslot, reservations, clients)
	} else if queue < 0 {
		timeslot.removeNfromQueue(r.Id, -queue)
	} else {
		timeslot.addToReservations(reserve, r.Id)
		timeslot.addToQueue(queue, r.Id)
	}
	timeslot.Reserved += reserve
	r.Event.Timeslots[r.Timeslot] = *timeslot
}

func (r *Reservation) removeFromReservationsAndQueue(
	timeslot *Timeslot,		reservations *Reservations,		clients *Clients,
) {
	timeslot.purgeFromQueue(r.Id)
	timeslot.purgeFromReservations(r.Id)
	promoteFromQueue(timeslot, reservations, clients)
	timeslot.Reserved = len(timeslot.Reservations)
	r.Event.Timeslots[r.Timeslot] = *timeslot
}

// Progress the Queue; Top off reserve slots with reservations from queue
func promoteFromQueue(timeslot *Timeslot, reservations *Reservations, clients *Clients) int {
	var promote = 0
	freeSlots := timeslot.hasFree()
	if freeSlots > 0 {
		queueSize := timeslot.QueueSize()
		promote = min(freeSlots, queueSize)
	}
	if promote > 0 {
		var fromQueue, err = timeslot.popFromQueue(promote)
		updatePromotedReservations(fromQueue, reservations, clients)
		if err != nil {
			log.Println("error: over indexed request from queue")
		}
		timeslot.appendToReservationsFromList(fromQueue)
	}
	return promote
}

func updatePromotedReservations(promoted []crypt.ID, reservations *Reservations, clients *Clients) {
	for _, id := range promoted {
		reservation, found := reservations.ByID[id]
		if !found {
			continue
		}
		reservation.Confirmed++
		reservations.update(reservation, clients)
		// TODO: Notify
	}
}
