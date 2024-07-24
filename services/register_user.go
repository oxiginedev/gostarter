package services

import (
	"context"
	"github/oxiginedev/gostarter/api/types"
	"github/oxiginedev/gostarter/internal/database"
	"github/oxiginedev/gostarter/internal/pkg/jwt"
)

type RegisterUserService struct {
	UserRepo database.UserRepository
	JWT      *jwt.JWT
	Data     *types.RegisterUser
}

func (rs *RegisterUserService) Run(ctx context.Context) (*database.User, *jwt.Token, error) {
	password := &database.Password{PlainText: rs.Data.Password}
	if err := password.GenerateHash(); err != nil {
		return nil, nil, err
	}

	user := &database.User{
		Name:     rs.Data.Name,
		Email:    rs.Data.Email,
		Password: password.Hash,
	}

	err := rs.UserRepo.CreateUser(ctx, user)
	if err != nil {
		return nil, nil, err
	}

	token, err := rs.JWT.GenerateAccessToken(user)
	if err != nil {
		return nil, nil, err
	}

	return user, token, nil
}
