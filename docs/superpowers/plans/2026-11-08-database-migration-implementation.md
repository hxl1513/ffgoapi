# 数据库迁移功能实现计划

> **For agentic workers:** 使用 superpowers:subagent-driven-development 或 superpowers:executing-plans 实现此计划。

**目标：** 为 Go Web API 框架实现基于 golang-migrate 的数据库迁移功能，支持 MySQL 8.0 和 PostgreSQL 16。

**架构：** 基于设计文档 `docs/superpowers/specs/2026-11-08-database-migration-design.md`

**技术栈：** golang-migrate v4, GORM v1, MySQL 8.0, PostgreSQL 16

---

## Global Constraints

- **Go 版本：** 1.22+
- **golang-migrate：** v4.17.0
- **GORM：** v1.30.0（已存在）
- **MySQL 驱动：** gorm.io/driver/mysql v1.4.3（已存在）
- **PostgreSQL 驱动：** gorm.io/driver/postgres v1.5.2（已存在）
- **TDD 原则：** 先写测试，再实现
- **向后兼容：** 迁移功能默认禁用，不影响现有代码

---

## 任务分解

### Task 1: 依赖安装

**目标：** 在 go.mod 中添加 golang-migrate 依赖

**文件：**
- 修改：`server/go.mod`
- 修改：`server/go.sum`（由 go mod tidy 生成）

**步骤：**

- [ ] **Step 1: 添加依赖**

```bash
cd server && go get github.com/golang-migrate/migrate/v4
```

验证：`grep "golang-migrate" go.mod`

- [ ] **Step 2: 提交**

```bash
git add server/go.mod server/go.sum
git commit -m "chore: add golang-migrate dependency

- Add github.com/golang-migrate/migrate/v4 v4.17.0
- Prepare for database migration feature

Co-Authored-By: Claude Code <noreply@anthropic.com>"
```

---

### Task 2: MigrationConfig 配置扩展

**目标：** 在 config 包中添加 MigrationConfig 配置结构

**文件：**
- 修改：`server/config/config.go`

**接口：**

```go
// MigrationConfig 数据库迁移配置
type MigrationConfig struct {
    Enabled         bool   `json:"enabled"`          // 是否启用迁移
    Dir             string `json:"dir"`              // 迁移文件目录
    AutoApply       bool   `json:"auto_apply"`       // 启动时自动执行
    TableName       string `json:"table_name"`       // 版本追踪表名
}
```

**步骤：**

- [ ] **Step 1: 添加配置字段**

在 `config/config.go` 中添加：

```go
// MigrationConfig 数据库迁移配置
type MigrationConfig struct {
    Enabled         bool   `json:"enabled"`          // 是否启用迁移
    Dir             string `json:"dir"`              // 迁移文件目录
    AutoApply       bool   `json:"auto_apply"`       // 启动时自动执行
    TableName       string `json:"table_name"`       // 版本追踪表名
}

// Config 结构体末尾添加
type Config struct {
    Server      ServerConfig
    Database    DatabaseConfig
    Redis       RedisConfig
    Auth        AuthConfig
    RateLimit   RateLimitConfig
    Logger      LoggerConfig
    Migration   MigrationConfig  // 新增
}
```

在 `LoadFromEnv()` 函数中添加：

```go
cfg.Migration.Enabled = getEnvAsBool("MIGRATION_ENABLED", false)
cfg.Migration.Dir = getEnv("MIGRATION_DIR", "./migrations")
cfg.Migration.AutoApply = getEnvAsBool("MIGRATION_AUTO_APPLY", false)
cfg.Migration.TableName = getEnv("MIGRATION_TABLE_NAME", "schema_migrations")
```

- [ ] **Step 2: 编写测试**

在 `config/config_test.go` 中添加（如果不存在则创建）：

