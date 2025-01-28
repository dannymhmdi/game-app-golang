package httpserver

import (
	"github.com/labstack/echo/v4"
	"mymodule/dto"
	"mymodule/pkg/richerr"
	"net/http"
)

func (s Server) userRegisterHandler(c echo.Context) error {
	bd := dto.RegisterRequest{}
	if bErr := c.Bind(&bd); bErr != nil {
		return echo.NewHTTPError(http.StatusBadRequest, bErr.Error())
	}

	createdUSer, rErr := s.userSvc.Register(bd)
	if rErr != nil {
		code, msg, op := richerr.CheckTypeErr(rErr)
		return echo.NewHTTPError(richerr.MapKindToHttpErr(code), echo.Map{"message": msg, "operation": op})
	}

	return c.JSON(http.StatusOK, createdUSer)

}

func (s Server) userLoginHandler(c echo.Context) error {

	bd := dto.LoginRequest{}
	if bErr := c.Bind(&bd); bErr != nil {
		return echo.NewHTTPError(http.StatusBadRequest, bErr.Error())
	}

	loginResp, lErr := s.userSvc.Login(bd)
	if lErr != nil {
		code, msg, op := richerr.CheckTypeErr(lErr)
		return echo.NewHTTPError(richerr.MapKindToHttpErr(code), echo.Map{"message": msg, "operation": op})
	}
	return c.JSON(http.StatusOK, loginResp)
}

func (s Server) userProfileHandler(c echo.Context) error {
	token := c.Request().Header.Get("Authorization")
	if token == "" {
		return echo.NewHTTPError(http.StatusUnauthorized, "token is empty")
	}

	claim, pErr := s.authSvc.ParseToken(token)
	if pErr != nil {
		code, msg, op := richerr.CheckTypeErr(pErr)
		return echo.NewHTTPError(richerr.MapKindToHttpErr(code), echo.Map{"message": msg, "operation": op})
	}

	userInfo, gErr := s.userSvc.GetUserProfile(claim.UserId)
	if gErr != nil {
		code, msg, op := richerr.CheckTypeErr(gErr)
		return echo.NewHTTPError(richerr.MapKindToHttpErr(code), echo.Map{"message": msg, "operation": op})
	}

	return c.JSON(http.StatusOK, userInfo)
}
