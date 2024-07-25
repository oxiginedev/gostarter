package services

import (
	"context"
	"errors"
	"github/oxiginedev/gostarter/api/types"
	"github/oxiginedev/gostarter/internal/database"
	"github/oxiginedev/gostarter/internal/models"
	"github/oxiginedev/gostarter/internal/pkg/jwt"
	"github/oxiginedev/gostarter/util"
	"net/http"
)

const ErrInvalidCredentials = "invalid email or password"

type LoginUserService struct {
	DB   *database.Connection
	JWT  *jwt.JWT
	Data *types.LoginUser
}

func (lu *LoginUserService) Run(ctx context.Context) (*models.User, *jwt.Token, error) {
	user, err := models.FindUserByEmailAddress(lu.DB, lu.Data.Email)
	if err != nil {
		if errors.Is(err, database.ErrRecordNotFound) {
			return nil, nil, newServiceError(http.StatusBadRequest, ErrInvalidCredentials, err)
		}
		return nil, nil, err
	}

	matches, err := util.CompareHashAndPassword(user.Password, lu.Data.Password)
	if err != nil {
		return nil, nil, err
	}

	if !matches {
		return nil, nil, newServiceError(http.StatusBadRequest, ErrInvalidCredentials, err)
	}

	token, err := lu.JWT.GenerateAccessToken(user)
	if err != nil {
		return nil, nil, err
	}

	return user, token, nil
}
