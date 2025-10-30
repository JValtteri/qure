package middleware

import (
	"log"

	"github.com/JValtteri/qure/server/internal/crypt"
	"github.com/JValtteri/qure/server/internal/state"
)


func MakeEvent(req EventCreationRequest) EventCreationResponse {
	auth := AuthenticateSession(AuthenticateRequest{req.SessionKey, req.Ip})
	if !auth.Authenticated || !auth.IsAdmin {
		return EventCreationResponse{crypt.ID("")}
	}
	id, err := state.CreateEvent(req.Event)
	if err != nil {
		log.Printf("Error creating event from JSON: %v\n", err)
	}
	return EventCreationResponse{id}
}
