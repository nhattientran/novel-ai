package middleware

import (
	"novel-ai/internal/auth"

	"github.com/gin-gonic/gin"
)

// Auth creates an authentication middleware with the given JWT secret
func Auth(jwtSecret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString, err := auth.ExtractTokenFromCookie(c)
		if err != nil {
			c.JSON(401, gin.H{"error": "unauthorized"})
			c.Abort()
			return
		}

		claims, err := auth.ValidateToken(tokenString, jwtSecret)
		if err != nil {
			c.JSON(401, gin.H{"error": "invalid token"})
			c.Abort()
			return
		}

		auth.SetUserContext(c, claims.UserID, claims.Role)
		c.Next()
	}
}

// RequireRole creates a role-based authorization middleware
func RequireRole(role string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userRole, exists := auth.GetUserRole(c)
		if !exists {
			c.JSON(401, gin.H{"error": "unauthorized"})
			c.Abort()
			return
		}

		if userRole != role {
			c.JSON(403, gin.H{"error": "forbidden"})
			c.Abort()
			return
		}

		c.Next()
	}
}

// OptionalAuth creates a middleware that sets user context if token exists,
// but doesn't require authentication
func OptionalAuth(jwtSecret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString, err := auth.ExtractTokenFromCookie(c)
		if err == nil {
			claims, err := auth.ValidateToken(tokenString, jwtSecret)
			if err == nil {
				auth.SetUserContext(c, claims.UserID, claims.Role)
			}
		}
		c.Next()
	}
}
