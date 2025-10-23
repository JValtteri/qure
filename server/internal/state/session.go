package state

import (
    "fmt"
    "github.com/JValtteri/qure/server/internal/utils"
    "github.com/JValtteri/qure/server/internal/crypt"
)

var MAX_SESSION_AGE Epoch = 60*60*24*30    // max age in seconds

type Session struct {
    key         crypt.Key     // Session cookie
    expiresDt   Epoch
    ip          IP      // IP should be stored hashed
}

func ResumeSession(sessionKey crypt.Key, resumeIP IP) error {
    client, found := getClient(clients.bySession, sessionKey)

    if !found {
        return fmt.Errorf("no session matching key found: %v", sessionKey)
    }

    if !isIPMatch(resumeIP, client, sessionKey) {
        removeSession(sessionKey)
        return fmt.Errorf("IP doesn't match stored IP")
    }

    return cullExpired(&client.sessions)
}

func AddSession(role string, email string, temp bool, ip IP) (crypt.Key, error) {
    // Generate a unique session key
    sessionKey, err := createUniqueKey(16, clients.bySession)
    if err != nil {
        return sessionKey, fmt.Errorf("error adding session %v", err)  // Should not be possible (random byte generation)
    }

    // Check if the email is already registered
    client, found := clients.byEmail[email]
    if found {
        appendSession(client, sessionKey, ip)
        return sessionKey, err
    }

    // Determine expiration time for the new client
    expire := calculateExpiration(temp)

    // Create a new client with the given details
    client, err = NewClient(role, email, expire, sessionKey)
    if err != nil {
        return crypt.Key("0"), err // Client not unique or such
    }

    appendSession(client, sessionKey, ip)
    return sessionKey, err
}


func appendSession(client *Client, sessionKey crypt.Key, ip IP) {
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

func calculateExpiration(temp bool) Epoch {
    var expire Epoch = 0
    if temp {
        expire = utils.EpochNow() + TEMP_CLIENT_AGE
    } else {
        expire-- // Set expire to maximum
    }
    return expire
}

func isIPMatch(resumeIP IP, client *Client, sessionKey crypt.Key) bool {
    storedIP := client.sessions[sessionKey].ip
    return storedIP == resumeIP
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
        err = removeSession(key)
    }
    return err
}

func removeSession(sessionKey crypt.Key) error {
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
