package main

import (
	"bytes"
	"log"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

func failOnError(err error, message string) {
	if err != nil {
		log.Fatalf("%s: %s", message, err)
	}
}

func main() {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	failOnError(err, "Failed to establish a connection to RabbitMQ!")
	defer conn.Close()

	channel, err := conn.Channel()
	failOnError(err, "Failed to open a channel!")
	defer channel.Close()

	queue, err := channel.QueueDeclare(
		"durable_work_queue", // Name of the queue/broker
		true,         // Durable
		false,        // Delete when used
		false,        // Exclusive
		false,        // No wait
		nil,          // Arguments
	)
	failOnError(err, "Failed to connect to queue!")

	messages, err := channel.Consume(
		queue.Name, // Queue name
		"",         // Consumer
		false,      // Auto-ack - if set to 'true' the message will dissapear if the worker dies
		false,      // Exclusive
		false,      // No-local
		false,      // No-wait
		nil,        // Args
	)
	failOnError(err, "Failed to register a consumer!")

	forever := make(chan bool)

	go func() {
		for message := range messages {
			log.Printf("[O] Received message: %s", message.Body)

			dotCount := bytes.Count(message.Body, []byte("."))
			t := time.Duration(dotCount)
			time.Sleep(t * time.Second)

			log.Printf("DONE")
			message.Ack(false) // acknowledge the sender that the task has been done
		}
	}()

	log.Printf(" [*] Waiting for message...To exit, press Ctrl + C")
	<-forever
}
