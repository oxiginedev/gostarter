package database

import "context"

type UserRepository interface {
	CreateUser(ctx context.Context, user *User) error
	// UpdateUser(ctx context.Context, user *User) error
	FindUserByEmail(ctx context.Context, email string) (*User, error)
	FindUserByID(ctx context.Context, id string) (*User, error)
}
