package state

import (
    "fmt"
    "github.com/JValtteri/qure/server/internal/crypt"
    "github.com/JValtteri/qure/server/internal/utils"
)

const TEMP_CLIENT_AGE Epoch = 60*60*24*30    // max age in seconds

func NewClient(role string, email string, password crypt.Key, temp bool, sessionKey crypt.Key) (*Client, error) {
    var client *Client
    var err error
    var expire Epoch = calculateExpiration(temp)
    if temp {
        sessionKey, err = createUniqueKey(SESSION_KEY_LENGTH, clients.bySession)
        if err != nil {
            return client, fmt.Errorf("Error creating key: %v\n", err)
        }
        client, err = createTempClient(expire, email)
    } else {
        if sessionKey == "" {
            return client, fmt.Errorf("missing session key for NewClient")
        }
        client, err = createNormalClient(email, expire, password, role, sessionKey)
    }
    if err != nil {
        return client, fmt.Errorf("error creating client: %v\n", err)
    }
    registerClient(client, sessionKey)
    return client, err
}

func RemoveClient(client *Client) {
    delete(clients.byEmail, client.email)
    delete(clients.byID, client.id)
}

func createNormalClient(email string, expire Epoch, password crypt.Key, role string, sessionKey crypt.Key) (*Client, error) {
    fmt.Println("FOO")
    if !uniqueEmail(email) {
        fmt.Println("NOQ")
        return nil, fmt.Errorf("error: client email not unique")
    }

    id, err := createUniqueID(16, clients.byID)
    if err != nil {
        return nil, fmt.Errorf("error: Creating a new client\n%v", err) // Should not be possible (random byte generation)
    }
    client := createClient(id, expire, email, password, role)
    return client, nil
}

func createTempClient(expire Epoch, email string) (*Client, error) {
    id, err := createUniqueID(16, clients.byID)
    if err != nil {
        return nil, fmt.Errorf("error: Creating a new ID\n%v", err) // Should not be possible (random byte generation)
    }
    pseudoEmail, err := createUniqueID(16, clients.byEmail)
    if err != nil {
        return nil, fmt.Errorf("error: Creating a new ID\n%v", err) // Should not be possible (random byte generation)
    }
    password := crypt.Key(id)
    client := createClient(id, expire, pseudoEmail, password, "temp")
    client.email = email
    return client, nil
}


func uniqueEmail(email string) bool {
    return unique(email, clients.byEmail)
}

func createClient(idBytes ID, expire Epoch, email string, password crypt.Key, role string) *Client {
    return &Client{
        id:         crypt.ID(idBytes),
        password:   crypt.Hash(password),
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

func calculateExpiration(temp bool) Epoch {
    var expire Epoch = 0
    if temp {
        expire = utils.EpochNow() + TEMP_CLIENT_AGE
    } else {
        expire-- // Set expire to maximum
    }
    return expire
}
