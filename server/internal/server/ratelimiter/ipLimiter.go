// Algorithm based on a tutorial by Vivek Alhat
// https://dev.to/vivekalhat/rate-limiting-for-beginners-what-it-is-and-how-to-build-one-in-go-955
//
// Original source code hosted here:
// https://github.com/VivekAlhat/go-rate-limiter

package ratelimiter

import (
	"log"
	"sync"
	"time"
)


type IP string

// IP based rate limiter
type IPLimiter struct {
	limitRule		LimitRule
	buckets			map[IP]*LimitRule
	lastResetTime	time.Time
	mutex			sync.Mutex
}


func (i *IPLimiter) GetBucket(ipAddress string) *LimitRule {
	var ip = IP(ipAddress)
	i.mutex.Lock()
	defer i.mutex.Unlock()

	limiter, exists := i.buckets[ip]
	if !exists {
		limiter = NewLimiter(i.limitRule.maxTokens, i.limitRule.tokensPerMinute)
		i.buckets[ip] = limiter
	}
	return limiter
}

type RateLimiterConfig struct {
	MaxTokens		float64
	TokensPerMinute	float64
	ResetMinutes	float32
	AlertLimit		uint64
}

func NewIPLimiterRule(c RateLimiterConfig) *IPLimiter {
	var ipLimiter = &IPLimiter{
		limitRule: LimitRule{
			maxTokens: 			c.MaxTokens,
			tokensPerMinute:	c.TokensPerMinute,
		},
		buckets:		make(map[IP]*LimitRule),
		lastResetTime:	time.Now(),
	}
	go ipLimiterWatchdog(ipLimiter, c.ResetMinutes, c.AlertLimit)		// Reset limiter every 60 minutes
	return ipLimiter
}

// Clears Limiter every (N minutes).
// This is to clear old entries from map
func ipLimiterWatchdog(ipLimiter *IPLimiter, resetIntervalMinutes float32, alertLimit uint64) {
	const MILLISECONDS_IN_MINUTE = 1000 * 60
	var tick = time.Millisecond * time.Duration(resetIntervalMinutes * MILLISECONDS_IN_MINUTE)
	for range time.NewTicker(tick).C {
		ipLimiter.mutex.Lock()					// LOCK
		for ip, r := range ipLimiter.buckets {
			if r.limiterHits > alertLimit {
				log.Printf("Rate limiter summary: %v: %v limiter hits\n", ip, r.limiterHits)
			}
		}
		ipLimiter.buckets = make(map[IP]*LimitRule)
		ipLimiter.lastResetTime = time.Now()
		ipLimiter.mutex.Unlock()				// UNLOCK
	}
}
