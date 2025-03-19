package redisMatchMaking

import (
	"context"
	"fmt"
	"github.com/labstack/gommon/log"
	"mymodule/entity"
	"time"
)

func (r RedisDB) DeleteOfflinePlayers(category entity.Category, players []entity.WaitingMember) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	redisKey := fmt.Sprintf("waitingList:%s", category)

	playerIDs := make([]any, 0)

	for _, offlinePlayer := range players {
		playerIDs = append(playerIDs, offlinePlayer.UserID)
	}

	if _, zErr := r.adaptor.Client().ZRem(ctx, redisKey, playerIDs...).Result(); zErr != nil {
		log.Errorf("failed to delete offline in (users-redisMatchMaking.DeleteOfflinePlayers):%v", zErr)
	}

}
