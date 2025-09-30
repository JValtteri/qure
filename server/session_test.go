package main

import (
    "testing"
    "log"
    "os"
)

func resetClients() {
    clients = Clients{
        raw:        make(map[ID]Client),
        bySession:  make(map[Key]*Client),
        byEmail:    make(map[string]*Client),
    }
}

func TestAddTempSession(t *testing.T) {
    log.SetOutput(os.Stdout)
    role := "test"
    email := "temp@example.com"
    ip := IP("0.0.0.0")
    temp := true
    expect := 16
    got, err := AddSession(role, email, temp, ip)
    if err != nil {
        t.Errorf("Expected: %v, Got: %v\n", nil, err)
    }
    if len(got) != expect {
        t.Errorf("Expected: %v, Got: %v\n", expect, got)
    }
}

func TestAddSessions(t *testing.T) {
    resetClients()
    log.SetOutput(os.Stdout)
    role := "test"
    email := "session@example.com"
    ip := IP("0.0.0.0")
    temp := false
    expect := 16
    got, err := AddSession(role, email, temp, ip)
    if err != nil {
        t.Errorf("Expected: %v, Got: %v\n", nil, err)
    }
    if len(got) != expect {
        t.Errorf("Expected: %v, Got: %v\n", expect, got)
    }
    got, err = AddSession(role, email, temp, ip)
    if err != nil {
        t.Errorf("Expected: %v, Got: %v\n", nil, err)
    }
    if len(got) != expect {
        t.Errorf("Expected: %v, Got: %v\n", expect, got)
    }

}

func TestResumeSession(t *testing.T) {
    resetClients()
    role := "test"
    email := "session@example.com"
    ip := IP("0.0.0.0")
    temp := false
    key, err := AddSession(role, email, temp, ip)
    t.Logf("Search Session Key: %v", key)
    if err != nil {
        t.Errorf("Expected: %v, Got: %v\n", nil, err)
    }
    err = ResumeSession(key, ip)
    _, found := clients.bySession[key]
    t.Logf("Found %v", found)
    for k, v := range clients.bySession {
        if len(k) > 4 {
            //t.Logf("Found  Session Key: %v", k)
            for i, j := range v.sessionsKeys {
                t.Logf("SessionKey: %v, %v", i, j.ip)
            }
        }
    }
    if err != nil {
        t.Errorf("Expected: '%v', Got: '%v'\n", nil, err)
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
