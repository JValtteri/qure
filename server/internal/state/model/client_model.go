package model

import (
	"fmt"
	"sync"
	"github.com/JValtteri/qure/server/internal/utils"
	"github.com/JValtteri/qure/server/internal/crypt"
)


var MAX_SESSION_AGE utils.Epoch = 60*60*24*30	// max age in seconds

var clientsLock sync.RWMutex = sync.RWMutex{}

type Client struct {
	Id				crypt.ID
	Password		crypt.Hash
	CreatedDt		utils.Epoch		// Unix timestamp
	ExpiresDt		utils.Epoch		// Unix timestamp, 0 = expire now, 0-- = keep indefinately
	Email			string
	Phone			string
	Role			string
	Sessions		map[crypt.Key]Session
	Reservations	[]*Reservation
}

func (t *Client) GetPasswordHash() crypt.Hash {
	clientsLock.RLock()
	defer clientsLock.RUnlock()
	return t.Password
}

func (t *Client) GetEmail() string {
	clientsLock.RLock()
	defer clientsLock.RUnlock()
	return t.Email
}

func (t *Client) GetRole() string {
	clientsLock.RLock()
	defer clientsLock.RUnlock()
	return t.Role
}

func (t *Client) IsAdmin() bool {
	clientsLock.RLock()
	defer clientsLock.RUnlock()
	return t.Role == "admin"
}

func (t *Client) GetReservations() []*Reservation {
	clientsLock.RLock()
	defer clientsLock.RUnlock()
	return t.Reservations
}

func (t *Client) AddReservation(res *Reservation) {
	t.Reservations = append(t.Reservations, res)
}

func (client *Client) AddSession(role string, email string, temp bool, fingerprint crypt.Hash, clients *Clients) (crypt.Key, error) {
	// Generate a unique session key
	sessionKey, err := CreateUniqueKey(SESSION_KEY_LENGTH, clients.BySession)
	if err != nil {
		return sessionKey, fmt.Errorf("error adding session %v", err)	// Should not be possible (random byte generation)
	}
	client.appendSession(sessionKey, fingerprint, clients)
	return sessionKey, err
}


func (client *Client) appendSession(sessionKey crypt.Key, fingerprint crypt.Hash, clients *Clients) {
	var session Session = Session{
		Key:			sessionKey,
		ExpiresDt:		utils.EpochNow() + MAX_SESSION_AGE,
		Fingerprint:	fingerprint,
	}
	clientsLock.Lock()
	defer clientsLock.Unlock()
	client.Sessions[sessionKey] = session
	clients.BySession[sessionKey] = client
}

type Clients struct {
	ByID		map[crypt.ID]*Client		// by client ID
	BySession	map[crypt.Key]*Client		// by session key
	ByEmail		map[string]*Client			// by session key
}

func (c *Clients) AddReservation(id crypt.ID, reservation *Reservation) error {
	clientsLock.Lock()
	defer clientsLock.Unlock()
	client, ok := c.ByID[id]
	if !ok {
		err := fmt.Errorf("no client found with ID <%v>", id)
		return err
	}
	client.AddReservation(reservation)
	return nil
}

func (c *Clients) GetClientBySession(sessionKey crypt.Key) (*Client, bool) {
	clientsLock.RLock()
	defer clientsLock.RUnlock()
	client, found := c.BySession[sessionKey]
	return client, found
}


func (c *Clients) RLock() {
	clientsLock.RLock()
}

func (c *Clients) RUnlock() {
	clientsLock.RUnlock()
}

func (c *Clients) Lock() {
	clientsLock.Lock()
}

func (c *Clients) Unlock() {
	clientsLock.Unlock()
}
