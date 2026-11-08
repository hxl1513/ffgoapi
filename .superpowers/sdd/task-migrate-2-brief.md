# Task 2 Brief: MigrationConfig 配置扩展

**目标：** 在 config 包中添加 MigrationConfig 配置结构和环境变量支持

**文件：**
- 修改：`server/config/config.go`
- 创建/修改：`server/config/config_test.go`（添加迁移配置测试）

---

## 上下文

这是数据库迁移功能的第二个任务。Task 1 已完成，golang-migrate 依赖已添加到 go.mod。

**前置任务：** Task Migrate-1（golang-migrate 依赖安装）✅

**工作目录：** `C:\workspace\go_project\bhz`

**Go 工具链：** `/c/go/bin/go`

**当前 git HEAD：** 6746c11

---

## 你的任务

### Step 1: 添加 MigrationConfig 结构

在 `config/config.go` 中添加：

```go
// MigrationConfig 数据库迁移配置
type MigrationConfig struct {
    Enabled         bool   `json:"enabled"`          // 是否启用迁移
    Dir             string `json:"dir"`              // 迁移文件目录
    AutoApply       bool   `json:"auto_apply"`       // 启动时自动执行
    TableName       string `json:"table_name"`       // 版本追踪表名
}
```

### Step 2: 修改 Config 结构体

在 `Config` 结构体末尾添加 `Migration` 字段：

```go
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

### Step 3: 在 LoadFromEnv() 中添加环境变量加载

在 `LoadFromEnv()` 函数末尾（在 `cfg.Logger.Dir` 之后）添加：

```go
cfg.Migration.Enabled = getEnvAsBool("MIGRATION_ENABLED", false)
cfg.Migration.Dir = getEnv("MIGRATION_DIR", "./migrations")
cfg.Migration.AutoApply = getEnvAsBool("MIGRATION_AUTO_APPLY", false)
cfg.Migration.TableName = getEnv("MIGRATION_TABLE_NAME", "schema_migrations")
```

### Step 4: 编写测试

在 `config/config_test.go` 中添加以下测试（如果文件不存在则创建）：

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
    if cfg.Migration.AutoApply {
        t.Errorf("Expected default Migration.AutoApply = false, got true")
    }
    if cfg.Migration.TableName != "schema_migrations" {
        t.Errorf("Expected default Migration.TableName = 'schema_migrations', got %s", cfg.Migration.TableName)
    }
}
```

### Step 5: 运行测试

```bash
cd C:/workspace/go_project/bhz/server && /c/go/bin/go test ./config -v -run TestMigrationConfig
```

**预期：** PASS（所有测试通过）

### Step 6: 提交

```bash
cd C:/workspace/go_project/bhz && git add server/config/config.go server/config/config_test.go
git commit -m "feat(config): add migration configuration

- Add MigrationConfig struct with Enabled, Dir, AutoApply, TableName fields
- Add MIGRATION_* environment variables
- Default: enabled=false, dir=./migrations, auto_apply=false, table=schema_migrations
- Add unit tests for migration config loading from environment

Co-Authored-By: Claude Code <noreply@anthropic.com>"
```

---

## 验证

- [ ] `config.go` 包含 `MigrationConfig` 结构体
- [ ] `Config` 结构体包含 `Migration` 字段
- [ ] `LoadFromEnv()` 加载 MIGRATION_* 环境变量
- [ ] 所有测试通过（`TestMigrationConfigFromEnv`, `TestMigrationConfigDefaults`）
- [ ] 代码已提交

---

## 报告

**报告文件：** `C:\workspace\go_project\bhz\.superpowers\sdd\task-migrate-2-report.md`

**报告必须包含：**
- 状态：DONE 或 NEEDS_CONTEXT
- git commit SHA
- 测试结果（实际输出）
- 任何问题或关注点

**提交：** 完成任务后请自行提交到 git。
