package utils

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gleb-korostelev/gophermart.git/internal/config"
	"github.com/gleb-korostelev/gophermart.git/internal/models"
	"github.com/gleb-korostelev/gophermart.git/tools/logger"
	"github.com/golang-jwt/jwt/v4"
)

func GenerateJWT(login string, jwtKeySecret string) (string, error) {
	expirationTime := time.Now().Add(config.TokenExpiration)
	claims := &models.Claims{
		Login: login,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(jwtKeySecret))

	return tokenString, err
}

func VerifyJWT(tokenString string, jwtKeySecret string) (*models.Claims, error) {
	claims := &models.Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return []byte(jwtKeySecret), nil
	})
	if err != nil {
		return nil, err
	}
	if !token.Valid {
		return nil, config.ErrTokenInvalid
	}

	return claims, nil
}

func SetJWTInCookie(w http.ResponseWriter, login string) {
	tokenString, err := GenerateJWT(login, config.JwtKeySecret)
	if err != nil {
		http.Error(w, "Failed to generate token", http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:  "auth",
		Value: tokenString,
	})
}

func GetLoginFromCookie(r *http.Request) (string, error) {
	cookie, err := r.Cookie("auth")
	if err != nil {
		return "", err
	}
	claims, err := VerifyJWT(cookie.Value, config.JwtKeySecret)
	if err != nil {
		return "", err
	}
	logger.Infof("Success authorized with login: %s", claims.Login)
	return claims.Login, nil
}
