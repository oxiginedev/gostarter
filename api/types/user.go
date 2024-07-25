package types

import (
	"github/oxiginedev/gostarter/internal/models"
	"github/oxiginedev/gostarter/internal/pkg/jwt"
	"github/oxiginedev/gostarter/pkg/validator"
)

type ContextKey string

func (c ContextKey) String() string {
	return string(c)
}

type RegisterUser struct {
	FirstName            string `json:"first_name" valid:"required"`
	LastName             string `json:"last_name" valid:"required"`
	Email                string `json:"email" valid:"required,email"`
	Password             string `json:"password" valid:"required"`
	PasswordConfirmation string `json:"password_confirmation" valid:"required"`
}

func (ru RegisterUser) Validate() error {
	return validator.Validate(ru)
}

type LoginUser struct {
	Email    string `json:"email" valid:"required,email"`
	Password string `json:"password" valid:"required"`
}

func (lu LoginUser) Validate() error {
	return validator.Validate(lu)
}

type LoginResponse struct {
	User  *models.User `json:"user"`
	Token *jwt.Token   `json:"token"`
}

type UserResponse struct {
	User *models.User `json:"user"`
}
