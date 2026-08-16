package core

import (
	"fmt"
	"log/slog"
	"sync"
	"time"
)

type CircuitBreakerState string

const (
	StateClosed   CircuitBreakerState = "CLOSED"
	StateOpen     CircuitBreakerState = "OPEN"
	StateHalfOpen CircuitBreakerState = "HALF_OPEN"
)

type CircuitBreaker struct {
	mu           sync.RWMutex
	state        CircuitBreakerState
	failureCount int
	threshold    int
	lastFailTime time.Time
	resetTimeout time.Duration
	logger       *slog.Logger
}

func NewCircuitBreaker(threshold int, resetTimeout time.Duration, logger *slog.Logger) *CircuitBreaker {
	if logger == nil {
		logger = GetLogger()
	}
	if threshold <= 0 {
		threshold = 5 // Open circuit after 5 consecutive failures / 429s
	}
	if resetTimeout <= 0 {
		resetTimeout = 30 * time.Second
	}
	return &CircuitBreaker{
		state:        StateClosed,
		threshold:    threshold,
		resetTimeout: resetTimeout,
		logger:       logger.With("component", "circuit_breaker"),
	}
}

func (cb *CircuitBreaker) AllowRequest() bool {
	cb.mu.Lock()
	defer cb.mu.Unlock()

	if cb.state == StateOpen {
		if time.Since(cb.lastFailTime) > cb.resetTimeout {
			cb.state = StateHalfOpen
			cb.logger.Info("Circuit Breaker switching from OPEN to HALF_OPEN, probing provider API...")
			return true
		}
		return false
	}
	return true
}

func (cb *CircuitBreaker) RecordResult(err error) {
	cb.mu.Lock()
	defer cb.mu.Unlock()

	if err != nil {
		cb.failureCount++
		cb.lastFailTime = time.Now()
		if cb.failureCount >= cb.threshold {
			cb.state = StateOpen
			cb.logger.Warn("Circuit Breaker OPENED due to consecutive API failures / 429 rate limits", "failure_count", cb.failureCount)
		}
	} else {
		cb.failureCount = 0
		if cb.state == StateHalfOpen {
			cb.state = StateClosed
			cb.logger.Info("Circuit Breaker CLOSED after successful probe request")
		}
	}
}

type RateLimiterBucket struct {
	tokens     float64
	maxTokens  float64
	refillRate float64 // Tokens per second
	lastRefill time.Time
}

type TenantAPIRateLimiter struct {
	mu      sync.RWMutex
	buckets map[string]*RateLimiterBucket
}

func NewTenantAPIRateLimiter() *TenantAPIRateLimiter {
	return &TenantAPIRateLimiter{
		buckets: make(map[string]*RateLimiterBucket),
	}
}

func (rl *TenantAPIRateLimiter) Allow(tenantID, provider string, limitPerSec float64) bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	key := tenantID + ":" + provider
	bucket, exists := rl.buckets[key]
	now := time.Now()

	if !exists {
		if limitPerSec <= 0 {
			limitPerSec = 10.0 // Default 10 requests per second per provider per tenant
		}
		bucket = &RateLimiterBucket{
			tokens:     limitPerSec,
			maxTokens:  limitPerSec,
			refillRate: limitPerSec,
			lastRefill: now,
		}
		rl.buckets[key] = bucket
	}

	// Refill tokens based on elapsed time
	elapsed := now.Sub(bucket.lastRefill).Seconds()
	bucket.tokens += elapsed * bucket.refillRate
	if bucket.tokens > bucket.maxTokens {
		bucket.tokens = bucket.maxTokens
	}
	bucket.lastRefill = now

	if bucket.tokens >= 1.0 {
		bucket.tokens -= 1.0
		return true
	}
	return false
}

type SmartMCPCacheEntry struct {
	Value     interface{}
	ExpiresAt time.Time
}

type SmartMCPCache struct {
	mu      sync.RWMutex
	entries map[string]SmartMCPCacheEntry
}

func NewSmartMCPCache() *SmartMCPCache {
	return &SmartMCPCache{
		entries: make(map[string]SmartMCPCacheEntry),
	}
}

func (c *SmartMCPCache) Get(key string) (interface{}, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	entry, ok := c.entries[key]
	if !ok || time.Now().After(entry.ExpiresAt) {
		return nil, false
	}
	return entry.Value, true
}

func (c *SmartMCPCache) Set(key string, val interface{}, ttl time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if ttl <= 0 {
		ttl = 5 * time.Minute // 5-minute default cache for idempotent tool calls
	}
	c.entries[key] = SmartMCPCacheEntry{
		Value:     val,
		ExpiresAt: time.Now().Add(ttl),
	}
}

func BuildCacheKey(tenantID, serverID, toolName string, args map[string]interface{}) string {
	return fmt.Sprintf("%s:%s:%s:%v", tenantID, serverID, toolName, args)
}
