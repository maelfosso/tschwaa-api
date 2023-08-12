package storage

import (
	"context"
	"database/sql"

	"tschwaa.com/api/models"
)

const getUserByUsername = `
SELECT id, phone, email, password, member_id
FROM users
WHERE (phone = $1) OR (email = $2)
`

type GetUserByUsernameParams struct {
	Phone string
	Email string
}

func (q *Queries) GetUserByUsername(ctx context.Context, arg GetUserByUsernameParams) (*models.User, error) {
	var user models.User

	if err := q.db.QueryRowContext(ctx, getUserByUsername, arg.Phone, arg.Email).Scan(&user.ID, &user.Phone, &user.Email, &user.Password, &user.MemberID); err == nil {
		return &user, nil
	} else if err == sql.ErrNoRows {
		return nil, nil
	} else {
		return nil, err
	}
}

func (q *Queries) DoesUserExist(ctx context.Context, phone string) (bool, error) {
	user, err := q.GetUserByUsername(ctx, GetUserByUsernameParams{Phone: phone, Email: phone})
	return (user != nil), err
}

const updateMemberUserID = `
UPDATE members SET user_id = $1 WHERE id = $2
`

type UpdateMemberUserIDParams struct {
	MemberID uint64
	UserID   uint64
}

func (q *Queries) UpdateMemberUserID(ctx context.Context, arg UpdateMemberUserIDParams) error {
	_, err := q.db.ExecContext(
		ctx,
		updateMemberUserID,
		arg.UserID, arg.MemberID,
	)
	return err
}

const createUser = `
INSERT INTO users(email, phone, password, token, member_id)
VALUES ($1, $2, $3, $4, $5)
RETURNING id, email, phone, password, token, member_id
`

type CreateUserParams struct {
	Phone    string
	Email    string
	Password string
	Token    string
	MemberID uint64
}

func (q *Queries) CreateUser(ctx context.Context, arg CreateUserParams) (*models.User, error) {
	var user *models.User

	err := q.db.QueryRowContext(
		ctx, createUser,
		arg.Email, arg.Phone, arg.Password, arg.Token, arg.MemberID,
	).Scan(
		&user.ID,
		&user.Phone,
		&user.Email,
		&user.Password,
		&user.Token,
		&user.MemberID,
	)
	return user, err
}

const getMemberByPhoneNumber = `
	SELECT id, first_name, last_name, sex, phone, email
	FROM members
	WHERE (phone = $1) OR (email = $2)
`

type GetMemberByUsernameParams struct {
	Phone string
	Email string
}

func (q *Queries) GetMemberByUsername(ctx context.Context, arg GetMemberByUsernameParams) (*models.Member, error) {
	var user models.Member

	if err := q.db.QueryRowContext(ctx, getMemberByPhoneNumber, arg.Phone, arg.Email).Scan(&user.ID, &user.FirstName, &user.LastName, &user.Sex, &user.Phone, &user.Email); err == nil {
		return &user, nil
	} else if err == sql.ErrNoRows {
		return nil, nil
	} else {
		return nil, err
	}
}

func (q *Queries) GetMemberByPhone(ctx context.Context, phone string) (*models.Member, error) {
	return q.GetMemberByUsername(ctx, GetMemberByUsernameParams{Phone: phone, Email: phone})
}

func (q *Queries) GetMemberByEmail(ctx context.Context, email string) (*models.Member, error) {
	return q.GetMemberByUsername(ctx, GetMemberByUsernameParams{Phone: email, Email: email})
}

const getMemberByID = `
	SELECT id, first_name, last_name, sex, phone, email
	FROM members
	WHERE (id = $1)
`

func (q *Queries) GetMemberByID(ctx context.Context, id uint64) (*models.Member, error) {
	var member models.Member

	if err := q.db.QueryRowContext(ctx, getMemberByID, id).Scan(&member.ID, &member.FirstName, &member.LastName, &member.Sex, &member.Phone, &member.Email); err == nil {
		return &member, nil
	} else if err == sql.ErrNoRows {
		return nil, nil
	} else {
		return nil, err
	}
}

const createMember = `
	INSERT INTO members(first_name, last_name, sex, email, phone)
	VALUES ($1, $2, $3, $4, $5)
	RETURNING id, first_name, last_name, sex, email, phone
`

type CreateMemberParams struct {
	FirstName string
	LastName  string
	Sex       string
	Email     string
	Phone     string
}

func (q *Queries) CreateMember(ctx context.Context, arg CreateMemberParams) (*models.Member, error) {
	var member *models.Member

	err := q.db.QueryRowContext(
		ctx,
		createMember,
		arg.FirstName, arg.LastName, arg.Sex, arg.Email, arg.Phone,
	).Scan(
		&member.ID,
		&member.FirstName,
		&member.LastName,
		&member.Sex,
		&member.Email,
		&member.Phone,
	)

	return member, err
}

const updateMember = `
	UPDATE members
	SET first_name = $1, last_name = $2, email = $3, sex = $4, phone = $5
	WHERE id = $6
`

type UpdateMemberParams struct {
	FirstName string
	LastName  string
	Email     string
	Sex       string
	Phone     string
	ID        uint64
}

func (q *Queries) UpdateMember(ctx context.Context, arg UpdateMemberParams) error {

	_, err := q.db.ExecContext(
		ctx, updateMember,
		arg.FirstName, arg.LastName, arg.Email, arg.Sex, arg.Phone, arg.ID,
	)
	if err != nil {
		return err
	}

	return nil
}
