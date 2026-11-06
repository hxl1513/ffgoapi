# Task 4 Report: Logger & Crypto Utilities

## Status: DONE

## Summary

Task 4 implements logging and cryptographic utilities for the go-web-api-framework project. All implementation files were already present in the repository from prior task work; this task's contributions were creating the test infrastructure, fixing the logger prefix format, and verifying all tests pass.

## Files Created / Verified

### `server/common/logger.go`
- `Logger` struct with four level-specific `log.Logger` instances (info, warn, error, debug)
- Thread-safe via `sync.Mutex` on every log call
- Level constants: `LogLevelDebug`, `LogLevelInfo`, `LogLevelWarn`, `LogLevelError`
- `NewLogger(output io.Writer, level string)` — constructor writing to any `io.Writer`
- `NewFileLogger(filePath, level)` — writes to both `os.Stdout` and a log file via `io.MultiWriter`
- Log format: `[LEVEL] RFC3339-timestamp | message`
- Level filtering: only logs at or above configured level (debug < info < warn < error)

### `server/common/utils.go`
- `IsEmpty(s string) bool` — true for empty or whitespace-only strings
- `IsNotEmpty(s string) bool` — inverse of `IsEmpty`

### `server/pkg/crypto/hash.go`
- `HashSHA256(s string) string` — SHA-256 hex digest (64-char lowercase hex)

### `server/pkg/crypto/jwt.go`
- `Claims` struct embedding `jwt.RegisteredClaims` with `UserID`, `Username`, `Role`
- `GenerateToken(secret, userID, username, role, expiry)` — HS256-signed JWT
- `ParseToken(secret, tokenString)` — validates and returns `*Claims`

### `server/common/logger_test.go`
- 7 tests: NewLogger creation, multi-level logging, debug-filter-at-info, error-filter, warn-filter, format prefix check, thread safety

### `server/pkg/crypto/crypto_test.go`
- 5 tests: SHA-256 correctness, token generation, token parsing, invalid-secret rejection, invalid-token rejection

### `server/common/utils_test.go`
- 2 tests (pre-existing): `IsEmpty` / `IsNotEmpty` table-driven tests

## Tests

```
=== RUN   TestNewLogger
--- PASS: TestNewLogger (0.00s)
=== RUN   TestLoggerLevels
--- PASS: TestLoggerLevels (0.01s)
=== RUN   TestLoggerInfoLevelFiltersDebug
--- PASS: TestLoggerInfoLevelFiltersDebug (0.00s)
=== RUN   TestLoggerErrorLevelFiltersInfoAndDebug
--- PASS: TestLoggerErrorLevelFiltersInfoAndDebug (0.00s)
=== RUN   TestLoggerWarnLevelFiltersDebug
--- PASS: TestLoggerWarnLevelFiltersDebug (0.00s)
=== RUN   TestLoggerFormat
--- PASS: TestLoggerFormat (0.00s)
=== RUN   TestLoggerThreadSafety
--- PASS: TestLoggerThreadSafety (0.00s)
=== RUN   TestIsEmpty
--- PASS: TestIsEmpty (0.00s)
=== RUN   TestIsNotEmpty
--- PASS: TestIsNotEmpty (0.00s)
PASS   common (9 tests)
=== RUN   TestHashSHA256
--- PASS: TestHashSHA256 (0.00s)
=== RUN   TestGenerateToken
--- PASS: TestGenerateToken (0.00s)
=== RUN   TestParseToken
--- PASS: TestParseToken (0.00s)
=== RUN   TestParseTokenInvalidSecret
--- PASS: TestParseTokenInvalidSecret (0.00s)
=== RUN   TestParseTokenInvalidToken
--- PASS: TestParseTokenInvalidToken (0.00s)
PASS   pkg/crypto (5 tests)
```

**Result: All 14 tests PASS**

## Commit

```
commit cc98d6532a0c84c6d2748fa206ac99ea2c74c068
Author: hxl1513 <1016155157@qq.com>
Date:   Fri Nov  6 12:54:28 2026 +0800

feat(common): add logger and crypto utilities

- Logger with level filtering (debug/info/warn/error)
- File + console logging support (NewFileLogger)
- Thread-safe logger with mutex protection
- SHA-256 hashing utility (HashSHA256)
- JWT token generation and parsing with Claims struct
- Unit tests for logger (level filtering, format, thread safety)
- Unit tests for crypto (hashing, JWT generate/parse/invalid)

Co-Authored-By: Claude Code <noreply@anthropic.com>
```

## Key Decisions

1. **Log prefix format**: `[LEVEL]` (no trailing space) — kept consistent with `[LEVEL] timestamp | message` pattern
2. **Thread safety**: `sync.Mutex` wraps every public logging method; tested with 10 concurrent goroutines
3. **JWT signing method**: HS256 (HMAC-SHA256) using `github.com/golang-jwt/jwt/v5` v5.2.0
4. **FormatArgs design**: Simple string-type-only variadic formatter to keep scope minimal
