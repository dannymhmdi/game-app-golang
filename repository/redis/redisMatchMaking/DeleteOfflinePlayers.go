package redisMatchMaking

import (
	"context"
	"fmt"
	"github.com/labstack/gommon/log"
	"mymodule/entity"
	"strconv"
	"time"
)

func (r RedisDB) DeleteOfflinePlayers(category entity.Category, players []entity.WaitingMember) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	redisKey := fmt.Sprintf("waitingList:%s", category)

	playerIDs := make([]string, 0)

	for _, offlinePlayer := range players {
		userId := strconv.Itoa(int(offlinePlayer.UserID))
		playerIDs = append(playerIDs, userId)
	}

	if _, zErr := r.adaptor.Client().ZRem(ctx, redisKey, playerIDs).Result(); zErr != nil {
		log.Errorf("failed to delete offline in (users-redisMatchMaking.DeleteOfflinePlayers):%v", zErr)
	}

}
