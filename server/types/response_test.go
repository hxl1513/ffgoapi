package types

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestResponseSuccess(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	type Dummy struct{ Name string }
	SuccessResponse(c, Dummy{Name: "hello"})

	if w.Code != http.StatusOK {
		t.Errorf("Status = %d; want %d", w.Code, http.StatusOK)
	}
	var resp Response[Dummy]
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("json.Unmarshal: %v", err)
	}
	if !resp.Success {
		t.Error("Success should be true")
	}
	if resp.Data.Name != "hello" {
		t.Errorf("Data.Name = %q; want %q", resp.Data.Name, "hello")
	}
}

func TestResponseSuccessWithMessage(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	SuccessResponseWithMessage(c, "all good", 42)

	if w.Code != http.StatusOK {
		t.Errorf("Status = %d; want %d", w.Code, http.StatusOK)
	}
	var resp Response[int]
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("json.Unmarshal: %v", err)
	}
	if resp.Message != "all good" {
		t.Errorf("Message = %q; want %q", resp.Message, "all good")
	}
	if resp.Data != 42 {
		t.Errorf("Data = %v; want %v", resp.Data, 42)
	}
}

func TestErrorResponse(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	appErr := NewAppError(ErrCodeNotFound, "missing", errors.New("not here"))
	ErrorResponse(c, appErr)

	if w.Code != http.StatusNotFound {
		t.Errorf("Status = %d; want %d", w.Code, http.StatusNotFound)
	}
	var resp Response[any]
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("json.Unmarshal: %v", err)
	}
	if resp.Success {
		t.Error("Success should be false")
	}
	if resp.Message != "missing" {
		t.Errorf("Message = %q; want %q", resp.Message, "missing")
	}
}

func TestPaginatedResponse(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	items := []string{"a", "b", "c"}
	PaginatedResponse(c, items, 25, 2, 10)

	if w.Code != http.StatusOK {
		t.Errorf("Status = %d; want %d", w.Code, http.StatusOK)
	}
	var resp PaginatedResult[string]
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("json.Unmarshal: %v", err)
	}
	if len(resp.Data) != 3 {
		t.Errorf("len(Data) = %d; want %d", len(resp.Data), 3)
	}
	if resp.Total != 25 {
		t.Errorf("Total = %d; want %d", resp.Total, 25)
	}
	if resp.Page != 2 {
		t.Errorf("Page = %d; want %d", resp.Page, 2)
	}
	if resp.PerPage != 10 {
		t.Errorf("PerPage = %d; want %d", resp.PerPage, 10)
	}
}
