# Task 5 Report: 示例迁移文件

**状态：** DONE

## Git Commit

- **Commit SHA：** `dc4aa99`
- **Commit Message：** `docs(migration): add example migration files`

## 创建的文件

| 文件 | 大小 | 描述 |
|------|------|------|
| `server/migration/example/20240101120000_init.up.sql` | 1043 bytes | 创建 users 表（MySQL 8.0 语法） |
| `server/migration/example/20240101120000_init.down.sql` | 183 bytes | 删除 users 表（回滚） |

## Up 迁移内容概要

- 创建 `users` 表，包含字段：id, username, email, password_hash, role, created_at, updated_at, deleted_at
- 包含三个索引：idx_email, idx_role, idx_created_at
- 支持软删除（deleted_at 字段）
- 角色字段带默认值（role=2 表示普通用户）
- 使用 utf8mb4 字符集和 InnoDB 引擎
- 所有字段均有中文注释

## Down 迁移内容概要

- `DROP TABLE IF EXISTS users;` 安全删除 users 表

## 问题与关注点

- 无。文件创建和提交均成功。
- 注意：当前 SQL 语法为 MySQL 8.0，若需支持 PostgreSQL 需另外提供对应语法版本（已在文件注释中注明）。
