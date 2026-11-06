# Task 12 Brief: DTOs and Integration Test

## Context

You are implementing Task 12 of 16. Tasks 1-11 are complete.

**Your task builds on:** All previous components

## Your Task

Create DTOs and an example user CRUD implementation to demonstrate framework usage.

**Working directory:** `C:\workspace\go_project\bhz`

## Files to Create

1. **`server/dto/request.go`** — Request DTOs
2. **`server/dto/response.go`** — Response DTOs  
3. **`server/dto/common.go`** — DTO conversion helpers
4. **`server/examples/users/user_model.go`** — Example user model
5. **`server/examples/users/user_service.go`** — Example user service
6. **`server/examples/users/user_controller.go`** — Example user controller
7. **`server/examples/main.go`** — Example main.go
8. **`server/examples/main_test.go`** — Integration test

## Requirements

### request.go

```go
type CreateUserRequest struct {
    Name     string `json:"name" validate:"required,min=2,max=100"`
    Email    string `json:"email" validate:"required,email"`
    Password string `json:"password" validate:"required,min=8"`
}

func (r *CreateUserRequest) Validate() error
```

### response.go

```go
type UserResponse struct {
    ID        uint   `json:"id"`
    Name      string `json:"name"`
    Email     string `json:"email"`
    Role      int    `json:"role"`
    CreatedAt string `json:"created_at"`
}
```

### common.go

```go
func FromUser(user *model.User) *UserResponse
```

### Example User Model

```go
type User struct {
    model.BaseModel
    Name     string `json:"name" gorm:"not null"`
    Email    string `json:"email" gorm:"uniqueIndex;not null"`
    Password string `json:"-" gorm:"not null"`
    Role     int    `json:"role" gorm:"default:2"`
}
```

### Example User Service

```go
type UserRepository struct {
    db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository
func (r *UserRepository) Create(ctx context.Context, user *User) error
func (r *UserRepository) FindByID(ctx context.Context, id int) (*User, error)

type UserService struct {
    service.BaseService
    repo *UserRepository
}

func NewUserService(db *model.DB) *UserService
func (s *UserService) CreateUser(ctx context.Context, req *dto.CreateUserRequest) (*dto.UserResponse, error)
```

Include `hashPassword` helper (simplified: `return "hashed_" + password`).

### Example User Controller

```go
type UserController struct {
    controller.BaseController
    service *UserService
}

func NewUserController(service *UserService) *UserController
func (c *UserController) RegisterRoutes(r *gin.RouterGroup)
func (c *UserController) CreateUser(ctx *gin.Context)
func (c *UserController) GetUser(ctx *gin.Context)
```

### Example main.go

Show how to wire everything together:
```go
func main() {
    cfg, _ := config.LoadFromEnv()
    srv, _ := server.NewServer(cfg)
    defer srv.Shutdown(context.Background())
    
    srv.RegisterRoutes(func(r *gin.RouterGroup) {
        userSvc := users.NewUserService(srv.DB())
        userCtrl := users.NewUserController(userSvc)
        userCtrl.RegisterRoutes(r)
    })
    
    srv.Start()
}
```

## Test Requirements

Integration test that creates a user and verifies it can be retrieved.

## Global Constraints

- Go 1.25+
- Follow existing patterns

## Test Verification

```bash
cd server && go test ./examples -v
```

## Commit Message

```
feat(dto): add data transfer objects and examples

- Define request and response DTOs
- Add validation integration
- Create example user CRUD implementation
- Add integration test demonstrating framework usage

Co-Authored-By: Claude Code <noreply@anthropic.com>
```

## Report

Write your report to: `C:\workspace\go_project\bhz\.superpowers\sdd\task-12-report.md`
