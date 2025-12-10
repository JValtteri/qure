package saveload

import (
	"log"
	"os"
	"testing"

	"github.com/JValtteri/qure/server/internal/crypt"
	"github.com/JValtteri/qure/server/internal/state/model"
	"github.com/JValtteri/qure/server/internal/testjson"
	"github.com/JValtteri/qure/server/internal/utils"
)


const testFileName = "test.gob"

// Client 1 data
var id1		= crypt.ID("id1")
var epoch1	= utils.Epoch(10)
var email1	= "first@emial.com"
var pass1	= crypt.Key("firstpass")
var role1	= "user"

// Client 2 data
var id2		= crypt.ID("id2")
var epoch2	= utils.Epoch(20)
var email2	= "second@emial.com"
var pass2	= crypt.Key("secondpass")
var role2	= "admin"

// Event data
var eventid = crypt.ID("event1")
var time	= utils.Epoch(111)
var size	= 5
var rsvd	= 1

// Reservation data
var resid	= crypt.ID("222")


func TestFailLoad(t *testing.T) {
	os.Remove(testFileName)
	_, err := LoadGob(testFileName)
	if err == nil {
		t.Fatal("Loading GOB should fail when GOB file doesn't exist")
	}
}

func TestSaveAndLoad(t *testing.T) {

	// --- Create data and save state --- //
	clients, events, reservations := initState()
	err := SaveGob(testFileName, &clients, &events, &reservations)
	if err != nil {
		t.Error(err)
	}

	// --- Make sure state is cleared before loading save --- //
	clients, events, reservations = clearState()
	testStateIsClean(t, &clients, &reservations, events)

	// --- Load state from save --- //
	store, err := LoadGob("test.gob")
	clients.ByID 		= store.Clients
	events		 		= store.Events
	reservations.ByID 	= store.Reservations
	if err != nil {
		t.Error(err)
	}

	// --- Test Basic Client Data --- //
	client1, auth := testClient1State(t, &clients)
	testClient2State(t, &clients, auth)
	// --- Test Reservation --- //
	res := testReservationState(t, client1, &reservations)
	// --- Test Event --- //
	testEventState(t, res, events)
	// -- Clean Up --- //
	os.Remove(testFileName)
}

func testStateIsClean(t *testing.T, clients *model.Clients, reservations *model.Reservations, events map[crypt.ID]model.Event) {
	_, ok := clients.ByID[id1]
	if ok {
		t.Fatal("Clients state wasn't cleared before loading save")
	}
	if len(reservations.ByID) > 0 {
		t.Fatal("Reservations state wasn't cleared before loading save")
	}
	if len(events) > 0 {
		t.Fatal("Events state wasn't cleared before loading save")
	}
}

func testClient1State(t *testing.T, clients *model.Clients) (*model.Client, bool) {
	client1, ok := clients.ByID[id1]
	if !ok {
		t.Fatal("Data wasnt loaded from GOB")
	}
	if client1.Email != email1 {
		t.Errorf("Expected: %v, Got: %v\n", email1, client1.Email)
	}
	if client1.ExpiresDt != epoch1 {
		t.Errorf("Expected: %v, Got: %v\n", epoch1, client1.ExpiresDt)
	}
	auth := crypt.CompareToHash(pass1, client1.Password)
	if !auth {
		t.Errorf("Expected: %v, Got: %v\n", true, auth)
	}
	if client1.Role != role1 {
		t.Errorf("Expected: %v, Got: %v\n", role1, client1.Role)
	}
	return client1, auth
}

func testClient2State(t *testing.T, clients *model.Clients, auth bool) {
	client2, ok := clients.ByID[id2]
	if !ok {
		t.Fatal("Data wasn't loaded from GOB")
	}
	if client2.Email != email2 {
		t.Errorf("Expected: %v, Got: %v\n", email2, client2.Email)
	}
	if client2.ExpiresDt != epoch2 {
		t.Errorf("Expected: %v, Got: %v\n", epoch2, client2.ExpiresDt)
	}
	auth = crypt.CompareToHash(pass2, client2.Password)
	if !auth {
		t.Errorf("Expected: %v, Got: %v\n", true, auth)
	}
	if client2.Role != role2 {
		t.Errorf("Expected: %v, Got: %v\n", role2, client2.Role)
	}
}

func testReservationState(t *testing.T, client1 *model.Client, reservations *model.Reservations) *model.Reservation {
	if len(client1.Reservations) < 1 {
		t.Fatal("Reservations wern't loaded from GOB")
	}
	res := client1.Reservations[0]
	if res.Id != resid {
		t.Errorf("Expected: %v, Got: %v\n", resid, res.Id)
	}
	if res.Size != size {
		t.Errorf("Expected: %v, Got: %v\n", size, res.Size)
	}
	if res.Timeslot != time {
		t.Errorf("Expected: %v, Got: %v\n", time, res.Timeslot)
	}
	if res.Client != client1.Id {
		t.Errorf("Expected: %v, Got: %v\n", client1.Id, res.Client)
	}
	if reservations.ByID[resid].Size != size {
		t.Errorf("Expected: %v, Got: %v\n", size, reservations.ByID[resid].Size)
	}
	return res
}

func testEventState(t *testing.T, res *model.Reservation, events map[crypt.ID]model.Event) {
	if res.Event.DtStart != 1735675270 {
		t.Errorf("Expected: %v, Got: %v\n", 1735675270, res.Event.DtStart)
	}
	if res.Event.Timeslots[time].Size != size {
		t.Errorf("Expected: %v, Got: %v\n", size, res.Event.Timeslots[time].Size)
	}
	if res.Event.Timeslots[time].Reservations[0] != res.Id {
		t.Errorf("Expected: %v, Got: %v\n", res.Id, res.Event.Timeslots[time].Reservations[0])
	}
	if len(events) < 1 {
		t.Fatal("Events wern't loaded from GOB")
	}
	if events[eventid].Name != "Test event" {
		t.Errorf("Expected: %v, Got: %v\n", "Test event", events[eventid].Name)
	}
}


func initState() (model.Clients, map[crypt.ID]model.Event, model.Reservations) {
	// Setup a fresh clear state
	clients, events, reservations := clearState()

	// Setup Clients
	client1 := model.CreateClient(id1, epoch1, email1, pass1, role1)
	client2 := model.CreateClient(id2, epoch2, email2, pass2, role2)

	// Create an event
	event := createTestEvent()

	// Create Reservation
	reservation := model.Reservation{
		Id:			resid,
		Client:		client1.Id,
		Size:		size,
		Confirmed:	size,
		Event:		&event,
		Timeslot:   time,
	}

	client1.Reservations = append(client1.Reservations, &reservation)

	clients.ByID[id1] = client1
	clients.ByID[id2] = client2
	reservations.ByID[resid] = reservation
	events[eventid] = event
	if len(events) < 1 {
		log.Fatal("Events wern't loaded from GOB")
	}
	return clients, events, reservations
}

func clearState() (model.Clients, map[crypt.ID]model.Event, model.Reservations) {
	var clients = model.Clients{
		ByID:		make(map[crypt.ID]*model.Client),
	}
	events := make(map[crypt.ID]model.Event)
	var reservations = model.Reservations{
		ByID:		make(map[crypt.ID]model.Reservation),
	}
	return clients, events, reservations
}

func createTestEvent() model.Event {
	var event model.Event
	utils.LoadJSON(testjson.EventJson, &event)
	event.Append(
		model.Timeslot{
			Size:         size,
			Reserved:     size,
			Reservations: []crypt.ID{resid},
		},
		time,
	)
	return event
}
