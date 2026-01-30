package middleware

import (
	"time"

	"github.com/AniketDubey199/online_shopping/model"
	"github.com/golang-jwt/jwt/v4"
)

func GenerateToken(user *model.User) (string, error) {
	secret := []byte("super-secret-key")
	methods := jwt.SigningMethodHS256
	claims := jwt.MapClaims{
		"userID":   user.ID.Hex(),
		"username": user.Username,
		"exp":      time.Now().Add(24 * time.Hour).Unix(),
	}
	token, err := jwt.NewWithClaims(methods, claims).SignedString(secret)
	if err != nil {
		return "", err
	}
	return token, err
}
