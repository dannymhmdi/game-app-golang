package matchMaking_handler

import (
	"github.com/labstack/echo/v4"
	"mymodule/delivery/httpserver/middleware"
)

func (h Handler) SetMatchMakingRoute(e *echo.Echo) {
	userGroup := e.Group("/game")
	userGroup.POST("/match-making", h.MatchMakingHandler, middleware.AuthMiddleWare(h.signingKey, h.authSvc))
}
