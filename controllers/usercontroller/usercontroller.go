package usercontroller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rianabd01/socialblog-be/models"
)

func Create(c *gin.Context) {
	var post models.User

	if err := c.ShouldBindJSON(&post); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	if err := models.DB.Create(&post).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to create post"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"post": post})
}