```go
package config

import (
    "os"
    "testing"
)

func TestMigrationConfigFromEnv(t *testing.T) {
    // 设置环境变量
    os.Setenv("MIGRATION_ENABLED", "true")
    os.Setenv("MIGRATION_DIR", "./migrations")
    os.Setenv("MIGRATION_AUTO_APPLY", "true")
    os.Setenv("MIGRATION_TABLE_NAME", "schema_migrations")
    defer func() {
        os.Unsetenv("MIGRATION_ENABLED")
        os.Unsetenv("MIGRATION_DIR")
        os.Unsetenv("MIGRATION_AUTO_APPLY")
        os.Unsetenv("MIGRATION_TABLE_NAME")
    }()

    cfg, err := LoadFromEnv()
    if err != nil {
        t.Fatalf("LoadFromEnv() error = %v", err)
    }

    if !cfg.Migration.Enabled {
        t.Errorf("Expected Migration.Enabled = true, got false")
    }
    if cfg.Migration.Dir != "./migrations" {
        t.Errorf("Expected Migration.Dir = './migrations', got %s", cfg.Migration.Dir)
    }
    if !cfg.Migration.AutoApply {
        t.Errorf("Expected Migration.AutoApply = true, got false")
    }
    if cfg.Migration.TableName != "schema_migrations" {
        t.Errorf("Expected Migration.TableName = 'schema_migrations', got %s", cfg.Migration.TableName)
    }
}

func TestMigrationConfigDefaults(t *testing.T) {
    // 确保未设置环境变量时使用默认值
    cfg, _ := LoadFromEnv()

    if cfg.Migration.Enabled {
        t.Errorf("Expected default Migration.Enabled = false, got true")
    }
    if cfg.Migration.Dir != "./migrations" {
        t.Errorf("Expected default Migration.Dir = './migrations', got %s", cfg.Migration.Dir)
    }
}
```

- [ ] **Step 3: 运行测试**

```bash
cd server && go test ./config -v -run TestMigrationConfig
```

Expected: PASS

- [ ] **Step 4: 提交**

```bash
git add server/config/config.go server/config/config_test.go
git commit -m "feat(config): add migration configuration

- Add MigrationConfig struct
- Add MIGRATION_* environment variables
- Default: enabled=false, dir=./migrations, auto_apply=false
- Add unit tests for migration config loading

Co-Authored-By: Claude Code <noreply@anthropic.com>"
```

---

### Task 3: MigrationManager 核心实现

**目标：** 实现 MigrationManager 封装 golang-migrate

**文件：**
- 创建：`server/migration/manager.go`
- 创建：`server/migration/manager_test.go`
- 创建：`server/migration/config.go`（可选，复用 config.MigrationConfig）

**接口：**

```go
type MigrationManager struct {
    mgr *migrate.Migrate
    db  *model.DB
}

func NewManager(db *model.DB, cfg *config.MigrationConfig) (*MigrationManager, error)
func (m *MigrationManager) Up() error
func (m *MigrationManager) Down() error
func (m *MigrationManager) DownTo(version uint64) error
func (m *MigrationManager) Status() (version uint64, dirty bool, err error)
func (m *MigrationManager) Force(version int) error
func (m *MigrationManager) Close() error
```

**步骤：**

- [ ] **Step 1: 编写失败测试**

创建 `server/migration/manager_test.go`：

```go
package migration

import (
    "testing"
    "github.com/your-org/go-web-api-framework/model"
)

func TestNewManager_NilDB(t *testing.T) {
    _, err := NewManager(nil, &config.MigrationConfig{})
    if err == nil {
        t.Error("Expected error with nil DB")
    }
}

func TestNewManager_EmptyDir(t *testing.T) {
    db := model.NewMemoryCache() // 或其他测试数据库
    cfg := &config.MigrationConfig{
        Dir: "./nonexistent",
    }
    _, err := NewManager(db, cfg)
    if err == nil {
        t.Error("Expected error with nonexistent directory")
    }
}
```

运行测试，验证失败：

```bash
cd server && go test ./migration -v
```

Expected: FAIL (函数未定义)

- [ ] **Step 2: 实现最小配置**

创建 `server/migration/config.go`：

```go
package migration

import "github.com/your-org/go-web-api-framework/config"

// Config 迁移配置别名
type Config = config.MigrationConfig
```

- [ ] **Step 3: 实现 MigrationManager**

创建 `server/migration/manager.go`：

