package main

import (
	"fmt"
	"log"
	"time"

	"github.com/streadway/amqp"
)

const (
	numOfMessages = 100000
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

	start := time.Now()

	for i := 0; i < numOfMessages; i++ {
		err = ch.Publish(
			"",
			q.Name,
			false,
			false,
			amqp.Publishing{
				ContentType: "text/plain",
				Body:        []byte("Benchmark"),
			},
		)

		failOnError(err, "Failed to publish a message")
	}

	fmt.Printf("%f req/s\n", float64(numOfMessages)/time.Since(start).Seconds())
}
