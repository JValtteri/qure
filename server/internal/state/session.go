package state

import (
    "fmt"
    "github.com/JValtteri/qure/server/internal/utils"
    "github.com/JValtteri/qure/server/internal/crypt"
)


var MAX_SESSION_AGE Epoch = 60*60*24*30    // max age in seconds
var SESSION_KEY_LENGTH = 20

type Session struct {
    key         crypt.Key     // Session cookie
    expiresDt   Epoch
    fingerprint crypt.Hash
}

func ResetClients() {
    clients = Clients{
        byID:       make(map[crypt.ID]*Client),
        bySession:  make(map[crypt.Key]*Client),
        byEmail:    make(map[string]*Client),
    }
}

func ResumeSession(sessionKey crypt.Key, resumeFingerprint string) (*Client, error) {
    client, found := getClient(clients.bySession, sessionKey)
    if !found {
        return client, fmt.Errorf("no session matching key found: %v", sessionKey)
    }
    if !fingerprintMatch(resumeFingerprint, client, sessionKey) {
        RemoveSession(sessionKey)
        return client, fmt.Errorf("fingerprint doesn't match stored fingerprint")
    }
    err := cullExpired(&client.sessions)
    return client, err
}

func (client *Client) AddSession(role string, email string, temp bool, fingerprint crypt.Hash) (crypt.Key, error) {
    // Generate a unique session key
    sessionKey, err := createUniqueKey(SESSION_KEY_LENGTH, clients.bySession)
    if err != nil {
        return sessionKey, fmt.Errorf("error adding session %v", err)  // Should not be possible (random byte generation)
    }
    client.appendSession(sessionKey, fingerprint)
    return sessionKey, err
}


func (client *Client) appendSession(sessionKey crypt.Key, fingerprint crypt.Hash) {
    var session Session = Session{
        key:        sessionKey,
        expiresDt:  utils.EpochNow() + MAX_SESSION_AGE,
        fingerprint:         fingerprint,
    }
    clients.withLock(func() {
        client.sessions[sessionKey] = session
        clients.bySession[sessionKey] = client
    })
}

func fingerprintMatch(resumeFingerprint string, client *Client, sessionKey crypt.Key) bool {
    storedFingerprint := client.sessions[sessionKey].fingerprint
    return crypt.CompareToHash(resumeFingerprint, storedFingerprint)
}

func getClient(structure map[crypt.Key]*Client, key crypt.Key) (*Client, bool) {
    clients.rLock()
    defer clients.rUnlock()
    client, found := structure[key]
    return client, found
}

func cullExpired(sessions *map[crypt.Key]Session) error {
    var err error
    for key, session := range *sessions {
        now := utils.EpochNow()
        if now < session.expiresDt {
            continue
        }
        err = RemoveSession(key)
    }
    return err
}

func RemoveSession(sessionKey crypt.Key) error {
    clients.Lock()
    defer clients.Unlock()
    client, found := clients.bySession[sessionKey]
    if !found {
        return fmt.Errorf("session remove error: session not found")
    }
    // We trust that client.sessions[sessionKey] matches clients.bySession
    delete(client.sessions, sessionKey)    // Remove from client's sessions
    delete(clients.bySession, sessionKey)  // Remove from globas sessions
    if len(client.sessions) == 0 {
        RemoveClient(client)
    }
    return nil
}
