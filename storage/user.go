package storage

import (
	"context"
	"database/sql"

	"go.uber.org/zap"
	"tschwaa.com/api/models"
)

func (d *Database) FindUserByUsername(ctx context.Context, phone, email string) (*models.User, error) {
	var user models.User

	query := `
		SELECT id, phone, email, password, member_id
		FROM users
		WHERE (phone = $1) OR (email = $2)
	`
	if err := d.DB.QueryRowContext(ctx, query, phone, email).Scan(&user.ID, &user.Phone, &user.Email, &user.Password, &user.MemberID); err == nil {
		return &user, nil
	} else if err == sql.ErrNoRows {
		return nil, nil
	} else {
		d.log.Info("Error FindMemberByUsername ", zap.Error(err))
		return nil, err
	}
}

func (d *Database) CreateUser(ctx context.Context, user models.User) (uint64, error) {
	tx, err := d.DB.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelSerializable})
	if err != nil {
		return 0, err
	}

	query := `
		INSERT INTO users(email, phone, password, token, member_id)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id
	`
	var lastInsertId uint64 = 0
	err = tx.QueryRowContext(
		ctx, query,
		user.Email, user.Phone, user.Password, user.Token, user.MemberID,
	).Scan(&lastInsertId)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	query = `
		UPDATE members SET user_id = $1 WHERE id = $2
	`
	_, err = tx.ExecContext(
		ctx, query,
		lastInsertId, user.MemberID,
	)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	if err := tx.Commit(); err != nil {
		return 0, err
	}

	return lastInsertId, err
}

func (d *Database) FindMemberByUsername(ctx context.Context, phone, email string) (*models.Member, error) {
	var user models.Member

	query := `
		SELECT id, first_name, last_name, sex, phone, email
		FROM members
		WHERE (phone = $1) OR (email = $2)
	`
	if err := d.DB.QueryRowContext(ctx, query, phone, email).Scan(&user.ID, &user.FirstName, &user.LastName, &user.Sex, &user.Phone, &user.Email); err == nil {
		return &user, nil
	} else if err == sql.ErrNoRows {
		return nil, nil
	} else {
		d.log.Info("Error FindMemberByUsername ", zap.Error(err))
		return nil, err
	}
}

func (d *Database) FindMemberByPhoneNumber(ctx context.Context, phone string) (*models.Member, error) {
	return d.FindMemberByUsername(ctx, phone, phone)
}

func (d *Database) FindMemberByEmail(ctx context.Context, email string) (*models.Member, error) {
	return d.FindMemberByUsername(ctx, email, email)
}

func (d *Database) FindMemberByID(ctx context.Context, id uint64) (*models.Member, error) {
	var member models.Member

	query := `
		SELECT id, first_name, last_name, sex, phone, email
		FROM members
		WHERE (id = $1)
	`
	if err := d.DB.QueryRowContext(ctx, query, id).Scan(&member.ID, &member.FirstName, &member.LastName, &member.Sex, &member.Phone, &member.Email); err == nil {
		return &member, nil
	} else if err == sql.ErrNoRows {
		return nil, nil
	} else {
		d.log.Info("Error FindMemberByUsername ", zap.Error(err))
		return nil, err
	}
}

func (d *Database) CreateMember(ctx context.Context, member models.Member) (uint64, error) {
	d.log.Info("create member", zap.String("firstname", member.FirstName), zap.String("lastname", member.LastName))
	query := `
		INSERT INTO members(first_name, last_name, sex, email, phone)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id
	`
	var lastInsertId uint64 = 0
	err := d.DB.QueryRowContext(
		ctx, query,
		member.FirstName, member.LastName, member.Sex, member.Email, member.Phone,
	).Scan(&lastInsertId)
	d.log.Info("Create Member : ", zap.Error(err))
	return lastInsertId, err
}

func (d *Database) UpdateMember(ctx context.Context, member models.Member) error {
	query := `
		UPDATE members
		SET first_name = $1, last_name = $2, email = $3, sex = $4, phone = $5
		WHERE id = $6
	`
	_, err := d.DB.ExecContext(
		ctx, query,
		member.FirstName, member.LastName, member.Email, member.Sex, member.Phone, member.ID,
	)
	if err != nil {
		return err
	}

	return nil
}
