# 数据库迁移指南

## 目录

- [快速开始](#快速开始)
- [迁移文件命名](#迁移文件命名)
- [编写迁移](#编写迁移)
- [执行迁移](#执行迁移)
- [回滚迁移](#回滚迁移)
- [跨数据库兼容性](#跨数据库兼容性)
- [CLI 命令](#cli-命令)
- [最佳实践](#最佳实践)

---

## 快速开始

### 1. 启用迁移功能

在 `.env` 文件中添加以下配置：

```env
MIGRATION_ENABLED=true
MIGRATION_DIR=./migrations
MIGRATION_AUTO_APPLY=false  # 开发环境可设为 true
MIGRATION_TABLE_NAME=schema_migrations
```

### 2. 创建迁移目录

```bash
mkdir -p migrations
```

### 3. 创建迁移文件

```bash
# 创建升迁脚本
touch migrations/20240101120000_create_users.up.sql

# 创建回滚脚本
touch migrations/20240101120000_create_users.down.sql
```

### 4. 执行迁移

```bash
go run cmd/migrate/main.go up
```

---

## 迁移文件命名

迁移文件采用 `golang-migrate` 标准命名格式：

```
{时间戳}_{描述}.{方向}.sql
```

**示例：**

```
20240101120000_create_users_table.up.sql
20240101120000_create_users_table.down.sql
20240102150000_add_roles.up.sql
20240102150000_add_roles.down.sql
```

- **时间戳** — 14 位数字，格式 `YYYYMMDDhhmmss`，用于确定迁移执行顺序
- **描述** — 简短英文描述，建议使用下划线分隔
- **方向** — `up`（升迁）或 `down`（回滚）
- **.sql** — 文件扩展名

---

## 编写迁移

### MySQL 8.0 示例

**Up 迁移（`20240101120000_create_users.up.sql`）：**

```sql
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

**Down 迁移（`20240101120000_create_users.down.sql`）：**

```sql
-- 回滚：删除用户表

DROP TABLE IF EXISTS users;
```

### PostgreSQL 16 示例

**Up 迁移（`20240101120000_create_users.up.sql`）：**

```sql
-- 创建用户表

CREATE TABLE IF NOT EXISTS users (
    id BIGSERIAL PRIMARY KEY,
    username VARCHAR(50) NOT NULL UNIQUE,
    email VARCHAR(255) NOT NULL UNIQUE,
    password_hash VARCHAR(255) NOT NULL,
    role SMALLINT NOT NULL DEFAULT 2,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL DEFAULT NULL
);

CREATE INDEX idx_email ON users (email);
CREATE INDEX idx_role ON users (role);
CREATE INDEX idx_created_at ON users (created_at);

COMMENT ON TABLE users IS '用户表';
```

**Down 迁移（`20240101120000_create_users.down.sql`）：**

```sql
-- 回滚：删除用户表

DROP TABLE IF EXISTS users;
```

---

## 执行迁移

### 方式一：CLI 命令（推荐）

```bash
# 执行所有待处理的迁移
go run cmd/migrate/main.go up

# 查看当前迁移状态
go run cmd/migrate/main.go status

# 回滚最后一个迁移
go run cmd/migrate/main.go down

# 强制设置版本（用于修复脏版本）
go run cmd/migrate/main.go force <version>
```

### 方式二：代码中调用

```go
package main

import (
    "log"

    "github.com/your-org/go-web-api-framework/config"
    "github.com/your-org/go-web-api-framework/model"
    "github.com/your-org/go-web-api-framework/migration"
)

func main() {
    cfg, err := config.LoadFromEnv()
    if err != nil {
        log.Fatal(err)
    }

    db, err := model.NewDB(&cfg.Database)
    if err != nil {
        log.Fatal(err)
    }

    mgr, err := migration.NewManager(db, &cfg.Migration)
    if err != nil {
        log.Fatal(err)
    }
    defer mgr.Close()

    // 查看当前版本
    version, dirty, err := mgr.Status()
    if err != nil {
        log.Printf("Status check failed: %v", err)
    } else {
        log.Printf("Current version: %d, dirty: %v", version, dirty)
    }

    // 执行迁移
    if err := mgr.Up(); err != nil && err != migrate.ErrNoChange {
        log.Fatal(err)
    }

    log.Println("Migration completed!")
}
```

### 方式三：自动执行

将 `MIGRATION_AUTO_APPLY` 设为 `true`，服务器启动时会自动执行所有待处理的迁移。

```env
MIGRATION_ENABLED=true
MIGRATION_AUTO_APPLY=true
```

---

## 回滚迁移

### 回滚最后一个迁移

```bash
go run cmd/migrate/main.go down
```

### 回滚到指定版本

```bash
go run cmd/migrate/main.go force <version>
```

将版本追踪表直接设置为指定版本，跳过该版本之后的所有迁移。适用于修复失败的迁移。

### 修复脏版本

当迁移执行失败后，版本追踪表会标记为 `dirty=true`。此时可以：

1. 检查迁移文件并修复问题
2. 执行 `force` 命令将版本标记为干净：

```bash
go run cmd/migrate/main.go force <last_good_version>
```

---

## 跨数据库兼容性

当需要同时支持 MySQL 和 PostgreSQL 时，有两种方案：

### 方案 A：条件 SQL（推荐）

使用 `golang-migrate` 的 `$if` 条件语法：

```sql
-- $if(DATABASE_TYPE == 'mysql')
CREATE TABLE users (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    username VARCHAR(50) NOT NULL
);
-- $elif(DATABASE_TYPE == 'postgres')
CREATE TABLE users (
    id BIGSERIAL PRIMARY KEY,
    username VARCHAR(50) NOT NULL
);
-- $fi
```

### 方案 B：分离文件

```
migrations/
├── mysql/
│   └── 20240101120000_create_users.up.sql
└── postgres/
    └── 20240101120000_create_users.up.sql
```

使用 `--path` 参数指定对应的数据库目录。

---

## CLI 命令

### 可用命令

| 命令 | 说明 |
|------|------|
| `go run cmd/migrate/main.go up` | 执行所有待处理的迁移 |
| `go run cmd/migrate/main.go down` | 回滚最后一个迁移 |
| `go run cmd/migrate/main.go status` | 查看当前迁移状态 |
| `go run cmd/migrate/main.go force <version>` | 强制设置版本（修复脏版本） |

### 配置项说明

| 环境变量 | 类型 | 默认值 | 说明 |
|----------|------|--------|------|
| `MIGRATION_ENABLED` | bool | `false` | 是否启用迁移功能 |
| `MIGRATION_DIR` | string | `./migrations` | 迁移文件目录 |
| `MIGRATION_AUTO_APPLY` | bool | `false` | 启动时自动执行迁移 |
| `MIGRATION_TABLE_NAME` | string | `schema_migrations` | 版本追踪表名 |

---

## 最佳实践

1. **始终提供 Down 迁移** — 每个 Up 迁移都必须有对应的 Down 迁移，以便紧急回滚
2. **使用事务** — 确保迁移文件的 DDL 语句放在事务中执行（golang-migrate 已自动处理）
3. **测试迁移** — 在部署到生产环境前，在测试数据库上完整测试升迁和回滚流程
4. **备份生产数据** — 在生产环境执行迁移前，务必先备份数据库
5. **小步提交** — 每次迁移只做一个变更，便于定位问题和回滚
6. **不要修改已执行的迁移** — 已执行的迁移文件禁止修改；如需变更，创建新的迁移文件
7. **版本号递增** — 使用递增的时间戳作为版本号，避免版本冲突
8. **注释清晰** — 在迁移文件头部添加注释，说明变更内容和原因
9. **生产环境禁用自动迁移** — `MIGRATION_AUTO_APPLY` 在生产环境应保持 `false`，使用 CLI 手动执行