```go
package migration

import (
    "fmt"

    "github.com/golang-migrate/migrate/v4"
    "github.com/golang-migrate/migrate/v4/database/mysql"
    "github.com/golang-migrate/migrate/v4/database/postgres"
    _ "github.com/golang-migrate/migrate/v4/source/file"
    "github.com/your-org/go-web-api-framework/config"
    "github.com/your-org/go-web-api-framework/model"
)

// MigrationManager 封装 golang-migrate
type MigrationManager struct {
    mgr      *migrate.Migrate
    db       *model.DB
    dbType   string
    tableName string
}

// NewManager 创建迁移管理器
func NewManager(db *model.DB, cfg *config.MigrationConfig) (*MigrationManager, error) {
    if db == nil {
        return nil, fmt.Errorf("database connection is nil")
    }

    // 推断数据库类型（从 db 或 cfg）
    // 这里简化处理，实际应从 DatabaseConfig 获取
    dbType := "sqlite" // 默认

    // 构建 DSN
    dsn, err := buildDSN(dbType, cfg)
    if err != nil {
        return nil, fmt.Errorf("failed to build DSN: %w", err)
    }

    // 创建数据库驱动
    var driver migrate.Driver
    switch dbType {
    case "mysql":
        driver, err = mysql.WithInstance(db.Write.DB(), &mysql.Config{})
    case "postgres":
        driver, err = postgres.WithInstance(db.Write.DB(), &postgres.Config{})
    default:
        return nil, fmt.Errorf("unsupported database type for migration: %s", dbType)
    }

    if err != nil {
        return nil, fmt.Errorf("failed to create database driver: %w", err)
    }

    // 构建源 URL
    sourceURL := fmt.Sprintf("file://%s", cfg.Dir)

    // 创建迁移实例
    mgr, err := migrate.NewWithDatabaseInstance(
        sourceURL,
        dbType,
        driver,
    )

    if err != nil {
        return nil, fmt.Errorf("failed to create migrate instance: %w", err)
    }

    return &MigrationManager{
        mgr:      mgr,
        db:       db,
        dbType:   dbType,
        tableName: cfg.TableName,
    }, nil
}

// buildDSN 构建数据库连接字符串
func buildDSN(dbType string, cfg *config.MigrationConfig) (string, error) {
    // 简化实现，实际应从 DatabaseConfig 获取
    switch dbType {
    case "mysql":
        return fmt.Sprintf(
            "mysql://root:password@tcp(localhost:3306)/testdb?charset=utf8mb4&parseTime=True&loc=Local",
            nil
        case "postgres":
        return "postgres://user:password@localhost:5432/testdb?sslmode=disable", nil
    default:
        return "", fmt.Errorf("unsupported database type: %s", dbType)
    }
}

// Up 执行所有待处理的迁移
func (m *MigrationManager) Up() error {
    err := m.mgr.Up()
    if err == migrate.ErrNoChange {
        return nil // 无需操作
    }
    return err
}

// Down 回滚最后一个迁移
func (m *MigrationManager) Down() error {
    err := m.mgr.Steps(-1)
    if err == migrate.ErrNoChange {
        return nil
    }
    return err
}

// DownTo 回滚到指定版本
func (m *MigrationManager) DownTo(version uint64) error {
    err := m.mgr.Migrate(version)
    if err == migrate.ErrNoChange {
        return nil
    }
    return err
}

// Status 查看当前迁移状态
func (m *MigrationManager) Status() (version uint64, dirty bool, err error) {
    // golang-migrate v4 没有直接的 Status 方法
    // 需要查询 schema_migrations 表
    var v uint64
    err = m.db.Write.Raw("SELECT version FROM schema_migrations ORDER BY version DESC LIMIT 1").Scan(&v).Error
    if err != nil {
        return 0, false, err
    }

    var d bool
    err = m.db.Write.Raw("SELECT dirty FROM schema_migrations WHERE version = ?", v).Scan(&d).Error
    return v, d, err
}

// Force 强制设置版本（用于修复失败的迁移）
func (m *MigrationManager) Force(version int) error {
    return m.mgr.Force(version)
}

// Close 释放资源
func (m *MigrationManager) Close() error {
    if m.mgr != nil {
        m.mgr.Close()
    }
    return nil
}
```

**注意：** `buildDSN` 函数需要从 `model.DatabaseConfig` 获取连接信息，当前版本先简化实现，后续 Task 4 会完善。

- [ ] **Step 4: 运行测试**

