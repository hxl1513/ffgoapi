package server

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/hxl1513/ffgoapi/config"
)

func TestRegisterRoutes(t *testing.T) {
	gin.SetMode(gin.TestMode)

	cfg := &config.Config{
		Server: config.ServerConfig{
			Port:         0,
			Mode:         "debug",
			ReadTimeout:  10,
			WriteTimeout: 10,
		},
		Database: config.DatabaseConfig{
			Type:   "sqlite",
			DBName: ":memory:",
		},
		Redis: config.RedisConfig{
			Enabled:  false,
			URL:      "",
			PoolSize: 10,
		},
		Auth: config.AuthConfig{
			SessionSecret: "test-secret",
		},
		Logger: config.LoggerConfig{
			Level: "info",
		},
	}

	srv, err := NewServer(cfg)
	if err != nil {
		t.Fatalf("NewServer() error = %v", err)
	}
	defer srv.Shutdown(nil)

	called := false
	srv.RegisterRoutes(func(r *gin.RouterGroup) {
		called = true
	})

	if !called {
		t.Error("Route function was not called")
	}

	// Verify routes are mounted under /api
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/test", nil)
	srv.Engine().ServeHTTP(w, req)
	if w.Code != http.StatusNotFound {
		t.Errorf("expected status 404 for unregistered route, got %d", w.Code)
	}
}
