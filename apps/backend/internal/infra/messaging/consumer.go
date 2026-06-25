package messaging

import (
	"encoding/json"
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Handler func(msg SyncMessage) error

type Consumer struct {
	ch      *amqp.Channel
	queue   string
	handler Handler
}

func NewConsumer(ch *amqp.Channel, queue string, handler Handler) *Consumer {
	return &Consumer{ch: ch, queue: queue, handler: handler}
}

func (c *Consumer) Start() error {
	msgs, err := c.ch.Consume(
		c.queue,
		"",     // consumer
		false,  // autoAck (manual ack)
		false,  // exclusive
		false,  // noLocal
		false,  // noWait
		nil,
	)
	if err != nil {
		return err
	}

	go func() {
		for d := range msgs {
			var msg SyncMessage
			if err := json.Unmarshal(d.Body, &msg); err != nil {
				log.Printf("[CONSUMER] unmarshal error: %v", err)
				d.Nack(false, false) // discard
				continue
			}

			if err := c.handler(msg); err != nil {
				log.Printf("[CONSUMER] handler error: %v", err)
				d.Nack(false, !d.Redelivered)
				continue
			}

			d.Ack(false)
		}
	}()

	return nil
}
