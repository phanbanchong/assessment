package auth

import (
	"os"

	"github.com/labstack/echo/v4"
)

func AuthMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			auth := c.Request().Header.Values("Authorization")
			if auth != nil && auth[0] == os.Getenv("AUTHORIZATION") {
				return next(c)
			}
			return echo.ErrUnauthorized
		}
	}
}
