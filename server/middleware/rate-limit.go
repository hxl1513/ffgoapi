package middleware

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/hxl1513/ffgoapi/model"
)

// RateLimitConfig holds rate limiting configuration.
type RateLimitConfig struct {
	Enabled         bool
	DefaultMax      int
	DefaultDuration int // seconds
}

// RateLimit returns a gin.HandlerFunc that enforces per-IP rate limiting using
// a sliding-window counter stored in cache. If the cache is unavailable the
// middleware fails open (allows the request).
func RateLimit(cfg *RateLimitConfig, cache model.Cache) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 1. If disabled, pass through
		if !cfg.Enabled {
			c.Next()
			return
		}

		// 2. Generate per-IP key
		ip := c.ClientIP()
		key := "rate_limit:" + ip

		// 3. Use cache Get / Set / Incr for sliding window
		val, err := cache.Get(key)
		if err != nil {
			// Key does not exist yet — initialise to 1 with expiry
			_ = cache.Set(key, "1", time.Duration(cfg.DefaultDuration)*time.Second)
			c.Next()
			return
		}

		// Key exists — increment
		err = cache.Incr(key, 1)
		if err != nil {
			// Cache error: fail open
			c.Next()
			return
		}

		// Parse current count from stored value
		count := parseInt(val) + 1

		// 4. Return 429 if limit exceeded
		if count > cfg.DefaultMax {
			c.JSON(http.StatusTooManyRequests, map[string]any{
				"success": false,
				"message": "rate limit exceeded",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

// parseInt converts a decimal string to int. Returns 0 on any parse error.
func parseInt(s string) int {
	n := 0
	for i := 0; i < len(s); i++ {
		if s[i] < '0' || s[i] > '9' {
			return 0
		}
		n = n*10 + int(s[i]-'0')
	}
	return n
}
