package matchStoreService

import (
	"context"
	"fmt"
	amqp "github.com/rabbitmq/amqp091-go"
	"mymodule/contract/broker"
	"mymodule/contract/golang/matchingPlayer"
	"mymodule/entity"
	"mymodule/params"
	"mymodule/pkg/slice"
	"time"
)

type Service struct {
	repository  MatchStoreRepositoryService
	msgConsumer broker.Consumer
}

type MatchStoreRepositoryService interface {
	StoreMatch(ctx context.Context, game entity.Game) (uint, error)
}

func New(repo MatchStoreRepositoryService, consumer broker.Consumer) *Service {
	return &Service{
		repository:  repo,
		msgConsumer: consumer,
	}
}

func (m Service) StoreMatch(ctx context.Context, req params.MatchStoreRequest) (*amqp.Connection, *amqp.Channel) {
	done := make(chan bool)
	deliveryMsg := make(chan matchingPlayer.MatchedPlayers)
	conn, ch := m.msgConsumer.Consume(ctx, "matchedPlayers_queue", done, deliveryMsg)
	go func() {
		for msg := range deliveryMsg {
			game := entity.Game{
				ID:        0,
				Category:  entity.Category(msg.Category),
				PlayersID: slice.Uint64ToUintMapper(msg.UserIds),
				StartTime: time.Now(),
			}
			matchId, sErr := m.repository.StoreMatch(ctx, game)
			if sErr != nil {
				fmt.Println("matchStoreService.StoreMatch:", sErr)
				continue
			}
			fmt.Println("match created with id:", matchId)
			done <- true
		}
	}()

	return conn, ch
}
