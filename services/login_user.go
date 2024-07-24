package services

import (
	"context"
	"errors"
	"github/oxiginedev/gostarter/api/types"
	"github/oxiginedev/gostarter/internal/database"
	"github/oxiginedev/gostarter/internal/pkg/jwt"
	"net/http"
)

const ErrInvalidCredentials = "invalid email or password"

type LoginUserService struct {
	UserRepo database.UserRepository
	JWT      *jwt.JWT
	Data     *types.LoginUser
}

func (lu *LoginUserService) Run(ctx context.Context) (*database.User, *jwt.Token, error) {
	user, err := lu.UserRepo.FindUserByEmail(ctx, lu.Data.Email)
	if err != nil {
		if errors.Is(err, database.ErrRecordNotFound) {
			return nil, nil, newServiceError(http.StatusBadRequest, ErrInvalidCredentials, err)
		}
		return nil, nil, err
	}

	password := &database.Password{
		PlainText: lu.Data.Password,
		Hash:      user.Password,
	}

	matches, err := password.Matches()
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
