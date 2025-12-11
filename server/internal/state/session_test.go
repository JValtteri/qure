package state

import (
	"os"
	"log"
	"testing"
	"github.com/JValtteri/qure/server/internal/crypt"
	"github.com/JValtteri/qure/server/internal/utils"
	c "github.com/JValtteri/qure/server/internal/config"
	"github.com/JValtteri/qure/server/internal/state/model"
)


func TestAddSessions(t *testing.T) {
    ResetClients()
    log.SetOutput(os.Stdout)
    role := "test"
    email := "session@example.com"
    fingerprint := crypt.Hash("0.0.0.0")
    temp := false
	expect := model.SESSION_KEY_LENGTH
    client, err := NewClient(role, email, crypt.Key("asdf"), temp)
    if err != nil {
        t.Fatalf("Expected: %v, Got: %v\n", nil, err)
    }
	got, err := client.AddSession(role, email, temp, fingerprint, &clients)
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

	got, err = client.AddSession(role, email, temp, fingerprint, &clients)
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

    err := RemoveSession("asd")
    if err == nil {
        t.Errorf("Expected: '%v', Got: '%v'\n", "error", err)
    }
}

func TestResumeSession(t *testing.T) {
    ResetClients()
    role := "test"
    email := "resume@example.com"
    fingerprint := "0.0.0.0"
    temp := false
    client, err := NewClient(role, email, crypt.Key("asdf"), temp)
    if err != nil {
        t.Fatalf("Expected: %v, Got: %v\n", nil, err)
    }
	key, err := client.AddSession(role, email, temp, crypt.GenerateHash(fingerprint), &clients)
    if err != nil {
        t.Errorf("Expected: %v, Got: %v\n", nil, err)
    }
    _, err = ResumeSession(key, fingerprint)
    if err != nil {
        t.Errorf("Expected: '%v', Got: '%v'\n", nil, err)
    }
}


func TestResumeSessionWithChangedIp(t *testing.T) {
    ResetClients()
    role := "test"
    email := "resume2@example.com"
    fingerprint0 := "0.0.0.0"
    fingerprint1 := "0.0.0.1"
    temp := false
    client, err := NewClient(role, email, crypt.Key("asdf"), temp)
    if err != nil {
        t.Fatalf("Expected: %v, Got: %v\n", nil, err)
    }
	key, err := client.AddSession(role, email, temp, crypt.GenerateHash(fingerprint0), &clients)
    if err != nil {
        t.Errorf("Expected: %v, Got: %v\n", nil, err)
    }
    _, err = ResumeSession(key, fingerprint1)
    if err == nil {
        t.Errorf("Expected: '%v', Got: '%v'\n", "error", err)
    }
}

func TestResumeSessionWithWrongKey(t *testing.T) {
    ResetClients()
    role := "test"
    email := "resum3@example.com"
    fingerprint := "0.0.0.0"
    temp := false
    client, err := NewClient(role, email, crypt.Key("asdf"), temp)
    if err != nil {
        t.Fatalf("Expected: %v, Got: %v\n", nil, err)
    }
	_, err = client.AddSession(role, email, temp, crypt.GenerateHash(fingerprint), &clients)
    if err != nil {
        t.Errorf("Expected: %v, Got: %v\n", nil, err)
    }
    _, err = ResumeSession("wrong key", fingerprint)
    if err == nil {
        t.Errorf("Expected: '%v', Got: '%v'\n", "error", err)
    }
}

func TestCullExpired(t *testing.T) {
    ResetClients()
	c.CONFIG.MAX_SESSION_AGE = 0
    role := "test"
    email := "cull@example.com"
    fingerprint := "0.0.0.0"
    temp := false
    client, err := NewClient(role, email, crypt.Key("asdf"), temp)
    if err != nil {
        t.Fatalf("Expected: %v, Got: %v\n", nil, err)
    }
    addPersistantSession(crypt.GenerateHash(fingerprint), client)
    addPersistantSession(crypt.GenerateHash(fingerprint), client)
	key, err := client.AddSession(role, email, temp, crypt.GenerateHash(fingerprint), &clients)
    if err != nil {
        t.Fatalf("Expected: %v, Got: %v\n", nil, err)
    }
    expectSessions := 3
	if len(client.Sessions) != expectSessions {
		t.Errorf("Expected: %v, Got: %v\n", expectSessions, len(client.Sessions))
	}

	err = cullExpired(&client.Sessions)
    if err != nil {
        t.Errorf("Expected: '%v', Got: '%v'\n", nil, err)
    }
	_, found := clients.GetClientBySession(key)
    expect := false
    if found != expect {
        t.Errorf("Expected: %v, Got: %v\n", expect, found)
    }
    _, err = ResumeSession(key, fingerprint)
    if err == nil {
        t.Errorf("Expected: '%v', Got: '%v'\n", "error", err)
    }
    expectSessions = 2
	if len(client.Sessions) != expectSessions {
		t.Errorf("Expected: %v, Got: %v\n", expectSessions, len(client.Sessions))
	}
}

func addPersistantSession(fingerprint crypt.Hash, client *model.Client) {
	sessionKey, _ := model.CreateUniqueKey(model.SESSION_KEY_LENGTH, clients.BySession)
	now := utils.EpochNow()
	var session model.Session = model.Session{
		Key:        sessionKey,
		ExpiresDt:  now + utils.Epoch(1000),
		Fingerprint:         fingerprint,
	}
	clients.Lock()
	defer clients.Unlock()
	client.Sessions[sessionKey] = session
	clients.BySession[sessionKey] = client
}
