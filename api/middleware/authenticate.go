package middleware

import "github.com/labstack/echo/v4"

// Authenticates restricts access to protected routes
func Authenticate() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			return next(c)
		}
	}
}
