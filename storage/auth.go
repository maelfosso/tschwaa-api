package storage

import (
	"context"
	"crypto/rand"
	"fmt"

	"github.com/fatih/structs"
	"go.uber.org/zap"
	"tschwaa.com/api/models"
	"tschwaa.com/api/services"
)

func createSecret() (string, error) {
	secret := make([]byte, 32)
	if _, err := rand.Read(secret); err != nil {
		return "", err
	}

	return fmt.Sprintf("%x", secret), nil
}

func (d *Database) Signup(ctx context.Context, member models.Member, user models.User) (string, error) {
	// Check if a user with the same email exist
	existingMember, err := d.FindMemberByUsername(ctx, user.Phone, user.Email)
	if err != nil || existingMember != nil {
		return "", fmt.Errorf("member with the email/phone already exists")
	}

	mID, err := d.CreateMember(ctx, member)
	if err != nil {
		return "", fmt.Errorf("error when creating the member: %w", err)
	}
	member.ID = mID
	user.MemberID = mID

	// Hash the password
	if err := user.HashPassword(); err != nil {
		return "", err
	}

	// Get the token - Next will have token for email and token for sms
	token, err := createSecret()
	if err != nil {
		return "", err
	}

	uID, err := d.CreateUser(ctx, user)
	if err != nil {
		return "", fmt.Errorf("error when creating the user: %w", err)
	}
	member.UserID = uID

	return token, err
}

func (d *Database) Signin(ctx context.Context, credentials models.SignInInputs) (*models.SignInResult, error) {
	existingUser, err := d.FindUserByUsername(ctx, credentials.Username, credentials.Username)
	if err != nil || existingUser == nil {
		return nil, fmt.Errorf("user with that username does not exist")
	}

	if existingUser.IsPasswordMatched(credentials.Password) {
		return nil, fmt.Errorf("the password is not correct")
	}

	existingMember, err := d.FindMemberByID(ctx, existingUser.MemberID)
	if err != nil || existingMember == nil {
		return nil, fmt.Errorf("member related to the user does not exist")
	}

	var signInResult models.SignInResult
	signInResult.Name = fmt.Sprintf("%s %s", existingMember.FirstName, existingMember.LastName)
	signInResult.Email = existingMember.Email
	signInResult.ID = existingMember.ID

	tokenString, err := services.GenerateJWTToken(structs.Map(&signInResult))
	if err != nil {
		return nil, err
	}

	d.log.Info("Sign In Token", zap.String("token", tokenString), zap.Any("jwt token", tokenString), zap.Error(err))
	signInResult.Token = tokenString
	return &signInResult, nil
}
