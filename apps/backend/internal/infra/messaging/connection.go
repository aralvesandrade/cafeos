package messaging

import (
	"log/slog"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Connection struct {
	url  string
	conn *amqp.Connection
	ch   *amqp.Channel
}

func NewConnection(url string) (*Connection, error) {
	c := &Connection{url: url}
	if err := c.connect(); err != nil {
		return nil, err
	}
	return c, nil
}

func (c *Connection) connect() error {
	conn, err := amqp.Dial(c.url)
	if err != nil {
		return err
	}
	c.conn = conn

	ch, err := conn.Channel()
	if err != nil {
		conn.Close()
		return err
	}
	c.ch = ch

	go c.handleReconnect()
	return nil
}

func (c *Connection) handleReconnect() {
	notify := c.conn.NotifyClose(make(chan *amqp.Error))
	<-notify
	slog.Warn("rabbitmq connection lost, reconnecting...")
	for {
		time.Sleep(3 * time.Second)
		if err := c.connect(); err != nil {
			slog.Error("rabbitmq reconnect failed, retrying", "error", err)
			continue
		}
		slog.Info("rabbitmq reconnected")
		return
	}
}

func (c *Connection) Channel() *amqp.Channel {
	return c.ch
}

func (c *Connection) Close() {
	if c.ch != nil {
		c.ch.Close()
	}
	if c.conn != nil {
		c.conn.Close()
	}
}
