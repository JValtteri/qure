package state

import (
    "sync"
    "fmt"
    "github.com/JValtteri/qure/server/internal/utils"
    "github.com/JValtteri/qure/server/internal/crypt"
)

var MAX_PENDIG_RESERVATION_TIME Epoch = 60*10   // seconds
var RESERVATION_OVERTIME Epoch = 60*60          // the time a reservation is kept past reservation start time

type Reservations struct {
    mu          sync.RWMutex
    byID        map[crypt.ID]Reservation
    byEmail     map[string]*Reservation
}

func (r *Reservations) rLock() {
    r.mu.RLock()
}

func (r *Reservations) rUnlock() {
    r.mu.RUnlock()
}

func (r *Reservations) Lock() {
    r.mu.Lock()
}

func (r *Reservations) Unlock() {
    r.mu.Unlock()
}

func (r *Reservations) append(res Reservation) {
    r.Lock()
    defer r.Unlock()
    r.byID[res.id] = res
    r.byEmail[res.client.email] = &res
    clients.AddReservation(res.client.id, &res)
}

var reservations Reservations = Reservations{
    byID:      make(map[crypt.ID]Reservation),
    byEmail:   make(map[string]*Reservation),
}

type Reservation struct {
    id           crypt.ID
    client       *Client
    size         int            // Party size
    confirmed    int            // Reserved size
    event        *Event
    timeslot     Epoch
    expiration   Epoch
}

// Propagets reservation OR get an error for not enough room
// implements partial reservations
func (r *Reservation) validate() error {
    if err := r.checkBasicValidity(); err != nil {
        return fmt.Errorf("invalid reservation (event/client)")
    }
    eventslock.Lock()
    defer eventslock.Unlock()

    slot := r.getTimeslot()
    if slot.isFull() {
        return fmt.Errorf("slot full")
    }
    freeSlots := slot.hasFree()
    r.confirmSlots(freeSlots)
    r.updateEventTimeslot()
    reservations.append(*r)                     // Adds reservation to master data
    return nil
}

func (r *Reservation) checkBasicValidity() error {
    if r.event == nil || r.client == nil {
        return fmt.Errorf("invalid reservation (event/client)")
    }
    return nil
}

func (r *Reservation) getTimeslot() Timeslot {
    return r.event.Timeslots[r.timeslot]
}

func (r *Reservation) confirmSlots(freeSlots int) {
    if freeSlots >= r.size {
        r.confirmed = r.size
    } else {
        r.confirmed = freeSlots
        r.expiration = utils.EpochNow() + MAX_PENDIG_RESERVATION_TIME
    }
}

func (r *Reservation) getEventTimeslot() Timeslot {
    return r.event.Timeslots[r.timeslot]
}

func (r *Reservation) updateEventTimeslot() {
    timeslot := r.getEventTimeslot()            // Gets timeslot
    timeslot.append(r)                          // Adds reservation to event
    r.event.Timeslots[r.timeslot] = timeslot    // Returns updated timeslot
}
