# Go Web API 基础框架设计文档

**项目名称**: go-web-api-framework
**版本**: v1.0.0
**日期**: 2026-11-06
**目标**: 从 new-api 项目抽取基础组件，构建通用 Go Web API 服务框架

---

## 1. 项目概述

### 1.1 项目目标

从 `C:\workspace\go_project\new-api` 中抽取对外提供 API 服务的基础组件，构建一个**标准企业级的 Go 语言 Web API 基础框架**。

该框架适用于：
- RESTful API 服务
- BFF（Backend For Frontend）
- 内部管理系统后端
- 微服务单体

### 1.2 功能范围

**标准企业级**包含以下组件：
- ✅ HTTP 服务器封装（基于 Gin）
- ✅ 配置管理（环境变量 + 文件配置）
- ✅ 统一日志系统
- ✅ 统一错误处理和响应格式
- ✅ 基础中间件（CORS、限流、认证、日志、Panic恢复等）
- ✅ 数据库抽象层（GORM + MySQL/PostgreSQL/SQLite/ClickHouse）
- ✅ 缓存抽象层（内存 + Redis，自动降级）
- ✅ 工具函数库
- ✅ 数据验证（validator）

---

## 2. 整体架构

### 2.1 设计原则

1. **分层架构**：Controller → Service → Model 清晰分层
2. **依赖倒置**：高层模块不依赖低层模块，都依赖抽象
3. **单一职责**：每个组件职责明确，易于测试和维护
4. **约定优于配置**：合理的默认配置，减少重复代码
5. **向后兼容**：框架升级不影响业务代码

### 2.2 层次说明

| 层次 | 职责 | 依赖 |
|------|------|------|
| **Controller** | 处理 HTTP 请求/响应、参数校验、调用 Service | Service |
| **Service** | 业务逻辑实现、事务管理 | Model |
| **Model** | 数据访问、数据库操作 | Database/Cache |
| **Middleware** | 横切关注点（认证、限流、日志等） | - |
| **Config** | 配置管理、环境变量 | - |
| **Common** | 工具函数、常量、类型定义 | - |

### 2.3 核心原则

- ✅ Controller 依赖 Service
- ✅ Service 依赖 Model
- ✅ Model 依赖 Database/Cache
- ❌ Controller 不直接依赖 Database
- ❌ Service 不依赖 Gin Context（通过参数传递）

---

## 3. 目录结构

```
server/                                 # 框架根目录
├── server.go                          # 主入口：初始化 + 启动HTTP服务器
├── router.go                          # 路由注册中心
├── middleware/                        # 中间件集合
│   ├── auth.go                       # 认证（Session + Token）
│   ├── cors.go                       # CORS 支持
│   ├── i18n.go                       # 国际化
│   ├── logger.go                     # 请求日志记录
│   ├── rate-limit.go                 # 限流（IP级 + 用户级）
│   ├── recover.go                    # Panic 恢复
│   ├── request-id.go                 # 请求ID生成
│   ├── request-body-limit.go         # 请求体大小限制
│   └── body-cleanup.go               # 请求体清理
├── controller/                        # 控制器层
│   └── base_controller.go            # 基础控制器（通用方法）
├── service/                           # 业务逻辑层
│   └── base_service.go               # 基础服务（通用方法）
├── model/                             # 数据模型层
│   ├── base_model.go                 # 基础模型（公共字段）
│   ├── database.go                    # 数据库配置和初始化
│   ├── cache.go                       # 缓存接口和实现
│   └── paginate.go                    # 分页辅助
├── dto/                               # 数据传输对象
│   ├── request.go                     # 请求 DTO
│   ├── response.go                    # 响应 DTO
│   └── common.go                      # 通用 DTO
├── config/                            # 配置管理
│   ├── config.go                      # 配置结构定义
│   ├── manager.go                     # 配置管理器
│   └── env.go                         # 环境变量工具
├── common/                            # 通用工具
│   ├── logger.go                      # 日志封装
│   ├── utils.go                       # 工具函数
│   ├── constants.go                   # 常量定义
│   ├── types.go                       # 公共类型
│   └── helpers.go                     # 辅助函数
├── pkg/                               # 第三方工具封装
│   ├── crypto/                        # 加密工具
│   │   ├── hash.go
│   │   └── jwt.go
│   └── validator/                     # 验证器
│       └── validate.go
├── types/                             # 全局公共类型
│   ├── errors.go                      # 错误类型定义
│   ├── response.go                    # 响应类型定义
│   └── common.go                      # 通用类型
└── go.mod                             # Go 模块定义

```

