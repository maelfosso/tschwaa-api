package storage

import (
	"context"
	"crypto/rand"
	"fmt"

	"golang.org/x/crypto/bcrypt"
	"tschwaa.com/api/model"
)

func hashPassword(password string) (string, error) {
	var passwordBytes = []byte(password)

	// Hash password with Bcrypt MinCost
	hashedPassword, err := bcrypt.GenerateFromPassword(passwordBytes, bcrypt.MinCost)

	// Cast the hashedPassword to string
	return string(hashedPassword), err
}

func isPasswordMatched(hashedPassword, currentPassword string) bool {
	err := bcrypt.CompareHashAndPassword(
		[]byte(hashedPassword),
		[]byte(currentPassword),
	)
	return err == nil
}

func createSecret() (string, error) {
	secret := make([]byte, 32)
	if _, err := rand.Read(secret); err != nil {
		return "", err
	}

	return fmt.Sprintf("%x", secret), nil
}

func (d *Database) Signup(ctx context.Context, user model.User) (string, error) {
	// Hash the password
	hashedPassword, err := hashPassword(user.Password)
	if err != nil {
		return "", err
	}
	user.Password = hashedPassword

	// Get the token - Next will have token for email and token for sms
	token, err := createSecret()
	if err != nil {
		return "", err
	}

	// Insert query
	query := `
		INSERT INTO users(firstname, lastname, phone, email, password, token)
		VALUES ($1, $2, $3, $4, $5, $6)
	`
	_, err = d.DB.ExecContext(ctx, query, user.Firstname, user.Lastname, user.Phone, user.Email, user.Password, token)
	return token, err
}
