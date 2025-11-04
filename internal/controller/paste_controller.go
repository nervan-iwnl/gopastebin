package controller

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gopastebin/internal/service"
)

const maxPasteSize = 200_000 // 200KB, подгони если надо

type PasteController struct {
	srv *service.PasteService
}

func NewPasteController(s *service.PasteService) *PasteController {
	return &PasteController{srv: s}
}

// POST /api/v1/pastes
func (p *PasteController) Create(c *gin.Context) {
	var dto service.CreatePasteDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": gin.H{"code": "invalid_body", "message": err.Error()}})
		return
	}
	if len(dto.Content) > maxPasteSize {
		c.JSON(http.StatusBadRequest, gin.H{"error": gin.H{"code": "too_large", "message": "content too big"}})
		return
	}

	user := c.MustGet("user")
	paste, err := p.srv.CreateUserPaste(c, user, dto)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": gin.H{"code": "internal_error", "message": err.Error()}})
		return
	}
	c.JSON(http.StatusOK, gin.H{"paste": paste})
}

// POST /api/v1/pastes/anon
func (p *PasteController) CreateAnon(c *gin.Context) {
	var dto service.CreatePasteDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": gin.H{"code": "invalid_body", "message": err.Error()}})
		return
	}
	if len(dto.Content) > maxPasteSize {
		c.JSON(http.StatusBadRequest, gin.H{"error": gin.H{"code": "too_large", "message": "content too big"}})
		return
	}

	paste, err := p.srv.CreateAnonPaste(c, dto)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": gin.H{"code": "internal_error", "message": err.Error()}})
		return
	}
	c.JSON(http.StatusOK, gin.H{"paste": paste})
}

// GET /api/v1/pastes/:slug
func (p *PasteController) Get(c *gin.Context) {
	slug := c.Param("slug")
	paste, content, err := p.srv.GetPasteWithContent(c, slug)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": gin.H{"code": "not_found", "message": "paste not found"}})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"paste":   paste,
		"content": content,
	})
}

// GET /api/v1/pastes/:slug/raw
func (p *PasteController) Raw(c *gin.Context) {
	slug := c.Param("slug")
	_, content, err := p.srv.GetPasteWithContent(c, slug)
	if err != nil {
		c.String(http.StatusNotFound, "not found")
		return
	}
	c.String(http.StatusOK, content)
}

// GET /api/v1/pastes/recent?limit=50&offset=0
func (p *PasteController) Recent(c *gin.Context) {
	limit, offset := parseLimitOffset(c, 50, 0)
	pastes, err := p.srv.GetRecentPublicPaged(c, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": gin.H{"code": "internal_error", "message": err.Error()}})
		return
	}
	c.JSON(http.StatusOK, gin.H{"pastes": pastes})
}

// DELETE /api/v1/pastes/:slug
func (p *PasteController) Delete(c *gin.Context) {
	slug := c.Param("slug")
	user := c.MustGet("user")
	if err := p.srv.DeletePaste(c, user, slug); err != nil {
		if err.Error() == "forbidden" {
			c.JSON(http.StatusForbidden, gin.H{"error": gin.H{"code": "forbidden", "message": "not owner"}})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": gin.H{"code": "internal_error", "message": err.Error()}})
		return
	}
	c.Status(http.StatusNoContent)
}

// PUT /api/v1/pastes/:slug
func (p *PasteController) Update(c *gin.Context) {
	slug := c.Param("slug")
	var dto service.UpdatePasteDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": gin.H{"code": "invalid_body", "message": err.Error()}})
		return
	}
	if dto.Content != nil && len(*dto.Content) > maxPasteSize {
		c.JSON(http.StatusBadRequest, gin.H{"error": gin.H{"code": "too_large", "message": "content too big"}})
		return
	}
	user := c.MustGet("user")
	paste, err := p.srv.UpdatePaste(c, user, slug, dto)
	if err != nil {
		if err.Error() == "forbidden" {
			c.JSON(http.StatusForbidden, gin.H{"error": gin.H{"code": "forbidden", "message": "not owner"}})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": gin.H{"code": "internal_error", "message": err.Error()}})
		return
	}
	c.JSON(http.StatusOK, gin.H{"paste": paste})
}

// GET /api/v1/me/pastes?limit=50&offset=0
func (p *PasteController) MyPastes(c *gin.Context) {
	user := c.MustGet("user")
	limit, offset := parseLimitOffset(c, 50, 0)
	pastes, err := p.srv.GetMyPastesPaged(c, user, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": gin.H{"code": "internal_error", "message": err.Error()}})
		return
	}
	c.JSON(http.StatusOK, gin.H{"pastes": pastes})
}

// GET /api/v1/me/folders
func (p *PasteController) MyFolders(c *gin.Context) {
	user := c.MustGet("user")
	folders, err := p.srv.GetMyFolders(c, user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": gin.H{"code": "internal_error", "message": err.Error()}})
		return
	}
	c.JSON(http.StatusOK, gin.H{"folders": folders})
}

// GET /api/v1/me/pastes/by-folder?folder=...
func (p *PasteController) MyPastesByFolder(c *gin.Context) {
	user := c.MustGet("user")
	folder := c.Query("folder")
	pastes, err := p.srv.GetMyPastesInFolder(c, user, folder)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": gin.H{"code": "internal_error", "message": err.Error()}})
		return
	}
	c.JSON(http.StatusOK, gin.H{"pastes": pastes})
}

func parseLimitOffset(c *gin.Context, defLimit, defOffset int) (int, int) {
	limitStr := c.DefaultQuery("limit", strconv.Itoa(defLimit))
	offsetStr := c.DefaultQuery("offset", strconv.Itoa(defOffset))

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		limit = defLimit
	}
	offset, err := strconv.Atoi(offsetStr)
	if err != nil || offset < 0 {
		offset = defOffset
	}
	return limit, offset
}
