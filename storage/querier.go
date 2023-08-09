package storage

import (
	"context"

	"tschwaa.com/api/models"
)

type Querier interface {
	// User & Member
	GetUserByUsername(ctx context.Context, arg GetUserByUsernameParams) (*models.User, error)
	DoesMemberExist(ctx context.Context, phone string) (bool, error)
	UpdateMemberUserID(ctx context.Context, arg UpdateMemberUserIDParams) error
	CreateUser(ctx context.Context, arg CreateUserParams) (*models.User, error)
	GetMemberByUsername(ctx context.Context, arg GetMemberByUsernameParams) (*models.Member, error)
	GetMemberByPhoneNumber(ctx context.Context, phone string) (*models.Member, error)
	GetMemberByEmail(ctx context.Context, email string) (*models.Member, error)
	GetMemberByID(ctx context.Context, id uint64) (*models.Member, error)
	CreateMember(ctx context.Context, arg CreateMemberParams) (*models.Member, error)
	UpdateMember(ctx context.Context, arg UpdateMemberParams) error
}

type QuerierTx interface {
	// User
	CreateUserWithMemberTx(ctx context.Context, arg CreateUserWithMemberParams) (uint64, error)
}

var _ Querier = (*Queries)(nil)
