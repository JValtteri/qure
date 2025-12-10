package state

import (
	"os"
	"sync"
	"testing"

	"github.com/JValtteri/qure/server/internal/crypt"
	"github.com/JValtteri/qure/server/internal/testjson"
	"github.com/JValtteri/qure/server/internal/utils"
)


const testFileName = "test.gob"
const email = "temp@example.com"
var testwg sync.WaitGroup

func TestSaveLoadAndReIndexClients(t *testing.T) {
	InitWaitGroup(&testwg)
	client1, err := NewClient("test", email, crypt.Key("asdf"), true)
	clientId := client1.Id
	if err != nil {
		t.Errorf("Expected: %v, Got: %v\n", nil, err)
	}

	Save(testFileName)
	ResetClients()
	Load(testFileName)

	client, ok := clients.ByEmail[email]
	if !ok {
		t.Errorf("Expected: %v, Got: %v\n", true, ok)
	}
	if client.Id != clientId {
		t.Errorf("Expected: %v, Got: %v\n", clientId, client.Id)
	}
	os.Remove(testFileName)
}

func TestSaveLoadAndReIndexReservations(t *testing.T) {
	InitWaitGroup(&testwg)
	ResetEvents()
	ResetClients()
	role := "test"
	fingerprint := "0.0.0.0"
	time := utils.Epoch(1100)
	size := 1
	temp := false
	timeslot := setTimeslot(5)
	client, err := NewClient(role, email, crypt.Key("asdf"), temp)
	if err != nil {
		t.Fatalf("Error in creating client: %v", err)
	}
	sessionKey, _ := client.AddSession(role, email, temp, crypt.GenerateHash(fingerprint), &clients)
	event := EventFromJson(testjson.EventJson)
	eventID, err := CreateEvent(event)
	event.Append(timeslot, time)
	MakeReservation(sessionKey, email, fingerprint, crypt.GenerateHash(fingerprint), size, eventID, time)

	Save(testFileName)
	ResetClients()
	Load(testFileName)
	if _, ok := reservations.ByEmail[email] ; !ok {
		t.Errorf("Expected: %v, Got: %v\n", true, ok)
	}
	if _, ok := clients.BySession[sessionKey] ; !ok {
		t.Errorf("Expected: %v, Got: %v\n", true, ok)
	}
	os.Remove(testFileName)
}

// Tests that save processes are thread safe
func TestMultipleSaves(t *testing.T) {
	InitWaitGroup(&testwg)
	defer os.Remove(testFileName)
	go Save(testFileName)
	go Save(testFileName)
	go Save(testFileName)
}
