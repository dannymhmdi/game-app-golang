package httpserver

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"mymodule/config"
	"mymodule/delivery/httpserver/user_handler"
)

type Server struct {
	config      config.Config
	userHandler user_handler.Handler
}

func New(cfg config.Config, userHandler user_handler.Handler) *Server {
	return &Server{
		config:      cfg,
		userHandler: userHandler,
	}
}

func (s Server) Serve() {
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.GET("/healthcheck", s.healthCheck)
	s.userHandler.SetRoute(e)
	//e.POST("users/register", s.userRegisterHandler)
	//e.POST("users/login", s.userLoginHandler)
	//e.GET("users/profile", s.userProfileHandler)
	e.Logger.Fatal(e.Start(fmt.Sprintf(":%s", s.config.HttpConfig.Port)))
}
