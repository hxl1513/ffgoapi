# Task 5 Brief: 示例迁移文件

**目标：** 创建示例迁移文件展示最佳实践

**文件：**
- 创建：`server/migration/example/20240101120000_init.up.sql`
- 创建：`server/migration/example/20240101120000_init.down.sql`

---

## 上下文

这是数据库迁移功能的第五个任务。Task 4（文档改进）已完成。

**前置任务：** Task Migrate-1 ✅, Task Migrate-2 ✅, Task Migrate-3 ⚠️, Task Migrate-4 ✅

**当前 git HEAD：** 3e23f8e

**工作目录：** `C:\workspace\go_project\bhz`

---

## 你的任务

### Step 1: 创建目录

```bash
mkdir -p C:/workspace/go_project/bhz/server/migration/example
```

### Step 2: 创建 Up 迁移文件

创建 `server/migration/example/20240101120000_init.up.sql`：

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

**注意：** 这是 MySQL 语法。如果需要支持 PostgreSQL，可以在文档中说明差异。

### Step 3: 创建 Down 迁移文件

创建 `server/migration/example/20240101120000_init.down.sql`：

```sql
-- ========================================
-- Rollback: 20240101120000_init
-- Description: 删除 users 表
-- ========================================

DROP TABLE IF EXISTS users;
```

### Step 4: 验证文件

```bash
ls -la C:/workspace/go_project/bhz/server/migration/example/
```

应看到两个 `.sql` 文件。

### Step 5: 提交

```bash
cd C:/workspace/go_project/bhz && git add server/migration/example/
git commit -m "docs(migration): add example migration files

- Add 20240101120000_init.up.sql (create users table with MySQL 8.0 syntax)
- Add 20240101120000_init.down.sql (drop users table)
- Include comments and field documentation
- Demonstrate best practices (indexes, comments, soft delete support)

Co-Authored-By: Claude Code <noreply@anthropic.com>"
```

---

## 验证

- [ ] 目录 `server/migration/example/` 存在
- [ ] `20240101120000_init.up.sql` 文件存在且内容正确
- [ ] `20240101120000_init.down.sql` 文件存在且内容正确
- [ ] 文件包含适当的注释
- [ ] 代码已提交

---

## 报告

**报告文件：** `C:\workspace\go_project\bhz\.superpowers\sdd\task-migrate-5-report.md`

**报告必须包含：**
- 状态：DONE
- git commit SHA
- 文件列表
- 任何问题或关注点

**提交：** 完成任务后请自行提交到 git。
