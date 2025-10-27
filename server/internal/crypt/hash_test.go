package crypt

import (
    "testing"
    "strings"
)


func TestHash(t *testing.T) {
    input := Key("password")
    length := 32
    got, err := GenerateHash(input)
    fragments := strings.Split(string(got), ",p=1$")
    if err != nil {
        t.Errorf("Expected: %v, Got: %v\n", "ok", err)
    }
    if len(fragments[1]) < length {
        t.Errorf("Expected: %v < %v\n", length, len(fragments[1]))
    }
}

func TestCompareHash(t *testing.T) {
    password := Key("password")
    expect := true
    hash, _ := GenerateHash(password)
    got, err := CompareToHash(password, hash)
    if err != nil {
        t.Errorf("Expected: %v, Got: %v\n", "ok", err)
    }
    if got == false {
        t.Errorf("Expected: %v, Got: %v\n", expect, got)
    }
}

func TestCompareBadHash(t *testing.T) {
    password := Key("password")
    wrongpass := Key("passw0rd")
    expect := false
    hash, _ := GenerateHash(password)
    badhash, _ := GenerateHash(wrongpass)
    t.Logf("\nHash: %v\n", hash)
    t.Logf("\nHash: %v\n", badhash)
    got, err := CompareToHash(password, badhash)
    if err != nil {
        t.Errorf("Expected: %v, Got: %v\n", "ok", err)
    }
    if got != false {
        t.Errorf("Expected: %v, Got: %v\n", expect, got)
    }
}
