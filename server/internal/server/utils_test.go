package server

import (
	"testing"
    "net/http"
    ware "github.com/JValtteri/qure/server/internal/middleware"
)


func TestSanitizeAllow(t *testing.T) {
    const input string = "asd 123-5"
    var got string = sanitize(input)
    if got != input {
        t.Errorf("Expected: %s, Got: %s\n", input, got)
    }
}

func TestSanitizeBlock(t *testing.T) {
    const input string = "<asd=1,23>"
    const expect string = "asd123"
    var got string = sanitize(input)
    if got != expect {
        t.Errorf("Expected: '%s', Got: '%s'\n", expect, got)
    }
}

func TestLoadBody(t *testing.T) {
    requestObj := ware.LoginRequest{}
    incomingRequest := ware.LoginRequest{
        User: "abc",
        Password: "123",
    }
    requestBodyWriter := makeWriter(incomingRequest)
    req, err := http.NewRequest("POST", "/api/user/login", requestBodyWriter)
    if err != nil {
        t.Fatalf("NewRequest returned an error: %v\n", err)
    }
    got, err := loadRequestBody(req, requestObj)
    if err != nil {
        t.Fatalf("LoadBody returned an error: %v\n", err)
    }
    if got.User != incomingRequest.User {
        t.Errorf("Expected: '%s', Got: '%s'\n", incomingRequest.User, got.User)
    }
}
