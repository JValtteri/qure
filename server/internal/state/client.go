package state

import (
    "fmt"
    "github.com/JValtteri/qure/server/internal/crypt"
    "github.com/JValtteri/qure/server/internal/utils"
)

const TEMP_CLIENT_AGE Epoch = 60*60*24*30    // max age in seconds

func NewClient(role string, email string, expire Epoch, sessionKey crypt.Key) (*Client, error) {
    if !uiniqueEmail(email) {
        return nil, fmt.Errorf("error: client email not unique")
    }

    id, err := createUniqueID(16, clients.byID)
    if err != nil {
        return nil, fmt.Errorf("error: Creating a new client\n%v", err) // Should not be possible (random byte generation)
    }
    client := createClient(id, expire, email, role)
    registerClient(client, sessionKey)

    return client, nil
}

func RemoveClient(client *Client) {
    delete(clients.byEmail, client.email)
    delete(clients.byID, client.id)
}


func uiniqueEmail(email string) bool {
    return unique(email, clients.byEmail)
}

func createClient(idBytes ID, expire Epoch, email string, role string) *Client {
    return &Client{
        id:         crypt.ID(idBytes),
        createdDt:  utils.EpochNow(),
        expiresDt:  expire,
        email:      email,
        phone:      "",
        role:       role,
        sessions:   make(map[crypt.Key]Session),
    }
}

func registerClient(client *Client, sessionKey crypt.Key) {
    clients.withLock(func() {
        clients.byID[client.id] = client;
        clients.bySession[sessionKey] = client;
        clients.byEmail[client.email] = client;
    })
}
