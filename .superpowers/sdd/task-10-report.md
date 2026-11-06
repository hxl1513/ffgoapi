# Task 10 Report: Service Base

## Status: DONE

## Summary

Implemented `BaseService` in `server/service/base_service.go` with full database and optional cache support, plus transaction handling.

## Files Created

- **`C:\workspace\go_project\bhz\server\service\base_service.go`** — Production code
- **`C:\workspace\go_project\bhz\server\service\base_service_test.go`** — Tests

## Implementation Details

### BaseService struct
- Embeds `*model.DB` for database access and an optional `model.Cache` interface.
- `NewBaseService(db *model.DB)` — constructor.
- `WithCache(cache model.Cache)` — fluent setter, returns `*BaseService` for chaining.
- `DB() / Cache()` — simple accessors.

### WithTransaction
- Begins a write transaction via `s.db.Write.WithContext(ctx).Begin()`.
- Wraps the returned `*gorm.DB` into `*model.DB{Write: tx, Read: tx}` before passing to the callback.
- Rolls back and returns the original callback error if `fn` returns an error; includes rollback error as a wrapped context if the rollback itself fails.
- Commits on success; returns a wrapped error on commit failure.

## Test Results

```
=== RUN   TestBaseService_DB
--- PASS: TestBaseService_DB (0.00s)
=== RUN   TestBaseService_Transaction
--- PASS: TestBaseService_Transaction (0.00s)
PASS
ok  	github.com/your-org/go-web-api-framework/service	2.019s
```

All tests pass on Go 1.25.0.

## Commit

```
a2b01da feat(service): add base service with transaction support

- Implement BaseService with database access
- Add WithTransaction for atomic operations
- Add optional cache support
- Add unit tests

Co-Authored-By: Claude Code <noreply@anthropic.com>
```

## Notes

- The only build fix required was wrapping the `*gorm.DB` returned by `.Begin()` into `*model.DB` before passing to the callback — `model.DB` is a struct wrapper, not an alias.
- `go.mod` declares `go 1.22`; the brief says Go 1.25+. The existing `go.mod` was left untouched (Task 1–9 already set it to 1.22 and it compiled fine under Go 1.25 runtime).
