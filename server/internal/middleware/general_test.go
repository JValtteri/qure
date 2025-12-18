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
	adminPassword := crypt.Key("adminpass")
	fingerprint := "0.0.0.0"

	adminClient, err := state.NewClient("admin", "admin", adminPassword, false)
	if err != nil {
		t.Fatalf("Error generating test-admin account:\n%v", err)
	}
	auth := checkPasswordAuthentication(adminClient, adminPassword, crypt.GenerateHash(fingerprint))

	// Make events
	newEvent := state.EventFromJson(testjson.EventJson)
	resp := MakeEvent(EventManipulationRequest{auth.SessionKey, fingerprint, "", newEvent})
	emptyEvent := state.EventFromJson([]byte("{}"))
	_   = MakeEvent(EventManipulationRequest{auth.SessionKey, fingerprint, "", emptyEvent})		// Decoy event
	if resp.EventID == crypt.ID("") {
		t.Fatal("Critical error making event")
	}

	// Modify Event
	modEvent, err := state.GetEvent(resp.EventID, true)
	if resp.Error != "" {
		t.Fatalf("Critical fetching event %v", err)
	}
	modEvent.Name = "Updated Event"
	resp = EditEvent(EventManipulationRequest{auth.SessionKey, fingerprint, "", modEvent})
	if resp.Error != "" {
		t.Fatalf("Critical error modifying event %v", resp.Error)
	}
	fail := EditEvent(EventManipulationRequest{crypt.Key("wrong"), "none", "", modEvent})
	if fail.Error == "" {
		t.Errorf("Error Unauthorized modification allowed %v", fail.Error)
	}
	// Make user
	email	 := "new@email"
	pass	 := crypt.Key("asdfghjk")
	size	 := 1
	got := Register(RegisterRequest{email, pass, crypt.GenerateHash(fingerprint)})
	if got.Error != "" {
		t.Fatalf("Expected: %v, Got: %v\n", nil, got.Error)
	}
	ress	 := GetUserReservatoions(UserReservationsRequest{got.SessionKey})
	if len(ress.Reservations) != 0 {
		t.Errorf("Expected: %v, Got: %v\n", 0, len(ress.Reservations))
	}
	// Check Events
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
		got.SessionKey, email, fingerprint, crypt.Hash(""), size, resp.EventID, utils.Epoch(1100),
	})
	if res.Error != "" {
		t.Fatalf("Expected: %v, Got: %v\n", nil, res.Error)
	}
	ress = GetUserReservatoions(UserReservationsRequest{got.SessionKey})
	if ress.Reservations[0].EventID != resp.EventID {
		t.Fatalf("Expected: %v, Got: %v\n", resp.EventID, ress.Reservations[0].EventID)
	}
	if len(ress.Reservations) != 1 {
		t.Fatalf("Expected: %v, Got: %v\n", 1, len(ress.Reservations))
	}
	if res.Id != ress.Reservations[0].Id {
		t.Fatalf("Expected: %v, Got: %v\n", res.Id, ress.Reservations[0].Id)
	}
	resp = DeleteEvent(
		EventManipulationRequest{auth.SessionKey, fingerprint, event.ID, model.Event{}},
	)
	if res.Error != "" {
		t.Fatalf("Event removal failed\n")
	}
}

func TestNotAdminMakeEvent(t *testing.T) {
	state.ResetEvents()
	state.ResetClients()
	fingerprint := "0.0.0.0"
	newEvent := state.EventFromJson(testjson.EventJson)
	resp := MakeEvent(EventManipulationRequest{crypt.Key("somekey"), fingerprint, "", newEvent})
	if resp.EventID != crypt.ID("") {
		t.Errorf("Expected: %v, Got: %v\n", "''", resp.EventID)
	}
}

func TestAdminListEvents(t *testing.T) {
	/* Setup events */
	openEvent := state.EventFromJson(testjson.EventJson)
    _, err := state.CreateEvent(openEvent)
	if err != nil {
		t.Fatalf("Got error: %v", err)
	}
	draftEvent := state.EventFromJson(testjson.EventJson)
	draftEvent.Draft = true;
	draftEvent.DtStart = utils.EpochNow()
    _, err = state.CreateEvent(draftEvent)
	if err != nil {
		t.Fatalf("Got error: %v", err)
	}
	/* Setup Admin */
	adminName := "admin@test"
	password := "asdfghjk"
	fingerprint := "authenticprint"
	state.NewClient("admin", adminName, crypt.Key(password), false)
	got := Login(LoginRequest{
		User: adminName,
		Password: crypt.Key(password),
		HashPrint: crypt.GenerateHash(fingerprint),
	})
	/* Get events as guest */
	events1 := GetEvents(EventRequest{})
	if len(events1) != 1 {
		t.Errorf("Expected: %v, Got: %v\n", 1, len(events1))
	}
	/* Get events as admin */
	events2 := GetEvents(EventRequest{
		SessionKey: got.SessionKey,
		Fingerprint: fingerprint,
	})
	if len(events2) != 2 {
		t.Errorf("Expected: %v, Got: %v\n", 2, len(events2))
	}
}
