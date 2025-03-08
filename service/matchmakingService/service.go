package matchmakingService

import (
	"fmt"
	"mymodule/entity"
	"mymodule/params"
	"mymodule/pkg/richerr"
	"time"
)

type MatchMakingRepositoryService interface {
	Enqueue(userId uint, category entity.Category) error
}

type Service struct {
	repository MatchMakingRepositoryService
	config     Config
}

type Config struct {
	Timeout time.Time
}

func New(repo MatchMakingRepositoryService, cfg Config) *Service {
	return &Service{
		repository: repo,
		config:     cfg,
	}
}

func (s Service) AddToWaitingList(req params.AddToWaitingListRequest) (params.AddToWaitingListResponse, error) {
	eErr := s.repository.Enqueue(req.UserId, req.Category)
	if eErr != nil {
		fmt.Println("kiri", eErr)
		return params.AddToWaitingListResponse{}, richerr.New().
			SetMsg(eErr.Error()).
			SetOperation("matchMakingService.AddToWaitingList").
			SetKind(richerr.KindUnexpected).
			SetWrappedErr(eErr)
	}
	return params.AddToWaitingListResponse{Message: "successfully add to list", Timeout: s.config.Timeout}, nil
}

func (s Service) MatchMaking() error {

	fmt.Println("matchmakin service run", time.Now())
	return fmt.Errorf("matchmaking service error")
}
