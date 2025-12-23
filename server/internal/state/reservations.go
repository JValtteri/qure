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
) model.Reservation {
	var client *model.Client
	// Try to resume session; if it fails, create a new temp client
	client, err := ResumeSession(sessionKey, fingerprint)
	if err != nil {
		client, _ = NewClient("guest", email, crypt.Key(""), true)								// Does not check for conflicting temp client. Both exist
		sessionKey, err = client.AddSession("guest", email, true, hashedFingerprint, &clients)	// WARNING! session marked as temporary here. This will need to be accounted for!
		if err != nil {
			return model.Reservation{Error: fmt.Sprintf("error creating a session for reservation: %v", err)}	// Should not be possible (random byte generation)
		}
	}

	// Fetch the event details using the provided event ID
	event, err := GetEvent(eventID, true)														// We're assuming that only those authorized have the event id.
	if err != nil {
		return model.Reservation{Error: fmt.Sprintf("event doesn't exist: %v", err)}
	}

	// Create a new reservation with the client and event details
	reservation, err := newReservation(client, &event, timeslot, size)
	if err != nil {
		return model.Reservation{Error: fmt.Sprintf("error creating a reservation: %v", err)}	// Should not be possible (random byte generation)
	}

	// Validate the newly created reservation
	err = reservation.Validate(&reservations, &clients)
	if err != nil {
		reservation.Error = fmt.Sprint(err)
	}
	reservation.Session = sessionKey							// This is to provide the session key in when a session is created simultaneously
	return reservation
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
	newID, err := model.CreateUniqueHumanReadableID(10, reservations.ByID)
	reservation := model.Reservation{
		Id:			crypt.ID(newID),
		Client:		client.Id,
		Size:		size,
		Confirmed:	0,
		Event:		event,
		Timeslot:	timeslot,
		Expiration:	timeslot + c.CONFIG.RESERVATION_OVERTIME,
	}
	return reservation, err
}

func reservationsFor(userID crypt.ID) []*model.Reservation {
	client, found := GetClientByID(userID)
	if !found {
		return nil
	}
	return client.GetReservations()
}
