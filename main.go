package main

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"golang.org/x/net/context"
	"mymodule/adaptor/redis"
	"mymodule/config"
	"mymodule/delivery/httpserver"
	"mymodule/delivery/httpserver/backOffice_handler"
	"mymodule/delivery/httpserver/matchMaking_handler"
	"mymodule/delivery/httpserver/user_handler"
	"mymodule/repository/mysql"
	"mymodule/repository/mysql/mysqlAccessControl"
	"mymodule/repository/mysql/mysqlUser"
	"mymodule/repository/redisMatchMaking"
	"mymodule/scheduler"
	"mymodule/service/authorizationService"
	"mymodule/service/authservice"
	"mymodule/service/backoffice"
	"mymodule/service/matchmakingService"
	"mymodule/service/userservice"
	"mymodule/validator/matchMakingValidator"
	"mymodule/validator/uservalidator"
	"os"
	"os/signal"
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
	config.Load()
	userHandler, backOfficeHandler, matchMakingHandler, matchmakingSvc := setUp()
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

	server := httpserver.New(cfg, *userHandler, *backOfficeHandler, *matchMakingHandler)
	done := make(chan bool)
	quit := make(chan os.Signal)

	schedulerOp := scheduler.New(matchmakingSvc)
	signal.Notify(quit, os.Interrupt)
	go func() {
		schedulerOp.Start(done)
	}()

	go func() {
		server.Serve()
	}()

	<-quit
	done <- true
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if sErr := server.Router.Shutdown(ctx); sErr != nil {
		fmt.Println(sErr)
	}
	<-ctx.Done()

	//time.Sleep(6 * time.Second)
	fmt.Println("app terminated gracefully")

}

func setUp() (*user_handler.Handler, *backOffice_handler.Handler, *matchMaking_handler.Handler, matchmakingService.Service) {
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

	redisAdaptorCfg := redis.Config{
		Host: "localhost",
		Port: 6380,
	}

	matchMakingCfg := matchmakingService.Config{
		Timeout: time.Now(),
	}
	authSvc := authservice.New(cfg)
	mysqlDB := mysql.New(dbConfig)
	userRepo := mysqlUser.New(mysqlDB)
	validator := uservalidator.New(userRepo)
	userSvc := userservice.New(userRepo, authSvc, *validator)
	authorizationRepo := mysqlAccessControl.New(mysqlDB)
	authorizationSvc := authorizationService.New(authorizationRepo)
	backOfficeSvc := backoffice.New()
	userHandler := user_handler.New(*authSvc, *userSvc, *validator, []byte(signingKey))
	backOfficeHandler := backOffice_handler.New(*backOfficeSvc, *authSvc, *authorizationSvc)
	matchMakerValidator := matchMakingValidator.New()
	redisAdaptor := redis.New(redisAdaptorCfg)
	matchMakingRepo := redisMatchMaking.New(redisAdaptor)
	matchMakingSvc := matchmakingService.New(matchMakingRepo, matchMakingCfg)
	waitingListHandler := matchMaking_handler.New(*matchMakingSvc, *authSvc, []byte(signingKey), *matchMakerValidator)
	return userHandler, backOfficeHandler, waitingListHandler, *matchMakingSvc
}
