package middleware

import (
	"crypto/rand"
	"encoding/hex"

	"github.com/gin-gonic/gin"
)

// RequestID returns a gin.HandlerFunc that injects a unique X-Request-ID into
// the response header and gin context when one is not already present.
func RequestID() gin.HandlerFunc {
	return func(c *gin.Context) {
		rid := c.GetHeader("X-Request-ID")
		if rid == "" {
			buf := make([]byte, 16)
			_, _ = rand.Read(buf)
			rid = hex.EncodeToString(buf)
		}
		c.Set("X-Request-ID", rid)
		c.Header("X-Request-ID", rid)
		c.Next()
	}
}
