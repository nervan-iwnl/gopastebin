package controllers

import (
	"fmt"
	"net/http"

	"gopastebin/src/db"
	"gopastebin/src/fb"
	"gopastebin/src/models"
	"gopastebin/src/utils"

	"github.com/gin-gonic/gin"
)

// TODO переписать используя firebase
func CreatePaste(c *gin.Context) {
	var paste models.Paste
	err := c.BindJSON(&paste)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if paste.Content == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Empty file"})
		return
	}

	user, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	paste.UserID = user.(*models.User).ID

	if paste.Slug == "" {
		var lnk string
		for {
			lnk = utils.GenerateRandomPasteLink()
			isUnique, err := db.IsPasteUnique(lnk)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "error while checking paste slug uniqueness"})
			}
			if isUnique {
				paste.Slug = lnk
				break
			}
		}
	}
	isUnique, err := db.IsPasteUnique(paste.Slug)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error while checking paste slug uniqueness"})
		return
	}
	if !isUnique {
		c.JSON(http.StatusBadRequest, gin.H{"error": "This slug is used by other user"})
	}

	path, err := fb.UploadFileToFirebase(fmt.Sprintf("%v", paste.UserID), paste.Slug, paste.Content)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to write file"})
		return
	}
	paste.Content = path
	err = db.CreatePaste(&paste)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Paste created successfully", "paste": paste})
}

// переписать используя firebase
func UpdatePaste(c *gin.Context) {
	slug := c.Param("slug")

	var paste models.Paste
	err := db.GetPasteBySlug(slug, &paste)
	if err != nil || paste.ID == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Paste not found"})
		return
	}

	user, exists := c.Get("user")
	if !exists || paste.UserID != user.(*models.User).ID {
		c.JSON(http.StatusForbidden, gin.H{"error": "Forbidden"})
		return
	}

	err = c.BindJSON(&paste)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = db.UpdatePaste(&paste)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update paste"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Paste updated successfully", "paste": paste})
}

func GetPaste(c *gin.Context) {
	slug := c.Param("slug")
	var paste models.Paste
	err := db.GetPasteBySlug(slug, &paste)
	if err != nil || paste.ID == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Paste not found"})
		return
	}
	content, err := fb.GetFileFromFirebase(fmt.Sprintf("%v", paste.UserID), paste.Slug)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error while reading file"})
		return
	}
	paste.Content = content
	c.JSON(http.StatusOK, gin.H{"paste": paste})
}
