package state

import (
    "fmt"
    "github.com/JValtteri/qure/server/internal/crypt"
)

func MakeReservation(sessionKey crypt.Key, email string, ip IP, size int, eventID crypt.ID, timeslot Epoch) (Reservation, error) {
    var client *Client
    // Try to resume session; if it fails, create a new one
    err := ResumeSession(sessionKey, ip)
    if err != nil {
        client, _ = NewClient("guest", email, crypt.Key(""), true, sessionKey)   // Do not check for conflicting temp client. Both exist
        sessionKey, err = client.AddSession("guest", email, true, ip) // WARNING! session marked as temporary here. This will need to be accounted for!
        if err != nil {
            return Reservation{}, fmt.Errorf("error creating a session for reservation: %v", err)   // Should not be possible (random byte generation)
        }
    }

    // Retrieve the client associated with the session key
    client, found := clients.getClientBySession(sessionKey)
    if !found {
        return Reservation{}, fmt.Errorf("client not found; possible data desynchronization")       // Should not be possible
    }

    // Fetch the event details using the provided event ID
    event, err := GetEvent(eventID)
    if err != nil {
        return Reservation{}, fmt.Errorf("event doesn't exist")
    }

    // Create a new reservation with the client and event details
    reservation, err := newReservation(client, &event, timeslot, size)
    if err != nil {
        return reservation, fmt.Errorf("error creating a reservation: %v", err)                     // Should not be possible (random byte generation)
    }

    // Validate the newly created reservation
    err = reservation.validate()
    return reservation, err
}

func newReservation(client *Client, event *Event, timeslot Epoch, size int) (Reservation, error) {
    newID, err := createUniqueHumanReadableID(10, reservations.byID)
    reservation := Reservation{
        id:         crypt.ID(newID),
        client:     client,
        size:       size,
        confirmed:  0,
        event:      event,
        timeslot:   timeslot,
        expiration: timeslot+RESERVATION_OVERTIME,
    }
    return reservation, err
}

func reservationsFor(userID crypt.ID) []*Reservation {
    return clients.GetReservations(userID)
}
