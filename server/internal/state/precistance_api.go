package state

import (
	"log"
	"sync"

	"github.com/JValtteri/qure/server/internal/crypt"
	"github.com/JValtteri/qure/server/internal/saveload"
	"github.com/JValtteri/qure/server/internal/state/model"
)


var wg *sync.WaitGroup
var isSetUp = false

type SaveStatus struct {
	mu			sync.RWMutex
	inProgress	bool
}

func (s *SaveStatus) Lock() {
	s.mu.Lock()
}

func (s *SaveStatus) Unlock() {
	s.mu.Unlock()
}

var saveStatus = SaveStatus{}


func InitWaitGroup(group *sync.WaitGroup) {
	isSetUp = true
	wg = group
}

// Save() is thread safe and can be called multiple times
// Only one instance will run
func Save(filename string) {
	if !isSetUp {
		log.Fatal("Precistance is not initialized. Run AddToWaitGroup() before running Save()")
	}
	wg.Add(1)
	defer wg.Done()
	saveInProgress := reserveFile()
	if saveInProgress {
		return
	}
	defer freeFile()
	saveload.SaveGob(filename, &clients, &events, &reservations)
}

func reserveFile() bool {
	saveStatus.Lock()
	defer saveStatus.Unlock()
	if saveStatus.inProgress {
		return true
	}
	saveStatus.inProgress = true
	return false
}

func freeFile() {
	saveStatus.Lock()
	defer saveStatus.Unlock()
	saveStatus.inProgress = false
}

func Load(filename string) {
	store, err := saveload.LoadGob(filename)
	if err != nil {
		log.Println(err)
	} else {
		clients.ByID		= store.Clients
		events				= store.Events
		reservations.ByID	= store.Reservations
		rebuildIndexes()
	}
}

func rebuildIndexes() {
	log.Println("Re-indexing state data...")
	restoreClients()								// Clients: bySession, byEmail
	reIndexReservations()							// Reservations: byEmail
}

func restoreClients() {
	for _, client := range clients.ByID {
		reIndexClientByEmail(client)
		reIndexClientBySession(client)
	}
}

func reIndexReservations() {
	var invalidReservations []crypt.ID
	for _, reservation := range reservations.ByID {
		if reservation.Client == "" {
			break
		}
		// This shouldn't be needed, but if the data becomes desynced, this
		// prevents saveload from crashing and cleans up the worst offenders
		if clients.ByID[reservation.Client] == nil {
			invalidReservations = append(invalidReservations, reservation.Id)
			if reservation.Event != nil {
				removeInvalidReservationFromEvent(reservation)
			}
			continue
		}
		clientEmail := clients.ByID[reservation.Client].Email
		reservations.ByEmail[clientEmail] = &reservation
	}
	// If any dangling reservations are encountered, they are removed
	for _, invalidReservationID := range invalidReservations {
		log.Println(invalidReservationID)
		delete(reservations.ByID, invalidReservationID)
	}
}

// This shouldn't be needed, but if the data becomes desynced,
// this is used to clean dangling reservations
func removeInvalidReservationFromEvent(reservation model.Reservation) {
	var timeslot = reservation.Event.Timeslots[reservation.Timeslot]
	var eventReservations = timeslot.Reservations
	var eventQueue = timeslot.Queue
	// remove invalid from Reservations and Queue
	eventReservations	= model.FilterFrom(eventReservations, reservation.Id)
	eventQueue			= model.FilterFrom(eventQueue, reservation.Id)
	// Substitute back modified lists
	timeslot.Reservations	= eventReservations
	timeslot.Queue 			= eventQueue
	timeslot.Reserved = len(eventReservations)
	reservation.Event.Timeslots[reservation.Timeslot] = timeslot
}

func reIndexClientBySession(client *model.Client) {
	for sessionKey := range client.Sessions {
		clients.BySession[sessionKey] = client
	}
}

func reIndexClientByEmail(client *model.Client) {
	clients.ByEmail[client.Email] = client
}
