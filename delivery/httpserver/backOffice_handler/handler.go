package backOffice_handler

import (
	"github.com/labstack/echo/v4"
	"mymodule/service/authorizationService"
	"mymodule/service/authservice"
	"mymodule/service/backoffice"
	"net/http"
)

type Handler struct {
	backOfficeSvc    backoffice.Service
	auth             authservice.Service
	authorizationSvc authorizationService.Service
}

func New(backofficeSvc backoffice.Service, auth authservice.Service, authorizationSvc authorizationService.Service) *Handler {
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
