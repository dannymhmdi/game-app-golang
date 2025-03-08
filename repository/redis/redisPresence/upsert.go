package redisPresence

import (
	"context"
	"fmt"
	"mymodule/pkg/richerr"
)

func (r RedisDB) UpsertUserStatus(ctx context.Context, userID uint, key string, timeStamp int64) error {

	_, sErr := r.adaptor.Client().Set(ctx, key, timeStamp, r.config.PresenceKeyExpirationTime).Result()
	if sErr != nil {
		fmt.Println("presence error :", sErr)
		return richerr.New().SetMsg(sErr.Error()).
			SetOperation("redisPresence.UpsertUserStatus").
			SetKind(richerr.KindUnexpected).
			SetWrappedErr(sErr)
	}
	return nil
}
