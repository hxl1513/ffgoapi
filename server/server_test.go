package server

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/your-org/go-web-api-framework/config"
	"github.com/your-org/go-web-api-framework/model"
)

func TestNewServer(t *testing.T) {
	cfg := &config.Config{
		Server: config.ServerConfig{
			Port:         0,
			Mode:         "debug",
			ReadTimeout:  10 * time.Second,
			WriteTimeout: 10 * time.Second,
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
	defer srv.Shutdown(context.Background())

	if srv.Engine() == nil {
		t.Error("Engine() returned nil")
	}

	// Verify health check endpoint is registered
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/health", nil)
	srv.Engine().ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Errorf("GET /health status = %d, want %d", w.Code, http.StatusOK)
	}

	// Verify db and cache are initialised
	if srv.db == nil {
		t.Error("db is nil")
	}
	if srv.cache == nil {
		t.Error("cache is nil")
	}

	_ = model.ErrCacheMiss // ensure model package referenced
}

func TestServerStart(t *testing.T) {
	cfg := &config.Config{
		Server: config.ServerConfig{
			Port:         0,
			Mode:         "debug",
			ReadTimeout:  10 * time.Second,
			WriteTimeout: 10 * time.Second,
		},
		Database: config.DatabaseConfig{
			Type:   "sqlite",
			DBName: ":memory:",
		},
		Redis: config.RedisConfig{
			Enabled: false,
			URL:     "",
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

	go func() {
		_ = srv.Start()
	}()

	// Give the server a moment to start
	time.Sleep(200 * time.Millisecond)

	// Verify server responds on the health endpoint via the engine directly
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/health", nil)
	srv.Engine().ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Errorf("GET /health status = %d, want %d", w.Code, http.StatusOK)
	}

	srv.Shutdown(context.Background())
}
