package middleware

import (
	"errors"
	"github/oxiginedev/gostarter/internal/database"
	"github/oxiginedev/gostarter/internal/models"
	"github/oxiginedev/gostarter/internal/pkg/jwt"
	"github/oxiginedev/gostarter/util"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
)

const AuthUserCtx = "authUser"

// Authenticates restricts access to protected routes
func Authenticate(jwt *jwt.JWT, db *database.Connection) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			authHeader := c.Request().Header.Get("Authorization")
			authParts := strings.Split(authHeader, " ")

			if len(authParts) != 2 {
				return c.JSON(http.StatusUnauthorized, util.BuildErrorResponse("invalid header structure", nil))
			}

			authToken := authParts[1]

			if util.IsStringEmpty(authToken) {
				return c.JSON(http.StatusUnauthorized, util.BuildErrorResponse("empty token", nil))
			}

			token, err := jwt.ValidateAccessToken(authToken)
			if err != nil {
				return c.JSON(http.StatusUnauthorized, util.BuildErrorResponse("unauthenticated", nil))
			}

			user, err := models.FindUserByID(db, token.UserID)
			if err != nil {
				if errors.Is(err, database.ErrRecordNotFound) {
					return c.JSON(http.StatusUnauthorized, util.BuildErrorResponse("unauthenticated", nil))
				}

				return err
			}

			c.Set(AuthUserCtx, user)

			return next(c)
		}
	}
}

func GetAuthUserFromContext(c echo.Context) *models.User {
	return c.Get(AuthUserCtx).(*models.User)
}
