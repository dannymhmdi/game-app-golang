package rabbitmq

import (
	"fmt"
	amqp "github.com/rabbitmq/amqp091-go"
	"log"
	"sync"
)

type Adaptor struct {
	rabbitClient *amqp.Connection
	channel      *amqp.Channel
	queue        amqp.Queue
	mu           sync.Mutex
	config       Config
}

type Config struct {
	Username string `koanf:"username"`
	Password string `koanf:"password"`
	Host     string `koanf:"host"`
	Port     uint   `koanf:"port"`
}

func (a *Adaptor) RabbitConn() *amqp.Connection {
	return a.rabbitClient
}

func (a *Adaptor) RabbitChannel() *amqp.Channel {
	return a.channel
}

//func New(cfg Config) *Adaptor {
//	conn, dErr := amqp.Dial(fmt.Sprintf("amqp://%s:%s@%s:%d/", cfg.Username, cfg.Password, cfg.Host, cfg.Port))
//	if dErr != nil {
//		log.Fatalf("failed to connect to RabbitMQ: %v", dErr)
//	}
//	ch, cErr := conn.Channel()
//	if cErr != nil {
//		log.Fatalf("failed to open rabbitmq channel:%+v", cErr)
//	}
//
//	queue, qErr := ch.QueueDeclare(
//		"matchedPlayers_queue", // name
//		true,                   // durable
//		false,                  // delete when unused
//		false,                  // exclusive
//		false,                  // no-wait
//		nil)
//
//	if qErr != nil {
//		log.Fatalf("failed to declare rabbitmq queue:%+v", qErr)
//	}
//	fmt.Println("Connected to RabbitMQ")
//	return &Adaptor{
//		rabbitClient: conn,
//		config:       cfg,
//		channel:      ch,
//		queue:        queue,
//	}
//}

// exchange version :
func New(cfg Config) *Adaptor {
	exchangeName := "match_events_exchange"
	conn, dErr := amqp.Dial(fmt.Sprintf("amqp://%s:%s@%s:%d/", cfg.Username, cfg.Password, cfg.Host, cfg.Port))
	if dErr != nil {
		log.Fatalf("failed to connect to RabbitMQ: %v", dErr)
	}
	ch, cErr := conn.Channel()
	if cErr != nil {
		log.Fatalf("failed to open rabbitmq channel:%+v", cErr)
	}

	eErr := ch.ExchangeDeclare(exchangeName, "fanout", true, false, false, false, nil)
	if eErr != nil {
		log.Fatalf("failed to declare rabbitmq exchange:%+v", eErr)
	}

	queues := []string{"matchedPlayers_queue", "notification_queue"}

	for _, queue := range queues {
		_, qErr := ch.QueueDeclare(
			queue, // name
			true,  // durable
			false, // delete when unused
			false, // exclusive
			false, // no-wait
			nil)

		if qErr != nil {
			log.Fatalf("failed to declare rabbitmq queue:%+v", qErr)
		}

		//We bind each queue to the exchange. Now, when a message is published to the exchange, all queues receive it.
		if qErr := ch.QueueBind(queue, "", exchangeName, false, nil); qErr != nil {
			log.Fatalf("failed to bind rabbitmq queue to exchange:%+v", qErr)
		}
	}

	fmt.Println("Connected to RabbitMQ")
	return &Adaptor{
		rabbitClient: conn,
		config:       cfg,
		channel:      ch,
	}
}
