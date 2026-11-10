package server

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hxl1513/ffgoapi/common"
	"github.com/hxl1513/ffgoapi/config"
	"github.com/hxl1513/ffgoapi/migration"
	"github.com/hxl1513/ffgoapi/middleware"
	"github.com/hxl1513/ffgoapi/model"
)

// Server holds the HTTP server dependencies.
type Server struct {
	engine   *gin.Engine
	config   *config.Config
	logger   *common.Logger
	db       *model.DB
	cache    model.Cache
	migrator *migration.MigrationManager  // 可选：迁移管理器
	httpSrv  *http.Server
}

// NewServer constructs and configures a new Server.
func NewServer(cfg *config.Config) (*Server, error) {
	// Set Gin mode based on configuration.
	switch cfg.Server.Mode {
	case "release":
		gin.SetMode(gin.ReleaseMode)
	case "debug":
		gin.SetMode(gin.DebugMode)
	default:
		gin.SetMode(gin.ReleaseMode)
	}

	// Create engine.
	engine := gin.New()

	// Initialise logger.
	logger := common.NewLogger(gin.DefaultWriter, cfg.Logger.Level)

	// Initialise database.
	db, err := model.NewDB(toModelDBConfig(&cfg.Database))
	if err != nil {
		return nil, fmt.Errorf("failed to initialise database: %w", err)
	}

	// Initialise cache.
	cache, err := model.NewCache(cfg.Redis.URL)
	if err != nil {
		return nil, fmt.Errorf("failed to initialise cache: %w", err)
	}

	srv := &Server{
		engine: engine,
		config: cfg,
		logger: logger,
		db:     db,
		cache:  cache,
	}

	// Register middleware.
	srv.setupMiddleware()

	// Register default routes (health check).
	srv.setupDefaultRoutes()

	// 初始化迁移管理器（可选功能）
	if cfg.Migration.Enabled {
		mgr, err := migration.NewManager(srv.db, &cfg.Database, &cfg.Migration)
		if err != nil {
			// 如果迁移功能失败，记录警告但不阻止服务器启动
			srv.logger.Warn("Failed to create migration manager: %v", err)
		} else {
			srv.migrator = mgr

			// 自动执行迁移（如果配置了 AutoApply）
			if cfg.Migration.AutoApply {
				if err := mgr.Up(); err != nil {
					srv.logger.Warn("Auto migration failed: %v", err)
				} else {
					srv.logger.Info("Auto migration completed successfully")
				}
			}
		}
	}

	return srv, nil
}

// toModelDBConfig converts config.DatabaseConfig to model.DatabaseConfig.
func toModelDBConfig(c *config.DatabaseConfig) *model.DatabaseConfig {
	return &model.DatabaseConfig{
		Type:         c.Type,
		Host:         c.Host,
		Port:         c.Port,
		User:         c.User,
		Password:     c.Password,
		DBName:       c.DBName,
		MaxIdleConns: c.MaxIdleConns,
		MaxOpenConns: c.MaxOpenConns,
		LogLevel:     c.LogLevel,
	}
}

// setupMiddleware registers all middleware on the engine.
func (s *Server) setupMiddleware() {
	s.engine.Use(middleware.Recover())
	s.engine.Use(middleware.RequestID())
	s.engine.Use(middleware.Logger(s.logger))
	s.engine.Use(middleware.CORS())
}

// setupDefaultRoutes registers built-in routes such as the health check.
func (s *Server) setupDefaultRoutes() {
	s.engine.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})
}

// Engine returns the underlying Gin engine.
func (s *Server) Engine() *gin.Engine {
	return s.engine
}

// Start starts the HTTP server.
func (s *Server) Start() error {
	addr := fmt.Sprintf(":%d", s.config.Server.Port)

	s.httpSrv = &http.Server{
		Addr:         addr,
		Handler:      s.engine,
		ReadTimeout:  s.config.Server.ReadTimeout,
		WriteTimeout: s.config.Server.WriteTimeout,
	}

	return s.httpSrv.ListenAndServe()
}

// Shutdown gracefully shuts down the server.
func (s *Server) Shutdown(ctx context.Context) error {
	if s.httpSrv != nil {
		if err := s.httpSrv.Shutdown(ctx); err != nil {
			return err
		}
	}
	if s.db != nil {
		if err := s.db.Close(); err != nil {
			return err
		}
	}
	// 关闭迁移管理器
	if s.migrator != nil {
		if err := s.migrator.Close(); err != nil {
			s.logger.Warn("Failed to close migration manager: %v", err)
		}
	}
	return nil
}

// GetMigrator returns the migration manager (if initialized).
func (s *Server) GetMigrator() *migration.MigrationManager {
	return s.migrator
}
