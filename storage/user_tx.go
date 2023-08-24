package storage

import (
	"context"
	"fmt"

	"tschwaa.com/api/helpers"
	"tschwaa.com/api/utils"
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
			return utils.Fail(
				fmt.Sprintf("error when creating an user %s", arg.Phone),
				"ERR_CRT_USR_MBR_01", err)
		}

		err = q.UpdateMemberUserID(ctx, UpdateMemberUserIDParams{UserID: user.ID, MemberID: user.MemberID})
		return utils.Fail(
			fmt.Sprintf("error when updating member[%d] user[%d]", user.MemberID, user.ID),
			"ERR_CRT_USR_MBR_01", err)
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
			return utils.Fail(
				"error when creating the member",
				"ERR_CRT_MBR_USR_01", err)
		}

		// Hash the password
		hashedPassword, err := helpers.HashPassword(arg.Password)
		if hashedPassword == "" || err != nil {
			return utils.Fail(
				"error when hashing the password",
				"ERR_CRT_MBR_USR_02", err)
		}

		// Get the token - Next will have token for email and token for sms
		token, err := helpers.CreateSecret()
		if err != nil {
			return utils.Fail(
				"error createSecret",
				"ERR_CRT_MBR_USR_03", err)
		}

		user, err := q.CreateUser(ctx, CreateUserParams{
			Phone:    arg.Phone,
			Email:    arg.Email,
			Password: arg.Password,
			Token:    token,
			MemberID: newMember.ID,
		})

		if err != nil {
			return utils.Fail(
				"error creating the user",
				"ERR_CRT_MBR_USR_04", err)
		}

		err = q.UpdateMemberUserID(ctx, UpdateMemberUserIDParams{UserID: user.ID, MemberID: user.MemberID})
		return utils.Fail(
			fmt.Sprintf("error updating the member[%d] with user[%d]", user.MemberID, user.ID),
			"ERR_CRT_MBR_USR_03", err)
	})

	return err
}
