package handlers

import (
	"net/http"

	"github.com/oxiginedev/gostarter/api/middleware"
	"github.com/oxiginedev/gostarter/api/types"
	"github.com/oxiginedev/gostarter/services/auth"
	"github.com/oxiginedev/gostarter/util"

	"github.com/labstack/echo/v4"
)

func (h *Handler) RegisterUser(c echo.Context) error {
	var newUser types.RegisterUserRequest

	err := c.Bind(&newUser)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	err = newUser.Validate()
	if err != nil {
		errs := util.BuildErrorResponse("The given data was invalid", err)
		return c.JSON(http.StatusUnprocessableEntity, errs)
	}

	rs := auth.RegisterUserService{
		DB:   h.DB,
		JWT:  h.JWT,
		Data: &newUser,
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

func (h *Handler) LoginUser(c echo.Context) error {
	var authUser types.LoginUserRequest

	if err := c.Bind(&authUser); err != nil {
		return c.JSON(http.StatusBadRequest, util.BuildErrorResponse("Bad request", nil))
	}

	err := authUser.Validate()
	if err != nil {
		errs := util.BuildErrorResponse("The given data was invalid", err)
		return c.JSON(http.StatusUnprocessableEntity, errs)
	}

	ls := auth.LoginUserService{
		DB:   h.DB,
		JWT:  h.JWT,
		Data: &authUser,
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

func (h *Handler) VerifyUserEmail(c echo.Context) error {
	token := c.Param("token")

	vu := &auth.VerifyUserEmailService{
		DB:    h.DB,
		Token: token,
	}

	err := vu.Run(c.Request().Context())
	if err != nil {
		return err
	}

	rd := util.BuildSuccessResponse("User email verified successfully", nil)
	return c.JSON(http.StatusOK, rd)
}

func (h *Handler) ResendVerificationEmail(c echo.Context) error {
	user := middleware.GetAuthUserFromContext(c)

	if user.HasVerifiedEmail() {
		rd := util.BuildSuccessResponse("Email already verified", nil)
		return c.JSON(http.StatusOK, rd)
	}

	// call email verification service

	rd := util.BuildSuccessResponse("Email verification code resent", nil)
	return c.JSON(http.StatusOK, rd)
}

func (h *Handler) ForgotPassword(c echo.Context) error {
	return nil
}

func (h *Handler) ResetPassword(c echo.Context) error {
	return nil
}
