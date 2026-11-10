# Go Web API Framework

A production-ready Go Web API framework built on [Gin](https://github.com/gin-gonic/gin) with a clean layered architecture, multi-database support, and sensible defaults for security, observability, and performance.

## Features

- **High Performance** — Built on the Gin framework for low-latency request handling.
- **Layered Architecture** — Clear separation between controllers, services, and data layers.
- **Environment-Based Configuration** — All settings loaded from environment variables via `config.LoadFromEnv()`.
- **Secure by Default** — Built-in JWT authentication middleware, request rate limiting, and CORS support.
- **Observable** — Structured request logging, request ID propagation, and panic recovery middleware included.
- **Multi-Database Support** — First-class drivers for SQLite, MySQL, and PostgreSQL via GORM.
- **Flexible Caching** — Pluggable `Cache` interface with Redis and in-memory implementations; automatic in-memory fallback when Redis is unavailable.
- **Request Validation** — Integrated struct validation powered by `go-playground/validator`.

## Quick Start

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
		r.GET("/hello", func(c *gin.Context) {
			c.JSON(200, gin.H{"message": "Hello, World!"})
		})
	})

	srv.Start()
}
```

Set environment variables before running (or use the defaults):

```bash
export PORT=3000
export GIN_MODE=release
export DB_TYPE=sqlite
export DB_NAME=app.db
```

Then run:

```bash
go run main.go
```

The server starts on `http://localhost:3000`. Visit `http://localhost:3000/hello` to confirm it works.

## Project Structure

```
server/
├── config/         # Configuration structs and env-variable loader
├── controller/     # HTTP handlers (base controller with JSON helpers)
├── middleware/      # Auth, CORS, rate-limit, request-ID, logger, recover
├── model/          # Database and cache abstractions (GORM + Cache interface)
├── service/        # Base service with transaction support
├── dto/            # Request and response DTOs
├── types/          # Shared types (AppError, Response, PaginatedResponse)
├── pkg/            # Crypto (JWT, hashing) and validator utilities
├── examples/       # Working example (user CRUD)
├── docs/           # Getting started guide and API docs
└── go.mod
```

For a deeper architectural overview, see the [design document](../../docs/superpowers/specs/2026-11-06-go-web-api-framework-design.md).

## 数据库迁移

框架基于 [golang-migrate](https://github.com/golang-migrate/migrate) 提供数据库迁移功能，支持 MySQL 8.0 和 PostgreSQL 16。

### 快速开始

在 `.env` 文件中添加：

```env
MIGRATION_ENABLED=true
MIGRATION_DIR=./migrations
```

创建迁移文件：

```sql
-- migrations/20240101120000_create_users.up.sql
CREATE TABLE users (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    username VARCHAR(50) NOT NULL UNIQUE,
    email VARCHAR(255) NOT NULL UNIQUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

### 执行迁移

```bash
go run cmd/migrate/main.go up
```

详细文档：[MIGRATION.md](docs/MIGRATION.md)

## License

MIT
