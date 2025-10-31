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

func TestLogin(t *testing.T) {
	if _, err := testLoginUser("login"); err != nil {
		t.Errorf("Error in response handler:\n %v\n", err)
	}
}

func TestEventLifesycle(t *testing.T) {
	setupFirstAdminUser("admin", deterministicKeyGenerator)
	client, err := state.NewClient("admin", "test-admin", "adminpasswordexample", false)
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
	_, err = testReserve(sessionKey, "test-admin", 1, state.ID(eventID), client.Id)
	if err != nil {
		t.Fatalf("Response handler:\n%v\n", err)
	}
}

func TestUnregisteredReservation(t *testing.T) {
	/*
	_, err = testEventLogin(state.ID(eventID))
	if err != nil {
		t.Fatalf("Response handler:\n%v\n", err)
	}
	*/
}
