package types

import (
	"github.com/oxiginedev/gostarter/internal/models"
	"github.com/oxiginedev/gostarter/internal/pkg/jwt"
	"github.com/oxiginedev/gostarter/pkg/validator"
)

type ContextKey string

func (c ContextKey) String() string {
	return string(c)
}

type RegisterUserRequest struct {
	FirstName            string `json:"first_name" valid:"required"`
	LastName             string `json:"last_name" valid:"required"`
	Email                string `json:"email" valid:"required,email"`
	Password             string `json:"password" valid:"required"`
	PasswordConfirmation string `json:"password_confirmation" valid:"required"`
}

func (ru RegisterUserRequest) Validate() error {
	return validator.Validate(ru)
}

type LoginUserRequest struct {
	Email    string `json:"email" valid:"required,email"`
	Password string `json:"password" valid:"required"`
}

func (lu LoginUserRequest) Validate() error {
	return validator.Validate(lu)
}

type LoginResponse struct {
	User  *models.User `json:"user"`
	Token *jwt.Token   `json:"token"`
}

type UserResponse struct {
	User *models.User `json:"user"`
}
