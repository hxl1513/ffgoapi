package controller

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/hxl1513/ffgoapi/types"
)

func TestBaseController_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("GET", "/test", nil)

	bc := &BaseController{}
	bc.Success(c, "test data")

	if w.Code != http.StatusOK {
		t.Errorf("Status = %d; want %d", w.Code, http.StatusOK)
	}
	var resp types.Response[interface{}]
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("json.Unmarshal: %v", err)
	}
	if !resp.Success {
		t.Error("Success should be true")
	}
	if resp.Data == nil {
		t.Error("Data should not be nil")
	}
}

func TestBaseController_SuccessWithMessage(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("GET", "/test", nil)

	bc := &BaseController{}
	bc.SuccessWithMessage(c, "all good", map[string]string{"key": "value"})

	if w.Code != http.StatusOK {
		t.Errorf("Status = %d; want %d", w.Code, http.StatusOK)
	}
	var resp types.Response[map[string]string]
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("json.Unmarshal: %v", err)
	}
	if !resp.Success {
		t.Error("Success should be true")
	}
	if resp.Message != "all good" {
		t.Errorf("Message = %q; want %q", resp.Message, "all good")
	}
}

func TestBaseController_Created(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("POST", "/test", nil)

	bc := &BaseController{}
	bc.Created(c, "created data")

	if w.Code != http.StatusCreated {
		t.Errorf("Status = %d; want %d", w.Code, http.StatusCreated)
	}
	var resp types.Response[interface{}]
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("json.Unmarshal: %v", err)
	}
	if !resp.Success {
		t.Error("Success should be true")
	}
}

func TestBaseController_BadRequest(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("GET", "/test", nil)

	bc := &BaseController{}
	bc.BadRequest(c, "invalid input", nil)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Status = %d; want %d", w.Code, http.StatusBadRequest)
	}
	var resp types.Response[interface{}]
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("json.Unmarshal: %v", err)
	}
	if resp.Success {
		t.Error("Success should be false")
	}
	if resp.Message != "invalid input" {
		t.Errorf("Message = %q; want %q", resp.Message, "invalid input")
	}
}

func TestBaseController_NotFound(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("GET", "/test", nil)

	bc := &BaseController{}
	bc.NotFound(c, "resource not found")

	if w.Code != http.StatusNotFound {
		t.Errorf("Status = %d; want %d", w.Code, http.StatusNotFound)
	}
	var resp types.Response[interface{}]
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("json.Unmarshal: %v", err)
	}
	if resp.Success {
		t.Error("Success should be false")
	}
	if resp.Message != "resource not found" {
		t.Errorf("Message = %q; want %q", resp.Message, "resource not found")
	}
}

func TestBaseController_Unauthorized(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("GET", "/test", nil)

	bc := &BaseController{}
	bc.Unauthorized(c, "unauthorized access")

	if w.Code != http.StatusUnauthorized {
		t.Errorf("Status = %d; want %d", w.Code, http.StatusUnauthorized)
	}
	var resp types.Response[interface{}]
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("json.Unmarshal: %v", err)
	}
	if resp.Success {
		t.Error("Success should be false")
	}
}

func TestBaseController_Forbidden(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("GET", "/test", nil)

	bc := &BaseController{}
	bc.Forbidden(c, "forbidden")

	if w.Code != http.StatusForbidden {
		t.Errorf("Status = %d; want %d", w.Code, http.StatusForbidden)
	}
	var resp types.Response[interface{}]
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("json.Unmarshal: %v", err)
	}
	if resp.Success {
		t.Error("Success should be false")
	}
}

func TestBaseController_InternalError(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("GET", "/test", nil)

	bc := &BaseController{}
	bc.InternalError(c, "something went wrong", nil)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("Status = %d; want %d", w.Code, http.StatusInternalServerError)
	}
	var resp types.Response[interface{}]
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("json.Unmarshal: %v", err)
	}
	if resp.Success {
		t.Error("Success should be false")
	}
	if resp.Message != "something went wrong" {
		t.Errorf("Message = %q; want %q", resp.Message, "something went wrong")
	}
}

