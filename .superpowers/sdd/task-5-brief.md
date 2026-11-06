# Task 5 Brief: Database Abstraction

## Context

You are implementing Task 5 of 16. Tasks 1-4 are complete.

**Your task builds on:** All previous tasks (types, config, logger, common utils)

## Your Task

Implement database abstraction layer with GORM, supporting SQLite, MySQL, and PostgreSQL.

**Working directory:** `C:\workspace\go_project\bhz`

## Files to Create

1. **`server/model/database.go`** — Database configuration and initialization
2. **`server/model/base_model.go`** — Base model and pagination
3. **`server/model/database_test.go`** — Unit tests

## Requirements (verbatim from plan)

### database.go

```go
type DatabaseConfig struct {
    Type         string `json:"type"`
    Host         string `json:"host"`
    Port         int    `json:"port"`
    User         string `json:"user"`
    Password     string `json:"password"`
    DBName       string `json:"dbname"`
    MaxIdleConns int    `json:"max_idle_conns"`
    MaxOpenConns int    `json:"max_open_conns"`
    LogLevel     string `json:"log_level"`
}

type DB struct {
    Write *gorm.DB
    Read  *gorm.DB
}

type BaseModel struct {
    ID        uint           `gorm:"primaryKey" json:"id"`
    CreatedAt time.Time      `json:"created_at"`
    UpdatedAt time.Time      `json:"updated_at"`
    DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

type PaginatedResult[T any] struct {
    Items      []T   `json:"items"`
    Total      int64 `json:"total"`
    Page       int   `json:"page"`
    PerPage    int   `json:"per_page"`
    TotalPages int   `json:"total_pages"`
}

func NewDB(cfg *DatabaseConfig) (*DB, error)
func (db *DB) Close() error
func Paginate[T any](items []T, total int64, page, perPage int) PaginatedResult[T]
```

### Database Drivers

Import and support:
- `gorm.io/driver/sqlite` — for SQLite
- `gorm.io/driver/mysql` — for MySQL
- `gorm.io/driver/postgres` — for PostgreSQL

### NewDB Implementation

```go
func NewDB(cfg *DatabaseConfig) (*DB, error) {
    var dialector gorm.Dialector
    
    switch cfg.Type {
    case "sqlite":
        dialector = sqlite.Open(cfg.DBName)
    case "mysql":
        dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
            cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.DBName)
        dialector = mysql.Open(dsn)
    case "postgres":
        dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
            cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.DBName)
        dialector = postgres.Open(dsn)
    default:
        return nil, fmt.Errorf("unsupported database type: %s", cfg.Type)
    }
    
    // Configure logger
    logLevel := logger.Silent
    switch cfg.LogLevel {
    case "error": logLevel = logger.Error
    case "warn": logLevel = logger.Warn
    case "info": logLevel = logger.Info
    }
    
    db, err := gorm.Open(dialector, &gorm.Config{
        Logger: logger.Default.LogMode(logLevel),
        NowFunc: func() time.Time {
            return time.Now().UTC()
        },
    })
    // ... connection pool config
}
```

### Test Requirements

```go
func TestNewDB_SQLite(t *testing.T) {
    cfg := &DatabaseConfig{
        Type:   "sqlite",
        DBName: ":memory:",
    }
    
    db, err := NewDB(cfg)
    if err != nil { t.Fatalf("NewDB() error = %v", err) }
    defer db.Close()
    
    if db.Write == nil { t.Error("DB.Write is nil") }
}

func TestBaseModelTimestamps(t *testing.T) {
    model := &BaseModel{}
    if model.ID != 0 { t.Error("ID should be zero value") }
}

func TestPaginate(t *testing.T) {
    items := []string{"a", "b", "c", "d", "e"}
    total := int64(len(items))
    
    result := Paginate(items, total, 1, 2)
    if len(result.Items) != 2 { t.Errorf("Items length = %d, want 2", len(result.Items)) }
    if result.Total != total { t.Errorf("Total = %d, want %d", result.Total, total) }
    if result.Page != 1 { t.Errorf("Page = %d, want 1", result.Page) }
    if result.PerPage != 2 { t.Errorf("PerPage = %d, want 2", result.PerPage) }
    if result.TotalPages != 3 { t.Errorf("TotalPages = %d, want 3", result.TotalPages) }
}
```

## Global Constraints

- Go 1.25+
- GORM v1.25+
- Time functions use UTC

## Test Verification

```bash
cd server && go test ./model -v
```

Expected: All tests PASS

## Commit Message

```
feat(model): add database abstraction layer

- Implement database abstraction supporting SQLite/MySQL/PostgreSQL
- Add BaseModel with common fields (ID, timestamps, soft delete)
- Implement pagination helper
- Add unit tests for database initialization and pagination

Co-Authored-By: Claude Code <noreply@anthropic.com>
```

## Report

Write your report to: `C:\workspace\go_project\bhz\.superpowers\sdd\task-5-report.md`
