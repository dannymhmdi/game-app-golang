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
	"mymodule/repository/redis/redisMatchMaking"
	"mymodule/repository/redis/redisPresence"
	"mymodule/scheduler"
	"mymodule/service/authService"
	"mymodule/service/authorizationService"
	"mymodule/service/backofficeService"
	"mymodule/service/matchmakingService"
	"mymodule/service/presenceService"
	"mymodule/service/userService"
	"mymodule/validator/matchMakingValidator"
	"mymodule/validator/uservalidator"
	"os"
	"os/signal"
	"time"
)

func main() {
	//logFile, sErr := logger.SetUpFile("errors.log")
	//if sErr != nil {
	//	log.Fatal("failed to setup logger file")
	//}
	//defer logFile.Close()
	//config.Load()
	userHandler, backOfficeHandler, matchMakingHandler, matchmakingSvc, appConfig := setUp()

	server := httpserver.New(appConfig, *userHandler, *backOfficeHandler, *matchMakingHandler)
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

	fmt.Println("app terminated gracefully")

}

func setUp() (*user_handler.Handler, *backOffice_handler.Handler, *matchMaking_handler.Handler, matchmakingService.Service, config.Config) {
	appConfig := config.Load()

	authSvc := authService.New(appConfig.AuthConfig)
	mysqlDB := mysql.New(appConfig.DbConfig)
	userRepo := mysqlUser.New(mysqlDB)
	validator := uservalidator.New(userRepo)
	userSvc := userService.New(userRepo, authSvc, *validator)
	authorizationRepo := mysqlAccessControl.New(mysqlDB)
	authorizationSvc := authorizationService.New(authorizationRepo)
	backOfficeSvc := backofficeService.New()
	backOfficeHandler := backOffice_handler.New(*backOfficeSvc, *authSvc, *authorizationSvc)
	matchMakerValidator := matchMakingValidator.New()
	redisAdaptor := redis.New(appConfig.RedisConfig)
	matchMakingRepo := redisMatchMaking.New(redisAdaptor)
	matchMakingSvc := matchmakingService.New(matchMakingRepo, appConfig.MatchMakingConfig)
	waitingListHandler := matchMaking_handler.New(*matchMakingSvc, *authSvc, []byte(appConfig.AuthConfig.SigningKey), *matchMakerValidator)
	presenceRepo := redisPresence.New(redisAdaptor, appConfig.RedisPresence)
	presenceSvc := presenceService.New(presenceRepo)
	userHandler := user_handler.New(*authSvc, *userSvc, *presenceSvc, *validator, []byte(appConfig.AuthConfig.SigningKey))
	return userHandler, backOfficeHandler, waitingListHandler, *matchMakingSvc, appConfig
}
