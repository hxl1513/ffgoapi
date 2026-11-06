package service

import (
	"context"
	"testing"

	"github.com/your-org/go-web-api-framework/model"
)

func TestBaseService_DB(t *testing.T) {
	cfg := &model.DatabaseConfig{
		Type:   "sqlite",
		DBName: ":memory:",
	}
	db, err := model.NewDB(cfg)
	if err != nil {
		t.Fatalf("NewDB() error = %v", err)
	}
	defer db.Close()

	svc := NewBaseService(db)
	if svc.DB() == nil {
		t.Error("DB() returned nil")
	}
}

func TestBaseService_Transaction(t *testing.T) {
	cfg := &model.DatabaseConfig{
		Type:   "sqlite",
		DBName: ":memory:",
	}
	db, err := model.NewDB(cfg)
	if err != nil {
		t.Fatalf("NewDB() error = %v", err)
	}
	defer db.Close()

	svc := NewBaseService(db)

	// Success path
	err = svc.WithTransaction(context.Background(), func(tx *model.DB) error {
		return nil
	})
	if err != nil {
		t.Errorf("WithTransaction() error = %v", err)
	}

	// Rollback path
	testErr := &testingError{msg: "test error"}
	err = svc.WithTransaction(context.Background(), func(tx *model.DB) error {
		return testErr
	})
	if err != testErr {
		t.Errorf("WithTransaction() error = %v, want %v", err, testErr)
	}
}

type testingError struct {
	msg string
}

func (e *testingError) Error() string { return e.msg }
