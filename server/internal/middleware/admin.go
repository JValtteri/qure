package middleware

import (
	"log"

	"github.com/JValtteri/qure/server/internal/crypt"
	"github.com/JValtteri/qure/server/internal/state"
)


func MakeEvent(eventJson []byte, isAdmin bool) crypt.ID {
	if !isAdmin {
		return crypt.ID("")
	}
	id, err := state.CreateEvent(eventJson)
	if err != nil {
		log.Printf("Error creating event from JSON: %v\n", err)
	}
	return id
}
