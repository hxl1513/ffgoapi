# Task 2 Implementer: Configuration Management

## Context

You are implementing Task 2 of 16 for the go-web-api-framework project.

**Completed before you:**
- Task 1: Project initialization — `server/go.mod` is created with all dependencies

**Your task builds on:** Task 1 (imports will be available)

## Your Task

Implement type-safe configuration loading with environment variable support.

**Working directory:** `C:\workspace\go_project\bhz`

## Files to Create

1. **`server/config/env.go`** — Environment variable helper functions
2. **`server/config/config.go`** — Configuration structures and LoadFromEnv()
3. **`server/config/config_test.go`** — Unit tests

## Exact Requirements

### env.go Helper Functions

```go
package config

import (
    "os"
    "strconv"
    "time"
)

func getEnv(key, defaultValue string) string
func getEnvAsInt(key string, defaultValue int) int
func getEnvAsBool(key string, defaultValue bool) bool
func getEnvAsDuration(key string, defaultValue time.Duration) time.Duration
```

### config.go Structures

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

func LoadFromEnv() (*Config, error)
func ValidateConfig(cfg *Config) error

type ValidationError struct {
    Field   string
    Message string
}

func (e *ValidationError) Error() string
```

### Default Values (in LoadFromEnv)

Use the helper functions with these defaults:

```
Server:
  Port:         3000
  Mode:         "release"
  ReadTimeout:  10s
  WriteTimeout: 10s

Database:
  Type:         "sqlite"
  Host:         "localhost"
  Port:         3306
  User:         "root"
  Password:     ""
  DBName:       "app.db"
  MaxIdleConns: 10
  MaxOpenConns: 100
  LogLevel:     "error"

Redis:
  Enabled:  false
  URL:      ""
  PoolSize: 10

Auth:
  SessionSecret: "" (empty - must be validated)
  TokenExpiry:   86400

RateLimit:
  Enabled:         true
  DefaultMax:      100
  DefaultDuration: 60

Logger:
  Level: "info"
  Dir:   "./logs"
```

### Validation

ValidateConfig must return error if `cfg.Auth.SessionSecret == ""`:

```go
func ValidateConfig(cfg *Config) error {
    if cfg.Auth.SessionSecret == "" {
        return &ValidationError{Field: "auth.session_secret", Message: "session secret is required"}
    }
    return nil
}
```

## Test Requirements

Create `server/config/config_test.go` with these exact tests:

```go
package config

import (
    "os"
    "testing"
    "time"
)

func TestLoadFromEnv(t *testing.T) {
    os.Clearenv()
    os.Setenv("PORT", "8080")
    os.Setenv("GIN_MODE", "test")
    os.Setenv("DB_TYPE", "sqlite")
    os.Setenv("DB_NAME", "test.db")
    os.Setenv("SESSION_SECRET", "test-secret")

    cfg, err := LoadFromEnv()
    if err != nil { t.Fatalf("LoadFromEnv() error = %v", err) }

    if cfg.Server.Port != 8080 { t.Errorf("Port = %d, want 8080", cfg.Server.Port) }
    if cfg.Server.Mode != "test" { t.Errorf("Mode = %s, want test", cfg.Server.Mode) }
    if cfg.Database.Type != "sqlite" { t.Errorf("Database.Type = %s, want sqlite", cfg.Database.Type) }
    if cfg.Database.DBName != "test.db" { t.Errorf("Database.DBName = %s, want test.db", cfg.Database.DBName) }
    if cfg.Auth.SessionSecret != "test-secret" { t.Errorf("Auth.SessionSecret = %s, want test-secret", cfg.Auth.SessionSecret) }
}

func TestLoadFromEnvDefaults(t *testing.T) {
    os.Clearenv()
    os.Setenv("SESSION_SECRET", "test-secret")

    cfg, err := LoadFromEnv()
    if err != nil { t.Fatalf("LoadFromEnv() error = %v", err) }

    if cfg.Server.Port != 3000 { t.Errorf("Port = %d, want 3000 (default)", cfg.Server.Port) }
    if cfg.Server.Mode != "release" { t.Errorf("Mode = %s, want release (default)", cfg.Server.Mode) }
    if cfg.Server.ReadTimeout != 10*time.Second { t.Errorf("ReadTimeout = %v, want 10s", cfg.Server.ReadTimeout) }
    if cfg.Database.Type != "sqlite" { t.Errorf("Database.Type = %s, want sqlite", cfg.Database.Type) }
    if cfg.RateLimit.DefaultMax != 100 { t.Errorf("RateLimit.DefaultMax = %d, want 100", cfg.RateLimit.DefaultMax) }
}

func TestValidateConfig(t *testing.T) {
    tests := []struct {
        name    string
        config  *Config
        wantErr bool
    }{
        {
            name: "valid config",
            config: &Config{
                Server:   ServerConfig{Port: 3000},
                Database: DatabaseConfig{Type: "sqlite", DBName: "test.db"},
                Auth:     AuthConfig{SessionSecret: "secret"},
            },
            wantErr: false,
        },
        {
            name: "missing session secret",
            config: &Config{
                Server:   ServerConfig{Port: 3000},
                Database: DatabaseConfig{Type: "sqlite", DBName: "test.db"},
                Auth:     AuthConfig{SessionSecret: ""},
            },
            wantErr: true,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            err := ValidateConfig(tt.config)
            if (err != nil) != tt.wantErr { t.Errorf("ValidateConfig() error = %v, wantErr %v", err, tt.wantErr) }
        })
    }
}
```

## Global Constraints

- Go version: 1.25+
- Package naming: lowercase
- Error handling: use `fmt.Errorf` with `%w` verb for wrapping

## TDD Process

1. Write the tests first (provided above)
2. Run tests to verify they fail
3. Implement the minimal code
4. Run tests to verify they pass
5. Commit

## Test Verification

```bash
cd server && go test ./config -v
```

Expected: All 3 tests PASS

## Commit Message

```
feat(config): implement configuration management

- Add Config struct with all configuration sections
- Implement LoadFromEnv() with environment variable support
- Add default values for all configuration options
- Implement config validation
- Add unit tests for environment loading and validation

Co-Authored-By: Claude Code <noreply@anthropic.com>
```

## Report

Write your report to: `C:\workspace\go_project\bhz\.superpowers\sdd\task-2-report.md`

Include:
1. Status: DONE, DONE_WITH_CONCERNS, NEEDS_CONTEXT, or BLOCKED
2. Commits created
3. Test summary (exact output from `go test ./config -v`)
4. Any concerns or observations
5. Self-review checklist
