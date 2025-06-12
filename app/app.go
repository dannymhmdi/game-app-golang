package app

import (
	"mymodule/adaptor/rabbitmq"
	"mymodule/config"
	"mymodule/delivery/httpserver/backOffice_handler"
	"mymodule/delivery/httpserver/matchMaking_handler"
	"mymodule/delivery/httpserver/user_handler"
	"mymodule/repository/mysql"
	"mymodule/service/matchStoreService"
	"mymodule/service/matchmakingService"
)

type App struct {
	UserHandler        *user_handler.Handler
	BackOfficeHandler  *backOffice_handler.Handler
	MatchMakingHandler *matchMaking_handler.Handler
	Services           Services
	Config             config.Config
	DB                 mysql.MysqlDB
	RabbitAdaptor      *rabbitmq.Adaptor
}

type Services struct {
	MatchmakingService matchmakingService.Service
	MatchStoreService  matchStoreService.Service
}
