package redisMatchMaking

import (
	"context"
	"fmt"
	"mymodule/entity"
)

func (r RedisDB) DeleteMatchedUsers(ctx context.Context, userIds []uint, category entity.Category) error {
	redisKey := fmt.Sprintf("waitingList:%s", category)
	playerIds := make([]string, len(userIds))
	for i, v := range userIds {
		id := fmt.Sprintf("%d", v)
		playerIds[i] = id
	}

	if _, zErr := r.adaptor.Client().ZRem(ctx, redisKey, playerIds).Result(); zErr != nil {
		return zErr
	}

	return nil
}
