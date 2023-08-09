package storage

import (
	"context"
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
