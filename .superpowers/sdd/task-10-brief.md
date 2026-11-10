# Task 10 Brief: Service Base

## Context

You are implementing Task 10 of 16. Tasks 1-9 are complete.

**Your task builds on:** model.DB (Task 5), model.Cache (Task 6)

## Your Task

Implement base service with common database operations and transaction support.

**Working directory:** `C:\workspace\go_project\bhz`

## Files to Create

1. **`server/service/base_service.go`** — Base service
2. **`server/service/base_service_test.go`** — Tests

## Requirements

### BaseService

```go
type BaseService struct {
    db    *model.DB
    cache model.Cache
}

func NewBaseService(db *model.DB) *BaseService
func (s *BaseService) WithCache(cache model.Cache) *BaseService
func (s *BaseService) DB() *model.DB
func (s *BaseService) Cache() model.Cache

func (s *BaseService) WithTransaction(ctx context.Context, fn func(tx *model.DB) error) error
```

### WithTransaction Implementation

```go
func (s *BaseService) WithTransaction(ctx context.Context, fn func(tx *model.DB) error) error {
    tx := s.db.Write.WithContext(ctx).Begin()
    if tx.Error != nil {
        return fmt.Errorf("failed to begin transaction: %w", tx.Error)
    }
    
    if err := fn(tx); err != nil {
        if rbErr := tx.Rollback().Error; rbErr != nil {
            return fmt.Errorf("transaction error: %v, rollback error: %w", err, rbErr)
        }
        return err
    }
    
    if err := tx.Commit().Error; err != nil {
        return fmt.Errorf("failed to commit transaction: %w", err)
    }
    
    return nil
}
```

### Test Requirements

```go
func TestBaseService_DB(t *testing.T) {
    cfg := &model.DatabaseConfig{
        Type:   "sqlite",
        DBName: ":memory:",
    }
    db, _ := model.NewDB(cfg)
    defer db.Close()
    
    svc := NewBaseService(db)
    if svc.DB() == nil { t.Error("DB() returned nil") }
}

func TestBaseService_Transaction(t *testing.T) {
    cfg := &model.DatabaseConfig{
        Type:   "sqlite",
        DBName: ":memory:",
    }
    db, _ := model.NewDB(cfg)
    defer db.Close()
    
    svc := NewBaseService(db)
    
    // Success
    err := svc.WithTransaction(context.Background(), func(tx *model.DB) error {
        return nil
    })
    if err != nil { t.Errorf("WithTransaction() error = %v", err) }
    
    // Rollback
    testErr := &testingError{msg: "test error"}
    err = svc.WithTransaction(context.Background(), func(tx *model.DB) error {
        return testErr
    })
    if err != testErr { t.Errorf("WithTransaction() error = %v, want %v", err, testErr) }
}

type testingError struct{ msg string }
func (e *testingError) Error() string { return e.msg }
```

## Global Constraints

- Go 1.25+
- Import model: `github.com/hxl1513/ffgoapi/model`

## Test Verification

```bash
cd server && go test ./service -v
```

Expected: All tests PASS

## Commit Message

```
feat(service): add base service with transaction support

- Implement BaseService with database access
- Add WithTransaction for atomic operations
- Add optional cache support
- Add unit tests

Co-Authored-By: Claude Code <noreply@anthropic.com>
```

## Report

Write your report to: `C:\workspace\go_project\bhz\.superpowers\sdd\task-10-report.md`
