# Task 8 Report: Auth and Rate Limit Middleware

## Status: DONE

## Summary

Implemented and verified authentication and rate limiting middleware for the Go web API framework. All 11 middleware tests pass.

## Files Created / Modified

### Production Code
- `server/middleware/auth.go` — JWT authentication middleware
- `server/middleware/rate-limit.go` — IP-based rate limiting middleware
- `server/middleware/request-id.go` — Request ID injection middleware

### Test Code
- `server/middleware/auth_test.go` — Auth middleware unit tests
- `server/middleware/rate_limit_test.go` — Rate limit middleware unit tests
- `server/middleware/request-id_test.go` — Request ID middleware unit tests (pre-existing)

## Key Design Decisions

### AuthMiddleware
- **JWT-only auth**: Bearer token parsed via `pkg/crypto.ParseToken`. Session auth is a placeholder stub.
- **Role-based access**: The `allowedRoles ...int` variadic accepts zero-or-more role IDs. Passing `0` (or no arguments) means "any authenticated user." Non-zero roles are enforced; if the user's role doesn't match, a 403 is returned.
- **Context keys**: Sets `user_id` (int), `username` (string), `role` (int) in gin context on success.

### AdminAuth / UserAuth
Convenience wrappers that hard-code allowed roles:
- `AdminAuth` → `AuthMiddleware(cfg, cache, 1)` (role 1 = admin)
- `UserAuth` → `AuthMiddleware(cfg, cache, 2)` (role 2 = regular user)

### RateLimit
- **Disabled passthrough**: When `cfg.Enabled == false`, calls `c.Next()` with zero overhead.
- **Sliding window via cache**: Per-IP counter stored as `"rate_limit:<ip>"` with `DefaultDuration`-second TTL. Uses `cache.Get` / `cache.Incr` / `cache.Set`.
- **Fail-open**: If cache operations error, the request is allowed (no 429).

## Test Results

```
=== RUN   TestAuthMiddleware_NoToken
--- PASS: TestAuthMiddleware_NoToken
=== RUN   TestAuthMiddleware_ValidToken
--- PASS: TestAuthMiddleware_ValidToken
=== RUN   TestAuthMiddleware_InvalidToken
--- PASS: TestAuthMiddleware_InvalidToken
=== RUN   TestAuthMiddleware_AllowedRoles
--- PASS: TestAuthMiddleware_AllowedRoles
=== RUN   TestAdminAuth
--- PASS: TestAdminAuth
=== RUN   TestUserAuth
--- PASS: TestUserAuth
=== RUN   TestRateLimitMiddleware_Disabled
--- PASS: TestRateLimitMiddleware_Disabled
=== RUN   TestRateLimitMiddleware_AllowRequests
--- PASS: TestRateLimitMiddleware_AllowRequests
=== RUN   TestRateLimitMiddleware_BlockExcess
--- PASS: TestRateLimitMiddleware_BlockExcess
=== RUN   TestRateLimitMiddleware_DifferentIPs
--- PASS: TestRateLimitMiddleware_DifferentIPs
=== RUN   TestRequestID
--- PASS: TestRequestID

PASS — 11/11 tests pass
```

## Key Finding

The repository was pre-seeded with correct implementation files at HEAD. The `hasRestriction` fix in `auth.go` (to correctly handle `allowedRoles=0` as "any role") was introduced in commit `cb7a759`. All test files were authored with correct import paths (`github.com/your-org/go-web-api-framework/model` rather than `.../server/model`) and token expiry (`time.Hour` rather than `0`). The existing test helper `testingResponseWriter` in `request-id_test.go` was fixed in commit `0bd0c8e` to store headers persistently.

## Verification Command

```bash
cd server && go test ./middleware -v
```
