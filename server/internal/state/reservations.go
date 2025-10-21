package state

import (
    "sync"
    "fmt"
    "github.com/JValtteri/qure/server/internal/utils"
)

var MAX_PENDIG_RESERVATION_TIME Epoch = 60*10   // seconds
var RESERVATION_OVERTIME Epoch = 60*60 // How long the reservation is kept in system after the reservation slot has started

type Reservations struct {
    mu          sync.RWMutex
    byID        map[ID]Reservation
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
    byID:      make(map[ID]Reservation),
    byEmail:   make(map[string]*Reservation),
}

type Reservation struct {
    id           ID
    client       *Client
    size         int            // Party size
    confirmed    int            // Reserved size
    event        *Event
    timeslot     Epoch
    expiration   Epoch
}

func (r *Reservation)validate() error {
    if r.event == nil || r.client == nil {
        return fmt.Errorf("invalid reservation (event/client)")
    }
    // check size/slot
    eventslock.Lock()
    defer eventslock.Unlock()
    targetSlot := r.event.Timeslots[r.timeslot]
    if targetSlot.isFull() {
        return fmt.Errorf("slot full")
    }
    freeSlots := targetSlot.hasFree()
    if freeSlots >= r.size {
        r.confirmed = r.size
    }
    if freeSlots < r.size {
        r.confirmed = freeSlots
        r.expiration = utils.EpochNow() + MAX_PENDIG_RESERVATION_TIME
    }
    timeslot := r.event.Timeslots[r.timeslot]   // Gets timeslot
    timeslot.append(r)                          // Adds reservation to event
    r.event.Timeslots[r.timeslot] = timeslot    // Returns updated timeslot
    reservations.append(*r)                     // Adds reservation to master data
    return nil
}

func MakeReservation(sessionKey Key, email string, ip IP, size int, eventID ID, timeslot Epoch) (Reservation, error) {
    var reservation Reservation
    err := ResumeSession(sessionKey, ip)
    if err != nil {
        sessionKey, err = AddSession("guest", email, true, ip) // WARNING! session marked as temporary here. This will need to be accounted for!
        if err != nil {
            return reservation, fmt.Errorf("error creating a session for reservation: %v", err) // Should not be possible (random byte generation)
        }
    }
    client, found := clients.getClient(sessionKey)
    if !found {
        return reservation, fmt.Errorf("client not found")  // Should not be possible (Data desync)
    }
    event, err := GetEvent(eventID)
    if err != nil {
        return reservation, fmt.Errorf("event doesn't exist")
    }
    reservation, err = newReservation(client, &event, timeslot, size)
    if err != nil {
        return reservation, fmt.Errorf("error creating a reservation: %v", err) // Should not be possible (random byte generation)
    }
    // Add reservation to every place
    // or get an error for not enough room
    // implement partial reservations
    // with option to respond
    err = reservation.validate()
    return reservation, err
}

func newReservation(client *Client, event *Event, timeslot Epoch, size int) (Reservation, error) {
    newID, err := createHumanReadableId(10)
    reservation := Reservation{
        id:         ID(newID),
        client:     client,
        size:       size,
        confirmed:  0,
        event:      event,
        timeslot:   timeslot,
        expiration: timeslot+RESERVATION_OVERTIME,
    }
    return reservation, err
}

func reservationsFor(userID ID) []*Reservation {
    return clients.GetReservations(userID)
}
