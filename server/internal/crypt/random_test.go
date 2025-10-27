package crypt

import (
    "testing"
)


func TestCreateRandomBytes(t *testing.T) {
    const expect int = 16
    got, err := randomBytes(expect)
    if err != nil {
        t.Errorf("Expected: %v, Got: %v\n", "ok", err)
    }
    if len(got) < expect {
        t.Errorf("Expected: %v, Got: %v\n", expect, len(got))
    }
}

func TestCreateRandomInts(t *testing.T) {
    const expect int = 16
    got, err := randomInts(expect, 25)
    if err != nil {
        t.Errorf("Expected: %v, Got: %v\n", "ok", err)
    }
    if len(got) < expect {
        t.Errorf("Expected: %v, Got: %v\n", expect, len(got))
    }
}

func TestAsciiOffser(t *testing.T) {
    var expect byte = 0
    expect = 65 // A
    input := 0
    got := asciiOffset(input)
    if byte(got) != expect {
        t.Errorf("Expected: %v, Got: %v %c\n", expect, byte(got), byte(got))
    }
    expect = 90 // Z
    input = 25
    got = asciiOffset(input)
    if byte(got) != expect {
        t.Errorf("Expected: %v, Got: %v %c\n", expect, byte(got), byte(got))
    }
    expect = 97 // a
    input = 26
    got = asciiOffset(input)
    if byte(got) != expect {
        t.Errorf("Expected: %v, Got: %v %c\n", expect, byte(got), byte(got))
    }
    expect = 48 // 0
    input = 52
    got = asciiOffset(input)
    if byte(got) != expect {
        t.Errorf("Expected: %v, Got: %v %c\n", expect, byte(got), byte(got))
    }
}

func TestCreateRandomChars(t *testing.T) {
    const expect int = 62
    got, err := randomChars(expect)
    if err != nil {
        t.Errorf("Expected: %v, Got: %v\n", "ok", err)
    }
    if len(got) < expect {
        t.Errorf("Expected: %v, Got: %v\n", expect, len(got))
    }
}

func TestCreateKey(t *testing.T) {
    const expect int = 16
    var got Key
    var err error
    got, err = CreateKey(&got, expect)
    if err != nil {
        t.Errorf("Expected: %v, Got: %v\n", "ok", err)
    }
    if len(got) != expect {
        t.Errorf("Expected: %v, Got: %v\n", expect, len(got))
    }
}

func TestCreateHumanreadableID(t *testing.T) {
    const expect int = 16
    var got ID
    var err error
    // Create lots of short ID's (1 char) to force a collision
    got, err = CreateHumanReadableKey(&got, expect)
    if err != nil {
        t.Errorf("Expected: %v, Got: %v\n", "ok", err)
    }
    if len(got) != expect {
        t.Errorf("Expected: %v, Got: %v\n", expect, len(got))
    }
}
