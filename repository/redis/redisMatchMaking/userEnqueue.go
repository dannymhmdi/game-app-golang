package redisMatchMaking

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"mymodule/entity"
	"mymodule/pkg/richerr"
	"time"
)

func (r RedisDB) Enqueue(userId uint, category entity.Category) error {
	zsetKey := fmt.Sprintf("waitingList:%s", category)

	ctx := context.Background()
	_, zErr := r.adaptor.Client().ZAdd(ctx, zsetKey, redis.Z{
		Score:  float64(time.Now().UnixMicro()),
		Member: userId,
	}).Result()

	if zErr != nil {
		return richerr.New().
			SetMsg(zErr.Error()).
			SetOperation("redisMatchMaking.Enqueue").
			SetKind(richerr.KindUnexpected)
	}

	return nil

}
