package handlers

import (
	"github/oxiginedev/gostarter/api/types"
	"github/oxiginedev/gostarter/services"
	"github/oxiginedev/gostarter/util"
	"net/http"

	"github.com/labstack/echo/v4"
)

func (h *Handler) RegisterUser(c echo.Context) error {
	var newUser types.RegisterUser

	err := c.Bind(&newUser)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	err = newUser.Validate()
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, err)
	}

	rs := services.RegisterUserService{
		Data: &newUser,
	}

	user, token, err := rs.Run(c.Request().Context())
	if err != nil {

	}

	data := &types.LoginResponse{
		User:  user,
		Token: token,
	}

	return c.JSON(http.StatusCreated, util.BuildSuccessResponse("Registration successful", data))
}
