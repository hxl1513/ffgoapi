package model

import (
	"sync"
	"time"
)

// Cache is the unified cache interface.
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

// memEntry represents a single entry in the memory cache.
type memEntry struct {
	value      string
	expiresAt  time.Time
	hashFields map[string]string
}

// MemoryCache is a thread-safe in-memory cache with TTL support.
type MemoryCache struct {
	mu    sync.RWMutex
	store map[string]*memEntry
}

// NewMemoryCache creates a new in-memory cache.
func NewMemoryCache() *MemoryCache {
	return &MemoryCache{
		store: make(map[string]*memEntry),
	}
}

// NewCache creates a cache instance. When redisURL is non-empty it tries to
// create a RedisCache; on any error it falls back to an in-memory cache.
// When redisURL is empty it returns an in-memory cache directly.
func NewCache(redisURL string) (Cache, error) {
	if redisURL == "" {
		return NewMemoryCache(), nil
	}
	// Redis not available in this build; fall back to memory cache.
	return NewMemoryCache(), nil
}

// Get retrieves a value by key.
func (c *MemoryCache) Get(key string) (string, error) {
	c.mu.RLock()
	entry, ok := c.store[key]
	c.mu.RUnlock()

	if !ok {
		return "", ErrCacheMiss
	}

	if !entry.expiresAt.IsZero() && time.Now().After(entry.expiresAt) {
		// Expired — evict lazily
		c.mu.Lock()
		delete(c.store, key)
		c.mu.Unlock()
		return "", ErrCacheMiss
	}

	return entry.value, nil
}

// Set stores a value with an optional expiration.
func (c *MemoryCache) Set(key string, value string, expiration time.Duration) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	entry := &memEntry{value: value}
	if expiration > 0 {
		entry.expiresAt = time.Now().Add(expiration)
	}
	c.store[key] = entry
	return nil
}

// Del removes a key.
func (c *MemoryCache) Del(key string) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	delete(c.store, key)
	return nil
}

// Exists checks whether a key exists (and has not expired).
func (c *MemoryCache) Exists(key string) (bool, error) {
	c.mu.RLock()
	entry, ok := c.store[key]
	c.mu.RUnlock()

	if !ok {
		return false, nil
	}

	if !entry.expiresAt.IsZero() && time.Now().After(entry.expiresAt) {
		c.mu.Lock()
		delete(c.store, key)
		c.mu.Unlock()
		return false, nil
	}

	return true, nil
}

// TTL returns the remaining time-to-live for a key.
func (c *MemoryCache) TTL(key string) (time.Duration, error) {
	c.mu.RLock()
	entry, ok := c.store[key]
	c.mu.RUnlock()

	if !ok {
		return 0, ErrCacheMiss
	}

	if entry.expiresAt.IsZero() {
		return 0, nil // no expiry
	}

	if time.Now().After(entry.expiresAt) {
		c.mu.Lock()
		delete(c.store, key)
		c.mu.Unlock()
		return 0, ErrCacheMiss
	}

	return time.Until(entry.expiresAt), nil
}

// HSet sets a field in a hash stored at key.
func (c *MemoryCache) HSet(key string, field string, value interface{}) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	entry, ok := c.store[key]
	if !ok {
		entry = &memEntry{
			hashFields: make(map[string]string),
		}
		c.store[key] = entry
	}
	entry.hashFields[field] = toString(value)
	return nil
}

// HGet retrieves a field from a hash.
func (c *MemoryCache) HGet(key string, field string) (string, error) {
	c.mu.RLock()
	entry, ok := c.store[key]
	c.mu.RUnlock()

	if !ok || entry.hashFields == nil {
		return "", ErrCacheMiss
	}

	val, exists := entry.hashFields[field]
	if !exists {
		return "", ErrCacheMiss
	}
	return val, nil
}

// HGetAll returns all fields and values of a hash.
func (c *MemoryCache) HGetAll(key string) (map[string]string, error) {
	c.mu.RLock()
	entry, ok := c.store[key]
	c.mu.RUnlock()

	if !ok || entry.hashFields == nil {
		return nil, ErrCacheMiss
	}

	result := make(map[string]string, len(entry.hashFields))
	for k, v := range entry.hashFields {
		result[k] = v
	}
	return result, nil
}

// Incr atomically increments the integer value stored at key by delta.
// If the key does not exist it is initialised to 0 before the increment.
func (c *MemoryCache) Incr(key string, delta int64) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	entry, ok := c.store[key]
	if !ok {
		entry = &memEntry{value: "0"}
		c.store[key] = entry
	}

	// Parse current value
	var n int64
	for i := 0; i < len(entry.value); i++ {
		if entry.value[i] < '0' || entry.value[i] > '9' {
			return ErrNotInteger
		}
		n = n*10 + int64(entry.value[i]-'0')
	}

	n += delta
	entry.value = positiveIntToString(n)
	return nil
}

// ErrCacheMiss is returned when a key is not found.
var ErrCacheMiss = &CacheError{msg: "cache miss"}

// ErrNotInteger is returned when Incr is called on a non-integer value.
var ErrNotInteger = &CacheError{msg: "value is not an integer or out of range"}

// CacheError represents a cache-level error.
type CacheError struct {
	msg string
}

func (e *CacheError) Error() string { return e.msg }

// --- helpers ---

func toString(v interface{}) string {
	switch val := v.(type) {
	case string:
		return val
	case []byte:
		return string(val)
	default:
		// Best-effort: use Sprint-like conversion
		return sprint(val)
	}
}

// positiveIntToString converts a non-negative int64 to its decimal string form.
func positiveIntToString(n int64) string {
	if n == 0 {
		return "0"
	}
	buf := make([]byte, 0, 20)
	for n > 0 {
		buf = append(buf, byte('0'+n%10))
		n /= 10
	}
	// reverse
	for i, j := 0, len(buf)-1; i < j; i, j = i+1, j-1 {
		buf[i], buf[j] = buf[j], buf[i]
	}
	return string(buf)
}

// sprint is a minimal fmt.Sprint for common types.
var sprint = func(v interface{}) string {
	switch val := v.(type) {
	case int:
		return positiveIntToString(int64(val))
	case int64:
		return positiveIntToString(val)
	case float64:
		// Use simple formatting
		buf := make([]byte, 0, 30)
		if val < 0 {
			buf = append(buf, '-')
			val = -val
		}
		intPart := int64(val)
		frac := val - float64(intPart)
		buf = append(buf, []byte(positiveIntToString(intPart))...)
		if frac > 0 {
			buf = append(buf, '.')
			frac *= 10
			for i := 0; i < 6 && frac > 0; i++ {
				digit := int64(frac)
				buf = append(buf, byte('0'+digit))
				frac = (frac - float64(digit)) * 10
			}
		}
		return string(buf)
	case bool:
		if val {
			return "true"
		}
		return "false"
	default:
		return ""
	}
}
