package main

import (
	amqp "github.com/rabbitmq/amqp091-go"
	"log"
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

	body := "Hello World"
	err = channel.Publish(
		"",         // Exchange
		queue.Name, // Routing Key
		false,      // Mandatory -- it can be 'undeliverable' if set to true and no queue match the routing key
		false,      // Immediate -- it can be 'undeliverable' if set to true and no consumer ready to consume the message
		amqp.Publishing{ // Message
			ContentType: "text/plain",
			Body:        []byte(body),
		})
	failOnError(err, "Failed to publish a message")

	log.Printf(" [X] Sent %s\n", body)
}
