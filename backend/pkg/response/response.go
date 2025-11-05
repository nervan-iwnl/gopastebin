package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func OK(c *gin.Context, data gin.H) {
	c.JSON(http.StatusOK, data)
}

func BadRequest(c *gin.Context, code, msg string) {
	c.JSON(http.StatusBadRequest, gin.H{
		"error": gin.H{
			"code":    code,
			"message": msg,
		},
	})
}

func FromError(c *gin.Context, err error) {
	c.JSON(http.StatusInternalServerError, gin.H{
		"error": gin.H{
			"code":    "internal_error",
			"message": err.Error(),
		},
	})
}
