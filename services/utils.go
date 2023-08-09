package services

import "golang.org/x/crypto/bcrypt"

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

func IsPasswordMatched(currentPassword, password string) bool {
	err := bcrypt.CompareHashAndPassword(
		[]byte(currentPassword),
		[]byte(password),
	)
	return err == nil
}
