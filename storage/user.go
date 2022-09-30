package storage

import (
	"context"
	"fmt"

	"go.uber.org/zap"
	"tschwaa.com/api/model"
)

func (d *Database) FindUserByUsername(ctx context.Context, phone, email string) (*model.User, error) {
	var user model.User

	query := `
		SELECT id, firstname, lastname, phone, email, password
		FROM users
		WHERE (phone = $1) OR (email = $2)
	`
	err := d.DB.QueryRowContext(ctx, query, phone, email).Scan(&user.ID, &user.Firstname, &user.Lastname, &user.Phone, &user.Email, &user.Password)
	return &user, err
}

func (d *Database) FindUserByEmail(ctx context.Context, email string) (*model.User, error) {
	var user model.User

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
