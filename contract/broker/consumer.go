package broker

import (
	"context"
	amqp "github.com/rabbitmq/amqp091-go"
	"mymodule/contract/golang/matchingPlayer"
)

//type Consumer interface {
//	Consume(ctx context.Context,done <-chan bool, deliveryMsg chan<- matchingPlayer.MatchedPlayers) (*amqp.Connection, *amqp.Channel)
//}

type Consumer interface {
	Consume(ctx context.Context, queueName string, done <-chan bool, deliveryMsg chan<- matchingPlayer.MatchedPlayers) (*amqp.Connection, *amqp.Channel)
}
