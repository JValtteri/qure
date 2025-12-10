package saveload

import (
	"os"
	"maps"
	"encoding/gob"

	"github.com/JValtteri/qure/server/internal/crypt"
	"github.com/JValtteri/qure/server/internal/state/model"
)


type PersistedStore struct {
	Clients			map[crypt.ID]*model.Client			// clients.byID
	Events			map[crypt.ID]model.Event			// events
	Reservations	map[crypt.ID]model.Reservation		// reservations.byID
}

func SaveGob(
	filename string,
	clients *model.Clients,
	events *map[crypt.ID]model.Event,
	reservations *model.Reservations,
) error {
	store := emptyStore()
	clients.RLock()
	reservations.RLock()

	defer clients.RUnlock()
	defer reservations.RUnlock()

	maps.Copy(store.Clients, clients.ByID)
	maps.Copy(store.Events, *events)
	maps.Copy(store.Reservations, reservations.ByID)

	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	enc := gob.NewEncoder(file)
	if err := enc.Encode(store); err != nil {
		return err
	}
	return nil
}

func LoadGob(filename string) (PersistedStore, error) {
	var store PersistedStore
	file, err := os.Open(filename)
	if err != nil {
		return store, err
	}
	defer file.Close()
	if err := gob.NewDecoder(file).Decode(&store); err != nil {
		return store, err
	}
	return store, nil
}

func emptyStore() *PersistedStore {
	return &PersistedStore {
		Clients:		make(map[crypt.ID]*model.Client),
		Events:			make(map[crypt.ID]model.Event),
		Reservations:	make(map[crypt.ID]model.Reservation),
	}
}

func init() {
	// Registers every concrete type that appears in the graph.
	// gob must know the type to serialize a pointer.
	gob.Register(&model.Client{})
	gob.Register(&model.Session{})
	gob.Register(&model.Event{})
	gob.Register(&model.Timeslot{})
	gob.Register(&model.Reservation{})
}
