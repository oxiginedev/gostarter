package handlers

import (
	"github/oxiginedev/gostarter/api/types"
	"github/oxiginedev/gostarter/internal/database/postgres"
	"github/oxiginedev/gostarter/services"
	"github/oxiginedev/gostarter/util"
	"net/http"

	"github.com/labstack/echo/v4"
)

func (h *Handler) LoginUser(c echo.Context) error {
	var authUser types.LoginUser

	if err := c.Bind(&authUser); err != nil {
		return c.JSON(http.StatusBadRequest, util.BuildErrorResponse("Bad request", nil))
	}

	err := authUser.Validate()
	if err != nil {
		errs := util.BuildErrorResponse("The given data was invalid", err)
		return c.JSON(http.StatusUnprocessableEntity, errs)
	}

	ls := services.LoginUserService{
		UserRepo: postgres.NewUserRepository(h.DB.GetDB()),
		JWT:      h.JWT,
		Data:     &authUser,
	}

	user, token, err := ls.Run(c.Request().Context())
	if err != nil {
		return err
	}

	data := &types.LoginResponse{
		User:  user,
		Token: token,
	}

	return c.JSON(http.StatusOK, util.BuildSuccessResponse("Login successful", data))
}
