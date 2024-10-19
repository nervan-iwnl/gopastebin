package utils

import (
	"errors"
	"gopastebin/src/db"
	"gopastebin/src/models"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
)

var jwtKey = []byte(os.Getenv("JWT_SECRET"))
var refreshTokenKey = []byte(os.Getenv("JWT_REFRESH_SECRET"))

func GenerateJWT(user *models.User) (accessToken string, refreshToken string, err error) {
	accessClaims := &jwt.MapClaims{
		"username": user.Username,
		"exp":      time.Now().Add(time.Minute * 10).Unix(), // 10m
	}
	accessTokenObj := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	accessToken, err = accessTokenObj.SignedString(jwtKey)
	if err != nil {
		return "", "", err
	}

	refreshClaims := &jwt.MapClaims{
		"username": user.Username,
		"exp":      time.Now().Add(time.Hour * 24 * 30).Unix(), // 30d
	}
	refreshTokenObj := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
	refreshToken, err = refreshTokenObj.SignedString(refreshTokenKey)
	if err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}

func VerifyJWT(tokenString string, isRefresh bool) (*models.User, error) {
	var key []byte
	if isRefresh {
		key = refreshTokenKey
	} else {
		key = jwtKey
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return key, nil
	})

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token")
	}

	user, err := db.GetUserByUsername(claims["username"].(string))
	if err != nil {
		return nil, err
	}
	return user, nil
}

func RefreshTokens(refreshTokenString string) (newAccessToken string, newRefreshToken string, err error) {
	user, err := VerifyJWT(refreshTokenString, true)
	if err != nil {
		return "", "", errors.New("invalid or expired refresh token")
	}

	newAccessToken, newRefreshToken, err = GenerateJWT(user)
	if err != nil {
		return "", "", errors.New("failed to generate new tokens")
	}

	return newAccessToken, newRefreshToken, nil
}
