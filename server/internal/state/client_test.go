package state

import (
    "testing"
    "github.com/JValtteri/qure/server/internal/crypt"
)


func TestCreateClient(t *testing.T) {
    email := "example@example.com"
    role := "test"
    _, err := NewClient(role, email, crypt.Key("asdf"), false)
    if err != nil {
        t.Errorf("Expected: %v, Got: %v\n", nil, err)
    }
    client, found := GetClientByEmail(email)
    if !found {
        t.Errorf("Expected: %v, Got: %v\n", "found", found)
    }
    email2 := client.GetEmail()
    if email2 != email {
        t.Errorf("Expected: %v, Got: %v\n", email, email2)
    }
    hash := client.GetPasswordHash()
    if len(hash) < 10 {
        t.Errorf("Expected: %v, Got: %v\n", "len(hash) > 10", len(hash))
    }
    role2 := client.GetRole()
    if role2 != role {
        t.Errorf("Expected: %v, Got: %v\n", role, role2)
    }
    admin := client.IsAdmin()
    if admin {
        t.Errorf("Expected: %v, Got: %v\n", "not", admin)
    }

}

func TestCreateTempClient(t *testing.T) {
    email := "temp@example.com"
    _, err := NewClient("test", email, crypt.Key("asdf"), true)
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
    id := client.Id
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
