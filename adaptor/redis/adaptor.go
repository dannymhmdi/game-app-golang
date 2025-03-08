package redis

import (
	"fmt"
	"github.com/redis/go-redis/v9"
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
