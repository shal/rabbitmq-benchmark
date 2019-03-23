package main

import (
	"fmt"
	"log"
	"time"

	"github.com/streadway/amqp"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

func main() {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"benchmark",
		false,
		false,
		false,
		false,
		nil,
	)
	failOnError(err, "Failed to declare a queue")

	msgs, err := ch.Consume(
		q.Name,
		"benchmark",
		true,
		false,
		false,
		true,
		nil,
	)
	failOnError(err, "Failed to register a consumer")

	wait := make(chan struct{})

	start := time.Now()

	read := func() {
		for range msgs {
			wait <- struct{}{}
		}
	}

	for i := 0; i < 2; i++ {
		go read()
	}

	res := 0
	for i := 0; i < 100000; i++ {
		<-wait
		res++
		fmt.Printf("%f req/s\n", float64(res)/time.Since(start).Seconds())
	}
}
