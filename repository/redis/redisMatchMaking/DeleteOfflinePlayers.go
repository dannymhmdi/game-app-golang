package redisMatchMaking

import (
	"context"
	"mymodule/entity"
	"strconv"
)

func (r RedisDB) DeleteOfflinePlayers(ctx context.Context, players []entity.WaitingMember) error {
	for _, waitingMember := range players {
		redisKey := strconv.Itoa(int(waitingMember.UserID))
		if _, zErr := r.adaptor.Client().ZRem(ctx, redisKey).Result(); zErr != nil {

			return zErr
		}
	}
	return nil
}
