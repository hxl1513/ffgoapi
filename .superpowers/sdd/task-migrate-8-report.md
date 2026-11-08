# Task 8 Report: 数据库迁移文档

**状态：** DONE

**git commit SHA：** `915a633`

---

## 文件列表

| 文件 | 操作 | 说明 |
|------|------|------|
| `server/docs/MIGRATION.md` | 新建 | 完整迁移指南（386 行） |
| `server/README.md` | 修改 | 新增"数据库迁移"章节（英文） |
| `server/README.zh-CN.md` | 修改 | 新增"数据库迁移"章节（中文） |

---

## 完成内容

### `server/docs/MIGRATION.md`

包含以下章节：
- **快速开始** — 环境变量配置、目录创建、执行迁移
- **迁移文件命名** — `{时间戳}_{描述}.{方向}.sql` 格式说明
- **编写迁移** — MySQL 8.0 和 PostgreSQL 16 完整 Up/Down 示例
- **执行迁移** — CLI、代码调用、自动执行三种方式
- **回滚迁移** — 回滚到指定版本、修复脏版本
- **跨数据库兼容性** — 条件 SQL 和分离文件两种方案
- **CLI 命令** — up/down/status/force 命令说明及配置项表
- **最佳实践** — 9 条实践建议

### `server/README.md` / `server/README.zh-CN.md`

在两个 README 的 "License" 之前新增"数据库迁移"章节，包含：
- 功能简介（基于 golang-migrate，支持 MySQL 8.0 / PostgreSQL 16）
- 快速开始示例（环境变量 + SQL 示例）
- 执行命令
- 链接到详细文档 `docs/MIGRATION.md`

---

## 问题与关注点

- 无阻塞问题，所有步骤顺利完成
- MIGRATION.md 内容完整，覆盖了设计文档第 11 节的所有示例和配置信息
