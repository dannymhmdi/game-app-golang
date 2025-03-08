package middleware

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"mymodule/params"
	"mymodule/pkg/richerr"
	"mymodule/service/authService"
	"mymodule/service/authorizationService"
	"net/http"
)

func AccessCheckMiddleware(authSvc authService.Service, authorizationSvc authorizationService.Service) echo.MiddlewareFunc {

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			tokenString := c.Request().Header.Get("Authorization")
			claims, pErr := authSvc.ParseToken(tokenString)

			if pErr != nil {
				return echo.NewHTTPError(http.StatusInternalServerError, pErr.Error())
			}

			//var bd dto.PermissionRequest

			bd := params.PermissionRequest{}
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
