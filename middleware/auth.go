package middleware

import (
	"context"
	"errors"
	"strings"

	"github.com/galihsatriawan/sample-consumer-spotify/internal/domain/repository"
	"github.com/labstack/echo"
)

type authMiddleware struct {
	authRepository repository.Auth
}
type AuthMiddleware interface {
	Auth(next echo.HandlerFunc) echo.HandlerFunc
}

func NewAuthMiddleware(authRepository repository.Auth) AuthMiddleware {
	return authMiddleware{authRepository: authRepository}
}
func (m authMiddleware) Auth(next echo.HandlerFunc) echo.HandlerFunc {
	return func(ec echo.Context) error {
		authHeader := ec.Request().Header.Get("Authorization")
		if !strings.Contains(authHeader, "Bearer") {
			ec.Error(errors.New("invalid auth"))
		}
		arrayToken := strings.Split(authHeader, " ")
		if len(arrayToken) != 2 {
			ec.Error(errors.New("invalid auth"))
		}
		tokenString := arrayToken[1]
		token, err := m.authRepository.Find(context.Background(), tokenString)
		if err != nil {
			ec.Error(errors.New("invalid auth"))
		}
		ec.Set("token", token)
		return next(ec)
	}
}
