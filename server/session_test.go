package main

import (
    "testing"
)

func TestCreateRandomBytes(t *testing.T) {
    const expect int = 16
    var got []byte
    var err error
    got, err = RandomBytes(expect)
    if err != nil {
        t.Errorf("Expected: %v, Got: %v\n", "ok", err)
    }
    if len(got) < expect {
        t.Errorf("Expected: %v, Got: %v\n", expect, len(got))
    }
}

func TestCreateRandomInts(t *testing.T) {
    const expect int = 16
    var got []int
    var err error
    got, err = RandomInts(expect, 25)
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
    var got string
    var err error
    got, err = RandomChars(expect)
    if err != nil {
        t.Errorf("Expected: %v, Got: %v\n", "ok", err)
    }
    if len(got) < expect {
        t.Errorf("Expected: %v, Got: %v\n", expect, len(got))
    }
}

func TestCreateUniqueID(t *testing.T) {
    const expect int = 16
    var got Key
    var err error
    got, err = createUniqueID(expect)
    if err != nil {
        t.Errorf("Expected: %v, Got: %v\n", "ok", err)
    }
    if len(got) != expect {
        t.Errorf("Expected: %v, Got: %v\n", expect, len(got))
    }
}

func TestCreateConflictingID(t *testing.T) {
    const expect int = 15
    var got Key
    var err error
    i := 0
    // Create lots of short ID's (1 char) to force a collision
    for i < 100000 {
        i++
        var client Client
        got, err = createUniqueID(1)
        clients[got] = client
        if err != nil {
            return
        }
    }
    t.Errorf("Expected: %v, Got: %v\n", "Duplicate Error", err)
}

func TestCreateClient(t *testing.T) {
    const expect int = 1 // 15
    var got Key
    var err error
    got, err = createUniqueID(16)
    if err != nil {
        t.Errorf("Expected: %v, Got: %v\n", "ok", err)
    }
    if len(got) < expect {
        t.Errorf("Expected: %v, Got: %v\n", expect, len(got))
    }
}

/*
func TestCreateSession(t *testing.T) {
    const input string = ""
    const expect string = ""
    var got string = CreateSession(input)
    if got != expect {
        t.Errorf("Expected: %v, Got: %v\n", expect, got)
    }
}
*/
