package state

import (
    "fmt"
    "github.com/JValtteri/qure/server/internal/utils"
)

var MAX_SESSION_AGE Epoch = 60*60*24*30    // max age in seconds

type Key string     // Session Key
type IP string      // IP address

type Session struct {
    key         Key     // Session cookie
    expiresDt   Epoch
    ip          IP      // IP should be stored hashed
}

func ResumeSession(sessionKey Key, resumeIP IP) error {
    client, found := clients.withRLock(func() (*Client, bool) {
        client, found := clients.bySession[sessionKey]
        return client, found
    })
    if !found {
        return fmt.Errorf("no session matching key found: %v", sessionKey)
    }
    storedIP := client.sessions[sessionKey].ip
    if storedIP != resumeIP {
        removeSession(sessionKey)
        return fmt.Errorf("IP doesn't match stored IP")
    }
    err := cullExpired(&client.sessions)
    return err
}

func AddSession(role string, email string, temp bool, ip IP) (Key, error) {
    sessionKey, err := createUniqueID(16, clients.bySession)
    if err != nil {
        return sessionKey, fmt.Errorf("error adding session %v", err)  // Should not be possible (random byte generation)
    }

    // Is the email registered alredy?
    client, found := clients.byEmail[email]
    if found {
        appendSession(client, sessionKey, ip)
        return sessionKey, err
    }

    var expire Epoch = 0
    if temp {
        expire = utils.EpochNow() + TEMP_CLIENT_AGE
    } else {
        expire-- // Set expire to maximum
    }

    client, err = NewClient(role, email, expire, sessionKey)
    appendSession(client, sessionKey, ip)
    if err != nil {
        return Key("0"), err // Client not unique or such
    }
    return sessionKey, err
}

func appendSession(client *Client, sessionKey Key, ip IP) {
    var session Session = Session{
        key:        sessionKey,
        expiresDt:  utils.EpochNow() + MAX_SESSION_AGE,
        ip:         ip,
    }
    clients.withLock(func() {
        client.sessions[sessionKey] = session
        clients.bySession[sessionKey] = client
    })
}

func cullExpired(sessions *map[Key]Session) error {
    var err error
    for key, session := range *sessions {
        now := utils.EpochNow()
        if now < session.expiresDt {
            continue
        }
        err = removeSession(key)
    }
    return err
}

func removeSession(sessionKey Key) error {
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
