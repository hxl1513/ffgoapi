# Task 4 Brief: Logger

## Context

You are implementing Task 4 of 16. Tasks 1-3 are complete.

**Working directory:** `C:\workspace\go_project\bhz`

## Files to Create

1. **`server/common/logger.go`** — Logger implementation
2. **`server/common/logger_test.go`** — Logger tests
3. **`server/pkg/crypto/hash.go`** — SHA-256 hashing
4. **`server/pkg/crypto/jwt.go`** — JWT token generation/parsing

## Requirements

### Logger

```go
type Logger struct {
    infoLogger  *log.Logger
    warnLogger  *log.Logger
    errorLogger *log.Logger
    debugLogger *log.Logger
    level       string
    mu          sync.Mutex
}

const (
    LogLevelDebug = "debug"
    LogLevelInfo  = "info"
    LogLevelWarn  = "warn"
    LogLevelError = "error"
)

func NewLogger(output io.Writer, level string) *Logger
func (l *Logger) Info(msg string, args ...any)
func (l *Logger) Warn(msg string, args ...any)
func (l *Logger) Error(msg string, args ...any)
func (l *Logger) Debug(msg string, args ...any)

// NewFileLogger creates logger writing to file + console
func NewFileLogger(filePath string, level string) (*Logger, *os.File, error)
```

**Log format:** `[LEVEL] timestamp | message`

**Level filtering:** Only log at or above configured level (debug < info < warn < error).

### SHA-256 Hashing

```go
// HashSHA256 hashes a string using SHA-256
func HashSHA256(s string) string
```

### JWT

```go
type Claims struct {
    UserID   int    `json:"user_id"`
    Username string `json:"username"`
    Role     int    `json:"role"`
    jwt.RegisteredClaims
}

func GenerateToken(secret string, userID int, username string, role int, expiry time.Duration) (string, error)
func ParseToken(secret, tokenString string) (*Claims, error)
```

## Test Requirements

### logger_test.go

```go
func TestNewLogger(t *testing.T) {
    buf := &bytes.Buffer{}
    logger := NewLogger(buf, "info")
    if logger == nil { t.Fatal("NewLogger() returned nil") }
}

func TestLoggerLevels(t *testing.T) {
    buf := &bytes.Buffer{}
    logger := NewLogger(buf, "debug")
    
    logger.Info("info message")
    logger.Warn("warn message")
    logger.Error("error message")
    logger.Debug("debug message")
    
    output := buf.String()
    if !contains(output, "info message") { t.Error(...) }
    if !contains(output, "warn message") { t.Error(...) }
    if !contains(output, "error message") { t.Error(...) }
    if !contains(output, "debug message") { t.Error(...) }
}

func TestLoggerInfoLevelFiltersDebug(t *testing.T) {
    buf := &bytes.Buffer{}
    logger := NewLogger(buf, "info")
    
    logger.Debug("should not appear")
    
    if contains(buf.String(), "should not appear") {
        t.Error("Debug message should not appear at info level")
    }
}
```

## Global Constraints

- Go 1.25+
- Thread-safe logger (use mutex for writes)
- Import path for jwt: `github.com/golang-jwt/jwt/v5`

## Test Verification

```bash
cd server && go test ./common -v && go test ./pkg/crypto -v
```

Expected: All tests PASS

## Commit Message

```
feat(common): add logger and crypto utilities

- Implement Logger with level filtering (debug/info/warn/error)
- Add file + console logging support
- Implement SHA-256 hashing utility
- Add JWT token generation and parsing
- Add unit tests for logger

Co-Authored-By: Claude Code <noreply@anthropic.com>
```

## Report

Write your report to: `C:\workspace\go_project\bhz\.superpowers\sdd\task-4-report.md`
