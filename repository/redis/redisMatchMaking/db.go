package redisMatchMaking

import (
	"mymodule/adaptor/redis"
)

type RedisDB struct {
	adaptor redis.Adaptor
}

func New(adaptor redis.Adaptor) RedisDB {
	return RedisDB{
		adaptor: adaptor,
	}
}
