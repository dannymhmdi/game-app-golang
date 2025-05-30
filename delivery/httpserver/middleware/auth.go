package middleware

import (
	"fmt"
	"mymodule/entity"
	"net/http"
	//echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"mymodule/service/authService"
)

func AuthMiddleWare(signKey []byte, authSvc authService.Service) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			accessToken := c.Request().Header.Get("Authorization")
			if accessToken == "" {
				c.Set("isAccTokenValid", false)
				return next(c)
			}

			claim, pErr := authSvc.ParseToken(accessToken)
			if pErr != nil {
				if pErr.Error() == "token has invalid claims: token is expired" {
					fmt.Println("is expired danny")
					c.Set("isAccTokenValid", false)

					return next(c)
				}

				return echo.NewHTTPError(http.StatusInternalServerError, echo.Map{"message": pErr.Error(), "operation": "middleware.AuthMiddleWare"})
			}

			//date, gErr := claim.GetExpirationTime()
			//if gErr != nil {
			//	return echo.NewHTTPError(http.StatusInternalServerError, echo.Map{"message": gErr.Error(), "operation": "middleware.AuthMiddleWare"})
			//}
			//
			//if date.UnixMilli() < time.Now().UnixMilli() {
			//	c.Set("isAccTokenValid", false)
			//	return next(c)
			//}

			c.Set("claim", claim)
			c.Set("isAccTokenValid", true)
			return next(c)
		}
	}
	//return echojwt.WithConfig(echojwt.Config{
	//	ContextKey:    "claim",
	//	SigningKey:    signKey,
	//	SigningMethod: "HS256",
	//	ParseTokenFunc: func(c echo.Context, auth string) (interface{}, error) {
	//		claim, pErr := authSvc.ParseToken(auth)
	//		if pErr != nil {
	//			return nil, pErr
	//		}
	//		return claim, nil
	//	},
	//})
}

func CheckRefreshToken(authSvc authService.Service) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if c.Get("isAccTokenValid").(bool) {
				return next(c)
			}
			cookie, cErr := c.Request().Cookie("refresh-token")
			if cErr != nil {
				return echo.NewHTTPError(http.StatusUnauthorized, echo.Map{"message": "failed to get refresh-token", "operation": "middleware.CheckRefreshToken"})
			}
			claim, pErr := authSvc.ParseToken(cookie.Value)
			if pErr != nil {
				return echo.NewHTTPError(http.StatusUnauthorized, echo.Map{"message": "failed to parse token", "operation": "middleware.CheckRefreshToken"})
			}
			c.Set("claim", claim)

			fmt.Printf("claim-refresh-token:%+v\n", claim)
			user := entity.User{
				Name: claim.Name,
				ID:   claim.UserId,
				Role: claim.Role,
			}
			generatedNewAccessToken, cErr := authSvc.CreateAccessToken(user)
			if cErr != nil {
				return echo.NewHTTPError(http.StatusUnauthorized, echo.Map{"message": "failed to create access token", "operation": "middleware.CheckRefreshToken"})
			}

			c.Set("generatedNewAccessToken", generatedNewAccessToken)
			return next(c)
		}
	}
}
