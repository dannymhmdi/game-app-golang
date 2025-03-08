package httpserver

import (
	"errors"
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"log"
	"mymodule/config"
	"mymodule/delivery/httpserver/backOffice_handler"
	"mymodule/delivery/httpserver/matchMaking_handler"
	"mymodule/delivery/httpserver/user_handler"
	"net/http"
)

type Server struct {
	config             config.Config
	userHandler        user_handler.Handler
	backOfficeHandler  backOffice_handler.Handler
	waitingListHandler matchMaking_handler.Handler
	Router             *echo.Echo
}

func New(cfg config.Config, userHandler user_handler.Handler, backOfficeHandler backOffice_handler.Handler, waitingListHandler matchMaking_handler.Handler) *Server {
	e := echo.New()
	return &Server{
		config:             cfg,
		userHandler:        userHandler,
		backOfficeHandler:  backOfficeHandler,
		waitingListHandler: waitingListHandler,
		Router:             e,
	}
}

func (s Server) Serve() {
	s.Router.Use(middleware.Logger())
	s.Router.Use(middleware.Recover())
	s.Router.GET("/healthcheck", s.healthCheck)
	s.userHandler.SetRoute(s.Router)
	s.backOfficeHandler.SetBackOfficeRoute(s.Router)
	s.waitingListHandler.SetMatchMakingRoute(s.Router)
	if err := s.Router.Start(fmt.Sprintf(":%s", s.config.HttpConfig.Port)); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatalf("failed tostart http server : %v", err)
	} else if errors.Is(err, http.ErrServerClosed) {
		fmt.Println("http server shutdown gracefully")
	}
}

//func (s Server) Serve() *echo.Echo {
//	e := echo.New()
//	e.Use(middleware.Logger())
//	e.Use(middleware.Recover())
//	e.GET("/healthcheck", s.healthCheck)
//	s.userHandler.SetRoute(e)
//	s.backOfficeHandler.SetBackOfficeRoute(e)
//	s.waitingListHandler.SetMatchMakingRoute(e)
//	//go func() {
//	//	e.Logger.Fatal(e.Start(fmt.Sprintf(":%s", s.config.HttpConfig.Port)))
//	//}()
//	go func() {
//		if err := e.Start(fmt.Sprintf(":%s", s.config.HttpConfig.Port)); err != nil && !errors.Is(err, http.ErrServerClosed) {
//			e.Logger.Fatal("Failed to start server:", err) // Fatal for unexpected errors
//		} else if errors.Is(err, http.ErrServerClosed) {
//			fmt.Println("Server shut down gracefully") // Info for graceful shutdown
//		}
//	}()
//	return e
//}
