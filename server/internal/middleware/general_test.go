package middleware

import (
	"testing"

	"github.com/JValtteri/qure/server/internal/crypt"
	"github.com/JValtteri/qure/server/internal/state"
)

func TestEventLifesycle(t *testing.T) {
	state.ResetEvents()
	state.ResetClients()
	isAdmin := false
	// Make events
	id := MakeEvent(state.EventJson, true)	// Must be admin to make event
	_   = MakeEvent([]byte("{}"), true)		// Decoy event
	if id == crypt.ID("") {
		t.Fatal("Critical error making event")
	}
	// Make user
	email	 := "new@email"
	pass	 := "asdfgh"
	ip		 := state.IP("0.0.0.0")
	size	 := 1
	key, msg := Register(email, pass, ip)
	if msg != "" {
		t.Fatalf("Expected: %v, Got: %v\n", nil, msg)
	}
	ress	 := GetUserReservatoions(key)
	if len(ress) != 0 {
		t.Errorf("Expected: %v, Got: %v\n", 0, len(ress))
	}
	// Check Events
	events := GetEvents(isAdmin)
	if len(events) != 2 {
		t.Errorf("Expected: %v, Got: %v\n", 2, len(events))
	}
	if events[0].DtEnd != 1735687830 {
		t.Errorf("Expected: %v, Got: %v\n", 1735687830, events[0].DtEnd)
	}
	event := GetEvent(id, isAdmin)
	if event.DtStart != 1735675270 {
		t.Errorf("Expected: %v, Got: %v\n", 1735675270, event.DtStart)
	}
	// Make reservation
	res      := MakeReservation(key, email, ip, size, id, state.Epoch(1100))
	if res.Error != "<nil>" {
		t.Fatalf("Expected: %v, Got: %v\n", nil, res.Error)
	}
	ress	 = GetUserReservatoions(key)
	if ress[0].Event.DtEnd != 1735687830 {
		t.Errorf("Expected: %v, Got: %v\n", 1735687830, ress[0].Event.DtEnd)
	}
	if len(ress) != 1 {
		t.Errorf("Expected: %v, Got: %v\n", 1, len(ress))
	}
	if res.Id != ress[0].Id {
		t.Errorf("Expected: %v, Got: %v\n", res.Id, ress[0].Id)
	}
}

func TestNotAdminMakeEvent(t *testing.T) {
	state.ResetEvents()
	state.ResetClients()
	isAdmin := false
	id := MakeEvent(state.EventJson, isAdmin)
	if id != state.ID("") {
		t.Errorf("Expected: %v, Got: %v\n", "''", id)
	}
}

