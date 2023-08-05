package storage

import (
	"context"

	"tschwaa.com/api/models"
)

func (d *Database) CreateUserIfNotExists(ctx context.Context, phoneNumber, name string) error {
	return nil
}

func (d *Database) CreateOTP(ctx context.Context, pinCode models.OTP) error {
	return nil
}

func (d *Database) SaveOTP(ctx context.Context, pinCode models.OTP) error {
	return nil
}

func (d *Database) CheckOTP(ctx context.Context, phoneNumber, pinCode string) (*models.OTP, error) {
	return nil, nil
}

func (d *Database) FindUserByPhoneNumber(ctx context.Context, phoneNumber string) (*models.User, error) {
	return nil, nil
}
