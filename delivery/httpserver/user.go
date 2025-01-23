package httpserver

import (
	"github.com/labstack/echo/v4"
	"mymodule/service/registerservice"
	"net/http"
)

func (s Server) userRegisterHandler(c echo.Context) error {
	bd := registerservice.RegisterRequest{}
	if bErr := c.Bind(&bd); bErr != nil {
		return echo.NewHTTPError(http.StatusBadRequest, bErr.Error())
	}

	createdUSer, rErr := s.userSvc.RegisterUser(bd)
	if rErr != nil {
		return echo.NewHTTPError(http.StatusBadRequest, rErr.Error())
	}

	return c.JSON(http.StatusOK, createdUSer)

}

func (s Server) userLoginHandler(c echo.Context) error {

	bd := registerservice.LoginRequest{}
	if bErr := c.Bind(&bd); bErr != nil {
		return echo.NewHTTPError(http.StatusBadRequest, bErr.Error())
	}

	loginResp, lErr := s.userSvc.Login(bd)
	if lErr != nil {
		return echo.NewHTTPError(http.StatusBadRequest, lErr.Error())
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
		return echo.NewHTTPError(http.StatusUnauthorized, pErr.Error())
	}

	userInfo, gErr := s.userSvc.GetUserProfile(claim.UserId)
	if gErr != nil {
		return echo.NewHTTPError(http.StatusBadRequest, gErr.Error())
	}

	return c.JSON(http.StatusOK, userInfo)
}
