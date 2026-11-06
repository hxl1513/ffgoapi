package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/your-org/go-web-api-framework/model"
	"github.com/your-org/go-web-api-framework/pkg/crypto"
)

func TestAuthMiddleware_NoToken(t *testing.T) {
	gin.SetMode(gin.TestMode)
	cache := model.NewMemoryCache()
	cfg := &AuthConfig{SessionSecret: "test-secret"}

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test", nil)
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	handler := AuthMiddleware(cfg, cache, 0)
	handler(c)

	if c.Writer.Status() != http.StatusUnauthorized {
		t.Errorf("Expected 401, got %d", c.Writer.Status())
	}
}

func TestAuthMiddleware_ValidToken(t *testing.T) {
	gin.SetMode(gin.TestMode)
	cache := model.NewMemoryCache()
	cfg := &AuthConfig{SessionSecret: "test-secret-key"}

	token, err := crypto.GenerateToken(cfg.SessionSecret, 1, "testuser", 2, time.Hour)
	if err != nil {
		t.Fatalf("GenerateToken() error: %v", err)
	}

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	handler := AuthMiddleware(cfg, cache, 0)
	handler(c)

	if c.Writer.Status() == http.StatusUnauthorized {
		t.Errorf("Expected non-401, got 401")
	}

	// User context should be set
	val, exists := c.Get("user_id")
	if !exists || val != 1 {
		t.Errorf("Expected user_id=1, got exists=%v val=%v", exists, val)
	}
	val, exists = c.Get("username")
	if !exists || val != "testuser" {
		t.Errorf("Expected username=testuser, got exists=%v val=%v", exists, val)
	}
	val, exists = c.Get("role")
	if !exists || val != 2 {
		t.Errorf("Expected role=2, got exists=%v val=%v", exists, val)
	}
}

func TestAuthMiddleware_InvalidToken(t *testing.T) {
	gin.SetMode(gin.TestMode)
	cache := model.NewMemoryCache()
	cfg := &AuthConfig{SessionSecret: "test-secret-key"}

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test", nil)
	req.Header.Set("Authorization", "Bearer invalid.token.here")
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	handler := AuthMiddleware(cfg, cache, 0)
	handler(c)

	if c.Writer.Status() != http.StatusUnauthorized {
		t.Errorf("Expected 401, got %d", c.Writer.Status())
	}
}

func TestAuthMiddleware_AllowedRoles(t *testing.T) {
	gin.SetMode(gin.TestMode)
	cache := model.NewMemoryCache()
	cfg := &AuthConfig{SessionSecret: "test-secret-key"}

	// Generate token with role=1
	token, err := crypto.GenerateToken(cfg.SessionSecret, 1, "testuser", 1, 0)
	if err != nil {
		t.Fatalf("GenerateToken() error: %v", err)
	}

	// AllowedRoles is 2,3 - role 1 should be rejected
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	handler := AuthMiddleware(cfg, cache, 2, 3)
	handler(c)

	if c.Writer.Status() != http.StatusForbidden {
		t.Errorf("Expected 403, got %d", c.Writer.Status())
	}
}

func TestAdminAuth(t *testing.T) {
	gin.SetMode(gin.TestMode)
	cache := model.NewMemoryCache()
	cfg := &AuthConfig{SessionSecret: "test-secret-key"}

	// Admin role = 1
	token, err := crypto.GenerateToken(cfg.SessionSecret, 1, "admin", 1, 0)
	if err != nil {
		t.Fatalf("GenerateToken() error: %v", err)
	}

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/admin", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	handler := AdminAuth(cfg, cache)
	handler(c)

	if c.Writer.Status() == http.StatusUnauthorized || c.Writer.Status() == http.StatusForbidden {
		t.Errorf("Admin auth should pass with admin role, got %d", c.Writer.Status())
	}
}

func TestUserAuth(t *testing.T) {
	gin.SetMode(gin.TestMode)
	cache := model.NewMemoryCache()
	cfg := &AuthConfig{SessionSecret: "test-secret-key"}

	// Regular user role = 2
	token, err := crypto.GenerateToken(cfg.SessionSecret, 2, "user", 2, 0)
	if err != nil {
		t.Fatalf("GenerateToken() error: %v", err)
	}

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/user", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	handler := UserAuth(cfg, cache)
	handler(c)

	if c.Writer.Status() == http.StatusUnauthorized || c.Writer.Status() == http.StatusForbidden {
		t.Errorf("User auth should pass with user role, got %d", c.Writer.Status())
	}
}
