# Task 1 Brief: 安装 golang-migrate 依赖

**目标：** 在 `server/go.mod` 中添加 golang-migrate 依赖

**文件：**
- 修改：`server/go.mod`
- 修改：`server/go.sum`（由 go mod tidy 生成）

---

## 上下文

这是数据库迁移功能的第一个任务，在现有的 Go Web API 框架基础上添加 golang-migrate 支持。

**前置任务：** 无（这是迁移功能的第一项任务）

**工作目录：** `C:\workspace\go_project\bhz`

**Go 工具链：** `/c/go/bin/go`（使用此完整路径执行所有 go 命令）

---

## 你的任务

### Step 1: 添加 golang-migrate 依赖

在 `server/` 目录下执行：

```bash
cd server && /c/go/bin/go get github.com/golang-migrate/migrate/v4
```

**预期结果：**
- `go.mod` 中新增 `github.com/golang-migrate/migrate/v4` 依赖
- `go.sum` 自动更新

### Step 2: 验证

验证依赖已添加：

```bash
cd server && grep "golang-migrate" go.mod
```

应看到类似：
```
github.com/golang-migrate/migrate/v4 v4.17.0
```

### Step 3: 提交

提交更改：

```bash
cd C:/workspace/go_project/bhz && git add server/go.mod server/go.sum
git commit -m "chore: add golang-migrate dependency

- Add github.com/golang-migrate/migrate/v4
- Prepare for database migration feature

Co-Authored-By: Claude Code <noreply@anthropic.com>"
```

---

## 验证

- [ ] `go.mod` 包含 `github.com/golang-migrate/migrate/v4`
- [ ] `go.sum` 已更新
- [ ] 代码已提交

---

## 报告

**报告文件：** `C:\workspace\go_project\bhz\.superpowers\sdd\task-migrate-1-report.md`

**报告必须包含：**
- 状态：DONE 或 NEEDS_CONTEXT
- git commit SHA
- go.mod 相关行（验证依赖已添加）
- 任何问题或关注点

**不要求：** 运行 `go mod tidy`（网络可能不可用）
