package main

import (
    "testing"
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
