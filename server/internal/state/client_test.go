package state

import (
    "testing"
    "github.com/JValtteri/qure/server/internal/crypt"
)


func TestCreateClient(t *testing.T) {
    _, err := NewClient("test", "example@example.com", crypt.Key("asdf"), false)
    if err != nil {
        t.Errorf("Expected: %v, Got: %v\n", nil, err)
    }
}

func TestCreateTempClient(t *testing.T) {
    _, err := NewClient("test", "temp@example.com", crypt.Key("asdf"), true)
    if err != nil {
        t.Errorf("Expected: %v, Got: %v\n", nil, err)
    }
}

func TestCreateDuplicateEmailClient(t *testing.T) {
    expect := "error"
    _, _ = NewClient("test", "example@example.com", crypt.Key("asdf"), false)
    _, err := NewClient("asd", "example@example.com", crypt.Key("asdf"), false)
    if err == nil {
        t.Errorf("Expected: %v, Got: %v\n", expect, err)
    }
}

func TestRemoveClient(t *testing.T) {
    email := "remove@this.com"
    _, _ = NewClient("test", email, crypt.Key("asdf"), false)
    client := clients.byEmail[email]
    if client.email != email {
        t.Errorf("Test error: Created client corrupt")
    }
    id := client.id
    if len(id) < 16 {
        t.Errorf("Test error: Created client corrupt")
    }
    RemoveClient(client)
    _, found := clients.byEmail[email]
    expect := false
    if found {
        t.Errorf("Expected: %v, Got: %v\n", expect, found)
    }
    _, found = clients.byID[id]
    if found {
        t.Errorf("Expected: %v, Got: %v\n", expect, found)
    }
}
