package main

import (
    "testing"
    "log"
    "os"
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
    got, err = createUniqueID(expect, clients)
    if err != nil {
        t.Errorf("Expected: %v, Got: %v\n", "ok", err)
    }
    if len(got) != expect {
        t.Errorf("Expected: %v, Got: %v\n", expect, len(got))
    }
}

func TestCreateConflictingID(t *testing.T) {
    var got Key
    var err error
    i := 0
    // Create lots of short ID's (1 char) to force a collision
    for i < 1000000 {
        i++
        var client Client
        got, err = createUniqueID(1, clients)
        clients[got] = client
        if err != nil {
            t.Logf("Tries before conflict %v", i)
            return
        }
    }
    t.Errorf("Expected: %v, Got: %v\n", "Duplicate Error", err)
}

func TestHumanReadableId(t *testing.T) {
    const expect int = 16 // 15
    var got Key
    var err error
    got, err = createHumanReadableId(expect)
    if err != nil {
        t.Errorf("Expected: %v, Got: %v\n", "ok", err)
    }
    if len(got) != expect {
        t.Errorf("Expected: %v, Got: %v\n", expect, len(got))
    }
}

func TestHumanReadableConflictId(t *testing.T) {
    var got Key
    var err error
    i := 0
    // Create lots of short ID's (1 char) to force a collision
    for i < 1000000 {
        i++
        var client Client
        got, err = createHumanReadableId(1)
        clientsbySession[got] = &client
        if err != nil {
            t.Logf("Tries before conflict %v", i)
            return
        }
    }
    t.Errorf("Expected: %v, Got: %v\n", "Duplicate Error", err)
}

func TestCreateClient(t *testing.T) {
    const expect bool = true
    got := NewClient("test", "example@example.com", 0, Key("000"))
    if got != expect {
        t.Errorf("Expected: %v, Got: %v\n", expect, got)
    }
}

func TestCreateDuplicateEmailClient(t *testing.T) {
    const expect bool = false
    _ = NewClient("test", "example@example.com", 0, Key("000"))
    got := NewClient("asd", "example@example.com", 0, Key("123"))
    if got != expect {
        t.Errorf("Expected: %v, Got: %v\n", expect, got)
    }

}

func TestAddTempSession(t *testing.T) {
    log.SetOutput(os.Stdout)
    role := "test"
    email := "temp@example.com"
    temp := true
    const expect bool = true
    got := AddSession(role, email, temp)
    if got != expect {
        t.Errorf("Expected: %v, Got: %v\n", expect, got)
    }
}

func TestAddSession(t *testing.T) {
    log.SetOutput(os.Stdout)
    role := "test"
    email := "session@example.com"
    temp := false
    const expect bool = true
    got := AddSession(role, email, temp)
    if got != expect {
        t.Errorf("Expected: %v, Got: %v\n", expect, got)
    }
}

/*
func TestCullExpired(t *testing.T) {
    cullExpired(&clientsbySession)
    if got != expect {
        t.Errorf("Expected: %v, Got: %v\n", expect, got)
    }
}
*/
