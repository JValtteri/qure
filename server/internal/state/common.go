package state

import (
    "fmt"
    "github.com/JValtteri/qure/server/internal/crypt"
)

type IP string      // IP address
type ID = crypt.ID

func unique[ K crypt.Key | crypt.ID | string , C *Client | Reservation ](key K, structure map[K]C) bool {
    clients.rLock()
    defer clients.rUnlock()
    _, notUnique := structure[key]
    return !notUnique
}

func createUniqueKey(length int, structure map[crypt.Key]*Client) (crypt.Key, error) {
    key, err := createUnique(length, structure)
    return key, err
}

func createUniqueID [ k crypt.ID | string ] (length int, structure map[k]*Client) (k, error) {
    key, err := createUnique(length, structure)
    return key, err
}

func createUniqueHumanReadableKey[ C *Client | Reservation ](length int, structure map[crypt.Key]C) (crypt.Key, error) {
    key, err := createUniqueHumanReadable(length, structure)
    return key, err
}

func createUniqueHumanReadableID[ C *Client | Reservation ](length int, structure map[crypt.ID]C) (crypt.ID, error) {
    key, err := createUniqueHumanReadable(length, structure)
    return key, err
}

func createUnique[ K crypt.Key | crypt.ID | string ](length int, structure map[K]*Client) (K, error) {
    var newId K = ""
    var err error
    var i int = 0
    var maxTries int = 5
    var keytype = K("")
    for i < maxTries {
        newId, err = crypt.CreateKey(&keytype, length)
        if unique(K(newId), structure) {
            return K(newId), err
        }
        i++
    }
    return K(newId), fmt.Errorf("failed to generate unique ID. Max tries (%v) exceeded \n%v", maxTries, err)
}

func createUniqueHumanReadable[ K crypt.Key | crypt.ID, C *Client | Reservation ](length int, structure map[K]C) (K, error) {
    var newId K = ""
    var err error
    var i int = 0
    var maxTries int = 5
    var keytype = K("")
    for i < maxTries {
        newId, err = crypt.CreateHumanReadableKey(&keytype, length)
        if unique(K(newId), structure) {
            return K(newId), err
        }
        i++
    }
    return K(newId), fmt.Errorf("failed to generate unique ID. Max tries (%v) exceeded \n%v", maxTries, err)
}
