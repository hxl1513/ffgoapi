package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

const defaultMaxBodyBytes = 4 << 20 // 4MB

// RequestBodyLimit returns a gin.HandlerFunc that wraps the request body
// with http.MaxBytesReader to enforce a maximum request body size.
// If maxBytes is <= 0, the default of 4MB is used.
func RequestBodyLimit(maxBytes int64) gin.HandlerFunc {
	return func(c *gin.Context) {
		if maxBytes <= 0 {
			maxBytes = defaultMaxBodyBytes
		}
		c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, maxBytes)
		c.Next()
		if c.Request.Body != nil {
			c.Request.Body.Close()
		}
	}
}
