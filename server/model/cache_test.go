package model

import (
	"testing"
	"time"
)

func TestMemoryCache(t *testing.T) {
	cache := NewMemoryCache()

	// Set and Get
	cache.Set("key1", "value1", 1*time.Minute)
	val, err := cache.Get("key1")
	if err != nil {
		t.Fatalf("Get() unexpected error: %v", err)
	}
	if val != "value1" {
		t.Errorf("Get() = %q, want %q", val, "value1")
	}

	// Exists
	exists, err := cache.Exists("key1")
	if err != nil {
		t.Fatalf("Exists() unexpected error: %v", err)
	}
	if !exists {
		t.Error("Exists() should return true for existing key")
	}

	// Del
	err = cache.Del("key1")
	if err != nil {
		t.Fatalf("Del() unexpected error: %v", err)
	}
	_, err = cache.Get("key1")
	if err == nil {
		t.Error("expected error after delete")
	}
}

func TestMemoryCacheTTL(t *testing.T) {
	cache := NewMemoryCache()
	cache.Set("key1", "value1", 100*time.Millisecond)

	ttl, err := cache.TTL("key1")
	if err != nil {
		t.Fatalf("TTL() unexpected error: %v", err)
	}
	if ttl <= 0 {
		t.Errorf("TTL() = %v, want positive duration", ttl)
	}

	time.Sleep(150 * time.Millisecond)
	_, err = cache.Get("key1")
	if err == nil {
		t.Error("expected error after TTL expiry")
	}
}

func TestMemoryCacheHSetHGet(t *testing.T) {
	cache := NewMemoryCache()
	err := cache.HSet("hash1", "field1", "value1")
	if err != nil {
		t.Fatalf("HSet() unexpected error: %v", err)
	}
	val, err := cache.HGet("hash1", "field1")
	if err != nil {
		t.Fatalf("HGet() unexpected error: %v", err)
	}
	if val != "value1" {
		t.Errorf("HGet() = %q, want %q", val, "value1")
	}
}

func TestMemoryCacheHGetAll(t *testing.T) {
	cache := NewMemoryCache()
	cache.HSet("hash1", "field1", "value1")
	cache.HSet("hash1", "field2", "value2")
	result, err := cache.HGetAll("hash1")
	if err != nil {
		t.Fatalf("HGetAll() unexpected error: %v", err)
	}
	if result["field1"] != "value1" {
		t.Errorf("field1 = %q, want %q", result["field1"], "value1")
	}
	if result["field2"] != "value2" {
		t.Errorf("field2 = %q, want %q", result["field2"], "value2")
	}
}

func TestMemoryCacheIncr(t *testing.T) {
	cache := NewMemoryCache()
	cache.Set("counter", "10", 0)
	err := cache.Incr("counter", 5)
	if err != nil {
		t.Fatalf("Incr() unexpected error: %v", err)
	}
	val, err := cache.Get("counter")
	if err != nil {
		t.Fatalf("Get() unexpected error: %v", err)
	}
	if val != "15" {
		t.Errorf("Incr() result = %q, want %q", val, "15")
	}
}

func TestMemoryCacheExistsMissing(t *testing.T) {
	cache := NewMemoryCache()
	exists, err := cache.Exists("nonexistent")
	if err != nil {
		t.Fatalf("Exists() unexpected error: %v", err)
	}
	if exists {
		t.Error("Exists() should return false for missing key")
	}
}

func TestMemoryCacheDelMissing(t *testing.T) {
	cache := NewMemoryCache()
	err := cache.Del("nonexistent")
	if err != nil {
		t.Errorf("Del() on nonexistent key should not error, got: %v", err)
	}
}
