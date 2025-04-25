## How Message Redelivery Actually Works
1. ### Normal Flow:
   - Producer sends message to queue

   - RabbitMQ delivers to a consumer

   - Consumer processes and acknowledges (ack)

    - Message is removed from queue

2. ### When No Acknowledgement Occurs:
   - If consumer disconnects without acknowledging

   - If consumer explicitly rejects (nack) with requeue=true

   - If consumer crashes during processing

   - RabbitMQ will redeliver the message to another consumer (or the same consumer if it reconnects)

### Key Scenarios
     

	
	
msg.Nack(false, false)	Message dropped or sent to DLX (no redelivery).
Consumer crashes without ack	RabbitMQ auto-requeues the message (like Nack(requeue=true)).

| Scenario       | Effect |
|----------------| ----- |
| msg.Ack(false) | 	Message removed from queue (success). |
| msg.Nack(false, true) | Message requeued (will be redelivered).  |
| msg.Nack(false, false)               |Message dropped or sent to DLX (no redelivery).|
| Consumer crashes without ack               |RabbitMQ auto-requeues the message (like Nack(requeue=true)).|