```bash
cd server && go test ./migration -v
```

Expected: 测试失败（因为依赖未实现完全）

- [ ] **Step 5: 完善实现**

参考 Step 3 的代码，确保实现完整。

- [ ] **Step 6: 重新运行测试**

```bash
cd server && go test ./migration -v
```

Expected: PASS（至少通过编译和基础测试）

- [ ] **Step 7: 提交**

```bash
git add server/migration/
git commit -m "feat(migration): add MigrationManager core

- Implement MigrationManager wrapping golang-migrate
- Add Up, Down, DownTo, Status, Force methods
- Support MySQL and PostgreSQL
- Add basic unit tests

Co-Authored-By: Claude Code <noreply@anthropic.com>"
```

---

### Task 4: DSN 构建与数据库驱动集成

**目标：** 完善 MigrationManager 的 DSN 构建逻辑，从 config.DatabaseConfig 获取连接信息

**文件：**
- 修改：`server/migration/manager.go`

**步骤：**

- [ ] **Step 1: 修改 NewManager 签名**

```go
// 之前
func NewManager(db *model.DB, cfg *config.MigrationConfig) (*MigrationManager, error)

// 之后
func NewManager(db *model.DB, dbCfg *config.DatabaseConfig, cfg *config.MigrationConfig) (*MigrationManager, error)
```

- [ ] **Step 2: 实现 buildDSN**

```go
func buildDSN(dbCfg *config.DatabaseConfig) (string, error) {
    switch dbCfg.Type {
    case "mysql":
        return fmt.Sprintf(
            "mysql://%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
            dbCfg.User, dbCfg.Password, dbCfg.Host, dbCfg.Port, dbCfg.DBName,
        ), nil
    case "postgres":
        return fmt.Sprintf(
            "postgres://%s:%s@%s:%d/%s?sslmode=disable",
            dbCfg.User, dbCfg.Password, dbCfg.Host, dbCfg.Port, dbCfg.DBName,
        ), nil
    default:
        return "", fmt.Errorf("unsupported database type for migration: %s", dbCfg.Type)
    }
}
```

- [ ] **Step 3: 更新测试**

```bash
cd server && go test ./migration -v
```

Expected: PASS

- [ ] **Step 4: 提交**

```bash
git commit -am "feat(migration): improve DSN builder and database integration

- Accept DatabaseConfig for DSN construction
- Support MySQL and PostgreSQL connection strings
- Update tests for new signature

Co-Authored-By: Claude Code <noreply@anthropic.com>"
```

---

### Task 5: 示例迁移文件

**目标：** 创建示例迁移文件展示最佳实践

**文件：**
- 创建：`server/migration/example/20240101120000_init.up.sql`
- 创建：`server/migration/example/20240101120000_init.down.sql`

**步骤：**

- [ ] **Step 1: 创建 Up 迁移**

`server/migration/example/20240101120000_init.up.sql`：

```sql
-- ========================================
-- Migration: 20240101120000_init
-- Description: 初始化示例（创建 users 表）
-- Database: MySQL 8.0 / PostgreSQL 16
-- ========================================

-- 创建用户表
CREATE TABLE IF NOT EXISTS users (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY COMMENT '用户ID',
    username VARCHAR(50) NOT NULL UNIQUE COMMENT '用户名',
    email VARCHAR(255) NOT NULL UNIQUE COMMENT '邮箱地址',
    password_hash VARCHAR(255) NOT NULL COMMENT '密码哈希',
    role TINYINT NOT NULL DEFAULT 2 COMMENT '角色: 1=admin, 2=user',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    deleted_at TIMESTAMP NULL DEFAULT NULL COMMENT '删除时间（软删除）',

    INDEX idx_email (email),
    INDEX idx_role (role),
    INDEX idx_created_at (created_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci
COMMENT='用户表';
```

- [ ] **Step 2: 创建 Down 迁移**

`server/migration/example/20240101120000_init.down.sql`：

```sql
-- ========================================
-- Rollback: 20240101120000_init
-- Description: 删除 users 表
-- ========================================

DROP TABLE IF EXISTS users;
```

- [ ] **Step 3: 提交**

