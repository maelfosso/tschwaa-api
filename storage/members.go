package storage

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"strconv"
	"time"

	"go.uber.org/zap"
	"tschwaa.com/api/models"
)

func (d *Database) GetOrganizationMembers(ctx context.Context, orgId uint64) ([]models.OrganizationMember, error) {
	members := []models.OrganizationMember{}

	query := `
		SELECT M.id, M.first_name, M.last_name, M.sex, M.phone, A.joined, A.joined_at, A.position, A.role, A.status
		FROM adhesions A INNER JOIN members M ON A.member_id = M.id
		WHERE A.organization_id = $1
	`
	rows, err := d.DB.QueryContext(ctx, query, orgId)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var id, firstName, lastName, sex, phoneNumber, position, role, status string
		var joined bool
		var joinedAt sql.NullTime
		if err := rows.Scan(&id, &firstName, &lastName, &sex, &phoneNumber, &joined, &joinedAt, &position, &role, &status); err != nil {
			return nil, fmt.Errorf("error when parsing the organization's members result", err)
		}

		i, _ := strconv.ParseUint(id, 10, 64)
		members = append(members, models.OrganizationMember{
			ID:        i,
			FirstName: firstName,
			LastName:  lastName,
			Sex:       sex,
			Phone:     phoneNumber,

			Joined:  joined,
			JointAt: joinedAt.Time,

			Position: position,
			Role:     role,
			Status:   status,
		})
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error when parsing the organization's members result")
	}

	return members, nil
}

func (d *Database) FindAdhesionByMemberAndOrg(ctx context.Context, memberId, orgId uint64) (*models.Adhesion, error) {
	var adhesion models.Adhesion

	query := `
		SELECT id, member_id, organization_id
		FROM adhesions
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

func (d *Database) FindAdhesionById(ctx context.Context, adhesionId uint64) (*models.Adhesion, error) {
	var adhesion models.Adhesion

	query := `
		SELECT id, member_id, organization_id
		FROM adhesions
		WHERE (id = $1)
	`
	if err := d.DB.QueryRowContext(ctx, query, adhesionId).Scan(&adhesion.ID, &adhesion.MemberID, &adhesion.OrgID); err == nil {
		return &adhesion, nil
	} else if err == sql.ErrNoRows {
		return nil, nil
	} else {
		d.log.Info("Error FindMemberByUsername ", zap.Error(err))
		return nil, err
	}
}

func (d *Database) CreateAdhesion(ctx context.Context, memberId, orgId uint64, joined bool) (uint64, error) {
	adhesion, err := d.FindAdhesionByMemberAndOrg(ctx, memberId, orgId)
	if err != nil {
		return 0, err
	}

	if adhesion != nil {
		return adhesion.ID, nil
	}

	log.Println("CreateAdhesion", memberId, orgId)
	var query string
	var lastInsertId uint64 = 0
	if joined {
		query = `
			INSERT INTO adhesions(member_id, organization_id, joined, joined_at)
			VALUES ($1, $2, $3, $4)
			RETURNING id
		`
		err = d.DB.QueryRowContext(ctx, query, memberId, orgId, joined, time.Now()).Scan(&lastInsertId)
	} else {
		query = `
			INSERT INTO adhesions(member_id, organization_id)
			VALUES ($1, $2)
			RETURNING id
		`
		err = d.DB.QueryRowContext(ctx, query, memberId, orgId).Scan(&lastInsertId)
	}

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
		err = d.DB.QueryRowContext(ctx, query, adhesionId).Scan(&mid)
		if err != nil && err != sql.ErrNoRows {
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

func (d *Database) GetInvitation(ctx context.Context, link string) (*models.Invitation, error) {
	var invitation models.Invitation

	query := `
		SELECT link, active, i.created_at,
			a.joined, a.member_id, a.organization_id,
			m.id, m.first_name, m.last_name, m.sex, m.phone, m.email, m.user_id,
			o.id, o.name, o.description
		FROM invitations i
		INNER JOIN adhesions a ON i.adhesion_id = a.id
		INNER JOIN members m ON a.member_id = m.id
		INNER JOIN organizations o ON a.organization_id = o.id
		WHERE link = $1
	`
	if err := d.DB.QueryRowContext(ctx, query, link).Scan(
		&invitation.Link, &invitation.Active, &invitation.CreatedAt,
		&invitation.Adhesion.Joined, &invitation.Adhesion.MemberID, &invitation.Adhesion.OrgID,
		&invitation.Member.ID, &invitation.Member.FirstName, &invitation.Member.LastName, &invitation.Member.Sex, &invitation.Member.Phone, &invitation.Member.Email, &invitation.Member.UserID,
		&invitation.Organization.ID, &invitation.Organization.Name, &invitation.Organization.Description,
	); err == nil {
		return &invitation, nil
	} else {
		return nil, err
	}
}

func (d *Database) DisableInvitation(ctx context.Context, link string) (uint64, error) {
	var adhesionId uint64

	query := `
		UPDATE invitations
		SET active = FALSE
		WHERE link = $1
		RETURNING adhesion_id
	`
	err := d.DB.QueryRowContext(ctx, query, link).Scan(&adhesionId)
	if err != nil && err != sql.ErrNoRows {
		return 0, err
	}

	return adhesionId, nil
}

func (d *Database) ApprovedAdhesion(ctx context.Context, adhesionID uint64) error {
	var mid sql.NullInt64
	query := `
		UPDATE adhesions
		SET joined = TRUE, joined_at = NOW()
		WHERE id = $1
	`
	err := d.DB.QueryRowContext(ctx, query, adhesionID).Scan(&mid)
	if err != nil && err != sql.ErrNoRows {
		return err
	}

	return nil
}
