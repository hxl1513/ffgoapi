package service

import (
	"context"
	"fmt"

	"github.com/hxl1513/ffgoapi/model"
)

// BaseService provides common database and cache operations.
type BaseService struct {
	db    *model.DB
	cache model.Cache
}

// NewBaseService creates a new BaseService with the given database connection.
func NewBaseService(db *model.DB) *BaseService {
	return &BaseService{db: db}
}

// WithCache attaches a cache implementation and returns the service for chaining.
func (s *BaseService) WithCache(cache model.Cache) *BaseService {
	s.cache = cache
	return s
}

// DB returns the underlying database handle.
func (s *BaseService) DB() *model.DB {
	return s.db
}

// Cache returns the attached cache implementation (may be nil).
func (s *BaseService) Cache() model.Cache {
	return s.cache
}

// WithTransaction executes fn inside a write transaction.
// It commits on success or rolls back on error.
func (s *BaseService) WithTransaction(ctx context.Context, fn func(tx *model.DB) error) error {
	tx := s.db.Write.WithContext(ctx).Begin()
	if tx.Error != nil {
		return fmt.Errorf("failed to begin transaction: %w", tx.Error)
	}

	txDB := &model.DB{Write: tx, Read: tx}

	if err := fn(txDB); err != nil {
		if rbErr := tx.Rollback().Error; rbErr != nil {
			return fmt.Errorf("transaction error: %v, rollback error: %w", err, rbErr)
		}
		return err
	}

	if err := tx.Commit().Error; err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}
