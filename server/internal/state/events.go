package state

import (
	"log"
	"fmt"
	"github.com/JValtteri/qure/server/internal/utils"
	"github.com/JValtteri/qure/server/internal/crypt"
	"github.com/JValtteri/qure/server/internal/state/model"
)


func CreateEvent(eventObj model.Event) (crypt.ID, error) {
	id := setId(&eventObj)
	model.Eventslock.Lock()
	defer model.Eventslock.Unlock()
    _, duplicate := events[id]
    if duplicate {
        return "0", fmt.Errorf("duplicate Event: %v", id)
    }
    events[eventObj.ID] = eventObj
    return id, nil
}

func GetEvent(id crypt.ID, isAdmin bool) (model.Event, error) {
	model.Eventslock.RLock()
	defer model.Eventslock.RUnlock()
    event, ok := events[id]
    if !ok {
		return model.Event{}, fmt.Errorf("event not found: %v", id)
    }
    if event.Draft && !isAdmin {
		return model.Event{}, fmt.Errorf("event not found: %v", id)
	}
    return event, nil
}

func GetEvents(isAdmin bool) []model.Event {
	var outEvents []model.Event
	model.Eventslock.RLock()
	defer model.Eventslock.RUnlock()
    for _, obj := range events {
        if !obj.Draft {
            outEvents = append(outEvents, obj)
        }
    }
    return outEvents
}

func RemoveEvent(id crypt.ID) bool {
	model.Eventslock.Lock()
	defer model.Eventslock.Unlock()
    _, ok := events[id]
    if !ok {
        return false
    }
    delete(events, id)
    return true
}

func ListEvents() {
    e := GetEvents(false)
    log.Println("Events: ", len(e))
    for _, obj := range e {
        log.Println("->   ", obj.Name, obj.ID)
    }
}

func setId(event *model.Event) crypt.ID {
    currentTime := utils.EpochNow()
    currentSeconds := currentTime % 60
    uID := event.DtStart + currentSeconds
    newID := crypt.ID(fmt.Sprintf("%v", uID))
    event.ID = newID
    return newID
}

// Only for testing low level functions
// handlers use middleware function
func EventFromJson(input []byte) model.Event {
	var event model.Event
	utils.LoadJSON(input, &event)
	return event
}
