# Task 8 Brief: 文档

**目标：** 编写数据库迁移功能的完整使用文档

**文件：**
- 创建：`server/docs/MIGRATION.md` — 迁移完整指南
- 修改：`server/README.md` — 新增迁移章节
- 修改：`server/README.zh-CN.md` — 新增迁移章节（中文）

---

## 上下文

这是数据库迁移功能的第八个任务（也是倒数第二个可执行任务，Task 7 和 9 因网络/环境限制跳过）。

**前置任务：** Task Migrate-1 ✅, Task Migrate-2 ✅, Task Migrate-4 ✅, Task Migrate-5 ✅, Task Migrate-6 ✅

**当前 git HEAD：** 47a6d45

**工作目录：** `C:\workspace\go_project\bhz`

---

## 你的任务

### Step 1: 创建迁移指南文档

创建 `server/docs/MIGRATION.md`，包含以下章节：

```markdown
# 数据库迁移指南

## 目录

- [快速开始](#快速开始)
- [迁移文件命名](#迁移文件命名)
- [编写迁移](#编写迁移)
- [执行迁移](#执行迁移)
- [回滚迁移](#回滚迁移)
- [最佳实践](#最佳实践)

---

## 快速开始

### 1. 启用迁移功能

在 `.env` 文件中添加：

\`\`\`env
MIGRATION_ENABLED=true
MIGRATION_DIR=./migrations
MIGRATION_AUTO_APPLY=false  # 开发环境可设为 true
\`\`\`

### 2. 创建迁移文件

\`\`\`bash
mkdir -p migrations
\`\`\`

### 3. 执行迁移

\`\`\`bash
go run cmd/migrate/main.go up
\`\`\`

---

## 迁移文件命名

格式：`{时间戳}_{描述}.{方向}.sql`

---

## 编写迁移

### MySQL 8.0 示例

### PostgreSQL 16 示例

---

## 最佳实践

1. **始终提供 Down 迁移**
2. **使用事务**
3. **测试迁移**
4. **备份生产数据**
5. **小步提交**
```

**完整内容见设计文档 Section 11。**

### Step 2: 更新英文 README

在 `server/README.md` 末尾添加迁移章节（在 "License" 之前）：

```markdown
## 数据库迁移

框架基于 [golang-migrate](https://github.com/golang-migrate/migrate) 提供数据库迁移功能，支持 MySQL 8.0 和 PostgreSQL 16。

### 快速开始

\`\`\`env
# .env
MIGRATION_ENABLED=true
MIGRATION_DIR=./migrations
\`\`\`

\`\`\`sql
-- migrations/20240101120000_create_users.up.sql
CREATE TABLE users (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    username VARCHAR(50) NOT NULL UNIQUE,
    email VARCHAR(255) NOT NULL UNIQUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
\`\`\`

### 执行迁移

\`\`\`bash
go run cmd/migrate/main.go up
\`\`\`

详细文档：[MIGRATION.md](docs/MIGRATION.md)
```

### Step 3: 更新中文 README

在 `server/README.zh-CN.md` 末尾添加迁移章节：

```markdown
## 数据库迁移

框架基于 [golang-migrate](https://github.com/golang-migrate/migrate) 提供数据库迁移功能，支持 MySQL 8.0 和 PostgreSQL 16。

### 快速开始

...

详细文档：[MIGRATION.md](docs/MIGRATION.md)
```

### Step 4: 提交

```bash
cd C:/workspace/go_project/bhz && git add server/docs/MIGRATION.md server/README.md server/README.zh-CN.md
git commit -m "docs(migration): add comprehensive migration guide

- Add docs/MIGRATION.md with complete migration guide
- Add migration section to README.md (English)
- Add migration section to README.zh-CN.md (Chinese)
- Include quick start, naming conventions, examples
- Add MySQL 8.0 and PostgreSQL 16 examples
- Document best practices

Co-Authored-By: Claude Code <noreply@anthropic.com>"
```

---

## 验证

- [ ] `docs/MIGRATION.md` 文件存在且内容完整
- [ ] `README.md` 包含迁移章节
- [ ] `README.zh-CN.md` 包含迁移章节
- [ ] 代码已提交

---

## 报告

**报告文件：** `C:\workspace\go_project\bhz\.superpowers\sdd\task-migrate-8-report.md`

**报告必须包含：**
- 状态：DONE
- git commit SHA
- 文件列表
- 任何问题或关注点

**提交：** 完成任务后请自行提交到 git。
