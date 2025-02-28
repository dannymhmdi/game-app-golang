package httpserver

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"mymodule/config"
	"mymodule/delivery/httpserver/backOffice_handler"
	"mymodule/delivery/httpserver/matchMaking_handler"
	"mymodule/delivery/httpserver/user_handler"
)

type Server struct {
	config             config.Config
	userHandler        user_handler.Handler
	backOfficeHandler  backOffice_handler.Handler
	waitingListHandler matchMaking_handler.Handler
}

func New(cfg config.Config, userHandler user_handler.Handler, backOfficeHandler backOffice_handler.Handler, waitingListHandler matchMaking_handler.Handler) *Server {
	return &Server{
		config:             cfg,
		userHandler:        userHandler,
		backOfficeHandler:  backOfficeHandler,
		waitingListHandler: waitingListHandler,
	}
}

func (s Server) Serve() {
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.GET("/healthcheck", s.healthCheck)
	s.userHandler.SetRoute(e)
	s.backOfficeHandler.SetBackOfficeRoute(e)
	s.waitingListHandler.SetMatchMakingRoute(e)
	e.Logger.Fatal(e.Start(fmt.Sprintf(":%s", s.config.HttpConfig.Port)))
}
