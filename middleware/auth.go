package middleware

import (
	"errors"

	"github.com/labstack/echo"
)

func Auth(next echo.HandlerFunc) echo.HandlerFunc {
	return func(ec echo.Context) error {
		apiKey := ec.Request().Header.Get("Authorization")
		if apiKey != "secretApiKey" {
			ec.Error(errors.New("invalid api-key"))
		}
		return next(ec)
	}
}
