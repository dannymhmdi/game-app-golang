package presenceService

import (
	"context"
	"fmt"
	"mymodule/params"
	"mymodule/pkg/richerr"
	"time"
)

type Service struct {
	repository PresenceRepositoryService
}

type PresenceRepositoryService interface {
	UpsertUserStatus(ctx context.Context, userID uint, key string, timeStamp int64) error
}

func New(repo PresenceRepositoryService) *Service {
	return &Service{
		repository: repo,
	}
}

func (s Service) Presence(ctx context.Context, req params.PresenseRequest) (params.PresenseResponse, error) {
	redisKey := fmt.Sprintf("presence:%d", req.UserId)
	uErr := s.repository.UpsertUserStatus(ctx, req.UserId, redisKey, time.Now().UnixMilli())
	if uErr != nil {
		return params.PresenseResponse{}, richerr.New().
			SetMsg(uErr.Error()).
			SetOperation("presenceService.Presence").
			SetWrappedErr(uErr)
		//SetKind()
	}

	return params.PresenseResponse{Message: "user presence upsert"}, nil
}
