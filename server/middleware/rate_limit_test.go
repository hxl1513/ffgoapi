package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/hxl1513/ffgoapi/model"
)

func TestRateLimitMiddleware_Disabled(t *testing.T) {
	gin.SetMode(gin.TestMode)
	cache := model.NewMemoryCache()
	cfg := &RateLimitConfig{
		Enabled:         false,
		DefaultMax:      1,
		DefaultDuration: 60,
	}

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test", nil)
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	handler := RateLimit(cfg, cache)
	handler(c)

	if c.Writer.Status() == http.StatusTooManyRequests {
		t.Errorf("Rate limit disabled, should not return 429, got %d", c.Writer.Status())
	}
}

func TestRateLimitMiddleware_AllowRequests(t *testing.T) {
	gin.SetMode(gin.TestMode)
	cache := model.NewMemoryCache()
	cfg := &RateLimitConfig{
		Enabled:         true,
		DefaultMax:      2,
		DefaultDuration: 60,
	}

	for i := 0; i < 2; i++ {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/test", nil)
		c, _ := gin.CreateTestContext(w)
		c.Request = req

		handler := RateLimit(cfg, cache)
		handler(c)

		if c.Writer.Status() == http.StatusTooManyRequests {
			t.Errorf("Request %d should be allowed, got 429", i+1)
		}
	}
}

func TestRateLimitMiddleware_BlockExcess(t *testing.T) {
	gin.SetMode(gin.TestMode)
	cache := model.NewMemoryCache()
	cfg := &RateLimitConfig{
		Enabled:         true,
		DefaultMax:      2,
		DefaultDuration: 60,
	}

	handler := RateLimit(cfg, cache)

	// First 2 requests should pass
	for i := 0; i < 2; i++ {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/test", nil)
		c, _ := gin.CreateTestContext(w)
		c.Request = req
		handler(c)
	}

	// 3rd request should be blocked
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test", nil)
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	handler(c)

	if c.Writer.Status() != http.StatusTooManyRequests {
		t.Errorf("Expected 429 after exceeding limit, got %d", c.Writer.Status())
	}
}

func TestRateLimitMiddleware_DifferentIPs(t *testing.T) {
	gin.SetMode(gin.TestMode)
	cache := model.NewMemoryCache()
	cfg := &RateLimitConfig{
		Enabled:         true,
		DefaultMax:      1,
		DefaultDuration: 60,
	}

	handler := RateLimit(cfg, cache)

	// IP A: first request allowed
	w1 := httptest.NewRecorder()
	req1, _ := http.NewRequest("GET", "/test", nil)
	req1.RemoteAddr = "10.0.0.1:1234"
	c1, _ := gin.CreateTestContext(w1)
	c1.Request = req1
	handler(c1)
	if w1.Code == http.StatusTooManyRequests {
		t.Errorf("First request from IP A should be allowed, got 429")
	}

	// IP B: first request should also be allowed
	w2 := httptest.NewRecorder()
	req2, _ := http.NewRequest("GET", "/test", nil)
	req2.RemoteAddr = "10.0.0.2:1234"
	c2, _ := gin.CreateTestContext(w2)
	c2.Request = req2
	handler(c2)
	if w2.Code == http.StatusTooManyRequests {
		t.Errorf("First request from IP B should be allowed, got 429")
	}
}
