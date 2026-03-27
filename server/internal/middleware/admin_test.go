package middleware

import (
	"testing"

	"github.com/JValtteri/qure/server/internal/utils"
	"github.com/JValtteri/qure/server/internal/crypt"
	"github.com/JValtteri/qure/server/internal/state"
	"github.com/JValtteri/qure/server/internal/testjson"
)

func TestNotAdminMakeEvent(t *testing.T) {
	state.ResetEvents()
	state.ResetClients()
	fingerprint := "0.0.0.0"
	newEvent := state.EventFromJson(testjson.EventJson)
	resp := MakeEvent(EventManipulationRequest{crypt.Key("somekey"), fingerprint, newEvent})
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

func TestAdminGetAllUsers(t *testing.T) {
	/* Setup Admin */
	var adminName = "admin@test"
	var password = "asdfghjk"
	var fingerprint = "authenticprint"
	state.NewClient("admin", adminName, crypt.Key(password), false)
	got := Login(LoginRequest{
		User: adminName,
		Password: crypt.Key(password),
		HashPrint: crypt.GenerateHash(fingerprint),
	})
	resList := ListAllUsers(AuthenticateRequest{"wrong key", fingerprint})
	if len(resList) != 0 {
		t.Errorf("Expected: %v, Got: %v\n", 0, resList)
	}
	resList = ListAllUsers(AuthenticateRequest{got.SessionKey, fingerprint})
	if len(resList) == 0 {
		t.Errorf("Expected: %v, Got: %v\n", "1 or more", len(resList))
	}
}

func TestGetEventReservations(t *testing.T) {
	state.ResetEvents()
	state.ResetClients()
	/* Setup Admin */
	var adminName = "admin@test"
	var password = "asdfghjk"
	var fingerprint = "authenticprint"
	state.NewClient("admin", adminName, crypt.Key(password), false)
	got := Login(LoginRequest{
		User: adminName,
		Password: crypt.Key(password),
		HashPrint: crypt.GenerateHash(fingerprint),
	})

	/* Setup Event */
	newEvent := state.EventFromJson(testjson.EventJson)
	eventResp := MakeEvent(EventManipulationRequest{got.SessionKey, fingerprint, newEvent})
	MakeReservation(ReserveRequest{
		crypt.ID(""), got.SessionKey, adminName, fingerprint, crypt.Hash(""), 1, eventResp.EventID, utils.Epoch(1100),
	})

	/* Test Wrong */
	resList := GetEventReservations(EventRequest{
		EventID:		eventResp.EventID,
		SessionKey:		"wrong key",
		Fingerprint:	fingerprint,
	})
	if len(resList) != 0 {
		t.Errorf("Expected: %v, Got: %v\n", 0, resList)
	}

	/* Test Right */
	resList = GetEventReservations(EventRequest{
		EventID:		eventResp.EventID,
		SessionKey:		got.SessionKey,
		Fingerprint:	fingerprint,
	})

	if len(resList) == 0 {
		t.Errorf("Expected: %v, Got: %v\n", "1 or more", len(resList))
	}
}

func TestStaffAuthority(t *testing.T) {
	state.ResetClients()
	var hasAuthority bool
	var err error
	/* Setup Admin */
	adminName := "admin@test"
	adminPass := crypt.Key("asdfghjk")
	adminFingerprint := "authenticprint"
	state.NewClient("admin", adminName, adminPass, false)
	/* Setup Staff */
	staffName := "staff@test"
	staffPass := crypt.Key("jvwnoivw")
	staffFingerprint := "authenticstaff"
	state.NewClient("staff", staffName, staffPass, false)
	/* Setup User */
	userName := "user@test"
	userPass := crypt.Key("7331")
	userFingerprint := "authenticuser"
	state.NewClient("guest", userName, userPass, false)
	/* Setup Sessions */
	adminAuth := Login(LoginRequest{
		User: adminName,
		Password: adminPass,
		HashPrint: crypt.GenerateHash(adminFingerprint),
	})
	staffAuth := Login(LoginRequest{
		User: staffName,
		Password: staffPass,
		HashPrint: crypt.GenerateHash(staffFingerprint),
	})
	userAuth := Login(LoginRequest{
		User: userName,
		Password: userPass,
		HashPrint: crypt.GenerateHash(userFingerprint),
	})
	// Tests
	hasAuthority, _ = staffAuthority(userAuth.SessionKey, userFingerprint)
	if hasAuthority {
		t.Fatalf("Regular user shouldn't have staff authority")
	}
	hasAuthority, err = staffAuthority(staffAuth.SessionKey, staffFingerprint)
	if !hasAuthority {
		t.Errorf("Staff user should have staff authority: %v", err)
	}
	hasAuthority, _ = staffAuthority(adminAuth.SessionKey, adminFingerprint)
	if !hasAuthority {
		t.Errorf("Admin user should have staff authority")
	}
}

func TestAdminAuthority(t *testing.T) {
	state.ResetClients()
	var hasAuthority bool
	var err error
	/* Setup Admin */
	adminName := "admin@test"
	adminPass := crypt.Key("asdfghjk")
	adminFingerprint := "authenticprint"
	state.NewClient("admin", adminName, adminPass, false)
	/* Wrong Name */
	adminAuth := Login(LoginRequest{
		User: "wr0ngnam3",
		Password: adminPass,
		HashPrint: crypt.GenerateHash(adminFingerprint),
	})
	hasAuthority, err = adminAuthority(adminAuth.SessionKey, adminFingerprint)
	if hasAuthority {
		t.Errorf("Expected: %v, Got: %v, %v\n", false, hasAuthority, err)
	}
	/* Wrong Password */
	adminAuth = Login(LoginRequest{
		User: adminName,
		Password: "t0ta11ywr0ngpassw0rd",
		HashPrint: crypt.GenerateHash(adminFingerprint),
	})
	hasAuthority, err = adminAuthority(adminAuth.SessionKey, adminFingerprint)
	if hasAuthority {
		t.Errorf("Expected: %v, Got: %v, %v\n", false, hasAuthority, err)
	}
	/* Correct Request */
	adminAuth = Login(LoginRequest{
		User: adminName,
		Password: adminPass,
		HashPrint: crypt.GenerateHash(adminFingerprint),
	})
	hasAuthority, err = adminAuthority(adminAuth.SessionKey, adminFingerprint)
	if !hasAuthority {
		t.Errorf("Expected: %v, Got: %v, %v\n", true, hasAuthority, err)
	}
}


func TestAdminRemoveUser(t *testing.T) {
	/* Setup Admin */
	adminName := "admin@test"
	adminPass := crypt.Key("asdfghjk")
	adminFingerprint := "authenticprint"
	state.NewClient("admin", adminName, adminPass, false)
	adminAuth := Login(LoginRequest{
		User: adminName,
		Password: adminPass,
		HashPrint: crypt.GenerateHash(adminFingerprint),
	})
	/* Setup User */
	user := "remove@example"
	pass := crypt.Key("12345678")
	fingerprint := "0.0.0.0"
	got := Register(RegisterRequest{user, pass, crypt.GenerateHash(fingerprint)})
	if got.Error != "" {
		t.Fatalf("Client wasn't created: %v", got.Error)
	}
	_, ok := state.GetClientByEmail(user)
	if !ok {
		t.Fatalf("Client wasn't found")
	}
	/* Test Wrong Session key */
	resp := AdminRemoveUser(EnhancedUserRequest{
		User: user,
		SessionKey: "wrong key",
		Fingerprint: adminFingerprint,
		HashPrint: crypt.GenerateHash(adminFingerprint),
		Password: adminPass,
	})
	if resp.Error == "" || resp.Success {
		t.Errorf("Expected: %v, Got: %v\n", "error", resp.Error)
	}
	auth := Login(LoginRequest{user, pass, crypt.GenerateHash(fingerprint)})
	if !auth.Authenticated {
		t.Errorf("Expected: %v, Got: %v\n", true, auth.Authenticated)
	}
	/* Test Wrong Password */
	resp = AdminRemoveUser(EnhancedUserRequest{
		User: user,
		SessionKey: adminAuth.SessionKey,
		Fingerprint: adminFingerprint,
		HashPrint: crypt.GenerateHash(adminFingerprint),
		Password: "wrong-pass",
	})
	if resp.Error == "" || resp.Success {
		t.Errorf("Expected: %v, Got: %v\n", "error", resp.Error)
	}
	auth = Login(LoginRequest{user, pass, crypt.GenerateHash(fingerprint)})
	if !auth.Authenticated {
		t.Errorf("Expected: %v, Got: %v\n", true, auth.Authenticated)
	}
	/* Test Wrong User */
	resp = AdminRemoveUser(EnhancedUserRequest{
		User: "nonexistant@user",
		SessionKey: adminAuth.SessionKey,
		Fingerprint: adminFingerprint,
		HashPrint: crypt.GenerateHash(adminFingerprint),
		Password: adminPass,
	})
	if resp.Error == "" || resp.Success {
		t.Errorf("Expected: %v, Got: %v\n", "error", resp.Error)
	}
	auth = Login(LoginRequest{user, pass, crypt.GenerateHash(fingerprint)})
	if !auth.Authenticated {
		t.Errorf("Expected: %v, Got: %v\n", true, auth.Authenticated)
	}
	/* Test Correct */
	resp = AdminRemoveUser(EnhancedUserRequest{
		User: user,
		SessionKey: adminAuth.SessionKey,
		Fingerprint: adminFingerprint,
		HashPrint: crypt.GenerateHash(adminFingerprint),
		Password: adminPass,
	})
	if resp.Error != "" {
		t.Errorf("Expected: %v, Got: %v\n", "''", resp.Error)
	}
	auth = Login(LoginRequest{user, pass, crypt.GenerateHash(fingerprint)})
	if auth.Authenticated {
		t.Errorf("Expected: %v, Got: %v\n", false, auth.Authenticated)
	}
}

func TestAdminChangeUserRole(t *testing.T) {
		/* Setup Admin */
	adminName := "admin@test"
	adminPass := crypt.Key("asdfghjk")
	adminFingerprint := "authenticprint"
	state.NewClient("admin", adminName, adminPass, false)
	adminAuth := Login(LoginRequest{
		User: adminName,
		Password: adminPass,
		HashPrint: crypt.GenerateHash(adminFingerprint),
	})
	/* Setup User */
	user := "remove@example"
	pass := crypt.Key("12345678")
	fingerprint := "0.0.0.0"
	got := Register(RegisterRequest{user, pass, crypt.GenerateHash(fingerprint)})
	if got.Error != "" {
		t.Fatalf("Client wasn't created: %v", got.Error)
	}
	_, ok := state.GetClientByEmail(user)
	if !ok {
		t.Fatalf("Client wasn't found")
	}
	/* Test Correct */
	var newRole = "test-role"
	var resp = AdminChangeUserRole(RoleChangeRequest{
		User: user,
		Role: newRole,
		SessionKey: adminAuth.SessionKey,
		Fingerprint: adminFingerprint,
		HashPrint: crypt.GenerateHash(adminFingerprint),
		Password: adminPass,
	})
	if resp.Error != "" || !resp.Success {
		t.Errorf("Expected: %v, Got: %v\n", "''", resp.Error)
	}
	auth := Login(LoginRequest{user, pass, crypt.GenerateHash(fingerprint)})
	if !auth.Authenticated {
		t.Errorf("Expected: %v, Got: %v\n", true, auth.Authenticated)
	}
	client, ok := state.GetClientByEmail(user)
	if !ok {
		t.Fatalf("Client wasn't found")
	}
	if client.Role != newRole {
		t.Errorf("Expected: %v, Got: %v\n", newRole, client.Role)
	}
	/* Test Wrong Password */
	resp = AdminChangeUserRole(RoleChangeRequest{
		User: user,
		Role: newRole,
		SessionKey: adminAuth.SessionKey,
		Fingerprint: adminFingerprint,
		HashPrint: crypt.GenerateHash(adminFingerprint),
		Password: "foobar",
	})
	if resp.Error == "" || resp.Success {
		t.Errorf("Expected: %v, Got: %v\n", "error", resp.Error)
	}
	/* Test Invalid Target User */
	resp = AdminChangeUserRole(RoleChangeRequest{
		User: "nonexistant@user",
		Role: newRole,
		SessionKey: adminAuth.SessionKey,
		Fingerprint: adminFingerprint,
		HashPrint: crypt.GenerateHash(adminFingerprint),
		Password: adminPass,
	})
	if resp.Error == "" || resp.Success {
		t.Errorf("Expected: %v, Got: %v\n", "error", resp.Error)
	}
}

func TestBadListUserReservations(t *testing.T) {
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
	// Test bad user
	ress := AdminListUserReservatoions(EnhancedUserRequest{
		User: "",
		SessionKey: adminAuth.SessionKey,
		Fingerprint: adminprint,
		HashPrint: crypt.Hash(adminprint),
		Password: adminPassword,
	})
	if len(ress) != 0 {
		t.Errorf("Expected: %v, Got: %v\n", 0, len(ress))
	}
	// Test bad auth
	ress = AdminListUserReservatoions(EnhancedUserRequest{
		User: adminName,
		SessionKey: adminAuth.SessionKey,
		Fingerprint: adminprint,
		HashPrint: crypt.Hash(adminprint),
		Password: "wrong pass",
	})
	if len(ress) != 0 {
		t.Errorf("Expected: %v, Got: %v\n", 0, len(ress))
	}
}
