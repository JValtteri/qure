package state

import (
    "sync"
    "github.com/JValtteri/qure/server/internal/crypt"
)

type Client struct {
    id           crypt.ID
    password     crypt.Hash
    createdDt    Epoch     // Unix timestamp
    expiresDt    Epoch     // Unix timestamp, 0 = expire now, 0-- = keep indefinately
    email        string
    phone        string
    role         string
    sessions map[crypt.Key]Session
    reservations []*Reservation
}

func (t *Client) AddReservation(res *Reservation) {
    t.reservations = append(t.reservations, res)
}

var clients Clients = Clients{
    byID:       make(map[crypt.ID]*Client),
    bySession:  make(map[crypt.Key]*Client),
    byEmail:    make(map[string]*Client),
}
type Clients struct {
    mu          sync.RWMutex
    byID        map[crypt.ID]*Client        // by client ID
    bySession   map[crypt.Key]*Client       // by session key
    byEmail     map[string]*Client    // by session key
}

func (c *Clients) AddReservation(id crypt.ID, reservation *Reservation) {
    c.withLock(func() {
        client := c.byID[id]
        client.AddReservation(reservation)
    })
}

func (c *Clients) withLock(fn func()) {
    c.mu.Lock()
    defer c.mu.Unlock()
    fn()
}

func (c *Clients) getClient(sessionKey crypt.Key) (*Client, bool) {
    c.mu.RLock()
    defer c.mu.RUnlock()
    client, found := clients.bySession[sessionKey]
    return client, found
}

func (c *Clients) GetReservations(id crypt.ID) []*Reservation {
    c. rLock()
    defer c.rUnlock()
    return c.byID[id].reservations
}

func (c *Clients) rLock() {
    c.mu.RLock()
}

func (c *Clients) rUnlock() {
    c.mu.RUnlock()
}

func (c *Clients) Lock() {
    c.mu.Lock()
}

func (c *Clients) Unlock() {
    c.mu.Unlock()
}
