package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/hxl1513/ffgoapi/model"
	"github.com/hxl1513/ffgoapi/pkg/crypto"
	"github.com/hxl1513/ffgoapi/types"
)

// AuthConfig holds the authentication configuration.
type AuthConfig struct {
	SessionSecret string
	TokenExpiry   int
}

// AuthMiddleware returns a gin.HandlerFunc that authenticates requests via JWT Bearer tokens.
// allowedRoles is a variadic list of permitted role IDs; passing 0 or none means any authenticated
// user is accepted. Session auth is a placeholder (always returns false for now).
func AuthMiddleware(cfg *AuthConfig, cache model.Cache, allowedRoles ...int) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 1. Try session auth (placeholder - not yet implemented)
		// Skipping; always falls through to Bearer token.

		// 2. Try Bearer token from Authorization header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, types.Response[any]{
				Success: false,
				Message: "missing authorization token",
			})
			c.Abort()
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
			c.JSON(http.StatusUnauthorized, types.Response[any]{
				Success: false,
				Message: "invalid authorization header format",
			})
			c.Abort()
			return
		}

		token := parts[1]
		if token == "" {
			c.JSON(http.StatusUnauthorized, types.Response[any]{
				Success: false,
				Message: "missing authorization token",
			})
			c.Abort()
			return
		}

		// 3. Parse JWT
		claims, err := crypto.ParseToken(cfg.SessionSecret, token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, types.Response[any]{
				Success: false,
				Message: "invalid or expired token",
			})
			c.Abort()
			return
		}

		// 4. Check role if non-zero allowedRoles specified
		hasRestriction := false
		for _, r := range allowedRoles {
			if r != 0 {
				hasRestriction = true
				break
			}
		}
		if hasRestriction {
			allowed := false
			for _, r := range allowedRoles {
				if claims.Role == r {
					allowed = true
					break
				}
			}
			if !allowed {
				c.JSON(http.StatusForbidden, types.Response[any]{
					Success: false,
					Message: "insufficient permissions",
				})
				c.Abort()
				return
			}
		}

		// 5. Set user context
		c.Set("user_id", claims.UserID)
		c.Set("username", claims.Username)
		c.Set("role", claims.Role)

		c.Next()
	}
}

// AdminAuth returns a middleware that only allows admin role (role == 1).
func AdminAuth(cfg *AuthConfig, cache model.Cache) gin.HandlerFunc {
	return AuthMiddleware(cfg, cache, 1)
}

// UserAuth returns a middleware that allows authenticated regular users (role == 2).
func UserAuth(cfg *AuthConfig, cache model.Cache) gin.HandlerFunc {
	return AuthMiddleware(cfg, cache, 2)
}
