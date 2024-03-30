package main

import (
	"context"
	"log"
	"os"
	"strings"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

//routingKey is ignored for fanout

func main() {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		log.Panicf("%s: %s", "failed to connect to rabbitMQ", err)
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		log.Panicf("%s: %s", "Failed to open a channel", err)
	}
	defer ch.Close()

	err = ch.ExchangeDeclare(
		"logs",   //exchange name
		"fanout", //exchange type other different types are direct topic headers fanout
		true,     //durable
		false,    //auto-deleted
		false,    //internal
		false,    //no-wait
		nil,      //arguments
	)
	if err != nil {
		log.Panicf("%s: %s", "failed to declare an exchange", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	body := bodyFrom(os.Args)
	err = ch.PublishWithContext(ctx,
		"logs", //exchange we just declared above
		"",     //routing key we have to provide in exchange but here we are not providing because our exchange type is `fanout` and routingKey is ignored in fanout exchange
		false,  //mandotory
		false,  //immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(body),
		})
	if err != nil {
		log.Panicf("%s: %s", "failed to publish a message", err)
	}
	log.Printf("[x] Sent %s", body)

}
func bodyFrom(args []string) string {
	var s string
	if (len(args) < 2) || os.Args[1] == "" {
		s = "hello"
	} else {
		s = strings.Join(args[1:], " ")
	}
	return s
}