---

## 4. 核心组件详细设计

### 4.1 HTTP 服务器 (`server.go`)

**职责**：
- 初始化 Gin 引擎
- 加载配置
- 初始化所有组件（数据库、缓存、日志等）
- 注册路由
- 启动 HTTP 服务器
- 优雅关闭

**核心结构**：

```go
type Server struct {
    engine *gin.Engine
    config *config.Config
    logger *common.Logger
    db     *model.DB
    cache  model.Cache
}

func NewServer(cfg *config.Config) (*Server, error)
func (s *Server) Engine() *gin.Engine
func (s *Server) Start() error
func (s *Server) Shutdown(ctx context.Context) error
```

**初始化流程**：

```
1. 加载配置
2. 初始化日志
3. 初始化数据库
4. 初始化缓存
5. 注册中间件
6. 注册路由
7. 启动服务器
```

---

### 4.2 配置管理 (`config/`)

**职责**：
- 统一管理所有配置项
- 支持环境变量、配置文件加载
- 提供类型安全的配置访问

**核心结构**：

```go
// Config 应用配置
type Config struct {
    Server   ServerConfig
    Database DatabaseConfig
    Redis    RedisConfig
    Auth     AuthConfig
    RateLimit RateLimitConfig
    Logger   LoggerConfig
}

// ServerConfig 服务器配置
type ServerConfig struct {
    Port         int           `json:"port" env:"PORT" envDefault:"3000"`
    Mode         string        `json:"mode" env:"GIN_MODE" envDefault:"release"`
    ReadTimeout  time.Duration `json:"read_timeout" envDefault:"10s"`
    WriteTimeout time.Duration `json:"write_timeout" envDefault:"10s"`
}

// DatabaseConfig 数据库配置
type DatabaseConfig struct {
    Type         string `json:"type" env:"DB_TYPE" envDefault:"sqlite"`
    Host         string `json:"host" env:"DB_HOST"`
    Port         int    `json:"port" env:"DB_PORT"`
    User         string `json:"user" env:"DB_USER"`
    Password     string `json:"password" env:"DB_PASSWORD"`
    DBName       string `json:"dbname" env:"DB_NAME"`
    MaxIdleConns int    `json:"max_idle_conns" envDefault:"10"`
    MaxOpenConns int    `json:"max_open_conns" envDefault:"100"`
    LogLevel     string `json:"log_level" envDefault:"error"`
}

// RedisConfig Redis 配置
type RedisConfig struct {
    Enabled  bool          `json:"enabled" env:"REDIS_ENABLED" envDefault:"false"`
    URL      string        `json:"url" env:"REDIS_CONN_STRING"`
    PoolSize int           `json:"pool_size" env:"REDIS_POOL_SIZE" envDefault:"10"`
}

// AuthConfig 认证配置
type AuthConfig struct {
    SessionSecret string `json:"session_secret" env:"SESSION_SECRET" required:"true"`
    TokenExpiry   int    `json:"token_expiry" env:"TOKEN_EXPIRY" envDefault:"86400"` // 秒
}

// RateLimitConfig 限流配置
type RateLimitConfig struct {
    Enabled          bool `json:"enabled" env:"RATE_LIMIT_ENABLED" envDefault:"true"`
    DefaultMax       int  `json:"default_max" env:"RATE_LIMIT_DEFAULT_MAX" envDefault:"100"`
    DefaultDuration  int  `json:"default_duration" env:"RATE_LIMIT_DEFAULT_DURATION" envDefault:"60"` // 秒
}
```

**配置加载优先级**：
1. 环境变量（最高）
2. 配置文件（config.yaml / config.json）
3. 默认值（最低）

---

### 4.3 统一日志 (`common/logger.go`)

**职责**：
- 统一的日志格式
- 支持文件 + 控制台输出
- 自动日志轮转
- 请求 ID 追踪

**核心结构**：

```go
type Logger struct {
    info  *log.Logger
    warn  *log.Logger
    error *log.Logger
    debug *log.Logger
    file  *os.File
}

func NewLogger(output io.Writer, level string) *Logger
func (l *Logger) Info(msg string, args ...any)
func (l *Logger) Warn(msg string, args ...any)
func (l *Logger) Error(msg string, args ...any)
func (l *Logger) Debug(msg string, args ...any)
```

