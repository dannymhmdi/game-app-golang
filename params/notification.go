package params

import (
amqp "github.com/rabbitmq/amqp091-go"
)

type NotficationRequest struct {

}

type NotficationResponse struct {
RabbitConnection *amqp.Connection
RabbitChannel *amqp.Channel
}