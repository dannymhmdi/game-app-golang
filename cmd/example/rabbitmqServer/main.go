package main

import (
	"fmt"
	amqp "github.com/rabbitmq/amqp091-go"
	"log"
	"time"
)

func main() {
	// 1. Connect to RabbitMQ
	conn, dErr := amqp.Dial("amqp://kalo:kalo@localhost:5672/")
	if dErr != nil {
		log.Fatalf("failed to connect to RabbitMQ: %v", dErr)
	}
	defer conn.Close()

	// 2. Create a channel
	ch, cErr := conn.Channel()
	if cErr != nil {
		log.Fatalf("failed to open a channel: %v", cErr)
	}
	defer ch.Close()

	// 3. Declare the same durable queue
	q, qErr := ch.QueueDeclare(
		"task_queue", // name
		true,         // durable
		false,        // delete when unused
		false,        // exclusive
		false,        // no-wait
		nil,          // arguments
	)

	if qErr != nil {
		log.Fatalf("failed to declare a durable queue: %v", qErr)
	}
	// 4. Fair dispatch (1 message per worker)
	err := ch.Qos(
		1,     // prefetch count
		0,     // prefetch size
		false, // global
	)
	if err != nil {
		log.Fatalf("failed to set QoS: %v", err)
	}

	queue, iErr := ch.QueueInspect("task_queue")
	if iErr != nil {
		log.Fatalf("Failed to inspect queue: %v", iErr)
	}

	fmt.Printf("queue data:%+v\n", queue.Messages)

	// 5. Consume messages
	msgs, cErr := ch.Consume(
		q.Name, // queue
		"",     // consumer
		false,  // auto-ack (false = manual ack)
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)

	if cErr != nil {
		log.Fatalf("failed to register a consumer: %v", cErr)
	}
	// 6. Process messages in a goroutine
	forever := make(chan bool)

	go func() {
		for msg := range msgs {
			log.Printf("Received: %s", msg.Body)
			// Simulate work
			time.Sleep(1 * time.Second)
			aErr := msg.Ack(false) // Manual acknowledgment
			if aErr != nil {
				fmt.Printf("failed to ack message: %v\n", aErr)

				continue
			}
			fmt.Printf("acked message")
		}
	}()

	<-forever
}