**日志格式**：
```
[INFO] 2026/11/06 - 15:04:05 | REQUEST_ID | message
[ERROR] 2026/11/06 - 15:04:05 | REQUEST_ID | message
```

---

### 4.4 统一错误处理 (`types/errors.go`)

**职责**：
- 统一的错误类型定义
- 错误码映射
- 错误响应格式化

**核心结构**：

```go
// AppError 应用错误
type AppError struct {
    Code    int    `json:"code"`              // 业务错误码
    Message string `json:"message"`            // 错误信息
    Err     error  `json:"-"`                  // 原始错误
}

// 常见错误码
const (
    ErrCodeBadRequest     = 4001 // 请求参数错误
    ErrCodeUnauthorized   = 4002 // 未认证
    ErrCodeForbidden      = 4003 // 无权限
    ErrCodeNotFound       = 4004 // 资源不存在
    ErrCodeValidation     = 4005 // 数据验证失败
    ErrCodeInternalServer = 5001 // 服务器内部错误
    ErrCodeDatabase       = 5002 // 数据库错误
    ErrCodeCache          = 5003 // 缓存错误
)

func NewAppError(code int, message string, err error) *AppError
func (e *AppError) Error() string
func (e *AppError) Unwrap() error
```

**统一响应格式**：

```go
// 成功响应
{
    "success": true,
    "message": "success",
    "data": { ... }
}

// 错误响应
{
    "success": false,
    "message": "error message",
    "code": 4001
}
```

---

### 4.5 统一响应 (`types/response.go`)

```go
// Response 统一响应结构
type Response[T any] struct {
    Success bool   `json:"success"`
    Message string `json:"message"`
    Data    T      `json:"data,omitempty"`
}

// PaginatedResponse 分页响应
type PaginatedResponse[T any] struct {
    Success bool       `json:"success"`
    Data    []T       `json:"data"`
    Total   int64     `json:"total"`
    Page    int       `json:"page"`
    PerPage int       `json:"per_page"`
}

// 响应辅助函数
func Success[T any](c *gin.Context, data T)
func SuccessWithMessage(c *gin.Context, message string, data any)
func Error(c *gin.Context, err *AppError)
func Paginated[T any](c *gin.Context, data []T, total int64, page, perPage int)
```

---

### 4.6 中间件设计

#### 4.6.1 认证中间件 (`middleware/auth.go`)

```go
// AuthMiddleware 认证中间件
// 支持 Session 认证和 Token 认证
func AuthMiddleware(roles ...int) gin.HandlerFunc

// AdminAuth 管理员认证
func AdminAuth() gin.HandlerFunc

// UserAuth 用户认证
func UserAuth() gin.HandlerFunc
```

**认证流程**：
```
1. 检查 Session
2. 检查 Authorization Header (Bearer Token)
3. 验证 Token
4. 检查角色权限
5. 设置用户信息到 Context
```

#### 4.6.2 限流中间件 (`middleware/rate-limit.go`)

```go
// RateLimit IP 级限流
func RateLimit(max int, duration int64) gin.HandlerFunc

// UserRateLimit 用户级限流（需在 Auth 之后）
func UserRateLimit(max int, duration int64) gin.HandlerFunc
```

**特性**：
- 滑动窗口算法
- 支持 Redis 和内存两种后端
- 自动降级（Redis 不可用时使用内存）

#### 4.6.3 其他中间件

| 中间件 | 职责 |
|--------|------|
| `cors.go` | CORS 支持，可配置允许的源、方法、头 |
| `logger.go` | 请求日志记录（Method、Path、Status、Latency） |
| `recover.go` | Panic 恢复，防止服务崩溃 |
| `request-id.go` | 生成唯一请求 ID，便于追踪 |
| `request-body-limit.go` | 限制请求体大小 |
| `body-cleanup.go` | 清理请求体存储（内存/磁盘） |
| `i18n.go` | 国际化支持 |

---

### 4.7 数据库抽象 (`model/`)

#### 4.7.1 数据库配置 (`database.go`)

```go
type DatabaseConfig struct {
    Type         string        `json:"type"`
    Host         string        `json:"host"`
    Port         int           `json:"port"`
    User         string        `json:"user"`
    Password     string        `json:"password"`
    DBName       string        `json:"dbname"`
    MaxIdleConns int           `json:"max_idle_conns"`
    MaxOpenConns int           `json:"max_open_conns"`
    LogLevel     string        `json:"log_level"`
}

type DB struct {
    Write *gorm.DB
    Read  *gorm.DB  // 可选：只读副本
}

func NewDB(cfg *DatabaseConfig) (*DB, error)
func (db *DB) Close() error
```

