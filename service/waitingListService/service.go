package waitingListService

import (
	"mymodule/dto"
	"mymodule/entity"
	"mymodule/pkg/richerr"
)

type MatchMakingRepositoryService interface {
	Enqueue(userId uint, category entity.Category) error
}

type Service struct {
	repository MatchMakingRepositoryService
}

type Config struct {
}

func New(repo MatchMakingRepositoryService) *Service {
	return &Service{
		repository: repo,
	}
}

func (s Service) AddToWaitingList(req dto.AddToWaitingListRequest) (dto.AddToWaitingListResponse, error) {
	eErr := s.repository.Enqueue(req.UserId, req.Category)
	if eErr != nil {
		return dto.AddToWaitingListResponse{}, richerr.New().
			SetMsg(eErr.Error()).
			SetOperation("matchMakingService.AddToWaitingList").
			SetKind(richerr.KindUnexpected).
			SetWrappedErr(eErr)
	}
	return dto.AddToWaitingListResponse{Message: "successfully add to list"}, nil
}
