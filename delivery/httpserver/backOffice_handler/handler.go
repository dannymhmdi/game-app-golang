package backOffice_handler

import (
	"github.com/labstack/echo/v4"
	"mymodule/service/authService"
	"mymodule/service/authorizationService"
	"mymodule/service/backofficeService"
	"net/http"
)

type Handler struct {
	backOfficeSvc    backofficeService.Service
	auth             authService.Service
	authorizationSvc authorizationService.Service
}

func New(backofficeSvc backofficeService.Service, auth authService.Service, authorizationSvc authorizationService.Service) *Handler {
	return &Handler{
		backOfficeSvc:    backofficeSvc,
		auth:             auth,
		authorizationSvc: authorizationSvc,
	}
}

func (h Handler) backOfficeListUsersHandler(c echo.Context) error {
	user, lErr := h.backOfficeSvc.ListUsers()
	if lErr != nil {

		return echo.NewHTTPError(http.StatusInternalServerError, lErr.Error())
	}
	return c.JSON(http.StatusOK, echo.Map{"message": "users list", "users": user})
}
