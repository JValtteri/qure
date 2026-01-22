package middleware

import (
	"fmt"
	"log"

	"github.com/JValtteri/qure/server/internal/crypt"
	"github.com/JValtteri/qure/server/internal/state"
	c "github.com/JValtteri/qure/server/internal/config"
	"github.com/JValtteri/qure/server/internal/state/model"
)


// Login with a regular account
func Login(rq LoginRequest) (Authentication) {
    client, found := state.GetClientByEmail(rq.User)
    if !found {
        log.Printf("User '%v' not found\n", rq.User)
        return Authentication{}
    }
    auth := checkPasswordAuthentication(client, rq.Password, rq.HashPrint)
    return auth
}

// Login using a reservation made without an account
// This feature has been removed as redundant
// func ReservationLogin(rq EventLogin) Authentication

func Logout(rq AuthenticateRequest) Authentication {
    auth := Authentication{}
    state.RemoveSession(rq.SessionKey)
    return auth
}

func AuthenticateSession(rq AuthenticateRequest) Authentication {
    auth := Authentication{}
    client, err := state.ResumeSession(rq.SessionKey, rq.Fingerprint)
    if err != nil {
        return auth
    }
    populateAuthObject(&auth, true, client.IsAdmin())
    auth.SessionKey = rq.SessionKey
    auth.User = client.GetEmail()
    return auth
}

func Register(rq RegisterRequest) RegistrationResponse {
    role := "guest"
    temp := false
    if len(rq.User) < c.CONFIG.MIN_USERNAME_LENGTH {
        return RegistrationResponse{
            Error: fmt.Sprintf("Username length must be at least %v characters", c.CONFIG.MIN_USERNAME_LENGTH),
        }
    }
    if len(rq.Password) < c.CONFIG.MIN_PASSWORD_LENGTH {
        return RegistrationResponse{
            Error: fmt.Sprintf("Password length must be at least %v characters", c.CONFIG.MIN_PASSWORD_LENGTH),
        }
    }
    client, err := state.NewClient(role, rq.User, rq.Password, temp)
    if err != nil {
        return RegistrationResponse{Error: fmt.Sprintf("%v", err)}
    }
    key, err := state.AddSession(client, role, rq.User, temp, rq.HashPrint)
    if err != nil {
        return RegistrationResponse{Error: fmt.Sprintf("%v", err)}
    }
    return RegistrationResponse{key, ""}
}

func MakeReservation(rq ReserveRequest) ReservationResponse {
    res := state.MakeReservation(
		rq.SessionKey,	rq.User,	rq.Fingerprint,
		rq.HashPrint,	rq.Size,	rq.EventID,		rq.Timeslot,
    )
	return reservationToResponse(res)	// Here a Reservation object is translated to a ReservationResponse
}

/*
func EditReservation(rq ReserveRequest) ReservationResponse {
    res := state.EditReservation(
		rq.SessionKey,	rq.User,	rq.Fingerprint,
		rq.HashPrint,	rq.Size,	rq.EventID,		rq.Timeslot,
	)
	return reservationToResponse(res)	// Here a Reservation object is translated to a ReservationResponse
}
*/

func ChangePassword(rq PasswordChangeRequest) PasswordChangeResponse {
	if len(rq.NewPassword) < c.CONFIG.MIN_PASSWORD_LENGTH {
		return PasswordChangeResponse{
			Error: fmt.Sprintf("Password length must be at least %v characters", c.CONFIG.MIN_PASSWORD_LENGTH),
		}
	}
	var failure = PasswordChangeResponse{
		Success: false,
		Error: "Authentication failed",
	}
	client, found := state.GetClientByEmail(rq.User)
	if !found {
		log.Printf("User '%v' not found\n", rq.User)
		return failure
	}
	auth := AuthenticateSession(AuthenticateRequest{rq.SessionKey, rq.Fingerprint})
	if !auth.Authenticated {
		return failure
	}
	auth = checkPasswordAuthentication(client, rq.Password, rq.HashPrint)
	if !auth.Authenticated {
		return failure
	}
	state.ChangeClientPassword(client, rq.NewPassword)
	// Reset all sessions as a precaution
	// but renew current sessuion
	client.ClearSessions()
	key, err := state.AddSession(client, client.GetRole(), client.GetEmail(), false, rq.HashPrint)
	if err != nil {
		return PasswordChangeResponse{ Error: fmt.Sprintf("%v", err)}
	}
	return PasswordChangeResponse{ Success: true, SessionKey: key }
}

func RemoveUser(rq RemovalRequest) SuccessResponse {
	var failure = SuccessResponse{
		Success: false,
		Error: "Authentication failed",
	}
	client, found := state.GetClientByEmail(rq.User)
	if !found {
		log.Printf("User '%v' not found\n", rq.User)
		return failure
	}
	auth := AuthenticateSession(AuthenticateRequest{rq.SessionKey, rq.Fingerprint})
	if !auth.Authenticated {
		return failure
	}
	auth = checkPasswordAuthentication(client, rq.Password, rq.HashPrint)
	if !auth.Authenticated {
		return failure
	}
	state.RemoveClient(client)
	return SuccessResponse{
		Success: true,
	}
}

func checkPasswordAuthentication(
	client *model.Client,  password crypt.Key,  hashedFingerprint crypt.Hash,
) Authentication {
    authorized := crypt.CompareToHash(password, client.GetPasswordHash())
    if !authorized {
        log.Printf("Client's '%v' password didn't match\n", client.GetEmail())
        return Authentication{}
    }
    auth := Authentication{}
    populateAuthObject(&auth, authorized, client.IsAdmin())
    key, err := state.AddSession(client, client.GetRole(), client.GetEmail(), false, hashedFingerprint)
    if err != nil {
        return Authentication{Error: fmt.Sprintf("%v", err)}
    }
    auth.User = client.GetEmail()
    auth.SessionKey = key
    return auth
}

func populateAuthObject(auth *Authentication, authorized bool, isAdmin bool) {
    auth.Authenticated = authorized
    auth.IsAdmin = isAdmin
}

func reservationToResponse(res model.Reservation) ReservationResponse {
    errorMsg := res.Error
    if errorMsg != "" {
		return ReservationResponse {Error: errorMsg}
	} else {
        errorMsg = ""
	}
	return ReservationResponse {
        Id:         res.Id,
        EventID:    res.Event.ID,
		ClientID:	res.Client,
        Size:       res.Size,
        Confirmed:  res.Confirmed,
        Timeslot:   res.Timeslot,
        Expiration: res.Expiration,
        Error:      errorMsg,
		Session:	res.Session,
		Event: Event{
			ID:			res.Event.ID,
			Name:		res.Event.Name,
			DtStart:	res.Event.DtStart,
			DtEnd:		res.Event.DtEnd,
		},
	}
}
