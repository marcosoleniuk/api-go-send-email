package middleware

import (
	"api-enviar-email-moleniuk/database"
	"github.com/gin-gonic/gin"
	"net/http"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		username, password, hasAuth := c.Request.BasicAuth()
		apiKey := c.GetHeader("X-API-KEY")

		if hasAuth {
			if !database.ValidateUser(username, password) {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Credenciais inválidas"})
				c.Abort()
				return
			}
		} else if apiKey != "" {
			if !database.ValidateAPIKey(apiKey) {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "API Key inválida"})
				c.Abort()
				return
			}
		} else {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Autenticação necessária"})
			c.Abort()
			return
		}

		c.Next()
	}
}
