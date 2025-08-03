package middleware

import (
	"sync"
	"time"

	"github.com/devdahcoder/golang-todo-api/pkg/errors"
	"github.com/gofiber/fiber/v3"
)

type RateLimiter struct {
	mu       sync.Mutex
	requests map[string][]time.Time
	limit    int
	window   time.Duration
}

func NewRateLimiter(limit int, window time.Duration) *RateLimiter {
	return &RateLimiter{
		requests: make(map[string][]time.Time),
		limit:    limit,
		window:   window,
	}
}

func (rl *RateLimiter) Cleanup() {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	now := time.Now()
	for ip, times := range rl.requests {
		var valid []time.Time
		for _, t := range times {
			if now.Sub(t) <= rl.window {
				valid = append(valid, t)
			}
		}
		if len(valid) == 0 {
			delete(rl.requests, ip)
		} else {
			rl.requests[ip] = valid
		}
	}
}

func (rl *RateLimiter) Middleware() fiber.Handler {
	go func() {
		for {
			time.Sleep(rl.window)
			rl.Cleanup()
		}
	}()

	return func(c fiber.Ctx) error {
		ip := c.IP()
		now := time.Now()

		rl.mu.Lock()
		times := rl.requests[ip]
		
		var valid []time.Time
		for _, t := range times {
			if now.Sub(t) <= rl.window {
				valid = append(valid, t)
			}
		}

		if len(valid) >= rl.limit {
			rl.mu.Unlock()
			return errors.TooManyRequestsError(c, "rate limit exceeded", nil)
		}

		valid = append(valid, now)
		rl.requests[ip] = valid
		rl.mu.Unlock()

		return c.Next()
	}
} 