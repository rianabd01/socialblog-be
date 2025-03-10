package blogcontroller

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/rianabd01/socialblog-be/internal/models"
	"github.com/rianabd01/socialblog-be/internal/server"
	"gorm.io/gorm"
)

func Index(c *gin.Context) {
	var blogs []models.Blog

	server.DB.Find(&blogs)

	c.JSON(http.StatusOK, gin.H{"blogs": blogs})
}

func ShowDetail(c *gin.Context) {
	var blog models.Blog

	id := c.Param("id")

	if err := server.DB.First(&blog, id).Error; err != nil {
		switch err {
		case gorm.ErrRecordNotFound:
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"messsage": "Data tidak ditemukan"})
			return
		default:
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			return
		}

	}

	c.JSON(http.StatusOK, gin.H{"blog": blog})
}

func Create(c *gin.Context) {
	var blog models.Blog
	var user models.User
	userID, _ := c.Get("user_id") // mendapat username dari middleware

	fmt.Println("userID", userID)

	if result := server.DB.Where("id = ?", userID).First(&user); result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}

	if err := c.ShouldBindJSON(&blog); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	// menetapkan OwnerID dari pengguna yang terotentikasi
	blog.UserID = user.ID

	if err := server.DB.Create(&blog).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to create blog"})
		return
	}

	// preload owner untuk menampilkan informasi pengguna di respons
	if err := server.DB.Preload("Owner").First(&blog, blog.ID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to preload owner data"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "blog has been created"})
}

func Update(c *gin.Context) {
	var blog models.Blog

	id := c.Param("id")

	if err := c.ShouldBindJSON(&blog); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	if server.DB.Model(&blog).Where("id = ?", id).Updates(&blog).RowsAffected == 0 {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "tidak ada data yang berubah"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Data berhasil di update"})
}

func Delete(c *gin.Context) {
	var blog models.Blog

	input := map[string]string{"id": "0"}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	id, _ := strconv.ParseInt(input["id"], 10, 64)

	if server.DB.Delete(&blog, id).RowsAffected == 0 {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "Tidak dapat menghapus data"})
	}

	c.JSON(http.StatusOK, gin.H{"message": "Data berhasil dihapus"})
}
