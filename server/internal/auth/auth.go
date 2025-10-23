package auth

import (
    //"github.com/JValtteri/qure/server/internal/state"
    "github.com/JValtteri/qure/server/internal/crypt"
)

type Authentication struct {
    Authenticated bool
    IsAdmin       bool
    SessionKey    crypt.ID
}

func AuthenticateLogin(user string, password string) Authentication {
    auth := Authentication{false, false, crypt.ID("")}
    // Check valid session
    // Check admin
    return auth
}

func AuthenticateSession(sessionKey crypt.ID) Authentication {
    auth := Authentication{false, false, crypt.ID("")}
    // Check valid session
    // Check admin
    return auth
}
