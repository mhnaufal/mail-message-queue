package main

import (
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

func failOnError(err error, message string) {
	if err != nil {
		log.Fatalf("%s: %s", message, err)
	}
}

func main() {
	// CONNECTION
	// establish a connection to the RabbitMQ
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	// CHANNEL
	// create a channel
	// is a virtual connection (a tiny connection inside TCP) inside the a TCP connection that previously has been created. Used to message the broker
	channel, err := conn.Channel()
	failOnError(err, "Failed to open channel")
	defer channel.Close()

	// declare a queue to be a bucket for accepting the messages
	queue, err := channel.QueueDeclare(
		"Greeting Queue", // Name of the queue/broker
		false,            // Durable
		false,            // Delete when used
		false,            // Exclusive
		false,            // No wait
		nil,              // Arguments
	)
	failOnError(err, "Failed to declare a queue")

	messages, err := channel.Consume(
		queue.Name, // Queue name
		"",         // Consumer
		true,       // Auto-ack
		false,      // Exclusive
		false,      // No-local
		false,      // No-wait
		nil,        // Args
	)
	failOnError(err, "Failed to register the consumer")

	forever := make(chan bool)

	go func() {
		for message := range messages {
			log.Printf("[O] Received message: %s", message.Body)
		}
	}()

	log.Printf(" [*] Waiting for message...To exit, press Ctrl + C")
	<-forever
}
