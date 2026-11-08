package config

import (
	"fmt"
	"time"
)

type Config struct {
	Server      ServerConfig
	Database    DatabaseConfig
	Redis       RedisConfig
	Auth        AuthConfig
	RateLimit   RateLimitConfig
	Logger      LoggerConfig
	Migration   MigrationConfig
}

type ServerConfig struct {
	Port         int           `json:"port"`
	Mode         string        `json:"mode"`
	ReadTimeout  time.Duration `json:"read_timeout"`
	WriteTimeout time.Duration `json:"write_timeout"`
}

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

type RedisConfig struct {
	Enabled  bool   `json:"enabled"`
	URL      string `json:"url"`
	PoolSize int    `json:"pool_size"`
}

type AuthConfig struct {
	SessionSecret string `json:"session_secret"`
	TokenExpiry   int    `json:"token_expiry"`
}

type RateLimitConfig struct {
	Enabled         bool `json:"enabled"`
	DefaultMax      int  `json:"default_max"`
	DefaultDuration int  `json:"default_duration"`
}

type LoggerConfig struct {
	Level string `json:"level"`
	Dir   string `json:"dir"`
}

// MigrationConfig 数据库迁移配置
type MigrationConfig struct {
	Enabled   bool   `json:"enabled"`    // 是否启用迁移
	Dir       string `json:"dir"`        // 迁移文件目录
	AutoApply bool   `json:"auto_apply"` // 启动时自动执行
	TableName string `json:"table_name"` // 版本追踪表名
}

func LoadFromEnv() (*Config, error) {
	cfg := &Config{}

	cfg.Server.Port = getEnvAsInt("PORT", 3000)
	cfg.Server.Mode = getEnv("GIN_MODE", "release")
	cfg.Server.ReadTimeout = getEnvAsDuration("READ_TIMEOUT", 10*time.Second)
	cfg.Server.WriteTimeout = getEnvAsDuration("WRITE_TIMEOUT", 10*time.Second)

	cfg.Database.Type = getEnv("DB_TYPE", "sqlite")
	cfg.Database.Host = getEnv("DB_HOST", "localhost")
	cfg.Database.Port = getEnvAsInt("DB_PORT", 3306)
	cfg.Database.User = getEnv("DB_USER", "root")
	cfg.Database.Password = getEnv("DB_PASSWORD", "")
	cfg.Database.DBName = getEnv("DB_NAME", "app.db")
	cfg.Database.MaxIdleConns = getEnvAsInt("DB_MAX_IDLE_CONNS", 10)
	cfg.Database.MaxOpenConns = getEnvAsInt("DB_MAX_OPEN_CONNS", 100)
	cfg.Database.LogLevel = getEnv("DB_LOG_LEVEL", "error")

	cfg.Redis.Enabled = getEnvAsBool("REDIS_ENABLED", false)
	cfg.Redis.URL = getEnv("REDIS_URL", "")
	cfg.Redis.PoolSize = getEnvAsInt("REDIS_POOL_SIZE", 10)

	cfg.Auth.SessionSecret = getEnv("SESSION_SECRET", "")
	cfg.Auth.TokenExpiry = getEnvAsInt("TOKEN_EXPIRY", 86400)

	cfg.RateLimit.Enabled = getEnvAsBool("RATE_LIMIT_ENABLED", true)
	cfg.RateLimit.DefaultMax = getEnvAsInt("RATE_LIMIT_DEFAULT_MAX", 100)
	cfg.RateLimit.DefaultDuration = getEnvAsInt("RATE_LIMIT_DEFAULT_DURATION", 60)

	cfg.Logger.Level = getEnv("LOG_LEVEL", "info")
	cfg.Logger.Dir = getEnv("LOG_DIR", "./logs")

	cfg.Migration.Enabled = getEnvAsBool("MIGRATION_ENABLED", false)
	cfg.Migration.Dir = getEnv("MIGRATION_DIR", "./migrations")
	cfg.Migration.AutoApply = getEnvAsBool("MIGRATION_AUTO_APPLY", false)
	cfg.Migration.TableName = getEnv("MIGRATION_TABLE_NAME", "schema_migrations")

	return cfg, nil
}

type ValidationError struct {
	Field   string
	Message string
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("validation error on field %s: %s", e.Field, e.Message)
}

func ValidateConfig(cfg *Config) error {
	if cfg.Auth.SessionSecret == "" {
		return &ValidationError{Field: "auth.session_secret", Message: "session secret is required"}
	}
	return nil
}
