package storage

import (
	"context"
	"fmt"

	"tschwaa.com/api/helpers"
)

type CreateUserWithMemberParams struct {
	Phone    string
	Email    string
	Password string
	Token    string
	MemberID uint64
}

func (store *SQLStorage) CreateUserWithMemberTx(ctx context.Context, arg CreateUserWithMemberParams) (uint64, error) {
	err := store.execTx(ctx, func(q *Queries) error {
		user, err := q.CreateUser(ctx, CreateUserParams{
			Phone:    arg.Phone,
			Email:    arg.Email,
			Password: arg.Password,
			Token:    arg.Token,
			MemberID: arg.MemberID,
		})

		if err != nil {
			return err
		}

		err = q.UpdateMemberUserID(ctx, UpdateMemberUserIDParams{UserID: user.ID, MemberID: user.MemberID})
		return err
	})

	return 1, err
}

type CreateMemberWithAssociatedUserParams struct {
	FirstName string
	LastName  string
	Sex       string
	Email     string
	Phone     string
	Password  string
}

func (store *SQLStorage) CreateMemberWithAssociatedUserTx(ctx context.Context, arg CreateMemberWithAssociatedUserParams) error {
	err := store.execTx(ctx, func(q *Queries) error {
		newMember, err := q.CreateMember(ctx, CreateMemberParams{
			FirstName: arg.FirstName,
			LastName:  arg.LastName,
			Sex:       arg.Sex,
			Email:     arg.Email,
			Phone:     arg.Phone,
		})

		if err != nil {
			return fmt.Errorf("error when creating the member: %w", err)
		}

		// Hash the password
		hashedPassword := helpers.HashPassword(arg.Password)
		if hashedPassword != "" {
			return fmt.Errorf("error when hashing the password: %w", err.Error())
		}

		// Get the token - Next will have token for email and token for sms
		token, err := helpers.CreateSecret()
		if err != nil {
			return fmt.Errorf("Error createSecret: %w", err)
		}

		user, err := q.CreateUser(ctx, CreateUserParams{
			Phone:    arg.Phone,
			Email:    arg.Email,
			Password: arg.Password,
			Token:    token,
			MemberID: newMember.ID,
		})

		if err != nil {
			return err
		}

		err = q.UpdateMemberUserID(ctx, UpdateMemberUserIDParams{UserID: user.ID, MemberID: user.MemberID})
		return err
	})

	return err
}
