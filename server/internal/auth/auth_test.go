package auth

import (
    "testing"
    //"github.com/JValtteri/qure/server/internal/state"
    "github.com/JValtteri/qure/server/internal/crypt"
)

func TestNotAuthenticateLogin(t *testing.T) {
    user := "example@example"
    pass := "asdfgh"
    expected := false
    got := AuthenticateLogin(user, pass)
    if got.Authenticated {
        t.Errorf("Expected: %v, Got: %v\n", expected, got.Authenticated)
    }
}

func TestNotAuthenticateSession(t *testing.T) {
    sessionkey := crypt.ID("123456")
    expected := false
    got := AuthenticateSession(sessionkey)
    if got.Authenticated {
        t.Errorf("Expected: %v, Got: %v\n", expected, got.Authenticated)
    }
}

/*
func TestAuthenticateSession(t *testing.T) {
    user := "example@example"
    pass := "asdfgh"
    expected := true
    admin := false
    expire := state.Epoch(0)
    expire--
    session := crypt.Key("1234")
    state.NewClient("user", user, expire, session)
    got := AuthenticateLogin(user, pass)
    if !got.Authenticated {
        t.Errorf("Expected: %v, Got: %v\n", expected, got.Authenticated)
    }
    if got.IsAdmin {
        t.Errorf("Expected: %v, Got: %v\n", admin, got.IsAdmin)
    }
}
*/
