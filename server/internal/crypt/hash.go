package crypt

import "github.com/alexedwards/argon2id"


type Hash string

var parms = &argon2id.Params{
    Memory: 19*1024,
    Iterations: 2,
    Parallelism: 1,
    SaltLength: 16,
    KeyLength: 32,
}

func GenerateHash [ K Key | string ](password K) (Hash, error) {
    key := string(password)
    hash, err := argon2id.CreateHash(key, parms)
    return Hash(hash), err
}

func CompareToHash [ K Key | string ](password K, hash Hash) (bool, error) {
    match, err := argon2id.ComparePasswordAndHash(string(password), string(hash))
    return match, err
}
