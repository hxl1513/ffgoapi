# Task 3 Brief: Common Types and Utilities

## Context

You are implementing Task 3 of 16 for the go-web-api-framework project. Tasks 1-2 are complete.

## Your Task

Define common types, error handling, utility functions, and constants used across the framework.

**Working directory:** `C:\workspace\go_project\bhz`

## Files to Create

1. **`server/types/common.go`** — Common types
2. **`server/types/errors.go`** — Error types and error codes
3. **`server/types/response.go`** — Response structures
4. **`server/types/errors_test.go`** — Error tests
5. **`server/common/utils.go`** — Utility functions
6. **`server/common/constants.go`** — Constants

## Requirements (verbatim from plan)

### Error Codes (server/types/errors.go)

```go
const (
    ErrCodeBadRequest     = 4001
    ErrCodeUnauthorized   = 4002
    ErrCodeForbidden      = 4003
    ErrCodeNotFound       = 4004
    ErrCodeValidation     = 4005
    ErrCodeInternalServer = 5001
    ErrCodeDatabase       = 5002
    ErrCodeCache          = 5003
)
```

### AppError Structure

```go
type AppError struct {
    Code    int    `json:"code"`
    Message string `json:"message"`
    Err     error  `json:"-"`
}

func NewAppError(code int, message string, err error) *AppError
func (e *AppError) Error() string
func (e *AppError) Unwrap() error
```

### Response Structures (server/types/response.go)

```go
type Response[T any] struct {
    Success bool   `json:"success"`
    Message string `json:"message"`
    Data    T      `json:"data,omitempty"`
}

type PaginatedResponse[T any] struct {
    Success bool       `json:"success"`
    Data    []T       `json:"data"`
    Total   int64     `json:"total"`
    Page    int       `json:"page"`
    PerPage int       `json:"per_page"`
}
```

### Response Helper Functions

```go
func SuccessResponse[T any](c *gin.Context, data T)
func SuccessResponseWithMessage[T any](c *gin.Context, message string, data T)
func ErrorResponse(c *gin.Context, err *AppError)
func PaginatedResponse[T any](c *gin.Context, data []T, total int64, page, perPage int)
```

### Common Type

```go
// server/types/common.go
type Empty struct{}
```

### Utilities (server/common/utils.go)

```go
func IsEmpty(s string) bool
func IsNotEmpty(s string) bool
```

### Constants (server/common/constants.go)

```go
const (
    DefaultPageSize = 10
    MaxPageSize     = 100
)
```

### Test Requirements

Write `server/types/errors_test.go`:

```go
func TestNewAppError(t *testing.T) {
    err := errors.New("underlying error")
    appErr := NewAppError(ErrCodeNotFound, "resource not found", err)

    if appErr.Code != ErrCodeNotFound { t.Errorf(...) }
    if appErr.Message != "resource not found" { t.Errorf(...) }
    if !errors.Is(appErr, err) { t.Error("AppError should wrap underlying error") }
}

func TestNewAppErrorWithoutWrapping(t *testing.T) {
    appErr := NewAppError(ErrCodeBadRequest, "invalid input", nil)
    // Verify no nil pointer dereference
}

func TestAppError_Unwrap(t *testing.T) {
    err := errors.New("wrapped error")
    appErr := NewAppError(ErrCodeDatabase, "db error", err)
    unwrapped := appErr.Unwrap()
    if unwrapped != err { t.Errorf(...) }
}
```

## Global Constraints

- Go version: 1.25+
- Package naming: lowercase, single word where possible
- Error handling: use `fmt.Errorf` with `%w` verb for wrapping

## Test Verification

```bash
cd server && go test ./types -v && go test ./common -v
```

Expected: All tests PASS

## Report

Write your report to: `C:\workspace\go_project\bhz\.superpowers\sdd\task-3-report.md`

Include:
1. Status: DONE, DONE_WITH_CONCERNS, NEEDS_CONTEXT, or BLOCKED
2. Commits created
3. Test summary
4. Any concerns
5. Self-review checklist
