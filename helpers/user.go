package helpers

import (
	"fmt"
	"math/rand"

	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) string {
	var passwordBytes = []byte(password)

	// Hash password with Bcrypt MinCost
	hashedPassword, err := bcrypt.GenerateFromPassword(passwordBytes, bcrypt.MinCost)
	if err != nil {
		return ""
	}

	// Cast the hashedPassword to string
	return string(hashedPassword)
}

func CreateSecret() (string, error) {
	secret := make([]byte, 32)
	if _, err := rand.Read(secret); err != nil {
		return "", err
	}

	return fmt.Sprintf("%x", secret), nil
}

func IsPasswordMatched(currentPassword, password string) bool {
	err := bcrypt.CompareHashAndPassword(
		[]byte(currentPassword),
		[]byte(password),
	)
	return err == nil
}
