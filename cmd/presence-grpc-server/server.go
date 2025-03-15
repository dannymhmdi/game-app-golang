package main

import (
	"mymodule/adaptor/redis"
	"mymodule/config"
	"mymodule/delivery/grpcserver/presenceserver"
	"mymodule/repository/redis/redisPresence"
	"mymodule/service/presenceService"
)

func main() {
	appConfig := config.Load()
	redisAdaptor := redis.New(appConfig.RedisConfig)
	presenceRepo := redisPresence.New(redisAdaptor, appConfig.RedisPresence)
	presenceSVc := presenceService.New(presenceRepo)
	presenceServer := presenceserver.New(*presenceSVc)
	presenceServer.Start()
}
