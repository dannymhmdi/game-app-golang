package redisMatchMaking

import (
	"context"
	"fmt"
	"mymodule/entity"
	"mymodule/pkg/richerr"
	"strconv"
)

const zsetName = "waitingList:soccer"

func (r RedisDB) GetCategoryWaitingList(ctx context.Context, category entity.Category) ([]entity.WaitingMember, error) {

	list, zErr := r.adaptor.Client().ZRangeWithScores(ctx, zsetName, 0, -1).Result()
	if len(list) == 0 {
		return nil, richerr.New().
			SetKind(richerr.KindUnexpected).
			SetOperation("redisMatchMaking.getCategoryWaitingList").
			SetMsg(fmt.Sprintf("%s category is empty", category))
	}
	if zErr != nil {
		fmt.Println("zErr", zErr)
		return nil, richerr.New().
			SetKind(richerr.KindUnexpected).
			SetOperation("redisMatchMaking.getCategoryWaitingList").
			SetMsg(zErr.Error())
	}

	waitingList := make([]entity.WaitingMember, 0)

	for _, user := range list {
		userID, aErr := strconv.Atoi(user.Member.(string))
		if aErr != nil {
			return nil, richerr.New().
				SetKind(richerr.KindUnexpected).
				SetOperation("redisMatchMaking.getCategoryWaitingList").
				SetMsg(aErr.Error())
		}

		waitingList = append(waitingList, entity.WaitingMember{
			UserID:    uint(userID),
			Timestamp: int64(user.Score),
			Category:  category,
		})

	}

	return waitingList, nil
}
