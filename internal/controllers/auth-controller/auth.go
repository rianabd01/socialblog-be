package authcontroller

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"gorm.io/gorm"

	"github.com/rianabd01/socialblog-be/internal/models"
	"github.com/rianabd01/socialblog-be/internal/server"
	"github.com/rianabd01/socialblog-be/internal/utils"
)

var (
	googleOauthConfig *oauth2.Config // declare as pointer
	oauthStateString  = "random"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file:", err) // print error env
	}

	googleOauthConfig = &oauth2.Config{
		RedirectURL:  "http://localhost:8080/auth/google/callback",
		ClientID:     os.Getenv("GOOGLE_ID"),
		ClientSecret: os.Getenv("GOOGLE_SECRET"),
		Scopes:       []string{"profile", "email"},
		Endpoint:     google.Endpoint,
	}

	if googleOauthConfig.ClientID == "" || googleOauthConfig.ClientSecret == "" {
		panic("GOOGLE_CLIENT_ID and GOOGLE_CLIENT_SECRET must be set in .env")
	}
}

func Signup(c *gin.Context) {
	var user models.User

	// Binding JSON
	if err := c.ShouldBindJSON(&user); err != nil {
		log.Printf("Invalid request payload: %v", err) // Log detail error
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid request data. Please check your input.",
		})
		return
	}

	// Validasi user
	if err := user.Validate(); err != nil {
		log.Printf("Validation failed: %v", err) // Log detail error
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	// Hash password
	hashedPassword, err := hashPassword(user.Password)
	if err != nil {
		log.Printf("Failed to hash password: %v", err) // Log error
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Something went wrong. Please try again later.",
		})
		return
	}
	user.Password = hashedPassword

	// Simpan ke database
	if err := server.DB.Create(&user).Error; err != nil {
		log.Printf("Database error: %v", err) // Log detail error
		// Cek apakah error karena duplikat
		if isDuplicateError(err) {
			c.JSON(http.StatusConflict, gin.H{
				"message": "Email or username already exists.",
			})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "Failed to create account. Please try again later.",
			})
		}
		return
	}

	// Response sukses
	c.JSON(http.StatusCreated, gin.H{
		"message": "Signup successful",
	})
}

func isDuplicateError(err error) bool {
	if err != nil {
		return strings.Contains(err.Error(), "duplicate") || strings.Contains(err.Error(), "unique")
	}
	return false
}

func Login(c *gin.Context) {
	var credentials struct {
		Username string `json:"username" validate:"required"`
		Password string `json:"password" validate:"required"`
	}

	if err := c.ShouldBindJSON(&credentials); err != nil {
		log.Printf("Invalid request payload: %v", err) // Log detail error
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid request data. Please check your input.",
		})
		return
	}

	var user models.User
	if err := server.DB.Where("username = ?", credentials.Username).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Invalid username or password"})
		return
	}

	if !checkPasswordHash(credentials.Password, user.Password) {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Invalid username or password"})
		return
	}

	jwtToken, err := utils.GenerateJWT(user, "manual")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "token invalid"})
		return
	}

	user.LastLogin = &gorm.DeletedAt{Time: time.Now()}
	server.DB.Save(&user)

	c.JSON(http.StatusOK, gin.H{
		"message": "login successful",
		"token":   jwtToken,
	})
}

func GoogleLogin(c *gin.Context) {
	url := googleOauthConfig.AuthCodeURL(oauthStateString,
		oauth2.AccessTypeOffline,
		oauth2.ApprovalForce)

	c.Redirect(http.StatusTemporaryRedirect, url)
}

func GoogleCallback(c *gin.Context) {
	code := c.Query("code")
	if code == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Code not found"})
		return
	}

	token, err := googleOauthConfig.Exchange(c, code)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to exchange token", "error": err.Error()})
		return
	}

	resp, err := http.Get(fmt.Sprintf("https://www.googleapis.com/oauth2/v2/userinfo?access_token=%s", token.AccessToken))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to fetch user info", "error": err.Error()})
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to fetch user info", "status": resp.Status})
		return
	}

	var userInfo struct {
		ID      string `json:"id"`
		Email   string `json:"email"`
		Name    string `json:"name"`
		Picture string `json:"picture"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&userInfo); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to decode user info", "error": err.Error()})
		return
	}

	// Simpan atau update user, tambahkan refresh token jika ada
	var user models.User
	if err := server.DB.Where("provider_id = ? OR email = ?", userInfo.ID, userInfo.Email).First(&user).Error; err != nil {
		user = models.User{
			ProviderID:   userInfo.ID,
			Email:        userInfo.Email,
			Name:         userInfo.Name,
			AvatarUrl:    userInfo.Picture,
			Provider:     "google",
			RefreshToken: token.RefreshToken,
			Verified:     true,
			Username:     userInfo.Email,
		}
		if err := server.DB.Create(&user).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to create user", "error": err.Error()})
			return
		}
	} else if token.RefreshToken != "" {
		server.DB.Model(&user).Update("refresh_token", token.RefreshToken)
	}

	// Generate JWT
	jwtToken, err := utils.GenerateJWT(user, "google")

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to generate JWT", "error": err.Error()})
		return
	}

	redirect := os.Getenv("AUTH_REDIRECT") + "?token=" + jwtToken

	c.Redirect(http.StatusFound, redirect)
}

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func checkPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
