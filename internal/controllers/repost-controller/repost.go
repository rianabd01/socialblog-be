package repostcontroller

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/rianabd01/socialblog-be/internal/models"
	"github.com/rianabd01/socialblog-be/internal/server"
	"gorm.io/gorm"
)

type RepostResponse struct {
	ID        uint         `json:"id"`
	User      UserResponse `json:"user"`
	Blog      BlogResponse `json:"blog"`
	Quote     string       `json:"quote"`
	LikeCount int          `json:"like_count"`
}

type UserResponse struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

type BlogResponse struct {
	ID    uint   `json:"id"`
	Title string `json:"title"`
	Body  string `json:"body"`
}

func Index(c *gin.Context) {
	var reposts []models.Repost
	server.DB.Preload("Owner").Preload("Blog").Find(&reposts)

	// Konversi ke response struct agar data lebih bersih
	var responses []RepostResponse
	for _, repost := range reposts {
		responses = append(responses, RepostResponse{
			ID:        repost.ID,
			User:      UserResponse{ID: repost.Owner.ID, Name: repost.Owner.Name},
			Blog:      BlogResponse{ID: repost.Blog.ID, Title: repost.Blog.Title, Body: repost.Blog.Body},
			Quote:     repost.Quote,
			LikeCount: repost.LikeCount,
		})
	}

	c.JSON(http.StatusOK, gin.H{"reposts": responses})
}

func ShowDetail(c *gin.Context) {
	var repost models.Repost

	id := c.Param("id")

	if err := server.DB.First(&repost, id).Error; err != nil {
		switch err {
		case gorm.ErrRecordNotFound:
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"messsage": "Data tidak ditemukan"})
			return
		default:
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			return
		}

	}

	c.JSON(http.StatusOK, gin.H{"repost": repost})
}

func Create(c *gin.Context) {
	var repost models.Repost
	var user models.User
	userID, _ := c.Get("user_id") // mendapat username dari middleware
	fmt.Println("userID", userID)

	if result := server.DB.Where("id = ?", userID).First(&user); result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}

	if err := c.ShouldBindJSON(&repost); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	// menetapkan OwnerID dari pengguna yang terotentikasi
	repost.UserID = user.ID

	if err := server.DB.Create(&repost).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to create blog"})
		return
	}

	// preload owner untuk menampilkan informasi pengguna di respons
	if err := server.DB.Preload("Blog").First(&repost, repost.ID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to preload blog data"})
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
