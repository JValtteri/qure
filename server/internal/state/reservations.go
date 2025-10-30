package state

import (
    "fmt"
    "github.com/JValtteri/qure/server/internal/crypt"
)


func ResetEvents() {
    events = make(map[crypt.ID]Event)
}

func MakeReservation(sessionKey crypt.Key, email string, ip IP, size int, eventID crypt.ID, timeslot Epoch) Reservation {
    var client *Client
    // Try to resume session; if it fails, create a new one
    client, err := ResumeSession(sessionKey, ip)
    if err != nil {
        client, _ = NewClient("guest", email, crypt.Key(""), true)    // Does not check for conflicting temp client. Both exist
        sessionKey, err = client.AddSession("guest", email, true, ip) // WARNING! session marked as temporary here. This will need to be accounted for!
        if err != nil {
            return Reservation{Error: fmt.Sprintf("error creating a session for reservation: %v", err)}   // Should not be possible (random byte generation)
        }
    }

    // Fetch the event details using the provided event ID
    event, err := GetEvent(eventID, true)                                                                 // We're assuming that only those authorized have the event id.
    if err != nil {
        return Reservation{Error: fmt.Sprintf("event doesn't exist: %v", err)}
    }

    // Create a new reservation with the client and event details
    reservation, err := newReservation(client, &event, timeslot, size)
    if err != nil {
        return Reservation{Error: fmt.Sprintf("error creating a reservation: %v", err)}                     // Should not be possible (random byte generation)
    }

    // Validate the newly created reservation
    err = reservation.validate()
    if err != nil {
        reservation.Error = fmt.Sprint(err)
    }
    return reservation
}

func newReservation(client *Client, event *Event, timeslot Epoch, size int) (Reservation, error) {
    newID, err := createUniqueHumanReadableID(10, reservations.byID)
    reservation := Reservation{
        Id:         crypt.ID(newID),
        Client:     client,
        Size:       size,
        Confirmed:  0,
        Event:      event,
        Timeslot:   timeslot,
        Expiration: timeslot+RESERVATION_OVERTIME,
    }
    return reservation, err
}

func reservationsFor(userID crypt.ID) []*Reservation {
    client, found := GetClientByID(userID)
    if !found {
        return nil
    }
    return client.GetReservations()
}
