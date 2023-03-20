package storage

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"strconv"

	"go.uber.org/zap"
	"tschwaa.com/api/model"
)

func (d *Database) GetOrganizationMembers(ctx context.Context, orgId uint64) ([]model.Member, error) {
	members := []model.Member{}

	query := `
		SELECT M.id, M.name, M.sex, M.phone, A.joined
		FROM adhesions A INNER JOIN members M ON A.member_id = M.id
		WHERE A.organization_id = $1
	`
	rows, err := d.DB.QueryContext(ctx, query, orgId)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var id, name, sex, phoneNumber string
		var joined bool
		if err := rows.Scan(&id, &name, &sex, &phoneNumber, &joined); err != nil {
			return nil, fmt.Errorf("error when parsing the organization's members result")
		}

		i, _ := strconv.ParseUint(id, 10, 64)
		members = append(members, model.Member{
			ID:     i,
			Name:   name,
			Sex:    sex,
			Phone:  phoneNumber,
			Joined: joined,
		})
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error when parsing the organization's members result")
	}

	return members, nil
}

func (d *Database) FindMemberByPhoneNumber(ctx context.Context, phone string) (*model.Member, error) {
	var member model.Member

	query := `
		SELECT id, name, sex, phone
		FROM members
		WHERE (phone = $1)
	`
	if err := d.DB.QueryRowContext(ctx, query, phone).Scan(&member.ID, &member.Name, &member.Sex, &member.Phone); err == nil {
		return &member, nil
	} else if err == sql.ErrNoRows {
		return nil, nil
	} else {
		d.log.Info("Error FindUserByUsername ", zap.Error(err))
		return nil, err
	}
}

func (d *Database) CreateMember(ctx context.Context, member model.Member) (uint64, error) {
	query := `
		INSERT INTO members(name, sex, phone)
		VALUES ($1, $2, $3)
		RETURNING id
	`
	var lastInsertId uint64 = 0
	err := d.DB.QueryRowContext(ctx, query, member.Name, member.Sex, member.Phone).Scan(&lastInsertId)
	return lastInsertId, err
}

func (d *Database) FindAdhesionByMemberAndOrg(ctx context.Context, memberId, orgId uint64) (*model.Adhesion, error) {
	var adhesion model.Adhesion

	query := `
		SELECT id, member_id, organization_id
		FROM members
		WHERE (member_id = $1 AND organization_id = $2)
	`
	if err := d.DB.QueryRowContext(ctx, query, memberId, orgId).Scan(&adhesion.ID, &adhesion.MemberID, &adhesion.OrgID); err == nil {
		return &adhesion, nil
	} else if err == sql.ErrNoRows {
		return nil, nil
	} else {
		d.log.Info("Error FindAdhesionByMemberAndOrg ", zap.Error(err))
		return nil, err
	}
}

func (d *Database) FindAdhesionById(ctx context.Context, adhesionId uint64) (*model.Adhesion, error) {
	var adhesion model.Adhesion

	query := `
		SELECT id, member_id, organization_id
		FROM members
		WHERE (id = $1)
	`
	if err := d.DB.QueryRowContext(ctx, query, adhesionId).Scan(&adhesion.ID, &adhesion.MemberID, &adhesion.OrgID); err == nil {
		return &adhesion, nil
	} else if err == sql.ErrNoRows {
		return nil, nil
	} else {
		d.log.Info("Error FindUserByUsername ", zap.Error(err))
		return nil, err
	}
}

func (d *Database) CreateAdhesion(ctx context.Context, memberId, orgId uint64) (uint64, error) {
	adhesion, err := d.FindAdhesionByMemberAndOrg(ctx, memberId, orgId)
	if err != nil {
		return 0, err
	}

	if adhesion != nil {
		return adhesion.ID, nil
	}

	query := `
		INSERT INTO adhesions(member_id, organization_id)
		VALUES ($1, $2)
		RETURNING id
	`
	log.Println("CreateAdhesion", memberId, orgId)
	var lastInsertId uint64 = 0
	err = d.DB.QueryRowContext(ctx, query, memberId, orgId).Scan(&lastInsertId)

	return lastInsertId, err
}

func (d *Database) CreateInvitation(ctx context.Context, link string, adhesionId uint64) (uint64, error) {
	var mid sql.NullInt64
	adhesion, err := d.FindAdhesionById(ctx, adhesionId)
	if err != nil {
		return 0, err
	}

	if adhesion != nil {
		query := `
			UPDATE invitations
			SET active = FALSE
			WHERE adhesion_id = $1 AND active = TRUE
		`
		err = d.DB.QueryRowContext(ctx, query, link, adhesionId).Scan(&mid)
		if err != nil {
			return 0, err
		}
	}

	query := `
		INSERT INTO invitations(link, adhesion_id)
		VALUES ($1, $2)
		RETURNING id
	`
	log.Println("CreateAdhesion", link, adhesionId)
	var lastInsertId uint64 = 0
	err = d.DB.QueryRowContext(ctx, query, link, adhesionId).Scan(&lastInsertId)

	return lastInsertId, err
}
