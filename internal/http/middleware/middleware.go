package middleware

import (
	"net/http"
	"tender-backend/internal/http/token"

	"github.com/gin-gonic/gin"
)

func JWTMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenStr := c.Request.Header.Get("Authorization")

		if tokenStr == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "Missing token"})
			c.Abort()
			return
		}

		// remove Bearer prefix
		tokenStr = tokenStr[7:]

		claims, err := token.VerifyJWT(tokenStr)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			c.Abort()
			return
		}

		c.Set("user_id", claims.UserID)
		c.Set("role", claims.Role)
		c.Next()
	}
}

func ClientMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		role := c.GetString("role")
		if role != "client" {
			c.JSON(http.StatusForbidden, gin.H{"error": "Only clients can access this resource"})
			c.Abort()
			return
		}
		c.Next()
	}
}

func ContractorMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		role := c.GetString("role")
		if role != "contractor" {
			c.JSON(http.StatusForbidden, gin.H{"error": "Only contractors can access this resource"})
			c.Abort()
			return
		}
		c.Next()
	}
}
