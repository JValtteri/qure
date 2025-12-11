package middleware

import (
	"fmt"
	"log"

	"github.com/JValtteri/qure/server/internal/crypt"
	"github.com/JValtteri/qure/server/internal/state"
)


func MakeEvent(req EventCreationRequest) EventCreationResponse {
	auth := AuthenticateSession(AuthenticateRequest{req.SessionKey, req.Fingerprint})
	if !auth.Authenticated || !auth.IsAdmin {
		return EventCreationResponse{
			crypt.ID(""), fmt.Sprintf(
				"Authentication failed: Auth: %v, Admin: %v, Key: %v, authError: %v",
				auth.Authenticated, auth.IsAdmin, req.SessionKey, auth.Error,
			)}
	}
	id, err := state.CreateEvent(req.Event)
	if err != nil {
		log.Printf("Error creating event from JSON: %v\n", err)
	}
	return EventCreationResponse{id, ""}
}

func EditEvent(req EventCreationRequest) EventCreationResponse {
	auth := AuthenticateSession(AuthenticateRequest{req.SessionKey, req.Fingerprint})
	if !auth.Authenticated || !auth.IsAdmin {
		return EventCreationResponse{
			crypt.ID(""), fmt.Sprintf(
				"Authentication failed: Auth: %v, Admin: %v, Key: %v, authError: %v",
				auth.Authenticated, auth.IsAdmin, req.SessionKey, auth.Error,
			)}
	}
	id, err := state.EditEvent(req.Event)
	if err != nil {
		return EventCreationResponse{id, fmt.Sprintf("Error editing event: %v\n", err)}
	}
	return EventCreationResponse{id, ""}
}
