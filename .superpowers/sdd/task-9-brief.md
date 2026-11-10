# Task 9 Brief: Controller Base

## Context

You are implementing Task 9 of 16. Tasks 1-8 are complete.

**Your task builds on:** types.AppError, types.Response, types.PaginatedResponse

## Your Task

Implement base controller with common HTTP handling utilities.

**Working directory:** `C:\workspace\go_project\bhz`

## Files to Create

1. **`server/controller/base_controller.go`** — Base controller
2. **`server/controller/base_controller_test.go`** — Tests

## Requirements

### BaseController

```go
type BaseController struct{}

// Response methods
func (bc *BaseController) Success(c *gin.Context, data interface{})
func (bc *BaseController) SuccessWithMessage(c *gin.Context, message string, data interface{})
func (bc *BaseController) Created(c *gin.Context, data interface{})
func (bc *BaseController) BadRequest(c *gin.Context, message string, err error)
func (bc *BaseController) NotFound(c *gin.Context, message string)
func (bc *BaseController) Unauthorized(c *gin.Context, message string)
func (bc *BaseController) Forbidden(c *gin.Context, message string)
func (bc *BaseController) InternalError(c *gin.Context, message string, err error)

// Parameter helpers
func (bc *BaseController) ParamInt(c *gin.Context, key string) int
func (bc *BaseController) QueryInt(c *gin.Context, key string, defaultValue int) int

// Binding helpers
func (bc *BaseController) BindJSON(c *gin.Context, obj interface{}) error
func (bc *BaseController) BindQuery(c *gin.Context, obj interface{}) error

// Pagination
func (bc *BaseController) Paginate(c *gin.Context, items interface{}, total int64, page, perPage int)

// No content
func (bc *BaseController) NoContent(c *gin.Context)
```

### Implementation Notes

- Success uses `types.SuccessResponse(c, data)`
- Created uses `c.JSON(201, ...)` with Success: true
- Error responses use `types.NewAppError()` then `types.ErrorResponse()`
- ParamInt: `strconv.Atoi(c.Param(key))` (ignore error, return 0)
- QueryInt: `c.Query(key)`, parse or return default
- Created returns HTTP 201, all others return HTTP 200

### Test Requirements

```go
func TestBaseController_Success(t *testing.T) {
    gin.SetMode(gin.TestMode)
    c, _ := gin.CreateTestContext(nil)
    c.Request, _ = http.NewRequest("GET", "/test", nil)
    
    bc := &BaseController{}
    bc.Success(c, "test data")
    
    assert.Equal(t, 200, c.Writer.Status())
}

func TestBaseController_BadRequest(t *testing.T) {
    // Test 200 status with error message
}

func TestBaseController_NotFound(t *testing.T) {
    // Test 404 response
}

func TestBaseController_InternalError(t *testing.T) {
    // Test 500 response
}
```

Use `github.com/stretchr/testify/assert`.

## Global Constraints

- Go 1.25+
- Import types: `github.com/hxl1513/ffgoapi/types`

## Test Verification

```bash
cd server && go test ./controller -v
```

Expected: All tests PASS

## Commit Message

```
feat(controller): add base controller with HTTP utilities

- Implement BaseController with common HTTP response helpers
- Add Success, Created, BadRequest, NotFound, InternalError methods
- Add parameter binding helpers (ParamInt, QueryInt)
- Add pagination helper
- Add unit tests

Co-Authored-By: Claude Code <noreply@anthropic.com>
```

## Report

Write your report to: `C:\workspace\go_project\bhz\.superpowers\sdd\task-9-report.md`
