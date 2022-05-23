package main

import (
	"log"
	"os"
	"strings"

	amqp "github.com/rabbitmq/amqp091-go"
)

func failOnError(err error, message string) {
	if err != nil {
		log.Fatalf("%s: %s", message, err)
	}
}

func main() {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	failOnError(err, "Failed to establish connection to the RabbitMQ server!")
	defer conn.Close() // get executed when the parent function return

	channel, err := conn.Channel()
	failOnError(err, "Failed to create a channel!")
	defer channel.Close()

	queue, err := channel.QueueDeclare(
		"durable_work_queue", // Name of the queue/broker
		true,         // Durable
		false,        // Delete when used
		false,        // Exclusive
		false,        // No wait
		nil,          // Arguments
	)
	failOnError(err, "Failed to declare a queue!")

	body := bodyFrom(os.Args)
	err = channel.Publish(
		"",         // Exchange
		queue.Name, // Routing key
		false,      // Mandatory
		false,      // Immediate
		amqp.Publishing{
			DeliveryMode: amqp.Persistent,
			ContentType:  "text/plain",
			Body:         []byte(body),
		})
	failOnError(err, "Failed to publish the message!")

	log.Printf(" [X] Sent %s\n", body)
}

func bodyFrom(args []string) string {
	var s string

	if (len(args) < 2) || os.Args[1] == "" {
		s = `hello...halo...hola...`
	} else {
		s = strings.Join(args[1:], " ")
	}

	return s
}
