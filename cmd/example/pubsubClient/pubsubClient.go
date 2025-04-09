package main

import (
	"context"
	"encoding/base64"
	"fmt"
	"google.golang.org/protobuf/proto"
	"mymodule/adaptor/redis"
	"mymodule/config"
	"mymodule/contract/golang/matchingPlayer"
)

func main() {
	topic := "matchMakingSvc:playerMatch"
	appConfig := config.Load()
	redisAdaptor := redis.New(appConfig.RedisConfig)
	subscriber := redisAdaptor.Client().Subscribe(context.Background(), topic)
	for {
		msg, rErr := subscriber.ReceiveMessage(context.Background())
		if rErr != nil {
			panic(rErr)
		}

		decodedPayLoad, dErr := base64.StdEncoding.DecodeString(msg.Payload)
		if dErr != nil {
			panic(dErr)
		}

		var p matchingPlayer.MatchedPlayers
		if uErr := proto.Unmarshal(decodedPayLoad, &p); uErr != nil {
			panic(uErr)
		}
		fmt.Printf("pubsub:%+v\n", p)
	}
}
