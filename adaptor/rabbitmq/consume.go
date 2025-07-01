package rabbitmq

import (
	"context"
	"encoding/base64"
	"fmt"
	amqp "github.com/rabbitmq/amqp091-go"
	"google.golang.org/protobuf/proto"
	"log"
	"mymodule/contract/golang/matchingPlayer"
)

//func (a *Adaptor) Consume(ctx context.Context, done <-chan bool, deliveryMsg chan<- matchingPlayer.MatchedPlayers) (*amqp.Connection, *amqp.Channel) {
//
//	// 4. Fair dispatch (1 message per worker)
//	err := a.channel.Qos(
//		10,    // prefetch count
//		0,     // prefetch size
//		false, // global
//	)
//	if err != nil {
//		log.Fatalf("failed to set QoS: %v", err)
//	}
//
//	queue, iErr := a.channel.QueueInspect("matchedPlayers_queue")
//	if iErr != nil {
//		log.Fatalf("Failed to inspect queue: %v", iErr)
//	}
//
//	fmt.Printf("queue data:%+v\n", queue)
//
//	// 5. Consume messages
//	msgs, cErr := a.channel.Consume(
//		a.queue.Name, // queue
//		"",           // consumer
//		false,        // auto-ack (false = manual ack)
//		false,        // exclusive
//		false,        // no-local
//		false,        // no-wait
//		nil,          // args
//	)
//
//	if cErr != nil {
//		log.Fatalf("failed to register a consumer: %v", cErr)
//	}
//	// 6. Process messages in a goroutine
//
//	go func() {
//		for msg := range msgs {
//			decodedPayload, dErr := base64.StdEncoding.DecodeString(string(msg.Body))
//			if dErr != nil {
//				log.Fatalf("failed to decode message: %v", dErr)
//			}
//			var matchedPlayers matchingPlayer.MatchedPlayers
//			uErr := proto.Unmarshal(decodedPayload, &matchedPlayers)
//			if uErr != nil {
//				log.Fatalf("failed to unmarshal message: %v", uErr)
//			}
//			deliveryMsg <- matchedPlayers
//			fmt.Printf("Received:%+v\n", matchedPlayers)
//			// Simulate work
//			//time.Sleep(1 * time.Second)
//			<-done
//			aErr := msg.Ack(false) // Manual acknowledgment
//			if aErr != nil {
//				fmt.Printf("failed to ack message: %v\n", aErr)
//
//				continue
//			}
//
//			fmt.Printf("message acknowledged")
//		}
//	}()
//
//	return a.rabbitClient, a.channel
//}

// exchange version:
func (a *Adaptor) Consume(ctx context.Context, queueName string, done <-chan bool, deliveryMsg chan<- matchingPlayer.MatchedPlayers) (*amqp.Connection, *amqp.Channel) {

	// 4. Fair dispatch (1 message per worker)
	err := a.channel.Qos(
		10,    // prefetch count
		0,     // prefetch size
		false, // global
	)
	if err != nil {
		log.Fatalf("failed to set QoS: %v", err)
	}

	queue, iErr := a.channel.QueueInspect(queueName)
	if iErr != nil {
		log.Fatalf("Failed to inspect queue: %v", iErr)
	}

	fmt.Printf("queue data:%+v\n", queue)

	// 5. Consume messages
	msgs, cErr := a.channel.Consume(
		queueName, // queue
		"",        // consumer
		false,     // auto-ack (false = manual ack)
		false,     // exclusive
		false,     // no-local
		false,     // no-wait
		nil,       // args
	)

	if cErr != nil {
		log.Fatalf("failed to register a consumer: %v", cErr)
	}
	// 6. Process messages in a goroutine

	go func() {
		for msg := range msgs {
			decodedPayload, dErr := base64.StdEncoding.DecodeString(string(msg.Body))
			if dErr != nil {
				log.Fatalf("failed to decode message: %v", dErr)
			}
			var matchedPlayers matchingPlayer.MatchedPlayers
			uErr := proto.Unmarshal(decodedPayload, &matchedPlayers)
			if uErr != nil {
				log.Fatalf("failed to unmarshal message: %v", uErr)
			}
			deliveryMsg <- matchedPlayers
			fmt.Printf("Received:%+v\n", matchedPlayers)
			// Simulate work
			//time.Sleep(1 * time.Second)
			msgSucceed:=<-done
			if msgSucceed {
				aErr := msg.Ack(false) // Manual acknowledgment
				if aErr != nil {
					fmt.Printf("failed to ack message: %v\n", aErr)
	
					continue
				}
				fmt.Printf("message acknowledged")
				continue
			}
			// aErr := msg.Ack(false) // Manual acknowledgment
			// if aErr != nil {
			// 	fmt.Printf("failed to ack message: %v\n", aErr)

			// 	continue
			// }

			fmt.Printf("message not acknowledged")
		}
	}()

	return a.rabbitClient, a.channel
}
