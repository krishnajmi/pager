package server

import (
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/kp/pager/login"
)

// AuthPermissionMiddleware checks for valid auth token and required permissions
// This middleware ensures that if any auth or permission check fails,
// the request is aborted and the API handler is not called
func AuthPermissionMiddleware(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// First check if this endpoint has any required permissions
		requiredPerms, exists := GetCachedPermissions(c.Request.Method, c.FullPath())
		if !exists || len(requiredPerms) == 0 {
			// No permissions required - skip all checks
			c.Next()
			return
		}

		// Step 1: Check authentication (only if permissions are required)
		authHeader := c.GetHeader("X-Auth-Token")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		claims, err := login.ValidateToken(tokenString)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			return
		}

		// Add claims to context
		c.Set("username", claims.Username)
		c.Set("user_type", claims.UserType)
		c.Set("permissions", claims.Permissions)

		// Log username and API call
		logger := c.Value("logger").(*slog.Logger)
		logger.Info("API request",
			slog.String("username", claims.Username),
			slog.String("api", c.Request.URL.Path),
			slog.String("method", c.Request.Method),
			slog.Time("time", time.Now()),
		)

		// Admin users bypass permission checks
		if claims.UserType == "admin" {
			c.Next()
			return
		}

		// Step 2: Check permissions
		for _, reqPerm := range requiredPerms {
			hasPerm := false
			for _, perm := range claims.Permissions {
				if perm == reqPerm {
					hasPerm = true
					break
				}
			}
			if !hasPerm {
				c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Insufficient permissions"})
				return
			}
		}

		// Only proceed to the API handler if all checks pass
		c.Next()
	}
}
