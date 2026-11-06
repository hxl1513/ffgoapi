package examples

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/your-org/go-web-api-framework/dto"
	"github.com/your-org/go-web-api-framework/examples/users"
	"github.com/your-org/go-web-api-framework/model"
	"github.com/your-org/go-web-api-framework/types"
)

// setupTestDB creates an in-memory SQLite database with the users table
// migrated and returns the *model.DB handle.  Callers must close it.
func setupTestDB(t *testing.T) *model.DB {
	t.Helper()
	db, err := model.NewDB(&model.DatabaseConfig{
		Type:   "sqlite",
		DBName: ":memory:",
	})
	if err != nil {
		t.Fatalf("NewDB() error = %v", err)
	}
	if err := db.Write.AutoMigrate(&model.User{}); err != nil {
		t.Fatalf("AutoMigrate() error = %v", err)
	}
	return db
}

// decodeResponse unmarshals the JSON body of w into target.
func decodeResponse(t *testing.T, w *httptest.ResponseRecorder, target interface{}) {
	t.Helper()
	if err := json.Unmarshal(w.Body.Bytes(), target); err != nil {
		t.Fatalf("json.Unmarshal: %v", err)
	}
}

// ── DTO validation ──────────────────────────────────────────────────────────

func TestCreateUserRequest_Validate(t *testing.T) {
	tests := []struct {
		name    string
		req     dto.CreateUserRequest
		wantErr bool
	}{
		{
			name:    "valid request",
			req:     dto.CreateUserRequest{Name: "Alice", Email: "alice@example.com", Password: "longenough"},
			wantErr: false,
		},
		{name: "empty name", req: dto.CreateUserRequest{Name: "", Email: "a@b.com", Password: "longenough"}, wantErr: true},
		{name: "short name", req: dto.CreateUserRequest{Name: "A", Email: "a@b.com", Password: "longenough"}, wantErr: true},
		{name: "invalid email", req: dto.CreateUserRequest{Name: "Bob", Email: "not-an-email", Password: "longenough"}, wantErr: true},
		{name: "short password", req: dto.CreateUserRequest{Name: "Bob", Email: "b@c.com", Password: "short"}, wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.req.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

// ── DTO conversion ──────────────────────────────────────────────────────────

func TestFromUser(t *testing.T) {
	user := &model.User{
		BaseModel: model.BaseModel{
			ID:        7,
			CreatedAt: time.Date(2024, 1, 2, 15, 4, 5, 0, time.UTC),
			UpdatedAt: time.Date(2024, 1, 2, 15, 4, 5, 0, time.UTC),
		},
		Name:     "Bob",
		Email:    "bob@example.com",
		Password: "hashed_xyz",
		Role:     1,
	}

	resp := dto.FromUser(user)
	if resp == nil {
		t.Fatal("FromUser() returned nil")
	}
	if resp.ID != 7 {
		t.Errorf("ID = %d; want 7", resp.ID)
	}
	if resp.Name != "Bob" {
		t.Errorf("Name = %q; want %q", resp.Name, "Bob")
	}
	if resp.Email != "bob@example.com" {
		t.Errorf("Email = %q; want %q", resp.Email, "bob@example.com")
	}
	if resp.Role != 1 {
		t.Errorf("Role = %d; want 1", resp.Role)
	}
	if resp.CreatedAt == "" {
		t.Error("CreatedAt is empty")
	}
}

func TestFromUser_Nil(t *testing.T) {
	if dto.FromUser(nil) != nil {
		t.Error("FromUser(nil) should return nil")
	}
}

// ── Service-level integration ────────────────────────────────────────────────

func TestUserService_CreateAndFindByID(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	svc := users.NewUserService(db)

	req := &dto.CreateUserRequest{
		Name:     "Alice Example",
		Email:    "alice@example.com",
		Password: "secret123",
	}

	// Create
	created, err := svc.CreateUser(serviceCtx(), req)
	if err != nil {
		t.Fatalf("CreateUser() error = %v", err)
	}
	if created.ID == 0 {
		t.Error("CreateUser() returned user with ID == 0")
	}
	if created.Name != req.Name {
		t.Errorf("Name = %q; want %q", created.Name, req.Name)
	}
	if created.Email != req.Email {
		t.Errorf("Email = %q; want %q", created.Email, req.Email)
	}
	if created.Role != 2 {
		t.Errorf("Role = %d; want 2", created.Role)
	}
	if created.CreatedAt == "" {
		t.Error("CreatedAt is empty")
	}

	// Retrieve
	found, err := svc.Repo.FindByID(serviceCtx(), created.ID)
	if err != nil {
		t.Fatalf("FindByID() error = %v", err)
	}
	if found == nil {
		t.Fatal("FindByID() returned nil")
	}
	if found.Name != req.Name {
		t.Errorf("Name = %q; want %q", found.Name, req.Name)
	}
	if found.Email != req.Email {
		t.Errorf("Email = %q; want %q", found.Email, req.Email)
	}
	if found.Password != "hashed_secret123" {
		t.Errorf("Password = %q; want %q", found.Password, "hashed_secret123")
	}
}

// ── Controller-level integration via gin test context ───────────────────────

func TestUserController_CreateUser(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db := setupTestDB(t)
	defer db.Close()

	svc := users.NewUserService(db)
	ctrl := users.NewUserController(svc)

	r := gin.New()
	ctrl.RegisterRoutes(r)

	body := []byte(`{"name":"Test User","email":"test@example.com","password":"mypassword"}`)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/users", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Fatalf("Status = %d; want %d; body = %s", w.Code, http.StatusCreated, w.Body.String())
	}

	var resp types.Response[dto.UserResponse]
	decodeResponse(t, w, &resp)
	if !resp.Success {
		t.Error("Success should be true")
	}
	if resp.Data.ID == 0 {
		t.Error("Response Data.ID should not be 0")
	}
	if resp.Data.Name != "Test User" {
		t.Errorf("Name = %q; want %q", resp.Data.Name, "Test User")
	}
	if resp.Data.Email != "test@example.com" {
		t.Errorf("Email = %q; want %q", resp.Data.Email, "test@example.com")
	}
	if resp.Data.Role != 2 {
		t.Errorf("Role = %d; want 2", resp.Data.Role)
	}
}

func TestUserController_GetUser(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db := setupTestDB(t)
	defer db.Close()

	svc := users.NewUserService(db)
	ctrl := users.NewUserController(svc)

	created, err := svc.CreateUser(serviceCtx(), &dto.CreateUserRequest{
		Name:     "FindMe",
		Email:    "findme@example.com",
		Password: "pass1234",
	})
	if err != nil {
		t.Fatalf("CreateUser() error = %v", err)
	}

	r := gin.New()
	ctrl.RegisterRoutes(r)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/users/"+itoa(int(created.ID)), nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("Status = %d; want %d; body = %s", w.Code, http.StatusOK, w.Body.String())
	}

	var resp types.Response[dto.UserResponse]
	decodeResponse(t, w, &resp)
	if !resp.Success {
		t.Error("Success should be true")
	}
	if resp.Data.ID != created.ID {
		t.Errorf("ID = %d; want %d", resp.Data.ID, created.ID)
	}
	if resp.Data.Name != "FindMe" {
		t.Errorf("Name = %q; want %q", resp.Data.Name, "FindMe")
	}
}

func TestUserController_GetUser_NotFound(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db := setupTestDB(t)
	defer db.Close()

	svc := users.NewUserService(db)
	ctrl := users.NewUserController(svc)

	r := gin.New()
	ctrl.RegisterRoutes(r)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/users/9999", nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Fatalf("Status = %d; want %d", w.Code, http.StatusNotFound)
	}
}

// ── helpers ─────────────────────────────────────────────────────────────────

// serviceCtx returns a minimal context.Context for service-layer calls.
func serviceCtx() *gin.Context {
	c, _ := gin.CreateTestContext(httptest.NewRecorder())
	return c
}

// itoa converts a non-negative integer to its decimal string form.
func itoa(n int) string {
	if n == 0 {
		return "0"
	}
	var b []byte
	for n > 0 {
		b = append(b, byte('0'+n%10))
		n /= 10
	}
	// reverse
	for i, j := 0, len(b)-1; i < j; i, j = i+1, j-1 {
		b[i], b[j] = b[j], b[i]
	}
	return string(b)
}
