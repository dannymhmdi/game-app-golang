package presenceService

import (
	"context"
	"fmt"
	"mymodule/entity"
	"mymodule/params"
	"mymodule/pkg/richerr"
	"time"
)

type Service struct {
	repository PresenceRepositoryService
}

type PresenceRepositoryService interface {
	UpsertUserStatus(ctx context.Context, key string, timeStamp int64) error
	CheckUserStatus(ctx context.Context, userIDs []uint) ([]entity.OnlinePlayer, error)
}

func New(repo PresenceRepositoryService) *Service {
	return &Service{
		repository: repo,
	}
}

func (s Service) Presence(ctx context.Context, req params.PresenseRequest) (params.PresenseResponse, error) {
	select {
	case <-ctx.Done():
		return params.PresenseResponse{}, richerr.New().
			SetMsg(ctx.Err().Error()).
			SetOperation("presenceService.Presence").
			SetKind(richerr.KindResponseTimeout)
	default:
		redisKey := fmt.Sprintf("presence:%d", req.UserId)
		uErr := s.repository.UpsertUserStatus(ctx, redisKey, time.Now().UnixMicro())
		if uErr != nil {
			return params.PresenseResponse{}, richerr.New().
				SetMsg(uErr.Error()).
				SetOperation("presenceService.Presence").
				SetWrappedErr(uErr)
			//SetKind()
		}

		return params.PresenseResponse{Message: "user presence upsert"}, nil

	}
	//redisKey := fmt.Sprintf("presence:%d", req.UserId)
	//uErr := s.repository.UpsertUserStatus(ctx, redisKey, time.Now().UnixMicro())
	//if uErr != nil {
	//	return params.PresenseResponse{}, richerr.New().
	//		SetMsg(uErr.Error()).
	//		SetOperation("presenceService.Presence").
	//		SetWrappedErr(uErr)
	//	//SetKind()
	//}
	//
	//return params.PresenseResponse{Message: "user presence upsert"}, nil
}

func (s Service) GetPresence(ctx context.Context, req params.GetPresenceRequest) (params.GetPresenceResponse, error) {
	onlinePlayers, cErr := s.repository.CheckUserStatus(ctx, req.UserIDs)
	if cErr != nil {
		return params.GetPresenceResponse{}, richerr.New().
			SetMsg(cErr.Error()).
			SetOperation("presenceService.GetPresence").
			SetKind(richerr.KindUnexpected)
	}

	return params.GetPresenceResponse{
		OnlinePlayers: onlinePlayers,
	}, nil

	//return params.GetPresenceResponse{
	//	OnlinePlayers: []entity.OnlinePlayer{
	//		{UserId: 1, Timestamp: 1742021804},
	//		{UserId: 2, Timestamp: 1742021804},
	//		{UserId: 3, Timestamp: 1742021804},
	//	},
	//}, nil

}
