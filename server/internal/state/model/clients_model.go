package model

import (
	"fmt"
	"github.com/JValtteri/qure/server/internal/utils"
	"github.com/JValtteri/qure/server/internal/crypt"
)

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

func (c *Clients) RemoveReservation(id crypt.ID, reservation *Reservation) error {
	clientsLock.Lock()
	defer clientsLock.Unlock()
	client, ok := c.ByID[id]
	if !ok {
		err := fmt.Errorf("no client found with ID <%v>", id)
		return err
	}
	client.RemoveReservation(reservation)
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


func CreateClient(idBytes crypt.ID, expire utils.Epoch, email string, password crypt.Key, role string) *Client {
	return &Client{
		Id:			crypt.ID(idBytes),
		Password:	crypt.GenerateHash(password),
		CreatedDt:	utils.EpochNow(),
		ExpiresDt:	expire,
		Email:		email,
		Phone:		"",
		Role:		role,
		Sessions:	make(map[crypt.Key]Session),
	}
}
