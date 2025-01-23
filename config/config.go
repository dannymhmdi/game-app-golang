package config

import "mymodule/service/authservice"

type HttpServer struct {
	Port string
}

type Config struct {
	HttpConfig HttpServer
	AuthConfig authservice.Config
}
