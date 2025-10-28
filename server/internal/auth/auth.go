package auth

import (
    "fmt"
    "github.com/JValtteri/qure/server/internal/state"
    "github.com/JValtteri/qure/server/internal/crypt"
)


type Authentication struct {
    Authenticated bool
    IsAdmin       bool
    SessionKey    crypt.Key
    Error         string
}

func Login(user string, password string, ip state.IP) (Authentication) {
    client, found := state.GetClientByEmail(user)
    if !found {
        return Authentication{}
    }
    auth := checkPasswordAuthentication(client, password, ip)
    return auth
}

func ReservationLogin(password string, ip state.IP) Authentication {
    client, found := state.GetClientByID(state.ID(password))
    if !found {
        return Authentication{}
    }
    auth := checkPasswordAuthentication(client, password, ip)
    return auth
}

func AuthenticateSession(sessionKey crypt.Key, ip state.IP) Authentication {
    auth := initAuthenticationResponse()
    client, err := state.ResumeSession(sessionKey, ip)
    if err != nil {
        return auth
    }
    populateAuthObject(&auth, true, client.IsAdmin())
    auth.SessionKey = sessionKey
    return auth
}

func Register(email string, password string, ip state.IP) (crypt.Key, string) {
    role := "guest"
    temp := false
    client, err := state.NewClient(role, email, crypt.Key(password), temp)
    if err != nil {
        return crypt.Key(""), fmt.Sprintf("%v", err)
    }
    key, err := client.AddSession(role, email, temp, ip)
    if err != nil {
        return crypt.Key(""), fmt.Sprintf("%v", err)
    }
    return key, ""
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

func checkPasswordAuthentication(client *state.Client, password string, ip state.IP) Authentication {
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
