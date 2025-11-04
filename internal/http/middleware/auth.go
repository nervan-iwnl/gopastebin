package middleware

import (
	"net/http"
	"time"

	"gopastebin/internal/service"

	"github.com/gin-gonic/gin"
)

func Auth(authSrv *service.AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		access, err := c.Cookie("access_token")
		if err != nil || access == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": gin.H{"code": "unauthorized", "message": "no token"}})
			c.Abort()
			return
		}

		user, err := authSrv.VerifyAccess(access)
		if err != nil {
			refresh, err2 := c.Cookie("refresh_token")
			if err2 != nil || refresh == "" {
				c.JSON(http.StatusUnauthorized, gin.H{"error": gin.H{"code": "unauthorized", "message": "invalid token"}})
				c.Abort()
				return
			}
			newAccess, newRefresh, user2, err3 := authSrv.RefreshTokens(refresh)
			if err3 != nil {
				c.JSON(http.StatusUnauthorized, gin.H{"error": gin.H{"code": "unauthorized", "message": "invalid refresh"}})
				c.Abort()
				return
			}
			c.SetCookie("access_token", newAccess, int((10 * time.Minute).Seconds()), "/", "", false, true)
			c.SetCookie("refresh_token", newRefresh, int((30*24*time.Hour).Seconds()), "/", "", false, true)
			c.Set("user", user2)
			c.Next()
			return
		}

		c.Set("user", user)
		c.Set("is_admin", user.IsAdmin) // ← вот это
		c.Next()
	}
}
