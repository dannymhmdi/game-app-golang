package rabbitmq

import (
	"fmt"
	amqp "github.com/rabbitmq/amqp091-go"
)

//func (a *Adaptor) Publish(event string, payLoad string) error {
//
//	//defer func() {
//	//	if cErr := a.channel.Close(); cErr != nil {
//	//		log.Fatalf("Failed to close a channel: %v", cErr)
//	//	}
//	//}()
//	a.mu.Lock()
//	defer a.mu.Unlock()
//
//	pErr := a.channel.Publish(
//		"",           // exchange (default)
//		a.queue.Name, // routing key (queue name)
//		false,        // mandatory
//		false,        // immediate
//		amqp.Publishing{
//			DeliveryMode: amqp.Persistent, // Make message persistent
//			ContentType:  "text/plain",
//			Body:         []byte(payLoad),
//		})
//
//	if pErr != nil {
//		return fmt.Errorf("failed to publish a message: %v", pErr)
//	}
//	fmt.Printf("Published a message")
//	return nil
//}

// exchange version:
func (a *Adaptor) Publish(event string, payLoad string, queue string) error {
	exchangeName := "match_events_exchange"
	//defer func() {
	//	if cErr := a.channel.Close(); cErr != nil {
	//		log.Fatalf("Failed to close a channel: %v", cErr)
	//	}
	//}()
	a.mu.Lock()
	defer a.mu.Unlock()

	pErr := a.channel.Publish(
		exchangeName, // exchange (default)
		queue,        // routing key (queue name)
		false,        // mandatory
		false,        // immediate
		amqp.Publishing{
			DeliveryMode: amqp.Persistent, // Make message persistent
			ContentType:  "text/plain",
			Body:         []byte(payLoad),
		})

	if pErr != nil {
		return fmt.Errorf("failed to publish a message: %v", pErr)
	}
	fmt.Printf("Published a message")
	return nil
}
