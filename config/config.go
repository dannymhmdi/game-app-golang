package config

import (
	"mymodule/repository/mysql"
	"mymodule/service/authservice"
)

type HttpServer struct {
	Port string
}

type Config struct {
	HttpConfig HttpServer
	AuthConfig authservice.Config
	DbConfig   mysql.Config
}
