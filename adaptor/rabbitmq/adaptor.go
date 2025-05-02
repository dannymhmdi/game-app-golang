package rabbitmq

import (
	"fmt"
	amqp "github.com/rabbitmq/amqp091-go"
	"log"
	"sync"
)

type Adaptor struct {
	rabbitClient *amqp.Connection
	config       Config
	mu           sync.Mutex
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

func New(cfg Config) *Adaptor {
	conn, dErr := amqp.Dial(fmt.Sprintf("amqp://%s:%s@%s:%d/", cfg.Username, cfg.Password, cfg.Host, cfg.Port))
	if dErr != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %v", dErr)
	}
	fmt.Println("Connected to RabbitMQ")
	return &Adaptor{
		rabbitClient: conn,
		config:       cfg,
	}
}

func (a *Adaptor) Publish(event string, payLoad string) {

	fmt.Println("enter match")
	a.mu.Lock()
	ch, cErr := a.rabbitClient.Channel()
	a.mu.Unlock()
	if cErr != nil {
		log.Fatalf("Failed to open a channel: %v", cErr)
	}
	defer func() {
		if cErr := ch.Close(); cErr != nil {
			log.Fatalf("Failed to close a channel: %v", cErr)
		}
	}()

	q, qErr := ch.QueueDeclare(
		"matchedPlayers_queue", // name
		true,                   // durable
		false,                  // delete when unused
		false,                  // exclusive
		false,                  // no-wait
		nil,                    // arguments
	)

	if qErr != nil {
		log.Fatalf("Failed to declare a durable queue: %v", qErr)
	}

	pErr := ch.Publish(
		"",     // exchange (default)
		q.Name, // routing key (queue name)
		false,  // mandatory
		false,  // immediate
		amqp.Publishing{
			DeliveryMode: amqp.Persistent, // Make message persistent
			ContentType:  "text/plain",
			Body:         []byte(payLoad),
		})

	if pErr != nil {
		log.Fatalf("Failed to publish a message: %v", pErr)
	}
	fmt.Printf("Published a message")
}

//"amqp://kalo:kalo@localhost:5672/"
