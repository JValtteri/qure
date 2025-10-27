package crypt

import(
    "fmt"
    "strings"
    "crypto/rand"
    "github.com/JValtteri/qure/server/internal/utils"
)


type Key string     // Session Key
type ID string      // Static ID

func CreateHumanReadableKey[ K Key | ID ](keytype *K, length int) (K, error) {
    var key string
    var err error
    maxTries := 5
    i := 0
    for i < maxTries {
        i++
        newKey, err := CreateKey(keytype, length*2)
        key = string(newKey)
        // Remove look-alike characters
        key = strings.ReplaceAll(string(key), "O", "")
        key = strings.ReplaceAll(string(key), "0", "")
        key = strings.ReplaceAll(string(key), "Q", "")
        key = strings.ReplaceAll(string(key), "I", "")
        key = strings.ReplaceAll(string(key), "l", "")
        key = strings.ReplaceAll(string(key), "1", "")
        if len(key) > length {
            return K(key[:length]), err
        }
    }
    return *keytype, fmt.Errorf("failed to generate unique ID. Max tries (%v) exceeded \n%v", maxTries, err)
}

func CreateKey[ K Key | ID | string ](newKey *K, length int) (K, error) {
    str, err := randomChars(length)
    if err != nil {
        return *newKey, fmt.Errorf("error Creating a key: %v", err)
    }
    key := K(str)
    return key, nil
}

// Returns a string containing random chars from [A..Z,a..z,0..9]
func randomChars(length int) (string, error) {
    ints, err := randomInts(length, 62)
    for i, v := range ints {
        ints[i] = asciiOffset(v)
    }
    bytes := utils.ItoB(ints)
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

func randomInts(length int, base int) ([]int, error) {
    ints := make([]int, length)
    bytes, err := randomBytes(length)
    for i, b := range bytes {
        ints[i] = int(b) % base
    }
    return ints, err
}

func randomBytes(length int) ([]byte, error) {
    buffer := make([]byte, length)
    _, err := rand.Read(buffer)
    if err != nil {
        return buffer, fmt.Errorf("error generating random bytes: %v", err) // Should not be possible (random byte generation)
    }
    return buffer, nil
}
