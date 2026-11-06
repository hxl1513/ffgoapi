package model

import (
	"testing"
)

func TestNewDB_SQLite(t *testing.T) {
	cfg := &DatabaseConfig{
		Type:   "sqlite",
		DBName: ":memory:",
	}

	db, err := NewDB(cfg)
	if err != nil {
		t.Fatalf("NewDB() error = %v", err)
	}
	defer db.Close()

	if db.Write == nil {
		t.Error("DB.Write is nil")
	}
	if db.Read == nil {
		t.Error("DB.Read is nil")
	}
}

func TestBaseModelTimestamps(t *testing.T) {
	model := &BaseModel{}
	if model.ID != 0 {
		t.Error("ID should be zero value")
	}
}

func TestPaginate(t *testing.T) {
	items := []string{"a", "b", "c", "d", "e"}
	total := int64(len(items))

	result := Paginate(items, total, 1, 2)
	if len(result.Items) != 2 {
		t.Errorf("Items length = %d, want 2", len(result.Items))
	}
	if result.Total != total {
		t.Errorf("Total = %d, want %d", result.Total, total)
	}
	if result.Page != 1 {
		t.Errorf("Page = %d, want 1", result.Page)
	}
	if result.PerPage != 2 {
		t.Errorf("PerPage = %d, want 2", result.PerPage)
	}
	if result.TotalPages != 3 {
		t.Errorf("TotalPages = %d, want 3", result.TotalPages)
	}
}
