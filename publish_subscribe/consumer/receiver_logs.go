package main

import (
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		log.Panicf("%s: %s", "failed to connect to rabbitMQ", err)
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		log.Panicf("%s: %s", "failed to open a channel", err)
	}
	defer ch.Close()
	//after establishing the connection we declared the exchange. This step is necessary as publishing to a non-existing exchange is forbidden.
	err = ch.ExchangeDeclare(
		"logs",   // name
		"fanout", // type
		true,     // durable
		false,    // auto-deleted
		false,    // internal
		false,    // no-wait
		nil,      // arguments
	)
	if err != nil {
		log.Panicf("%s: %s", "failed to declare a queue", err)
	}
	//The messages will be lost if no queue is bound to the exchange yet

	q, err := ch.QueueDeclare(
		"",    // name
		false, // durable
		false, // delete when unused
		true,  // exclusive
		false, // no-wait
		nil,   // arguments
	)
	if err != nil {
		log.Panicf("%s: %s", "failed to declare a queue", err)
	}

	//We've already created a fanout exchange and a queue. Now we need to tell the exchange to send messages to our queue. That relationship between exchange and a queue is called a binding.
	//From now on the logs exchange will append messages to our queue.
	err = ch.QueueBind(
		q.Name, //queue name
		"",     //routing key
		"logs", //exchange
		false,
		nil,
	)
	if err != nil {
		log.Panicf("%s: %s", "failed to bind queue", err)
	}

	msgs, err := ch.Consume(
		q.Name, //queue name binded with the logs exchange
		"",     //consumerr
		true,   //autto-ack
		false,  //exclusive
		false,  //no-local
		false,  //no-wait
		nil,    //arguments
	)
	if err != nil {
		log.Panicf("%s: %s", "failed to register a consumer", err)
	}
	var forever chan struct{}

	go func() {
		for msg := range msgs {
			log.Printf(" [x] %s", msg.Body)
		}
	}()
	log.Printf(" [*] Waiting for logs. To exit press CTRL+C")
	<-forever
}
