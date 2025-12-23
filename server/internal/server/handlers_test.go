package server

import (
	"testing"

	"github.com/JValtteri/qure/server/internal/crypt"
	"github.com/JValtteri/qure/server/internal/state"
)


func TestGetEvents(t *testing.T) {
	if _, err := testGetEvents(); err != nil {
		t.Logf("Error in response handler: %v\n", err)
	}
}

func TestGetEvent(t *testing.T) {
	if _, err := testGetEvent(crypt.ID("nothing")); err != nil {
		t.Errorf("Expected: '%s', Got: '%v'\n", "error", err)
	}
}

func TestRegisterUser(t *testing.T) {
	if _, err := testRegisterUser("test"); err != nil {
		t.Errorf("Error in response handler:\n %v\n", err)
	}
}

func TestRegisterUserTwice(t *testing.T) {
	testRegisterDuplicateUser("double")
	if _, err := testRegisterDuplicateUser("double"); err != nil {
		t.Errorf("Error in response handler:\n %v\n", err)
	}
}

func TestResumeSession(t *testing.T) {
	if _, err := testResumeSession("test"); err != nil {
		t.Errorf("Error in response handler:\n %v\n", err)
	}
}

func TestEventLifesycle(t *testing.T) {
	setupFirstAdminUser("admin", deterministicKeyGenerator)
	setupFirstAdminUser("admin", deterministicKeyGenerator)	// Nothing should happen on second call
	_, err := state.NewClient("admin", "test-admin", "adminpasswordexample", false)
	if err != nil {
		t.Fatalf("Error generating test-admin account:\n%v", err)
	}
	sessionKey, err := testLoginAdmin("test-admin")
	if err != nil {
		t.Fatalf("Response handler:\n%v\n", err)
	}
	eventID, err := testMakeEvent(sessionKey)
	if err != nil {
		t.Fatalf("Response handler:\n%v\n", err)
	}
	if len(eventID) < 9 {
		t.Fatalf("Unexpected EventID: %v\n", eventID)
	}
	_, err = testReserve(sessionKey, "test-admin", 1, crypt.ID(eventID))
	if err != nil {
		t.Fatalf("Response handler:\n%v\n", err)
	}
	//if clientID != string(client.Id) {								// These tests should be replaced with something
	//	t.Errorf("Expected: %v, Got: %v\n", clientID, client.Id)
	//}

	// Test Unregistered Reservation and Login
	tempResID, err := testReserve("no-key", "anonymous@account.not", 1, crypt.ID(eventID))
	if err != nil {
		t.Fatalf("Response handler:\n%v\n", err)
	}
	//if tempClientID == string(client.Id) {							// These tests should be replaced with something
	//	t.Errorf("Temp client was given admin's client ID:\n%v\n", tempClientID)
	//}

	tempSessionKey, err := testLoginUser("anonymous@account.not", crypt.Key(tempResID))
	if err != nil {
		t.Fatalf("Response handler:\n%v\n", err)
	}

	newUserSession, err := testRegisterUser("thirduser")
	if err != nil {
		t.Fatalf("Response handler:\n%v\n", err)
	}
	// admins reservations
	_, err = testUserReservations(crypt.Key(sessionKey), crypt.ID(eventID))
	if err != nil {
		t.Errorf("Response handler:\n%v\n", err)
	}
	// temp user reservations
	_, err = testUserReservations(crypt.Key(tempSessionKey), crypt.ID(eventID))
	if err != nil {
		t.Errorf("Response handler:\n%v\n", err)
	}
	_, err = testUserReservations(crypt.Key(newUserSession), crypt.ID(eventID))
	if err == nil {
		t.Errorf("This user shouldn't have reservations!\n")
	}
}

