package httpserver

import (
	"github.com/labstack/echo/v4"
	"net/http"
	"time"
)

func (s Server) healthCheck(c echo.Context) error {
	ctx := c.Request().Context()
	select {
	case <-time.After(8 * time.Second):
		return echo.NewHTTPError(http.StatusRequestTimeout, echo.Map{"msg": "timeout"})
	case <-ctx.Done():
		return echo.NewHTTPError(http.StatusOK, echo.Map{"msg": "connected to server successfully!"})
	}

}
