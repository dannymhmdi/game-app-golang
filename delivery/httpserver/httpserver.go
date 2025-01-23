package httpserver

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"mymodule/config"
	"mymodule/service/authservice"
	"mymodule/service/registerservice"
)

type Server struct {
	config  config.Config
	authSvc authservice.Service
	userSvc registerservice.Service
}

func New(cfg config.Config, authSvc authservice.Service, userSvc registerservice.Service) *Server {
	return &Server{
		config:  cfg,
		authSvc: authSvc,
		userSvc: userSvc,
	}
}

func (s Server) Serve() {
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.GET("/healthcheck", s.healthCheck)
	e.POST("users/register", s.userRegisterHandler)
	e.POST("users/login", s.userLoginHandler)
	e.GET("users/profile", s.userProfileHandler)
	e.Logger.Fatal(e.Start(fmt.Sprintf(":%s", s.config.HttpConfig.Port)))
}
