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

func (r *Reservations) append(res Reservation, clients *Clients) error {
	r.Lock()
	defer r.Unlock()
	r.ByID[res.Id] = res
	clientEmail := clients.ByID[res.Client].Email
	r.ByEmail[clientEmail] = &res
	return clients.AddReservation(res.Client, &res)
}
