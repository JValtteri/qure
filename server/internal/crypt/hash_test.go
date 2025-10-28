package crypt

import (
    "testing"
    "strings"
)


func TestHash(t *testing.T) {
    input := Key("password")
    length := 32
    got := GenerateHash(input)
    fragments := strings.Split(string(got), ",p=1$")
    if len(fragments[1]) < length {
        t.Errorf("Expected: %v < %v\n", length, len(fragments[1]))
    }
}

func TestCompareHash(t *testing.T) {
    password := Key("password")
    expect := true
    hash := GenerateHash(password)
    got := CompareToHash(password, hash)
    if got == false {
        t.Errorf("Expected: %v, Got: %v\n", expect, got)
    }
}

func TestCompareBadHash(t *testing.T) {
    password := Key("password")
    wrongpass := Key("passw0rd")
    expect := false
    hash := GenerateHash(password)
    badhash := GenerateHash(wrongpass)
    t.Logf("\nHash: %v\n", hash)
    t.Logf("\nHash: %v\n", badhash)
    got := CompareToHash(password, badhash)
    if got != false {
        t.Errorf("Expected: %v, Got: %v\n", expect, got)
    }
}
