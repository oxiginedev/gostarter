package services

import (
	"context"
	"github/oxiginedev/gostarter/api/types"
	"github/oxiginedev/gostarter/internal/database"
	"github/oxiginedev/gostarter/internal/models"
	"github/oxiginedev/gostarter/internal/pkg/jwt"
)

type RegisterUserService struct {
	DB   *database.Connection
	JWT  *jwt.JWT
	Data *types.RegisterUser
}

func (rs *RegisterUserService) Run(ctx context.Context) (*models.User, *jwt.Token, error) {
	user := &models.User{
		FirstName: rs.Data.FirstName,
		LastName:  rs.Data.LastName,
		Email:     rs.Data.Email,
		Password:  rs.Data.Password,
	}

	err := models.CreateUser(rs.DB, user)
	if err != nil {
		return nil, nil, err
	}

	token, err := rs.JWT.GenerateAccessToken(user)
	if err != nil {
		return nil, nil, err
	}

	return user, token, nil
}
