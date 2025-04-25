package redis

import (
	"context"
	"encoding/base64"
	"fmt"
	"github.com/redis/go-redis/v9"
	"google.golang.org/protobuf/proto"
	"mymodule/contract/golang/matchingPlayer"
	"mymodule/entity"
	"mymodule/pkg/slice"
	"time"
)

type Adaptor struct {
	client *redis.Client
	config Config
}

type Config struct {
	Addr string `koanf:"addr"`
	DB   int    `koanf:"db"`
	Port uint   `koanf:"port"`
	Host string `koanf:"host"`
}

func (a Adaptor) Client() *redis.Client {
	return a.client
}

func New(config Config) Adaptor {
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", config.Host, config.Port),
		Password: "", // No password set
		DB:       0,  // Use default DB
	})
	return Adaptor{
		client: client,
		config: config,
	}

}

func (a Adaptor) PublishMsgToPubSub(ctx context.Context, mu entity.MatchedPlayers) {
	topic := "matchMakingSvc:playerMatch"
	protoMu := matchingPlayer.MatchedPlayers{
		UserIds:   slice.UintToUint64Mapper(mu.UserIDs),
		Category:  string(mu.Category),
		Timestamp: mu.Timestamp,
	}

	payLoad, mErr := proto.Marshal(&protoMu)
	if mErr != nil {
		return
	}

	payloadToString := base64.StdEncoding.EncodeToString(payLoad)
	a.client.Subscribe(ctx, topic, payloadToString)
}

func (a Adaptor) Publish(event string, payLoad string) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	a.client.Publish(ctx, event, payLoad)
}