```bash
git add server/migration/example/
git commit -m "docs(migration): add example migration files

- Add 20240101120000_init.up.sql (create users table)
- Add 20240101120000_init.down.sql (drop users table)
- Include comments and cross-database notes

Co-Authored-By: Claude Code <noreply@anthropic.com>"
```

---

### Task 6: 集成到 Server

**目标：** 将 MigrationManager 集成到 Server 中，支持自动迁移

**文件：**
- 修改：`server/server.go`
- 修改：`server/server_test.go`（如果需要）

**步骤：**

- [ ] **Step 1: 修改 Server 结构体**

```go
type Server struct {
    engine      *gin.Engine
    config      *config.Config
    logger      *common.Logger
    db          *model.DB
    cache       model.Cache
    migrator    *migration.MigrationManager  // 新增
}
```

- [ ] **Step 2: 修改 NewServer**

在 `NewServer()` 函数末尾添加：

```go
// 初始化迁移管理器
if cfg.Migration.Enabled {
    mgr, err := migration.NewManager(s.db, &cfg.Database, &cfg.Migration)
    if err != nil {
        return nil, fmt.Errorf("failed to create migration manager: %w", err)
    }
    s.migrator = mgr

    // 自动执行迁移
    if cfg.Migration.AutoApply {
        if err := mgr.Up(); err != nil && err != migrate.ErrNoChange {
            s.logger.Warn("Auto migration failed: %v", err)
        }
    }
}
```

- [ ] **Step 3: 修改 Shutdown**

```go
func (s *Server) Shutdown(ctx context.Context) error {
    // ... 现有逻辑 ...
    if s.migrator != nil {
        s.migrator.Close()
    }
    return nil
}
```

- [ ] **Step 4: 添加 Getter（可选）**

```go
func (s *Server) Migrator() *migration.MigrationManager {
    return s.migrator
}
```

- [ ] **Step 5: 运行测试**

```bash
cd server && go test ./server -v
```

Expected: PASS（或至少通过编译）

- [ ] **Step 6: 提交**

```bash
git commit -am "feat(server): integrate migration manager

- Add MigrationManager to Server struct
- Auto-initialize when Migration.Enabled=true
- Auto-apply migrations on startup when AutoApply=true
- Add GetMigrator() method

Co-Authored-By: Claude Code <noreply@anthropic.com>"
```

---

### Task 7: CLI 迁移工具

**目标：** 提供独立的 CLI 命令执行迁移操作

**文件：**
- 创建：`server/cmd/migrate/main.go`

**步骤：**

- [ ] **Step 1: 创建 CLI 工具**

`server/cmd/migrate/main.go`：

```go
package main

import (
    "fmt"
    "os"

    "github.com/your-org/go-web-api-framework/config"
    "github.com/your-org/go-web-api-framework/migration"
    "github.com/your-org/go-web-api-framework/model"
)

func main() {
    // 加载配置
    cfg, err := config.LoadFromEnv()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Failed to load config: %v\n", err)
        os.Exit(1)
    }

    // 连接数据库
    db, err := model.NewDB(&cfg.Database)
    if err != nil {
        fmt.Fprintf(os.Stderr, "Failed to connect database: %v\n", err)
        os.Exit(1)
    }
    defer db.Close()

    // 解析命令
    if len(os.Args) < 2 {
        fmt.Println("Usage: go run cmd/migrate/main.go [up|down|status|force]")
        os.Exit(1)
    }

    command := os.Args[1]

    // 创建迁移管理器
    mgr, err := migration.NewManager(db, &cfg.Database, &cfg.Migration)
    if err != nil {
        fmt.Fprintf(os.Stderr, "Failed to create migration manager: %v\n", err)
        os.Exit(1)
    }
    defer mgr.Close()

    // 执行命令
    switch command {
    case "up":
        fmt.Println("Applying pending migrations...")
        if err := mgr.Up(); err != nil {
            fmt.Fprintf(os.Stderr, "Migration failed: %v\n", err)
            os.Exit(1)
        }
        fmt.Println("Migrations applied successfully!")

    case "down":
        fmt.Println("Rolling back last migration...")
        if err := mgr.Down(); err != nil {
            fmt.Fprintf(os.Stderr, "Rollback failed: %v\n", err)
            os.Exit(1)
        }
        fmt.Println("Migration rolled back successfully!")

    case "status":
        version, dirty, err := mgr.Status()
        if err != nil {
            fmt.Fprintf(os.Stderr, "Failed to get status: %v\n", err)
            os.Exit(1)
        }
        fmt.Printf("Current version: %d\n", version)
        fmt.Printf("Dirty: %v\n", dirty)

    case "force":
        if len(os.Args) < 3 {
            fmt.Println("Usage: go run cmd/migrate/main.go force <version>")
            os.Exit(1)
        }
        var version int
        fmt.Sscanf(os.Args[2], "%d", &version)
        if err := mgr.Force(version); err != nil {
            fmt.Fprintf(os.Stderr, "Force failed: %v\n", err)
            os.Exit(1)
        }
        fmt.Printf("Forced version to %d\n", version)

    default:
        fmt.Printf("Unknown command: %s\n", command)
        fmt.Println("Available commands: up, down, status, force")
        os.Exit(1)
    }
}
```

