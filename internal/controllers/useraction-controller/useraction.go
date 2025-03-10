package useractioncontroller

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rianabd01/socialblog-be/internal/models"
	"github.com/rianabd01/socialblog-be/internal/server"
	"gorm.io/gorm"
)

// Inisialisasi Redis client
var ctx = context.Background()
var redisClient = server.RedisClient

func LikeContent(c *gin.Context) {
	var user models.User
	userID, exists := c.Get("user_id") // Mengambil user_id dari middleware

	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	// Cek apakah user ada di database
	if err := server.DB.Where("id = ?", userID).First(&user).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	type LikeRequest struct {
		PostID *uint `json:"post_id"`
		BlogID *uint `json:"blog_id"`
	}

	var input LikeRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if input.PostID == nil && input.BlogID == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "post_id atau blog_id harus diisi"})
		return
	}

	var like models.Like
	var countField string

	if input.PostID != nil {
		like = models.Like{UserID: &user.ID, PostID: input.PostID}
		countField = "like_count"
	} else {
		like = models.Like{UserID: &user.ID, BlogID: input.BlogID}
		countField = "like_count"
	}

	// Cek apakah user sudah pernah like
	var existingLike models.Like
	if err := server.DB.Where("user_id = ? AND (post_id = ? OR blog_id = ?)", user.ID, input.PostID, input.BlogID).First(&existingLike).Error; err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Already liked"})
		return
	}

	// Insert like ke database
	if err := server.DB.Create(&like).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to like"})
		return
	}

	// Update jumlah like pada tabel yang sesuai
	if input.PostID != nil {
		if err := server.DB.Model(&models.Post{}).Where("id = ?", *input.PostID).Update(countField, gorm.Expr(countField+" + ?", 1)).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update like count"})
			return
		}
		// redisClient.SAdd(ctx, "user:"+strconv.Itoa(int(user.ID))+":liked_posts", *input.PostID) // Cache di Redis
	} else {
		if err := server.DB.Model(&models.Blog{}).Where("id = ?", *input.BlogID).Update(countField, gorm.Expr(countField+" + ?", 1)).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update like count"})
			return
		}
		// redisClient.SAdd(ctx, "user:"+strconv.Itoa(int(user.ID))+":liked_blogs", *input.BlogID) // Cache di Redis
	}

	c.JSON(http.StatusOK, gin.H{"message": "Liked successfully"})
}

func UnlikeContent(c *gin.Context) {
	var user models.User
	userID, exists := c.Get("user_id") // Mengambil user_id dari middleware

	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	// Cek apakah user ada di database
	if err := server.DB.Where("id = ?", userID).First(&user).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	type UnlikeRequest struct {
		PostID *uint `json:"post_id"`
		BlogID *uint `json:"blog_id"`
	}

	var input UnlikeRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if input.PostID == nil && input.BlogID == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "post_id atau blog_id harus diisi"})
		return
	}

	// Cek apakah user sudah melakukan like sebelumnya
	var existingLike models.Like
	if err := server.DB.Where("user_id = ? AND (post_id = ? OR blog_id = ?)", user.ID, input.PostID, input.BlogID).First(&existingLike).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Like not found"})
		return
	}

	// Hapus like dari database
	if err := server.DB.Delete(&existingLike).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to unlike"})
		return
	}

	// Update jumlah like pada tabel yang sesuai
	var countField string
	if input.PostID != nil {
		countField = "like_count"
		if err := server.DB.Model(&models.Post{}).Where("id = ?", *input.PostID).Update(countField, gorm.Expr(countField+" - ?", 1)).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update like count"})
			return
		}
		// redisClient.SRem(ctx, "user:"+strconv.Itoa(int(user.ID))+":liked_posts", *input.PostID) // Hapus dari cache Redis
	} else {
		countField = "like_count"
		if err := server.DB.Model(&models.Blog{}).Where("id = ?", *input.BlogID).Update(countField, gorm.Expr(countField+" - ?", 1)).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update like count"})
			return
		}
		// redisClient.SRem(ctx, "user:"+strconv.Itoa(int(user.ID))+":liked_blogs", *input.BlogID) // Hapus dari cache Redis
	}

	c.JSON(http.StatusOK, gin.H{"message": "Unliked successfully"})
}
