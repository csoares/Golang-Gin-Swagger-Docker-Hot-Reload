package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"projetoapi/internal/infrastructure/jwt"
)

// AuthMiddleware creates a new authorization middleware
func AuthMiddleware(jwtService jwt.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		if !jwtService.ValidateToken(c) {
			c.JSON(http.StatusUnauthorized, gin.H{"status": http.StatusUnauthorized, "message": "Not authorized"})
			c.Abort()
			return
		}

		username, err := jwtService.GetUsernameFromToken(c)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"status": http.StatusUnauthorized, "message": "Not authorized"})
			c.Abort()
			return
		}

		c.Set("username", username)
		c.Next()
	}
}
