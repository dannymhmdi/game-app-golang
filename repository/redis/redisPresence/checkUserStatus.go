package redisPresence

import (
	"context"
	"fmt"
	"mymodule/entity"
	"strconv"
	"time"
)

func (r RedisDB) CheckUserStatus(ctx context.Context, userIDs []uint) ([]entity.OnlinePlayer, error) {

	onlinePlayers := make([]entity.OnlinePlayer, 0)

	for _, id := range userIDs {
		redisKey := fmt.Sprintf("presence:%d", id)
		res, gErr := r.adaptor.Client().Get(ctx, redisKey).Result()

		if gErr != nil {
			continue
		}

		timeStamp, aErr := strconv.Atoi(res)
		if aErr != nil {
			fmt.Println("failed to convert timestamp to number", aErr)

			continue
		}

		if int64(timeStamp) < time.Now().Add(-30*time.Second).UnixMicro() {
			//r.adaptor.Client().ZRem(ctx,)
			continue
		}

		onlinePlayers = append(onlinePlayers, entity.OnlinePlayer{
			UserId:    id,
			Timestamp: int64(timeStamp),
		})

	}

	return onlinePlayers, nil

}
