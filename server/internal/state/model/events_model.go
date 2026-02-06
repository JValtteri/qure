package model

import (
	"slices"
	"sync"

	"github.com/JValtteri/qure/server/internal/crypt"
	"github.com/JValtteri/qure/server/internal/utils"
)

var Eventslock sync.RWMutex = sync.RWMutex{}

type Event struct {
	ID					crypt.ID
	Name				string
	ShortDescription	string
	LongDescription		string
	Draft				bool
	DtStart				utils.Epoch
	DtEnd				utils.Epoch
	StaffSlots			int
	Staff				int
	Timeslots			map[utils.Epoch]Timeslot
}


func addNElementsToList(number int, item crypt.ID, list []crypt.ID) []crypt.ID {
	items := slices.Repeat([]crypt.ID{item}, number)
	return append(list, items...)
}
