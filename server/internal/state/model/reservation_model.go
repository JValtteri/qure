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
func (r *Reservation) Validate(reservations *Reservations, clients *Clients) error {
    if err := r.checkBasicValidity(); err != nil {
        return fmt.Errorf("invalid reservation (event/client)")
    }
	Eventslock.Lock()
	defer Eventslock.Unlock()

    slot := r.getTimeslot()
    if slot.isFull() {
        return fmt.Errorf("slot full")
    }
    freeSlots := slot.hasFree()
    r.confirmSlots(freeSlots)
    r.updateEventTimeslot()
	err := reservations.append(*r, clients)			// Adds reservation to master data

    return err
}

func (r *Reservation) checkBasicValidity() error {
    if r.Event == nil || r.Client == "" {
        return fmt.Errorf("invalid reservation (event/client)")
    }
    return nil
}

func (r *Reservation) getTimeslot() Timeslot {
    return r.Event.Timeslots[r.Timeslot]
}

func (r *Reservation) confirmSlots(freeSlots int) {
    if freeSlots >= r.Size {
        r.Confirmed = r.Size
    } else {
        r.Confirmed = freeSlots
		r.Expiration = utils.EpochNow() + c.CONFIG.MAX_PENDIG_RESERVATION_TIME
    }
}

func (r *Reservation) getEventTimeslot() Timeslot {
    return r.Event.Timeslots[r.Timeslot]
}

func (r *Reservation) updateEventTimeslot() {
    timeslot := r.getEventTimeslot()            // Gets timeslot
    timeslot.append(r)                          // Adds reservation to event
    r.Event.Timeslots[r.Timeslot] = timeslot    // Returns updated timeslot
}
