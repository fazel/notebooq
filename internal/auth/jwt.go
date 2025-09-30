package auth

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func JWTMiddleware(secret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		auth := c.GetHeader("Authorization")
		if auth == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "missing auth header"})
			return
		}
		parts := strings.SplitN(auth, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid auth header"})
			return
		}
		tokenStr := parts[1]
		// parse token (simple check). For brevity we skip verifying claims library here
		// In production use github.com/golang-jwt/jwt/v5 to parse and validate.
		c.Set("user_id_from_token", tokenStr) // placeholder: handlers should parse token and extract user id
		c.Next()
	}
}
