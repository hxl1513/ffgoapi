# Task 13 Brief: Validator Integration

## Context

You are implementing Task 13 of 16. Tasks 1-12 are complete.

## Your Task

Integrate validator package into controller base.

**Working directory:** `C:\workspace\go_project\bhz`

## Files to Create

1. **`server/pkg/validator/validate.go`** — Validator wrapper
2. **`server/pkg/validator/validate_test.go`** — Tests

## Requirements

### validate.go

```go
type Validator struct {
    validate *validator.Validate
}

func New() *Validator
func (v *Validator) ValidateStruct(s interface{}) error
func (v *Validator) Validate(i interface{}) error
func (v *Validator) ValidateField(field interface{}, tag string) error
```

### Update BaseController

Add to `server/controller/base_controller.go`:

```go
import "github.com/hxl1513/ffgoapi/pkg/validator"

func (bc *BaseController) ValidateAndBind(c *gin.Context, obj interface{}) error {
    if err := c.ShouldBind(obj); err != nil {
        return err
    }
    v := validator.New()
    return v.ValidateStruct(obj)
}
```

### Test Requirements

```go
func TestNewValidator(t *testing.T) {
    v := New()
    if v == nil { t.Fatal("New() returned nil") }
}

func TestValidateStruct(t *testing.T) {
    v := New()
    
    type TestStruct struct {
        Name  string `validate:"required,min=3"`
        Email string `validate:"required,email"`
    }
    
    valid := &TestStruct{Name: "John", Email: "john@example.com"}
    if err := v.ValidateStruct(valid); err != nil { t.Errorf(...) }
    
    invalid := &TestStruct{Name: "AB", Email: "invalid-email"}
    if err := v.ValidateStruct(invalid); err == nil { t.Error("Invalid passed") }
}
```

## Global Constraints

- Go 1.25+
- Import: `github.com/go-playground/validator/v10`

## Test Verification

```bash
cd server && go test ./pkg/validator -v && go test ./controller -v
```

Expected: All tests PASS

## Commit Message

```
feat(validator): integrate validation into controllers

- Add Validator wrapper for go-playground/validator
- Add ValidateAndBind method to BaseController
- Add unit tests for validator

Co-Authored-By: Claude Code <noreply@anthropic.com>
```

## Report

Write your report to: `C:\workspace\go_project\bhz\.superpowers\sdd\task-13-report.md`
