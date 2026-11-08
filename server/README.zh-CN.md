# Go Web API 框架

一个基于 [Gin](https://github.com/gin-gonic/gin) 构建的生产级 Go Web API 框架，采用清晰的分层架构，支持多数据库，并为安全性、可观测性和性能提供开箱即用的配置。

## ✨ 特性

- **🚀 高性能** — 基于 Gin 框架，提供低延迟的请求处理能力
- **🏗️ 分层架构** — 控制器、服务、数据层职责清晰分离
- **🔧 环境配置** — 通过环境变量管理所有配置，支持 `.env` 文件
- **🛡️ 安全可靠** — 内置 JWT 认证、请求限流、CORS 支持
- **📊 可观测性** — 结构化日志、请求 ID 追踪、 panic 恢复
- **🗄️ 多数据库支持** — 通过 GORM 支持 SQLite、MySQL、PostgreSQL
- **💾 灵活缓存** — 可插拔的 Cache 接口，支持 Redis 和内存实现，Redis 不可用时自动降级
- **✅ 请求验证** — 集成 `go-playground/validator` 进行请求体验证

## 🚀 快速开始

### 安装

```bash
go mod init your-org/your-app
go get github.com/your-org/go-web-api-framework
```

**要求：** Go 1.22 或更高版本

### 基础配置

创建 `.env` 文件（或在 shell 中导出环境变量）：

```env
# 服务器配置
PORT=3000
GIN_MODE=release
READ_TIMEOUT=10s
WRITE_TIMEOUT=10s

# 数据库配置（默认 SQLite — 无需额外服务）
DB_TYPE=sqlite
DB_NAME=app.db

# 可选：MySQL / PostgreSQL
# DB_TYPE=mysql
# DB_HOST=localhost
# DB_PORT=3306
# DB_USER=root
# DB_PASSWORD=
# DB_NAME=myapp

# Redis（可选 — 留空使用内存缓存）
# REDIS_ENABLED=true
# REDIS_URL=redis://localhost:6379

# 认证配置（生产环境请使用强密钥）
SESSION_SECRET=your-secret-key-here
TOKEN_EXPIRY=86400

# 限流配置（默认启用）
RATE_LIMIT_ENABLED=true
RATE_LIMIT_DEFAULT_MAX=100
RATE_LIMIT_DEFAULT_DURATION=60

# 日志配置
LOG_LEVEL=info
LOG_DIR=./logs
```

### 第一个 API

```go
package main

import (
	"log"

	"github.com/your-org/go-web-api-framework/config"
	"github.com/your-org/go-web-api-framework/server"
)

func main() {
	// 1. 加载配置
	cfg, err := config.LoadFromEnv()
	if err != nil {
		log.Fatal(err)
	}

	// 2. 创建服务器
	srv, err := server.NewServer(cfg)
	if err != nil {
		log.Fatal(err)
	}

	// 3. 注册路由
	srv.RegisterRoutes(func(r *gin.RouterGroup) {
		r.GET("/hello", func(c *gin.Context) {
			c.JSON(200, gin.H{"message": "Hello, World!"})
		})
	})

	// 4. 启动服务
	srv.Start()
}
```

运行：

```bash
go run main.go
```

服务器将在 `http://localhost:3000` 启动。访问 `http://localhost:3000/hello` 验证服务正常。

## 📁 项目结构

```
server/
├── config/         # 配置结构和环境变量加载
├── controller/     # HTTP 控制器（内置 JSON 辅助方法）
├── middleware/      # 认证、CORS、限流、请求 ID、日志、恢复
├── model/          # 数据库和缓存抽象（GORM + Cache 接口）
├── service/        # 基础服务（支持事务）
├── dto/            # 请求和响应 DTO
├── types/          # 共享类型（AppError、Response、PaginatedResponse）
├── pkg/            # 加密（JWT、哈希）和验证工具
├── examples/       # 完整示例（用户 CRUD）
├── docs/           # 入门指南和 API 文档
└── go.mod
```

## 🗺️ 架构概览

```
┌─────────────┐
│   Client    │
└──────┬──────┘
       │
┌──────▼──────────────────┐
│   Middleware Layer      │
│  • CORS                 │
│  • Request ID           │
│  • Logger               │
│  • Auth (JWT)           │
│  • Rate Limit           │
│  • Request Body Limit   │
│  • Panic Recovery       │
└──────┬──────────────────┘
       │
┌──────▼──────────────────┐
│   Controller Layer      │  ← HTTP 处理
│  (BaseController)       │
└──────┬──────────────────┘
       │
┌──────▼──────────────────┐
│    Service Layer        │  ← 业务逻辑
│  (BaseService)          │
└──────┬──────────────────┘
       │
┌──────▼──────────────────┐
│     Model Layer         │  ← 数据访问
│  • Database (GORM)      │
│  • Cache (Redis/Mem)    │
│  • BaseModel            │
└─────────────────────────┘
```

## 📖 核心组件

### 配置管理 (`config`)

通过环境变量加载配置，支持 `.env` 文件：

```go
cfg, err := config.LoadFromEnv()
```

### 统一响应 (`types`)

泛型响应类型，简化 API 开发：

```go
// 成功响应
types.SuccessResponse(c, user)

// 错误响应
types.ErrorResponse(c, types.NewAppError(
    types.ErrCodeBadRequest,
    "invalid input",
    err,
))

// 分页响应
types.PaginatedResponse(c, users, total, page, perPage)
```

### 基础控制器 (`controller`)

嵌入 `BaseController` 快速构建控制器：

```go
type UserController struct {
    controller.BaseController
    service *UserService
}

func (c *UserController) GetUser(ctx *gin.Context) {
    id := c.ParamInt(ctx, "id")
    user, err := c.service.GetUser(ctx.Request.Context(), id)
    if err != nil {
        c.NotFound(ctx, "user not found")
        return
    }
    c.Success(ctx, user)
}
```

### 基础服务 (`service`)

提供数据库访问和事务支持：

```go
type UserService struct {
    service.BaseService
}

// 使用事务
err := svc.WithTransaction(ctx, func(tx *model.DB) error {
    return tx.Write.Create(&user).Error
})
```

### 数据模型 (`model`)

- **BaseModel** — 包含 ID、时间戳、软删除
- **Database** — 支持 SQLite / MySQL / PostgreSQL
- **Cache** — Redis 和内存实现，自动降级

### 中间件 (`middleware`)

- **RequestID** — 为每个请求生成唯一 ID
- **CORS** — 跨域资源共享
- **Logger** — 请求日志记录
- **Recover** — Panic 恢复
- **Auth** — JWT Bearer Token 认证
- **AdminAuth** — 仅管理员访问
- **UserAuth** — 仅普通用户访问
- **RateLimit** — 基于 IP 的请求限流
- **RequestBodyLimit** — 请求体大小限制

### JWT 认证 (`pkg/crypto`)

```go
// 生成 Token
token, err := crypto.GenerateToken(secret, userID, username, role, expiry)

// 解析 Token
claims, err := crypto.ParseToken(secret, tokenString)
```

### 请求验证 (`pkg/validator`)

```go
// 在控制器中使用
type CreateUserRequest struct {
    Name     string `json:"name" validate:"required,min=2,max=100"`
    Email    string `json:"email" validate:"required,email"`
    Password string `json:"password" validate:"required,min=8"`
}

func (c *UserController) CreateUser(ctx *gin.Context) {
    var req dto.CreateUserRequest
    if err := c.ValidateAndBind(ctx, &req); err != nil {
        c.BadRequest(ctx, "invalid request", err)
        return
    }
    // ...
}
```

## 🔐 认证与授权

### JWT Token

```go
// 配置
cfg := &config.AuthConfig{
    SessionSecret: "your-secret-key",
    TokenExpiry:   86400, // 秒
}

// 使用中间件
r.GET("/admin", AdminAuth(cfg, cache), adminHandler)
r.GET("/user", UserAuth(cfg, cache), userHandler)
```

### 角色说明

- `1` — 管理员（Admin）
- `2` — 普通用户（User）

## 🚦 限流

基于 Redis / 内存的滑动窗口限流：

```go
cfg := &config.RateLimitConfig{
    Enabled:         true,
    DefaultMax:      100,              // 最大请求数
    DefaultDuration: 60,               // 时间窗口（秒）
}
```

## 📝 完整示例

框架包含一个完整的用户 CRUD 示例，展示了如何结合使用各层组件。详见 [examples/](examples/) 目录。

运行示例：

```bash
cd examples
go test -v
```

## 🔧 开发工具命令

### 运行测试

```bash
cd server
go test ./...
```

### 构建

```bash
cd server
go build ./...
```

### 代码检查

```bash
go vet ./...
```

## 📚 更多文档

- **[入门指南](docs/GETTING_STARTED.md)** — 从零开始的详细教程
- **[设计文档](../../docs/superpowers/specs/2026-11-06-go-web-api-framework-design.md)** — 架构设计与技术选型
- **[示例代码](examples/)** — 用户 CRUD 完整示例

## 🤝 贡献

欢迎提交 Issue 和 Pull Request！

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

## 许可证

MIT License
