package utils

import (
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var jwtKey = []byte(os.Getenv("JWT_SECRET_KEY"))

type JWTClaim struct {
	Username string `json:"username"`
	jwt.RegisteredClaims
}

func GenerateJWT(userID uint, email, source string) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"email":   email,
		"source":  source,
		"exp":     time.Now().Add(30 * 24 * time.Hour).Unix(), // kedaluarsa 30 hari
		"iat":     time.Now().Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
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
