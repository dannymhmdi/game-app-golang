package backOffice_handler

import (
	"github.com/labstack/echo/v4"
	"mymodule/delivery/httpserver/middleware"
)

func (h Handler) SetBackOfficeRoute(e *echo.Echo) {
	backOfficeGroup := e.Group("/backOffice/users")
	backOfficeGroup.GET("/", h.backOfficeListUsersHandler, middleware.AccessCheckMiddleware(h.auth, h.authorizationSvc))
}
