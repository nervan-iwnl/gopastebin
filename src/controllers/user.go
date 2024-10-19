package controllers

import (
	"net/http"

	"gopastebin/src/db"
	"gopastebin/src/models"

	"github.com/gin-gonic/gin"
)

func GetUserProfile(c *gin.Context) {
	username := c.Param("username")

	userProfile, err := db.GetUserByUsername(username)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	userResponse := models.UserProfile{
		Username:  userProfile.Username,
		Email:     userProfile.Email,
		Avatar:    userProfile.Avatar,
		Confirmed: userProfile.Confirmed,
	}

	c.JSON(http.StatusOK, gin.H{"user_profile": userResponse})
}

func UpdateUserProfile(c *gin.Context) { // TODO catch errs
	var updateUser models.User
	err := c.BindJSON(&updateUser)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	existingUser := user.(*models.User)
	existingUser.Username = updateUser.Username
	existingUser.Avatar = updateUser.Avatar

	err = db.UpdateUser(existingUser)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user profile"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Profile updated successfully", "user": existingUser})
}
