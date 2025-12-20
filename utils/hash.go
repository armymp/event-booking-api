package utils

import (
	"log/slog"

	"golang.org/x/crypto/bcrypt"
)

const bcryptCost = 14

func HashPassword(password string) (string, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcryptCost)
	if err != nil {
		slog.Error("Failed to hash password", "error", err)
		return "", err
	}

	return string(hashed), err
}

func CheckPasssword(hashedPassword, plainPassword string) bool {
	err := bcrypt.CompareHashAndPassword(
		[]byte(hashedPassword),
		[]byte(plainPassword),
	)
	return err == nil
}
