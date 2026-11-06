# Task 2 Brief: Configuration Management

## Context

You are implementing Task 2 of 16 for the go-web-api-framework project. Task 1 (Project Initialization) is complete — the Go module is set up at `C:\workspace\go_project\bhz\server\go.mod`.

## Your Task

Implement type-safe configuration loading with environment variable support in `C:\workspace\go_project\bhz\server`.

**Working directory:** `C:\workspace\go_project\bhz`

## Files to Create

1. **`server/config/env.go`** — Environment variable helper functions
2. **`server/config/config.go`** — Configuration structures and LoadFromEnv()
3. **`server/config/config_test.go`** — Unit tests

## Requirements (verbatim from plan)

### Config Structures

```go
type Config struct {
    Server   ServerConfig
    Database DatabaseConfig
    Redis    RedisConfig
    Auth     AuthConfig
    RateLimit RateLimitConfig
    Logger   LoggerConfig
}

type ServerConfig struct {
    Port         int           `json:"port"`
    Mode         string        `json:"mode"`
    ReadTimeout  time.Duration `json:"read_timeout"`
    WriteTimeout time.Duration `json:"write_timeout"`
}

type DatabaseConfig struct {
    Type         string `json:"type"`
    Host         string `json:"host"`
    Port         int    `json:"port"`
    User         string `json:"user"`
    Password     string `json:"password"`
    DBName       string `json:"dbname"`
    MaxIdleConns int    `json:"max_idle_conns"`
    MaxOpenConns int    `json:"max_open_conns"`
    LogLevel     string `json:"log_level"`
}

type RedisConfig struct {
    Enabled  bool   `json:"enabled"`
    URL      string `json:"url"`
    PoolSize int    `json:"pool_size"`
}

type AuthConfig struct {
    SessionSecret string `json:"session_secret"`
    TokenExpiry   int    `json:"token_expiry"`
}

type RateLimitConfig struct {
    Enabled         bool `json:"enabled"`
    DefaultMax      int  `json:"default_max"`
    DefaultDuration int  `json:"default_duration"`
}

type LoggerConfig struct {
    Level string `json:"level"`
    Dir   string `json:"dir"`
}
```

### Environment Variable Defaults

```go
Server:
  Port:         getEnvAsInt("PORT", 3000)
  Mode:         getEnv("GIN_MODE", "release")
  ReadTimeout:  getEnvAsDuration("READ_TIMEOUT", 10*time.Second)
  WriteTimeout: getEnvAsDuration("WRITE_TIMEOUT", 10*time.Second)

Database:
  Type:         getEnv("DB_TYPE", "sqlite")
  Host:         getEnv("DB_HOST", "localhost")
  Port:         getEnvAsInt("DB_PORT", 3306)
  User:         getEnv("DB_USER", "root")
  Password:     getEnv("DB_PASSWORD", "")
  DBName:       getEnv("DB_NAME", "app.db")
  MaxIdleConns: getEnvAsInt("DB_MAX_IDLE_CONNS", 10)
  MaxOpenConns: getEnvAsInt("DB_MAX_OPEN_CONNS", 100)
  LogLevel:     getEnv("DB_LOG_LEVEL", "error")

Redis:
  Enabled:  getEnvAsBool("REDIS_ENABLED", false)
  URL:      getEnv("REDIS_CONN_STRING", "")
  PoolSize: getEnvAsInt("REDIS_POOL_SIZE", 10)

Auth:
  SessionSecret: getEnv("SESSION_SECRET", "")
  TokenExpiry:   getEnvAsInt("TOKEN_EXPIRY", 86400)

RateLimit:
  Enabled:         getEnvAsBool("RATE_LIMIT_ENABLED", true)
  DefaultMax:      getEnvAsInt("RATE_LIMIT_DEFAULT_MAX", 100)
  DefaultDuration: getEnvAsInt("RATE_LIMIT_DEFAULT_DURATION", 60)

Logger:
  Level: getEnv("LOG_LEVEL", "info")
  Dir:   getEnv("LOG_DIR", "./logs")
```

### Helper Functions (env.go)

```go
func getEnv(key, defaultValue string) string
func getEnvAsInt(key string, defaultValue int) int
func getEnvAsBool(key string, defaultValue bool) bool
func getEnvAsDuration(key string, defaultValue time.Duration) time.Duration
```

### Validation

```go
func ValidateConfig(cfg *Config) error
```

Returns error if `Auth.SessionSecret` is empty.

### Test Requirements

Write tests in `server/config/config_test.go`:

1. **TestLoadFromEnv** — Set environment variables, verify they're loaded correctly
2. **TestLoadFromEnvDefaults** — Verify default values when env vars not set
3. **TestValidateConfig** — Test validation (valid config passes, missing secret fails)

## Global Constraints

- Go version: 1.25+
- Package naming: lowercase, single word where possible
- Error handling: use `fmt.Errorf` with `%w` verb for wrapping

## Test Verification

After implementation, run:
```bash
cd server && go test ./config -v
```

Expected: All tests PASS

## Report

Write your report to: `C:\workspace\go_project\bhz\.superpowers\sdd\task-2-report.md`

Include:
1. Status: DONE, DONE_WITH_CONCERNS, NEEDS_CONTEXT, or BLOCKED
2. Commits created (git log --oneline -3)
3. Test summary: exact test command output
4. Any concerns or observations
5. Self-review checklist results
