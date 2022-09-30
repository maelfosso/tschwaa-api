package storage

import (
	"context"
	"crypto/rand"
	"fmt"

	"github.com/fatih/structs"
	"tschwaa.com/api/model"
	"tschwaa.com/api/services"
)

func createSecret() (string, error) {
	secret := make([]byte, 32)
	if _, err := rand.Read(secret); err != nil {
		return "", err
	}

	return fmt.Sprintf("%x", secret), nil
}

func (d *Database) Signup(ctx context.Context, user model.User) (string, error) {
	// Hash the password
	if err := user.HashPassword(); err != nil {
		return "", err
	}

	// Check if a user with the same email exist
	existingUser, err := d.FindUserByUsername(ctx, user.Phone, user.Email)
	if err != nil || existingUser != nil {
		return "", fmt.Errorf("user with the email/phone already exists")
	}

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

func (d *Database) Signin(ctx context.Context, credentials model.SignInCredentials) (*model.SignInResult, error) {
	// Check if the user exists
	var user = model.User{
		Password: credentials.Password,
	}
	var hashedPassword string

	query := `
		SELECT id, firstname, lastname, phone, email, password
		FROM users
		WHERE (phone = $1) OR (email = $1)
	`
	err := d.DB.QueryRowContext(ctx, query, credentials.Username).Scan(&user.ID, &user.Firstname, &user.Lastname, &user.Phone, &user.Email, &hashedPassword)
	if err != nil {
		return nil, fmt.Errorf("user with that username does not exist")
	}

	if !user.IsPasswordMatched(hashedPassword) {
		return nil, fmt.Errorf("the password is not correct")
	}

	user.Password = ""

	var signInResult model.SignInResult
	signInResult.Name = fmt.Sprintf("%s %s", user.Firstname, user.Lastname)
	signInResult.Email = user.Email

	_, tokenString, _ := services.TokenAuth.Encode(structs.Map(&signInResult))
	signInResult.Token = tokenString
	return &signInResult, nil
}
