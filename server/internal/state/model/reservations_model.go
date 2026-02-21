package model

import (
	"sync"

	"github.com/JValtteri/qure/server/internal/crypt"
)

type Reservations struct {
	mu      sync.RWMutex
	ByID    map[crypt.ID]Reservation
	ByEmail map[string]*Reservation
}

func (r *Reservations) Lock() {
	r.mu.Lock()
}

func (r *Reservations) Unlock() {
	r.mu.Unlock()
}

func (r *Reservations) RLock() {
	r.mu.RLock()
}

func (r *Reservations) RUnlock() {
	r.mu.RUnlock()
}

// Updates reservation to reservations
//
// Updates:
//
// - reservations.ByID
//
// - reservations.ByEmail
//
// - client.Reservations
func (r *Reservations) update(res Reservation, clients *Clients) error {
	r.Lock()
	defer r.Unlock()
	if res.Size > 0 {
		r.ByID[res.Id] = res
		clientEmail := clients.ByID[res.Client].Email
		r.ByEmail[clientEmail] = &res
		return clients.AddReservation(res.Client, &res)
	} else {
		clientEmail := clients.ByID[res.Client].Email
		delete(r.ByID, res.Id)
		delete(r.ByEmail, clientEmail)
		return clients.RemoveReservation(res.Client, &res)
	}
}
