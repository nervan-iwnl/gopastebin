package middleware

import (
	"log"
	"time"

	"github.com/gin-gonic/gin"
)

func Logging() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()
		log.Println("Request:", c.Request.Method, c.Request.URL, "Duration:", time.Since(start))
	}
}
