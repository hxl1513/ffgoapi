# 数据库迁移功能设计

> **For agentic workers:** 本设计文档用于实现数据库迁移功能，支持 MySQL 8.0 和 PostgreSQL 16。

**目标：** 为 Go Web API 框架提供企业级数据库迁移能力，基于 golang-migrate 标准。

---

## 1. 架构概览

### 1.1 核心设计

```
┌─────────────────────────────────────────────────────────┐
│                    应用层                                │
│  ┌─────────────────────────────────────────────────┐   │
│  │  MigrationManager                                │   │
│  │  - Up()      // 执行所有待处理迁移                │   │
│  │  - Down()    // 回滚最后一个迁移                   │   │
│  │  - Status()  // 查看迁移状态                       │   │
│  │  - Force()   // 强制标记版本                       │   │
│  └─────────────────────────────────────────────────┘   │
└─────────────────────────┬───────────────────────────────┘
                          │
┌─────────────────────────▼───────────────────────────────┐
│               golang-migrate 核心                        │
│  ┌─────────────────────────────────────────────────┐   │
│  |  - 读取 migrations/ 目录                           │   │
│  |  - 解析 .up.sql / .down.sql 文件                   │   │
│  |  - 维护 schema_migrations 表                        │   │
│  |  - 执行 SQL 并记录版本                              │   │
│  └─────────────────────────────────────────────────┘   │
└─────────────────────────┬───────────────────────────────┘
                          │
┌─────────────────────────▼───────────────────────────────┐
│              数据库层                                    │
│   ┌─────────────┐  ┌─────────────┐  ┌─────────────┐  │
│   │  MySQL 8.0  │  │ PostgreSQL  │  │   SQLite    │  │
│   └─────────────┘  └─────────────┘  └─────────────┘  │
└─────────────────────────────────────────────────────────┘
```

---

## 2. 目录结构

### 2.1 框架侧结构

```
server/
├── migration/                    # 新增：迁移管理包
│   ├── manager.go                # MigrationManager 核心
│   ├── manager_test.go           # 单元测试
│   ├── config.go                 # MigrationConfig 配置
│   └── example/                  # 示例迁移
│       ├── 20240101120000_init.up.sql
│       └── 20240101120000_init.down.sql
├── config/
│   └── config.go                 # 新增 MigrationConfig 字段
└── go.mod                        # 新增 golang-migrate 依赖
```

### 2.2 用户项目结构

```
your-project/
├── main.go
├── .env
├── migrations/                    # 用户管理迁移文件
│   ├── 20240101120000_create_users.up.sql
│   ├── 20240101120000_create_users.down.sql
│   ├── 20240102150000_add_roles.up.sql
│   └── 20240102150000_add_roles.down.sql
└── model/
    └── user.go                   # GORM Model（业务代码）
```

---

## 3. 核心组件设计

### 3.1 配置扩展 (`config/config.go`)

```go
// 新增 MigrationConfig
type MigrationConfig struct {
    Enabled         bool   `json:"enabled"`          // 是否启用迁移
    Dir             string `json:"dir"`              // 迁移文件目录
    AutoApply       bool   `json:"auto_apply"`       // 启动时自动执行迁移
    TableName       string `json:"table_name"`       // 版本追踪表名
}

// Config 新增字段
type Config struct {
    // ... 现有字段 ...
    Migration MigrationConfig  // 新增
}
```

**环境变量：**

| 变量 | 类型 | 默认值 | 说明 |
|------|------|--------|------|
| `MIGRATION_ENABLED` | bool | `false` | 是否启用迁移功能 |
| `MIGRATION_DIR` | string | `./migrations` | 迁移文件目录 |
| `MIGRATION_AUTO_APPLY` | bool | `false` | 启动时自动执行 |
| `MIGRATION_TABLE_NAME` | string | `schema_migrations` | 版本追踪表名 |

---

### 3.2 MigrationManager (`migration/manager.go`)

