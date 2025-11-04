package controller

import (
	"net/http"

	"gopastebin/internal/service"
	"gopastebin/pkg/response"

	"github.com/gin-gonic/gin"
)

type AuthController struct {
	srv *service.AuthService
}

func NewAuthController(s *service.AuthService) *AuthController {
	return &AuthController{srv: s}
}

func (a *AuthController) Register(c *gin.Context) {
	var req service.RegisterDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "invalid_body", err.Error())
		return
	}
	user, err := a.srv.Register(req)
	if err != nil {
		response.FromError(c, err)
		return
	}
	response.OK(c, gin.H{"user": user})
}

func (a *AuthController) Login(c *gin.Context) {
	var req service.LoginDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "invalid_body", err.Error())
		return
	}
	access, refresh, user, err := a.srv.Login(req)
	if err != nil {
		if err.Error() == "email_not_verified" {
			response.BadRequest(c, "email_not_verified", "please verify your email")
			return
		}
		response.FromError(c, err)
		return
	}

	c.SetCookie("access_token", access, 600, "/", "", false, true)
	c.SetCookie("refresh_token", refresh, 60*60*24*30, "/", "", false, true)

	response.OK(c, gin.H{"user": user})
}

func (a *AuthController) Me(c *gin.Context) {
	access, err := c.Cookie("access_token")
	if err != nil || access == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": gin.H{"code": "unauthorized", "message": "no token"}})
		return
	}
	user, err := a.srv.VerifyAccess(access)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": gin.H{"code": "unauthorized", "message": "invalid token"}})
		return
	}
	response.OK(c, gin.H{"user": user})
}

func (a *AuthController) Verify(c *gin.Context) {
	email := c.Query("email")
	code := c.Query("code")
	if email == "" || code == "" {
		response.BadRequest(c, "invalid_params", "email and code required")
		return
	}
	if err := a.srv.VerifyEmail(c, email, code); err != nil {
		response.FromError(c, err)
		return
	}
	response.OK(c, gin.H{"status": "verified"})
}
