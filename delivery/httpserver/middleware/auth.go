package middleware

import (
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"mymodule/service/authservice"
)

func AuthMiddleWare(signKey []byte, authSvc authservice.Service) echo.MiddlewareFunc {
	return echojwt.WithConfig(echojwt.Config{
		ContextKey:    "claim",
		SigningKey:    signKey,
		SigningMethod: "HS256",
		ParseTokenFunc: func(c echo.Context, auth string) (interface{}, error) {
			claim, pErr := authSvc.ParseToken(auth)
			if pErr != nil {
				return nil, pErr
			}
			return claim, nil
		},
	})
}
