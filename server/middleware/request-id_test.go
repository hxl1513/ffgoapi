package middleware

import (
	"net/http"
	"testing"

	"github.com/gin-gonic/gin"
)

type testingResponseWriter struct {
	header http.Header
}

func newTestingResponseWriter() *testingResponseWriter {
	return &testingResponseWriter{header: http.Header{}}
}

func (w *testingResponseWriter) Header() http.Header { return w.header }
func (w *testingResponseWriter) Write([]byte) (int, error) { return 0, nil }
func (w *testingResponseWriter) WriteHeader(int) {}

func TestRequestID(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := newTestingResponseWriter()
	req, _ := http.NewRequest("GET", "/test", nil)
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	handler := RequestID()
	handler(c)

	requestID := c.GetString("X-Request-ID")
	if requestID == "" { t.Error("RequestID not set") }

	if w.Header().Get("X-Request-ID") == "" { t.Error("RequestID not in header") }
}
