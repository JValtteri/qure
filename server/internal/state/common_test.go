package state

import (
    "testing"
    "github.com/JValtteri/qure/server/internal/crypt"
)

func TestCreateUniqueID(t *testing.T) {
    const expect int = 16
    var got ID
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
    var got ID
    var err error
    i := 0
    // Create lots of short ID's (1 char) to force a collision
    for i < 1000000 {
        i++
        var client Client
        got, err = createUniqueID(1, clients.byID)
        clients.byID[got] = &client
        if err != nil {
            t.Logf("Tries before conflict %v", i)
            return
        }
    }
    t.Errorf("Expected: %v, Got: %v\n", "Duplicate Error", err)
}

func TestHumanReadableId(t *testing.T) {
    const expect int = 16 // 15
    var got crypt.ID
    var err error
    got, err = createUniqueHumanReadableID(expect, clients.byID)
    if err != nil {
        t.Errorf("Expected: %v, Got: %v\n", "ok", err)
    }
    if len(got) != expect {
        t.Errorf("Expected: %v, Got: %v\n", expect, len(got))
    }
}

func TestHumanReadableConflictId(t *testing.T) {
    var got crypt.Key
    var err error
    i := 0
    // Create lots of short ID's (1 char) to force a collision
    for i < 1000000 {
        i++
        var client Client
        got, err = createUniqueHumanReadableKey(1, clients.bySession)
        clients.bySession[got] = &client
        if err != nil {
            t.Logf("Tries before conflict %v", i)
            return
        }
    }
    t.Errorf("Expected: %v, Got: %v\n", "Duplicate Error", err)
}
