package model

import (
	"fmt"
	"sync"
	"github.com/JValtteri/qure/server/internal/utils"
	"github.com/JValtteri/qure/server/internal/crypt"
	c "github.com/JValtteri/qure/server/internal/config"
)


var clientsLock sync.RWMutex = sync.RWMutex{}

type Client struct {
	Id				crypt.ID
	Password		crypt.Hash
	CreatedDt		utils.Epoch		// Unix timestamp
	ExpiresDt		utils.Epoch		// Unix timestamp, 0 = expire now, 0-- = keep indefinately
	IsTemporary		bool
	Email			string
	Phone			string
	Role			string
	Sessions		map[crypt.Key]Session
	Reservations	[]*Reservation
}

func (c *Client) GetPasswordHash() crypt.Hash {
	clientsLock.RLock()
	defer clientsLock.RUnlock()
	return c.Password
}

func (c *Client) GetEmail() string {
	clientsLock.RLock()
	defer clientsLock.RUnlock()
	return c.Email
}

func (c *Client) GetRole() string {
	clientsLock.RLock()
	defer clientsLock.RUnlock()
	return c.Role
}

func (c *Client) IsAdmin() bool {
	clientsLock.RLock()
	defer clientsLock.RUnlock()
	return c.Role == "admin"
}

func (c *Client) GetReservations() []*Reservation {
	clientsLock.RLock()
	defer clientsLock.RUnlock()
	return c.Reservations
}

func (c *Client) AddReservation(res *Reservation) {
	c.Reservations = append(c.Reservations, res)
}

func (c *Client) RemoveReservation(res *Reservation) {
	// Filters from list all instances of targetID
	var filtered = []*Reservation{}
	for _, value := range c.Reservations {
		if value.Id != res.Id {
			filtered = append(filtered, value)
		}
	}
	c.Reservations = filtered
}

func (client *Client) AddSession(
	role string,	email string,	temp bool,
	fingerprint		crypt.Hash,		clients *Clients,
) (crypt.Key, error) {
	// Generate a unique session key
	sessionKey, err := CreateUniqueKey(c.CONFIG.SESSION_KEY_LENGTH, clients.BySession)
	if err != nil {
		return sessionKey, fmt.Errorf("error adding session %v", err)	// Should not be possible (random byte generation)
	}
	client.appendSession(sessionKey, fingerprint, clients)
	return sessionKey, err
}

func (c *Client) ClearSessions(){
	c.Sessions = make(map[crypt.Key]Session)
}

func (client *Client) appendSession(sessionKey crypt.Key, fingerprint crypt.Hash, clients *Clients) {
	var session Session
	if c.CONFIG.EXTRA_STRICT_SESSIONS {
		session = Session{
			Key:			sessionKey,
			ExpiresDt:		utils.EpochNow() + c.CONFIG.MAX_SESSION_AGE,
			Fingerprint:	fingerprint,
		}
	} else {
		session = Session{
			Key:			crypt.Key(fmt.Sprintf("%s%s", sessionKey, fingerprint)),
			ExpiresDt:		utils.EpochNow() + c.CONFIG.MAX_SESSION_AGE,
			Fingerprint:	fingerprint,
		}
	}
	clientsLock.Lock()
	defer clientsLock.Unlock()
	client.Sessions[sessionKey] = session
	clients.BySession[sessionKey] = client
}
