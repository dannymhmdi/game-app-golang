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

//echojwt.WithConfig(echojwt.Config{
//ContextKey:    "claim",
//SigningKey:    h.authSignKey,
//SigningMethod: "HS256",
//ParseTokenFunc: func(c echo.Context, auth string) (interface{}, error) {
//claim, pErr := h.authSvc.ParseToken(auth)
//if pErr != nil {
//return nil, pErr
//}
//return claim, nil
//},
//})
