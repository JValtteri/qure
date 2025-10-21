package state

import (
    "testing"
    "log"
    "io"
)

func TestCreateUniqueID(t *testing.T) {
    const expect int = 16
    var got Key
    var err error
    got, err = createUniqueID(expect, clients.byID)
    if err != nil {
        t.Errorf("Expected: %v, Got: %v\n", "ok", err)
    }
    if len(got) != expect {
        t.Errorf("Expected: %v, Got: %v\n", expect, len(got))
    }
}

func TestCreateConflictingID(t *testing.T) {
    var got Key
    var err error
    i := 0
    // Create lots of short ID's (1 char) to force a collision
    for i < 1000000 {
        i++
        var client Client
        got, err = createUniqueID(1, clients.byID)
        clients.byID[ID(got)] = &client
        if err != nil {
            t.Logf("Tries before conflict %v", i)
            return
        }
    }
    t.Errorf("Expected: %v, Got: %v\n", "Duplicate Error", err)
}

func TestHumanReadableId(t *testing.T) {
    const expect int = 16 // 15
    var got Key
    var err error
    got, err = createHumanReadableId(expect)
    if err != nil {
        t.Errorf("Expected: %v, Got: %v\n", "ok", err)
    }
    if len(got) != expect {
        t.Errorf("Expected: %v, Got: %v\n", expect, len(got))
    }
}

func TestHumanReadableConflictId(t *testing.T) {
    var got Key
    var err error
    i := 0
    // Create lots of short ID's (1 char) to force a collision
    for i < 1000000 {
        i++
        var client Client
        got, err = createHumanReadableId(1)
        clients.bySession[got] = &client
        if err != nil {
            t.Logf("Tries before conflict %v", i)
            return
        }
    }
    t.Errorf("Expected: %v, Got: %v\n", "Duplicate Error", err)
}

func TestCreateClient(t *testing.T) {
    _, err := NewClient("test", "example@example.com", 0, Key("000"))
    if err != nil {
        t.Errorf("Expected: %v, Got: %v\n", nil, err)
    }
}

func TestCreateDuplicateEmailClient(t *testing.T) {
    log.SetOutput(io.Discard)
    expect := "error"
    _, _ = NewClient("test", "example@example.com", 0, Key("000"))
    _, err := NewClient("asd", "example@example.com", 0, Key("123"))
    if err == nil {
        t.Errorf("Expected: %v, Got: %v\n", expect, err)
    }
}

func TestRemoveClient(t *testing.T) {
    email := "remove@this.com"
    _, _ = NewClient("test", email, 0, Key("999"))
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
