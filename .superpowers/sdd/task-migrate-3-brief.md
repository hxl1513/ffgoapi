# Task 3 Brief: MigrationManager 核心实现

**目标：** 实现 MigrationManager 封装 golang-migrate，提供 Up/Down/Status/Force 等核心方法

**文件：**
- 创建：`server/migration/manager.go`
- 创建：`server/migration/manager_test.go`

---

## 上下文

这是数据库迁移功能的第三个任务。Task 1（依赖）和 Task 2（配置）已完成。

**前置任务：** Task Migrate-1 ✅, Task Migrate-2 ✅

**工作目录：** `C:\workspace\go_project\bhz`

**Go 工具链：** `/c/go/bin/go`

**当前 git HEAD：** 38d855a

---

## 你的任务

### Step 1: 创建 migration 包目录

```bash
mkdir -p C:/workspace/go_project/bhz/server/migration
```

### Step 2: 编写失败测试

创建 `server/migration/manager_test.go`，先写测试：

```go
package migration

import (
    "testing"
    "github.com/hxl1513/ffgoapi/model"
    "github.com/hxl1513/ffgoapi/config"
)

// TestNewManager_NilDB 测试 nil DB 参数
func TestNewManager_NilDB(t *testing.T) {
    _, err := NewManager(nil, &config.DatabaseConfig{}, &config.MigrationConfig{})
    if err == nil {
        t.Error("Expected error with nil DB")
    }
}

// TestMigrationManager_Up_NoMigrations 测试无迁移文件时 Up 的行为
func TestMigrationManager_Up_NoMigrations(t *testing.T) {
    // 这个测试需要真实的数据库连接，暂时跳过
    t.Skip("Integration test - requires real database")
}
```

运行测试，验证失败（因为 `NewManager` 未定义）：

```bash
cd C:/workspace/go_project/bhz/server && /c/go/bin/go test ./migration -v
```

Expected: FAIL (undefined: NewManager)

### Step 3: 实现 MigrationManager

创建 `server/migration/manager.go`：

