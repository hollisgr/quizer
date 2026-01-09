package middleware

import (
	"github.com/gin-gonic/gin"
)

type TokenManager interface {
	ParseToken(tokenStr string) (int, error)
}

func AuthMiddleware(tm TokenManager) gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenStr := c.GetHeader("Authorization")
		userID, err := tm.ParseToken(tokenStr)
		if err != nil {
			c.AbortWithStatus(401)
			return
		}

		c.Set("user_id", userID)
		c.Next()
	}
}
