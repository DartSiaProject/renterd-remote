package middlewares

import (
	"net/http"
	constants "renterd-remote/constant"
	"renterd-remote/services/auth"

	"github.com/gin-gonic/gin"
)

func JwtAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		err := auth.VerifyToken(auth.ExtractToken(c))
		if err != nil {
			c.IndentedJSON(http.StatusUnauthorized, gin.H{"error": err.Error(), "message": constants.Unauthorized})
			c.Abort()
			return
		}
		c.Next()
	}
}