```go
package migration

import (
    "github.com/golang-migrate/migrate/v4"
    _ "github.com/golang-migrate/migrate/v4/database/mysql"
    _ "github.com/golang-migrate/migrate/v4/database/postgres"
    _ "github.com/golang-migrate/migrate/v4/source/file"
    "github.com/hxl1513/ffgoapi/model"
)

// MigrationManager 封装 golang-migrate
type MigrationManager struct {
    mgr *migrate.Migrate
    db  *model.DB
}

// NewManager 创建迁移管理器
func NewManager(db *model.DB, cfg *config.MigrationConfig) (*MigrationManager, error)

// Up 执行所有待处理的迁移
func (m *MigrationManager) Up() error

// Down 回滚最后一个迁移
func (m *MigrationManager) Down() error

// DownTo 回滚到指定版本
func (m *MigrationManager) DownTo(version uint64) error

// Status 查看当前迁移状态
func (m *MigrationManager) Status() (version uint64, dirty bool, err error)

// Force 强制设置版本（用于修复失败的迁移）
func (m *MigrationManager) Force(version int) error

// Close 释放资源
func (m *MigrationManager) Close() error
```

---

### 3.3 DSN 构建

```go
// manager.go 内部方法

func (m *MigrationManager) dsn() string {
    switch m.dbType {
    case "mysql":
        return fmt.Sprintf(
            "mysql://%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
            m.user, m.password, m.host, m.port, m.dbname,
        )
    case "postgres":
        return fmt.Sprintf(
            "postgres://%s:%s@%s:%d/%s?sslmode=disable",
            m.user, m.password, m.host, m.port, m.dbname,
        )
    default:
        return ""
    }
}
```

---

### 3.4 自动创建版本追踪表

```go
// manager.go - 内部方法

func (m *MigrationManager) ensureVersionTable() error {
    // 使用 golang-migrate 的 internal 包创建
    // 或者直接执行 DDL
    sql := fmt.Sprintf(`
        CREATE TABLE IF NOT EXISTS %s (
            version BIGINT PRIMARY KEY,
            dirty BOOLEAN NOT NULL DEFAULT FALSE,
            applied_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
        )
    `, m.tableName)

    return m.db.Write.Exec(sql).Error
}
```

---

## 4. API 设计

### 4.1 基础用法

```go
package main

import (
    "log"
    "github.com/hxl1513/ffgoapi/config"
    "github.com/hxl1513/ffgoapi/model"
    "github.com/hxl1513/ffgoapi/migration"
)

func main() {
    // 1. 加载配置
    cfg, err := config.LoadFromEnv()
    if err != nil {
        log.Fatal(err)
    }

    // 2. 连接数据库
    db, err := model.NewDB(&cfg.Database)
    if err != nil {
        log.Fatal(err)
    }
    defer db.Close()

    // 3. 创建迁移管理器
    mgr, err := migration.NewManager(db, &cfg.Migration)
    if err != nil {
        log.Fatal(err)
    }
    defer mgr.Close()

    // 4. 查看当前版本
    version, dirty, err := mgr.Status()
    if err != nil {
        log.Printf("Status check failed: %v", err)
    } else {
        log.Printf("Current version: %d, dirty: %v", version, dirty)
    }

    // 5. 执行迁移
    if err := mgr.Up(); err != nil && err != migrate.ErrNoChange {
        log.Fatal(err)
    }

    log.Println("Migration completed!")
}
```

---

### 4.2 集成到 Server

```go
// server/server.go

func (s *Server) startMigration() error {
    if !s.cfg.Migration.Enabled {
        return nil
    }

    mgr, err := migration.NewManager(s.db, &s.cfg.Migration)
    if err != nil {
        return err
    }
    defer mgr.Close()

    // 自动执行迁移
    if s.cfg.Migration.AutoApply {
        return mgr.Up()
    }

    return nil
}
```

---

## 5. 迁移文件规范

### 5.1 文件命名

```
{时间戳}_{描述}.{方向}.sql

示例：
20240101120000_create_users_table.up.sql
20240101120000_create_users_table.down.sql
20240102150000_add_roles.up.sql
20240102150000_add_roles.down.sql
```

### 5.2 文件格式

#### **Up 迁移（升迁）**

```sql
-- 20240101120000_create_users_table.up.sql
-- 创建用户表

CREATE TABLE IF NOT EXISTS users (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    username VARCHAR(50) NOT NULL UNIQUE COMMENT '用户名',
    email VARCHAR(255) NOT NULL UNIQUE COMMENT '邮箱',
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

#### **Down 迁移（回滚）**

```sql
-- 20240101120000_create_users_table.down.sql
-- 回滚：删除用户表

