package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Rabbitmq struct {
	hostname string
	user     string
	password string
	port     string
	queue    string
}

func (r *Rabbitmq) RabbitMqGetMessages(bb *Busyboi) {
	conn, err := amqp.Dial(fmt.Sprintf("amqp://%s:%s@%s:%s/", r.user, r.password, r.hostname, r.port))
	if err != nil {
		log.Panicf("Failed to connect to queue: %s", err)
	}

	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		log.Panicf("Failed to create channel: %s", err)
	}

	defer ch.Close()

	q, err := ch.QueueDeclare(
		r.queue, // name
		false,   // durable
		false,   // delete when unused
		false,   // exclusive
		false,   // no-wait
		nil,     // arguments
	)

	if err != nil {
		log.Panicf("Failed to declare queue: %s", err)
	}

	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)

	if err != nil {
		log.Panicf("Failed to consume: %s", err)
	}

	for {
		bb.queueMsgs <- <-msgs
	}
}

func (r *Rabbitmq) RabbitMqAddMessages(job JobConfig) {
	conn, err := amqp.Dial(fmt.Sprintf("amqp://%s:%s@%s:%s/", r.user, r.password, r.hostname, r.port))
	if err != nil {
		log.Panicf("Failed to connect to queue: %s", err)
	}

	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		log.Panicf("Failed to create channel: %s", err)
	}

	defer ch.Close()

	q, err := ch.QueueDeclare(
		r.queue, // name
		false,   // durable
		false,   // delete when unused
		false,   // exclusive
		false,   // no-wait
		nil,     // arguments
	)

	if err != nil {
		log.Panicf("Failed to declare queue: %s", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	body, _ := json.Marshal(job)

	err = ch.PublishWithContext(ctx,
		"",     // exchange
		q.Name, // routing key
		false,  // mandatory
		false,
		amqp.Publishing{
			DeliveryMode: amqp.Persistent,
			ContentType:  "text/plain",
			Body:         body,
		})

	if err != nil {
		log.Panicf("Failed to publish: %s", err)
	}
}