func TestBaseController_ParamInt(t *testing.T) {
	gin.SetMode(gin.TestMode)
	c, _ := gin.CreateTestContext(nil)

	// Test valid integer param
	c.Params = gin.Params{gin.Param{Key: "id", Value: "42"}}
	bc := &BaseController{}
	result := bc.ParamInt(c, "id")
	if result != 42 {
		t.Errorf("ParamInt = %d; want 42", result)
	}

	// Test invalid integer param (should return 0)
	c.Params = gin.Params{gin.Param{Key: "id", Value: "abc"}}
	result = bc.ParamInt(c, "id")
	if result != 0 {
		t.Errorf("ParamInt with invalid value = %d; want 0", result)
	}
}

func TestBaseController_QueryInt(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("GET", "/test?page=3", nil)

	bc := &BaseController{}
	result := bc.QueryInt(c, "page", 1)
	if result != 3 {
		t.Errorf("QueryInt = %d; want 3", result)
	}

	// Test missing param returns default
	result = bc.QueryInt(c, "missing", 10)
	if result != 10 {
		t.Errorf("QueryInt with missing key = %d; want 10", result)
	}

	// Test invalid integer returns default
	c.Request, _ = http.NewRequest("GET", "/test?size=abc", nil)
	result = bc.QueryInt(c, "size", 20)
	if result != 20 {
		t.Errorf("QueryInt with invalid value = %d; want 20", result)
	}
}

func TestBaseController_BindJSON(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("POST", "/test", strings.NewReader(`{"name":"test"}`))
	c.Request.Header.Set("Content-Type", "application/json")

	type payload struct {
		Name string `json:"name"`
	}
	bc := &BaseController{}
	var p payload
	err := bc.BindJSON(c, &p)
	if err != nil {
		t.Fatalf("BindJSON error = %v; want nil", err)
	}
	if p.Name != "test" {
		t.Errorf("Name = %q; want %q", p.Name, "test")
	}
}

func TestBaseController_BindQuery(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("GET", "/test?name=hello&count=5", nil)

	type query struct {
		Name  string `form:"name"`
		Count int    `form:"count"`
	}
	bc := &BaseController{}
	var q query
	err := bc.BindQuery(c, &q)
	if err != nil {
		t.Fatalf("BindQuery error = %v; want nil", err)
	}
	if q.Name != "hello" {
		t.Errorf("Name = %q; want %q", q.Name, "hello")
	}
	if q.Count != 5 {
		t.Errorf("Count = %d; want 5", q.Count)
	}
}

func TestBaseController_Paginate(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("GET", "/test?page=2&per_page=10", nil)

	items := []string{"a", "b", "c"}
	bc := &BaseController{}
	bc.Paginate(c, items, 25, 2, 10)

	if w.Code != http.StatusOK {
		t.Errorf("Status = %d; want %d", w.Code, http.StatusOK)
	}
	var resp types.PaginatedResult[string]
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("json.Unmarshal: %v", err)
	}
	if !resp.Success {
		t.Error("Success should be true")
	}
	if len(resp.Data) != 3 {
		t.Errorf("len(Data) = %d; want 3", len(resp.Data))
	}
	if resp.Total != 25 {
		t.Errorf("Total = %d; want 25", resp.Total)
	}
	if resp.Page != 2 {
		t.Errorf("Page = %d; want 2", resp.Page)
	}
	if resp.PerPage != 10 {
		t.Errorf("PerPage = %d; want 10", resp.PerPage)
	}
}

func TestBaseController_NoContent(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("DELETE", "/test/1", nil)

	bc := &BaseController{}
	bc.NoContent(c)

	if w.Code != http.StatusNoContent {
		t.Errorf("Status = %d; want %d", w.Code, http.StatusNoContent)
	}
}