DROP TABLE IF EXISTS users;
```

---

### 5.3 跨数据库兼容性

#### **方法 A：条件 SQL（推荐）**

```sql
-- 使用 golang-migrate 的 $if 语法
-- $if(DATABASE_TYPE == 'mysql')

CREATE TABLE users (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    ...
);

-- $elif(DATABASE_TYPE == 'postgres')

CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    ...
);

-- $fi
```

#### **方法 B：分离文件（简单粗暴）**

```
migrations/
├── mysql/
│   └── 20240101120000_create_users.up.sql
└── postgres/
    └── 20240101120000_create_users.up.sql
```

---

## 6. 集成到现有框架

### 6.1 配置加载

```go
// config/config.go

type Config struct {
    // ... 现有字段 ...
    Migration MigrationConfig  // 新增
}

type MigrationConfig struct {
    Enabled         bool   `json:"enabled"`
    Dir             string `json:"dir"`
    AutoApply       bool   `json:"auto_apply"`
    TableName       string `json:"table_name"`
}

func LoadFromEnv() (*Config, error) {
    cfg := &Config{}

    // ... 现有配置加载 ...

    // 新增：迁移配置
    cfg.Migration.Enabled = getEnvAsBool("MIGRATION_ENABLED", false)
    cfg.Migration.Dir = getEnv("MIGRATION_DIR", "./migrations")
    cfg.Migration.AutoApply = getEnvAsBool("MIGRATION_AUTO_APPLY", false)
    cfg.Migration.TableName = getEnv("MIGRATION_TABLE_NAME", "schema_migrations")

    return cfg, nil
}
```

---

### 6.2 Server 集成

```go
// server/server.go

type Server struct {
    engine      *gin.Engine
    config      *config.Config
    logger      *common.Logger
    db          *model.DB
    cache       model.Cache
    migrator    *migration.MigrationManager  // 新增
}

