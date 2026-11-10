# Getting Started

This guide walks you through setting up a new project with the Go Web API Framework, from installation to a running server with a working endpoint.

## 1. Installation

Add the framework as a module dependency in your Go project:

```bash
go mod init your-org/your-app
go get github.com/hxl1513/ffgoapi
```

The framework requires **Go 1.22** or later.

## 2. Configuration

All framework settings are loaded from environment variables. Create a `.env` file (or export them in your shell):

```env
# Server
PORT=3000
GIN_MODE=release
READ_TIMEOUT=10s
WRITE_TIMEOUT=10s

# Database (SQLite by default — no external service needed)
DB_TYPE=sqlite
DB_NAME=app.db

# Optional: MySQL / PostgreSQL
# DB_TYPE=mysql
# DB_HOST=localhost
# DB_PORT=3306
# DB_USER=root
# DB_PASSWORD=
# DB_NAME=myapp

# Redis (optional — omit to use in-memory cache)
# REDIS_ENABLED=true
# REDIS_URL=redis://localhost:6379

# Authentication (set a strong secret in production)
SESSION_SECRET=your-secret-key-here
TOKEN_EXPIRY=86400

# Rate limiting (enabled by default)
RATE_LIMIT_ENABLED=true
RATE_LIMIT_DEFAULT_MAX=100
RATE_LIMIT_DEFAULT_DURATION=60

# Logging
LOG_LEVEL=info
LOG_DIR=./logs
```

Load configuration in your `main` function with `config.LoadFromEnv()`.

### Configuration Reference

| Variable | Type | Default | Description |
|---|---|---|---|
| `PORT` | int | `3000` | HTTP listen port |
| `GIN_MODE` | string | `release` | Gin mode: `debug` or `release` |
| `READ_TIMEOUT` | duration | `10s` | HTTP read timeout |
| `WRITE_TIMEOUT` | duration | `10s` | HTTP write timeout |
| `DB_TYPE` | string | `sqlite` | Database type: `sqlite`, `mysql`, `postgres` |
| `DB_NAME` | string | `app.db` | Database name or file path |
| `SESSION_SECRET` | string | *(required)* | JWT signing secret |
| `TOKEN_EXPIRY` | int | `86400` | JWT token expiry in seconds |
| `LOG_LEVEL` | string | `info` | Log level: `debug`, `info`, `warn`, `error` |

## 3. Creating Your First Endpoint

A minimal `main.go` using the framework:

```go
package main

import (
	"log"

	"github.com/hxl1513/ffgoapi/config"
	"github.com/hxl1513/ffgoapi/server"
)

func main() {
	cfg, err := config.LoadFromEnv()
	if err != nil {
		log.Fatal(err)
	}

	srv, err := server.NewServer(cfg)
	if err != nil {
		log.Fatal(err)
	}

	srv.RegisterRoutes(func(r *gin.RouterGroup) {
		r.GET("/health", func(c *gin.Context) {
			c.JSON(200, gin.H{"status": "ok"})
		})

		r.GET("/hello/:name", func(c *gin.Context) {
			name := c.Param("name")
			c.JSON(200, gin.H{"greeting": "Hello, " + name + "!"})
		})
	})

	if err := srv.Start(); err != nil {
		log.Fatal(err)
	}
}
```

### Using the Base Controller

For production routes, extend `BaseController` to get JSON helpers without repeating response formatting:

```go
package controller

import (
	"github.com/hxl1513/ffgoapi/controller"
	"github.com/hxl1513/ffgoapi/dto"
	"github.com/hxl1513/ffgoapi/types"
)

type GreetingController struct{}

func NewGreetingController() *GreetingController {
	return &GreetingController{}
}

func (gc *GreetingController) SayHello(c *controller.BaseController, req dto.GreetingRequest) {
	c.Success(c, types.Response{
		Success: true,
		Data:    "Hello, " + req.Name + "!",
	})
}
```

## 4. Adding a Model

Define your data model as a GORM-compatible struct in the `model` package:

```go
package model

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Name     string `json:"name" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"-" validate:"required,min=8"`
}
```

Create a corresponding DTO for request/response payloads:

```go
package dto

type CreateUserRequest struct {
	Name     string `json:"name" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
}

type UserResponse struct {
	ID    uint   `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}
```

Implement a service that handles business logic backed by the database:

```go
package service

import (
	"context"

	"github.com/hxl1513/ffgoapi/model"
)

type UserService struct {
	*BaseService
}

func NewUserService(db *model.DB) *UserService {
	return &UserService{BaseService: NewBaseService(db)}
}

func (us *UserService) CreateUser(ctx context.Context, user *model.User) error {
	return us.DB().Write.WithContext(ctx).Create(user).Error
}
```

Wire everything up in your route registration:

```go
srv.RegisterRoutes(func(r *gin.RouterGroup) {
	userController := controller.NewUserController(userService)
	r.POST("/users", userController.Create)
	r.GET("/users/:id", userController.GetByID)
	r.GET("/users", userController.List)
})
```

## 5. Running the Server

```bash
go run .
```

The server listens on the port defined by `PORT` (default `3000`).

### Health Check

```bash
curl http://localhost:3000/health
```

Expected response:

```json
{"status":"ok"}
```

### Structured Logging

The framework writes structured logs to `LOG_DIR` (default `./logs`). Set `LOG_LEVEL=debug` to see per-request details including request IDs, HTTP method, path, status code, and latency.

### Graceful Shutdown

When the server receives a `SIGINT` or `SIGTERM` signal it stops accepting new connections and waits for in-flight requests to complete before shutting down.

## Next Steps

- Explore the **[user CRUD example](examples/README.md)** for a complete working implementation.
- Read the [design document](../../docs/superpowers/specs/2026-11-06-go-web-api-framework-design.md) for a full architectural overview.
- Add [authentication middleware](middleware/auth.go) to protect routes.
- Add [rate limiting](middleware/rate-limit.go) to defend against abuse.
- Swap the cache backend to Redis by setting `REDIS_ENABLED=true` in the environment.
