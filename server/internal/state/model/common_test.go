package model

import (
	"testing"
	"github.com/JValtteri/qure/server/internal/crypt"
)


var testClients Clients = Clients{
	ByID:		make(map[crypt.ID]*Client),
	BySession:	make(map[crypt.Key]*Client),
	ByEmail:	make(map[string]*Client),
}

func TestCreateUniqueID(t *testing.T) {
    const expect int = 16
	var got crypt.ID
    var err error
	got, err = CreateUniqueID(expect, testClients.ByID)
    if err != nil {
        t.Errorf("Expected: %v, Got: %v\n", "ok", err)
    }
    if len(got) != expect {
        t.Errorf("Expected: %v, Got: %v\n", expect, len(got))
    }
}

func TestCreateConflictingID(t *testing.T) {
	var got crypt.ID
    var err error
    i := 0
    // Create lots of short ID's (1 char) to force a collision
    for i < 1000000 {
        i++
		var client Client
		got, err = CreateUniqueID(1, testClients.ByID)
		testClients.ByID[got] = &client
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
	got, err = CreateUniqueHumanReadableID(expect, testClients.ByID)
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
		got, err = CreateUniqueHumanReadableKey(1, testClients.BySession)
		testClients.BySession[got] = &client
        if err != nil {
            t.Logf("Tries before conflict %v", i)
            return
        }
    }
    t.Errorf("Expected: %v, Got: %v\n", "Duplicate Error", err)
}
