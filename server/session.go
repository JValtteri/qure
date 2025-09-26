package main

import (
    "crypto/rand"
    "errors"
    "time"
    "fmt"
    "log"
)

var clients map[Key]Client = make(map[Key]Client)

type Key string   // Session Key
//type ID string    // Static ID

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

func NewClient(role string, email string, phone string, expire uint) {
    var client Client
    id, err := createUniqueID(16)
    if err != nil {
        log.Printf("Error: Creating a new client\n%v", err)
    }
    client.id = id
    client.createdDt = uint(time.Now().Unix())
    client.expiresDt = expire
    client.email = email
    client.phone = phone
    client.role = role
    client.sessionsKeys = make(map[Key]Epoch)
    // client.reservations = []  // make sure it's empty
}

func CreateSession(role string, email string) Client {
    // TODO
    var client Client
    return client
}

func addSession(role string, email string) bool {
    //newSession := createSession(role, email)
    // TODO
    return false
}

func createUniqueID(length int) (Key, error) {
    var newId string = ""
    var err error
    var i int = 0
    var maxTries int = 5
    for i < maxTries {
        newId, err = RandomChars(length)
        if unique(Key(newId)) {
            return Key(newId), err
        }
        i++
    }
    return Key(newId), errors.New(fmt.Sprintf("Error: Failed to generate unique ID. Max tries (%v) exceeded \n%v\n", maxTries, err))
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
    ints := make([]int, length, length)
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

func unique(id Key) bool {
    _, notUnique := clients[id]
    return !notUnique
}
