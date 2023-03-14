package storage

import (
	"context"
	"fmt"
	"strconv"

	"tschwaa.com/api/model"
)

func (d *Database) GetOrganizationMembers(ctx context.Context, orgId uint64) ([]model.Member, error) {
	members := []model.Member{}

	query := `
		SELECT M.id, M.name, M.sex, M.phone_number
		FROM adhesions A INNER JOIN members M ON A.member_id = M.id
		WHERE A.organization_id = $1
	`
	rows, err := d.DB.QueryContext(ctx, query, orgId)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var id, name, sex, phoneNumber string
		if err := rows.Scan(&id, &name, &sex, &phoneNumber); err != nil {
			return nil, fmt.Errorf("error when parsing the organization's members result")
		}

		i, _ := strconv.ParseUint(id, 10, 64)
		members = append(members, model.Member{
			ID:          i,
			Name:        name,
			Sex:         sex,
			PhoneNumber: phoneNumber,
		})
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error when parsing the organization's members result")
	}

	return members, nil
}
