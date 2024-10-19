package middleware

import (
	"gopastebin/src/utils"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func Auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		accessToken, err := c.Cookie("access_token")
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Access token missing"})
			c.Abort()
			return
		}

		refreshToken, err := c.Cookie("refresh_token")
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Refresh token missing"})
			c.Abort()
			return
		}

		user, err := utils.VerifyJWT(accessToken, false)
		if err != nil {
			log.Println("Access token invalid or expired:", err)

			newAccessToken, newRefreshToken, err := utils.RefreshTokens(refreshToken)
			if err != nil {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired refresh token"})
				c.Abort()
				return
			}

			c.SetCookie("access_token", newAccessToken, int((time.Minute * 10).Seconds()), "/", "", false, true)
			c.SetCookie("refresh_token", newRefreshToken, int((time.Hour * 24 * 30).Seconds()), "/", "", false, true)

			user, err = utils.VerifyJWT(newAccessToken, false)
			if err != nil {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid access token after refresh"})
				c.Abort()
				return
			}
		}

		c.Set("user", user)
		log.Println("Authenticated user:", user.Username)
		c.Next()
	}
}
