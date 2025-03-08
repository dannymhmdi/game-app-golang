package middleware

import (
	"context"
	"github.com/labstack/echo/v4"
	"mymodule/params"
	"mymodule/pkg/richerr"
	"mymodule/service/authService"
	"mymodule/service/presenceService"
	"net/http"
	"time"
)

func PresenceMiddleWare(authSvc authService.Service, presenceSvc presenceService.Service) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			token := c.Request().Header.Get("Authorization")
			claim, pErr := authSvc.ParseToken(token)
			if pErr != nil {
				return echo.NewHTTPError(http.StatusUnauthorized, echo.Map{"message": pErr.Error()})
			}
			req := params.PresenseRequest{
				UserId:    claim.UserId,
				Timestamp: time.Now().UnixMilli(),
			}
			res, pErr := presenceSvc.Presence(context.Background(), req)
			if pErr != nil {
				code, msg, op := richerr.CheckTypeErr(pErr)
				return echo.NewHTTPError(richerr.MapKindToHttpErr(code), echo.Map{"message": msg, "operation": op})
			}

			return c.JSON(http.StatusOK, res)
		}
	}
}
