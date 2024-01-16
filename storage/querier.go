package storage

import (
	"context"

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
	ListOrganizationsCreatedBy(ctx context.Context, createdBy uint64) ([]*models.Organization, error)
	// Session
	GetCurrentSession(ctx context.Context, organizationID uint64) (*models.Session, error)
	GetSession(ctx context.Context, arg GetSessionParams) (*models.Session, error)
	NoSessionInProgress(ctx context.Context, organizationID uint64) error
	CreateSession(ctx context.Context, arg CreateSessionParams) (*models.Session, error)
	// Members of session
	ListAllMembersOfSession(ctx context.Context, arg ListAllMembersOfSessionParams) ([]*models.MembersOfSession, error)
	RemoveMemberFromSession(ctx context.Context, arg RemoveMemberFromSessionParams) error
	RemoveAllMembersFromSession(ctx context.Context, arg RemoveAllMembersFromSessionParams) error
	AddMemberToSession(ctx context.Context, arg AddMemberToSessionParams) (*models.MembersOfSession, error)
	// Session Place
	CreateSessionPlace(ctx context.Context, arg CreateSessionPlaceParams) (*models.SessionPlace, error)
	CreateSessionPlaceGivenVenue(ctx context.Context, arg CreateSessionPlaceGivenVenueParams) (*models.SessionPlacesGivenVenue, error)
	CreateSessionPlaceMemberHome(ctx context.Context, sessionPlaceID uint64) (*models.SessionPlacesMemberHome, error)
	CreateSessionPlaceOnline(ctx context.Context, arg CreateSessionPlaceOnlineParams) (*models.SessionPlacesOnline, error)
	DeleteSessionPlaceGivenVenue(ctx context.Context, id uint64) error
	DeleteSessionPlaceMemberHome(ctx context.Context, id uint64) error
	DeleteSessionPlaceOnline(ctx context.Context, id uint64) error
	UpdateSessionPlace(ctx context.Context, arg UpdateSessionPlaceParams) (*models.SessionPlace, error)
	GetSessionPlace(ctx context.Context, sessionID uint64) (*models.SessionPlace, error)
	GetSessionPlaceGiveVenue(ctx context.Context, sessionPlaceID uint64) (*models.SessionPlacesGivenVenue, error)
	GetSessionPlaceMemberHome(ctx context.Context, sessionPlaceID uint64) (*models.SessionPlacesMemberHome, error)
	GetSessionPlaceOnline(ctx context.Context, sessionPlaceID uint64) (*models.SessionPlacesOnline, error)
	// Membership
	DoesMembershipExist(ctx context.Context, arg DoesMembershipExistParams) (*models.Membership, error)
	DoesMembershipConcernOrganization(ctx context.Context, arg DoesMembershipConcernOrganizationParams) (*models.Membership, error)
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
	// Session
	CreateSessionTx(ctx context.Context, arg CreateSessionParams) (*models.Session, error)
	// Members of session
	UpdateSessionMembersTx(ctx context.Context, arg UpdateSessionMembersParams) ([]*models.MembersOfSession, error)
	// Membership
	CreateInvitationTx(ctx context.Context, arg CreateMembershipInvitationParams) (*models.Organization, error)
	// Invitation
	ApprovedInvitationTx(ctx context.Context, link string) error
}

var _ Querier = (*Queries)(nil)
