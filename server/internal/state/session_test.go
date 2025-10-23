package state

import (
    "os"
    "log"
    "testing"
    "github.com/JValtteri/qure/server/internal/crypt"
)

func resetClients() {
    clients = Clients{
        bySession:  make(map[crypt.Key]*Client),
        byEmail:    make(map[string]*Client),
    }
}

func TestAddTempSession(t *testing.T) {
    log.SetOutput(os.Stdout)
    role := "test"
    email := "temp@example.com"
    ip := IP("0.0.0.0")
    temp := true
    expect := 16
    got, err := AddSession(role, email, temp, ip)
    if err != nil {
        t.Errorf("Expected: %v, Got: %v\n", nil, err)
    }
    if len(got) != expect {
        t.Errorf("Expected: %v, Got: %v\n", expect, got)
    }
}

func TestAddSessions(t *testing.T) {
    log.SetOutput(os.Stdout)
    role := "test"
    email := "session@example.com"
    ip := IP("0.0.0.0")
    temp := false
    expect := 16
    got, err := AddSession(role, email, temp, ip)
    if err != nil {
        t.Errorf("Expected: %v, Got: %v\n", nil, err)
    }
    if len(got) != expect {
        t.Errorf("Expected: %v, Got: %v\n", expect, got)
    }

    _, found := clients.byEmail[email]
    if !found {
        t.Errorf("Expected: %v, Got: %v\n", "found", found)
    }
    _, found = clients.bySession[got]
    if !found {
        t.Errorf("Expected: %v, Got: %v\n", "found", found)
    }

    got, err = AddSession(role, email, temp, ip)
    if err != nil {
        t.Errorf("Expected: %v, Got: %v\n", nil, err)
    }
    if len(got) != expect {
        t.Errorf("Expected: %v, Got: %v\n", expect, got)
    }
}

func TestRemovingNonexistentSession(t *testing.T) {
    log.SetOutput(os.Stdout)

    err := removeSession("asd")
    if err == nil {
        t.Errorf("Expected: '%v', Got: '%v'\n", "error", err)
    }
}

func TestResumeSession(t *testing.T) {
    role := "test"
    email := "session@example.com"
    ip := IP("0.0.0.0")
    temp := false
    key, err := AddSession(role, email, temp, ip)
    if err != nil {
        t.Errorf("Expected: %v, Got: %v\n", nil, err)
    }
    err = ResumeSession(key, ip)
    if err != nil {
        t.Errorf("Expected: '%v', Got: '%v'\n", nil, err)
    }
}

func TestResumeSessionWithChangedIp(t *testing.T) {
    role := "test"
    email := "session@example.com"
    ip0 := IP("0.0.0.0")
    ip1 := IP("0.0.0.1")
    temp := false
    key, err := AddSession(role, email, temp, ip0)
    if err != nil {
        t.Errorf("Expected: %v, Got: %v\n", nil, err)
    }
    err = ResumeSession(key, ip1)
    if err == nil {
        t.Errorf("Expected: '%v', Got: '%v'\n", "error", err)
    }
}

func TestResumeSessionWithWrongKey(t *testing.T) {
    role := "test"
    email := "session@example.com"
    ip := IP("0.0.0.0")
    temp := false
    _, err := AddSession(role, email, temp, ip)
    if err != nil {
        t.Errorf("Expected: %v, Got: %v\n", nil, err)
    }
    err = ResumeSession("wrong key", ip)
    if err == nil {
        t.Errorf("Expected: '%v', Got: '%v'\n", "error", err)
    }
}

func TestCullExpired(t *testing.T) {
    MAX_SESSION_AGE = 0
    role := "test"
    email := "session@example.com"
    ip := IP("0.0.0.0")
    temp := false
    key, _ := AddSession(role, email, temp, ip)

    client := clients.bySession[key]
    err := cullExpired(&client.sessions)
    if err != nil {
        t.Errorf("Expected: '%v', Got: '%v'\n", nil, err)
    }
    _, found := clients.bySession[key]
    expect := false
    if found != expect {
        t.Errorf("Expected: %v, Got: %v\n", expect, found)
    }
    client = clients.byEmail[email]
    _, found = client.sessions[key]
    if found != expect {
        t.Errorf("Expected: %v, Got: %v\n", expect, found)
    }
}

func TestCullExpiredCompletely(t *testing.T) {
    MAX_SESSION_AGE = 0
    role := "test"
    email := "expired@example.com"
    ip := IP("0.0.0.0")
    temp := false
    key, _ := AddSession(role, email, temp, ip)

    client := clients.bySession[key]
    err := cullExpired(&client.sessions)
    if err != nil {
        t.Errorf("Expected: '%v', Got: '%v'\n", nil, err)
    }
    _, found := clients.bySession[key]
    expect := false
    if found != expect {
        t.Errorf("Expected: %v, Got: %v\n", expect, found)
    }
    _, found = clients.byEmail[email]
    if found != expect {
        t.Errorf("Expected: %v, Got: %v\n", expect, found)
    }
}
