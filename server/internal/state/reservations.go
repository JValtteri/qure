package state

import (
    "fmt"
    "github.com/JValtteri/qure/server/internal/crypt"
)


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
    event, err := GetEvent(eventID)
    if err != nil {
        return Reservation{Error: fmt.Sprintf("event doesn't exist")}
    }

    // Create a new reservation with the client and event details
    reservation, err := newReservation(client, &event, timeslot, size)
    if err != nil {
        return Reservation{Error: fmt.Sprintf("error creating a reservation: %v", err)}                     // Should not be possible (random byte generation)
    }

    // Validate the newly created reservation
    reservation.Error = fmt.Sprint(reservation.validate())
    return reservation
}

func newReservation(client *Client, event *Event, timeslot Epoch, size int) (Reservation, error) {
    newID, err := createUniqueHumanReadableID(10, reservations.byID)
    reservation := Reservation{
        Id:         crypt.ID(newID),
        Client:     client,
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
