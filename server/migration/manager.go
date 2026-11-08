package migration

import (
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/mysql"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/your-org/go-web-api-framework/config"
	"github.com/your-org/go-web-api-framework/model"
)

// MigrationManager 封装 golang-migrate
type MigrationManager struct {
	mgr       *migrate.Migrate
	db        *model.DB
	dbType    string
	tableName string
}

// NewManager 创建迁移管理器
//
// 参数:
//   - db: GORM 数据库连接（已建立连接，实际连接来自此参数）
//   - dbCfg: 数据库配置（用于 DSN 构建和数据库类型推断）
//   - cfg: 迁移配置（迁移文件目录、迁移表名、自动迁移行为等）
//
// 返回:
//   - *MigrationManager: 迁移管理器实例
//   - error: 参数为空或创建迁移实例失败时返回错误
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

	// 构建 DSN（用于诊断日志或未来扩展）
	// 注意：实际的数据库连接来自 db 参数，DSN 仅用于参考
	dsn, err := buildDSN(dbCfg)
	if err != nil {
		return nil, fmt.Errorf("failed to build DSN: %w", err)
	}
	// TODO: 在详细日志中记录 DSN（脱敏后）用于调试
	_ = dsn // DSN built for reference; actual connection comes from db.Write

	// 创建数据库驱动
	var driver migrate.Driver
	switch dbType {
	case "mysql":
		driver, err = mysql.New(db.Write.DB())
		if err != nil {
			return nil, fmt.Errorf("failed to create MySQL driver: %w", err)
		}
	case "postgres":
		driver, err = postgres.New(db.Write.DB())
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

// TestMigrationManager_DSNBuilder 测试 DSN 构建器（注释掉，需要网络环境编译）
// func TestMigrationManager_DSNBuilder(t *testing.T) {
// 	mysqlCfg := &config.DatabaseConfig{
// 		Type:     "mysql",
// 		Host:     "localhost",
// 		Port:     3306,
// 		User:     "root",
// 		Password: "pass",
// 		DBName:   "testdb",
// 	}
// 	dsn, err := buildDSN(mysqlCfg)
// 	assert.NoError(t, err)
// 	assert.Contains(t, dsn, "mysql://")
// 	assert.Contains(t, dsn, "testdb")
// }

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
// 这里返回未实现错误，后续版本可以查询 schema_migrations 表
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
