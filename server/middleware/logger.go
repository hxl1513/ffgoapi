package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/hxl1513/ffgoapi/common"
)

// Logger returns a gin.HandlerFunc that logs request details using the provided logger.
func Logger(logger *common.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery

		c.Next()

		latency := time.Since(start)
		clientIP := c.ClientIP()
		method := c.Request.Method
		statusCode := c.Writer.Status()

		logger.Info("[GIN] %s | %d | %v | %s | %s %s",
			clientIP, statusCode, latency, method, path, query)
	}
}
