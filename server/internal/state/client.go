package state

import (
    "fmt"
    "sync"
    "github.com/JValtteri/qure/server/internal/crypt"
    "github.com/JValtteri/qure/server/internal/utils"
)

const TEMP_CLIENT_AGE Epoch = 60*60*24*30    // max age in seconds

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

func (c *Clients) withRLock(fn func() (*Client, bool) ) (*Client, bool) {
    c.mu.RLock()
    defer c.mu.RUnlock()
    a, b := fn()
    return a, b
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

func (c *Clients) AddReservation(id crypt.ID, reservation *Reservation) {
    c.withLock(func() {
        client := c.byID[id]
        client.AddReservation(reservation)
    })
}

func (c *Clients) GetReservations(id crypt.ID) []*Reservation {
    c. rLock()
    defer c.rUnlock()
    return c.byID[id].reservations
}

func NewClient(role string, email string, expire Epoch, sessionKey crypt.Key) (*Client, error) {
    var client Client
    uiniqueEmail := unique(email, clients.byEmail)
    if !uiniqueEmail {
        return &client, fmt.Errorf("error: client email not unique")
    } else {
        kId, err := createUniqueID(16, clients.byID)
        id := crypt.ID(kId)
        if err != nil {
            return &client, fmt.Errorf("error: Creating a new client\n%v", err) // Should not be possible (random byte generation)
        }
        client.id = crypt.ID(id)
        client.createdDt = utils.EpochNow()
        client.expiresDt = expire
        client.email = email
        client.phone = ""
        client.role = role
        client.sessions = make(map[crypt.Key]Session)
        // client.reservations = []  // make sure it's empty
        clients.withLock(func() {
                clients.byID[id] = &client;
                clients.byID[id] = &client;
                clients.bySession[sessionKey] = &client;
                clients.byEmail[email] = &client;
        })
        return &client, nil
    }
}

func RemoveClient(client *Client) {
    delete(clients.byEmail, client.email)
    delete(clients.byID, client.id)
}
