package main

import (
    "time"
    "fmt"
)

const MAX_SESSION_AGE uint = 60*60*24*30    // max age in seconds


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
    storedIP := client.sessionsKeys[sessionKey].ip
    if storedIP != resumeIP {
        // TODO: expire session
        return fmt.Errorf("IP doesn't match stored IP")
    }
    //cullExpired(&client.sessionsKeys)
    return nil
}

func AddSession(role string, email string, temp bool, ip IP) (Key, error) {
    sessionKey, err := createUniqueID(16, clients.bySession)
    if err != nil {
        return sessionKey, fmt.Errorf("error adding session %v", err)
    }

    // Is the email registered alredy?
    client, found := clients.byEmail[email]
    if found {
        appendSession(client, sessionKey, ip)
        return sessionKey, err
    }

    var expire Epoch = 0
    if temp {
        expire = Epoch(uint(time.Now().Unix()) + TEMP_CLIENT_AGE)
    } else {
        expire-- // Set expire to maximum
    }

    client, err = NewClient(role, email, expire, sessionKey)
    appendSession(client, sessionKey, ip)
    if err != nil {
        return Key("0"), err
    }
    return sessionKey, err
}

func appendSession(client *Client, sessionKey Key, ip IP) {
    var session Session = Session{
        key:        sessionKey,
        expiresDt:  Epoch(uint(time.Now().Unix()) + MAX_SESSION_AGE),
        ip:         ip,
    }
    clients.withLock(func() {
        client.sessionsKeys[sessionKey] = session
    })
}

func cullExpired(sessions *map[Key]Epoch) {
    for key, expire := range *sessions {
        now := Epoch(time.Now().Unix())
        if now < expire {
            continue
        }
        delete(*sessions, key)          // Remove from Clients sessions
        client := clients.bySession[key]
        delete(clients.bySession, key)   // Remove from globas sessions
        if len(client.sessionsKeys) == 0 {
            RemoveClient(client)
        }
    }
}

