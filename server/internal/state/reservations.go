package state

import (
	"fmt"
	"github.com/JValtteri/qure/server/internal/utils"
	"github.com/JValtteri/qure/server/internal/crypt"
	c "github.com/JValtteri/qure/server/internal/config"
	"github.com/JValtteri/qure/server/internal/state/model"
)


func MakeReservation(
	sessionKey			crypt.Key,
	email 				string,
	fingerprint 		string,
	hashedFingerprint	crypt.Hash,
	size 				int,
	eventID 			crypt.ID,
	timeslot 			utils.Epoch,
	reservationID		crypt.ID,		// If ID is given, attempt to modifier the reservation
) model.Reservation {
	var temp bool = false
	var isNewReservation bool = reservationID == ""
	var client *model.Client

	// Fetch the event details using the provided event ID
	event, err := GetEvent(eventID, true)														// We're assuming that only those authorized have the event id.
	if err != nil {
		return model.Reservation{Error: fmt.Sprintf("event doesn't exist: %v", err)}
	}

	// Try to resume session; if it fails, create a new temp client
	client, err = ResumeSession(sessionKey, fingerprint)
	if err != nil {
		if isNewReservation {
			// Check reservation size is valid
			if size < 1 {
				return model.Reservation{Error: fmt.Sprintf("invalid size: '%v'", size)}
			}
			temp = true
			client, sessionKey, err = newTempClient(email, hashedFingerprint, client, sessionKey)
			if err != nil {
				return model.Reservation{Error: fmt.Sprint(err)}
			}
		} else {
			return model.Reservation{Error: fmt.Sprintf("couldn't validate session: %v", err)}
		}
	}

	// Create a new reservation with the client and event details
	reservation, err := newReservation(client, &event, timeslot, size)
	if err != nil {
		return model.Reservation{Error: fmt.Sprintf("error creating a reservation: %v", err)}	// Should not be possible (random byte generation)
	}

	// Create a new ID if none was given
	if isNewReservation {
		reservationID, err = model.CreateUniqueHumanReadableID(10, reservations.ByID)
	}
	reservation.Id = crypt.ID(reservationID)

	// Validate the newly created reservation
	if isNewReservation {
		err = reservation.Register(&reservations, &clients)
	} else {
		err = reservation.Amend(&reservations, &clients)
	}
	if err != nil {
		reservation.Error = fmt.Sprint(err)
	}

	// Make user password same as reservation ID for new temp user
	if temp {
		ChangeClientPassword(client, crypt.Key(reservation.Id))
	}

	reservation.Session = sessionKey							// This is to provide the session key in when a session is created simultaneously
	return reservation
}

func CancelReservation(sessionKey			crypt.Key,
	email 				string,
	fingerprint 		string,
	hashedFingerprint	crypt.Hash,
	size 				int,
	eventID 			crypt.ID,
	timeslot 			utils.Epoch,
	reservationID		crypt.ID,		// If ID is given, attempt to modifier the reservation
) model.Reservation {
	var isNewReservation bool = reservationID == ""
	if isNewReservation {
		return model.Reservation{Error: fmt.Sprintln("No change")}
	}
	var client *model.Client

	// Fetch the event details using the provided event ID
	event, err := GetEvent(eventID, true)														// We're assuming that only those authorized have the event id.
	if err != nil {
		return model.Reservation{Error: fmt.Sprintf("event doesn't exist: %v", err)}
	}

	// Try to resume session; if it fails, create a new temp client
	client, err = ResumeSession(sessionKey, fingerprint)
	if err != nil {
		return model.Reservation{Error: fmt.Sprintf("couldn't validate session: %v", err)}
	}

	// Create a new reservation object with the client and event details
	reservation, err := newReservation(client, &event, timeslot, size)
	if err != nil {
		return model.Reservation{Error: fmt.Sprintf("error creating a reservation: %v", err)}	// Should not be possible (random byte generation)
	}

	reservation.Id = crypt.ID(reservationID)
	err = reservation.Cancel(&reservations, &clients)
	if err != nil {
		reservation.Error = fmt.Sprint(err)
	}

	reservation.Session = sessionKey							// This is to provide the session key in when a session is created simultaneously
	return reservation
}

func newTempClient(
	email				string,
	hashedFingerprint	crypt.Hash,
	client 				*model.Client,
	sessionKey			crypt.Key,
) (*model.Client, crypt.Key, error) {
	client, err := NewClient("guest", email, crypt.Key(""), true)								// Does not check for conflicting temp client. Both exist
	if err != nil {
		return client, crypt.Key(""), fmt.Errorf("error creating new client for reservation: %v", err)
	}
	sessionKey, err = client.AddSession("guest", email, true, hashedFingerprint, &clients)	// WARNING! session marked as temporary here. This will need to be accounted for!
	if err != nil {
		return client, crypt.Key(""), fmt.Errorf("error creating a session for reservation: %v", err)	// Should not be possible (random byte generation)
	}
	return client, sessionKey, nil
}

func newReservation(
	client		*model.Client,
	event		*model.Event,
	timeslot	utils.Epoch,
	size 		int,
) (model.Reservation, error) {
	if client == nil || event == nil {
		return model.Reservation{}, fmt.Errorf("error creating reservation: Client or Event is <nil>")
	}
	reservation := model.Reservation{
		Id:			crypt.ID(""),
		Client:		client.Id,
		Size:		size,
		Confirmed:	0,
		Event:		event,
		Timeslot:	timeslot,
		Expiration:	timeslot + c.CONFIG.RESERVATION_OVERTIME,
	}
	return reservation, nil
}

func reservationsFor(userID crypt.ID) []*model.Reservation {
	client, found := GetClientByID(userID)
	if !found {
		return nil
	}
	return client.GetReservations()
}
