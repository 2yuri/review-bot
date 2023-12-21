package middlewares

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func ApiKey() gin.HandlerFunc {
	return func(c *gin.Context) {
		key := c.GetHeader("x-api-key")

		if key != "awesomekey" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"message": "Invalid api key",
			})

			return
		}

		c.Next()
	}
}
