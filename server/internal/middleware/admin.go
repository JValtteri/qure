package middleware

import (
	"fmt"
	"log"

	"github.com/JValtteri/qure/server/internal/crypt"
	"github.com/JValtteri/qure/server/internal/state"
)


func MakeEvent(req EventManipulationRequest) EventManipulationResponse {
	authorized, err := adminAuthority(req.SessionKey, req.Fingerprint)
	if !authorized {
		return EventManipulationResponse{crypt.ID(""), fmt.Sprintf("%v", err)}
	}
	id, err := state.CreateEvent(req.Event)
	if err != nil {
		log.Printf("Error creating event from JSON: %v\n", err)
	}
	return EventManipulationResponse{id, ""}
}

func EditEvent(req EventManipulationRequest) EventManipulationResponse {
	authorized, err := adminAuthority(req.SessionKey, req.Fingerprint)
	if !authorized {
		return EventManipulationResponse{crypt.ID(""), fmt.Sprintf("%v", err)}
	}
	id, err := state.EditEvent(req.Event)
	if err != nil {
		return EventManipulationResponse{id, fmt.Sprintf("Error editing event: %v\n", err)}
	}
	return EventManipulationResponse{id, ""}
}

func DeleteEvent(req EventManipulationRequest) EventManipulationResponse {
	authorized, err := adminAuthority(req.SessionKey, req.Fingerprint)
	if !authorized {
		return EventManipulationResponse{crypt.ID(""), fmt.Sprintf("%v", err)}
	}
	ok := state.RemoveEvent(req.Event.ID)
	if !ok {
		return EventManipulationResponse{"", "Event not found"}
	}
	return EventManipulationResponse{"", ""}
}

func GetEventReservations(req EventRequest) []ReservationResponse {
	var response []ReservationResponse

	authorized, err := adminAuthority(req.SessionKey, req.Fingerprint)
	if !authorized {
		return response
	}
	resss, err := state.GetEventReservations(req.EventID, authorized)
	for _, res := range resss {
		response = append(response, reservationToResponse(res))
	}

	if err != nil {
		return response
	}

	return response
}



// Checks for valid admin authority
func adminAuthority(sessionKey crypt.Key, fingerprint string) (bool, error) { //} req EventManipulationRequest) (bool, EventManipulationResponse) {
	auth := AuthenticateSession(AuthenticateRequest{sessionKey, fingerprint})
	if !auth.Authenticated || auth.Role != "admin" {
		return false, fmt.Errorf(
			"Authentication failed: Auth: %v, Role: %v, Key: %v, authError: %v",
			auth.Authenticated, auth.Role, sessionKey, auth.Error,
		)
	}
	return true, nil
}
