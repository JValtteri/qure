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
    auth := Authentication{
        Authenticated:  false,
        IsAdmin:        false,
        SessionKey:     crypt.Key(""),
    }
    client, found := state.GetClientByEmail(user)
    if !found {
        return auth
    }
    authorized := crypt.CompareToHash(password, client.GetPasswordHash())
    auth.Authenticated = authorized
    auth.IsAdmin = client.IsAdmin()
    key, err := client.AddSession(client.GetRole(), client.GetEmail(), false, ip)
    if err != nil {
        auth.Error = fmt.Sprintf("%v", err)
        return auth
    }
    auth.SessionKey = key
    return auth
}

func ReservationLogin(password string, ip state.IP) Authentication {
    auth := Authentication{
        Authenticated:  false,
        IsAdmin:        false,
        SessionKey:     crypt.Key(""),
    }
    client, found := state.GetClientByID(state.ID(password))
    if !found {
        return auth
    }
    authorized := crypt.CompareToHash(password, client.GetPasswordHash())
    auth.Authenticated = authorized
    key, err := client.AddSession(client.GetRole(), client.GetEmail(), true, ip)
    if err != nil {
        auth.Error = fmt.Sprintf("%v", err)
        return auth
    }
    auth.SessionKey = key
    auth.IsAdmin = client.IsAdmin()
    return auth
}

func AuthenticateSession(sessionKey crypt.Key, ip state.IP) Authentication {
    auth := Authentication{
        Authenticated:  false,
        IsAdmin:        false,
        SessionKey:     crypt.Key(""),
    }
    client, err := state.ResumeSession(sessionKey, ip)
    if err != nil {
        return auth
    }
    auth.Authenticated = true
    auth.SessionKey = sessionKey
    auth.IsAdmin = client.IsAdmin()
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
    return key, ""
}
