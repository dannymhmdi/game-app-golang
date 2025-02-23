package user_handler

import (
	"github.com/labstack/echo/v4"
	"mymodule/delivery/httpserver/middleware"
)

func (h Handler) SetRoute(e *echo.Echo) {
	userGroup := e.Group("/users")
	userGroup.POST("/register", h.userRegisterHandler)
	userGroup.POST("/login", h.userLoginHandler)
	userGroup.GET("/profile", h.userProfileHandler, middleware.AuthMiddleWare(h.authSignKey, h.authSvc))
}