**支持数据库**：
- SQLite（开发/测试）
- MySQL（生产）
- PostgreSQL（生产）
- ClickHouse（日志/分析）

#### 4.7.2 基础模型 (`base_model.go`)

```go
type BaseModel struct {
    ID        uint           `gorm:"primaryKey" json:"id"`
    CreatedAt time.Time      `json:"created_at"`
    UpdatedAt time.Time      `json:"updated_at"`
    DeletedAt gorm.DeletedAt `gorm:"index" json:"-"` // 软删除
}

type BaseService struct {
    db *gorm.DB
}

func NewBaseService(db *gorm.DB) *BaseService
func (s *BaseService) DB() *gorm.DB
```

#### 4.7.3 缓存抽象 (`cache.go`)

```go
type Cache interface {
    Get(key string) (string, error)
    Set(key string, value string, expiration time.Duration) error
    Del(key string) error
    Exists(key string) (bool, error)
    TTL(key string) (time.Duration, error)
    HSet(key string, field string, value interface{}) error
    HGet(key string, field string) (string, error)
    HGetAll(key string) (map[string]string, error)
    Incr(key string, delta int64) error
}

type MemoryCache struct{ /* ... */ }
type RedisCache struct{ /* ... */ }

func NewCache(cfg *RedisConfig) (Cache, error)
```

**特性**：
- 统一接口
- Redis 不可用时自动降级到内存
- 支持 Hash、Incr 等常用操作

---

### 4.8 数据验证 (`pkg/validator/`)

```go
type Validator struct {
    validate *validator.Validate
}

func NewValidator() *Validator
func (v *Validator) Validate(i interface{}) error
func (v *Validator) ValidateStruct(s interface{}) error
```

**集成到 Controller**：

```go
func (c *BaseController) ValidateAndBind(i interface{}) error {
    if err := c.ShouldBind(i); err != nil {
        return err
    }
    return c.Validator().ValidateStruct(i)
}
```

---

## 5. 核心流程设计

### 5.1 服务器启动流程

```
main.go
  → config.Load()                    // 1. 加载配置
  → server.New(cfg)                  // 2. 创建服务器实例
    → logger.New()                   // 3. 初始化日志
    → model.NewDB(cfg.Database)      // 4. 初始化数据库
    → model.NewCache(cfg.Redis)      // 5. 初始化缓存
  → server.SetupMiddleware()         // 6. 注册中间件
  → server.SetupRoutes()             // 7. 注册路由
  → server.Start()                   // 8. 启动服务
```

### 5.2 请求处理流程

```
HTTP Request
  → Middleware: RequestID             // 生成请求ID
  → Middleware: Logger                // 记录请求日志
  → Middleware: CORS                  // CORS 处理
  → Middleware: RateLimit             // 限流检查
  → Middleware: Auth                  // 认证
  → Middleware: BodyLimit             // 请求体限制
  → Controller: Bind & Validate      // 参数绑定和验证
  → Service: Business Logic          // 业务逻辑
  → Model: Database Operation        // 数据操作
  → Controller: Response             // 返回响应
  → Middleware: BodyCleanup          // 清理资源
```

---

## 6. 与源项目的映射关系

| 框架组件 | 源项目路径 | 迁移说明 |
|---------|-----------|---------|
| `server.go` | `main.go` | 抽取初始化逻辑 |
| `middleware/` | `middleware/` | 直接迁移（调整包路径） |
| `controller/base_controller.go` | `controller/` | 抽取通用方法 |
| `service/base_service.go` | `service/` | 抽取通用方法 |
| `model/` | `model/` + `common/database.go` | 整合数据库初始化 |
| `config/` | `setting/config/` | 重构配置管理 |
| `common/` | `common/` + `constant/` | 抽取工具函数和常量 |
| `types/` | `types/` | 错误和响应类型 |

---

## 7. 技术选型

| 组件 | 技术选型 | 版本要求 | 理由 |
|------|---------|---------|------|
| Web 框架 | Gin | v1.9+ | 高性能、生态成熟、中间件丰富 |
| 数据库 ORM | GORM | v1.25+ | 支持多数据库、生态完善 |
| 数据库驱动 | MySQL/PostgreSQL/SQLite/ClickHouse | - | 覆盖主要使用场景 |
| 缓存 | Redis + 内存 | - | 可选 Redis，无 Redis 自动降级 |
| 配置加载 | go-playground/validator + 环境变量 | v10+ | 强大灵活、支持自定义规则 |
| 日志 | 标准库 log + 自定义封装 | - | 轻量、可控 |
| 限流 | 滑动窗口算法 | - | 实现简单、效果良好 |
| 加密 | golang.org/x/crypto | - | 官方维护、安全性高 |

