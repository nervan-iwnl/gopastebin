package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func Admin() gin.HandlerFunc {
	return func(c *gin.Context) {
		isAdmin, ok := c.Get("is_admin")
		if !ok || isAdmin != true {
			c.JSON(http.StatusForbidden, gin.H{
				"error": gin.H{
					"code":    "forbidden",
					"message": "admin only",
				},
			})
			c.Abort()
			return
		}
		c.Next()
	}
}
