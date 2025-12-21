package utils

import (
	"errors"
	"time"

	"github.com/armymp/event-booking-api/config"
	"github.com/golang-jwt/jwt/v5"
)

const tokenExpiryHours = 2

type AuthClaims struct {
	Email string `json:"email"`
	UserID int64 `json:"userId"`
	jwt.RegisteredClaims
}

func GenerateToken(email string, userId int64) (string, error) {
	secret := config.AppConfig.JWT.Secret
	if secret == "" {
		return "", errors.New("JWT secret not configured.")
	}
	claims := AuthClaims {
		Email: email,
		UserID: userId,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(
				time.Now().Add(time.Hour * tokenExpiryHours),
			),
			IssuedAt: jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

func VerifyToken(token string) error{
	secret := config.AppConfig.JWT.Secret
	parsedToken, err := jwt.Parse(token, func(t *jwt.Token) (any, error) {
		_, ok := t.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, errors.New("Unexpected signing method.")
		}
		return []byte(secret), nil
	})
	if err != nil {
		return errors.New("Could not parse token.")
	}

	isValidToken := parsedToken.Valid
	if !isValidToken {
		return errors.New("Invalid token")
	}

	// uncomment later when we need to extract email, userId
	// claims, ok := parsedToken.Claims.(jwt.MapClaims)
	// if !ok {
	// 	return errors.New("Invalid token claims.")
	// }
	// email := claims["email"].(string)
	// userId := claims["userId"].(int64)

	return nil
}