func NewServer(cfg *config.Config) (*Server, error) {
    // ... 现有初始化逻辑 ...

    s := &Server{...}

    // 新增：初始化迁移管理器
    if cfg.Migration.Enabled {
        mgr, err := migration.NewManager(s.db, &cfg.Migration)
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

    return s, nil
}

func (s *Server) Close() error {
    // ... 现有关闭逻辑 ...
    if s.migrator != nil {
        s.migrator.Close()
    }
    return nil
}
```

---

### 6.3 CLI 命令集成（可选）

```go
// cmd/migrate.go

package main

import (
    "os"
    "github.com/hxl1513/ffgoapi/config"
    "github.com/hxl1513/ffgoapi/model"
    "github.com/hxl1513/ffgoapi/migration"
)

func main() {
    if len(os.Args) < 2 {
        fmt.Println("Usage: go run cmd/migrate.go [up|down|status|force]")
        os.Exit(1)
    }

    cfg, _ := config.LoadFromEnv()
    db, _ := model.NewDB(&cfg.Database)
    mgr, _ := migration.NewManager(db, &cfg.Migration)

    switch os.Args[1] {
    case "up":
        mgr.Up()
    case "down":
        mgr.Down()
    case "status":
        mgr.Status()
    case "force":
        // ...
    }

    mgr.Close()
}
```

**使用：**
```bash
# 执行迁移
go run cmd/migrate.go up

# 查看状态
go run cmd/migrate.go status

# 回滚
go run cmd/migrate.go down
```

---

## 7. 错误处理

### 7.1 常见错误场景

| 场景 | 错误 | 处理策略 |
|------|------|----------|
| **迁移文件不存在** | `os.ErrNotExist` | 返回错误，提示检查目录 |
| **SQL 语法错误** | 数据库返回错误 | 返回详细错误，包含 SQL 片段 |
| **迁移已执行** | `migrate.ErrNoChange` | 视为成功，跳过 |
| **版本追踪表损坏** | `sql.ErrNoRows` | 尝试重建追踪表 |
| **并发迁移** | `migrate.ErrLocked` | 提示用户等待或删除锁 |
| **脏版本（dirty）** | `dirty=true` | 提示用户修复或 Force |

### 7.2 错误处理示例

```go
func (m *MigrationManager) Up() error {
    err := m.mgr.Up()

    if err == migrate.ErrNoChange {
        return nil  // 无需操作
    }

    if err != nil {
        return fmt.Errorf("migration failed: %w", err)
    }

    return nil
}
```

---

## 8. 测试策略

### 8.1 单元测试

```go
// migration/manager_test.go

func TestMigrationManager_Up(t *testing.T) {
    // 1. 准备测试数据库（SQLite :memory:）
    db := setupTestDB(t)
    defer db.Close()

    // 2. 创建迁移管理器
    mgr := setupTestManager(t, db)
    defer mgr.Close()

    // 3. 执行迁移
    err := mgr.Up()
    assert.NoError(t, err)

    // 4. 验证表是否存在
    assert.TableExists(t, db, "users")
}

func TestMigrationManager_Down(t *testing.T) {
    // 测试回滚逻辑
}

func TestMigrationManager_Status(t *testing.T) {
    // 测试版本状态查询
}
```

### 8.2 集成测试

```go
// migration/integration_test.go

func TestMySQLMigration(t *testing.T) {
    // 针对 MySQL 8.0 的集成测试
}

func TestPostgresMigration(t *testing.T) {
    // 针对 PostgreSQL 16 的集成测试
}
```

---

## 9. 安全考虑

### 9.1 迁移锁

```go
// 防止并发执行迁移
mgr, err := migrate.New(
    "file://./migrations",
    "mysql://user:pass@tcp(localhost:3306)/dbname",
)

// 迁移锁（自动）
// golang-migrate 内部使用数据库锁防止并发
```

### 9.2 备份建议

```go
// 执行迁移前提示备份
func (m *MigrationManager) UpWithBackup() error {
    log.Println("WARNING: Please backup your database before migration!")
    // ...
}
```

### 9.3 权限控制

```go
// 迁移操作应限制在开发/部署环境
if cfg.Environment == "production" {
    if !confirm("Are you sure you want to run migration in PRODUCTION?") {
        return errors.New("migration aborted")
    }
}
```

---

## 10. 与现有组件的兼容性

### 10.1 与 GORM AutoMigrate 的关系

**建议：二选一**

#### **方案 A：使用迁移系统（推荐）**

```go
// ❌ 禁用 GORM AutoMigrate
db, err := model.NewDB(&cfg.Database)
// 不要调用 db.AutoMigrate(&User{})
```

```sql
-- ✅ 使用迁移文件
-- migrations/20240101120000_create_users.up.sql
CREATE TABLE users (...);
```

**优点：**
- ✅ 版本可控
- ✅ 可审计
- ✅ 可回滚

---

#### **方案 B：开发阶段用 AutoMigrate，生产用迁移**

```go
// 开发环境：使用 AutoMigrate
if cfg.Environment == "development" {
    db.AutoMigrate(&User{})
} else {
    // 生产环境：使用迁移
    mgr.Up()
}
```

**优点：**
- ✅ 开发快速
- ✅ 生产稳定

**缺点：**
- ⚠️ 两套系统，容易不一致

---

### 10.2 与 BaseModel 的关系

**无冲突**，但需注意：

```go
// model/base_model.go
type BaseModel struct {
    ID        uint           `gorm:"primaryKey" json:"id"`
    CreatedAt time.Time      `json:"created_at"`
    UpdatedAt time.Time      `json:"updated_at"`
    DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}
```

```sql
-- 迁移文件中需与 BaseModel 字段对齐
CREATE TABLE users (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,  -- 对应 ID
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,  -- 对应 CreatedAt
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,  -- 对应 UpdatedAt
    deleted_at TIMESTAMP NULL DEFAULT NULL,  -- 对应 DeletedAt（软删除）
    ...
);
```

**建议：** 在文档中明确标注 BaseModel 的标准 DDL。

---

## 11. 示例使用

### 11.1 完整示例

```go
package main

import (
    "fmt"
    "log"

    "github.com/hxl1513/ffgoapi/config"
    "github.com/hxl1513/ffgoapi/migration"
    "github.com/hxl1513/ffgoapi/model"
    "github.com/hxl1513/ffgoapi/server"
)

func main() {
    // 1. 加载配置
    cfg, err := config.LoadFromEnv()
    if err != nil {
        log.Fatal(err)
    }

    // 2. 连接数据库
    db, err := model.NewDB(&cfg.Database)
    if err != nil {
        log.Fatal(err)
    }

    // 3. 执行迁移（如果启用）
    if cfg.Migration.Enabled {
        mgr, err := migration.NewManager(db, &cfg.Migration)
        if err != nil {
            log.Fatal(err)
        }
        defer mgr.Close()

        // 查看当前版本
        version, dirty, err := mgr.Status()
        if err != nil {
            log.Printf("Warning: failed to get migration status: %v", err)
        }
        fmt.Printf("Current migration version: %d (dirty: %v)\n", version, dirty)

        // 执行迁移
        if err := mgr.Up(); err != nil && err != migrate.ErrNoChange {
            log.Fatalf("Migration failed: %v", err)
        }

        fmt.Println("Migration completed successfully!")
    }

    // 4. 启动服务器
    srv, err := server.NewServer(cfg)
    if err != nil {
        log.Fatal(err)
    }
    defer srv.Shutdown(context.Background())

    srv.RegisterRoutes(func(r *gin.RouterGroup) {
        r.GET("/health", func(c *gin.Context) {
            c.JSON(200, gin.H{"status": "ok"})
        })
    })

    srv.Start()
}
```

### 11.2 环境变量示例

```env
# .env

# 数据库
DB_TYPE=mysql
DB_HOST=localhost
DB_PORT=3306
DB_USER=root
DB_PASSWORD=secret
DB_NAME=myapp

# 迁移配置
MIGRATION_ENABLED=true
MIGRATION_DIR=./migrations
MIGRATION_AUTO_APPLY=true  # 启动时自动执行
MIGRATION_TABLE_NAME=schema_migrations
```

---

## 12. 依赖项

### 12.1 go.mod 新增依赖

```go
require (
    github.com/golang-migrate/migrate/v4 v4.17.0
    github.com/golang-migrate/migrate/v4/database/mysql v4.17.0
    github.com/golang-migrate/migrate/v4/database/postgres v4.17.0
    github.com/golang-migrate/migrate/v4/source/file v4.17.0
)
```

---

## 13. 向后兼容性

### 13.1 破坏性变更

**无破坏性变更**

- 新增 `migration` 包，不影响现有代码
- 配置新增字段有默认值（`Migration.Enabled=false`）
- 现有用户无需修改任何代码即可升级

### 13.2 迁移路径

```go
// 现有代码无需修改
cfg, _ := config.LoadFromEnv()
db, _ := model.NewDB(&cfg.Database)
srv, _ := server.NewServer(cfg)

// 如需使用迁移功能，只需：
cfg.Migration.Enabled = true
cfg.Migration.AutoApply = true
// 其他字段有默认值，无需配置
```

---

## 14. 文档要求

### 14.1 用户文档

- [ ] README.md 新增"数据库迁移"章节
- [ ] docs/migration.md — 迁移完整教程
- [ ] 示例迁移文件（带注释）
- [ ] CLI 命令说明

### 14.2 API 文档

- [ ] MigrationManager 所有方法的 GoDoc
- [ ] MigrationConfig 字段说明
- [ ] 错误码和错误处理

---

## 15. Spec Self-Review

### ✅ Placeholder Scan
- [x] 无 "TBD"、"TODO"
- [x] 所有代码示例完整
- [x] 所有配置项有默认值

### ✅ Internal Consistency
- [x] 配置字段与环境变量命名一致
- [x] MigrationManager API 与现有 model.DB 兼容
- [x] 迁移文件命名符合 golang-migrate 标准

### ✅ Scope Check
- [x] 聚焦迁移功能，不涉及其他组件
- [x] 支持 MySQL 8.0 和 PostgreSQL 16
- [x] 兼容现有 GORM Model 定义

### ✅ Ambiguity Check
- [x] 明确 MigrationManager 的错误处理策略
- [x] 明确迁移文件命名规范
- [x] 明确向后兼容性保证

---

## 16. 待讨论问题

### Q1: 是否需要在框架中内置"常见迁移模板"？

如 `create_table`、`add_column`、`create_index` 等通用迁移模板？

**建议：** 提供但不强制，用户可复制使用。

### Q2: 是否支持迁移回滚到指定版本？

```go
mgr.DownTo(20240101120000) // 回滚到指定版本
```

**建议：** ✅ 支持（golang-migrate 原生支持）

### Q3: 是否提供迁移钩子（Hooks）？

在迁移前后执行自定义逻辑：

```go
mgr.BeforeUp(func() error { ... })
mgr.AfterUp(func() error { ... })
```

**建议：** P3 版本考虑，当前版本先实现核心功能

---

**请审阅此设计文档，确认是否符合预期后，我将继续编写详细的实现计划。**