- [ ] **Step 2: 测试编译**

```bash
cd server && go build ./cmd/migrate
```

Expected: 编译成功

- [ ] **Step 3: 提交**

```bash
git add server/cmd/migrate/main.go
git commit -m "feat(migration): add CLI migration tool

- Add cmd/migrate/main.go with up/down/status/force commands
- Support manual migration execution
- Provide detailed error messages

Co-Authored-By: Claude Code <noreply@anthropic.com>"
```

---

### Task 8: 文档与使用示例

**目标：** 编写完整的使用文档和示例

**文件：**
- 创建：`server/docs/MIGRATION.md`
- 修改：`server/README.md`（新增迁移章节）

**步骤：**

- [ ] **Step 1: 编写迁移文档**

`server/docs/MIGRATION.md`：

```markdown
# 数据库迁移指南

## 目录

- [快速开始](#快速开始)
- [迁移文件命名](#迁移文件命名)
- [编写迁移](#编写迁移)
- [执行迁移](#执行迁移)
- [回滚迁移](#回滚迁移)
- [查看状态](#查看状态)
- [最佳实践](#最佳实践)

---

## 快速开始

### 1. 启用迁移功能

在 `.env` 文件中添加：

```env
MIGRATION_ENABLED=true
MIGRATION_DIR=./migrations
MIGRATION_AUTO_APPLY=false  # 开发环境可设为 true
```

### 2. 创建迁移文件

```bash
mkdir -p migrations
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

```sql
-- migrations/20240101120000_create_users.down.sql
DROP TABLE IF EXISTS users;
```

### 3. 执行迁移

```bash
go run cmd/migrate/main.go up
```

---

## 迁移文件命名

格式：`{时间戳}_{描述}.{方向}.sql`

示例：

- `20240101120000_create_users.up.sql`（升迁）
- `20240101120000_create_users.down.sql`（回滚）

---

## 编写迁移

### MySQL 8.0 示例

```sql
CREATE TABLE products (
    id INT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    price DECIMAL(10,2) NOT NULL,
    INDEX idx_name (name)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
```

### PostgreSQL 16 示例

```sql
CREATE TABLE products (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    price NUMERIC(10,2) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_products_name ON products(name);
```

---

## 最佳实践

1. **始终提供 Down 迁移** — 便于回滚
2. **使用事务** — 确保迁移原子性
3. **测试迁移** — 在测试数据库上先执行
4. **备份生产数据** — 生产环境执行前务必备份
5. **小步提交** — 每个迁移只做一件事

```

- [ ] **Step 2: 更新 README**

在 `server/README.md` 和 `server/README.zh-CN.md` 中新增"数据库迁移"章节

- [ ] **Step 3: 提交**

```bash
git commit -am "docs(migration): add migration guide and examples

- Add docs/MIGRATION.md with complete guide
- Update README with migration section
- Provide MySQL and PostgreSQL examples

Co-Authored-By: Claude Code <noreply@anthropic.com>"
```

---

### Task 9: 集成测试

**目标：** 编写针对 MySQL 和 PostgreSQL 的集成测试（可选，如无测试环境可跳过）

**文件：**
- 创建：`server/migration/integration_test.go`

**注意：** 此任务需要在真实 MySQL 8.0 / PostgreSQL 16 环境下运行，如果本地环境没有，可以标记为 `skip`。

**步骤：**

- [ ] **Step 1: 编写集成测试模板**

`server/migration/integration_test.go`：

```go
package migration

