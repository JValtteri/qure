package model

import (
	"sync"
	"fmt"

	"github.com/JValtteri/qure/server/internal/utils"
	"github.com/JValtteri/qure/server/internal/crypt"
	c "github.com/JValtteri/qure/server/internal/config"
)


type Reservations struct {
	mu			sync.RWMutex
	ByID		map[crypt.ID]Reservation
	ByEmail		map[string]*Reservation
}

func (r *Reservations) Lock() {
    r.mu.Lock()
}

func (r *Reservations) Unlock() {
    r.mu.Unlock()
}

func (r *Reservations) RLock() {
    r.mu.RLock()
}

func (r *Reservations) RUnlock() {
    r.mu.RUnlock()
}

func (r *Reservations) append(res Reservation, clients *Clients) error {
    r.Lock()
    defer r.Unlock()
	r.ByID[res.Id] = res
    clientEmail := clients.ByID[res.Client].Email
	r.ByEmail[clientEmail] = &res
    return clients.AddReservation(res.Client, &res)
}

type Reservation struct {
	Id			crypt.ID
	Client		crypt.ID
	Size		int			// Party size
	Confirmed	int			// Reserved size
	Event		*Event
	Timeslot	utils.Epoch
	Expiration	utils.Epoch
	Session		crypt.Key   // Not stored, but sent as part of a response. Neede when a session is created simultaneously
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
    if timeslot.isFull() {
        return fmt.Errorf("slot full")
    }
    err := r.propagateReservation(timeslot, reservations, clients)			// Adds reservation to master data

    return err
}

func (r *Reservation) Amend(reservations *Reservations, clients *Clients) error {
	if err := r.checkBasicValidity(); err != nil {
        return fmt.Errorf("invalid reservation (event/client)")
    }

	// Target event is valid check
	oldReservation, err := r.checkModificationValidity(reservations)
	if err != nil {
        return fmt.Errorf("invalid reservation (event/client)")
    }

	Eventslock.Lock()
	defer Eventslock.Unlock()

	timeslot := oldReservation.getTimeslot()
	var additionalSlots = r.Size - oldReservation.Confirmed - timeslot.countInQueue(r.Id)

	fmt.Printf("New slots: %v = %v - %v - %v\n", additionalSlots, r.Size, oldReservation.Confirmed, timeslot.countInQueue(r.Id))  // DEBUG

	if r.Size == oldReservation.Size {
		return fmt.Errorf("no change was made")
	} else if r.Size < oldReservation.Confirmed {
		timeslot.purgeReservations(r.Id)						// Remove old reservation
		timeslot.removeFromQueue(r.Id)							// Clear queue
		r.propagateReservation(timeslot, reservations, clients) // Validate new reservation
	} else if additionalSlots > 0 {
		timeslot.addToQueue(additionalSlots, r.Id)				// add additionalSlots to queue
	} else {
		timeslot.removeNfromQueue(r.Id, -additionalSlots)		// remove additionalSlots from queue
	}

	// TODO: Propagate Updated Reservation

	return err
}

func (r *Reservation) checkBasicValidity() error {
    if r.Event == nil || r.Client == "" {
        return fmt.Errorf("invalid reservation (event/client)")
    }
    return nil
}

func (r *Reservation) checkModificationValidity(reservations *Reservations) (Reservation, error) {
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
		r.Expiration = utils.EpochNow() + c.CONFIG.MAX_PENDIG_RESERVATION_TIME
    }
}

// Adds reservation to event and returns the updated timeslot
func (r *Reservation) updateEventTimeslot() {
	timeslot := r.getTimeslot()					// Gets timeslot
    timeslot.append(r)                          // Adds reservation to event
    r.Event.Timeslots[r.Timeslot] = timeslot    // Returns updated timeslot
}
