package main

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	amqp "github.com/rabbitmq/amqp091-go"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"mymodule/adaptor/presence"
	"mymodule/adaptor/rabbitmq"
	"mymodule/adaptor/redis"
	"mymodule/config"
	"mymodule/delivery/httpserver"
	"mymodule/delivery/httpserver/backOffice_handler"
	"mymodule/delivery/httpserver/matchMaking_handler"
	"mymodule/delivery/httpserver/user_handler"
	"mymodule/logger"
	"mymodule/params"
	"mymodule/repository/mysql"
	"mymodule/repository/mysql/mysqlAccessControl"
	"mymodule/repository/mysql/mysqlAuth"
	"mymodule/repository/mysql/mysqlMatchStore"
	"mymodule/repository/mysql/mysqlUser"
	"mymodule/repository/redis/redisMatchMaking"
	"mymodule/repository/redis/redisPresence"
	"mymodule/scheduler"
	"mymodule/service/authService"
	"mymodule/service/authorizationService"
	"mymodule/service/backofficeService"
	"mymodule/service/matchStoreService"
	"mymodule/service/matchmakingService"
	"mymodule/service/presenceService"
	"mymodule/service/userService"
	"mymodule/validator/authValidator"
	"mymodule/validator/matchMakingValidator"
	"mymodule/validator/uservalidator"
	"os"
	"os/signal"
	"time"
)

func main() {
	loggerCfg := logger.Config{
		Production: false, // or true for production
		LogFile:    "logger/logfiles/app.log",
		MaxSize:    1,    // MB
		MaxBackups: 5,    // files
		MaxAge:     30,   // days
		Compress:   true, // compress rotated files
	}
	logger.Init(loggerCfg)

	defer logger.Sync()
	conn, dErr := grpc.Dial(":8086", grpc.WithInsecure())
	if dErr != nil {
		panic(dErr)
	}

	defer conn.Close()
	userHandler, backOfficeHandler, matchMakingHandler, matchmakingSvc, matchStoreSvc, appConfig, db := setUp(conn)
	defer db.NewConn().Close()
	server := httpserver.New(appConfig, *userHandler, *backOfficeHandler, *matchMakingHandler)
	done := make(chan bool)
	quit := make(chan os.Signal)
	rabbitConn := make(chan *amqp.Connection, 1)
	rabbitCh := make(chan *amqp.Channel, 1)
	schedulerOp := scheduler.New(matchmakingSvc)
	signal.Notify(quit, os.Interrupt)
	go func() {
		schedulerOp.Start(done)
	}()

	go func() {
		server.Serve()
	}()

	go func() {
		conn, ch := matchStoreSvc.StoreMatch(context.Background(), params.MatchStoreRequest{})
		rabbitConn <- conn
		rabbitCh <- ch
	}()

	rabbitConnection := <-rabbitConn
	rabbitChannel := <-rabbitCh
	defer rabbitConnection.Close()
	defer rabbitChannel.Close()
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

func setUp(conn *grpc.ClientConn) (*user_handler.Handler, *backOffice_handler.Handler, *matchMaking_handler.Handler, matchmakingService.Service, matchStoreService.Service, config.Config, mysql.MysqlDB) {
	appConfig := config.Load()

	mysqlDB := mysql.New(appConfig.DbConfig)
	authRepo := mysqlAuth.New(*mysqlDB)
	authSvc := authService.New(appConfig.AuthConfig, authRepo)
	userRepo := mysqlUser.New(mysqlDB)
	validator := uservalidator.New(userRepo)
	authenticationValidator := authValidator.New(authRepo)
	userSvc := userService.New(userRepo, authSvc, *validator)
	authorizationRepo := mysqlAccessControl.New(mysqlDB)
	authorizationSvc := authorizationService.New(authorizationRepo)
	backOfficeSvc := backofficeService.New()
	backOfficeHandler := backOffice_handler.New(*backOfficeSvc, *authSvc, *authorizationSvc)
	matchMakerValidator := matchMakingValidator.New()
	redisAdaptor := redis.New(appConfig.RedisConfig)
	matchMakingRepo := redisMatchMaking.New(redisAdaptor)
	presenceRepo := redisPresence.New(redisAdaptor, appConfig.RedisPresence)
	presenceSvc := presenceService.New(presenceRepo)

	presenceAdaptor := presence.New(conn)
	rabbitAdaptor := rabbitmq.New(appConfig.RabbitMqConfig)
	matchMakingSvc := matchmakingService.New(matchMakingRepo, *presenceAdaptor, redisAdaptor, rabbitAdaptor, appConfig.MatchMakingConfig)
	waitingListHandler := matchMaking_handler.New(*matchMakingSvc, *authSvc, []byte(appConfig.AuthConfig.SigningKey), *matchMakerValidator)

	matchStoreRepo := mysqlMatchStore.New(*mysqlDB)
	matchStoreSvc := matchStoreService.New(matchStoreRepo, rabbitAdaptor)
	userHandler := user_handler.New(*authSvc, *userSvc, *presenceSvc, *validator, authenticationValidator, []byte(appConfig.AuthConfig.SigningKey))
	return userHandler, backOfficeHandler, waitingListHandler, *matchMakingSvc, *matchStoreSvc, appConfig, *mysqlDB
}

//todo
// refresh-token implementaion
//otp process when game created to matched users
