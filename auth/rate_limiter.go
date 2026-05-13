package auth

import (
	"sync"
	"time"
)

// RateLimiter implements rate limiting for authentication endpoints
type RateLimiter struct {
	attempts map[string]int
	mu       sync.Mutex
}

// NewRateLimiter creates a new rate limiter
func NewRateLimiter() *RateLimiter {
	return &RateLimiter{
		attempts: make(map[string]int),
	}
}

// AllowRequest checks if a request from the given IP is allowed
func (rl *RateLimiter) AllowRequest(ip string) bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()
	
	count := rl.attempts[ip]
	if count >= 5 {
		return false
	}
	
	rl.attempts[ip] = count + 1
	
	// Reset counter after 1 minute
	go func() {
		time.AfterFunc(1*time.Minute, func() {
			rl.mu.Lock()
			delete(rl.attempts, ip)
			rl.mu.Unlock()
		})
	}()
	
	return true
}