package middleware

import (
	"testing"

	"github.com/JValtteri/qure/server/internal/crypt"
	"github.com/JValtteri/qure/server/internal/state"
	"github.com/JValtteri/qure/server/internal/testjson"
	c "github.com/JValtteri/qure/server/internal/config"
)


func TestRegistrationLogin(t *testing.T) {
    loginRequest := LoginRequest{
        User: "first@example",
        Password: crypt.Key("asdfghjk"),
    }
    expected := false
    got := Login(loginRequest)
    if got.Authenticated {
        t.Errorf("Expected: %v, Got: %v\n", expected, got.Authenticated)
    }

    registerRequest := RegisterRequest{
        User: "second@example",
        Password: crypt.Key("asdfghjk"),
        HashPrint: crypt.GenerateHash("0.0.0.0"),
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
	rq := LoginRequest{
		User: "not",
		Password: "foo",
		HashPrint: "none",
    }
    expected := false
	got := Login(rq)
    if got.Authenticated {
        t.Errorf("Expected: %v, Got: %v\n", expected, got.Authenticated)
    }
}

func TestReservationLogin(t *testing.T) {
    expected := true
    fingerprint := "0.0.0.0"
    event := state.EventFromJson(testjson.EventJson)
    eventID, err := state.CreateEvent(event)
    if err != nil {
        t.Fatalf("Unexpected error in creating event: %v", err)
    }
    reserveRequest := ReserveRequest{
        SessionKey: crypt.Key(""),
        User: "reserve@example",
        Fingerprint: fingerprint,
        Size: 1,
        EventID: eventID,
        Timeslot: 1100,
    }
    res := MakeReservation(reserveRequest)
	eventLogin := LoginRequest{
		User: reserveRequest.User,
		Password: crypt.Key(res.Id),
		HashPrint: crypt.GenerateHash(reserveRequest.Fingerprint),
	}
	got := Login(eventLogin)
    if !got.Authenticated {
        t.Errorf("Expected: %v, Got: %v\n", expected, got.Authenticated)
    }
}

func TestNotAuthenticateSession(t *testing.T) {
    sessionkey := crypt.Key("123456")
    fingerprint := "0.0.0.0"
    expected := false
    got := AuthenticateSession(AuthenticateRequest{sessionkey, fingerprint})
    if got.Authenticated {
        t.Errorf("Expected: %v, Got: %v\n", expected, got.Authenticated)
    }
}

func TestRegisterAuthenticateAndLogout(t *testing.T) {
    user := "example@example"
    pass := crypt.Key("asdfghjk")
    fingerprint := "0.0.0.0"
    expected := ""
    got := Register(RegisterRequest{user, pass, crypt.GenerateHash(fingerprint)})
    if got.Error != "" {
        t.Errorf("Expected: %v, Got: %v\n", expected, got.Error)
    }
    auth := AuthenticateSession(AuthenticateRequest{got.SessionKey, fingerprint})
    if !auth.Authenticated {
        t.Errorf("Expected: %v, Got: %v\n", "Authenticated", auth.Error)
    }
    _ = Logout(AuthenticateRequest{got.SessionKey, fingerprint})
    auth = AuthenticateSession(AuthenticateRequest{got.SessionKey, fingerprint})
    if auth.Authenticated {
        t.Errorf("Expected: %v, Got: %v\n", "Not Authenticated", auth.Authenticated)
    }
}

func TestDuplicateRegister(t *testing.T) {
    user := "example@example"
    pass := crypt.Key("asdfghjk")
    fingerprint := crypt.Hash("0.0.0.0")
    expected := "Some error"
    _ = Register(RegisterRequest{user, pass, fingerprint})
    got := Register(RegisterRequest{user, pass, fingerprint})
    if got.Error == "" {
        t.Errorf("Expected: %v, Got: %v\n", expected, got.Error)
    }
}

func TestShortRegisters(t *testing.T) {
    user := "long@example"
    pass := crypt.Key("asdfghjk")
    fingerprint := crypt.Hash("0.0.0.0")
    expected := "Some error"
    got := Register(RegisterRequest{"1", pass, fingerprint})
    if got.Error == "" {
        t.Errorf("Expected: %v, Got: %v\n", expected, got.Error)
    }
    got = Register(RegisterRequest{user, crypt.Key("1"), fingerprint})
    if got.Error == "" {
        t.Errorf("Expected: %v, Got: %v\n", expected, got.Error)
    }
}

func TestCancelReservation(t *testing.T) {
	state.ResetEvents()
	var expected = 1
	var expectedNo = 0
	var fingerprint = "0.0.0.0"
	event := state.EventFromJson(testjson.EventJson)
	event.ID += crypt.ID("cancelTestEvent126123")
	eventID, err := state.CreateEvent(event)
	if err != nil {
		t.Fatalf("Expected: %v, Got: %v\n", "no error", err)
	}
	reserveRequest := ReserveRequest{
		SessionKey: crypt.Key(""),
		User: "cancel@example",
		Fingerprint: fingerprint,
		Size: 1,
		EventID: eventID,
		Timeslot: 1100,
	}
	res := MakeReservation(reserveRequest)
	if res.Error != "" {
		t.Fatalf("Expected: %v, Got: %v\n", "no error", res.Error)
	}
	if res.Confirmed != expected {
		t.Errorf("Expected: %v, Got: %v\n", expected, res.Confirmed)
	}
	reserveRequest.Id = res.Id
	reserveRequest.SessionKey = res.Session
	res = CancelReservation(reserveRequest)
	if res.Error != "" {
		t.Fatalf("Expected: %v, Got: %v\n", "no error", res.Error)
	}
	if res.Confirmed != expectedNo {
		t.Errorf("Expected: %v, Got: %v\n", expectedNo, res.Confirmed)
	}
}

func TestLogin(t *testing.T) {
    user := "login@example"
    pass := crypt.Key("asdfghjk")
    fingerprint := "0.0.0.0"
    expected := ""
    got := Register(RegisterRequest{user, pass, crypt.GenerateHash(fingerprint)})
    if got.Error != "" {
        t.Errorf("Expected: %v, Got: %v\n", expected, got.Error)
    }
    auth := Login(LoginRequest{user, pass, crypt.GenerateHash(fingerprint)})
    if !auth.Authenticated {
        t.Errorf("Expected: %v, Got: %v\n", "Authenticated", auth.Error)
    }
    if auth.IsAdmin {
        t.Errorf("Expected: %v, Got: %v\n", "guest", "admin")
    }
    auth = Login(LoginRequest{user, "wrong", crypt.GenerateHash(fingerprint)})
    if auth.Authenticated {
        t.Errorf("Expected: %v, Got: %v\n", "No Auth", auth.Authenticated)
    }
    auth = Login(LoginRequest{user, pass, crypt.GenerateHash("wrong")})
    if !auth.Authenticated {
        t.Errorf("Expected: %v, Got: %v\n", "Auth", auth.Authenticated)
    }
}

func TestChangePassword(t *testing.T) {
    c.CONFIG.MIN_PASSWORD_LENGTH = 3
	user := "change@example"
	pass := crypt.Key("12345678")
	newPass := crypt.Key("654321")
	fingerprint := "0.0.0.0"
	got := Register(RegisterRequest{user, pass, crypt.GenerateHash(fingerprint)})
	if got.Error != "" {
		t.Fatalf("Client wasn't created: %v", got.Error)
	}
	_, ok := state.GetClientByEmail(user)
	if !ok {
		t.Fatalf("Client wasn't found")
	}
	ChangePassword(PasswordChangeRequest{
		User: user,
		SessionKey: got.SessionKey,
		Fingerprint: fingerprint,
		HashPrint: crypt.GenerateHash(fingerprint),
		Password: pass,
		NewPassword: newPass,
	})
	auth := Login(LoginRequest{user, newPass, crypt.GenerateHash(fingerprint)})
	if !auth.Authenticated {
		t.Errorf("Expected: %v, Got: %v\n", "Authenticated", auth.Error)
	}
}

func TestRemoveUser(t *testing.T) {
	user := "remove@example"
	pass := crypt.Key("12345678")
	fingerprint := "0.0.0.0"
	got := Register(RegisterRequest{user, pass, crypt.GenerateHash(fingerprint)})
	if got.Error != "" {
		t.Fatalf("Client wasn't created: %v", got.Error)
	}
	_, ok := state.GetClientByEmail(user)
	if !ok {
		t.Fatalf("Client wasn't found")
	}
	RemoveUser(RemovalRequest{
		User: user,
		SessionKey: got.SessionKey,
		Fingerprint: fingerprint,
		HashPrint: crypt.GenerateHash(fingerprint),
		Password: pass,
	})
	auth := Login(LoginRequest{user, pass, crypt.GenerateHash(fingerprint)})
	if auth.Authenticated {
		t.Errorf("Expected: %v, Got: %v\n", false, auth.Authenticated)
	}
}
