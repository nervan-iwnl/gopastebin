package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gopastebin/internal/service"
)

type SettingsController struct {
	settings *service.AppSettingsService
}

func NewSettingsController(s *service.AppSettingsService) *SettingsController {
	return &SettingsController{settings: s}
}

// GET /api/v1/settings/storage
func (h *SettingsController) GetStorage(c *gin.Context) {
	cur := h.settings.GetStorageProvider()
	c.JSON(http.StatusOK, gin.H{"storage": cur})
}

type setStorageReq struct {
	Storage string `json:"storage"`
}

// POST /api/v1/settings/storage
func (h *SettingsController) SetStorage(c *gin.Context) {
	// тут проверка что админ
	isAdmin, _ := c.Get("is_admin")
	if isAdmin != true {
		c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
		return
	}

	var req setStorageReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid_body"})
		return
	}

	if err := h.settings.SetStorageProvider(req.Storage); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"ok": true})
}
