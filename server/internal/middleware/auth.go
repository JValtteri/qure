package middleware

import (
	"fmt"

	"github.com/JValtteri/qure/server/internal/crypt"
	"github.com/JValtteri/qure/server/internal/state"
)


func Login(rq LoginRequest) (Authentication) {
    client, found := state.GetClientByEmail(rq.User)
    if !found {
        return Authentication{}
    }
    auth := checkPasswordAuthentication(client, rq.Password, rq.Ip)
    return auth
}

func ReservationLogin(rq EventLogin) Authentication {
    client, found := state.GetClientByID(state.ID(rq.EventID))
    if !found {
        return Authentication{}
    }
    auth := checkPasswordAuthentication(client, rq.EventID, rq.Ip)
    return auth
}

func AuthenticateSession(rq AuthenticateRequest) Authentication {
    auth := initAuthenticationResponse()
    client, err := state.ResumeSession(rq.SessionKey, rq.Ip)
    if err != nil {
        return auth
    }
    populateAuthObject(&auth, true, client.IsAdmin())
    auth.SessionKey = rq.SessionKey
    return auth
}

func Register(rq RegisterRequest) RegistrationResponse {
    role := "guest"
    temp := false
    client, err := state.NewClient(role, rq.User, rq.Password, temp)
    if err != nil {
        return RegistrationResponse{Error: fmt.Sprintf("%v", err)}
    }
    key, err := client.AddSession(role, rq.User, temp, rq.Ip)
    if err != nil {
        return RegistrationResponse{Error: fmt.Sprintf("%v", err)}
    }
    return RegistrationResponse{key, ""}
}

func MakeReservation(rq ReserveRequest) state.Reservation {
    return state.MakeReservation(rq.SessionKey, rq.Email, rq.Ip, rq.Size, rq.EventId, rq.Timeslot)
}


func initAuthenticationResponse() Authentication {
    auth := Authentication{
        Authenticated: false,
        IsAdmin:       false,
        SessionKey:    crypt.Key(""),
        Error:         "",
    }
    return auth
}

func checkPasswordAuthentication(client *state.Client, password crypt.Key, ip state.IP) Authentication {
    auth := initAuthenticationResponse()
    authorized := crypt.CompareToHash(password, client.GetPasswordHash())
    if !authorized {
        return auth
    }
    populateAuthObject(&auth, authorized, client.IsAdmin())
    key, err := client.AddSession(client.GetRole(), client.GetEmail(), false, ip)
    if err != nil {
        auth.Error = fmt.Sprintf("%v", err)
        return auth
    }
    auth.SessionKey = key
    return auth
}

func populateAuthObject(auth *Authentication, authorized bool, isAdmin bool) {
    auth.Authenticated = authorized
    auth.IsAdmin = isAdmin
}
