### In RabbitMQ when we define exchange is it neccessary to define queue for each time publish a message or it automatically publish message to queues which has defined in exchange?

```golang
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
```