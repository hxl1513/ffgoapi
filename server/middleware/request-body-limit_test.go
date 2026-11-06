package middleware

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

// readBodyHandler reads the entire request body into memory. When the body
// exceeds the MaxBytesReader limit, the Read will return *http.MaxBytesError
// and gin will respond with 413.
func readBodyHandler(c *gin.Context) {
	buf, err := c.GetRawData()
	if err != nil {
		c.AbortWithStatus(http.StatusRequestEntityTooLarge)
		return
	}
	c.String(http.StatusOK, "read %d bytes", len(buf))
}

func TestRequestBodyLimit(t *testing.T) {
	gin.SetMode(gin.TestMode)

	largeBody := bytes.Repeat([]byte("a"), 1024*1024) // 1MB

	r := gin.New()
	r.Use(RequestBodyLimit(512 * 1024)) // 512KB limit
	r.POST("/test", readBodyHandler)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/test", bytes.NewReader(largeBody))
	req.ContentLength = int64(len(largeBody))
	r.ServeHTTP(w, req)

	if w.Code != http.StatusRequestEntityTooLarge {
		t.Errorf("Expected 413, got %d", w.Code)
	}
}

func TestRequestBodyLimit_UnderLimit(t *testing.T) {
	gin.SetMode(gin.TestMode)

	body := bytes.Repeat([]byte("a"), 1024) // 1KB

	r := gin.New()
	r.Use(RequestBodyLimit(512 * 1024)) // 512KB limit
	r.POST("/test", readBodyHandler)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/test", bytes.NewReader(body))
	req.ContentLength = int64(len(body))
	r.ServeHTTP(w, req)

	if w.Code == http.StatusRequestEntityTooLarge {
		t.Errorf("Expected no 413 for body under limit, got %d", w.Code)
	}
}

func TestRequestBodyLimit_Default(t *testing.T) {
	gin.SetMode(gin.TestMode)

	largeBody := bytes.Repeat([]byte("a"), 5*1024*1024) // 5MB

	r := gin.New()
	r.Use(RequestBodyLimit(0)) // 0 triggers default 4MB
	r.POST("/test", readBodyHandler)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/test", bytes.NewReader(largeBody))
	req.ContentLength = int64(len(largeBody))
	r.ServeHTTP(w, req)

	if w.Code != http.StatusRequestEntityTooLarge {
		t.Errorf("Expected 413 for body exceeding default 4MB limit, got %d", w.Code)
	}
}

func TestBodyCleanup(t *testing.T) {
	gin.SetMode(gin.TestMode)

	r := gin.New()
	r.Use(BodyCleanup())
	r.POST("/test", func(c *gin.Context) {
		c.String(http.StatusOK, "ok")
	})

	body := bytes.NewReader([]byte("test body"))
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/test", body)
	req.ContentLength = int64(len("test body"))
	r.ServeHTTP(w, req)

	if w.Code == http.StatusRequestEntityTooLarge {
		t.Errorf("Expected no 413, got %d", w.Code)
	}

	if w.Code != http.StatusOK {
		t.Logf("Unexpected status %d, body: %s", w.Code, w.Body.String())
	}
}

func TestRequestBodyLimit_LargeBody(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// 10MB body with 1MB limit
	largeBody := bytes.Repeat([]byte("x"), 10*1024*1024)

	r := gin.New()
	r.Use(RequestBodyLimit(1024 * 1024)) // 1MB limit
	r.POST("/test", readBodyHandler)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/test", bytes.NewReader(largeBody))
	req.ContentLength = int64(len(largeBody))
	r.ServeHTTP(w, req)

	if w.Code != http.StatusRequestEntityTooLarge {
		t.Errorf("Expected 413, got %d", w.Code)
	}
}
