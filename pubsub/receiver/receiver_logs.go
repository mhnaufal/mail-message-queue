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
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	failOnError(err, "Failed to establish a connection to RabbitMQ!")
	defer conn.Close()

	channel, err := conn.Channel()
	failOnError(err, "Failed to open a channel!")
	defer channel.Close()

	queue, err := channel.QueueDeclare(
		"",    // name - RabbitMQ will generate a random string for the name and it immediately deleted when the connection is closed
		false, // durable
		false, // delete when unused
		true,  // exclusive
		false, // no-wait
		nil,   // arguments
	)
	failOnError(err, "Failed to declare a queue")

	// Binding - relation between queue and exchange
	err = channel.QueueBind(
		queue.Name, // queue name
		"",         // routing key
		"logs",     // exchange
		false,
		nil,
	)
	failOnError(err, "Failed to register a binding")

	msgs, err := channel.Consume(
		queue.Name, // queue
		"",         // consumer
		true,       // auto-ack
		false,      // exclusive
		false,      // no-local
		false,      // no-wait
		nil,        // args
	)
	failOnError(err, "Failed to register a consumer")

	forever := make(chan bool)

	go func() {
		for d := range msgs {
			log.Printf(" [x] %s", d.Body)
		}
	}()

	log.Printf(" [*] Waiting for logs. To exit press CTRL+C")
	<-forever
}
