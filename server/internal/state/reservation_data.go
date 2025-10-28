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

func (r *Reservations) Lock() {
    r.mu.Lock()
}

func (r *Reservations) Unlock() {
    r.mu.Unlock()
}

func (r *Reservations) append(res Reservation) error {
    r.Lock()
    defer r.Unlock()
    r.byID[res.Id] = res
    r.byEmail[res.Client.email] = &res
    return clients.AddReservation(res.Client.Id, &res)
}

var reservations Reservations = Reservations{
    byID:      make(map[crypt.ID]Reservation),
    byEmail:   make(map[string]*Reservation),
}

type Reservation struct {
    Id           crypt.ID
    Client       *Client
    Size         int            // Party size
    Confirmed    int            // Reserved size
    Event        *Event
    Timeslot     Epoch
    Expiration   Epoch
    Error        string
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
    err := reservations.append(*r)                     // Adds reservation to master data

    return err
}

func (r *Reservation) checkBasicValidity() error {
    if r.Event == nil || r.Client == nil {
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
        r.Expiration = utils.EpochNow() + MAX_PENDIG_RESERVATION_TIME
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
