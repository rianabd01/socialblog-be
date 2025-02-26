package utils

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/rianabd01/socialblog-be/internal/models"
)

var jwtKey = []byte(os.Getenv("JWT_SECRET_KEY"))

type JWTClaim struct {
	UserID    uint   `json:"user_id"`
	Username  string `json:"username"`
	Email     string `json:"email"`
	Name      string `json:"name"`
	AvatarUrl string `json:"avatar_url"`
	Provider  string `json:"provider"`
	jwt.RegisteredClaims
}

func GenerateJWT(user models.User, source string) (string, error) {
	claims := JWTClaim{
		UserID:    user.ID,
		Username:  user.Username,
		Email:     MaskEmail(user.Email),
		Name:      user.Name,
		AvatarUrl: user.AvatarUrl,
		Provider:  source,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(30 * 24 * time.Hour)), // 30 hari
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	fmt.Println("claims:", claims)
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
