package ratelimiter

import (
	"time"
	"testing"
)


func TestGetBucket(t *testing.T) {
	var limiter = NewLimiter(60000, 1)
	var ipLimiter = IPLimiter{
		limitRule:		*limiter,
		buckets:		make(map[IP]*LimitRule),
		lastResetTime:	time.Now(),
	}
	var bucket1 = ipLimiter.GetBucket("first")
	var bucket2 = ipLimiter.GetBucket("first")
	if bucket1.tokens != bucket2.tokens {
		t.Fatalf("Expected: '%v', Got: '%v'\n", bucket2.tokens, bucket1.tokens)
	}
	bucket2.tokens = 0
	var bucket3 = ipLimiter.GetBucket("second")
	if bucket2.tokens == 60000 {
		t.Fatalf("Expected: '%v', Got: '%v'\n", "!60000", bucket2.tokens)
	}
	if bucket3.tokens != 60000 {
		t.Fatalf("Expected: '%v', Got: '%v'\n", "60000", bucket3.tokens)
	}
}

func TestWatchdog(t *testing.T) {
	var limiter = NewLimiter(60000, 1)
	var ipLimiter = IPLimiter{
		limitRule:		*limiter,
		buckets:		make(map[IP]*LimitRule),
		lastResetTime:	time.Now(),
	}
	var bucket = ipLimiter.GetBucket("first")
	bucket.tokens = 0
	go ipLimiterWatchdog(&ipLimiter, (1/(60000.0-10)), 90)
	bucket = ipLimiter.GetBucket("first")
	if bucket.tokens == 60000 {
		t.Fatalf("Expected: '%v', Got: '%v'\n", "!60000", bucket.tokens)
	}
	time.Sleep(time.Millisecond * 4)
	bucket = ipLimiter.GetBucket("first")
	if bucket.tokens != 60000 {
		t.Fatalf("Expected: '%v', Got: '%v'\n", "!60000", bucket.tokens)
	}
	bucket.tokens = 0
	bucket = ipLimiter.GetBucket("first")
	if bucket.tokens == 60000 {
		t.Fatalf("Expected: '%v', Got: '%v'\n", "!60000", bucket.tokens)
	}
	time.Sleep(time.Millisecond * 4)
	bucket = ipLimiter.GetBucket("first")
	if bucket.tokens != 60000 {
		t.Fatalf("Expected: '%v', Got: '%v'\n", "60000", bucket.tokens)
	}
}

func TestWatchdogLimiterReport(t *testing.T) {
	var limiter = NewLimiter(60000, 1)
	var ipLimiter = IPLimiter{
		limitRule:		*limiter,
		buckets:		make(map[IP]*LimitRule),
		lastResetTime:	time.Now(),
	}
	var bucket = ipLimiter.GetBucket("first")
	bucket.tokens = 0
	go ipLimiterWatchdog(&ipLimiter, (1/(60000.0-10)), 90)
	for range 110 {
		bucket.Allow()
	}
	time.Sleep(time.Millisecond * 10)
	bucket = ipLimiter.GetBucket("first")
	if bucket.tokens != 60000 {
		t.Fatalf("Expected: '%v', Got: '%v'\n", "!60000", bucket.tokens)
	}
	t.Log("^^^^^^ Above should be Rate Limit Summary ^^^^^^")
}

func TestCreateNewIPLimiterRule(t *testing.T) {
	ipLimiter := NewIPLimiterRule(RateLimiterConfig{
		MaxTokens:			6000,
		TokensPerMinute:	1,
		ResetMinutes:		(1/(60000.0-10)),
		AlertLimit:			90,
	})
	var bucket = ipLimiter.GetBucket("first")
	bucket.tokens = 0
	for range 110 {
		bucket.Allow()
	}
	time.Sleep(time.Millisecond * 30)
	bucket = ipLimiter.GetBucket("first")
	if !bucket.Allow() {
		t.Fatalf("Expected: '%v', Got: '%v'\n", "ok", "not")
	}
}
