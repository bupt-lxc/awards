package middleware

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

// RateLimiter 基于令牌桶的内存限流器
type RateLimiter struct {
	mu       sync.Mutex
	buckets  map[string]*tokenBucket
	rate     float64 // 每秒允许的请求数
	capacity int     // 桶容量
}

type tokenBucket struct {
	tokens   float64
	lastTime time.Time
}

// NewRateLimiter 创建限流器
func NewRateLimiter(requestsPerMinute, capacity int) *RateLimiter {
	return &RateLimiter{
		buckets:  make(map[string]*tokenBucket),
		rate:     float64(requestsPerMinute) / 60.0,
		capacity: capacity,
	}
}

func (rl *RateLimiter) allow(key string) bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	bucket, exists := rl.buckets[key]
	now := time.Now()

	if !exists {
		rl.buckets[key] = &tokenBucket{
			tokens:   float64(rl.capacity) - 1,
			lastTime: now,
		}
		return true
	}

	// 计算补充的令牌
	elapsed := now.Sub(bucket.lastTime).Seconds()
	bucket.tokens += elapsed * rl.rate
	if bucket.tokens > float64(rl.capacity) {
		bucket.tokens = float64(rl.capacity)
	}
	bucket.lastTime = now

	if bucket.tokens >= 1 {
		bucket.tokens--
		return true
	}

	return false
}

// 清理过期桶（每5分钟清理一次）
func (rl *RateLimiter) startCleanup() {
	go func() {
		rl.mu.Lock()
		rl.buckets = make(map[string]*tokenBucket)
		rl.mu.Unlock()
	}()
	time.AfterFunc(5*time.Minute, rl.startCleanup)
}

var (
	ipLimiter    *RateLimiter
	tokenLimiter *RateLimiter
	onceRL       sync.Once
)

func initLimiters() {
	onceRL.Do(func() {
		ipLimiter = NewRateLimiter(30, 30)    // 30 req/min, burst 30
		tokenLimiter = NewRateLimiter(60, 60) // 60 req/min, burst 60

		time.AfterFunc(5*time.Minute, ipLimiter.startCleanup)
		time.AfterFunc(5*time.Minute, tokenLimiter.startCleanup)
	})
}

// IPRateLimit 基于 IP 的限流中间件（用于公开查询接口）
func IPRateLimit() gin.HandlerFunc {
	initLimiters()
	return func(c *gin.Context) {
		ip := c.ClientIP()
		if !ipLimiter.allow(ip) {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"error": "请求过于频繁，请稍后再试",
			})
			c.Abort()
			return
		}
		c.Next()
	}
}

// TokenRateLimit 基于 Token 的限流中间件（用于认证接口）
func TokenRateLimit() gin.HandlerFunc {
	initLimiters()
	return func(c *gin.Context) {
		key, exists := c.Get("username")
		if !exists {
			key = c.ClientIP()
		}
		keyStr, ok := key.(string)
		if !ok {
			keyStr = c.ClientIP()
		}

		if !tokenLimiter.allow(keyStr) {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"error": "请求过于频繁，请稍后再试",
			})
			c.Abort()
			return
		}
		c.Next()
	}
}
