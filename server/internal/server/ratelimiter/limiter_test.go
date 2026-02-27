package ratelimiter

import (
	"time"
	"testing"
)


func TestRefillTokens (t *testing.T) {
	var rule = NewLimiter(60000, 60000)
	rule.tokens = 0
	if got := rule.tokens ; got != 0 {
		t.Fatalf("Expected: '%v', Got: '%v'\n", 0, got)
	}
	time.Sleep(time.Millisecond * 5)
	rule.refillTokens()
	if got := rule.tokens ; got > 5.5 {
		t.Fatalf("Expected: '%v', Got: '%v'\n", "~1", got)
	}
	if got := rule.tokens ; got < 4.5 {
		t.Fatalf("Expected: '%v', Got: '%v'\n", "~1", got)
	}
}

func TestAllowLogic (t *testing.T) {
	var rule = NewLimiter(60000, 60000)
	rule.tokens = 0
	if got := rule.tokens ; got != 0 {
		t.Fatalf("Expected: '%v', Got: '%v'\n", 0, got)
	}
	if rule.Allow() {
		t.Fatalf("Expected: '%v', Got: '%v'\n", false, true)
	}
	time.Sleep(time.Millisecond * 2)
	if !rule.Allow() {
		t.Fatalf("Expected: '%v', Got: '%v'\n", !false, !true)
	}
}

func TestFullTokens (t *testing.T) {
	var rule = NewLimiter(60000, 60000)
	if got := rule.tokens ; got != 60000 {
		t.Fatalf("Expected: '%v', Got: '%v'\n", 60000, got)
	}
	time.Sleep(time.Millisecond * 2)
	rule.refillTokens()
	if got := rule.tokens ; got != 60000 {
		t.Fatalf("Expected: '%v', Got: '%v'\n", 60000, got)
	}
}
