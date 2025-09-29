package main

import (
    "crypto/rand"
    "fmt"
    "log"
    "strings"
    "time"
)

const MAX_SESSION_AGE uint = 60*60*24*30    // max age in seconds
const TEMP_CLIENT_AGE uint = 60*60*24*30    // max age in seconds

type Clients interface {
    Client | *Client
}

var clients map[Key]Client            = make(map[Key]Client)       // by client ID
var clientsbySession map[Key]*Client  = make(map[Key]*Client)      // by session key
var clientsbyEmail map[string]*Client = make(map[string]*Client)   // by session key

type Key string     // Session Key
//type ID string    // Static ID
type IP string      // IP address

type Session struct {
    key         Key
    expiresDt   uint
    ip          IP
}

type Client struct {
    id           Key
    createdDt    uint     // Unix timestamp
    expiresDt    uint     // Unix timestamp, 0 = expire now, 0-- = keep indefinately
    email        string
    phone        string
    role         string
    sessionsKeys map[Key]Epoch
    reservations []*Reservation
}

func ResumeSession(sessionKey Key) {
    client, found := clientsbySession[sessionKey]
    if found {
        cullExpired(&client.sessionsKeys)
    }
}

func NewClient(role string, email string, expire uint, sessionKey Key) bool {
    uiniqueEmail := unique(email, clientsbyEmail)
    if !uiniqueEmail {
        log.Printf("Error: client email not unique\n")
        return false
    } else {
        var client Client
        id, err := createUniqueID(16, clients)
        if err != nil {
            log.Printf("Error: Creating a new client\n%v\n", err)
        }
        client.id = id
        client.createdDt = uint(time.Now().Unix())
        client.expiresDt = expire
        client.email = email
        client.phone = ""
        client.role = role
        client.sessionsKeys = make(map[Key]Epoch)
        // client.reservations = []  // make sure it's empty

        // Add client to
        clients[id] = client
        clientsbySession[sessionKey] = &client
        clientsbyEmail[email] = &client
        return true
    }
}

func appendSession(client *Client, sessionKey Key, role string, email string) {
    client.sessionsKeys[sessionKey] = Epoch(uint(time.Now().Unix()) + MAX_SESSION_AGE)
}

func AddSession(role string, email string, temp bool) bool {
    sessionKey, err := createUniqueID(16, clients)
    if err != nil {
        log.Printf("Error adding session %v\n", err)
        return false
    }
    // Is the email registered alredy?
    client, found := clientsbyEmail[email]
    var expire uint = 0
    if found {
        appendSession(client, sessionKey, role, email)
    }
    if temp {
        expire = uint(time.Now().Unix()) + TEMP_CLIENT_AGE
    } else {
        expire-- // Set expire to maximum
    }
    ok := NewClient(role, email, expire, sessionKey)
    return ok
}

func createHumanReadableId(length int) (Key, error) {
    var newID Key
    var id string
    var err error
    maxTries := 5
    i := 0
    for i < maxTries {
        i++
        newID, err = createUniqueID(length*2, clientsbySession)
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

func createUniqueID[V Clients](length int, structure map[Key]V) (Key, error) {
    var newId string = ""
    var err error
    var i int = 0
    var maxTries int = 5
    for i < maxTries {
        newId, err = RandomChars(length)
        if unique(Key(newId), structure) {
            return Key(newId), err
        }
        i++
    }
    return Key(newId), fmt.Errorf("failed to generate unique ID. Max tries (%v) exceeded \n%v", maxTries, err)
}

// Returns a string containing random chars from [A..Z,a..z,0..9]
func RandomChars(length int) (string, error) {
    ints, err := RandomInts(length, 62)
    for i, v := range ints {
        ints[i] = asciiOffset(v)
    }
    bytes := itob(ints)
    return string(bytes), err
}

func asciiOffset(v int) int {
        if v < 26 {
            return v+65
        } else if v < 52 {
            return v-26+97
        } else {
            return v-52+48
        }
}

func RandomInts(length int, base int) ([]int, error) {
    ints := make([]int, length)
    bytes, err := RandomBytes(length)
    for i, b := range bytes {
        ints[i] = int(b) % base
    }
    return ints, err
}

func RandomBytes(length int) ([]byte, error) {
    buffer := make([]byte, length)
    _, err := rand.Read(buffer)
    if err != nil {
        fmt.Println("error:", err)
        return buffer, err
    }
    return buffer, nil
}

func unique[ V Clients, K Key | string ](id K, structure map[K]V) bool {
    _, notUnique := structure[id]
    return !notUnique
}

func cullExpired(sessions *map[Key]Epoch) {
    for key, expire := range *sessions {
        now := Epoch(time.Now().Unix())
        if now < expire {
            continue
        }
        delete(*sessions, key)          // Remove from Clients sessions
        client := clientsbySession[key]
        delete(clientsbySession, key)   // Remove from globas sessions
        if len(client.sessionsKeys) == 0 {
            RemoveClient(client)
        }
    }
}

func RemoveClient(client *Client) {
    delete(clientsbyEmail, client.email)
    delete(clients, client.id)
}

