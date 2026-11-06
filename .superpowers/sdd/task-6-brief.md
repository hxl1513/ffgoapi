# Task 6 Brief: Cache Abstraction

## Context

You are implementing Task 6 of 16. Tasks 1-5 are complete.

**Your task builds on:** model.DB from Task 5

## Your Task

Implement cache abstraction with in-memory and Redis backends, with auto-fallback.

**Working directory:** `C:\workspace\go_project\bhz`

## Files to Create

1. **`server/model/cache.go`** — Cache interface and implementations
2. **`server/model/cache_test.go`** — Cache tests

## Requirements

### Cache Interface

```go
type Cache interface {
    Get(key string) (string, error)
    Set(key string, value string, expiration time.Duration) error
    Del(key string) error
    Exists(key string) (bool, error)
    TTL(key string) (time.Duration, error)
    HSet(key string, field string, value interface{}) error
    HGet(key string, field string) (string, error)
    HGetAll(key string) (map[string]string, error)
    Incr(key string, delta int64) error
}
```

### MemoryCache Implementation

- Thread-safe map with RWMutex
- TTL support (check on Get)
- Simple key-value storage

### RedisCache Implementation

- Use `github.com/go-redis/redis/v8`
- Ping to test connection on init
- Return error if Redis unavailable (let caller fall back to memory)

### NewCache Factory Function

```go
func NewCache(redisURL string) (Cache, error) {
    if redisURL == "" {
        return NewMemoryCache(), nil
    }
    // Try Redis, fall back to memory on error
}
```

### Test Requirements

```go
func TestMemoryCache(t *testing.T) {
    cache := NewMemoryCache()
    
    // Set and Get
    cache.Set("key1", "value1", 1*time.Minute)
    val, _ := cache.Get("key1")
    if val != "value1" { t.Errorf(...) }
    
    // Exists
    exists, _ := cache.Exists("key1")
    if !exists { t.Error("should exist") }
    
    // Del
    cache.Del("key1")
    _, err := cache.Get("key1")
    if err == nil { t.Error("expected error after delete") }
}

func TestMemoryCacheTTL(t *testing.T) {
    cache := NewMemoryCache()
    cache.Set("key1", "value1", 100*time.Millisecond)
    time.Sleep(150 * time.Millisecond)
    _, err := cache.Get("key1")
    if err == nil { t.Error("expected error after TTL") }
}

func TestMemoryCacheHSetHGet(t *testing.T) {
    cache := NewMemoryCache()
    cache.HSet("hash1", "field1", "value1")
    val, _ := cache.HGet("hash1", "field1")
    if val != "value1" { t.Errorf(...) }
}

func TestMemoryCacheIncr(t *testing.T) {
    cache := NewMemoryCache()
    cache.Set("counter", "10", 0)
    cache.Incr("counter", 5)
    val, _ := cache.Get("counter")
    if val != "15" { t.Errorf(...) }
}
```

## Global Constraints

- Go 1.25+
- Thread-safe (use sync.RWMutex)
- Redis client: github.com/go-redis/redis/v8

## Test Verification

```bash
cd server && go test ./model -v -run TestMemoryCache
```

Expected: All cache tests PASS

## Commit Message

```
feat(model): add cache abstraction layer

- Define Cache interface for unified cache access
- Implement MemoryCache with TTL support
- Implement RedisCache with auto-fallback to MemoryCache
- Add HSet/HGet/HGetAll operations
- Add Incr with automatic initialization
- Add comprehensive unit tests

Co-Authored-By: Claude Code <noreply@anthropic.com>
```

## Report

Write your report to: `C:\workspace\go_project\bhz\.superpowers\sdd\task-6-report.md`
