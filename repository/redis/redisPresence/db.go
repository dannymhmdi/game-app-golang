package redisPresence

import (
	"mymodule/adaptor/redis"
	"time"
)

type Config struct {
	PresenceKeyExpirationTime time.Duration `koanf:"presence_key_expiration_time"`
}

type RedisDB struct {
	adaptor redis.Adaptor
	config  Config
}

func New(adaptor redis.Adaptor, cfg Config) RedisDB {
	return RedisDB{
		adaptor: adaptor,
		config:  cfg,
	}
}
