package middleware

import (
    "testing"
    "github.com/JValtteri/qure/server/internal/state"
    "github.com/JValtteri/qure/server/internal/crypt"
)


func TestRegistrationLogin(t *testing.T) {
    loginRequest := LoginRequest{
        User: "first@example",
        Password: crypt.Key("asdfghjk"),
        Ip: state.IP("0.0.0.0"),
    }
    expected := false
    got := Login(loginRequest)
    if got.Authenticated {
        t.Errorf("Expected: %v, Got: %v\n", expected, got.Authenticated)
    }

    registerRequest := RegisterRequest{
        User: "second@example",
        Password: crypt.Key("asdfghjk"),
        Ip: state.IP("0.0.0.0"),
    }
    reg := Register(registerRequest)
    if reg.Error != "" {
        t.Errorf("Expected: %v, Got: %v\n", "''", reg.Error)
    }

    got = Login(loginRequest)
    if got.Authenticated {
        t.Errorf("Expected: %v, Got: %v\n", expected, got.Authenticated)
    }
}

func TestNotLogin(t *testing.T) {
    rq := LoginRequest{
        User: "example@example",
        Password: crypt.Key("asdfghjk"),
        Ip: state.IP("0.0.0.0"),
    }
    expected := false
    got := Login(rq)
    if got.Authenticated {
        t.Errorf("Expected: %v, Got: %v\n", expected, got.Authenticated)
    }
    if got.IsAdmin {
        t.Errorf("Expected: %v, Got: %v\n", expected, got.IsAdmin)
    }
    if got.SessionKey != crypt.Key("") {
        t.Errorf("Expected: %v, Got: %v\n", "''", got.SessionKey)
    }
    if got.Error != "" {
        t.Errorf("Expected: %v, Got: %v\n", "''", got.Error)
    }
}

func TestNotReservationLogin(t *testing.T) {
    rq := EventLogin{
        EventID: crypt.Key("asdfgh"),
        Ip: state.IP("0.0.0.0"),
    }
    expected := false
    got := ReservationLogin(rq)
    if got.Authenticated {
        t.Errorf("Expected: %v, Got: %v\n", expected, got.Authenticated)
    }
}

func TestReservationLogin(t *testing.T) {
    expected := true
    ip := state.IP("0.0.0.0")
    event := state.EventFromJson(state.EventJson)
    eventID, err := state.CreateEvent(event)
    if err != nil {
        t.Fatalf("Unexpected error in creating event: %v", err)
    }
    reserveRequest := ReserveRequest{
        SessionKey: crypt.Key(""),
        User: "reserve@example",
        Ip: ip,
        Size: 1,
        EventID: eventID,
        Timeslot: 1100,
    }
    res := MakeReservation(reserveRequest)
    eventLogin := EventLogin{
        EventID: crypt.Key(res.ClientID),
        Ip: ip,
    }
    got := ReservationLogin(eventLogin)
    if !got.Authenticated {
        t.Errorf("Expected: %v, Got: %v\n", expected, got.Authenticated)
    }
}

func TestNotAuthenticateSession(t *testing.T) {
    sessionkey := crypt.Key("123456")
    ip := state.IP("0.0.0.0")
    expected := false
    got := AuthenticateSession(AuthenticateRequest{sessionkey, ip})
    if got.Authenticated {
        t.Errorf("Expected: %v, Got: %v\n", expected, got.Authenticated)
    }
}

func TestRegisterAndAuthenticate(t *testing.T) {
    user := "example@example"
    pass := crypt.Key("asdfghjk")
    ip := state.IP("0.0.0.0")
    expected := ""
    got := Register(RegisterRequest{user, pass, ip})
    if got.Error != "" {
        t.Errorf("Expected: %v, Got: %v\n", expected, got.Error)
    }
    auth := AuthenticateSession(AuthenticateRequest{got.SessionKey, ip})
    if !auth.Authenticated {
        t.Errorf("Expected: %v, Got: %v\n", "Authenticated", auth.Error)
    }
}

func TestDuplicateRegister(t *testing.T) {
    user := "example@example"
    pass := crypt.Key("asdfghjk")
    ip := state.IP("0.0.0.0")
    expected := "Some error"
    _ = Register(RegisterRequest{user, pass, ip})
    got := Register(RegisterRequest{user, pass, ip})
    if got.Error == "" {
        t.Errorf("Expected: %v, Got: %v\n", expected, got.Error)
    }
}

func TestShortRegisters(t *testing.T) {
    user := "long@example"
    pass := crypt.Key("asdfghjk")
    ip := state.IP("0.0.0.0")
    expected := "Some error"
    got := Register(RegisterRequest{"1", pass, ip})
    if got.Error == "" {
        t.Errorf("Expected: %v, Got: %v\n", expected, got.Error)
    }
    got = Register(RegisterRequest{user, crypt.Key("1"), ip})
    if got.Error == "" {
        t.Errorf("Expected: %v, Got: %v\n", expected, got.Error)
    }
}

func TestLogin(t *testing.T) {
    user := "login@example"
    pass := crypt.Key("asdfghjk")
    ip := state.IP("0.0.0.0")
    expected := ""
    got := Register(RegisterRequest{user, pass, ip})
    if got.Error != "" {
        t.Errorf("Expected: %v, Got: %v\n", expected, got.Error)
    }
    auth := Login(LoginRequest{user, pass, ip})
    if !auth.Authenticated {
        t.Errorf("Expected: %v, Got: %v\n", "Authenticated", auth.Error)
    }
    if auth.IsAdmin {
        t.Errorf("Expected: %v, Got: %v\n", "guest", "admin")
    }
    auth = Login(LoginRequest{user, "wrong", ip})
    if auth.Authenticated {
        t.Errorf("Expected: %v, Got: %v\n", "No Auth", auth.Authenticated)
    }
    auth = Login(LoginRequest{user, pass, state.IP("wrong")})
    if !auth.Authenticated {
        t.Errorf("Expected: %v, Got: %v\n", "Auth", auth.Authenticated)
    }
}
