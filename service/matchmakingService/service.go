package matchmakingService

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"mymodule/entity"
	"mymodule/params"
	"mymodule/pkg/richerr"
	"mymodule/pkg/timestamp"
	"sync"
	"time"
)

type MatchMakingRepositoryService interface {
	Enqueue(userId uint, category entity.Category) error
	GetCategoryWaitingList(ctx context.Context, category entity.Category) ([]entity.WaitingMember, error)
	DeleteOfflinePlayers(category entity.Category, players []entity.WaitingMember)
}

// we implement this interface in Presence service because we need to check Presence service database but define
// interface in matchmaking service in general we implement our client requierments via interface
type GetPresenceClient interface {
	GetPresence(ctx context.Context, request params.GetPresenceRequest) (params.GetPresenceResponse, error)
}

type MsgPublisher interface {
	PublishMsgToPubSub(ctx context.Context, mu entity.MatchedPlayers)
}

type Service struct {
	repository        MatchMakingRepositoryService
	getPresenceClient GetPresenceClient
	config            Config
	client            *redis.Client
	publisher         MsgPublisher
}

type Config struct {
	Timeout time.Time
}

func New(repo MatchMakingRepositoryService, getPresenceClient GetPresenceClient, publisher MsgPublisher, cfg Config) *Service {
	return &Service{
		repository:        repo,
		getPresenceClient: getPresenceClient,
		config:            cfg,
		publisher:         publisher,
	}
}

func (s Service) AddToWaitingList(req params.AddToWaitingListRequest) (params.AddToWaitingListResponse, error) {
	eErr := s.repository.Enqueue(req.UserId, req.Category)
	if eErr != nil {
		return params.AddToWaitingListResponse{}, richerr.New().
			SetMsg(eErr.Error()).
			SetOperation("matchMakingService.AddToWaitingList").
			SetKind(richerr.KindUnexpected).
			SetWrappedErr(eErr)
	}
	return params.AddToWaitingListResponse{Message: "successfully add to list", Timeout: s.config.Timeout}, nil
}

func (s Service) MatchMaking(ctx context.Context, categories []entity.Category, req params.MatchMakingRequest) (params.MatchMakingResponse, error) {
	var wg sync.WaitGroup
	for _, category := range categories {
		wg.Add(1)

		go s.MatchMaker(ctx, category, &wg)
	}
	wg.Wait()

	return params.MatchMakingResponse{}, nil
}

func (s Service) MatchMaker(ctx context.Context, category entity.Category, wg *sync.WaitGroup) (params.MatchMakingResponse, error) {
	defer wg.Done()
	waitingList, gErr := s.repository.GetCategoryWaitingList(ctx, category)
	if gErr != nil {
		fmt.Printf("%s category is empty\n", category)
		return params.MatchMakingResponse{}, richerr.New().
			SetWrappedErr(gErr).
			SetOperation("matchMakingService.MatchMaking").
			SetMsg(gErr.Error())
	}

	waitingListIDs := make([]uint, 0)

	for _, waitingMember := range waitingList {
		waitingListIDs = append(waitingListIDs, waitingMember.UserID)
	}

	response, pErr := s.getPresenceClient.GetPresence(ctx, params.GetPresenceRequest{
		UserIDs: waitingListIDs,
	})

	if pErr != nil {
		return params.MatchMakingResponse{}, richerr.New().
			SetWrappedErr(pErr).
			SetOperation("matchMakingService.MatchMaking").
			SetKind(richerr.KindUnexpected).
			SetMsg(pErr.Error())
	}

	//delete offline players
	playersToDelete := make([]entity.WaitingMember, 0)

	for _, waiter := range waitingList {
		isExist := false
		for _, onlinePlayer := range response.OnlinePlayers {
			if waiter.UserID == onlinePlayer.UserId {
				isExist = true

				break
			}
		}

		if !isExist {
			playersToDelete = append(playersToDelete, waiter)
		}
	}

	go s.repository.DeleteOfflinePlayers(category, playersToDelete)

	matchedUsers := make([]entity.MatchedPlayers, 0)

	for i := 0; i < len(response.OnlinePlayers)-1; i = i + 2 {
		//todo check can matching players concurrently

		matchedIDs := []uint{response.OnlinePlayers[i].UserId, response.OnlinePlayers[i+1].UserId}

		mu := entity.MatchedPlayers{
			UserIDs:   matchedIDs,
			Category:  category,
			Timestamp: timestamp.Time(),
		}

		//todo matchedUsers should be deleted
		matchedUsers = append(matchedUsers, mu)

		go s.publisher.PublishMsgToPubSub(ctx, mu)
		//save created game for matched users ids in database & remove matched ids from zset
	}

	return params.MatchMakingResponse{MatchedUsers: matchedUsers}, nil
}
