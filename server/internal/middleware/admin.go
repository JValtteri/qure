package middleware

import (
	"fmt"
	"log"

	"github.com/JValtteri/qure/server/internal/crypt"
	"github.com/JValtteri/qure/server/internal/state"
	"github.com/JValtteri/qure/server/internal/state/model"
)


func MakeEvent(rq EventManipulationRequest) EventManipulationResponse {
	authorized, err := adminAuthority(rq.SessionKey, rq.Fingerprint)
	if !authorized {
		return EventManipulationResponse{crypt.ID(""), fmt.Sprintf("%v", err)}
	}
	id, err := state.CreateEvent(rq.Event)
	if err != nil {
		log.Printf("Error creating event from JSON: %v\n", err)
	}
	return EventManipulationResponse{id, ""}
}

func EditEvent(rq EventManipulationRequest) EventManipulationResponse {
	authorized, err := adminAuthority(rq.SessionKey, rq.Fingerprint)
	if !authorized {
		return EventManipulationResponse{crypt.ID(""), fmt.Sprintf("%v", err)}
	}
	id, err := state.EditEvent(rq.Event)
	if err != nil {
		return EventManipulationResponse{id, fmt.Sprintf("Error editing event: %v\n", err)}
	}
	return EventManipulationResponse{id, ""}
}

func DeleteEvent(rq EventManipulationRequest) EventManipulationResponse {
	authorized, err := adminAuthority(rq.SessionKey, rq.Fingerprint)
	if !authorized {
		return EventManipulationResponse{crypt.ID(""), fmt.Sprintf("%v", err)}
	}
	ok := state.RemoveEvent(rq.Event.ID)
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
	eventReservations, err := state.GetEventReservations(req.EventID)
	if err != nil {
		return response
	}
	for _, res := range eventReservations {
		response = append(response, reservationToResponse(res))
	}

	return response
}

// Requests list of all users. Access attempts to PII are logged to comply with GDPR Article 33:1,3 and Article 34:6
func ListAllUsers(rq AuthenticateRequest) []model.Client {
	var response []model.Client
	var reqUser, found = state.GetClientBySession(rq.SessionKey)
	var action = "list of all users"
	var username = ""

	if found {
		username = reqUser.Email
	}
	authorized, _ := adminAuthority(rq.SessionKey, rq.Fingerprint)
	if !authorized {
		log.Print(logGdprFormat(username, action, false))
		return response
	}
	log.Print(logGdprFormat(username, action, true))
	var clients = state.GetAllClients(authorized)

	return clients
}

func AdminListUserReservatoions(rq EnhancedUserRequest) []ReservationResponse {
	var failure []ReservationResponse
	var action = fmt.Sprintf("list user: '%s's' reservations", rq.User)

	authorized, username, _ := enhancedAdminAuthority(rq.Password, rq.SessionKey, rq.HashPrint)
	if !authorized {
		log.Print(logGdprFormat(username, action, false))
		return failure
	}
	log.Print(logGdprFormat(username, action, true))

	client, found := state.GetClientByEmail(rq.User)
	var response []ReservationResponse
	if !found {
		return response
	}
	reservations := client.GetReservations()
	for _, value := range(reservations) {
		response = append(response, reservationToResponse(*value))
	}
	return response
}

func AdminChangeUserRole(rq RoleChangeRequest) SuccessResponse {
	var failure = SuccessResponse{
		Success: false,
		Error: "Authentication failed",
	}
	var action = fmt.Sprintf("change user: '%s's' role to '%s'", rq.User, rq.Role)

	authorized, username, _ := enhancedAdminAuthority(rq.Password, rq.SessionKey, rq.HashPrint)
	if !authorized {
		log.Print(logGdprFormat(username, action, false))
		return failure
	}
	log.Print(logGdprFormat(username, action, true))

	client, found := state.GetClientByEmail(rq.User)
	if !found {
		failure.Error = fmt.Sprintf("User '%v' not found\n", rq.User)
		return failure
	}
	state.ChangeClientRole(client, rq.Role)

	return SuccessResponse{
		Success: true,
	}
}

func AdminRemoveUser(rq EnhancedUserRequest) SuccessResponse {
	var failure = SuccessResponse{
		Success: false,
		Error: "Authentication failed",
	}
	var action = fmt.Sprintf("delete user: '%s's' role", rq.User)

	authorized, username, _ := enhancedAdminAuthority(rq.Password, rq.SessionKey, rq.HashPrint)
	if !authorized {
		log.Print(logGdprFormat(username, action, false))
		return failure
	}
	log.Print(logGdprFormat(username, action, true))
	client, found := state.GetClientByEmail(rq.User)
	if !found {
		failure.Error = fmt.Sprintf("User '%v' not found\n", rq.User)
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
		log.Printf(
			"Admin authentication failed: Auth: %v, Role: %v, Key: %v, authError: %v\n",
			auth.Authenticated, auth.Role, sessionKey, auth.Error,
		)
		return false, fmt.Errorf("Authentication failed")
	}
	return true, nil
}

// Checks for valid admin authority and password
func enhancedAdminAuthority(password crypt.Key, sessionKey crypt.Key, hashPrint crypt.Hash) (bool, string, error) {
	admin, found := state.GetClientBySession(sessionKey)
	if !found {
		log.Printf(
			"Admin authentication failed: session not found Key: %v\n", sessionKey)
		return false, "", fmt.Errorf("Authentication failed")
	}
	var auth = checkPasswordAuthentication(admin, password, hashPrint)
	if !auth.Authenticated || auth.Role != "admin" {
		log.Printf(
			"Admin authentication failed: Auth: %v, Role: %v, Key: %v, authError: %v\n",
			auth.Authenticated, auth.Role, sessionKey, auth.Error,
		)
		return false, admin.Email, fmt.Errorf("Authentication failed")
	}
	return true, admin.Email, nil
}

func logGdprFormat(username string, requestedAction string, success bool) string {
	if username == "" {
		username = "[Unknown]"
	}
	var successTxt = "[AUTHENTICATION FAILED]"
	if success {
		successTxt = ""
	}
	return fmt.Sprintf("[GDPR]: '%v' requested %s. %v\n", username, requestedAction, successTxt)
}
