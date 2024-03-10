package util

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

// HashPassword computes the bcrypt hash of the password
func HashPassword(password string) (string, error) {
	hassedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("failed to ahs password: %v", err)
	}
	return string(hassedPassword), nil
}

// CheckPassword checks if the provided password os correct
func CheckPassword(password, hashedPassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}
