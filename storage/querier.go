package storage

import (
	"context"
	"database/sql"

	"tschwaa.com/api/models"
)

type Querier interface {
	// User & Member
	GetUserByUsername(ctx context.Context, arg GetUserByUsernameParams) (*models.User, error)
	DoesUserExist(ctx context.Context, phone string) (bool, error)
	UpdateMemberUserID(ctx context.Context, arg UpdateMemberUserIDParams) error
	CreateUser(ctx context.Context, arg CreateUserParams) (*models.User, error)
	GetMemberByUsername(ctx context.Context, arg GetMemberByUsernameParams) (*models.Member, error)
	GetMemberByPhone(ctx context.Context, phone string) (*models.Member, error)
	GetMemberByEmail(ctx context.Context, email string) (*models.Member, error)
	GetMemberByID(ctx context.Context, id uint64) (*models.Member, error)
	CreateMember(ctx context.Context, arg CreateMemberParams) (*models.Member, error)
	UpdateMember(ctx context.Context, arg UpdateMemberParams) error
	// Otp
	CreateOTP(ctx context.Context, arg CreateOTPParams) (*models.Otp, error)
	DeactivateOTP(ctx context.Context, id uint64) error
	GetActiveOTPFromPhone(ctx context.Context, phone string) (*models.Otp, error)
	CheckOTP(ctx context.Context, arg CheckOTPParams) (*models.Otp, error)
	// Organization
	CreateOrganization(ctx context.Context, arg CreateOrganizationParams) (*models.Organization, error)
	GetOrganization(ctx context.Context, id uint64) (*models.Organization, error)
	ListOrganizationOfMember(ctx context.Context, memberID uint64) ([]*models.Organization, error)
	ListOrganizations(ctx context.Context) ([]*models.Organization, error)
	ListOrganizationsCreatedBy(ctx context.Context, createdBy sql.NullInt32) ([]*models.Organization, error)
	// Membership
	DoesMembershipExist(ctx context.Context, arg DoesMembershipExistParams) (*models.Membership, error)
	CreateMembership(ctx context.Context, arg CreateMembershipParams) (*models.Membership, error)
	GetMembersFromOrganization(ctx context.Context, organizationID uint64) ([]*models.OrganizationMember, error)
	GetMembership(ctx context.Context, id uint64) (*models.Membership, error)
	ApprovedMembership(ctx context.Context, id uint64) (*models.Membership, error)
	// Invitation
	CreateInvitation(ctx context.Context, arg CreateInvitationParams) (*models.Invitation, error)
	GetInvitation(ctx context.Context, link string) (*models.Invitation, error)
	GetInvitationLinkFromMembership(ctx context.Context, membershipId uint64) (string, error)
	DesactivateInvitation(ctx context.Context, membershipID uint64) error
	DesactivateInvitationFromLink(ctx context.Context, link string) (*models.Invitation, error)
}

type QuerierTx interface {
	// User
	CreateUserWithMemberTx(ctx context.Context, arg CreateUserWithMemberParams) (uint64, error)
	CreateMemberWithAssociatedUserTx(ctx context.Context, arg CreateMemberWithAssociatedUserParams) error
	// Otp
	CreateOTPTx(ctx context.Context, arg CreateOTPParams) (*models.Otp, error)
	// Organization
	CreateOrganizationWithMembershipTx(ctx context.Context, arg CreateOrganizationParams) (*models.Organization, error)
	// Membership
	CreateInvitationTx(ctx context.Context, arg CreateMembershipInvitationParams) (*models.Organization, error)
	// Invitation
	ApprovedInvitationTx(ctx context.Context, link string) error
}

var _ Querier = (*Queries)(nil)
