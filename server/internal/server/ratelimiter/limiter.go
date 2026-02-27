// Algorithm based on a tutorial by Vivek Alhat
// https://dev.to/vivekalhat/rate-limiting-for-beginners-what-it-is-and-how-to-build-one-in-go-955
//
// Original source code hosted here:
// https://github.com/VivekAlhat/go-rate-limiter

package ratelimiter

import (
	"sync"
	"time"
)


// Rate limiter implementation based on token bucket algorithm
type LimitRule struct {
	tokens			float64
	maxTokens		float64
	tokensPerMinute	float64
	lastRefillTime	time.Time
	limiterHits		uint64
	mutex			sync.Mutex
}

func (r *LimitRule) refillTokens() {
	now := time.Now()
	duration := now.Sub(r.lastRefillTime).Minutes()
	r.lastRefillTime = now

	r.tokens += duration * r.tokensPerMinute
	if r.tokens > r.maxTokens {
		r.tokens = r.maxTokens
	}
}

func (r *LimitRule) Allow() bool {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	r.refillTokens()
	if r.tokens >= 1 {
		r.tokens--
		return true
	}
	r.limiterHits += 1
	return false
}

func NewLimiter(maxTokens float64, tokensPerMinute float64) *LimitRule {
	return &LimitRule{
		tokens:				maxTokens,
		maxTokens:			maxTokens,
		tokensPerMinute:	tokensPerMinute,
		limiterHits:		0,
		lastRefillTime:		time.Now(),
	}
}
