package storage

import (
	"context"
	"crypto/rand"
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/fatih/structs"
	"github.com/maelfosso/jwtauth"
	"tschwaa.com/api/model"
)

var JwtKey string = "schwaa"

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

func (d *Database) Signin(ctx context.Context, credentials model.SignInCredentials) (string, error) {
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
		return "", fmt.Errorf("user with that username does not exist")
	}

	if !user.IsPasswordMatched(hashedPassword) {
		return "", fmt.Errorf("the password is not correct")
	}

	user.Password = ""
	tokenAuth := jwtauth.New("HS512", []byte("schwaa"), nil)
	_, tokenString, _ := tokenAuth.Encode(structs.Map(&user))

	return tokenString, nil
	// claims := jwt.NewWithClaims(jwt.SigningMethodHS512, model.JwtClaims{
	// 	User: user,
	// 	StandardClaims: jwt.StandardClaims{
	// 		ExpiresAt: time.Now().Add(time.Hour * 5).Unix(),
	// 	},
	// })

	// token, err := claims.SignedString([]byte(JwtKey))

	// if err != nil {
	// 	return "", fmt.Errorf("could not generate jwt token")
	// }
	// return token, nil
}

func (d *Database) VerifyToken(signedToken string) (*model.User, error) {
	token, err := jwt.ParseWithClaims(
		signedToken,
		&model.JwtClaims{},
		func(t *jwt.Token) (interface{}, error) {
			return []byte(JwtKey), nil
		},
	)
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*model.JwtClaims)
	if !ok {
		return nil, fmt.Errorf("couldn't parse the jwt claims")
	}

	if claims.ExpiresAt < time.Now().Local().Unix() {
		return nil, fmt.Errorf("jwt token is expired")
	}

	return &claims.User, nil
}
