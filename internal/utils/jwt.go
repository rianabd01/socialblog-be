package utils

import (
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var jwtKey = []byte(os.Getenv("JWT_SECRET_KEY"))

type JWTClaim struct {
	UserID   uint   `json:"user_id"`
	Username string `json:"username"`
	jwt.RegisteredClaims
}

func GenerateJWT(userID uint, username, source string) (string, error) {
	claims := jwt.MapClaims{
		"user_id":  userID,
		"username": username,
		"source":   source,
		"exp":      time.Now().Add(30 * 24 * time.Hour).Unix(), // kedaluarsa 30 hari
		"iat":      time.Now().Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	fmt.Println("klem", claims)
	return token.SignedString(jwtKey)
}

func ValidateToken(signedToken string) (claims *JWTClaim, err error) {
	token, err := jwt.ParseWithClaims(
		signedToken,
		&JWTClaim{},
		func(token *jwt.Token) (interface{}, error) {
			return []byte(jwtKey), nil
		},
	)
	if err != nil {
		return
	}
	claims, ok := token.Claims.(*JWTClaim)
	if !ok {
		err = errors.New("couldn't parse claims")
		return
	}
	if claims.ExpiresAt.Before(time.Now()) {
		err = errors.New("token expired")
		return
	}
	return claims, nil
}

func MaskEmail(email string) string {
	if email == "" {
		return ""
	}

	// Pisahkan email menjadi bagian sebelum dan sesudah '@'
	parts := strings.Split(email, "@")
	if len(parts) != 2 {
		return email // Kembali ke email asli jika format salah
	}

	username := parts[0]
	domain := parts[1]

	// Mask username: tampilkan 3 karakter awal, sisanya jadi '*'
	maskedUsername := ""
	if len(username) <= 3 {
		maskedUsername = username
	} else {
		maskedUsername = username[:3] + strings.Repeat("*", len(username)-3)
	}

	// Mask domain: sembunyikan semua kecuali bagian setelah titik terakhir
	domainParts := strings.Split(domain, ".")
	if len(domainParts) < 2 {
		return maskedUsername + "@" + domain
	}
	maskedDomain := strings.Repeat("*", len(domain)-len(domainParts[len(domainParts)-1])-1) + "." + domainParts[len(domainParts)-1]

	return maskedUsername + "@" + maskedDomain
}
