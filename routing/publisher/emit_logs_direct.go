package main

//we're going to make it possible to subscribe only to a subset of the messages. For example,
//we will be able to direct only critical error messages to the log file (to save disk space), while still being able to print all of the log messages on the console.

//Direct exchange
//Our logging system from the previous publish/subscribe broadcasts all messages to all consumers. We want to extend that to allow filtering messages based on their severity.
// For example we may want the script which is writing log messages to the disk to only receive critical errors, and not waste disk space on warning or info log messages.
//We were using a fanout exchange in publish/subscribe, which doesn't give us much flexibility - it's only capable of mindless broadcasting.
//We will use a direct exchange instead. The routing algorithm behind a direct exchange is simple - a message goes to the queues whose binding key exactly matches the routing key of the message.

import (
	"context"
	"log"
	"os"
	"strings"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", msg, err)
	}
}

func main() {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	failOnError(err, "failed to connect to rabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "failed to open a channel")
	defer ch.Close()

	err = ch.ExchangeDeclare(
		"logs_direct", //exchange namme
		"direct",      //exchange type
		true,          //durable
		false,         //auto-deleted
		false,         //internal
		false,         //no-wait
		nil,           //arguments
	)
	failOnError(err, "failed to declare exchange")
	body := bodyFrom(os.Args)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = ch.PublishWithContext(ctx,
		"logs_direct",         //exchange name
		severityFrom(os.Args), //routing key
		false,                 //mandatory
		false,                 //immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(body),
		})

	failOnError(err, "Failed to publish a message")

	log.Printf(" [x] Sent %s", body)
}

func bodyFrom(args []string) string {
	var s string
	if (len(args) < 3) || os.Args[2] == "" {
		s = "hello"
	} else {
		s = strings.Join(args[2:], " ")
	}
	return s
}

func severityFrom(args []string) string {
	var s string
	if (len(args) < 2) || os.Args[1] == "" {
		s = "info"
	} else {
		s = os.Args[1]
	}
	return s
}
