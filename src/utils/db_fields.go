package utils

import (
	"fmt"
	"gopastebin/src/db"
	"net/http"
	"regexp"

	"github.com/gin-gonic/gin"
)

func IsFieldUnique(c *gin.Context, value string, field string) bool {
	var err error

	var isUniqueField bool
	switch field {
	case "email":
		isUniqueField, err = db.IsEmailUnique(value)
	case "username":
		isUniqueField, err = db.IsUsernameUnique(value)
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to get user by %s from db", field)})
		return false
	} else if !isUniqueField {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("This %s is already in use", field)})
		return false
	}

	return true
}

func ValidateEmail(email string) bool {
	const emailRegex = `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	re := regexp.MustCompile(emailRegex)
	return re.MatchString(email)
}

func ValidateUsername(username string) bool {
	var validUsername = regexp.MustCompile(`^[a-zA-Z0-9_-]{4,12}$`)
	return validUsername.MatchString(username)
}

func ValidatePassword(password string) bool {
	lengthCheck := regexp.MustCompile(`^.{8,20}$`)
	if !lengthCheck.MatchString(password) {
		return false
	}

	upperCaseCheck := regexp.MustCompile(`[A-Z]`)
	if !upperCaseCheck.MatchString(password) {
		return false
	}

	digitCheck := regexp.MustCompile(`\d`)
	if !digitCheck.MatchString(password) {
		return false
	}

	lowerCaseCheck := regexp.MustCompile(`[a-z]`)

	return lowerCaseCheck.MatchString(password)
}