---

## 8. 使用示例

### 8.1 快速开始

```go
// main.go
package main

import (
    "log"
    "your-org/go-web-api-framework/server"
    "your-org/go-web-api-framework/config"
    "your-org/go-web-api-framework/controller"
)

func main() {
    // 1. 加载配置
    cfg, err := config.Load()
    if err != nil {
        log.Fatalf("Failed to load config: %v", err)
    }

    // 2. 创建服务器
    srv, err := server.New(cfg)
    if err != nil {
        log.Fatalf("Failed to create server: %v", err)
    }

    // 3. 注册自定义路由
    srv.RegisterRoutes(func(r *gin.RouterGroup) {
        r.GET("/hello", controller.Hello)
        r.POST("/users", controller.CreateUser)
    })

    // 4. 启动服务
    if err := srv.Start(); err != nil {
        log.Fatalf("Failed to start server: %v", err)
    }
}
```

### 8.2 Controller 示例

```go
// controller/user_controller.go
package controller

import (
    "net/http"
    "your-org/go-web-api-framework/service"
    "your-org/go-web-api-framework/dto"
    "your-org/go-web-api-framework/types"
)

type UserController struct {
    BaseController
    userService *service.UserService
}

func NewUserController(userService *service.UserService) *UserController {
    return &UserController{
        userService: userService,
    }
}

// GetUser 获取用户信息
// @Summary 获取用户信息
// @Description 根据 ID 获取用户信息
// @Tags users
// @Accept json
// @Produce json
// @Param id path int true "用户 ID"
// @Success 200 {object} types.Response[dto.UserResponse]
// @Failure 404 {object} types.Response[any]
// @Router /users/{id} [get]
func (c *UserController) GetUser(ctx *gin.Context) {
    id := c.ParamInt(ctx, "id")
    user, err := c.userService.GetByID(ctx, id)
    if err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            c.NotFound(ctx, "用户不存在")
            return
        }
        c.InternalError(ctx, "获取用户失败", err)
        return
    }
    c.Success(ctx, user)
}

// CreateUser 创建用户
func (c *UserController) CreateUser(ctx *gin.Context) {
    var req dto.CreateUserRequest
    if err := c.ValidateAndBind(ctx, &req); err != nil {
        c.BadRequest(ctx, "参数错误", err)
        return
    }

    user, err := c.userService.Create(ctx, &req)
    if err != nil {
        c.InternalError(ctx, "创建用户失败", err)
        return
    }
    c.Created(ctx, user)
}
```

### 8.3 Service 示例

```go
// service/user_service.go
package service

import (
    "your-org/go-web-api-framework/model"
    "your-org/go-web-api-framework/dto"
)

type UserService struct {
    BaseService
    userRepo *model.UserRepository
}

func NewUserService(db *model.DB) *UserService {
    return &UserService{
        userRepo: model.NewUserRepository(db.Write),
    }
}

func (s *UserService) GetByID(ctx context.Context, id int) (*dto.UserResponse, error) {
    user, err := s.userRepo.FindByID(ctx, id)
    if err != nil {
        return nil, err
    }
    return dto.FromUser(user), nil
}

func (s *UserService) Create(ctx context.Context, req *dto.CreateUserRequest) (*dto.UserResponse, error) {
    user := &model.User{
        Name:     req.Name,
        Email:    req.Email,
        Password: hashPassword(req.Password),
    }

    if err := s.userRepo.Create(ctx, user); err != nil {
        return nil, err
    }

    return dto.FromUser(user), nil
}
```

### 8.4 Model 示例

```go
// model/user.go
package model

import "gorm.io/gorm"

type User struct {
    BaseModel
    Name     string `json:"name" gorm:"not null"`
    Email    string `json:"email" gorm:"uniqueIndex;not null"`
    Password string `json:"-" gorm:"not null"`
    Role     int    `json:"role" gorm:"default:1"`
}

type UserRepository struct {
    db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
    return &UserRepository{db: db}
}

func (r *UserRepository) FindByID(ctx context.Context, id int) (*User, error) {
    var user User
    err := r.db.WithContext(ctx).First(&user, id).Error
    return &user, err
}

func (r *UserRepository) Create(ctx context.Context, user *User) error {
    return r.db.WithContext(ctx).Create(user).Error
}
```

