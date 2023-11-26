package app

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func NewAccessToken(userID int) *jwt.Token {
	now := time.Now()
	return jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Issuer:    "chirpy-access",
		IssuedAt:  jwt.NewNumericDate(now.UTC()),
		ExpiresAt: jwt.NewNumericDate(now.Add(1 * time.Hour).UTC()),
		Subject:   fmt.Sprintf("%d", userID),
	})
}

func NewRefreshToken(userID int) *jwt.Token {
	now := time.Now()
	return jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Issuer:    "chirpy-refresh",
		IssuedAt:  jwt.NewNumericDate(now.UTC()),
		ExpiresAt: jwt.NewNumericDate(now.Add(60 * (24 * time.Hour)).UTC()),
		Subject:   fmt.Sprintf("%d", userID),
	})
}
