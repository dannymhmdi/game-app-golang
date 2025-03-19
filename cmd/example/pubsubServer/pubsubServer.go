package main

import (
	"context"
	"encoding/base64"
	"google.golang.org/protobuf/proto"
	"mymodule/adaptor/redis"
	"mymodule/config"
	"mymodule/contract/golang/matchingPlayer"
	"time"
)

func main() {
	topic := "matchMakingSvc:playerMatch"
	appConfig := config.Load()
	redisAdaptor := redis.New(appConfig.RedisConfig)

	for {
		protoMu := matchingPlayer.MatchedPlayers{
			UserIds:   []uint64{1, 2},
			Category:  "soccer",
			Timestamp: time.Now().UnixMicro(),
		}
		payLoad, mErr := proto.Marshal(&protoMu)
		if mErr != nil {
			panic(mErr)
		}

		encodePayLoadToString := base64.StdEncoding.EncodeToString(payLoad)
		redisAdaptor.Client().Publish(context.Background(), topic, encodePayLoadToString)
		time.Sleep(5 * time.Second)
	}

}
