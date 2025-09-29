package main


import (
)

var reservations []Reservation
var reservationsByClient map[string]*Reservation = make(map[string]*Reservation)

type Reservation struct {
    clientID     string
    email        string
    eventID      string
    id           int      // not sure about this
}

/*
func createReservation() bool {
    // TODO
    return false
}

func findReservations(userID string) []*Reservation {
    // TODO
    return nil
}

*/

