package state

import (
    "os"
    "log"
    "testing"
    "github.com/JValtteri/qure/server/internal/crypt"
    "github.com/JValtteri/qure/server/internal/utils"
)


func TestAddSessions(t *testing.T) {
    ResetClients()
    log.SetOutput(os.Stdout)
    role := "test"
    email := "session@example.com"
    ip := IP("0.0.0.0")
    temp := false
    expect := SESSION_KEY_LENGTH
    client, err := NewClient(role, email, crypt.Key("asdf"), temp)
    if err != nil {
        t.Fatalf("Expected: %v, Got: %v\n", nil, err)
    }
    got, err := client.AddSession(role, email, temp, ip)
    if err != nil {
        t.Errorf("Expected: %v, Got: %v\n", nil, err)
    }
    if len(got) != expect {
        t.Errorf("Expected: %v, Got: %v\n", expect, len(got))
    }

    _, found := GetClientByEmail(email)
    if !found {
        t.Errorf("Expected: %v, Got: %v\n", "found", found)
    }
    _, found = GetClientBySession(got)
    if !found {
        t.Errorf("Expected: %v, Got: %v\n", "found", found)
    }

    got, err = client.AddSession(role, email, temp, ip)
    if err != nil {
        t.Errorf("Expected: %v, Got: %v\n", nil, err)
    }
    if len(got) != expect {
        t.Errorf("Expected: %v, Got: %v\n", expect, len(got))
    }
}

func TestRemovingNonexistentSession(t *testing.T) {
    ResetClients()
    log.SetOutput(os.Stdout)

    err := removeSession("asd")
    if err == nil {
        t.Errorf("Expected: '%v', Got: '%v'\n", "error", err)
    }
}

func TestResumeSession(t *testing.T) {
    ResetClients()
    role := "test"
    email := "resume@example.com"
    ip := IP("0.0.0.0")
    temp := false
    client, err := NewClient(role, email, crypt.Key("asdf"), temp)
    if err != nil {
        t.Fatalf("Expected: %v, Got: %v\n", nil, err)
    }
    key, err := client.AddSession(role, email, temp, ip)
    if err != nil {
        t.Errorf("Expected: %v, Got: %v\n", nil, err)
    }
    _, err = ResumeSession(key, ip)
    if err != nil {
        t.Errorf("Expected: '%v', Got: '%v'\n", nil, err)
    }
}


func TestResumeSessionWithChangedIp(t *testing.T) {
    ResetClients()
    role := "test"
    email := "resume2@example.com"
    ip0 := IP("0.0.0.0")
    ip1 := IP("0.0.0.1")
    temp := false
    client, err := NewClient(role, email, crypt.Key("asdf"), temp)
    if err != nil {
        t.Fatalf("Expected: %v, Got: %v\n", nil, err)
    }
    key, err := client.AddSession(role, email, temp, ip0)
    if err != nil {
        t.Errorf("Expected: %v, Got: %v\n", nil, err)
    }
    _, err = ResumeSession(key, ip1)
    if err == nil {
        t.Errorf("Expected: '%v', Got: '%v'\n", "error", err)
    }
}

func TestResumeSessionWithWrongKey(t *testing.T) {
    ResetClients()
    role := "test"
    email := "resum3@example.com"
    ip := IP("0.0.0.0")
    temp := false
    client, err := NewClient(role, email, crypt.Key("asdf"), temp)
    if err != nil {
        t.Fatalf("Expected: %v, Got: %v\n", nil, err)
    }
    _, err = client.AddSession(role, email, temp, ip)
    if err != nil {
        t.Errorf("Expected: %v, Got: %v\n", nil, err)
    }
    _, err = ResumeSession("wrong key", ip)
    if err == nil {
        t.Errorf("Expected: '%v', Got: '%v'\n", "error", err)
    }
}

func TestCullExpired(t *testing.T) {
    ResetClients()
    MAX_SESSION_AGE = 0
    role := "test"
    email := "cull@example.com"
    ip := IP("0.0.0.0")
    temp := false
    client, err := NewClient(role, email, crypt.Key("asdf"), temp)
    if err != nil {
        t.Fatalf("Expected: %v, Got: %v\n", nil, err)
    }
    addPersistantSession(ip, client)
    addPersistantSession(ip, client)
    key, err := client.AddSession(role, email, temp, ip)
    if err != nil {
        t.Fatalf("Expected: %v, Got: %v\n", nil, err)
    }
    expectSessions := 3
    if len(client.sessions) != expectSessions {
        t.Errorf("Expected: %v, Got: %v\n", expectSessions, len(client.sessions))
    }

    err = cullExpired(&client.sessions)
    if err != nil {
        t.Errorf("Expected: '%v', Got: '%v'\n", nil, err)
    }
    _, found := clients.getClientBySession(key)
    expect := false
    if found != expect {
        t.Errorf("Expected: %v, Got: %v\n", expect, found)
    }
    _, err = ResumeSession(key, ip)
    if err == nil {
        t.Errorf("Expected: '%v', Got: '%v'\n", "error", err)
    }
    expectSessions = 2
    if len(client.sessions) != expectSessions {
        t.Errorf("Expected: %v, Got: %v\n", expectSessions, len(client.sessions))
    }
}

func addPersistantSession(ip IP, client *Client) {
    sessionKey, _ := createUniqueKey(SESSION_KEY_LENGTH, clients.bySession)
    now := utils.EpochNow()
    var session Session = Session{
        key:        sessionKey,
        expiresDt:  now + Epoch(1000),
        ip:         ip,
    }
    clients.withLock(func() {
        client.sessions[sessionKey] = session
        clients.bySession[sessionKey] = client
    })
}

func TestCullExpiredCompletely(t *testing.T) {
    ResetClients()
    MAX_SESSION_AGE = 0
    role := "test"
    email := "expired@example.com"
    ip := IP("0.0.0.0")
    temp := false
    client, err := NewClient(role, email, crypt.Key("asdf"), temp)
    if err != nil {
        t.Fatalf("Expected: %v, Got: %v\n", nil, err)
    }
    _, _ = client.AddSession(role, email, temp, ip)
    key, _ := client.AddSession(role, email, temp, ip)

    client, _ = clients.getClientBySession(key)
    id := client.Id
    err = cullExpired(&client.sessions)
    if err != nil {
        t.Errorf("Expected: '%v', Got: '%v'\n", nil, err)
    }
    _, found := clients.getClientBySession(key)
    expect := false
    if found != expect {
        t.Errorf("Expected: %v, Got: %v\n", expect, found)
    }
    _, found = GetClientByID(id)
    if found != expect {
        t.Errorf("Expected: %v, Got: %v\n", expect, found)
    }
    expectSessions := 0
    if len(client.sessions) != expectSessions {
        t.Errorf("Expected: %v, Got: %v\n", expectSessions, len(client.sessions))
    }
}
