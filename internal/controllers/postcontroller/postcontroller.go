package postcontroller

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/rianabd01/socialblog-be/internal/models"
	"github.com/rianabd01/socialblog-be/internal/server"
	"gorm.io/gorm"
)

func Index(c *gin.Context) {
	var posts []models.Post

	server.DB.Find(&posts)

	c.JSON(http.StatusOK, gin.H{"posts": posts})
}

func ShowDetail(c *gin.Context) {
	var post models.Post

	id := c.Param("id")

	if err := server.DB.First(&post, id).Error; err != nil {
		switch err {
		case gorm.ErrRecordNotFound:
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"messsage": "Data tidak ditemukan"})
			return
		default:
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			return
		}

	}

	c.JSON(http.StatusOK, gin.H{"post": post})
}

func Create(c *gin.Context) {
	var post models.Post
	var user models.User
	username, _ := c.Get("username") // mendapat username dari middleware

	if result := server.DB.Where("username = ?", username).First(&user); result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}

	if err := c.ShouldBindJSON(&post); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	// menetapkan OwnerID dari pengguna yang terotentikasi
	post.OwnerID = user.ID

	if err := server.DB.Create(&post).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to create post"})
		return
	}

	// preload owner untuk menampilkan informasi pengguna di respons
	if err := server.DB.Preload("Owner").First(&post, post.ID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to preload owner data"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "post has been created"})
}

func Update(c *gin.Context) {
	var post models.Post

	id := c.Param("id")

	if err := c.ShouldBindJSON(&post); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	if server.DB.Model(&post).Where("id = ?", id).Updates(&post).RowsAffected == 0 {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "tidak ada data yang berubah"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Data berhasil di update"})
}

func Delete(c *gin.Context) {
	var post models.Post

	input := map[string]string{"id": "0"}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	id, _ := strconv.ParseInt(input["id"], 10, 64)

	if server.DB.Delete(&post, id).RowsAffected == 0 {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "Tidak dapat menghapus data"})
	}

	c.JSON(http.StatusOK, gin.H{"message": "Data berhasil dihapus"})
}
