package handlers

import (
	"net/http"

	"github.com/oxiginedev/gostarter/api/middleware"
	"github.com/oxiginedev/gostarter/api/types"
	"github.com/oxiginedev/gostarter/util"

	"github.com/labstack/echo/v4"
)

func (h *Handler) GetUser(c echo.Context) error {
	authUser := middleware.GetAuthUserFromContext(c)

	data := &types.UserResponse{User: authUser}
	return c.JSON(http.StatusOK, util.BuildSuccessResponse("User retrieved successfully", data))
}
