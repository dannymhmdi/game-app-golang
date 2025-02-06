package main

import (
	_ "github.com/go-sql-driver/mysql"
	"mymodule/config"
	"mymodule/delivery/httpserver"
	"mymodule/delivery/httpserver/user_handler"
	"mymodule/repository/mysql"
	"mymodule/service/authservice"
	"mymodule/service/userservice"
	"mymodule/validator/uservalidator"
	"time"
)

const (
	signingKey             string        = "tokenpass"
	accessTokenExpireTime  time.Duration = time.Hour * 1
	refreshTokenExpireTime time.Duration = time.Hour * 24 * 7
	refreshSubject                       = "rt"
	accessSubject                        = "at"
)

func main() {
	//logFile, sErr := logger.SetUpFile("errors.log")
	//if sErr != nil {
	//	log.Fatal("failed to setup logger file")
	//}
	//defer logFile.Close()
	userHandler := setUp()
	cfg := config.Config{
		HttpConfig: config.HttpServer{Port: "8080"},
		AuthConfig: authservice.Config{
			SigningKey:             signingKey,
			AccessTokenExpireTime:  accessTokenExpireTime,
			RefreshTokenExpireTime: refreshTokenExpireTime,
			RefreshSubject:         refreshSubject,
			AccessSubject:          accessSubject,
		},
		DbConfig: mysql.Config{
			Username: "gameapp",
			Password: "gameappt0lk2o20",
			Host:     "localhost",
			Port:     3308,
			DbName:   "gameapp_db",
		},
	}
	server := httpserver.New(cfg, *userHandler)

	server.Serve()
}

func setUp() *user_handler.Handler {
	cfg := authservice.Config{
		SigningKey:             signingKey,
		AccessTokenExpireTime:  accessTokenExpireTime,
		RefreshTokenExpireTime: refreshTokenExpireTime,
		RefreshSubject:         refreshSubject,
		AccessSubject:          accessSubject,
	}

	dbConfig := mysql.Config{
		Username: "gameapp",
		Password: "gameappt0lk2o20",
		Host:     "localhost",
		Port:     3308,
		DbName:   "gameapp_db",
	}
	authSvc := authservice.New(cfg)
	mysqlRepo := mysql.New(dbConfig)
	validator := uservalidator.New(mysqlRepo)
	userSvc := userservice.New(mysqlRepo, authSvc, *validator)

	userHandler := user_handler.New(*authSvc, *userSvc, *validator)
	return userHandler
}