```go
package migration

import (
    "fmt"

    "github.com/golang-migrate/migrate/v4"
    "github.com/golang-migrate/migrate/v4/database/mysql"
    "github.com/golang-migrate/migrate/v4/database/postgres"
    _ "github.com/golang-migrate/migrate/v4/source/file"
    "github.com/hxl1513/ffgoapi/config"
    "github.com/hxl1513/ffgoapi/model"
)

// MigrationManager 封装 golang-migrate
type MigrationManager struct {
    mgr       *migrate.Migrate
    db        *model.DB
    dbType    string
    tableName string
}

// NewManager 创建迁移管理器
func NewManager(db *model.DB, dbCfg *config.DatabaseConfig, cfg *config.MigrationConfig) (*MigrationManager, error) {
    if db == nil {
        return nil, fmt.Errorf("database connection is nil")
    }

    if cfg == nil {
        return nil, fmt.Errorf("migration config is nil")
    }

    // 推断数据库类型
    dbType := dbCfg.Type
    if dbType == "" {
        dbType = "sqlite"
    }

    // 构建 DSN
    dsn, err := buildDSN(dbCfg)
    if err != nil {
        return nil, fmt.Errorf("failed to build DSN: %w", err)
    }

    // 创建数据库驱动
    var driver migrate.Driver
    switch dbType {
    case "mysql":
        driver, err = mysql.WithInstance(db.Write.DB(), &mysql.Config{})
        if err != nil {
            return nil, fmt.Errorf("failed to create MySQL driver: %w", err)
        }
    case "postgres":
        driver, err = postgres.WithInstance(db.Write.DB(), &postgres.Config{})
        if err != nil {
            return nil, fmt.Errorf("failed to create PostgreSQL driver: %w", err)
        }
    default:
        return nil, fmt.Errorf("unsupported database type for migration: %s (only mysql and postgres are supported)", dbType)
    }

    // 构建源 URL（迁移文件目录）
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
        mgr:       mgr,
        db:        db,
        dbType:    dbType,
        tableName: cfg.TableName,
    }, nil
}

// buildDSN 构建数据库连接字符串（仅用于日志或诊断，实际连接已存在）
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
        return "", fmt.Errorf("unsupported database type: %s", dbCfg.Type)
    }
}

// Up 执行所有待处理的迁移
func (m *MigrationManager) Up() error {
    err := m.mgr.Up()
    if err == migrate.ErrNoChange {
        return nil // 无需操作
    }
    if err != nil {
        return fmt.Errorf("failed to apply migrations: %w", err)
    }
    return nil
}

// Down 回滚最后一个迁移
func (m *MigrationManager) Down() error {
    err := m.mgr.Steps(-1)
    if err == migrate.ErrNoChange {
        return nil
    }
    if err != nil {
        return fmt.Errorf("failed to rollback migration: %w", err)
    }
    return nil
}

// DownTo 回滚到指定版本
func (m *MigrationManager) DownTo(version uint64) error {
    err := m.mgr.Migrate(version)
    if err == migrate.ErrNoChange {
        return nil
    }
    if err != nil {
        return fmt.Errorf("failed to migrate to version %d: %w", version, err)
    }
    return nil
}

// Status 查看当前迁移状态
// 注意：golang-migrate v4 没有直接的 Status 方法
// 这里返回 ErrNotImplemented，后续版本可以查询 schema_migrations 表
func (m *MigrationManager) Status() (version uint64, dirty bool, err error) {
    // TODO: 实现查询 schema_migrations 表
    return 0, false, fmt.Errorf("Status() not implemented yet")
}

// Force 强制设置版本（用于修复失败的迁移）
func (m *MigrationManager) Force(version int) error {
    if err := m.mgr.Force(version); err != nil {
        return fmt.Errorf("failed to force version to %d: %w", version, err)
    }
    return nil
}

// Close 释放资源
func (m *MigrationManager) Close() error {
    if m.mgr != nil {
        m.mgr.Close()
    }
    return nil
}
```

### Step 4: 运行测试

```bash
cd C:/workspace/go_project/bhz/server && /c/go/bin/go test ./migration -v
```

**预期：**
- 测试通过编译
- TestNewManager_NilDB PASS
- TestMigrationManager_Up_NoMigrations SKIP（集成测试，跳过）

### Step 5: 完善测试

根据测试结果调整实现，确保所有测试通过。

### Step 6: 提交

```bash
cd C:/workspace/go_project/bhz && git add server/migration/
git commit -m "feat(migration): add MigrationManager core

- Implement MigrationManager wrapping golang-migrate
- Add NewManager with database driver support (MySQL, PostgreSQL)
- Implement Up, Down, DownTo, Status, Force, Close methods
- Add DSN builder for MySQL and PostgreSQL
- Add basic unit tests
- Return ErrNoChange when no migrations pending

Co-Authored-By: Claude Code <noreply@anthropic.com>"
```

---

## 验证

- [ ] `migration` 包创建成功
- [ ] `NewManager` 接受 `(db *model.DB, dbCfg *config.DatabaseConfig, cfg *config.MigrationConfig)`
- [ ] `Up()` 执行迁移，处理 `ErrNoChange`
- [ ] `Down()` 回滚最后一个迁移
- [ ] `DownTo(version)` 回滚到指定版本
- [ ] `Status()` 返回未实现错误（占位）
- [ ] `Force(version)` 强制设置版本
- [ ] `Close()` 释放资源
- [ ] 测试通过（至少通过编译和基础测试）
- [ ] 代码已提交

---

## 报告

**报告文件：** `C:\workspace\go_project\bhz\.superpowers\sdd\task-migrate-3-report.md`

**报告必须包含：**
- 状态：DONE 或 NEEDS_CONTEXT
- git commit SHA
- 测试结果（实际输出）
- 任何问题或关注点（特别是 Status() 方法未实现的问题）

**提交：** 完成任务后请自行提交到 git。
