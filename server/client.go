package main

import (
    "strings"
    "sync"
    "fmt"
)

const TEMP_CLIENT_AGE Epoch = 60*60*24*30    // max age in seconds
type ID string      // Static ID

type ClientLike interface {
    Client | *Client
}

type Client struct {
    id           ID        // Effectively password. Should be stored hashed
    createdDt    Epoch     // Unix timestamp
    expiresDt    Epoch     // Unix timestamp, 0 = expire now, 0-- = keep indefinately
    email        string
    phone        string
    role         string
    sessions map[Key]Session
    reservations []*Reservation
}

func (t *Client) AddReservation(res *Reservation) {
    t.reservations = append(t.reservations, res)
}

var clients Clients = Clients{
    byID:       make(map[ID]*Client),
    bySession:  make(map[Key]*Client),
    byEmail:    make(map[string]*Client),
}
type Clients struct {
    mu          sync.RWMutex
    byID        map[ID]*Client        // by client ID
    bySession   map[Key]*Client       // by session key
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

func (c *Clients) getClient(sessionKey Key) (*Client, bool) {
    c.mu.RLock()
    defer c.mu.RUnlock()
    client, found := clients.bySession[sessionKey]
    return client, found
}

func (c *Clients) AddReservation(id ID, reservation *Reservation) {
    c.withLock(func() {
        client := c.byID[id]
        client.AddReservation(reservation)
        //c.raw[id] = client
    })
}

func (c *Clients) GetReservations(id ID) []*Reservation {
    c. rLock()
    defer c.rUnlock()
    return c.byID[id].reservations
}

func NewClient(role string, email string, expire Epoch, sessionKey Key) (*Client, error) {
    var client Client
    uiniqueEmail := unique(email, clients.byEmail)
    if !uiniqueEmail {
        return &client, fmt.Errorf("error: client email not unique")
    } else {
        kId, err := createUniqueID(16, clients.byID)
        id := ID(kId)
        if err != nil {
            return &client, fmt.Errorf("error: Creating a new client\n%v", err) // Should not be possible (random byte generation)
        }
        client.id = ID(id)
        client.createdDt = EpochNow()
        client.expiresDt = expire
        client.email = email
        client.phone = ""
        client.role = role
        client.sessions = make(map[Key]Session)
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

func createHumanReadableId(length int) (Key, error) {
    var newID Key
    var id string
    var err error
    maxTries := 5
    i := 0
    for i < maxTries {
        i++
        newID, err = createUniqueID(length*2, clients.byID)
        id = string(newID)
        // Remove look-alike characters
        id = strings.ReplaceAll(string(id), "O", "")
        id = strings.ReplaceAll(string(id), "0", "")
        id = strings.ReplaceAll(string(id), "Q", "")
        id = strings.ReplaceAll(string(id), "I", "")
        id = strings.ReplaceAll(string(id), "l", "")
        id = strings.ReplaceAll(string(id), "1", "")
        if len(id) > length {
            return Key(id[:length]), err
        }
    }
    return newID, fmt.Errorf("failed to generate unique ID. Max tries (%v) exceeded \n%v", maxTries, err)
}

func createUniqueID[V ClientLike, K Key | ID ](length int, structure map[K]V) (Key, error) {
    var newId string = ""
    var err error
    var i int = 0
    var maxTries int = 5
    for i < maxTries {
        newId, err = RandomChars(length)
        if unique(K(newId), structure) {
            return Key(newId), err
        }
        i++
    }
    return Key(newId), fmt.Errorf("failed to generate unique ID. Max tries (%v) exceeded \n%v", maxTries, err)
}

func unique[ V ClientLike, K Key | ID | string ](id K, structure map[K]V) bool {
    clients.rLock()
    defer clients.rUnlock()
    _, notUnique := structure[id]
    return !notUnique
}
