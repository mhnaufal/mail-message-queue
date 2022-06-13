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
	failOnError(err, "Failed to establish a connection to RabbitMQ!")
	defer conn.Close()

	channel, err := conn.Channel()
	failOnError(err, "Failed to open a channel!")
	defer channel.Close()

	err = channel.ExchangeDeclare(
		"logs",   // name
		"fanout", // type
		true,     // durable
		false,    // auto detect
		false,    // internal
		false,    // no-wait
		nil,      // argument
	)
	failOnError(err, "Failed to declare an exchange")

	body := bodyFrom(os.Args)
	err = channel.Publish(
		"logs", // exchange
		"",     // routing key
		false,  // mandatory
		false,  // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(body),
		})

	failOnError(err, "Failed to publish a message")

	log.Printf(" [x] Sent %s", body)
}

func bodyFrom(args []string) string {
	var s string

	if (len(args) < 2) || os.Args[1] == "" {
		s = `📧 Messages...message...message...`
	} else {
		s = strings.Join(args[1:], " ")
	}

	return s
}
