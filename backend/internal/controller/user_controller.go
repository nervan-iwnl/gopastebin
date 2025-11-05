package controller

import (
	"gopastebin/internal/service"
	"gopastebin/pkg/response"

	"github.com/gin-gonic/gin"
)

type UserController struct {
	srv *service.UserService
}

func NewUserController(s *service.UserService) *UserController {
	return &UserController{srv: s}
}

func (u *UserController) PublicProfile(c *gin.Context) {
	username := c.Param("username")
	profile, err := u.srv.GetPublicProfile(c, username)
	if err != nil {
		response.FromError(c, err)
		return
	}
	response.OK(c, profile)
}

func (u *UserController) Me(c *gin.Context) {
	user := c.MustGet("user")
	me, err := u.srv.GetMe(c, user)
	if err != nil {
		response.FromError(c, err)
		return
	}
	response.OK(c, gin.H{"user": me})
}

func (u *UserController) UpdateMe(c *gin.Context) {
	user := c.MustGet("user")
	var req service.UpdateUserDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "invalid_body", err.Error())
		return
	}
	me, err := u.srv.UpdateMe(c, user, req)
	if err != nil {
		response.FromError(c, err)
		return
	}
	response.OK(c, gin.H{"user": me})
}
