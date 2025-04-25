package main

import (
	"encoding/json"
	"fmt"
	amqp "github.com/rabbitmq/amqp091-go"
	"log"
	"mymodule/entity"
	"time"
)

func main() {
	fmt.Println("kiri")
	// 1. Connect to RabbitMQ
	conn, dErr := amqp.Dial("amqp://kalo:kalo@localhost:5672/")
	if dErr != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %v", dErr)
	}

	defer conn.Close()

	// 2. Create a channel
	ch, cErr := conn.Channel()
	if cErr != nil {
		log.Fatalf("Failed to create a channel: %v", cErr)
	}
	defer ch.Close()

	// 3. Declare a durable queue (survives server restarts)
	q, qErr := ch.QueueDeclare(
		"task_queue", // name
		true,         // durable
		false,        // delete when unused
		false,        // exclusive
		false,        // no-wait
		nil,          // arguments
	)
	if qErr != nil {
		log.Fatalf("Failed to declare a durable queue: %v", qErr)
	}

	bd, mErr := json.Marshal(entity.MatchedPlayers{
		UserIDs:   []uint{7, 2},
		Category:  "soccer",
		Timestamp: time.Now().UnixMicro(),
	})

	if mErr != nil {
		panic(mErr)
	}
	// 5. Publish a message
	for {
		err := ch.Publish(
			"",     // exchange (default)
			q.Name, // routing key (queue name)
			false,  // mandatory
			false,  // immediate
			amqp.Publishing{
				DeliveryMode: amqp.Persistent, // Make message persistent
				ContentType:  "application/json",
				Body:         bd,
			})

		if err != nil {
			log.Fatalf("Failed to publish a message :%v", err)
		}
		fmt.Println("message published")
		time.Sleep(3 * time.Second)
	}

}
