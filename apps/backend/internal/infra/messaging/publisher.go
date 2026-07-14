package messaging

import (
	"context"
	"encoding/json"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

type SyncMessage struct {
	ClientID        string      `json:"client_id"`
	EventType       string      `json:"event_type"`
	OrganizationID  string      `json:"organization_id"`
	Payload         interface{} `json:"payload"`
	ClientTimestamp string      `json:"client_timestamp"`
	PublishedAt     time.Time   `json:"published_at"`
}

type Publisher struct {
	ch     *amqp.Channel
	queues []string
}

func NewPublisher(ch *amqp.Channel) *Publisher {
	return &Publisher{ch: ch}
}

func (p *Publisher) DeclareQueue(name string) error {
	_, err := p.ch.QueueDeclare(
		name,  // name
		true,  // durable
		false, // autoDelete
		false, // exclusive
		false, // noWait
		nil,   // args
	)
	if err != nil {
		return err
	}
	return nil
}

func (p *Publisher) Publish(ctx context.Context, queue string, msg SyncMessage) error {
	body, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	return p.ch.PublishWithContext(ctx,
		"",    // exchange
		queue, // routing key
		true,  // mandatory
		false, // immediate
		amqp.Publishing{
			DeliveryMode: amqp.Persistent,
			ContentType:  "application/json",
			Body:         body,
			Timestamp:    time.Now(),
		},
	)
}
