package middleware

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"runtime/debug"
)

// Recover returns a gin.HandlerFunc that recovers from panics and returns a 500 error.
func Recover() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				debug.PrintStack()
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
					"success": false,
					"message": fmt.Sprintf("panic: %v", err),
				})
			}
		}()
		c.Next()
	}
}
