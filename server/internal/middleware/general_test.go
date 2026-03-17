package middleware

import (
	"testing"

	"github.com/JValtteri/qure/server/internal/utils"
	"github.com/JValtteri/qure/server/internal/crypt"
	"github.com/JValtteri/qure/server/internal/state"
	"github.com/JValtteri/qure/server/internal/state/model"
	"github.com/JValtteri/qure/server/internal/testjson"
)

func TestGetInvalidEvent(t *testing.T) {
	state.ResetEvents()
	state.ResetClients()
	event := GetEvent(EventRequest{EventID: crypt.ID("no-id")})
	if len(event.ID) > 1 {
		t.Errorf("Expected: len(event.ID) < 1, Got: %v\n", len(event.ID))
	}
}

func TestEventLifesycle(t *testing.T) {
	state.ResetEvents()
	state.ResetClients()
	// Make admin
	var adminName = "admin@test"
	var adminPassword = crypt.Key("adminpass")
	var adminprint = "adminadminadmin"

	_, err := state.NewClient("admin", adminName, adminPassword, false)
	if err != nil {
		t.Fatalf("Error generating test-admin account:\n%v", err)
	}
	var adminAuth = Login(LoginRequest{
		User: adminName,
		Password: crypt.Key(adminPassword),
		HashPrint: crypt.GenerateHash(adminprint),
	})
	// Make events
	newEvent := state.EventFromJson(testjson.EventJson)
	resp := MakeEvent(EventManipulationRequest{adminAuth.SessionKey, adminprint, newEvent})
	emptyEvent := state.EventFromJson([]byte("{}"))
	_   = MakeEvent(EventManipulationRequest{adminAuth.SessionKey, adminprint, emptyEvent})		// Decoy event
	if resp.EventID == crypt.ID("") {
		t.Fatal("Critical error making event")
	}
	// Modify Event
	modEvent, err := state.GetEvent(resp.EventID, true)
	if resp.Error != "" {
		t.Fatalf("Critical fetching event %v", err)
	}
	modEvent.Name = "Updated Event"
	resp = EditEvent(EventManipulationRequest{adminAuth.SessionKey, adminprint, modEvent})
	if resp.Error != "" {
		t.Fatalf("Critical error modifying event %v", resp.Error)
	}
	fail := EditEvent(EventManipulationRequest{crypt.Key("wrong"), "none", modEvent})
	if fail.Error == "" {
		t.Errorf("Error Unauthorized modification allowed %v", fail.Error)
	}
	// Make user
	email		:= "new@email"
	pass		:= crypt.Key("asdfghjk")
	fingerprint	:= "0.0.0.0"
	size	 	:= 1
	got := Register(RegisterRequest{email, pass, crypt.GenerateHash(fingerprint)})
	if got.Error != "" {
		t.Fatalf("Expected: %v, Got: %v\n", nil, got.Error)
	}
	// Tests
	ress	 := GetUserReservatoions(UserReservationsRequest{got.SessionKey})
	if len(ress) != 0 {
		t.Errorf("Expected: %v, Got: %v\n", 0, len(ress))
	}
	events := GetEvents(EventRequest{})
	if len(events) != 2 {
		t.Errorf("Expected: %v, Got: %v\n", 2, len(events))
	}
	// Check both events are accounted for
	countA := 0
	countB := 0
	for i := range(events) {
		if events[i].DtEnd == 1735687830 {
			countA++
		} else {
			countB++
		}
	}
	if countA != 1 || countB != 1 {
		t.Fatalf("Expected: %v-%v, Got: %v-%v\n", 1, 1, countA, countB)
	}
	event := GetEvent(EventRequest{EventID: resp.EventID})
	if event.DtStart != 1735675270 {
		t.Fatalf("Expected: %v, Got: %v\n", 1735675270, event.DtStart)
	}
	// Make reservation
	res := MakeReservation(ReserveRequest{
		crypt.ID(""), got.SessionKey, email, fingerprint, crypt.Hash(""), size, resp.EventID, utils.Epoch(1100),
	})
	if res.Error != "" {
		t.Fatalf("Expected: %v, Got: %v\n", nil, res.Error)
	}
	// Tests
	ress = GetUserReservatoions(UserReservationsRequest{got.SessionKey})
	if ress[0].EventID != resp.EventID {
		t.Fatalf("Expected: %v, Got: %v\n", resp.EventID, ress[0].EventID)
	}
	if len(ress) != 1 {
		t.Fatalf("Expected: %v, Got: %v\n", 1, len(ress))
	}
	if res.Id != ress[0].Id {
		t.Fatalf("Expected: %v, Got: %v\n", res.Id, ress[0].Id)
	}
	// Admin check
	ress = AdminListUserReservatoions(EnhancedUserRequest{
		User: email,
		SessionKey: adminAuth.SessionKey,
		Fingerprint: adminprint,
		HashPrint: crypt.Hash(adminprint),
		Password: adminPassword,
	})
	if len(ress) != 1 {
		t.Fatalf("Expected: %v, Got: %v\n", 1, len(ress))
	}
	if res.Id != ress[0].Id {
		t.Fatalf("Expected: %v, Got: %v\n", res.Id, ress[0].Id)
	}

	// Delete Event
	resp = DeleteEvent(
		EventManipulationRequest{adminAuth.SessionKey, fingerprint, model.Event{ID: event.ID}},
	)
	if resp.Error != "" {
		t.Fatalf("Event removal failed\n")
	}
	resp = DeleteEvent(
		EventManipulationRequest{adminAuth.SessionKey, fingerprint, model.Event{ID: event.ID}},
	)
	if resp.Error == "" {
		t.Fatalf("%v\n", resp.Error)
	}
	// admin check
	ress = AdminListUserReservatoions(EnhancedUserRequest{
		User: email,
		SessionKey: adminAuth.SessionKey,
		Fingerprint: adminprint,
		HashPrint: crypt.Hash(adminprint),
		Password: adminPassword,
	})
	// TODO: Delete reservations for deleted event
	t.Log("TODO: Delete reservations for deleted event")
	//if len(ress) != 0 {
	//	t.Errorf("Expected: %v, Got: %v\n", 0, len(ress))
	//}
}
