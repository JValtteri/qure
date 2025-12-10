package state

import (
	"github.com/JValtteri/qure/server/internal/crypt"
	"github.com/JValtteri/qure/server/internal/state/model"
)


var clients model.Clients = model.Clients{
	ByID:		make(map[crypt.ID]*model.Client),
	BySession:	make(map[crypt.Key]*model.Client),
	ByEmail:	make(map[string]*model.Client),
}

var events map[crypt.ID]model.Event = make(map[crypt.ID]model.Event)

var reservations model.Reservations = model.Reservations{
	ByID:		make(map[crypt.ID]model.Reservation),
	ByEmail:	make(map[string]*model.Reservation),
}


func init() {
	Load("db.gob")		// Resumes previously saved state
}

func ResetClients() {
	clients = model.Clients{
		ByID:		make(map[crypt.ID]*model.Client),
		BySession:	make(map[crypt.Key]*model.Client),
		ByEmail:	make(map[string]*model.Client),
	}
}

func ResetEvents() {
	events = make(map[crypt.ID]model.Event)
}
