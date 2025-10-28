package auth

import (
    "testing"
    "github.com/JValtteri/qure/server/internal/state"
    "github.com/JValtteri/qure/server/internal/crypt"
)


func TestNotLogin(t *testing.T) {
    user := "example@example"
    pass := "asdfgh"
    ip := state.IP("0.0.0.0")
    expected := false
    got := Login(user, pass, ip)
    if got.Authenticated {
        t.Errorf("Expected: %v, Got: %v\n", expected, got.Authenticated)
    }
}

func TestNotReservationLogin(t *testing.T) {
    pass := "asdfgh"
    expected := false
    ip := state.IP("0.0.0.0")
    got := ReservationLogin(pass, ip)
    if got.Authenticated {
        t.Errorf("Expected: %v, Got: %v\n", expected, got.Authenticated)
    }
}

func TestNotAuthenticateSession(t *testing.T) {
    sessionkey := crypt.Key("123456")
    ip := state.IP("0.0.0.0")
    expected := false
    got := AuthenticateSession(sessionkey, ip)
    if got.Authenticated {
        t.Errorf("Expected: %v, Got: %v\n", expected, got.Authenticated)
    }
}

func TestRegisterAndAuthenticate(t *testing.T) {
    user := "example@example"
    pass := "asdfgh"
    ip := state.IP("0.0.0.0")
    expected := ""
    key, errMsg := Register(user, pass, ip)
    if errMsg != "" {
        t.Errorf("Expected: %v, Got: %v\n", expected, errMsg)
    }
    auth := AuthenticateSession(key, ip)
    if !auth.Authenticated {
        t.Errorf("Expected: %v, Got: %v\n", "Authenticated", auth.Error)
    }
}

func TestDuplicateRegister(t *testing.T) {
    user := "example@example"
    pass := "asdfgh"
    ip := state.IP("0.0.0.0")
    expected := "Some error"
    _, _ = Register(user, pass, ip)
    _, got := Register(user, pass, ip)
    if got == "" {
        t.Errorf("Expected: %v, Got: %v\n", expected, got)
    }
}

func TestLogin(t *testing.T) {
    user := "login@example"
    pass := "asdfgh"
    ip := state.IP("0.0.0.0")
    expected := ""
    _, errMsg := Register(user, pass, ip)
    if errMsg != "" {
        t.Errorf("Expected: %v, Got: %v\n", expected, errMsg)
    }
    auth := Login(user, pass, ip)
    if !auth.Authenticated {
        t.Errorf("Expected: %v, Got: %v\n", "Authenticated", auth.Error)
    }
    if auth.IsAdmin {
        t.Errorf("Expected: %v, Got: %v\n", "guest", "admin")
    }
    auth = Login(user, "wrong", ip)
    if auth.Authenticated {
        t.Errorf("Expected: %v, Got: %v\n", "No Auth", auth.Authenticated)
    }
    auth = Login(user, pass, state.IP("wrong"))
    if !auth.Authenticated {
        t.Errorf("Expected: %v, Got: %v\n", "Auth", auth.Authenticated)
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
