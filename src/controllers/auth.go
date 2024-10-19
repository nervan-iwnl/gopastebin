package controllers

import (
	"gopastebin/src/db"
	"gopastebin/src/models"
	"gopastebin/src/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Register(c *gin.Context) {
	var user models.User
	err := c.BindJSON(&user)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if user.Username == "" || user.Email == "" || user.Password == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Username, email, and password are required"})
		return
	}

	if !utils.ValidateEmail(user.Email) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid email format"})
		return
	}

	if !utils.ValidateUsername(user.Username) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Username must be between 4 and 13 characters long and can only contain letters, numbers, underscores, and hyphens."})
		return
	}

	if !utils.ValidatePassword(user.Password) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Password must be between 8 and 20 characters long, and must contain at least one uppercase and lowercase letter and one digit."})
		return
	}

	if !utils.IsFieldUnique(c, user.Email, "email") || !utils.IsFieldUnique(c, user.Username, "username") {
		return
	}

	user.Password = utils.HashPassword(user.Password, utils.GenerateSalt())

	err = db.CreateUser(&user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	_, token, err := utils.GenerateJWT(&user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	utils.SendConfirmationEmail(user.Email, token)
	c.JSON(http.StatusOK, gin.H{"message": "Profile created successfully"})
}

func Login(c *gin.Context) {
	var credentials struct {
		Login    string `json:"login"`
		Password string `json:"password"`
	}
	err := c.BindJSON(&credentials)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var user *models.User
	if utils.ValidateEmail(credentials.Login) {
		user, err = db.GetUserByEmail(credentials.Login)
		if err != nil || !utils.CheckPassword(credentials.Password, user.Password) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
			return
		}
	} else if utils.ValidateUsername(credentials.Login) {
		user, err = db.GetUserByUsername(credentials.Login)
		if err != nil || !utils.CheckPassword(credentials.Password, user.Password) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username or password"})
			return
		}
	} else {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username or password"})
		return
	}

	accessToken, refreshToken, err := utils.GenerateJWT(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate tokens"})
		return
	}

	c.SetCookie("access_token", accessToken, 600, "/", "", false, true) // 10m

	c.SetCookie("refresh_token", refreshToken, 3600*24*30, "/", "", false, true) // 30d

	c.JSON(http.StatusOK, gin.H{"message": "Login successful"})
}

func ConfirmEmail(c *gin.Context) {
	token := c.Param("token")
	email, err := utils.VerifyJWT(token, true)
	if err != nil || token == "" {
		c.JSON(http.StatusNotFound, gin.H{"error": "Incorrect token"})
		return
	}

	user, err := db.GetUserByEmail(email.Email)
	if err != nil || user == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	user.Confirmed = true
	err = db.UpdateUser(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to confirm email"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Email confirmed successfully"})
}
