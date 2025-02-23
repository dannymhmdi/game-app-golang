package middleware

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"mymodule/dto"
	"mymodule/pkg/richerr"
	"mymodule/service/authorization"
	"mymodule/service/authservice"
	"net/http"
)

func AccessCheckMiddleware(authSvc authservice.Service, authorizationSvc authorization.Service) echo.MiddlewareFunc {

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			tokenString := c.Request().Header.Get("Authorization")
			claims, pErr := authSvc.ParseToken(tokenString)

			if pErr != nil {
				return echo.NewHTTPError(http.StatusInternalServerError, pErr.Error())
			}

			//var bd dto.PermissionRequest

			bd := dto.PermissionRequest{}
			if bErr := c.Bind(&bd); bErr != nil {
				return echo.NewHTTPError(http.StatusInternalServerError, bErr.Error())
			}
			fmt.Printf("%+v\n", bd)
			isAllowed, cErr := authorizationSvc.CheckAccess(claims.UserId, claims.Role, bd.PermissionTitles...)
			if cErr != nil {
				fmt.Println("check", cErr)
				code, msg, operation := richerr.CheckTypeErr(cErr)
				return echo.NewHTTPError(richerr.MapKindToHttpErr(code), echo.Map{"message": msg, "operation": operation})
			}

			if !isAllowed {
				return echo.ErrUnauthorized
			}

			return next(c)
		}
	}
}
