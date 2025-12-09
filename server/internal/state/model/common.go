package model

import (
	"fmt"
	"github.com/JValtteri/qure/server/internal/crypt"
)


type IP string			// IP address

func CreateUniqueKey(length int, structure map[crypt.Key]*Client) (crypt.Key, error) {
    key, err := createUnique(length, structure)
    return key, err
}

func CreateUniqueID [ k crypt.ID | string ] (length int, structure map[k]*Client) (k, error) {
    key, err := createUnique(length, structure)
    return key, err
}

func CreateUniqueHumanReadableKey[ C *Client | Reservation ](length int, structure map[crypt.Key]C) (crypt.Key, error) {
    key, err := createUniqueHumanReadable(length, structure)
    return key, err
}

func CreateUniqueHumanReadableID[ C *Client | Reservation ](length int, structure map[crypt.ID]C) (crypt.ID, error) {
    key, err := createUniqueHumanReadable(length, structure)
    return key, err
}

func Unique[ K crypt.Key | crypt.ID | string , C *Client | Reservation ](key K, structure map[K]C) bool {
	clientsLock.RLock()
	defer clientsLock.RUnlock()
	_, notUnique := structure[key]
	return !notUnique
}

func createUnique[ K crypt.Key | crypt.ID | string ](length int, structure map[K]*Client) (K, error) {
    var newId K = ""
    var err error
    var i int = 0
    var maxTries int = 5
    var keytype = K("")
    for i < maxTries {
        newId, err = crypt.CreateKey(&keytype, length)
		if Unique(K(newId), structure) {
            return K(newId), err
        }
        i++
    }
    return K(newId), fmt.Errorf("failed to generate unique ID. Max tries (%v) exceeded \n%v", maxTries, err)
}

func createUniqueHumanReadable[
	K crypt.Key | crypt.ID, C *Client | Reservation ](length int, structure map[K]C) (K, error) {
    var newId K = ""
    var err error
    var i int = 0
    var maxTries int = 5
    var keytype = K("")
    for i < maxTries {
        newId, err = crypt.CreateHumanReadableKey(&keytype, length)
		if Unique(K(newId), structure) {
            return K(newId), err
        }
        i++
    }
    return K(newId), fmt.Errorf("failed to generate unique ID. Max tries (%v) exceeded \n%v", maxTries, err)
}
