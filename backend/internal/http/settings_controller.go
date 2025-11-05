package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gopastebin/internal/service"
)

type SettingsController struct {
	srv *service.AppSettingsService
}

func NewSettingsController(s *service.AppSettingsService) *SettingsController {
	return &SettingsController{srv: s}
}

func (h *SettingsController) GetStorage(c *gin.Context) {
	cur := h.srv.GetStorageProvider()
	c.JSON(http.StatusOK, gin.H{"storage": cur})
}

type setStorageReq struct {
	Storage string `json:"storage"`
}

func (h *SettingsController) SetStorage(c *gin.Context) {
	var req setStorageReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid_body"})
		return
	}
	if err := h.srv.SetStorageProvider(req.Storage); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"ok": true})
}
