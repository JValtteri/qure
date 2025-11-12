package crypt

import (
    "log"
    "github.com/alexedwards/argon2id"
)


type Hash string

var parms = &argon2id.Params{
    Memory: 19*1024,
    Iterations: 2,
    Parallelism: 1,
    SaltLength: 16,
    KeyLength: 32,
}

func GenerateHash [ K Key | string ](password K) Hash {
    key := string(password)
    hash, err := argon2id.CreateHash(key, parms)
    if err != nil {
        log.Printf("Error in hash generation: %v\n", err)
    }
    return Hash(hash)
}

func CompareToHash [ K Key | string ](key K, hash Hash) bool {
    match, err := argon2id.ComparePasswordAndHash(string(key), string(hash))
    if err != nil {
        log.Printf("Error in hash comparison: %v\n", err)
    }
    return match
}
