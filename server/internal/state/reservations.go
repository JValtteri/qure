package state

import (
	"fmt"
	"github.com/JValtteri/qure/server/internal/utils"
	"github.com/JValtteri/qure/server/internal/crypt"
	c "github.com/JValtteri/qure/server/internal/config"
	"github.com/JValtteri/qure/server/internal/state/model"
)


// Handles both NEW and AMENDED reservations
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
	var isNewReservation bool = reservationID == ""
	if !isNewReservation && !validID(reservationID) {
		return createErrorResponse(fmt.Errorf("Invalid reservation ID: '%v'", reservationID))
	}
	if !isValidSize(size) {
		return createErrorResponse(fmt.Errorf("invalid size: '%v'", size))
	}

	client, sessionKey, err := getSessionOrCreateTempClient(sessionKey, fingerprint, email, hashedFingerprint, isNewReservation)
	if err != nil {
		return createErrorResponse(err)
	}
	reservation, err := newReservationObject(client, eventID, timeslot, size)
	if err != nil {
		return createErrorResponse(fmt.Errorf("error creating a reservation: %v", err)) // Should not be possible (random byte generation)
	}

	if isNewReservation {
		reservationID, err = model.CreateUniqueHumanReadableID(10, reservations.ByID)
	}
	reservation.Id = crypt.ID(reservationID)

	err = saveOrUpdateReservation(&reservation, &reservations, &clients, isNewReservation)
	if err != nil {
		reservation.Error = fmt.Sprint(err)
	}

	if isTempClient(client) {
		ChangeClientPassword(client, crypt.Key(reservation.Id))
	}

	reservation.Session = sessionKey // This is to provide the session key in when a session is created simultaneously
	return reservation
}

func validID(reservationID crypt.ID) bool {
	_, valid := reservations.ByID[reservationID]
	return valid
}

func isValidSize(size int) bool {
	return size >= 1
}

func createErrorResponse(err error) model.Reservation {
	return model.Reservation{Error: fmt.Sprintf("%v", err)}
}

func getSessionOrCreateTempClient(
	sessionKey crypt.Key,
	fingerprint string,
	email string,
	hashedFingerprint crypt.Hash,
	createNewTempClient bool,
) (*model.Client, crypt.Key, error) {

	client, err := ResumeSession(sessionKey, fingerprint)
	if err != nil {
		if createNewTempClient {
			tempClient, sessionKey, err := newTempClient(email, hashedFingerprint, client, sessionKey)
			return tempClient, sessionKey, err
		}
		return nil, "", fmt.Errorf("couldn't validate session: %v", err)
	}

	return client, sessionKey, nil
}

// Directs process to either Register new reservation or Amend existing
func saveOrUpdateReservation(reservation *model.Reservation, reservations *model.Reservations, clients *model.Clients, isNew bool) error {
	var err error
	if isNew {
		err = reservation.Register(reservations, clients)
		return err
	} else {
		err = reservation.Amend(reservations, clients)
	}
	return err
}

func isTempClient(client *model.Client) bool {
	return client.IsTemporary
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
		return model.Reservation{Error: fmt.Sprintln("Missing reservation ID")}
	}
	if !validID(reservationID) {
		return createErrorResponse(fmt.Errorf("Invalid reservation ID: '%v'", reservationID))
	}
	var client *model.Client

	// Try to resume session; if it fails, create a new temp client
	client, err := ResumeSession(sessionKey, fingerprint)
	if err != nil {
		return model.Reservation{Error: fmt.Sprintf("couldn't validate session: %v", err)}
	}

	// Create a new reservation object with the client and event details
	reservation, err := newReservationObject(client, eventID, timeslot, size)
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

func newReservationObject(
	client		*model.Client,
	eventID		crypt.ID,
	timeslot	utils.Epoch,
	size 		int,
) (model.Reservation, error) {
	event, err := GetEvent(eventID, true)														// We're assuming that only those authorized have the event id.
	if err != nil {
		return model.Reservation{}, fmt.Errorf("event doesn't exist: %v", err)
	}
	if client == nil {
		return model.Reservation{}, fmt.Errorf("error creating reservation: Client <nil>")
	}
	reservation := model.Reservation{
		Id:			crypt.ID(""),
		Client:		client.Id,
		Size:		size,
		Confirmed:	0,
		Event:		&event,
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