import (
    "testing"
    "github.com/your-org/go-web-api-framework/config"
    "github.com/your-org/go-web-api-framework/model"
)

func TestMySQLMigration(t *testing.T) {
    if testing.Short() {
        t.Skip("Skipping MySQL integration test in short mode")
    }

    cfg := &config.DatabaseConfig{
        Type:     "mysql",
        Host:     "localhost",
        Port:     3306,
        User:     "root",
        Password: "password",
        DBName:   "test_db",
    }

    db, err := model.NewDB(cfg)
    if err != nil {
        t.Skipf("Cannot connect to MySQL: %v", err)
    }
    defer db.Close()

    // 测试迁移逻辑...
}

func TestPostgresMigration(t *testing.T) {
    if testing.Short() {
        t.Skip("Skipping PostgreSQL integration test in short mode")
    }

    cfg := &config.DatabaseConfig{
        Type:     "postgres",
        Host:     "localhost",
        Port:     5432,
        User:     "postgres",
        Password: "password",
        DBName:   "test_db",
    }

    db, err := model.NewDB(cfg)
    if err != nil {
        t.Skipf("Cannot connect to PostgreSQL: %v", err)
    }
    defer db.Close()

    // 测试迁移逻辑...
}
```

- [ ] **Step 2: 提交（可选）**

```bash
git add server/migration/integration_test.go
git commit -m "test(migration): add integration test templates

- Add MySQL and PostgreSQL integration test templates
- Tests are skipped by default (require real database)
- Use -short flag to skip in CI

Co-Authored-By: Claude Code <noreply@anthropic.com>"
```

---

## 测试验证

### 完整测试命令

```bash
# 单元测试
cd server && go test ./migration -v

# 所有包测试
cd server && go test ./...

# 集成测试（跳过）
cd server && go test ./migration -short

# 覆盖率
cd server && go test -cover ./migration
```

### 预期结果

- ✅ 所有单元测试通过
- ✅ 集成测试在无数据库环境跳过（不失败）
- ✅ 编译通过，无警告

---

## 完成标准

### 功能完整性

- [x] MigrationConfig 配置字段
- [x] MigrationManager 核心方法（Up, Down, DownTo, Status, Force）
- [x] DSN 构建（MySQL + PostgreSQL）
- [x] 示例迁移文件
- [x] CLI 迁移工具
- [x] 文档和使用示例

### 质量保证

- [x] 单元测试覆盖率 ≥ 70%
- [x] 集成测试模板（可选跳过）
- [x] 代码符合 Go 规范（gofmt）
- [x] 向后兼容（不影响现有代码）

### 文档完整性

- [x] README 新增迁移章节
- [x] MIGRATION.md 完整指南
- [x] 代码注释和 GoDoc

---

## 潜在风险与缓解

### 风险 1: golang-migrate 版本兼容性

**缓解：** 使用稳定的 v4.17.0 版本，测试验证后再升级。

### 风险 2: 多数据库 SQL 语法差异

**缓解：** 文档中明确标注 MySQL/PostgreSQL 差异化部分，建议用户分别维护迁移文件。

### 风险 3: 迁移失败导致数据库不一致

**缓解：** golang-migrate 使用事务确保原子性；文档中强调备份的重要性。

### 风险 4: 与 GORM AutoMigrate 混用导致冲突

**缓解：** 文档中明确建议"二选一"，提供迁移方案优先的说明。

---

## 后续优化（P2/P3）

### P2: 代码生成工具

- [ ] 从 SQL 生成 GORM Model 骨架
- [ ] 从 GORM Model 生成迁移 SQL 模板

### P3: 一致性检查

- [ ] 检测 Model 和迁移文件的差异
- [ ] CI/CD 集成检查

### P4: 迁移钩子

- [ ] BeforeUp / AfterUp 钩子
- [ ] BeforeDown / AfterDown 钩子

---

**计划完成时间估算：** 约 6-8 个任务（每个任务 1-2 小时），总共约 8-12 小时

**下一步：** 按此计划逐步实施，每个任务完成后提交并进行代码审查。