---

## 9. 配置示例

### 9.1 环境变量配置 (`.env`)

```env
# 服务器配置
PORT=3000
GIN_MODE=release

# 数据库配置
DB_TYPE=mysql
DB_HOST=localhost
DB_PORT=3306
DB_USER=root
DB_PASSWORD=<SECRET_f5f68bae>ord
DB_NAME=myapp

# Redis 配置（可选）
REDIS_ENABLED=true
REDIS_CONN_STRING=redis://localhost:6379/0

# 认证配置
SESSION_SECRET=your-random-secret-key-here

# 限流配置
RATE_LIMIT_ENABLED=true
RATE_LIMIT_DEFAULT_MAX=100
RATE_LIMIT_DEFAULT_DURATION=60

# 日志配置
LOG_LEVEL=info
LOG_DIR=./logs
```

### 9.2 配置文件 (`config.yaml`)

```yaml
server:
  port: 3000
  mode: release
  read_timeout: 10s
  write_timeout: 10s

database:
  type: mysql
  host: localhost
  port: 3306
  user: root
  password: password
  dbname: myapp
  max_idle_conns: 10
  max_open_conns: 100
  log_level: error

redis:
  enabled: false
  url: redis://localhost:6379/0
  pool_size: 10

auth:
  session_secret: "your-secret-key"
  token_expiry: 86400

rate_limit:
  enabled: true
  default_max: 100
  default_duration: 60

logger:
  level: info
  dir: ./logs
```

---

## 10. 开发规范

### 10.1 Controller 规范

- ✅ Controller 只处理 HTTP 相关逻辑
- ✅ 参数绑定和验证放在 Controller
- ✅ 调用 Service 执行业务逻辑
- ❌ 不包含业务逻辑
- ❌ 不直接操作数据库

### 10.2 Service 规范

- ✅ Service 实现业务逻辑
- ✅ 接收 `context.Context` 参数
- ✅ 返回明确的错误类型
- ❌ 不依赖 Gin Context
- ❌ 不处理 HTTP 请求/响应

### 10.3 Model 规范

- ✅ Model 包含数据定义和数据访问方法
- ✅ 使用 Repository 模式封装数据库操作
- ✅ 所有方法接收 `context.Context`
- ❌ 不包含业务逻辑

### 10.4 错误处理规范

```go
// ✅ 使用 AppError
return types.NewAppError(types.ErrCodeNotFound, "用户不存在", err)

// ✅ 包装底层错误
if err != nil {
    return types.NewAppError(types.ErrCodeDatabase, "查询失败", err)
}

// ❌ 直接返回 error
return err

// ❌ 直接返回 string
return "error message"
```

---

## 11. 后续扩展

### 11.1 可选的扩展模块

| 模块 | 说明 | 优先级 |
|------|------|--------|
| 任务调度 | 基于 Cron 的定时任务 | 高 |
| 消息队列 | RabbitMQ / Kafka 集成 | 中 |
| 监控指标 | Prometheus + Grafana | 中 |
| 链路追踪 | OpenTelemetry | 中 |
| 灰度发布 | 特性开关、流量控制 | 低 |
| 限流算法 | 令牌桶、漏桶算法 | 低 |
| 分布式锁 | Redis 分布式锁 | 中 |

### 11.2 性能优化方向

1. **连接池优化**：数据库、Redis 连接池调优
2. **缓存策略**：二级缓存、缓存预热
3. **请求合并**：HTTP 请求合并减少数据库查询
4. **异步处理**：耗时任务异步化
5. **数据库优化**：读写分离、分库分表支持

---

## 12. 总结

本框架通过**分层架构**和**组件化设计**，提供了构建企业级 Go Web API 服务的完整基础设施。

**核心优势**：
1. **开箱即用**：包含 Web 服务开发所需的所有基础组件
2. **高度可扩展**：支持自定义中间件、插件化架构
3. **类型安全**：全链路类型安全，减少运行时错误
4. **易于测试**：依赖注入 + 接口抽象，单元测试友好
5. **生产就绪**：包含日志、监控、错误处理等生产级特性

**下一步**：
- 根据本文档编写详细实现计划
- 逐步实现各组件
- 编写单元测试和集成测试
- 编写使用文档和示例
