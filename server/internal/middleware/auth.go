package middleware

import (
	"fmt"
    "log"

	"github.com/JValtteri/qure/server/internal/crypt"
	"github.com/JValtteri/qure/server/internal/state"
)


var MIN_USERNAME_LENGTH int = 4
var MIN_PASSWORD_LENGTH int = 8

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
func ReservationLogin(rq EventLogin) Authentication {
    client, found := state.GetClientByID(state.ID(rq.EventID))
    if !found {
        log.Printf("Reservation '%v' not found\n", rq.EventID)
        return Authentication{}
    }
    auth := checkPasswordAuthentication(client, rq.EventID, rq.HashPrint)
    return auth
}

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
    if len(rq.User) < MIN_USERNAME_LENGTH {
        return RegistrationResponse{Error: fmt.Sprintf("Username length must be at least %v characters", MIN_USERNAME_LENGTH)}
    }
    if len(rq.Password) < MIN_PASSWORD_LENGTH {
        return RegistrationResponse{Error: fmt.Sprintf("Password length must be at least %v characters", MIN_PASSWORD_LENGTH)}
    }
    client, err := state.NewClient(role, rq.User, rq.Password, temp)
    if err != nil {
        return RegistrationResponse{Error: fmt.Sprintf("%v", err)}
    }
    key, err := client.AddSession(role, rq.User, temp, rq.HashPrint)
    if err != nil {
        return RegistrationResponse{Error: fmt.Sprintf("%v", err)}
    }
    return RegistrationResponse{key, ""}
}

func MakeReservation(rq ReserveRequest) Reservation {
    res := state.MakeReservation(rq.SessionKey, rq.User, rq.Fingerprint, rq.HashPrint, rq.Size, rq.EventID, rq.Timeslot)
    return reservationToResponse(res)
}

func checkPasswordAuthentication(client *state.Client, password crypt.Key, hashedFingerprint crypt.Hash) Authentication {
    authorized := crypt.CompareToHash(password, client.GetPasswordHash())
    if !authorized {
        fmt.Println("Password doesn't match")
        log.Printf("Client's '%v' password didn't match\n", client.GetEmail())
        return Authentication{}
    }
    auth := Authentication{}
    populateAuthObject(&auth, authorized, client.IsAdmin())
    key, err := client.AddSession(client.GetRole(), client.GetEmail(), false, hashedFingerprint)
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

func reservationToResponse(res state.Reservation) Reservation {
    errorMsg := res.Error
    if errorMsg == "<nil>" {
        errorMsg = ""
    }
    return Reservation{
        Id:         res.Id,
        EventID:    res.Event.ID,
        ClientID:   res.Client.Id,
        Size:       res.Size,
        Confirmed:  res.Confirmed,
        Timeslot:   res.Timeslot,
        Expiration: res.Expiration,
        Error:      errorMsg,
    }
}
