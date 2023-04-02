package storage

import (
	"context"
	"database/sql"
	"fmt"

	"go.uber.org/zap"
	"tschwaa.com/api/models"
)

func (d *Database) FindUserByUsername(ctx context.Context, phone, email string) (*models.User, error) {
	var user models.User

	query := `
		SELECT id, firstname, lastname, phone, email, password
		FROM users
		WHERE (phone = $1) OR (email = $2)
	`
	if err := d.DB.QueryRowContext(ctx, query, phone, email).Scan(&user.ID, &user.Firstname, &user.Lastname, &user.Phone, &user.Email, &user.Password); err == nil {
		return &user, nil
	} else if err == sql.ErrNoRows {
		return nil, nil
	} else {
		d.log.Info("Error FindUserByUsername ", zap.Error(err))
		return nil, err
	}
}

func (d *Database) FindUserByEmail(ctx context.Context, email string) (*models.User, error) {
	var user models.User

	query := `
		SELECT id, firstname, lastname, phone, email
		FROM users
		WHERE email = $1
	`
	err := d.DB.QueryRowContext(ctx, query, email).Scan(&user.ID, &user.Firstname, &user.Lastname, &user.Phone, &user.Email)
	if err != nil {
		d.log.Info("Error FindUserByEmail", zap.Error(err))
		return nil, fmt.Errorf("user with that email does not exist")
	}

	return &user, nil
}
