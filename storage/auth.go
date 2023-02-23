package storage

import (
	"context"
	"crypto/rand"
	"fmt"

	"github.com/fatih/structs"
	"go.uber.org/zap"
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
	existingUser, err := d.FindUserByUsername(ctx, credentials.Username, credentials.Username)
	if err != nil || existingUser == nil {
		return nil, fmt.Errorf("user with that username does not exist")
	}

	if existingUser.IsPasswordMatched(credentials.Password) {
		return nil, fmt.Errorf("the password is not correct")
	}

	var signInResult model.SignInResult
	signInResult.Name = fmt.Sprintf("%s %s", existingUser.Firstname, existingUser.Lastname)
	signInResult.Email = existingUser.Email
	// d.log.Info("Sign In Result", zap.Any("sign-i n-r", signInResult))
	// d.log.Info("JWT Secret", zap.String("secret", os.Getenv("JWT_SECRET")))

	// res := jwtauth.New("HS256", []byte(os.Getenv("JWT_SECRET")), nil)
	// d.log.Info("INstances", zap.Any("res", res), zap.Any("tokenauth", services.TokenAuth))

	// _, ts, te := res.Encode((structs.Map(&signInResult)))
	// d.log.Info("JWT GetOS", zap.Any("Os Getenv x2", ts), zap.Error(te))

	// _, rt, re := services.TokenAuth.Encode(map[string]interface{}{"user_id": 123})
	// d.log.Info("JWT GetOS x2", zap.Any("Os Getenv x2", rt), zap.Error(re))

	// d.log.Info("Struct Map", zap.Any("map sir", structs.Map(&signInResult)), zap.Any("jwt", services.TokenAuth))
	// t, tokenString, err := services.TokenAuth.Encode(map[string]interface{}{"user_id": 123}) // (structs.Map(&signInResult))

	// _, tokenString, err := res.Encode(structs.Map(&signInResult)) // (structs.Map(&signInResult))
	tokenString, err := services.GenerateJWTToken(structs.Map(&signInResult))
	if err != nil {
		return nil, err
	}

	d.log.Info("Sign In Token", zap.String("token", tokenString), zap.Any("jwt token", tokenString), zap.Error(err))
	signInResult.Token = tokenString
	return &signInResult, nil
}
