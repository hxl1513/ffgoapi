package middleware

import (
	"github.com/gin-gonic/gin"
)

// BodyCleanup returns a gin.HandlerFunc that ensures the request body
// is closed after the request has been processed, preventing resource leaks.
func BodyCleanup() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
		if c.Request.Body != nil {
			c.Request.Body.Close()
		}
	}
}
