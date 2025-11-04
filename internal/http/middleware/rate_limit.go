package middleware

import (
	"net"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

var (
	visitors = make(map[string]*rate.Limiter)
	mu       sync.Mutex
)

func getVisitor(ip string, rps int) *rate.Limiter {
	mu.Lock()
	defer mu.Unlock()
	lim, ok := visitors[ip]
	if !ok {
		lim = rate.NewLimiter(rate.Every(time.Minute/time.Duration(rps)), rps)
		visitors[ip] = lim
	}
	return lim
}

func RateLimit(rps int) gin.HandlerFunc {
	return func(c *gin.Context) {
		ip, _, err := net.SplitHostPort(c.Request.RemoteAddr)
		if err != nil {
			ip = c.ClientIP()
		}
		lim := getVisitor(ip, rps)
		if !lim.Allow() {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"error": gin.H{
					"code":    "rate_limit",
					"message": "too many requests",
				},
			})
			c.Abort()
			return
		}
		c.Next()
	}
}
