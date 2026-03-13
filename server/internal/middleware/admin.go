package middleware

import (
	"fmt"
	"log"

	"github.com/JValtteri/qure/server/internal/crypt"
	"github.com/JValtteri/qure/server/internal/state"
	"github.com/JValtteri/qure/server/internal/state/model"
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
	eventReservations, err := state.GetEventReservations(req.EventID, authorized)
	if err != nil {
		return response
	}
	for _, res := range eventReservations {
		response = append(response, reservationToResponse(res))
	}

	return response
}

// Requests list of all users. Access attempts to PII are logged to comply with GDPR Article 33:1,3 and Article 34:6
func ListAllUsers(req AuthenticateRequest) []model.Client {
	var response []model.Client
	var reqUser, found = state.GetClientBySession(req.SessionKey)
	var username = "[Unknown]"
	var successTxt = "[AUTHENTICATION FAILED]"
	if found {
		username = reqUser.Email
		successTxt = ""
	}
	var gdprLogTxt = fmt.Sprintf("[GDPR]: '%v' requested list of all users. %v\n", username, successTxt)

	authorized, _ := adminAuthority(req.SessionKey, req.Fingerprint)
	if !authorized {
		log.Print(gdprLogTxt)
		return response
	}
	log.Print(gdprLogTxt)
	var clients = state.GetAllClients(authorized)

	return clients
}

func AdminChangeUserRole(rq RoleChangeRequest) SuccessResponse {
	var failure = SuccessResponse{
		Success: false,
		Error: "Authentication failed",
	}
	authorized, _ := adminAuthority(rq.SessionKey, rq.Fingerprint)
	if !authorized {
		return failure
	}
	admin, found := state.GetClientBySession(rq.SessionKey)
	if !found {
		log.Printf("User '%v' not found\n", rq.User)
		return failure
	}
	// Require additional password confirmation to delete a user
	var auth = checkPasswordAuthentication(admin, rq.Password, rq.HashPrint)
	if !auth.Authenticated {
		return failure
	}
	client, found := state.GetClientByEmail(rq.User)
	if !found {
		log.Printf("User '%v' not found\n", rq.User)
		return failure
	}
	state.ChangeClientRole(client, rq.Role)

	return SuccessResponse{
		Success: true,
	}
}

func AdminRemoveUser(rq RemovalRequest) SuccessResponse {
	var failure = SuccessResponse{
		Success: false,
		Error: "Authentication failed",
	}
	authorized, _ := adminAuthority(rq.SessionKey, rq.Fingerprint)
	if !authorized {
		return failure
	}
	admin, found := state.GetClientBySession(rq.SessionKey)
	if !found {
		log.Printf("User '%v' not found\n", rq.User)
		return failure
	}
	// Require additional password confirmation to delete a user
	var auth = checkPasswordAuthentication(admin, rq.Password, rq.HashPrint)
	if !auth.Authenticated {
		return failure
	}
	client, found := state.GetClientByEmail(rq.User)
	if !found {
		log.Printf("User '%v' not found\n", rq.User)
		return failure
	}
	state.RemoveClient(client)

	return SuccessResponse{
		Success: true,
	}
}

// Checks for valid admin authority
func adminAuthority(sessionKey crypt.Key, fingerprint string) (bool, error) {
	auth := AuthenticateSession(AuthenticateRequest{sessionKey, fingerprint})
	if !auth.Authenticated || auth.Role != "admin" {
		return false, fmt.Errorf(
			"Authentication failed: Auth: %v, Role: %v, Key: %v, authError: %v",
			auth.Authenticated, auth.Role, sessionKey, auth.Error,
		)
	}
	return true, nil
}
