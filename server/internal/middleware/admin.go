package middleware

import (
	"fmt"
	"log"

	"github.com/JValtteri/qure/server/internal/crypt"
	"github.com/JValtteri/qure/server/internal/state"
)


func MakeEvent(req EventManipulationRequest) EventManipulationResponse {
	authorized, response := adminAuthority(req)
	if !authorized {
		return response
	}
	id, err := state.CreateEvent(req.Event)
	if err != nil {
		log.Printf("Error creating event from JSON: %v\n", err)
	}
	return EventManipulationResponse{id, ""}
}

func EditEvent(req EventManipulationRequest) EventManipulationResponse {
	authorized, response := adminAuthority(req)
	if !authorized {
		return response
	}
	id, err := state.EditEvent(req.Event)
	if err != nil {
		return EventManipulationResponse{id, fmt.Sprintf("Error editing event: %v\n", err)}
	}
	return EventManipulationResponse{id, ""}
}

func DeleteEvent(req EventManipulationRequest) EventManipulationResponse {
	authorized, response := adminAuthority(req)
	if !authorized {
		return response
	}
	ok := state.RemoveEvent(req.Event.ID)
	if !ok {
		return EventManipulationResponse{"", "Event not found"}
	}
	return EventManipulationResponse{"", ""}
}

// Checks for valid abmin authority
func adminAuthority(req EventManipulationRequest) (bool, EventManipulationResponse) {
	auth := AuthenticateSession(AuthenticateRequest{req.SessionKey, req.Fingerprint})
	if !auth.Authenticated || !auth.IsAdmin {
		return false, EventManipulationResponse{
			crypt.ID(""), fmt.Sprintf(
				"Authentication failed: Auth: %v, Admin: %v, Key: %v, authError: %v",
				auth.Authenticated, auth.IsAdmin, req.SessionKey, auth.Error,
			)}
	}
	return true, EventManipulationResponse{}
}
