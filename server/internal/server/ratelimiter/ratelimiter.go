// Algorithm based on a tutorial by Vivek Alhat
// https://dev.to/vivekalhat/rate-limiting-for-beginners-what-it-is-and-how-to-build-one-in-go-955
//
// Original source code hosted here:
// https://github.com/VivekAlhat/go-rate-limiter

package ratelimiter

import (
	"fmt"
	"net"
	"net/http"
)


func RateLimiter(limiter *IPLimiter, handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ip, ok := getClientIP(r, w)
		if !ok {
			return
		}
		ipLimiter := limiter.GetBucket(ip)
		if ipLimiter.Allow() {
			handler(w, r)
		} else {
			http.Error(w, fmt.Sprintln("Rate Limit Exceeded"), http.StatusTooManyRequests)
		}
	}
}

func getClientIP(r *http.Request, w http.ResponseWriter) (string, bool) {
	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		http.Error(w, "Invalid IP", http.StatusInternalServerError)
		return "", false
	}
	return ip, true
}
