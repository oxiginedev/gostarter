package handlers

import (
	"github/oxiginedev/gostarter/api/middleware"
	"github/oxiginedev/gostarter/api/types"
	"github/oxiginedev/gostarter/internal/database/postgres"
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
		errs := util.BuildErrorResponse("The given data was invalid", err)
		return c.JSON(http.StatusUnprocessableEntity, errs)
	}

	rs := services.RegisterUserService{
		UserRepo: postgres.NewUserRepository(h.DB.GetDB()),
		JWT:      h.JWT,
		Data:     &newUser,
	}

	user, token, err := rs.Run(c.Request().Context())
	if err != nil {
		return err
	}

	data := &types.LoginResponse{
		User:  user,
		Token: token,
	}

	return c.JSON(http.StatusCreated, util.BuildSuccessResponse("Registration successful", data))
}

func (h *Handler) GetUser(c echo.Context) error {
	authUser := middleware.GetAuthUserFromContext(c)

	data := &types.UserResponse{User: authUser}
	return c.JSON(http.StatusOK, util.BuildSuccessResponse("User retrieved successfully", data))
}
