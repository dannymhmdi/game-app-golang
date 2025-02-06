package user_handler

import "github.com/labstack/echo/v4"

func (h Handler) SetRoute(e *echo.Echo) {
	e.POST("users/register", h.userRegisterHandler)
	e.POST("users/login", h.userLoginHandler)
	e.GET("users/profile", h.userProfileHandler)
}
