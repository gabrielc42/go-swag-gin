package middleware

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

func hello(auths ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		fmt.Println("Print from middleware")
		c.Next()
	}
}

func respondWithError(c *gin.Context, code int, message interface{}) {
	c.AbortWithStatusJSON(code, gin.H{
		"error":  message,
		"status": false})
}

func tokenAuthMiddleware() gin.HandlerFunc {
	// todo
}

func verifyToken(t string) (jwt.MapClaims, bool) {
	// todo
}
